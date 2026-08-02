package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apiextensions"
	batchv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/batch/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// service renders the MLflow Service `<name>` (this kind's front door —
// ClusterIP by default; exposure composes from first-class kinds over the
// exported handle).
func service(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (*kubernetescorev1.Service, error) {
	serviceSpec := locals.Spec.GetService()
	serviceType := serviceSpec.GetType()
	if serviceType == "" {
		serviceType = "ClusterIP"
	}

	metadata := &kubernetesmeta.ObjectMetaArgs{
		Name:      pulumi.String(locals.Name),
		Namespace: pulumi.String(locals.Namespace),
		Labels:    pulumi.ToStringMap(locals.Labels),
	}
	if len(serviceSpec.GetAnnotations()) > 0 {
		metadata.Annotations = pulumi.ToStringMap(serviceSpec.GetAnnotations())
	}

	createdService, err := kubernetescorev1.NewService(ctx, "mlflow-service",
		&kubernetescorev1.ServiceArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(metadata),
			Spec: &kubernetescorev1.ServiceSpecArgs{
				Type:     pulumi.String(serviceType),
				Selector: pulumi.ToStringMap(locals.SelectorLabels),
				Ports: kubernetescorev1.ServicePortArray{
					&kubernetescorev1.ServicePortArgs{
						Name:       pulumi.String("http"),
						Port:       pulumi.Int(vars.ServerPort),
						TargetPort: pulumi.String("http"),
						Protocol:   pulumi.String("TCP"),
					},
				},
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create mlflow service")
	}
	return createdService, nil
}

// pvcs renders the module-owned PersistentVolumeClaims the arms need:
// `<name>-data` (sqlite backend — tracking db + auth db) and
// `<name>-artifacts` (PVC artifact arm). Both ReadWriteOnce — the CEL
// contract caps replicas at 1 whenever either exists.
func pvcs(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) ([]pulumi.Resource, error) {
	var created []pulumi.Resource

	newPvc := func(resourceName, pvcName, size, storageClass string) (pulumi.Resource, error) {
		pvcSpec := &kubernetescorev1.PersistentVolumeClaimSpecArgs{
			AccessModes: pulumi.ToStringArray([]string{"ReadWriteOnce"}),
			Resources: &kubernetescorev1.VolumeResourceRequirementsArgs{
				Requests: pulumi.StringMap{"storage": pulumi.String(size)},
			},
		}
		if storageClass != "" {
			pvcSpec.StorageClassName = pulumi.String(storageClass)
		}
		return kubernetescorev1.NewPersistentVolumeClaim(ctx, resourceName,
			&kubernetescorev1.PersistentVolumeClaimArgs{
				Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
					Name:      pulumi.String(pvcName),
					Namespace: pulumi.String(locals.Namespace),
					Labels:    pulumi.ToStringMap(locals.Labels),
				}),
				Spec: pvcSpec,
			}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
	}

	if locals.DataPvcEnabled {
		size := "5Gi"
		storageClass := ""
		if sqlite := locals.Spec.GetBackendStore().GetSqlitePvc(); sqlite != nil {
			if sqlite.GetStorageSize() != "" {
				size = sqlite.GetStorageSize()
			}
			storageClass = sqlite.GetStorageClass()
		}
		createdPvc, err := newPvc("data-pvc", locals.Name+vars.DataPvcSuffix, size, storageClass)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create data pvc")
		}
		created = append(created, createdPvc)
	}

	if locals.ArtifactsPvcEnabled {
		size := "10Gi"
		storageClass := ""
		if pvcArm := locals.Spec.GetArtifactStore().GetPvc(); pvcArm != nil {
			if pvcArm.GetStorageSize() != "" {
				size = pvcArm.GetStorageSize()
			}
			storageClass = pvcArm.GetStorageClass()
		}
		createdPvc, err := newPvc("artifacts-pvc", locals.Name+vars.ArtifactsPvcSuffix, size, storageClass)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create artifacts pvc")
		}
		created = append(created, createdPvc)
	}

	return created, nil
}

