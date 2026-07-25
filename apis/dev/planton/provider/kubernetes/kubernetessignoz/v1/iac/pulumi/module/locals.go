package module

import (
	"fmt"
	"strconv"

	kubernetessignozv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetessignoz/v1"
	"github.com/plantonhq/planton/apis/dev/planton/shared/cloudresourcekind"
	"github.com/plantonhq/planton/pkg/iac/pulumi/pulumimodule/provider/kubernetes/kuberneteslabelkeys"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Locals holds computed values derived from the stack input for use across
// the module. Every resolution here has an exact twin in the Terraform
// module's locals.tf — keep them in lockstep.
type Locals struct {
	Spec *kubernetessignozv1.KubernetesSignozSpec

	// Resource-identity labels stamped on the module-created satellites
	// (the namespace and the clickhouse-auth Secret — never injected into
	// the chart's own resources; Helm owns those).
	Labels map[string]string

	// Namespace SigNoz installs into (resolved literal from the spec's
	// value-or-ref).
	Namespace string

	// Helm release name — metadata.name, NOT a fixed chart name: several
	// SigNoz instances can coexist. fullnameOverride pins every chart
	// child name to it.
	ReleaseName string

	// Chart version resolved to the pinned default when unset, so both
	// engines install the same chart whether or not the platform's
	// defaulting middleware ran.
	ChartVersion string

	// Whether the external-clickhouse arm is declared (empty oneof = the
	// bundled arm with defaults — the appliance posture).
	IsExternal bool

	// The bundled installation's fullname (`<name>-clickhouse`) — pins
	// the CHI name AND the client Service name the chart's own helpers
	// derive.
	ClickHouseFullname string

	// The module-owned Secret exporting the generated bundled credential
	// (`<name>-clickhouse-auth`, keys username/password). Empty contract
	// on the external arm — the declared Secret is exported instead.
	ClickHouseAuthSecretName string

	// ---- exported composition handles (see outputs.go) -----------------
	Service                      string
	KubeEndpoint                 string
	PortForwardCommand           string
	OtelCollectorService         string
	OtlpGrpcEndpoint             string
	OtlpHttpEndpoint             string
	ClickHouseEndpoint           string
	ClickHouseUsername           string
	ClickHousePasswordSecretName string
	ClickHousePasswordSecretKey  string
}

// initializeLocals extracts and transforms spec fields into module-local
// values.
func initializeLocals(_ *pulumi.Context, stackInput *kubernetessignozv1.KubernetesSignozStackInput) *Locals {
	target := stackInput.Target
	spec := target.Spec

	labels := map[string]string{
		kuberneteslabelkeys.Resource:     strconv.FormatBool(true),
		kuberneteslabelkeys.ResourceName: target.Metadata.Name,
		kuberneteslabelkeys.ResourceKind: cloudresourcekind.CloudResourceKind_KubernetesSignoz.String(),
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

	namespace := spec.Namespace.GetValue()
	releaseName := target.Metadata.Name

	isExternal := spec.GetExternalClickhouse() != nil

	clickHouseFullname := fmt.Sprintf("%s-clickhouse", releaseName)
	clickHouseAuthSecretName := fmt.Sprintf("%s-clickhouse-auth", releaseName)

	// The signoz Service carries the fullname itself; the collector
	// Service appends the component name — both pinned via
	// fullnameOverride.
	service := releaseName
	kubeEndpoint := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", service, namespace, vars.ServerHttpPort)
	portForward := fmt.Sprintf("kubectl port-forward svc/%s -n %s %d:%d",
		service, namespace, vars.ServerHttpPort, vars.ServerHttpPort)

	collectorService := fmt.Sprintf("%s-otel-collector", releaseName)
	otlpGrpc := fmt.Sprintf("%s.%s.svc.cluster.local:%d", collectorService, namespace, vars.OtlpGrpcPort)
	otlpHttp := fmt.Sprintf("http://%s.%s.svc.cluster.local:%d", collectorService, namespace, vars.OtlpHttpPort)

	clickHouseEndpoint := fmt.Sprintf("%s.%s.svc.cluster.local:%d",
		clickHouseFullname, namespace, vars.ClickHouseTcpPort)
	clickHouseUsername := vars.BundledClickHouseUser
	passwordSecretName := clickHouseAuthSecretName
	passwordSecretKey := "password"
	if isExternal {
		external := spec.GetExternalClickhouse()
		tcpPort := int32(vars.ClickHouseTcpPort)
		if external.TcpPort != nil {
			tcpPort = external.GetTcpPort()
		}
		clickHouseEndpoint = fmt.Sprintf("%s:%d", external.Host.GetValue(), tcpPort)
		clickHouseUsername = external.Username
		passwordSecretName = external.PasswordSecret.SecretName.GetValue()
		passwordSecretKey = external.PasswordSecret.SecretKey
	}

	return &Locals{
		Spec:                         spec,
		Labels:                       labels,
		Namespace:                    namespace,
		ReleaseName:                  releaseName,
		ChartVersion:                 chartVersion,
		IsExternal:                   isExternal,
		ClickHouseFullname:           clickHouseFullname,
		ClickHouseAuthSecretName:     clickHouseAuthSecretName,
		Service:                      service,
		KubeEndpoint:                 kubeEndpoint,
		PortForwardCommand:           portForward,
		OtelCollectorService:         collectorService,
		OtlpGrpcEndpoint:             otlpGrpc,
		OtlpHttpEndpoint:             otlpHttp,
		ClickHouseEndpoint:           clickHouseEndpoint,
		ClickHouseUsername:           clickHouseUsername,
		ClickHousePasswordSecretName: passwordSecretName,
		ClickHousePasswordSecretKey:  passwordSecretKey,
	}
}
