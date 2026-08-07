package module

import (
	"strings"

	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// authnSecret materializes declared pre-shared API keys as the
// module-owned `<metadata.name>-authn-keys` Opaque Secret, data key
// `keys` = the comma-joined key list (the chart's keysSecret contract:
// the Deployment reads OPENFGA_AUTHN_PRESHARED_KEYS from exactly that
// key). This Secret is the ONLY place the key material lands — the chart
// values carry just its NAME (authn.preshared.keysSecret); the chart's
// plain-values key list (authn.preshared.keys), which would render every
// key into the Deployment manifest, is deliberately never used.
//
// Returns nil when authn is unset, oidc, or rides an existing Secret
// (the user owns that one). Created in the release namespace before the
// release — a secretKeyRef can only read Secrets in the workload's own
// namespace. Terraform twin: kubernetes_secret_v1.authn_keys with count.
func authnSecret(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) (*kubernetescorev1.Secret, error) {
	if locals.AuthnKeysSecretName == "" {
		return nil, nil
	}

	createdSecret, err := kubernetescorev1.NewSecret(ctx,
		locals.AuthnKeysSecretName,
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.AuthnKeysSecretName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			}),
			Type: pulumi.String("Opaque"),
			StringData: pulumi.StringMap{
				vars.AuthnKeysSecretDataKey: pulumi.String(
					strings.Join(locals.Spec.GetAuthn().GetPreshared().GetKeys(), ",")),
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependsOn...)...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create %s secret", locals.AuthnKeysSecretName)
	}

	return createdSecret, nil
}
