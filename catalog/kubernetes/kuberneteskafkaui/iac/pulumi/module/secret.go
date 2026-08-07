package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// consoleSecret materializes the LITERAL credentials the spec declares as
// the "<metadata.name>-secrets" Opaque Secret — today that is exactly the
// console login password, under the "console-user-password" key the
// KAFKA_UI_USER_PASSWORD secret mapping points at (helm_release.go).
//
// LOGIN_FORM supports exactly one account (Spring Boot's default security
// user — see the rendering comment in helm_release.go), which is why the
// spec models a single `user` and this Secret stores one password.
//
// Referenced credentials (sasl / schema registry / Connect password_secret
// entries) are NOT copied here — their mappings point at the source
// Secrets directly. This Secret is the only place the declared literal
// lands; it never transits chart values. Returns nil when auth is not
// declared. Terraform twin: kubernetes_secret_v1.console with the same
// name, key, and contents.
func consoleSecret(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	if !locals.AuthEnabled {
		return nil, nil
	}

	data := map[string]string{
		consolePasswordSecretKey: locals.Spec.GetAuth().GetUser().GetPassword(),
	}

	createdSecret, err := kubernetescorev1.NewSecret(ctx, locals.ConsoleSecretName,
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.ConsoleSecretName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			}),
			Type:       pulumi.String("Opaque"),
			StringData: pulumi.ToStringMap(data),
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependsOn...)...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create console secret")
	}

	return createdSecret, nil
}
