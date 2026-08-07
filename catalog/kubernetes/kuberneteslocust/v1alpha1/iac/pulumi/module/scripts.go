package module

import (
	"github.com/pkg/errors"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// scriptConfigMaps materializes the module-owned ConfigMaps the pods
// mount at startup (release dependencies):
//
//   - `<name>-locustfile` + `<name>-lib` (inline arm only): the user's
//     locustfile as `main.py` and the supporting modules. The chart
//     mounts them at the locustfile path (+`/lib`); content changes
//     roll the pods through the module's checksum annotation (the
//     chart's own checksums cover only chart-rendered ConfigMaps).
//   - `<name>-web-auth` (login on): the module-owned login backend the
//     master loads alongside the locustfile (webauth.go — the chart's
//     extraConfigMaps seam mounts it; the master's `-f` argument names
//     it).
//
// The existing-ConfigMaps arm creates no script ConfigMaps — the chart
// values reference the user's own.
func scriptConfigMaps(ctx *pulumi.Context,
	locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependsOn []pulumi.ResourceOption,
) ([]pulumi.Resource, error) {
	var createdResources []pulumi.Resource

	newConfigMap := func(resourceName, configMapName string, data map[string]string) (pulumi.Resource, error) {
		created, err := kubernetescorev1.NewConfigMap(ctx, resourceName,
			&kubernetescorev1.ConfigMapArgs{
				Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
					Name:      pulumi.String(configMapName),
					Namespace: pulumi.String(locals.Namespace),
					Labels:    pulumi.ToStringMap(locals.Labels),
				}),
				Data: pulumi.ToStringMap(data),
			}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependsOn...)...)
		if err != nil {
			return nil, err
		}
		return created, nil
	}

	if inline := locals.Spec.GetLoadTest().GetInline(); inline != nil {
		created, err := newConfigMap("locustfile", locals.LocustfileConfigMap, map[string]string{
			locals.LocustfileName: inline.GetLocustfileContent(),
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create locustfile configmap")
		}
		createdResources = append(createdResources, created)

		if locals.LibConfigMap != "" {
			created, err := newConfigMap("lib", locals.LibConfigMap, inline.GetLibFiles())
			if err != nil {
				return nil, errors.Wrap(err, "failed to create lib configmap")
			}
			createdResources = append(createdResources, created)
		}
	}

	if locals.WebLoginEnabled {
		created, err := newConfigMap("web-auth-code", locals.WebAuthCodeName, map[string]string{
			"planton_auth.py": webAuthBackendPy,
		})
		if err != nil {
			return nil, errors.Wrap(err, "failed to create web-auth code configmap")
		}
		createdResources = append(createdResources, created)
	}

	return createdResources, nil
}