// gcCronJob renders the `mlflow gc` CronJob (`<name>-gc`) — permanent
// removal of soft-deleted runs/experiments older than the retention
// window. Same image, same backend env; artifact-store credentials ride
// along so gc can delete the artifacts too.
func gcCronJob(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) error {
	args := []string{
		"mlflow", "gc",
		"--older-than", locals.GcOlderThan,
	}
	// The backend URI must be a pod argument here — `mlflow gc` takes
	// it as a flag. It rides env expansion ($(VAR) — resolved by the
	// kubelet from the Secret-sourced env var) so the URI itself still
	// never appears in the rendered spec.
	args = append(args, "--backend-store-uri", "$(MLFLOW_BACKEND_STORE_URI)")

	podSpec := &kubernetescorev1.PodSpecArgs{
		RestartPolicy: pulumi.String("OnFailure"),
		Containers: kubernetescorev1.ContainerArray{
			&kubernetescorev1.ContainerArgs{
				Name:    pulumi.String("mlflow-gc"),
				Image:   pulumi.String(locals.Image),
				Command: pulumi.ToStringArray(args),
				Env:     serverEnv(locals),
			},
		},
	}
	if scheduling := locals.Spec.GetScheduling(); scheduling != nil {
		if len(scheduling.GetNodeSelector()) > 0 {
			podSpec.NodeSelector = pulumi.ToStringMap(scheduling.GetNodeSelector())
		}
		if len(scheduling.GetTolerations()) > 0 {
			podSpec.Tolerations = tolerationsArgs(scheduling.GetTolerations())
		}
	}
	// The sqlite arm's database lives on the data PVC — gc must mount
	// it (RWO: the CronJob only runs cleanly on shapes where the server
	// pod and the job can share the volume's node; the postgres arm is
	// the real gc story).
	if locals.DataPvcEnabled {
		podSpec.Volumes = kubernetescorev1.VolumeArray{
			&kubernetescorev1.VolumeArgs{
				Name: pulumi.String("data"),
				PersistentVolumeClaim: &kubernetescorev1.PersistentVolumeClaimVolumeSourceArgs{
					ClaimName: pulumi.String(locals.Name + vars.DataPvcSuffix),
				},
			},
		}
		containers := podSpec.Containers.(kubernetescorev1.ContainerArray)
		containers[0].(*kubernetescorev1.ContainerArgs).VolumeMounts = kubernetescorev1.VolumeMountArray{
			&kubernetescorev1.VolumeMountArgs{
				Name:      pulumi.String("data"),
				MountPath: pulumi.String(vars.DataMountPath),
			},
		}
	}

	_, err := batchv1.NewCronJob(ctx, "mlflow-gc-cronjob",
		&batchv1.CronJobArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.Name + vars.GcCronJobSuffix),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			}),
			Spec: &batchv1.CronJobSpecArgs{
				Schedule:          pulumi.String(locals.GcSchedule),
				ConcurrencyPolicy: pulumi.String("Forbid"),
				JobTemplate: &batchv1.JobTemplateSpecArgs{
					Spec: &batchv1.JobSpecArgs{
						BackoffLimit: pulumi.Int(2),
						Template: &kubernetescorev1.PodTemplateSpecArgs{
							Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
								Labels: pulumi.ToStringMap(locals.Labels),
							}),
							Spec: podSpec,
						},
					},
				},
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
	if err != nil {
		return errors.Wrap(err, "failed to create gc cronjob")
	}
	return nil
}

// serviceMonitor renders a monitoring.coreos.com ServiceMonitor
// (`<name>-metrics`) for operator-based scraping. Requires the Prometheus
// operator CRDs on the cluster — deploying without them fails loudly (by
// design; the spec documents the prerequisite).
func serviceMonitor(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) error {
	_, err := apiextensions.NewCustomResource(ctx, "mlflow-service-monitor",
		&apiextensions.CustomResourceArgs{
			ApiVersion: pulumi.String("monitoring.coreos.com/v1"),
			Kind:       pulumi.String("ServiceMonitor"),
			Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.Name + vars.ServiceMonitorSuffix),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			}),
			OtherFields: map[string]interface{}{
				"spec": map[string]interface{}{
					"selector": map[string]interface{}{
						"matchLabels": locals.SelectorLabels,
					},
					"endpoints": []interface{}{
						map[string]interface{}{
							"port":     "http",
							"path":     "/metrics",
							"interval": "30s",
						},
					},
				},
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
	if err != nil {
		return errors.Wrap(err, "failed to create service monitor")
	}
	return nil
}
