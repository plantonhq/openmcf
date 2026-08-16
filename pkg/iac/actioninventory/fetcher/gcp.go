package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/plantonhq/planton/pkg/iac/actioninventory"
	"github.com/plantonhq/planton/pkg/iac/permissions"
)

// queryTestablePermissionsURL is IAM's inventory of every permission name
// GCP defines -- the same list `gcloud iam list-testable-permissions`
// reads. The query anchors on a resource, but permission NAMES are
// global: permissions of services not enabled on the anchor project are
// still listed (flagged apiDisabled), so any project answers the
// existence question.
const queryTestablePermissionsURL = "https://iam.googleapis.com/v1/permissions:queryTestablePermissions"

// refreshGcp rewrites the committed GCP snapshot
// (pkg/iac/actioninventory/gcp.yaml) from IAM's testable-permissions
// inventory, scoped to exactly the services (the first dotted segment,
// e.g. "container" in "container.clusters.create") the committed
// permissions manifests reference. The credential contract mirrors the
// GCP price fetcher's ADC arm: the operator's own application-default
// credentials mint the bearer token (network in make, never in CI). The
// query anchors on the operator's active gcloud project.
func refreshGcp(repoRoot string) error {
	services, err := referencedGcpServices(repoRoot)
	if err != nil {
		return err
	}
	if len(services) == 0 {
		fmt.Println("no GCP permissions in any permissions manifest -- gcp snapshot not written")
		return nil
	}

	token, err := gcpAccessToken()
	if err != nil {
		return err
	}
	project, err := gcloudProject()
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 60 * time.Second}
	permissionsByService, err := fetchTestablePermissions(client, token, project)
	if err != nil {
		return err
	}

	today := time.Now().UTC().Format("2006-01-02")
	inv := &actioninventory.Inventory{Provider: "gcp"}
	for _, service := range services {
		names := permissionsByService[service]
		if len(names) == 0 {
			return fmt.Errorf("service %q (used by a permissions manifest) defines no permissions in GCP's IAM inventory -- the service segment itself is wrong", service)
		}
		sort.Strings(names)
		inv.Services = append(inv.Services, actioninventory.Service{
			Prefix:      service,
			SourceURL:   queryTestablePermissionsURL,
			RetrievedOn: today,
			Actions:     names,
		})
	}

	out := filepath.Join(repoRoot, "pkg", "iac", "actioninventory", actioninventory.GcpFileName)
	if err := os.WriteFile(out, []byte(actioninventory.Render(inv)), 0o644); err != nil {
		return err
	}
	fmt.Printf("wrote %d service permission list(s) to %s\n", len(inv.Services), out)
	return nil
}

// referencedGcpServices collects the distinct GCP service segments named
// by every committed permissions manifest, sorted.
func referencedGcpServices(repoRoot string) ([]string, error) {
	discovered, err := permissions.Discover(repoRoot)
	if err != nil {
		return nil, err
	}
	set := map[string]bool{}
	for provider, components := range discovered {
		for _, component := range components {
			manifest, err := permissions.Load(repoRoot, provider, component)
			if err != nil {
				return nil, err
			}
			for _, group := range manifest.GetSpec().GetGcp().GetGroups() {
				for _, permission := range group.GetPermissions() {
					service, _, found := strings.Cut(permission, ".")
					if !found {
						return nil, fmt.Errorf("%s/%s: gcp permission %q has no service segment", provider, component, permission)
					}
					set[service] = true
				}
			}
		}
	}
	services := make([]string, 0, len(set))
	for service := range set {
		services = append(services, service)
	}
	sort.Strings(services)
	return services, nil
}

// gcpAccessToken mints a bearer token from the operator's own
// application-default credentials -- the fetcher's credential contract,
// shared with the GCP price fetcher. Missing credentials are a hard error
// naming the way in.
func gcpAccessToken() (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("gcloud", "auth", "application-default", "print-access-token")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("the IAM testable-permissions API needs a credential: run `gcloud auth application-default login` so the fetcher can mint an access token: %v: %s",
			err, strings.TrimSpace(stderr.String()))
	}
	token := strings.TrimSpace(stdout.String())
	if token == "" {
		return "", fmt.Errorf("`gcloud auth application-default print-access-token` returned an empty token")
	}
	return token, nil
}

// gcloudProject reads the operator's active gcloud project -- the query's
// resource anchor (any project answers; see queryTestablePermissionsURL).
func gcloudProject() (string, error) {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command("gcloud", "config", "get-value", "project")
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("reading the active gcloud project (`gcloud config get-value project`): %v: %s",
			err, strings.TrimSpace(stderr.String()))
	}
	project := strings.TrimSpace(stdout.String())
	if project == "" || project == "(unset)" {
		return "", fmt.Errorf("no active gcloud project -- run `gcloud config set project <project-id>` (any project; it only anchors the query)")
	}
	return project, nil
}

// fetchTestablePermissions pages through IAM's full permission inventory
// for the anchor project and returns the names grouped by service segment
// (the part before the first dot), stripped of that segment.
func fetchTestablePermissions(client *http.Client, token, project string) (map[string][]string, error) {
	permissionsByService := map[string][]string{}
	seen := map[string]bool{}
	pageToken := ""
	for {
		body, err := json.Marshal(map[string]any{
			"fullResourceName": "//cloudresourcemanager.googleapis.com/projects/" + project,
			"pageSize":         1000,
			"pageToken":        pageToken,
		})
		if err != nil {
			return nil, err
		}
		request, err := http.NewRequest(http.MethodPost, queryTestablePermissionsURL, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		request.Header.Set("Authorization", "Bearer "+token)
		request.Header.Set("Content-Type", "application/json")
		response, err := client.Do(request)
		if err != nil {
			return nil, err
		}
		var doc struct {
			Permissions []struct {
				Name string `json:"name"`
			} `json:"permissions"`
			NextPageToken string `json:"nextPageToken"`
		}
		decodeErr := json.NewDecoder(response.Body).Decode(&doc)
		response.Body.Close()
		if response.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("POST %s: %s", queryTestablePermissionsURL, response.Status)
		}
		if decodeErr != nil {
			return nil, fmt.Errorf("POST %s: decoding: %w", queryTestablePermissionsURL, decodeErr)
		}
		for _, permission := range doc.Permissions {
			service, name, found := strings.Cut(permission.Name, ".")
			if !found || name == "" || seen[permission.Name] {
				continue
			}
			seen[permission.Name] = true
			permissionsByService[service] = append(permissionsByService[service], name)
		}
		if doc.NextPageToken == "" {
			return permissionsByService, nil
		}
		pageToken = doc.NextPageToken
	}
}
