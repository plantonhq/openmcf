package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi-random/sdk/v4/go/random"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The generation-shape arguments are ignored after creation so an
// IMPORTED credential never silently regenerates: rotation stays an
// explicit verb, never plan fallout. Terraform twin: random_password
// lifecycle ignore_changes on the same arguments.
var generationShapeIgnores = pulumi.IgnoreChanges([]string{
	"length", "special", "upper", "lower", "numeric",
	"minLower", "minNumeric", "minSpecial", "minUpper", "overrideSpecial",
})

// webAuthSecret materializes the `<name>-auth` Secret when the web-UI
// login is on — the only place the credential lives:
//
//   - `username` / `password`: the login credential (the exported
//     handle operators, clients and verifiers read).
//   - `flask-secret-key`: the Flask session-signing key the login
//     backend loads — generated once and stable, so sessions survive
//     pod restarts (a per-start random would log every user out on
//     every roll).
//
// All three reach the pods as Secret-projected FILES through the
// chart's mount_external_secret seam (values.go) — never as
// environment variables, rendered values or process arguments.
func webAuthSecret(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) ([]pulumi.Resource, error) {
	if !locals.WebLoginEnabled {
		return nil, nil
	}

	// Letters+digits only: operators type this at the login form —
	// symbol-free avoids quoting bugs; the larger length compensates
	// the smaller alphabet.
	loginPassword, err := random.NewRandomPassword(ctx, "web-ui-password",
		&random.RandomPasswordArgs{
			Length:     pulumi.Int(24),
			Special:    pulumi.Bool(false),
			MinUpper:   pulumi.Int(2),
			MinLower:   pulumi.Int(2),
			MinNumeric: pulumi.Int(2),
		},
		generationShapeIgnores)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate web-ui password")
	}

	// The session-signing key never leaves the pod filesystem — 64
	// letters+digits.
	flaskSecretKey, err := random.NewRandomPassword(ctx, "flask-secret-key",
		&random.RandomPasswordArgs{
			Length:  pulumi.Int(64),
			Special: pulumi.Bool(false),
		},
		generationShapeIgnores)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate flask secret key")
	}

	created, err := kubernetescorev1.NewSecret(ctx, "auth-secret",
		&kubernetescorev1.SecretArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.AuthSecretName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			}),
			Type: pulumi.String("Opaque"),
			StringData: pulumi.StringMap{
				vars.AuthUsernameKey:       pulumi.String(locals.WebUsername),
				vars.AuthPasswordKey:       loginPassword.Result,
				vars.AuthFlaskSecretKeyKey: flaskSecretKey.Result,
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependsOn...)...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create web-ui auth secret")
	}

	return []pulumi.Resource{created}, nil
}
