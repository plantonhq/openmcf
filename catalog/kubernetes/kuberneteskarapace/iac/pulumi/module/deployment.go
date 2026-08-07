package module

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	kubernetesv1 "github.com/plantonhq/planton/catalog/kubernetes"
	appsv1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/apps/v1"
	kubernetescorev1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/core/v1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// registryDeployment renders the schema-registry Deployment
// (`<metadata.name>`): one `karapace` container running upstream's
// registry entrypoint with the KARAPACE_KARAPACE_REGISTRY role flag,
// configured entirely through KARAPACE_* environment variables (the
// pydantic settings env mechanism — env_prefix "karapace_", so config
// key X becomes KARAPACE_<uppercased X>).
func registryDeployment(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (*appsv1.Deployment, error) {
	spec := locals.Spec

	envVars := kubernetescorev1.EnvVarArray{
		// Role flags: this process serves ONLY the schema-registry API.
		envVar("KARAPACE_KARAPACE_REGISTRY", "true"),
		envVar("KARAPACE_KARAPACE_REST", "false"),
	}
	envVars = append(envVars, commonEnv(locals, locals.RegistryPort)...)

	// Registry behavior knobs (config.py: topic_name, replication_factor,
	// compatibility, group_id, master_election_strategy).
	envVars = append(envVars,
		envVar("KARAPACE_TOPIC_NAME", locals.TopicName),
	)
	if spec.GetRegistry() != nil && spec.GetRegistry().ReplicationFactor != nil {
		envVars = append(envVars,
			envVar("KARAPACE_REPLICATION_FACTOR", strconv.Itoa(int(*spec.GetRegistry().ReplicationFactor))))
	}
	envVars = append(envVars,
		envVar("KARAPACE_COMPATIBILITY", locals.Compatibility),
		envVar("KARAPACE_GROUP_ID", locals.GroupId),
		envVar("KARAPACE_MASTER_ELECTION_STRATEGY", locals.MasterElectionStrategy),
	)

	volumes, volumeMounts := kafkaTlsVolumes(locals)

	// Registry-side TLS serving. The advertised protocol must follow the
	// serving scheme: the leader coordinator publishes
	// `<advertised_protocol>://<advertised_hostname>:<port>` as the
	// master URL, and followers forward writes to it — an https server
	// advertising http would break follower forwarding.
	if spec.ServerTls != nil {
		certKey := spec.ServerTls.GetCertificate()
		if certKey == "" {
			certKey = "tls.crt"
		}
		keyKey := spec.ServerTls.GetKey()
		if keyKey == "" {
			keyKey = "tls.key"
		}
		volumes = append(volumes, &kubernetescorev1.VolumeArgs{
			Name: pulumi.String("server-tls"),
			Secret: &kubernetescorev1.SecretVolumeSourceArgs{
				SecretName: pulumi.String(spec.ServerTls.SecretName.GetValue()),
			},
		})
		volumeMounts = append(volumeMounts, &kubernetescorev1.VolumeMountArgs{
			Name:      pulumi.String("server-tls"),
			MountPath: pulumi.String(vars.ServerTlsMountPath),
			ReadOnly:  pulumi.Bool(true),
		})
		envVars = append(envVars,
			envVar("KARAPACE_ADVERTISED_PROTOCOL", "https"),
			envVar("KARAPACE_SERVER_TLS_CERTFILE", fmt.Sprintf("%s/%s", vars.ServerTlsMountPath, certKey)),
			envVar("KARAPACE_SERVER_TLS_KEYFILE", fmt.Sprintf("%s/%s", vars.ServerTlsMountPath, keyKey)),
		)
	}

	// HTTP-layer authentication (registry API only). Both mechanisms
	// leave /_health unauthenticated (config.py skip-auth paths), so the
	// probes below keep working.
	if auth := spec.GetHttpAuthentication(); auth != nil {
		if basic := auth.GetBasic(); basic != nil {
			authfileKey := basic.GetKey()
			if authfileKey == "" {
				authfileKey = "authfile.json"
			}
			volumes = append(volumes, &kubernetescorev1.VolumeArgs{
				Name: pulumi.String("authfile"),
				Secret: &kubernetescorev1.SecretVolumeSourceArgs{
					SecretName: pulumi.String(basic.SecretName.GetValue()),
				},
			})
			volumeMounts = append(volumeMounts, &kubernetescorev1.VolumeMountArgs{
				Name:      pulumi.String("authfile"),
				MountPath: pulumi.String(vars.AuthfileMountPath),
				ReadOnly:  pulumi.Bool(true),
			})
			envVars = append(envVars,
				envVar("KARAPACE_REGISTRY_AUTHFILE", fmt.Sprintf("%s/%s", vars.AuthfileMountPath, authfileKey)))
		}
		if oidc := auth.GetOidc(); oidc != nil {
			envVars = append(envVars,
				envVar("KARAPACE_SASL_OAUTHBEARER_AUTHENTICATION_ENABLED", "true"),
				envVar("KARAPACE_SASL_OAUTHBEARER_JWKS_ENDPOINT_URL", oidc.JwksEndpointUrl),
			)
			if oidc.ExpectedIssuer != "" {
				envVars = append(envVars,
					envVar("KARAPACE_SASL_OAUTHBEARER_EXPECTED_ISSUER", oidc.ExpectedIssuer))
			}
			if oidc.ExpectedAudience != "" {
				envVars = append(envVars,
					envVar("KARAPACE_SASL_OAUTHBEARER_EXPECTED_AUDIENCE", oidc.ExpectedAudience))
			}
		}
	}

	podSpec := &kubernetescorev1.PodSpecArgs{
		Containers: kubernetescorev1.ContainerArray{
			karapaceContainer(locals, vars.RegistryCommand, locals.RegistryPort, envVars, volumeMounts, spec.Resources),
		},
	}
	if len(volumes) > 0 {
		podSpec.Volumes = volumes
	}
	// Scheduling knobs are registry-scoped per the spec contract
	// ("Node selector for the registry pods"); the REST-proxy role
	// carries only replicas/port/resources.
	if len(spec.NodeSelector) > 0 {
		podSpec.NodeSelector = pulumi.ToStringMap(spec.NodeSelector)
	}
	if len(spec.Tolerations) > 0 {
		podSpec.Tolerations = buildTolerations(spec.Tolerations)
	}

	return createDeployment(ctx, locals, kubernetesProvider, dependencies,
		locals.RegistryName, locals.RegistryReplicas, locals.RegistrySelectorLabels, podSpec)
}

// restProxyDeployment renders the optional REST-proxy Deployment
// (`<metadata.name>-rest`): the SAME image with the role flags flipped
// (KARAPACE_KARAPACE_REST=true) and upstream's REST entrypoint, wired to
// the registry Service for schema lookups and to the same Kafka cluster
// (identical connection env and TLS mounts) for produce/consume.
func restProxyDeployment(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (*appsv1.Deployment, error) {
	envVars := kubernetescorev1.EnvVarArray{
		envVar("KARAPACE_KARAPACE_REST", "true"),
		envVar("KARAPACE_KARAPACE_REGISTRY", "false"),
	}
	envVars = append(envVars, commonEnv(locals, locals.RestPort)...)

	// Schema-registry wiring: the proxy resolves schemas through the
	// registry Service (scheme follows the registry's server_tls
	// posture).
	envVars = append(envVars,
		envVar("KARAPACE_REGISTRY_SCHEME", locals.Scheme),
		envVar("KARAPACE_REGISTRY_HOST", fmt.Sprintf("%s.%s.svc.cluster.local", locals.RegistryName, locals.Namespace)),
		envVar("KARAPACE_REGISTRY_PORT", strconv.Itoa(locals.RegistryPort)),
	)

	volumes, volumeMounts := kafkaTlsVolumes(locals)

	podSpec := &kubernetescorev1.PodSpecArgs{
		Containers: kubernetescorev1.ContainerArray{
			karapaceContainer(locals, vars.RestCommand, locals.RestPort, envVars, volumeMounts, locals.Spec.GetRestProxy().GetResources()),
		},
	}
	if len(volumes) > 0 {
		podSpec.Volumes = volumes
	}

	return createDeployment(ctx, locals, kubernetesProvider, dependencies,
		locals.RestName, locals.RestReplicas, locals.RestSelectorLabels, podSpec)
}

// commonEnv renders the env vars both roles share: serve address, the
// per-pod advertised identity, the Kafka connection, and SASL
// credentials.
func commonEnv(locals *Locals, port int) kubernetescorev1.EnvVarArray {
	envVars := kubernetescorev1.EnvVarArray{
		envVar("KARAPACE_HOST", "0.0.0.0"),
		envVar("KARAPACE_PORT", strconv.Itoa(port)),

		// PER-POD advertised hostname via the downward API. Upstream's
		// compose reference gives every instance its OWN identity
		// (KARAPACE_ADVERTISED_HOSTNAME = that container's hostname, one
		// per replica) — never a shared name: the leader publishes
		// `advertised_protocol://advertised_hostname:port` through the
		// consumer group and followers forward writes to it, so a shared
		// (Service) name would make followers forward to themselves.
		// config.py falls back to `host` when unset, which is 0.0.0.0
		// here — so it must be set explicitly. The POD IP is each pod's
		// resolvable self-identity — the Kubernetes twin of compose's
		// container hostname (a Deployment pod's bare NAME does not
		// resolve in cluster DNS, so a follower forwarding writes to
		// the leader by name would fail; the IP is directly reachable
		// pod-to-pod).
		&kubernetescorev1.EnvVarArgs{
			Name: pulumi.String("KARAPACE_ADVERTISED_HOSTNAME"),
			ValueFrom: &kubernetescorev1.EnvVarSourceArgs{
				FieldRef: &kubernetescorev1.ObjectFieldSelectorArgs{
					FieldPath: pulumi.String("status.podIP"),
				},
			},
		},

		envVar("KARAPACE_BOOTSTRAP_URI", locals.Spec.Kafka.BootstrapServers.GetValue()),
		envVar("KARAPACE_SECURITY_PROTOCOL", locals.SecurityProtocol),
		envVar("KARAPACE_LOG_LEVEL", locals.LogLevel),
	}

	// Kafka TLS file paths point into the Secret mounts rendered by
	// kafkaTlsVolumes (config.py: ssl_cafile / ssl_certfile /
	// ssl_keyfile).
	if tls := locals.Spec.GetKafka().GetTls(); tls != nil {
		caKey := tls.GetCaCertificate()
		if caKey == "" {
			caKey = "ca.crt"
		}
		envVars = append(envVars,
			envVar("KARAPACE_SSL_CAFILE", fmt.Sprintf("%s/%s", vars.KafkaCaMountPath, caKey)))

		if tls.ClientCertSecretName.GetValue() != "" {
			certKey := tls.GetClientCertificate()
			if certKey == "" {
				certKey = "user.crt"
			}
			keyKey := tls.GetClientKey()
			if keyKey == "" {
				keyKey = "user.key"
			}
			envVars = append(envVars,
				envVar("KARAPACE_SSL_CERTFILE", fmt.Sprintf("%s/%s", vars.KafkaClientCertMountPath, certKey)),
				envVar("KARAPACE_SSL_KEYFILE", fmt.Sprintf("%s/%s", vars.KafkaClientCertMountPath, keyKey)),
			)
		}
	}

	// SASL credentials. The password ALWAYS arrives via secretKeyRef —
	// either the referenced existing Secret or the module-materialized
	// one — never as a plaintext env value (see saslPasswordSecret).
	if sasl := locals.Spec.GetKafka().GetSasl(); sasl != nil {
		envVars = append(envVars,
			envVar("KARAPACE_SASL_MECHANISM", sasl.Mechanism),
			envVar("KARAPACE_SASL_PLAIN_USERNAME", sasl.Username),
			&kubernetescorev1.EnvVarArgs{
				Name: pulumi.String("KARAPACE_SASL_PLAIN_PASSWORD"),
				ValueFrom: &kubernetescorev1.EnvVarSourceArgs{
					SecretKeyRef: &kubernetescorev1.SecretKeySelectorArgs{
						Name: pulumi.String(locals.SaslPasswordSecretName),
						Key:  pulumi.String(locals.SaslPasswordSecretKey),
					},
				},
			},
		)
	}

	return envVars
}

// kafkaTlsVolumes renders the Kafka-side TLS Secret mounts shared by both
// roles: the CA to trust at /etc/karapace/kafka-ca and, for mutual-TLS
// listeners, the client identity at /etc/karapace/kafka-cert.
func kafkaTlsVolumes(locals *Locals) (kubernetescorev1.VolumeArray, kubernetescorev1.VolumeMountArray) {
	volumes := kubernetescorev1.VolumeArray{}
	mounts := kubernetescorev1.VolumeMountArray{}

	tls := locals.Spec.GetKafka().GetTls()
	if tls == nil {
		return volumes, mounts
	}

	volumes = append(volumes, &kubernetescorev1.VolumeArgs{
		Name: pulumi.String("kafka-ca"),
		Secret: &kubernetescorev1.SecretVolumeSourceArgs{
			SecretName: pulumi.String(tls.CaSecretName.GetValue()),
		},
	})
	mounts = append(mounts, &kubernetescorev1.VolumeMountArgs{
		Name:      pulumi.String("kafka-ca"),
		MountPath: pulumi.String(vars.KafkaCaMountPath),
		ReadOnly:  pulumi.Bool(true),
	})

	if tls.ClientCertSecretName.GetValue() != "" {
		volumes = append(volumes, &kubernetescorev1.VolumeArgs{
			Name: pulumi.String("kafka-cert"),
			Secret: &kubernetescorev1.SecretVolumeSourceArgs{
				SecretName: pulumi.String(tls.ClientCertSecretName.GetValue()),
			},
		})
		mounts = append(mounts, &kubernetescorev1.VolumeMountArgs{
			Name:      pulumi.String("kafka-cert"),
			MountPath: pulumi.String(vars.KafkaClientCertMountPath),
			ReadOnly:  pulumi.Bool(true),
		})
	}

	return volumes, mounts
}

// karapaceContainer assembles the single `karapace` container for a role.
//
// PROBES ON /_health: the engine's health endpoint is in config.py's
// skip-auth path list, so it answers without credentials even when HTTP
// authentication is enabled — the one path that is always safe to probe.
// The upstream image's own Docker HEALTHCHECK curls the same path. The
// probe scheme follows server_tls (kubelet probes hit the pod directly,
// so the probe must speak whatever the server serves; HTTPS probes skip
// certificate verification, which is what makes cert-manager-issued
// certs with Service SANs work).
func karapaceContainer(locals *Locals, command []string, port int,
	envVars kubernetescorev1.EnvVarArray,
	volumeMounts kubernetescorev1.VolumeMountArray,
	resources *kubernetesv1.ContainerResources,
) *kubernetescorev1.ContainerArgs {
	probeScheme := "HTTP"
	// server_tls covers the registry role only; the REST proxy always
	// serves plain HTTP on its own port.
	if locals.Spec.ServerTls != nil && port == locals.RegistryPort {
		probeScheme = "HTTPS"
	}

	probe := &kubernetescorev1.ProbeArgs{
		HttpGet: &kubernetescorev1.HTTPGetActionArgs{
			Path:   pulumi.String(vars.HealthCheckPath),
			Port:   pulumi.Int(port),
			Scheme: pulumi.String(probeScheme),
		},
		// Give the engine time to replay the schemas topic before
		// liveness can restart it.
		InitialDelaySeconds: pulumi.Int(10),
		PeriodSeconds:       pulumi.Int(10),
		TimeoutSeconds:      pulumi.Int(5),
		FailureThreshold:    pulumi.Int(6),
	}

	container := &kubernetescorev1.ContainerArgs{
		Name:    pulumi.String("karapace"),
		Image:   pulumi.String(locals.Image),
		Command: pulumi.ToStringArray(command),
		Ports: kubernetescorev1.ContainerPortArray{
			&kubernetescorev1.ContainerPortArgs{
				Name:          pulumi.String("http"),
				ContainerPort: pulumi.Int(port),
			},
		},
		Env:            envVars,
		ReadinessProbe: probe,
		LivenessProbe:  probe,
	}
	if len(volumeMounts) > 0 {
		container.VolumeMounts = volumeMounts
	}
	if resourceArgs := containerResources(resources); resourceArgs != nil {
		container.Resources = resourceArgs
	}

	return container
}

// createDeployment renders one apps/v1 Deployment with the module's
// shared identity conventions (selector = immutable per-role identity,
// object and pod labels = governance labels + per-role app label).
func createDeployment(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
	name string, replicas int, selectorLabels map[string]string,
	podSpec *kubernetescorev1.PodSpecArgs,
) (*appsv1.Deployment, error) {
	labels := mergedLabels(locals, selectorLabels)

	opts := append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)
	createdDeployment, err := appsv1.NewDeployment(ctx,
		name,
		&appsv1.DeploymentArgs{
			Metadata: &kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(name),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(labels),
			},
			Spec: &appsv1.DeploymentSpecArgs{
				Replicas: pulumi.Int(replicas),
				Selector: &kubernetesmeta.LabelSelectorArgs{
					MatchLabels: pulumi.ToStringMap(selectorLabels),
				},
				Template: &kubernetescorev1.PodTemplateSpecArgs{
					Metadata: &kubernetesmeta.ObjectMetaArgs{
						Labels: pulumi.ToStringMap(labels),
					},
					Spec: podSpec,
				},
			},
		}, opts...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create %s deployment", name)
	}

	return createdDeployment, nil
}

