package module

import (
	clickhousekeeperv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/altinityoperator/kubernetes/clickhouse_keeper/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// keeper renders the managed ClickHouseKeeperInstallation `<name>-keeper`
// when the coordination contract calls for one (see Locals.DeployKeeper).
// The CHI references it by name through its native keeper wiring; the
// operator resolves the ensemble endpoints itself.
func keeper(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	if !locals.DeployKeeper {
		return nil, nil
	}

	keeperSpec := locals.Spec.GetCoordination().GetKeeper()

	replicas := vars.DefaultKeeperReplicas
	if keeperSpec.GetReplicas() > 0 {
		replicas = int(keeperSpec.GetReplicas())
	}
	diskSize := vars.DefaultKeeperDiskSize
	if keeperSpec.GetDiskSize() != "" {
		diskSize = keeperSpec.GetDiskSize()
	}

	// The Keeper container is declared explicitly so the image pins to
	// the resource's own version line instead of the operator's fallback
	// (`latest`) — Keeper images are published in lockstep with server
	// releases and the protocol is compatible across them. The image
	// must NEVER be dropped from this container: an explicit container
	// entry suppresses the operator's default-image injection entirely —
	// verified live: a pod template carrying only resources produced a
	// StatefulSet the API server rejected with `containers[0].image:
	// Required value`, and the keeper never came up.
	container := pulumi.Map{
		"name":  pulumi.String("clickhouse-keeper"),
		"image": pulumi.String(vars.KeeperImageRepo + ":" + locals.Spec.GetVersion()),
	}
	if resources := containerResourcesMap(keeperSpec.GetResources()); resources != nil {
		container["resources"] = resources
	}

	return clickhousekeeperv1.NewClickHouseKeeperInstallation(ctx, locals.KeeperName,
		&clickhousekeeperv1.ClickHouseKeeperInstallationArgs{
			Metadata: kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.KeeperName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			},
			Spec: clickhousekeeperv1.ClickHouseKeeperInstallationSpecArgs{
				Configuration: clickhousekeeperv1.ClickHouseKeeperInstallationSpecConfigurationArgs{
					Clusters: clickhousekeeperv1.ClickHouseKeeperInstallationSpecConfigurationClustersArray{
						clickhousekeeperv1.ClickHouseKeeperInstallationSpecConfigurationClustersArgs{
							Name: pulumi.String(vars.KeeperClusterName),
							Layout: clickhousekeeperv1.ClickHouseKeeperInstallationSpecConfigurationClustersLayoutArgs{
								ReplicasCount: pulumi.Int(replicas),
							},
						},
					},
				},
				Defaults: clickhousekeeperv1.ClickHouseKeeperInstallationSpecDefaultsArgs{
					Templates: clickhousekeeperv1.ClickHouseKeeperInstallationSpecDefaultsTemplatesArgs{
						PodTemplate:             pulumi.String("keeper"),
						DataVolumeClaimTemplate: pulumi.String("data"),
					},
				},
				Templates: clickhousekeeperv1.ClickHouseKeeperInstallationSpecTemplatesArgs{
					PodTemplates: clickhousekeeperv1.ClickHouseKeeperInstallationSpecTemplatesPodTemplatesArray{
						clickhousekeeperv1.ClickHouseKeeperInstallationSpecTemplatesPodTemplatesArgs{
							Name: pulumi.String("keeper"),
							Spec: pulumi.Map{
								"containers": pulumi.Array{container},
							},
						},
					},
					VolumeClaimTemplates: clickhousekeeperv1.ClickHouseKeeperInstallationSpecTemplatesVolumeClaimTemplatesArray{
						clickhousekeeperv1.ClickHouseKeeperInstallationSpecTemplatesVolumeClaimTemplatesArgs{
							Name: pulumi.String("data"),
							Spec: pvcSpecMap(diskSize, keeperSpec.GetStorageClass().GetValue()),
						},
					},
				},
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
}
