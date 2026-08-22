package module

import (
	"github.com/pkg/errors"
	"sigs.k8s.io/yaml"
)

// buildHelmValues renders the typed spec into the chart's values map, then
// merges the spec's helm_values escape hatch over it with Helm `-f`
// semantics (maps deep-merge with the later document winning, lists
// replace).
//
// SECRET DISCIPLINE (load-bearing): the enrollment block always renders
// the chart's existingSecret form — the module-created `<name>-token`
// Secret's NAME, never the token itself. Rendered values land in Helm's
// release Secret where anyone with release-history read access can see
// them; a Secret name reveals nothing.
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values,
// yamlencode(enrollment re-pin)] and the provider merges the documents in
// exactly this order. Keep every typed mapping below in lockstep with the
// Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec

	values := map[string]interface{}{
		// Pin the chart's fullname to the resource name: the default
		// helper would name every child `<name>-planton-runner`, which
		// wastes name budget and decouples child names from the resource
		// (the deployed Deployment is simply metadata.name).
		"fullnameOverride": locals.ReleaseName,
		"image": map[string]interface{}{
			"repository": spec.GetImageRepository(),
			"tag":        spec.GetRunnerVersion(),
		},
		"enrollment": enrollmentBlock(locals),
	}

	// ---- container sizing --------------------------------------------------------
	// Rendered ONLY when customized: the chart's own defaults (requests
	// 100m/256Mi, limits 1/1Gi) are the documented baseline, and an empty
	// requests/limits map would REPLACE them with nothing.
	if resources := spec.GetResources(); resources != nil {
		resourcesBlock := map[string]interface{}{}
		if requests := resources.GetRequests(); requests != nil {
			resourcesBlock["requests"] = cpuMemoryBlock(requests.GetCpu(), requests.GetMemory())
		}
		if limits := resources.GetLimits(); limits != nil {
			resourcesBlock["limits"] = cpuMemoryBlock(limits.GetCpu(), limits.GetMemory())
		}
		if len(resourcesBlock) > 0 {
			values["resources"] = resourcesBlock
		}
	}

	// ---- build worker ------------------------------------------------------------
	if build := spec.GetBuild(); build != nil && build.GetEnabled() {
		buildBlock := map[string]interface{}{
			"enabled": true,
		}
		if build.GetTektonNamespace() != "" {
			buildBlock["tektonNamespace"] = build.GetTektonNamespace()
		}
		values["build"] = buildBlock
	}

	// ---- escape hatch (merged LAST, Helm -f semantics) ----------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as YAML")
		}
		values = mergeMaps(values, overrides)
	}

	// The enrollment block re-pinned AFTER the merge — the one deliberate
	// exception to the escape hatch's last-word contract (twin of the
	// Terraform module's third values document). The secret wiring (a
	// Secret NAME, never inline token material) is load-bearing; letting
	// an override move it would break the secret discipline.
	values["enrollment"] = enrollmentBlock(locals)

	return values, nil
}

// enrollmentBlock renders the chart's enrollment contract: the token via
// the existingSecret form, the runner's registration name, and the
// control-plane endpoint (omitted when unset — the runner's built-in
// hosted default then applies, and the chart renders no
// PLANTON_RUNNER_ENDPOINT env at all).
func enrollmentBlock(locals *Locals) map[string]interface{} {
	block := map[string]interface{}{
		"existingSecret":    locals.TokenSecretName,
		"existingSecretKey": vars.TokenSecretKey,
		"runnerName":        locals.RunnerName,
	}
	if endpoint := locals.Spec.GetControlPlaneEndpoint(); endpoint != "" {
		block["endpoint"] = endpoint
	}
	return block
}

func cpuMemoryBlock(cpu string, memory string) map[string]interface{} {
	block := map[string]interface{}{}
	if cpu != "" {
		block["cpu"] = cpu
	}
	if memory != "" {
		block["memory"] = memory
	}
	return block
}