// envVar renders one literal-valued env var.
func envVar(name, value string) *kubernetescorev1.EnvVarArgs {
	return &kubernetescorev1.EnvVarArgs{
		Name:  pulumi.String(name),
		Value: pulumi.String(value),
	}
}

// containerResources converts the shared ContainerResources shape into
// Pulumi resource requirements; nil passes through so Kubernetes applies
// no requests/limits.
func containerResources(resources *kubernetesv1.ContainerResources) *kubernetescorev1.ResourceRequirementsArgs {
	if resources == nil {
		return nil
	}

	args := &kubernetescorev1.ResourceRequirementsArgs{}
	if limits := cpuMemoryMap(resources.Limits); len(limits) > 0 {
		args.Limits = pulumi.ToStringMap(limits)
	}
	if requests := cpuMemoryMap(resources.Requests); len(requests) > 0 {
		args.Requests = pulumi.ToStringMap(requests)
	}
	return args
}

func cpuMemoryMap(cpuMemory *kubernetesv1.CpuMemory) map[string]string {
	if cpuMemory == nil {
		return nil
	}
	result := map[string]string{}
	if cpuMemory.Cpu != "" {
		result["cpu"] = cpuMemory.Cpu
	}
	if cpuMemory.Memory != "" {
		result["memory"] = cpuMemory.Memory
	}
	return result
}

// buildTolerations converts the shared WorkloadToleration shape,
// rendering only the fields the spec sets.
func buildTolerations(tolerations []*kubernetesv1.WorkloadToleration) kubernetescorev1.TolerationArray {
	result := make(kubernetescorev1.TolerationArray, 0, len(tolerations))
	for _, t := range tolerations {
		tolArgs := &kubernetescorev1.TolerationArgs{}
		if t.Key != "" {
			tolArgs.Key = pulumi.String(t.Key)
		}
		if t.Operator != "" {
			tolArgs.Operator = pulumi.String(t.Operator)
		}
		if t.Value != "" {
			tolArgs.Value = pulumi.String(t.Value)
		}
		if t.Effect != "" {
			tolArgs.Effect = pulumi.String(t.Effect)
		}
		if t.TolerationSeconds != nil {
			tolArgs.TolerationSeconds = pulumi.Int(int(*t.TolerationSeconds))
		}
		result = append(result, tolArgs)
	}
	return result
}
