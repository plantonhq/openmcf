package module

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// webhookGCLabelSelector selects every webhook configuration Kyverno
// registers at runtime. Source-verified at the pinned chart: the same
// label the chart's delete-webhooks helper (and the manual unstick
// command on the spec) use.
const webhookGCLabelSelector = "webhook.kyverno.io/managed-by=kyverno"

// webhookGCConfigMap creates the destroy-ordered sentinel ConfigMap
// (Terraform twin: kubernetes_config_map_v1.webhook_gc). The helm Release
// depends on it so destroy tears the release down first; the
// BeforeDelete hook then deletes any stranded kyverno-* webhook
// configurations the chart's broken delete-webhooks helper left behind
// (upstream kyverno/kyverno#16492 — ValidatingAdmissionPolicies deleted
// instead of ValidatingWebhookConfigurations).
//
// Delete hooks require `pulumi destroy --run-program` so the program
// that registered the hook is available; the E2E runner and the Planton
// Pulumi stack runner pass that flag on destroy.
func webhookGCConfigMap(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.Resource,
) (*kubernetescorev1.ConfigMap, error) {
	gcHook, err := ctx.RegisterResourceHook("kyverno-webhook-gc", func(args *pulumi.ResourceHookArgs) error {
		return deleteKyvernoWebhooks()
	}, nil)
	if err != nil {
		return nil, errors.Wrap(err, "registering kyverno webhook-gc hook")
	}

	opts := []pulumi.ResourceOption{
		pulumi.Provider(kubernetesProvider),
		pulumi.ResourceHooks(&pulumi.ResourceHookBinding{
			BeforeDelete: []*pulumi.ResourceHook{gcHook},
		}),
	}
	if len(dependencies) > 0 {
		opts = append(opts, pulumi.DependsOn(dependencies))
	}

	return kubernetescorev1.NewConfigMap(ctx, locals.WebhookGCConfigMapName,
		&kubernetescorev1.ConfigMapArgs{
			Metadata: &kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.WebhookGCConfigMapName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			},
			Data: pulumi.ToStringMap(map[string]string{
				"purpose": "kyverno-webhook-gc",
			}),
		}, opts...)
}

// deleteKyvernoWebhooks removes every runtime-registered Kyverno webhook
// configuration by label. Tolerates absence — a clean chart uninstall
// (once upstream fixes delete-webhooks) leaves nothing to delete.
func deleteKyvernoWebhooks() error {
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBE_CONFIG_PATH")
	}
	args := []string{
		"delete",
		"validatingwebhookconfiguration,mutatingwebhookconfiguration",
		"-l", webhookGCLabelSelector,
		"--ignore-not-found=true",
	}
	if kubeconfig != "" {
		args = append([]string{"--kubeconfig", kubeconfig}, args...)
	}
	cmd := exec.Command("kubectl", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		msg := string(out)
		if len(msg) > 240 {
			msg = msg[:240]
		}
		return fmt.Errorf("kyverno webhook gc: %w: %s", err, msg)
	}
	return nil
}
