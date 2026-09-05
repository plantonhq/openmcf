package workloadpod

import (
	"github.com/pkg/errors"
	kubernetesv1 "github.com/plantonhq/planton/catalog/kubernetes"
	"github.com/plantonhq/planton/pkg/kubernetes/dockerconfigjson"
)

// ImagePullSecretDataKey is the one key a `kubernetes.io/dockerconfigjson` Secret
// carries.
const ImagePullSecretDataKey = ".dockerconfigjson"

// BuildImagePullSecretData turns the pod's declared `image_registries` into the
// data of the ONE workload-scoped image-pull Secret — the twin of
// CollectLiteralEnvSecrets for the `<workload>-env-secrets` Secret. It returns nil
// when the pod declares no registry, so callers create the Secret only when the
// spec asked for one (a same-cloud registry the cluster's own identity reaches is
// never listed, and gets no Secret).
//
// The password reaches this code already resolved: the orchestrator replaces the
// manifest's `$secret/<slug>` reference with the value inside the customer's
// infrastructure before the module runs, so no plaintext ever sits in a record.
func BuildImagePullSecretData(pod *kubernetesv1.WorkloadPod) (map[string]string, error) {
	if pod == nil || len(pod.ImageRegistries) == 0 {
		return nil, nil
	}
	auths := make([]dockerconfigjson.Auth, 0, len(pod.ImageRegistries))
	for _, r := range pod.ImageRegistries {
		if r == nil {
			continue
		}
		auths = append(auths, dockerconfigjson.Auth{
			Server:   r.GetServer(),
			Username: r.GetUsername(),
			Password: r.GetPassword(),
			Email:    r.GetEmail(),
		})
	}
	doc, err := dockerconfigjson.Encode(auths)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build the image-pull secret from pod.image_registries")
	}
	return map[string]string{ImagePullSecretDataKey: doc}, nil
}
