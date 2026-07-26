package module

import (
	kubernetesrabbitmqoperatorv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetesrabbitmqoperator/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// deploymentTransformation returns the yaml.Transformation that applies the
// spec's typed overrides onto the release manifest's operator Deployment —
// every other document applies verbatim (faithful distribution). The
// Terraform twin performs the identical patch through locals merges.
//
// Patched surfaces (all optional):
//   - OPERATOR_SCOPE_NAMESPACE env (watch_namespaces, comma-joined; absent
//     = the operator watches ALL namespaces, the upstream default),
//   - DEFAULT_RABBITMQ_IMAGE / DEFAULT_USER_UPDATER_IMAGE env (fleet-wide
//     image defaults for air-gapped clusters),
//   - the operator container image (operator_image),
//   - the operator container resources,
//   - pod nodeSelector / tolerations / imagePullSecrets.
func deploymentTransformation(spec *kubernetesrabbitmqoperatorv1.KubernetesRabbitMqOperatorSpec) func(state map[string]interface{}, opts ...pulumi.ResourceOption) {
	return func(state map[string]interface{}, _ ...pulumi.ResourceOption) {
		if state["kind"] != "Deployment" {
			return
		}
		metadata, _ := state["metadata"].(map[string]interface{})
		if metadata == nil || metadata["name"] != vars.DeploymentName {
			return
		}

		deploymentSpec, _ := state["spec"].(map[string]interface{})
		template, _ := deploymentSpec["template"].(map[string]interface{})
		podSpec, _ := template["spec"].(map[string]interface{})
		if podSpec == nil {
			return
		}

		containers, _ := podSpec["containers"].([]interface{})
		if len(containers) > 0 {
			// The release manifest ships exactly one container ("operator");
			// patch it in place.
			container, _ := containers[0].(map[string]interface{})
			if container != nil {
				if extraEnv := extraEnvEntries(spec); len(extraEnv) > 0 {
					env, _ := container["env"].([]interface{})
					container["env"] = append(env, extraEnv...)
				}
				if image := resolveOperatorImage(spec); image != "" {
					container["image"] = image
				}
				if resources := resourcesMap(spec); resources != nil {
					container["resources"] = resources
				}
			}
		}

		if len(spec.GetNodeSelector()) > 0 {
			nodeSelector := map[string]interface{}{}
			for key, value := range spec.GetNodeSelector() {
				nodeSelector[key] = value
			}
			podSpec["nodeSelector"] = nodeSelector
		}
		if tolerations := tolerationsSlice(spec); tolerations != nil {
			podSpec["tolerations"] = tolerations
		}
		if names := pullSecretNames(spec); len(names) > 0 {
			pullSecrets := []interface{}{}
			for _, name := range names {
				pullSecrets = append(pullSecrets, map[string]interface{}{"name": name})
			}
			podSpec["imagePullSecrets"] = pullSecrets
		}
	}
}

// extraEnvEntries builds the env entries appended to the release
// manifest's own (which carry OPERATOR_NAMESPACE via fieldRef).
func extraEnvEntries(spec *kubernetesrabbitmqoperatorv1.KubernetesRabbitMqOperatorSpec) []interface{} {
	entries := []interface{}{}
	if len(spec.GetWatchNamespaces()) > 0 {
		entries = append(entries, map[string]interface{}{
			"name":  "OPERATOR_SCOPE_NAMESPACE",
			"value": joinComma(spec.GetWatchNamespaces()),
		})
	}
	if spec.GetDefaultRabbitmqImage() != "" {
		entries = append(entries, map[string]interface{}{
			"name":  "DEFAULT_RABBITMQ_IMAGE",
			"value": spec.GetDefaultRabbitmqImage(),
		})
	}
	if spec.GetDefaultUserUpdaterImage() != "" {
		entries = append(entries, map[string]interface{}{
			"name":  "DEFAULT_USER_UPDATER_IMAGE",
			"value": spec.GetDefaultUserUpdaterImage(),
		})
	}
	return entries
}

func joinComma(values []string) string {
	out := ""
	for i, value := range values {
		if i > 0 {
			out += ","
		}
		out += value
	}
	return out
}

// pullSecretNames joins image_pull_secrets with operator_image's own
// pull_secret_name, deduplicated — a private image override naturally
// travels with its own pull secret (Terraform twin:
// image_pull_secret_names in locals.tf).
func pullSecretNames(spec *kubernetesrabbitmqoperatorv1.KubernetesRabbitMqOperatorSpec) []string {
	names := append([]string{}, spec.GetImagePullSecrets()...)
	if extra := spec.GetOperatorImage().GetPullSecretName(); extra != "" {
		seen := false
		for _, name := range names {
			if name == extra {
				seen = true
				break
			}
		}
		if !seen {
			names = append(names, extra)
		}
	}
	return names
}

// resolveOperatorImage joins repo:tag when either is set; empty keeps the
// release manifest's pinned image.
func resolveOperatorImage(spec *kubernetesrabbitmqoperatorv1.KubernetesRabbitMqOperatorSpec) string {
	image := spec.GetOperatorImage()
	if image.GetRepo() == "" && image.GetTag() == "" {
		return ""
	}
	if image.GetRepo() == "" || image.GetTag() == "" {
		return image.GetRepo() + image.GetTag()
	}
	return image.GetRepo() + ":" + image.GetTag()
}

// resourcesMap renders the operator container resources override; nil
// keeps the release manifest's defaults (200m/500Mi requests and limits).
func resourcesMap(spec *kubernetesrabbitmqoperatorv1.KubernetesRabbitMqOperatorSpec) map[string]interface{} {
	resources := spec.GetResources()
	if resources == nil {
		return nil
	}

	out := map[string]interface{}{}
	if requests := quantityEntries(resources.GetRequests().GetCpu(), resources.GetRequests().GetMemory()); requests != nil {
		out["requests"] = requests
	}
	if limits := quantityEntries(resources.GetLimits().GetCpu(), resources.GetLimits().GetMemory()); limits != nil {
		out["limits"] = limits
	}
	if len(out) == 0 {
		return nil
	}
	return out
}

func quantityEntries(cpu string, memory string) map[string]interface{} {
	quantities := map[string]interface{}{}
	if cpu != "" {
		quantities["cpu"] = cpu
	}
	if memory != "" {
		quantities["memory"] = memory
	}
	if len(quantities) == 0 {
		return nil
	}
	return quantities
}

// tolerationsSlice renders the pod tolerations override; nil when empty.
func tolerationsSlice(spec *kubernetesrabbitmqoperatorv1.KubernetesRabbitMqOperatorSpec) []interface{} {
	if len(spec.GetTolerations()) == 0 {
		return nil
	}

	out := []interface{}{}
	for _, toleration := range spec.GetTolerations() {
		entry := map[string]interface{}{}
		if toleration.GetKey() != "" {
			entry["key"] = toleration.GetKey()
		}
		if toleration.GetOperator() != "" {
			entry["operator"] = toleration.GetOperator()
		}
		if toleration.GetValue() != "" {
			entry["value"] = toleration.GetValue()
		}
		if toleration.GetEffect() != "" {
			entry["effect"] = toleration.GetEffect()
		}
		if toleration.TolerationSeconds != nil {
			entry["tolerationSeconds"] = int(toleration.GetTolerationSeconds())
		}
		out = append(out, entry)
	}
	return out
}
