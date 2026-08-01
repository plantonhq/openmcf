package module

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	kubernetesv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// serverArgs renders the `mlflow server` command line — the shape follows
// upstream's own deployment reference (the in-repo chart at the pin):
// host/port/workers as flags, the artifact destination as a flag (the
// server PROXIES artifact traffic — clients never carry store
// credentials), basic-auth via --app-name, metrics via
// --expose-prometheus. The backend-store URI deliberately rides the
// MLFLOW_BACKEND_STORE_URI env var (from the module's Secret), never a
// pod argument — arguments are readable in every rendered pod spec.
func serverArgs(locals *Locals) []string {
	args := []string{
		"mlflow", "server",
		"--host", "0.0.0.0",
		"--port", strconv.Itoa(vars.ServerPort),
		"--workers", strconv.Itoa(locals.Workers),
		"--artifacts-destination", locals.ArtifactDestination,
		"--serve-artifacts",
	}
	if locals.AuthEnabled {
		args = append(args, "--app-name", "basic-auth")
	}
	if locals.MetricsEnabled {
		args = append(args, "--expose-prometheus", vars.MetricsExportPathFlag)
	}
	args = append(args, locals.Spec.GetExtraArgs()...)
	return args
}

// serverEnv renders the server container's environment: the backend URI
// (Secret-sourced on database arms, a literal sqlite path on the PVC
// arm), the auth config path, the artifact-store credentials (all
// Secret-sourced), and the spec's extra env maps.
func serverEnv(locals *Locals) kubernetescorev1.EnvVarArray {
	envVars := kubernetescorev1.EnvVarArray{}

	if locals.BackendType == "sqlite" {
		envVars = append(envVars, envVar("MLFLOW_BACKEND_STORE_URI", locals.SqliteBackendUri))
	} else {
		envVars = append(envVars, secretEnvVar("MLFLOW_BACKEND_STORE_URI",
			locals.BackendUriSecretName, vars.BackendUriKey))
	}

	if locals.AuthEnabled {
		envVars = append(envVars, envVar("MLFLOW_AUTH_CONFIG_PATH",
			fmt.Sprintf("%s/%s", vars.AuthConfigMountPath, vars.AuthConfigFileName)))
	}

	// Artifact-store credentials — env contracts of MLflow's own
	// artifact clients (boto3 / google-cloud-storage / azure sdk).
	artifactStore := locals.Spec.GetArtifactStore()
	switch locals.ArtifactType {
	case "s3_compatible":
		s3 := artifactStore.GetS3Compatible()
		credentials := s3.GetCredentialsSecret()
		secretName := credentials.GetSecretName().GetValue()
		accessKeyIdKey := credentials.GetAccessKeyIdKey()
		if accessKeyIdKey == "" {
			accessKeyIdKey = "admin_access_key_id"
		}
		secretAccessKeyKey := credentials.GetSecretAccessKeyKey()
		if secretAccessKeyKey == "" {
			secretAccessKeyKey = "admin_secret_access_key"
		}
		envVars = append(envVars,
			envVar("MLFLOW_S3_ENDPOINT_URL", s3.GetEndpoint().GetValue()),
			secretEnvVar("AWS_ACCESS_KEY_ID", secretName, accessKeyIdKey),
			secretEnvVar("AWS_SECRET_ACCESS_KEY", secretName, secretAccessKeyKey),
		)
	case "aws_s3":
		s3 := artifactStore.GetAwsS3()
		envVars = append(envVars, envVar("AWS_DEFAULT_REGION", s3.GetRegion()))
		if credentials := s3.GetCredentialsSecret(); credentials != nil {
			accessKeyIdKey := credentials.GetAccessKeyIdKey()
			if accessKeyIdKey == "" {
				accessKeyIdKey = "access_key_id"
			}
			secretAccessKeyKey := credentials.GetSecretAccessKeyKey()
			if secretAccessKeyKey == "" {
				secretAccessKeyKey = "secret_access_key"
			}
			envVars = append(envVars,
				secretEnvVar("AWS_ACCESS_KEY_ID", credentials.GetSecretName(), accessKeyIdKey),
				secretEnvVar("AWS_SECRET_ACCESS_KEY", credentials.GetSecretName(), secretAccessKeyKey),
			)
		}
		// Keyless arm: no env at all — the pod's ambient identity
		// (IRSA / Pod Identity) carries the credentials.
	case "gcs":
		if credentials := artifactStore.GetGcs().GetCredentialsSecret(); credentials != nil {
			secretKey := credentials.GetSecretKey()
			if secretKey == "" {
				secretKey = "credentials.json"
			}
			envVars = append(envVars, envVar("GOOGLE_APPLICATION_CREDENTIALS",
				fmt.Sprintf("%s/%s", vars.GcsCredentialsMountPath, secretKey)))
		}
		// Keyless arm: Workload Identity carries the credentials.
	case "azure_blob":
		azure := artifactStore.GetAzureBlob()
		credentials := azure.GetCredentialsSecret()
		secretKey := credentials.GetSecretKey()
		if secretKey == "" {
			secretKey = "access_key"
		}
		envVars = append(envVars,
			secretEnvVar("AZURE_STORAGE_ACCESS_KEY", credentials.GetSecretName(), secretKey))
	}

	// Extra env: plain values first, then Secret-sourced entries — both
	// the spec's escape hatches for tuning and integrations beyond the
	// typed arms.
	for name, value := range locals.Spec.GetExtraEnv() {
		envVars = append(envVars, envVar(name, value))
	}
	for name, ref := range locals.Spec.GetExtraEnvFromSecret() {
		secretKey := ref.GetSecretKey()
		if secretKey == "" {
			secretKey = "password"
		}
		envVars = append(envVars, secretEnvVar(name, ref.GetSecretName(), secretKey))
	}

	return envVars
}

