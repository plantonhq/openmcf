package module

import (
	"fmt"

	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// authSecret materializes the declared admin password as the
// "<metadata.name>-auth" Opaque Secret with the chart's contract: ONE key,
// NEO4J_AUTH, whose value is "neo4j/<password>". The chart consumes it via
// neo4j.passwordFromSecret and LOOKS THE SECRET UP AT TEMPLATE TIME (its
// neo4j.secretName helper fails the install when the Secret is missing or
// lacks NEO4J_AUTH), so this Secret must exist BEFORE the Helm release —
// main.go wires the explicit dependency. Because passwordFromSecret is set,
// the chart renders no auth Secret of its own — this Secret is the only
// place the credential lands, and it never transits chart values. Returns
// nil when the password arm is not declared (existing_secret and absent-auth
// arms create nothing). Terraform twin: kubernetes_secret_v1.auth with the
// same name, key, and contents.
func authSecret(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	if !locals.CreateAuthSecret {
		return nil, nil
	}

	// Marked secret so the value is encrypted in the Pulumi state — twin
	// of the Terraform module's sensitive() wrap.
	neo4jAuth := pulumi.ToSecret(
		pulumi.String(fmt.Sprintf("neo4j/%s", locals.AdminPassword)),
	).(pulumi.StringOutput)

	createdSecret, err := kubernetescorev1.NewSecret(ctx, locals.AuthSecretName,
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.AuthSecretName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			}),
			Type: pulumi.String("Opaque"),
			StringData: pulumi.StringMap{
				"NEO4J_AUTH": neo4jAuth,
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependsOn...)...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create auth secret")
	}

	return createdSecret, nil
}
