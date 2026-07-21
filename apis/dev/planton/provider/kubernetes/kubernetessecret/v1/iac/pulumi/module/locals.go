package module

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	kubernetessecretv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetessecret/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds all derived configuration and state for the module
type Locals struct {
	// Context for Pulumi operations
	Ctx *pulumi.Context

	// Stack input containing the target resource
	StackInput *kubernetessecretv1.KubernetesSecretStackInput

	// Target secret resource
	Target *kubernetessecretv1.KubernetesSecret

	// Spec from the target
	Spec *kubernetessecretv1.KubernetesSecretSpec

	// Secret name
	SecretName string

	// Secret namespace
	SecretNamespace string

	// Combined labels (spec labels + standard labels)
	Labels map[string]string

	// Combined annotations (spec annotations)
	Annotations map[string]string

	// Whether the secret is immutable
	Immutable bool

	// Kubernetes secret type string (e.g., "Opaque", "kubernetes.io/tls")
	SecretType string

	// Secret data map (stringData keys to plain-string values)
	SecretData map[string]string

	// Secret binary data map (data keys to base64-encoded values)
	SecretBinaryData map[string]string
}

// initializeLocals creates and populates the Locals struct
func initializeLocals(ctx *pulumi.Context, stackInput *kubernetessecretv1.KubernetesSecretStackInput) (*Locals, error) {
	locals := &Locals{
		Ctx:        ctx,
		StackInput: stackInput,
		Target:     stackInput.Target,
		Spec:       stackInput.Target.Spec,
	}

	locals.SecretName = stackInput.Target.Spec.Name

	// Namespace is a value-or-ref; by module execution time any reference has been
	// resolved into the literal value. Fall back to "default" when unset, mirroring
	// kubectl behavior without a namespace flag.
	locals.SecretNamespace = stackInput.Target.Spec.Namespace.GetValue()
	if locals.SecretNamespace == "" {
		locals.SecretNamespace = "default"
	}

	locals.Immutable = stackInput.Target.Spec.Immutable

	// Build labels
	locals.Labels = buildLabels(locals)

	// Build annotations
	locals.Annotations = buildAnnotations(locals)

	// Compute secret type and data from the oneof variant
	variant, err := computeSecretTypeAndData(locals.Spec)
	if err != nil {
		return nil, fmt.Errorf("failed to compute secret type and data: %w", err)
	}
	locals.SecretType = variant.secretType
	locals.SecretData = variant.stringData
	locals.SecretBinaryData = variant.binaryData

	// Some secret types mandate annotations (service-account-token binds to its
	// ServiceAccount via annotation); merge them into the user-provided ones.
	for k, v := range variant.annotations {
		locals.Annotations[k] = v
	}

	return locals, nil
}

// buildLabels combines spec labels with standard labels
func buildLabels(locals *Locals) map[string]string {
	labels := make(map[string]string)

	// Add standard labels
	labels["managed-by"] = "planton"
	labels["resource"] = locals.Target.Metadata.Name
	labels["resource-kind"] = "KubernetesSecret"

	// Add spec labels
	for k, v := range locals.Spec.Labels {
		labels[k] = v
	}

	return labels
}

// buildAnnotations combines spec annotations
func buildAnnotations(locals *Locals) map[string]string {
	annotations := make(map[string]string)

	// Add spec annotations
	for k, v := range locals.Spec.Annotations {
		annotations[k] = v
	}

	return annotations
}

// secretVariant holds everything derived from the oneof secret_data selection:
// the Kubernetes secret type plus the data and annotations that type mandates.
type secretVariant struct {
	secretType string

	// stringData carries plain-string values (Kubernetes stringData semantics)
	stringData map[string]string

	// binaryData carries base64-encoded values (Kubernetes data semantics)
	binaryData map[string]string

	// annotations mandated by the secret type (e.g., the service-account binding)
	annotations map[string]string
}

// computeSecretTypeAndData determines the Kubernetes secret type and constructs
// the data and annotation maps based on which oneof variant is set in the spec.
func computeSecretTypeAndData(spec *kubernetessecretv1.KubernetesSecretSpec) (*secretVariant, error) {
	switch data := spec.SecretData.(type) {
	case *kubernetessecretv1.KubernetesSecretSpec_Opaque:
		return &secretVariant{
			secretType: "Opaque",
			stringData: data.Opaque.Data,
			binaryData: data.Opaque.BinaryData,
		}, nil

	case *kubernetessecretv1.KubernetesSecretSpec_Tls:
		return &secretVariant{
			secretType: "kubernetes.io/tls",
			stringData: map[string]string{
				"tls.crt": data.Tls.TlsCrt,
				"tls.key": data.Tls.TlsKey,
			},
		}, nil

	case *kubernetessecretv1.KubernetesSecretSpec_DockerConfigJson:
		dockerConfigJSON, err := buildDockerConfigJSON(data.DockerConfigJson)
		if err != nil {
			return nil, fmt.Errorf("failed to build docker config json: %w", err)
		}
		return &secretVariant{
			secretType: "kubernetes.io/dockerconfigjson",
			stringData: map[string]string{
				".dockerconfigjson": dockerConfigJSON,
			},
		}, nil

	case *kubernetessecretv1.KubernetesSecretSpec_BasicAuth:
		return &secretVariant{
			secretType: "kubernetes.io/basic-auth",
			stringData: map[string]string{
				"username": data.BasicAuth.Username,
				"password": data.BasicAuth.Password,
			},
		}, nil

	case *kubernetessecretv1.KubernetesSecretSpec_SshAuth:
		return &secretVariant{
			secretType: "kubernetes.io/ssh-auth",
			stringData: map[string]string{
				"ssh-privatekey": data.SshAuth.SshPrivateKey,
			},
		}, nil

	case *kubernetessecretv1.KubernetesSecretSpec_ServiceAccountToken:
		// No data is set: the cluster's token controller populates token/ca.crt/namespace
		// after creation. The annotation binds the token to its ServiceAccount.
		return &secretVariant{
			secretType: "kubernetes.io/service-account-token",
			annotations: map[string]string{
				"kubernetes.io/service-account.name": data.ServiceAccountToken.ServiceAccountName.GetValue(),
			},
		}, nil

	default:
		return nil, fmt.Errorf("no secret data variant set in spec")
	}
}

// dockerConfigAuth represents the auth entry for a single registry
type dockerConfigAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
	Auth     string `json:"auth"`
}

// dockerConfigJSON represents the top-level .dockerconfigjson structure
type dockerConfigJSON struct {
	Auths map[string]dockerConfigAuth `json:"auths"`
}

// buildDockerConfigJSON constructs the .dockerconfigjson JSON string from structured fields
func buildDockerConfigJSON(data *kubernetessecretv1.KubernetesSecretDockerConfigJsonData) (string, error) {
	auth := base64.StdEncoding.EncodeToString(
		[]byte(data.Username + ":" + data.Password),
	)

	config := dockerConfigJSON{
		Auths: map[string]dockerConfigAuth{
			data.RegistryServer: {
				Username: data.Username,
				Password: data.Password,
				Email:    data.Email,
				Auth:     auth,
			},
		},
	}

	jsonBytes, err := json.Marshal(config)
	if err != nil {
		return "", fmt.Errorf("failed to marshal docker config json: %w", err)
	}

	return string(jsonBytes), nil
}
