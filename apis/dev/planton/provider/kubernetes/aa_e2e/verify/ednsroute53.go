package verify

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// ExternalDnsRoute53Verifier proves REAL record writes against Route 53: on
// top of the install assertions, it creates a verifier-owned Service source
// (hostname + explicit target annotations — the target annotation publishes
// endpoints regardless of service type, verified in the pinned controller
// source), asserts the A record AND the TXT ownership record materialize in
// the batch's hosted zone via the AWS API, then deletes the source and
// asserts the records are gone again (the sync policy's full loop). The
// zone needs no public delegation: external-dns writes through the API, and
// the API is where the proof reads.
//
// Batch coordinates come from the bootstrap-exported environment
// (PLANTON_E2E_ROUTE53_ZONE_ID / _ZONE_NAME); AWS CLI calls inherit the
// batch profile from the process environment.
type ExternalDnsRoute53Verifier struct {
	Namespace     string
	ComponentName string
}

func (v *ExternalDnsRoute53Verifier) VerifyExists(ctx context.Context, kubeconfig string) error {
	// Install contract first — identical to the plain install verifier.
	install := &ExternalDnsInstallVerifier{Namespace: v.Namespace, ComponentName: v.ComponentName}
	if err := install.VerifyExists(ctx, kubeconfig); err != nil {
		return err
	}

	zoneID := os.Getenv("PLANTON_E2E_ROUTE53_ZONE_ID")
	zoneName := os.Getenv("PLANTON_E2E_ROUTE53_ZONE_NAME")
	if zoneID == "" || zoneName == "" {
		return errors.New("PLANTON_E2E_ROUTE53_ZONE_ID / _ZONE_NAME unset — the batch bootstrap exports them")
	}
	recordName := "e2e-proof." + zoneName
	fmt.Printf("  [verify] behavioral record writes: %q must appear in zone %s\n", recordName, zoneID)

	// The source: hostname says WHAT to publish, target says WHERE it
	// points (TEST-NET-3 documentation address — nothing routes there, and
	// nothing needs to: the API read is the assertion).
	source := fmt.Sprintf(`apiVersion: v1
kind: Service
metadata:
  name: e2e-edns-source
  namespace: %s
  annotations:
    external-dns.alpha.kubernetes.io/hostname: %s
    external-dns.alpha.kubernetes.io/target: 203.0.113.42
spec:
  type: ClusterIP
  ports:
    - port: 80
`, v.Namespace, recordName)
	sourceFile, err := writeTempManifest(source)
	if err != nil {
		return err
	}
	defer os.Remove(sourceFile)
	if err := v.kubectl(ctx, kubeconfig, "apply", "-f", sourceFile); err != nil {
		return errors.Wrap(err, "failed to apply source service")
	}
	defer func() {
		_ = v.kubectl(context.Background(), kubeconfig, "delete", "service", "e2e-edns-source",
			"-n", v.Namespace, "--ignore-not-found")
	}()

	// Write half: the A record and its TXT ownership sibling (the registry
	// that makes the record deletable by exactly this instance).
	if err := v.waitForRecords(ctx, zoneID, recordName, true, 4*time.Minute); err != nil {
		return errors.Wrap(err, "external-dns never wrote the records to Route 53")
	}
	fmt.Printf("  [verify] A + TXT ownership records present in Route 53 — a real provider write\n")

	// Delete half: sync policy must remove what it owns once the source is
	// gone.
	if err := v.kubectl(ctx, kubeconfig, "delete", "service", "e2e-edns-source",
		"-n", v.Namespace, "--ignore-not-found"); err != nil {
		return err
	}
	if err := v.waitForRecords(ctx, zoneID, recordName, false, 4*time.Minute); err != nil {
		return errors.Wrap(err, "external-dns never removed the records after the source was deleted (sync policy)")
	}
	fmt.Printf("  [verify] records removed after source deletion — the full write-and-own loop is proven\n")
	return nil
}

func (v *ExternalDnsRoute53Verifier) VerifyAbsent(ctx context.Context, kubeconfig string) error {
	install := &ExternalDnsInstallVerifier{Namespace: v.Namespace, ComponentName: v.ComponentName}
	return install.VerifyAbsent(ctx, kubeconfig)
}

// waitForRecords polls the hosted zone via the AWS CLI until records for
// recordName are present (want=true) or gone (want=false). Both the A
// record and the TXT ownership record must agree — half-written state is a
// failure in either direction.
func (v *ExternalDnsRoute53Verifier) waitForRecords(ctx context.Context, zoneID, recordName string, want bool, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	var last string
	for time.Now().Before(deadline) {
		out, err := exec.CommandContext(ctx, "aws", "route53", "list-resource-record-sets",
			"--hosted-zone-id", zoneID, "--output", "json").CombinedOutput()
		if err != nil {
			last = fmt.Sprintf("aws err=%v: %s", err, strings.TrimSpace(string(out)))
		} else {
			var payload struct {
				ResourceRecordSets []struct {
					Name string `json:"Name"`
					Type string `json:"Type"`
				} `json:"ResourceRecordSets"`
			}
			if jsonErr := json.Unmarshal(out, &payload); jsonErr != nil {
				last = fmt.Sprintf("json err=%v", jsonErr)
			} else {
				hasA, hasTxt := false, false
				for _, rrs := range payload.ResourceRecordSets {
					if strings.TrimSuffix(rrs.Name, ".") == recordName {
						switch rrs.Type {
						case "A":
							hasA = true
						case "TXT":
							hasTxt = true
						}
					}
					// external-dns also writes prefixed TXT registry forms
					// (a-<name>, cname-<name>); ownership counts either way.
					if strings.Contains(rrs.Name, recordName) && rrs.Type == "TXT" {
						hasTxt = true
					}
				}
				if want && hasA && hasTxt {
					return nil
				}
				if !want && !hasA && !hasTxt {
					return nil
				}
				last = fmt.Sprintf("hasA=%v hasTxt=%v want=%v", hasA, hasTxt, want)
			}
		}
		time.Sleep(10 * time.Second)
	}
	return errors.Errorf("record state never converged (last: %s)", last)
}

func (v *ExternalDnsRoute53Verifier) kubectl(ctx context.Context, kubeconfig string, args ...string) error {
	full := append([]string{"--kubeconfig", kubeconfig}, args...)
	if out, err := exec.CommandContext(ctx, "kubectl", full...).CombinedOutput(); err != nil {
		return errors.Errorf("kubectl %s: %v: %s", strings.Join(args, " "), err, string(out))
	}
	return nil
}
