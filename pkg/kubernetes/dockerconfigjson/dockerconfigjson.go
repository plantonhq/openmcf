// Package dockerconfigjson renders the `.dockerconfigjson` document a
// `kubernetes.io/dockerconfigjson` Secret carries: one `auths` record per registry
// server, each with the username, the password, and the base64 `auth` pair the
// kubelet actually reads.
//
// It is the single encoder for every module that materializes a registry login
// into a Secret — the KubernetesSecret kind's docker arm and the workload kinds'
// `pod.image_registries` — so the document's shape is defined once. The Terraform
// twins render the same shape with `jsonencode`; their comments name this package.
package dockerconfigjson

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sort"
)

// Auth is one registry login: the server it is for and how to authenticate.
type Auth struct {
	// Server is the registry host the kubelet matches an image reference against
	// ("ghcr.io", "123456789012.dkr.ecr.us-east-1.amazonaws.com"); Docker Hub's login
	// is keyed by "https://index.docker.io/v1/".
	Server   string
	Username string
	Password string
	// Email is optional; most registries ignore it and it is omitted when empty.
	Email string
}

type authRecord struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
	Auth     string `json:"auth"`
}

type document struct {
	Auths map[string]authRecord `json:"auths"`
}

// Encode renders the `.dockerconfigjson` document for the given logins.
//
// Output is deterministic for a given input (encoding/json sorts map keys), so
// the Secret's data is stable across runs and never shows a spurious diff. A
// server named twice is refused: two records cannot share one `auths` key, and
// silently keeping the last one would hand the kubelet a login the author did not
// intend. Spec validation already refuses the duplicate up front; this check keeps
// a hand-built input honest and repeats the same sentence.
func Encode(auths []Auth) (string, error) {
	doc := document{Auths: make(map[string]authRecord, len(auths))}
	for _, a := range auths {
		if _, seen := doc.Auths[a.Server]; seen {
			return "", fmt.Errorf("two image registry entries name the same registry server %q — a workload holds one login per registry; merge them into one entry", a.Server)
		}
		doc.Auths[a.Server] = authRecord{
			Username: a.Username,
			Password: a.Password,
			Email:    a.Email,
			Auth:     base64.StdEncoding.EncodeToString([]byte(a.Username + ":" + a.Password)),
		}
	}
	out, err := json.Marshal(doc)
	if err != nil {
		return "", fmt.Errorf("failed to marshal docker config json: %w", err)
	}
	return string(out), nil
}

// Servers lists the registry servers of the given logins in sorted order — the
// same order the encoded document carries them — for log lines and tests.
func Servers(auths []Auth) []string {
	servers := make([]string, 0, len(auths))
	for _, a := range auths {
		servers = append(servers, a.Server)
	}
	sort.Strings(servers)
	return servers
}