// serverDeployment renders the tracking-server Deployment `<name>`.
func serverDeployment(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (*appsv1.Deployment, error) {
	spec := locals.Spec

	volumes := kubernetescorev1.VolumeArray{}
	volumeMounts := kubernetescorev1.VolumeMountArray{}

	if locals.DataPvcEnabled {
		volumes = append(volumes, &kubernetescorev1.VolumeArgs{
			Name: pulumi.String("data"),
			PersistentVolumeClaim: &kubernetescorev1.PersistentVolumeClaimVolumeSourceArgs{
				ClaimName: pulumi.String(locals.Name + vars.DataPvcSuffix),
			},
		})
		volumeMounts = append(volumeMounts, &kubernetescorev1.VolumeMountArgs{
			Name:      pulumi.String("data"),
			MountPath: pulumi.String(vars.DataMountPath),
		})
	}
	if locals.ArtifactsPvcEnabled {
		volumes = append(volumes, &kubernetescorev1.VolumeArgs{
			Name: pulumi.String("artifacts"),
			PersistentVolumeClaim: &kubernetescorev1.PersistentVolumeClaimVolumeSourceArgs{
				ClaimName: pulumi.String(locals.Name + vars.ArtifactsPvcSuffix),
			},
		})
		volumeMounts = append(volumeMounts, &kubernetescorev1.VolumeMountArgs{
			Name:      pulumi.String("artifacts"),
			MountPath: pulumi.String(vars.ArtifactsMountPath),
		})
	}
	if locals.AuthEnabled {
		volumes = append(volumes, &kubernetescorev1.VolumeArgs{
			Name: pulumi.String("auth-config"),
			Secret: &kubernetescorev1.SecretVolumeSourceArgs{
				SecretName: pulumi.String(locals.AuthConfigSecretName),
			},
		})
		volumeMounts = append(volumeMounts, &kubernetescorev1.VolumeMountArgs{
			Name:      pulumi.String("auth-config"),
			MountPath: pulumi.String(vars.AuthConfigMountPath),
			ReadOnly:  pulumi.Bool(true),
		})
	}
	if locals.ArtifactType == "gcs" {
		if credentials := spec.GetArtifactStore().GetGcs().GetCredentialsSecret(); credentials != nil {
			volumes = append(volumes, &kubernetescorev1.VolumeArgs{
				Name: pulumi.String("gcs-credentials"),
				Secret: &kubernetescorev1.SecretVolumeSourceArgs{
					SecretName: pulumi.String(credentials.GetSecretName()),
				},
			})
			volumeMounts = append(volumeMounts, &kubernetescorev1.VolumeMountArgs{
				Name:      pulumi.String("gcs-credentials"),
				MountPath: pulumi.String(vars.GcsCredentialsMountPath),
				ReadOnly:  pulumi.Bool(true),
			})
		}
	}

	container := &kubernetescorev1.ContainerArgs{
		Name:    pulumi.String("mlflow"),
		Image:   pulumi.String(locals.Image),
		Command: pulumi.ToStringArray(serverArgs(locals)),
		Env:     serverEnv(locals),
		Ports: kubernetescorev1.ContainerPortArray{
			&kubernetescorev1.ContainerPortArgs{
				Name:          pulumi.String("http"),
				ContainerPort: pulumi.Int(vars.ServerPort),
				Protocol:      pulumi.String("TCP"),
			},
		},
		// /health is the server's own unauthenticated health contract
		// (upstream's deployment reference probes it) — it answers even
		// with basic-auth on.
		LivenessProbe: &kubernetescorev1.ProbeArgs{
			HttpGet: &kubernetescorev1.HTTPGetActionArgs{
				Path: pulumi.String("/health"),
				Port: pulumi.String("http"),
			},
			InitialDelaySeconds: pulumi.Int(15),
			PeriodSeconds:       pulumi.Int(20),
		},
		ReadinessProbe: &kubernetescorev1.ProbeArgs{
			HttpGet: &kubernetescorev1.HTTPGetActionArgs{
				Path: pulumi.String("/health"),
				Port: pulumi.String("http"),
			},
			InitialDelaySeconds: pulumi.Int(5),
			PeriodSeconds:       pulumi.Int(10),
		},
	}
	if resources := resourcesArgs(spec.GetServer().GetResources()); resources != nil {
		container.Resources = resources
	}
	if len(volumeMounts) > 0 {
		container.VolumeMounts = volumeMounts
	}

	podSpec := &kubernetescorev1.PodSpecArgs{
		Containers: kubernetescorev1.ContainerArray{container},
	}
	if len(volumes) > 0 {
		podSpec.Volumes = volumes
	}
	if scheduling := spec.GetScheduling(); scheduling != nil {
		if len(scheduling.GetNodeSelector()) > 0 {
			podSpec.NodeSelector = pulumi.ToStringMap(scheduling.GetNodeSelector())
		}
		if len(scheduling.GetTolerations()) > 0 {
			podSpec.Tolerations = tolerationsArgs(scheduling.GetTolerations())
		}
	}

	// Strategy follows the volume truth (upstream's own rule): any RWO
	// PVC binds one pod, so updates must Recreate; stateless shapes
	// roll.
	strategy := &appsv1.DeploymentStrategyArgs{Type: pulumi.String("RollingUpdate")}
	if locals.DataPvcEnabled || locals.ArtifactsPvcEnabled {
		strategy = &appsv1.DeploymentStrategyArgs{Type: pulumi.String("Recreate")}
	}

	createdDeployment, err := appsv1.NewDeployment(ctx, "mlflow-deployment",
		&appsv1.DeploymentArgs{
			Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.Name),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
			}),
			Spec: &appsv1.DeploymentSpecArgs{
				Replicas: pulumi.Int(locals.Replicas),
				Strategy: strategy,
				Selector: &kubernetesmeta.LabelSelectorArgs{
					MatchLabels: pulumi.ToStringMap(locals.SelectorLabels),
				},
				Template: &kubernetescorev1.PodTemplateSpecArgs{
					Metadata: kubernetesmeta.ObjectMetaPtrInput(&kubernetesmeta.ObjectMetaArgs{
						Labels: pulumi.ToStringMap(mergeStringMaps(locals.Labels, locals.SelectorLabels)),
					}),
					Spec: podSpec,
				},
			},
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create mlflow deployment")
	}
	return createdDeployment, nil
}

