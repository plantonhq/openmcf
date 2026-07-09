// Package kubernetesannotationkeys defines the manifest annotation keys that carry
// Kubernetes-specific platform behavior.
//
// Platform-behavior signals live in metadata.annotations, never metadata.labels:
// labels are derived into cloud-provider tags by planton IaC modules, so a platform
// key there would leak internal detail onto the user's real cloud resources.
package kubernetesannotationkeys

const (
	// DockerConfigJsonFileAnnotationKey specifies the file path containing docker config JSON for image pull secret
	DockerConfigJsonFileAnnotationKey = "kubernetes.planton.dev/docker-config-json-file"

	// KubeContextAnnotationKey specifies the kubectl context to use for Kubernetes deployments
	KubeContextAnnotationKey = "kubernetes.planton.dev/context"
)
