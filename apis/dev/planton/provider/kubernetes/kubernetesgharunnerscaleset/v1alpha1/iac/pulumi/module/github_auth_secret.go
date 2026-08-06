package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// githubAuthSecret materializes the declared GitHub credential (PAT or
// GitHub App) as the `<name>-github-auth` Secret with the chart's own
// pre-defined-secret key contract — `github_token` for a PAT;
// `github_app_id` / `github_app_installation_id` /
// `github_app_private_key` for an App. The chart then reads the Secret BY
// NAME, so credential material never rides rendered chart values (the
// secret discipline). Returns nil on the existing-Secret arm (the user
// owns that Secret). Terraform twin: kubernetes_secret_v1 with count.
func githubAuthSecret(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (*kubernetescorev1.Secret, error) {
	if !locals.MaterializeAuthSecret {
		return nil, nil
	}

	stringData := pulumi.StringMap{}
	if pat := locals.Spec.GetAuth().GetPat(); pat != nil {
		stringData["github_token"] = pulumi.String(pat.GetToken())
	}
	if app := locals.Spec.GetAuth().GetGithubApp(); app != nil {
		stringData["github_app_id"] = pulumi.String(app.GetAppId())
		stringData["github_app_installation_id"] = pulumi.String(app.GetInstallationId())
		stringData["github_app_private_key"] = pulumi.String(app.GetPrivateKey())
	}

	createdSecret, err := kubernetescorev1.NewSecret(ctx,
		locals.GithubAuthSecretName,
		&kubernetescorev1.SecretArgs{
			Metadata: &kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.GithubAuthSecretName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			},
			StringData: stringData,
		},
		append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create %s secret", locals.GithubAuthSecretName)
	}

	return createdSecret, nil
}
