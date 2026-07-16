package kubecontext

import (
	"github.com/plantonhq/planton/pkg/kubernetes/kubernetesannotationkeys"
	"github.com/plantonhq/planton/pkg/reflection/metadatareflect"
	"google.golang.org/protobuf/proto"
)

// ExtractFromManifest extracts the kubectl context from manifest annotations.
// Returns:
//   - The context name if the annotation exists
//   - Empty string if the annotation is not present (uses default context from kubeconfig)
func ExtractFromManifest(manifest proto.Message) string {
	annotations := metadatareflect.ExtractAnnotations(manifest)
	if annotations == nil {
		return ""
	}

	context, ok := annotations[kubernetesannotationkeys.KubeContextAnnotationKey]
	if !ok {
		return ""
	}

	return context
}
