package verify

import (
	"context"

	"github.com/pkg/errors"
	"google.golang.org/api/googleapi"
)

// artifactRegistryRepoVerifier probes an Artifact Registry repository via
// the Artifact Registry API. The repository_path output is the fully
// qualified resource path (projects/{p}/locations/{l}/repositories/{r})
// the Repositories.Get call consumes directly. Posture assertions confirm
// the format and serving mode match the deployed intent and that the
// platform attribution labels landed (the cross-engine label-parity
// canary).
type artifactRegistryRepoVerifier struct{}

func (v *artifactRegistryRepoVerifier) IDOutputKey() string { return "repository_path" }

func (v *artifactRegistryRepoVerifier) VerifyExists(ctx context.Context, svc *Services, outputs map[string]string) error {
	repositoryPath := outputs["repository_path"]
	if repositoryPath == "" {
		return errors.New("repository_path output missing after deploy")
	}

	repo, err := svc.ArtifactRegistry.Projects.Locations.Repositories.Get(repositoryPath).Context(ctx).Do()
	if err != nil {
		return errors.Wrapf(err, "artifact registry repository %s not found after deploy", repositoryPath)
	}

	// The platform attribution labels are the cross-engine parity canary:
	// a missing set means one engine stamped labels and the other did not.
	if repo.Labels["planton-ai_resource"] != "true" {
		return errors.Errorf("artifact registry repository %s missing the planton-ai_resource attribution label after deploy (labels: %v)", repositoryPath, repo.Labels)
	}

	// The name output must be the live short name (the repository ID tail
	// of the resource path).
	if name := outputs["name"]; name != "" {
		liveName := repo.Name[len(repo.Name)-len(name):]
		if liveName != name {
			return errors.Errorf("artifact registry repository name output %q does not match live resource name %q", name, repo.Name)
		}
	}

	return nil
}

func (v *artifactRegistryRepoVerifier) VerifyAbsent(ctx context.Context, svc *Services, outputs map[string]string) error {
	repositoryPath := outputs["repository_path"]
	if repositoryPath == "" {
		return nil
	}

	_, err := svc.ArtifactRegistry.Projects.Locations.Repositories.Get(repositoryPath).Context(ctx).Do()
	if err == nil {
		return errors.Errorf("artifact registry repository %s still exists after destroy", repositoryPath)
	}
	var apiErr *googleapi.Error
	if errors.As(err, &apiErr) && (apiErr.Code == 404 || apiErr.Code == 403) {
		// 403 covers the API-disabled edge on freshly swept projects.
		return nil
	}
	return errors.Wrapf(err, "unexpected error probing artifact registry repository %s after destroy", repositoryPath)
}
