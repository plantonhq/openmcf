package module

import (
	"fmt"
	"strconv"

	kubernetesopenbaov1alpha1 "github.com/plantonhq/planton/catalog/kubernetes/kubernetesopenbao/v1alpha1"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/kuberneteslabelkeys"
	"github.com/plantonhq/planton/shared/cloudresourcekind"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Server modes (the spec's mode oneof; unset = standalone — the chart's
// own default).
const (
	modeDev        = "dev"
	modeStandalone = "standalone"
	modeHa         = "ha"
)

// Locals holds computed values derived from the stack input. Every
// resolution here has an exact twin in the Terraform module's locals.tf —
// keep them in lockstep.
type Locals struct {
	Spec *kubernetesopenbaov1alpha1.KubernetesOpenBaoSpec

	// Resource-identity labels stamped on module-created satellites
	// (namespace, seal-credentials Secret) — never injected into the
	// chart's own resources; Helm owns those.
	Labels map[string]string

	// Namespace the server installs into (resolved literal).
	Namespace string

	// ReleaseName is metadata.name; fullnameOverride is pinned to it,
	// so every chart-derived name hangs off this value.
	ReleaseName string

	ChartVersion string

	// Mode resolved from the spec oneof (dev / standalone / ha).
	Mode string

	// Raft peer count (ha mode; 1 otherwise).
	Replicas int

	// http or https, following tls.enabled — drives the synthesized
	// listener config, retry_join addresses, and the exported endpoint.
	Scheme string

	// True when a TLS certificate Secret is mounted.
	TlsEnabled bool

	// Resolved TLS Secret name ("" when TLS is off).
	TlsSecretName string

	// Name of the module-owned seal-credentials Secret ("" when the
	// declared seal arm carries no credential material — the keyless
	// posture).
	SealCredentialsSecretName string

	// The synthesized server config HCL (empty in dev mode — dev
	// ignores config).
	BaoConfigHcl string
}

// initializeLocals extracts and transforms spec fields into module-local
// values.
func initializeLocals(_ *pulumi.Context, stackInput *kubernetesopenbaov1alpha1.KubernetesOpenBaoStackInput) *Locals {
	target := stackInput.Target
	spec := target.Spec

	labels := map[string]string{
		kuberneteslabelkeys.Resource:     strconv.FormatBool(true),
		kuberneteslabelkeys.ResourceName: target.Metadata.Name,
		kuberneteslabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_KubernetesOpenBao.String(),
	}
	if target.Metadata.Id != "" {
		labels[kuberneteslabelkeys.ResourceId] = target.Metadata.Id
	}
	if target.Metadata.Org != "" {
		labels[kuberneteslabelkeys.Organization] = target.Metadata.Org
	}
	if target.Metadata.Env != "" {
		labels[kuberneteslabelkeys.Environment] = target.Metadata.Env
	}

	chartVersion := spec.GetChartVersion()
	if chartVersion == "" {
		chartVersion = vars.DefaultChartVersion
	}

	// Resolve the mode oneof; unset = standalone (the chart default).
	mode := modeStandalone
	replicas := 1
	if spec.GetServer() != nil {
		switch {
		case spec.GetServer().GetDev() != nil:
			mode = modeDev
		case spec.GetServer().GetHa() != nil:
			mode = modeHa
			replicas = 3
			if spec.GetServer().GetHa().Replicas != nil {
				replicas = int(spec.GetServer().GetHa().GetReplicas())
			}
		}
	}

	tlsEnabled := spec.GetTls().GetEnabled()
	tlsSecretName := ""
	if tlsEnabled {
		tlsSecretName = spec.GetTls().GetCertSecretName().GetValue()
	}
	scheme := "http"
	if tlsEnabled {
		scheme = "https"
	}

	sealCredentialsSecretName := ""
	if sealSecretData(spec) != nil {
		sealCredentialsSecretName = target.Metadata.Name + vars.SealCredentialsSecretSuffix
	}

	locals := &Locals{
		Spec:                      spec,
		Labels:                    labels,
		Namespace:                 spec.Namespace.GetValue(),
		ReleaseName:               target.Metadata.Name,
		ChartVersion:              chartVersion,
		Mode:                      mode,
		Replicas:                  replicas,
		Scheme:                    scheme,
		TlsEnabled:                tlsEnabled,
		TlsSecretName:             tlsSecretName,
		SealCredentialsSecretName: sealCredentialsSecretName,
	}

	locals.BaoConfigHcl = renderBaoConfigHcl(locals)

	return locals
}

// apiEndpoint is the in-cluster URL clients (external-secrets stores,
// cert-manager Vault issuers) point at.
func (l *Locals) apiEndpoint() string {
	return fmt.Sprintf("%s://%s.%s.svc.cluster.local:%d", l.Scheme, l.ReleaseName, l.Namespace, vars.ApiPort)
}
