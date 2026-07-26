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
// SECRET DISCIPLINE (load-bearing): githubConfigSecret always renders as
// a Secret NAME (the chart's pre-defined-secret form) — the user's own
// Secret on the existing-Secret arm, the module-materialized
// `<name>-github-auth` on the declared arms. Credential material never
// rides rendered chart values.
//
// PARITY: the Terraform module reaches the same result natively — its
// helm_release passes values = [yamlencode(typed values), helm_values]
// and the provider merges the documents in exactly this order. Keep every
// typed mapping below in lockstep with the Terraform module's locals.
func buildHelmValues(locals *Locals) (map[string]interface{}, error) {
	spec := locals.Spec
	values := map[string]interface{}{
		"githubConfigUrl":    spec.GetGithubConfigUrl(),
		"githubConfigSecret": locals.GithubAuthSecretName,
	}

	// The GitHub-visible fleet name — rendered explicitly (never left to
	// the release-name default) so the exported runs-on handle and the
	// rendered chart agree by construction.
	values["runnerScaleSetName"] = locals.RunnerScaleSetName

	if spec.GetRunnerGroup() != "" {
		values["runnerGroup"] = spec.GetRunnerGroup()
	}
	if spec.MinRunners != nil {
		values["minRunners"] = int(spec.GetMinRunners())
	}
	if spec.MaxRunners != nil {
		values["maxRunners"] = int(spec.GetMaxRunners())
	}

	// ---- container mode ---------------------------------------------------------
	if mode := spec.GetContainerMode(); mode != nil {
		modeBlock := map[string]interface{}{
			"type": mode.GetMode(),
		}
		if volume := mode.GetKubernetesWorkVolume(); volume != nil {
			modeBlock["kubernetesModeWorkVolumeClaim"] = map[string]interface{}{
				"accessModes":      []interface{}{"ReadWriteOnce"},
				"storageClassName": volume.StorageClass.GetValue(),
				"resources": map[string]interface{}{
					"requests": map[string]interface{}{
						"storage": volume.GetSize(),
					},
				},
			}
		}
		values["containerMode"] = modeBlock
	}

	// ---- the runner container -----------------------------------------------------
	// Rendered ONLY when customized: Helm values LISTS replace (never
	// merge), so any override must re-state the chart's own container
	// contract — name `runner` (the chart applies its mode wiring to the
	// container with exactly that name) and the run.sh command.
	if runner := spec.GetRunner(); runner != nil && (runner.GetImage() != "" || runner.GetResources() != nil) {
		image := runner.GetImage()
		if image == "" {
			image = vars.DefaultRunnerImage
		}
		container := map[string]interface{}{
			"name":    "runner",
			"image":   image,
			"command": []interface{}{vars.RunnerCommand},
		}
		if resources := resourcesBlock(runner.GetResources()); resources != nil {
			container["resources"] = resources
		}
		values["template"] = map[string]interface{}{
			"spec": map[string]interface{}{
				"containers": []interface{}{container},
			},
		}
	}

	// ---- proxy ----------------------------------------------------------------------
	if proxy := spec.GetProxy(); proxy != nil {
		proxyBlock := map[string]interface{}{}
		if http := proxy.GetHttp(); http != nil {
			proxyBlock["http"] = proxyServerBlock(http.GetUrl(), http.GetCredentialSecretName())
		}
		if https := proxy.GetHttps(); https != nil {
			proxyBlock["https"] = proxyServerBlock(https.GetUrl(), https.GetCredentialSecretName())
		}
		if len(proxy.GetNoProxy()) > 0 {
			noProxy := make([]interface{}, 0, len(proxy.GetNoProxy()))
			for _, host := range proxy.GetNoProxy() {
				noProxy = append(noProxy, host)
			}
			proxyBlock["noProxy"] = noProxy
		}
		if len(proxyBlock) > 0 {
			values["proxy"] = proxyBlock
		}
	}

	// ---- GitHub Enterprise Server private CA ------------------------------------------
	if tls := spec.GetGithubServerTls(); tls != nil {
		key := tls.GetKey()
		if key == "" {
			key = "ca.crt"
		}
		tlsBlock := map[string]interface{}{
			"certificateFrom": map[string]interface{}{
				"configMapKeyRef": map[string]interface{}{
					"name": tls.ConfigMapName.GetValue(),
					"key":  key,
				},
			},
		}
		if tls.GetRunnerMountPath() != "" {
			tlsBlock["runnerMountPath"] = tls.GetRunnerMountPath()
		}
		values["githubServerTLS"] = tlsBlock
	}

	// ---- explicit controller reference (single-namespace fenced controllers) ----------
	if controller := spec.GetControllerServiceAccount(); controller != nil {
		values["controllerServiceAccount"] = map[string]interface{}{
			"namespace": controller.GetNamespace(),
			"name":      controller.GetName(),
		}
	}

	// ---- escape hatch (merged LAST, Helm -f semantics) ---------------------------------
	if spec.GetHelmValues() != "" {
		overrides := map[string]interface{}{}
		if err := yaml.Unmarshal([]byte(spec.GetHelmValues()), &overrides); err != nil {
			return nil, errors.Wrap(err, "failed to parse helm_values as YAML")
		}
		values = mergeMaps(values, overrides)
	}

	// githubConfigSecret re-pinned AFTER the merge — the one deliberate
	// exception to the escape hatch's last-word contract (twin of the
	// Terraform module's third values document). The credential contract
	// (a Secret NAME, never inline material) is load-bearing; letting an
	// override move it would break the secret discipline.
	values["githubConfigSecret"] = locals.GithubAuthSecretName

	return values, nil
}

func proxyServerBlock(url string, credentialSecretName string) map[string]interface{} {
	block := map[string]interface{}{
		"url": url,
	}
	if credentialSecretName != "" {
		block["credentialSecretRef"] = credentialSecretName
	}
	return block
}
