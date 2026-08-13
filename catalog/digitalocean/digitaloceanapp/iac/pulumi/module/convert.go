package module

import (
	"strings"

	"github.com/pkg/errors"
	do "github.com/plantonhq/planton/catalog/digitalocean"
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// providerEnum maps a proto enum identifier (run_and_build_time, docr,
// pre_deploy) to the provider's uppercase token (RUN_AND_BUILD_TIME, DOCR,
// PRE_DEPLOY). Unspecified values become empty so the caller can omit the
// argument and let the provider default apply.
func providerEnum(enumName string) string {
	if enumName == "" || strings.HasSuffix(enumName, "_unspecified") {
		return ""
	}
	return strings.ToUpper(enumName)
}

func strPtr(s string) pulumi.StringPtrInput {
	if s == "" {
		return nil
	}
	return pulumi.StringPtr(s)
}

func intPtrFromUint32(p *uint32) pulumi.IntPtrInput {
	if p == nil {
		return nil
	}
	v := int(*p)
	return pulumi.IntPtr(v)
}

func destinationsSet(d *do.DigitalOceanAppAlertDestinations) bool {
	if d == nil {
		return false
	}
	return len(d.GetEmails()) > 0 || len(d.GetSlackWebhooks()) > 0
}

const destGap = "PARITY-EXCEPTION: alert destinations (emails / slack webhooks) are modeled and Terraform wires them; the Pulumi DigitalOcean SDK v4.49.0 has no destinations field on app or component alerts. Re-evaluate when the SDK exposes alert destinations."

func envTriple(e *do.DigitalOceanAppEnvVar) (key, value, typ, scope string) {
	key = e.GetKey()
	if e.GetSecret() != "" {
		value = e.GetSecret()
		typ = "SECRET"
	} else {
		value = e.GetPlaintext()
		typ = "GENERAL"
	}
	scope = providerEnum(e.GetScope().String())
	if scope == "" {
		scope = "RUN_AND_BUILD_TIME"
	}
	return
}

func appEnvs(envs []*do.DigitalOceanAppEnvVar) digitalocean.AppSpecEnvArray {
	out := digitalocean.AppSpecEnvArray{}
	for _, e := range envs {
		key, value, typ, scope := envTriple(e)
		out = append(out, digitalocean.AppSpecEnvArgs{
			Key:   pulumi.String(key),
			Value: pulumi.String(value),
			Type:  pulumi.String(typ),
			Scope: pulumi.String(scope),
		})
	}
	return out
}

func serviceEnvs(envs []*do.DigitalOceanAppEnvVar) digitalocean.AppSpecServiceEnvArray {
	out := digitalocean.AppSpecServiceEnvArray{}
	for _, e := range envs {
		key, value, typ, scope := envTriple(e)
		out = append(out, digitalocean.AppSpecServiceEnvArgs{
			Key:   pulumi.String(key),
			Value: pulumi.String(value),
			Type:  pulumi.String(typ),
			Scope: pulumi.String(scope),
		})
	}
	return out
}

func workerEnvs(envs []*do.DigitalOceanAppEnvVar) digitalocean.AppSpecWorkerEnvArray {
	out := digitalocean.AppSpecWorkerEnvArray{}
	for _, e := range envs {
		key, value, typ, scope := envTriple(e)
		out = append(out, digitalocean.AppSpecWorkerEnvArgs{
			Key:   pulumi.String(key),
			Value: pulumi.String(value),
			Type:  pulumi.String(typ),
			Scope: pulumi.String(scope),
		})
	}
	return out
}

func jobEnvs(envs []*do.DigitalOceanAppEnvVar) digitalocean.AppSpecJobEnvArray {
	out := digitalocean.AppSpecJobEnvArray{}
	for _, e := range envs {
		key, value, typ, scope := envTriple(e)
		out = append(out, digitalocean.AppSpecJobEnvArgs{
			Key:   pulumi.String(key),
			Value: pulumi.String(value),
			Type:  pulumi.String(typ),
			Scope: pulumi.String(scope),
		})
	}
	return out
}

func staticSiteEnvs(envs []*do.DigitalOceanAppEnvVar) digitalocean.AppSpecStaticSiteEnvArray {
	out := digitalocean.AppSpecStaticSiteEnvArray{}
	for _, e := range envs {
		key, value, typ, scope := envTriple(e)
		out = append(out, digitalocean.AppSpecStaticSiteEnvArgs{
			Key:   pulumi.String(key),
			Value: pulumi.String(value),
			Type:  pulumi.String(typ),
			Scope: pulumi.String(scope),
		})
	}
	return out
}

func functionEnvs(envs []*do.DigitalOceanAppEnvVar) digitalocean.AppSpecFunctionEnvArray {
	out := digitalocean.AppSpecFunctionEnvArray{}
	for _, e := range envs {
		key, value, typ, scope := envTriple(e)
		out = append(out, digitalocean.AppSpecFunctionEnvArgs{
			Key:   pulumi.String(key),
			Value: pulumi.String(value),
			Type:  pulumi.String(typ),
			Scope: pulumi.String(scope),
		})
	}
	return out
}

func componentAlerts(alerts []*do.DigitalOceanAppComponentAlert) error {
	for _, a := range alerts {
		if destinationsSet(a.GetDestinations()) {
			return errors.New(destGap)
		}
	}
	return nil
}

func healthCheck(h *do.DigitalOceanAppHealthCheck) *digitalocean.AppSpecServiceHealthCheckArgs {
	if h == nil {
		return nil
	}
	return &digitalocean.AppSpecServiceHealthCheckArgs{
		Port:                intPtrFromUint32(h.Port),
		HttpPath:            strPtr(h.GetHttpPath()),
		InitialDelaySeconds: intPtrFromUint32(h.InitialDelaySeconds),
		PeriodSeconds:       intPtrFromUint32(h.PeriodSeconds),
		TimeoutSeconds:      intPtrFromUint32(h.TimeoutSeconds),
		SuccessThreshold:    intPtrFromUint32(h.SuccessThreshold),
		FailureThreshold:    intPtrFromUint32(h.FailureThreshold),
	}
}
