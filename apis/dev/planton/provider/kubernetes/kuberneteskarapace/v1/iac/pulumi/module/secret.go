package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// saslPasswordSecret materializes a literal spec.kafka.sasl.password into
// the module-owned Secret `<metadata.name>-sasl` (key "password").
//
// WHY A SECRET AND NOT A PLAIN ENV VALUE: the pod spec is readable by
// anyone with get-pod RBAC (and lands in audit logs, kubectl describe
// output, and controller caches); Secret VALUES have their own,
// stricter ACL. Materializing the declared password and wiring
// KARAPACE_SASL_PLAIN_PASSWORD through a secretKeyRef keeps the
// credential out of every one of those surfaces. When the spec instead
// references an existing Secret (password_secret — the
// KubernetesKafkaUser composition path), no Secret is created here and
// the env var references that Secret directly.
//
// Returns nil when the spec carries no literal password.
func saslPasswordSecret(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (*kubernetescorev1.Secret, error) {
	if !locals.CreateSaslSecret {
		return nil, nil
	}

	secretArgs := &kubernetescorev1.SecretArgs{
		Metadata: &kubernetesmeta.ObjectMetaArgs{
			Name:      pulumi.String(locals.SaslSecretName),
			Namespace: pulumi.String(locals.Namespace),
			Labels:    pulumi.ToStringMap(locals.Labels),
		},
		Type: pulumi.String("Opaque"),
		StringData: pulumi.StringMap{
			"password": pulumi.String(locals.Spec.Kafka.Sasl.Password),
		},
	}

	opts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)
	createdSecret, err := kubernetescorev1.NewSecret(ctx,
		locals.SaslSecretName,
		secretArgs,
		opts...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create sasl password secret")
	}

	return createdSecret, nil
}