// envVar renders a plain-value environment variable.
func envVar(name, value string) *kubernetescorev1.EnvVarArgs {
	return &kubernetescorev1.EnvVarArgs{
		Name:  pulumi.String(name),
		Value: pulumi.String(value),
	}
}

// secretEnvVar renders a Secret-sourced environment variable — the
// leak-free path for credential material into the container.
func secretEnvVar(name, secretName, secretKey string) *kubernetescorev1.EnvVarArgs {
	return &kubernetescorev1.EnvVarArgs{
		Name: pulumi.String(name),
		ValueFrom: &kubernetescorev1.EnvVarSourceArgs{
			SecretKeyRef: &kubernetescorev1.SecretKeySelectorArgs{
				Name: pulumi.String(secretName),
				Key:  pulumi.String(secretKey),
			},
		},
	}
}

// resourcesArgs renders ContainerResources into the k8s resources shape.
func resourcesArgs(r *kubernetesv1.ContainerResources) *kubernetescorev1.ResourceRequirementsArgs {
	if r == nil {
		return nil
	}
	requirements := &kubernetescorev1.ResourceRequirementsArgs{}
	hasAny := false
	if q := r.GetRequests(); q != nil && (q.GetCpu() != "" || q.GetMemory() != "") {
		requests := pulumi.StringMap{}
		if q.GetCpu() != "" {
			requests["cpu"] = pulumi.String(q.GetCpu())
		}
		if q.GetMemory() != "" {
			requests["memory"] = pulumi.String(q.GetMemory())
		}
		requirements.Requests = requests
		hasAny = true
	}
	if l := r.GetLimits(); l != nil && (l.GetCpu() != "" || l.GetMemory() != "") {
		limits := pulumi.StringMap{}
		if l.GetCpu() != "" {
			limits["cpu"] = pulumi.String(l.GetCpu())
		}
		if l.GetMemory() != "" {
			limits["memory"] = pulumi.String(l.GetMemory())
		}
		requirements.Limits = limits
		hasAny = true
	}
	if !hasAny {
		return nil
	}
	return requirements
}

// tolerationsArgs renders WorkloadTolerations into the k8s tolerations
// shape.
func tolerationsArgs(tolerations []*kubernetesv1.WorkloadToleration) kubernetescorev1.TolerationArray {
	out := kubernetescorev1.TolerationArray{}
	for _, t := range tolerations {
		entry := &kubernetescorev1.TolerationArgs{}
		if t.GetKey() != "" {
			entry.Key = pulumi.String(t.GetKey())
		}
		if t.GetOperator() != "" {
			entry.Operator = pulumi.String(t.GetOperator())
		}
		if t.GetValue() != "" {
			entry.Value = pulumi.String(t.GetValue())
		}
		if t.GetEffect() != "" {
			entry.Effect = pulumi.String(t.GetEffect())
		}
		if t.TolerationSeconds != nil {
			entry.TolerationSeconds = pulumi.Int(int(t.GetTolerationSeconds()))
		}
		out = append(out, entry)
	}
	return out
}

// mergeStringMaps overlays b on a into a fresh map.
func mergeStringMaps(a, b map[string]string) map[string]string {
	out := make(map[string]string, len(a)+len(b))
	for k, v := range a {
		out[k] = v
	}
	for k, v := range b {
		out[k] = v
	}
	return out
}
