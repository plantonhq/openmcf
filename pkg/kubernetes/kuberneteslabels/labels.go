package kuberneteslabels

const (
	// DockerConfigJsonFileLabelKey specifies the file path containing docker config JSON for image pull secret
	DockerConfigJsonFileLabelKey = "kubernetes.openmcf.org/docker-config-json-file"

	// KubeContextLabelKey specifies the kubectl context to use for Kubernetes deployments
	KubeContextLabelKey = "kubernetes.openmcf.org/context"
)
