// Package kubernetesannotationkeys defines the manifest annotation keys that carry
// Kubernetes-specific platform behavior.
//
// Platform-behavior signals live in metadata.annotations, never metadata.labels:
// labels are derived into cloud-provider tags by planton IaC modules, so a platform
// key there would leak internal detail onto the user's real cloud resources.
package kubernetesannotationkeys

const (
	// DockerConfigJsonFileAnnotationKey names a docker-config JSON file on the
	// machine that RUNS the module, read at apply time to build an image pull
	// secret. This is a break-glass path for an expert applying a module from
	// their own laptop: the file lives on the operator's disk, never inside the
	// module, so it is deliberately outside the module self-containment
	// invariant (a module reads nothing beside its own directory). It has no
	// meaning in the runner or the Pulumi binary lane, where the platform
	// resolves pull secrets through its own seams.
	DockerConfigJsonFileAnnotationKey = "kubernetes.planton.dev/docker-config-json-file"

	// KubeContextAnnotationKey specifies the kubectl context to use for Kubernetes deployments
	KubeContextAnnotationKey = "kubernetes.planton.dev/context"
)
