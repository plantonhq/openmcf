package module

import (
	"strconv"

	kubernetesprovider "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes"
	kubernetessolrv1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetessolr/v1alpha1"
	solrv1beta1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/solroperator/kubernetes/solr/v1beta1"
	kubernetesmeta "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// createSolrCloud renders the solr.apache.org/v1beta1 SolrCloud resource
// with the typed crd2pulumi SDK (field/structure drift against the pinned
// CRD fails at compile time). Unset optionals are omitted entirely so the
// apiserver applies the CRD's own defaults — presence discipline mirrors
// the Terraform module's null-prune rendering byte for byte.
func createSolrCloud(ctx *pulumi.Context, locals *Locals,
	kubernetesProvider pulumi.ProviderResource,
	dependencies []pulumi.ResourceOption,
) (pulumi.Resource, error) {
	spec := locals.Spec

	// replicas / log level carry proto defaults (3 / INFO) so the module
	// renders the same resource whether or not the platform's defaulting
	// middleware ran — same coalescing as the Terraform variables'
	// optional() defaults.
	replicas := int32(3)
	if spec.Replicas != nil {
		replicas = spec.GetReplicas()
	}
	logLevel := spec.GetLogLevel()
	if logLevel == "" {
		logLevel = "INFO"
	}

	solrCloudSpec := solrv1beta1.SolrCloudSpecArgs{
		Replicas:           pulumi.Int(int(replicas)),
		SolrImage:          buildSolrImage(spec),
		SolrLogLevel:       pulumi.String(logLevel),
		SolrAddressability: buildSolrAddressability(spec),
	}

	if spec.GetJavaMem() != "" {
		solrCloudSpec.SolrJavaMem = pulumi.String(spec.GetJavaMem())
	}
	if spec.GetSolrOpts() != "" {
		solrCloudSpec.SolrOpts = pulumi.String(spec.GetSolrOpts())
	}
	if spec.GetGcTune() != "" {
		solrCloudSpec.SolrGCTune = pulumi.String(spec.GetGcTune())
	}

	// An empty zookeeper block renders NOTHING — the operator then
	// defaults to a provided 3-node ensemble on its own.
	if zookeeperRef := buildZookeeperRef(spec.GetZookeeper()); zookeeperRef != nil {
		solrCloudSpec.ZookeeperRef = zookeeperRef
	}

	if dataStorage := buildDataStorage(spec.GetStorage()); dataStorage != nil {
		solrCloudSpec.DataStorage = dataStorage
	}

	if podOptions := buildPodOptions(spec); podOptions != nil {
		solrCloudSpec.CustomSolrKubeOptions = solrv1beta1.SolrCloudSpecCustomSolrKubeOptionsArgs{
			PodOptions: podOptions,
		}
	}

	if updateStrategy := buildUpdateStrategy(spec.GetUpdateStrategy()); updateStrategy != nil {
		solrCloudSpec.UpdateStrategy = updateStrategy
	}

	// pdb_enabled is optional-with-default: absence already means enabled
	// to the operator, so only an EXPLICIT value renders.
	if availability := spec.GetAvailability(); availability != nil && availability.PdbEnabled != nil {
		solrCloudSpec.Availability = solrv1beta1.SolrCloudSpecAvailabilityArgs{
			PodDisruptionBudget: solrv1beta1.SolrCloudSpecAvailabilityPodDisruptionBudgetArgs{
				Enabled: pulumi.Bool(availability.GetPdbEnabled()),
			},
		}
	}

	if scaling := buildScaling(spec.GetScaling()); scaling != nil {
		solrCloudSpec.Scaling = scaling
	}

	if tls := buildSolrTls(spec.GetTls()); tls != nil {
		solrCloudSpec.SolrTLS = tls
	}

	if security := buildSolrSecurity(spec.GetSecurity()); security != nil {
		solrCloudSpec.SolrSecurity = security
	}

	if repositories := buildBackupRepositories(spec.GetBackupRepositories()); len(repositories) > 0 {
		solrCloudSpec.BackupRepositories = repositories
	}

	if len(spec.GetSolrModules()) > 0 {
		solrCloudSpec.SolrModules = pulumi.ToStringArray(spec.GetSolrModules())
	}
	if len(spec.GetAdditionalLibs()) > 0 {
		solrCloudSpec.AdditionalLibs = pulumi.ToStringArray(spec.GetAdditionalLibs())
	}

	return solrv1beta1.NewSolrCloud(ctx, locals.ClusterName,
		&solrv1beta1.SolrCloudArgs{
			Metadata: kubernetesmeta.ObjectMetaArgs{
				Name:      pulumi.String(locals.ClusterName),
				Namespace: pulumi.String(locals.Namespace),
				Labels:    pulumi.ToStringMap(locals.Labels),
				// BACKGROUND deletion, explicitly: the OPERATOR owns the
				// SolrCloud's cascade (its finalizer deletes the
				// StatefulSet, services and the provided ZooKeeper). The
				// provider's DEFAULT propagation is Foreground, which
				// DEADLOCKS the teardown — verified live: the
				// foregroundDeletion finalizer waits for the child
				// ZookeeperCluster while the zookeeper-operator keeps
				// reconciling it back to life. Terraform twin:
				// delete_cascade = "Background" on kubectl_manifest.
				Annotations: pulumi.StringMap{
					"pulumi.com/deletionPropagationPolicy": pulumi.String("background"),
				},
			},
			Spec: solrCloudSpec,
		}, append([]pulumi.ResourceOption{pulumi.Provider(kubernetesProvider)}, dependencies...)...)
}

// buildSolrImage always renders: the tag is the spec's required version and
// the repository falls back to the official "solr" image. imagePullPolicy
// is deliberately omitted (operator default).
func buildSolrImage(spec *kubernetessolrv1alpha1.KubernetesSolrSpec) solrv1beta1.SolrCloudSpecSolrImagePtrInput {
	repository := spec.GetImageRepository()
	if repository == "" {
		repository = "solr"
	}
	return solrv1beta1.SolrCloudSpecSolrImageArgs{
		Repository: pulumi.String(repository),
		Tag:        pulumi.String(spec.GetVersion()),
	}
}

func buildZookeeperRef(zookeeper *kubernetessolrv1alpha1.KubernetesSolrZookeeper) solrv1beta1.SolrCloudSpecZookeeperRefPtrInput {
	if external := zookeeper.GetExternal(); external != nil {
		return solrv1beta1.SolrCloudSpecZookeeperRefArgs{
			ConnectionInfo: solrv1beta1.SolrCloudSpecZookeeperRefConnectionInfoArgs{
				InternalConnectionString: pulumi.String(external.GetConnectionString()),
				Chroot:                   pulumi.String(chrootOrDefault(external.GetChroot())),
			},
		}
	}

	provided := zookeeper.GetProvided()
	if provided == nil {
		return nil
	}

	zkReplicas := int32(3)
	if provided.Replicas != nil {
		zkReplicas = provided.GetReplicas()
	}

	providedArgs := solrv1beta1.SolrCloudSpecZookeeperRefProvidedArgs{
		Replicas: pulumi.Int(int(zkReplicas)),
		Chroot:   pulumi.String(chrootOrDefault(provided.GetChroot())),
	}

	if persistence := provided.GetPersistence(); persistence != nil {
		persistenceSpec := solrv1beta1.SolrCloudSpecZookeeperRefProvidedPersistenceSpecArgs{}
		if persistence.GetSize() != "" {
			persistenceSpec.Resources = solrv1beta1.SolrCloudSpecZookeeperRefProvidedPersistenceSpecResourcesArgs{
				Requests: pulumi.Map{"storage": pulumi.String(persistence.GetSize())},
			}
		}
		if persistence.GetStorageClass().GetValue() != "" {
			persistenceSpec.StorageClassName = pulumi.String(persistence.GetStorageClass().GetValue())
		}
		providedArgs.Persistence = solrv1beta1.SolrCloudSpecZookeeperRefProvidedPersistenceArgs{
			Spec: persistenceSpec,
		}
	}

	if resources := provided.GetResources(); resources != nil {
		limits, requests := resourceMaps(resources)
		providedArgs.ZookeeperPodPolicy = solrv1beta1.SolrCloudSpecZookeeperRefProvidedZookeeperPodPolicyArgs{
			Resources: solrv1beta1.SolrCloudSpecZookeeperRefProvidedZookeeperPodPolicyResourcesArgs{
				Limits:   limits,
				Requests: requests,
			},
		}
	}

	return solrv1beta1.SolrCloudSpecZookeeperRefArgs{
		Provided: providedArgs,
	}
}

func buildDataStorage(storage *kubernetessolrv1alpha1.KubernetesSolrStorage) solrv1beta1.SolrCloudSpecDataStoragePtrInput {
	if persistent := storage.GetPersistent(); persistent != nil {
		reclaimPolicy := persistent.GetReclaimPolicy()
		if reclaimPolicy == "" {
			reclaimPolicy = "Retain"
		}
		pvcSpec := solrv1beta1.SolrCloudSpecDataStoragePersistentPvcTemplateSpecArgs{
			Resources: solrv1beta1.SolrCloudSpecDataStoragePersistentPvcTemplateSpecResourcesArgs{
				Requests: pulumi.Map{"storage": pulumi.String(persistent.GetSize())},
			},
		}
		if persistent.GetStorageClass().GetValue() != "" {
			pvcSpec.StorageClassName = pulumi.String(persistent.GetStorageClass().GetValue())
		}
		return solrv1beta1.SolrCloudSpecDataStorageArgs{
			Persistent: solrv1beta1.SolrCloudSpecDataStoragePersistentArgs{
				ReclaimPolicy: pulumi.String(reclaimPolicy),
				PvcTemplate: solrv1beta1.SolrCloudSpecDataStoragePersistentPvcTemplateArgs{
					Spec: pvcSpec,
				},
			},
		}
	}

	if ephemeral := storage.GetEphemeral(); ephemeral != nil {
		emptyDir := solrv1beta1.SolrCloudSpecDataStorageEphemeralEmptyDirArgs{}
		if ephemeral.GetSizeLimit() != "" {
			emptyDir.SizeLimit = pulumi.String(ephemeral.GetSizeLimit())
		}
		return solrv1beta1.SolrCloudSpecDataStorageArgs{
			Ephemeral: solrv1beta1.SolrCloudSpecDataStorageEphemeralArgs{
				EmptyDir: emptyDir,
			},
		}
	}

	return nil
}

// buildPodOptions gathers the Solr node scheduling/resource knobs into the
// operator's customSolrKubeOptions.podOptions — rendered only when any is
// set.
func buildPodOptions(spec *kubernetessolrv1alpha1.KubernetesSolrSpec) solrv1beta1.SolrCloudSpecCustomSolrKubeOptionsPodOptionsPtrInput {
	podOptions := solrv1beta1.SolrCloudSpecCustomSolrKubeOptionsPodOptionsArgs{}
	hasAny := false

	if resources := spec.GetResources(); resources != nil {
		limits, requests := resourceMaps(resources)
		podOptions.Resources = solrv1beta1.SolrCloudSpecCustomSolrKubeOptionsPodOptionsResourcesArgs{
			Limits:   limits,
			Requests: requests,
		}
		hasAny = true
	}
	if len(spec.GetNodeSelector()) > 0 {
		podOptions.NodeSelector = pulumi.ToStringMap(spec.GetNodeSelector())
		hasAny = true
	}
	if len(spec.GetTolerations()) > 0 {
		tolerations := solrv1beta1.SolrCloudSpecCustomSolrKubeOptionsPodOptionsTolerationsArray{}
		for _, toleration := range spec.GetTolerations() {
			tolerationArgs := solrv1beta1.SolrCloudSpecCustomSolrKubeOptionsPodOptionsTolerationsArgs{}
			if toleration.GetKey() != "" {
				tolerationArgs.Key = pulumi.String(toleration.GetKey())
			}
			if toleration.GetOperator() != "" {
				tolerationArgs.Operator = pulumi.String(toleration.GetOperator())
			}
			if toleration.GetValue() != "" {
				tolerationArgs.Value = pulumi.String(toleration.GetValue())
			}
			if toleration.GetEffect() != "" {
				tolerationArgs.Effect = pulumi.String(toleration.GetEffect())
			}
			if toleration.TolerationSeconds != nil {
				tolerationArgs.TolerationSeconds = pulumi.Int(int(toleration.GetTolerationSeconds()))
			}
			tolerations = append(tolerations, tolerationArgs)
		}
		podOptions.Tolerations = tolerations
		hasAny = true
	}

	if !hasAny {
		return nil
	}
	return podOptions
}

// buildSolrAddressability always renders: podPort carries a proto default
// (8983) and the external block models the operator's own Ingress /
// ExternalDNS exposure when declared.
func buildSolrAddressability(spec *kubernetessolrv1alpha1.KubernetesSolrSpec) solrv1beta1.SolrCloudSpecSolrAddressabilityPtrInput {
	podPort := int32(8983)
	if spec.PodPort != nil {
		podPort = spec.GetPodPort()
	}

	addressability := solrv1beta1.SolrCloudSpecSolrAddressabilityArgs{
		PodPort: pulumi.Int(int(podPort)),
	}

	if external := spec.GetExternal(); external != nil {
		externalArgs := solrv1beta1.SolrCloudSpecSolrAddressabilityExternalArgs{
			Method:     pulumi.String(external.GetMethod()),
			DomainName: pulumi.String(external.GetDomainName()),
		}
		if external.GetUseExternalAddress() {
			externalArgs.UseExternalAddress = pulumi.Bool(true)
		}
		if external.GetHideCommon() {
			externalArgs.HideCommon = pulumi.Bool(true)
		}
		if external.GetHideNodes() {
			externalArgs.HideNodes = pulumi.Bool(true)
		}
		addressability.External = externalArgs
	}

	return addressability
}

func buildUpdateStrategy(updateStrategy *kubernetessolrv1alpha1.KubernetesSolrUpdateStrategy) solrv1beta1.SolrCloudSpecUpdateStrategyPtrInput {
	if updateStrategy == nil {
		return nil
	}

	method := updateStrategy.GetMethod()
	if method == "" {
		method = "Managed"
	}

	updateStrategyArgs := solrv1beta1.SolrCloudSpecUpdateStrategyArgs{
		Method: pulumi.String(method),
	}

	managed := solrv1beta1.SolrCloudSpecUpdateStrategyManagedArgs{}
	hasManaged := false
	if updateStrategy.GetMaxPodsUnavailable() != "" {
		managed.MaxPodsUnavailable = intOrString(updateStrategy.GetMaxPodsUnavailable())
		hasManaged = true
	}
	if updateStrategy.GetMaxShardReplicasUnavailable() != "" {
		managed.MaxShardReplicasUnavailable = intOrString(updateStrategy.GetMaxShardReplicasUnavailable())
		hasManaged = true
	}
	if hasManaged {
		updateStrategyArgs.Managed = managed
	}

	if updateStrategy.GetRestartSchedule() != "" {
		updateStrategyArgs.RestartSchedule = pulumi.String(updateStrategy.GetRestartSchedule())
	}

	return updateStrategyArgs
}

// buildScaling renders only the EXPLICITLY set flags — both default true
// upstream, so absence already means "move replicas, don't drop them".
func buildScaling(scaling *kubernetessolrv1alpha1.KubernetesSolrScaling) solrv1beta1.SolrCloudSpecScalingPtrInput {
	if scaling == nil {
		return nil
	}

	scalingArgs := solrv1beta1.SolrCloudSpecScalingArgs{}
	hasAny := false
	if scaling.VacatePodsOnScaleDown != nil {
		scalingArgs.VacatePodsOnScaleDown = pulumi.Bool(scaling.GetVacatePodsOnScaleDown())
		hasAny = true
	}
	if scaling.PopulatePodsOnScaleUp != nil {
		scalingArgs.PopulatePodsOnScaleUp = pulumi.Bool(scaling.GetPopulatePodsOnScaleUp())
		hasAny = true
	}

	if !hasAny {
		return nil
	}
	return scalingArgs
}

func buildSolrTls(tls *kubernetessolrv1alpha1.KubernetesSolrTls) solrv1beta1.SolrCloudSpecSolrTLSPtrInput {
	if tls == nil {
		return nil
	}

	clientAuth := tls.GetClientAuth()
	if clientAuth == "" {
		clientAuth = "None"
	}

	tlsArgs := solrv1beta1.SolrCloudSpecSolrTLSArgs{
		Pkcs12Secret: solrv1beta1.SolrCloudSpecSolrTLSPkcs12SecretArgs{
			Name: pulumi.String(tls.GetPkcs12Secret().GetName()),
			Key:  pulumi.String(tls.GetPkcs12Secret().GetKey()),
		},
		KeyStorePasswordSecret: solrv1beta1.SolrCloudSpecSolrTLSKeyStorePasswordSecretArgs{
			Name: pulumi.String(tls.GetKeystorePasswordSecret().GetName()),
			Key:  pulumi.String(tls.GetKeystorePasswordSecret().GetKey()),
		},
		ClientAuth: pulumi.String(clientAuth),
	}

	if truststore := tls.GetTruststoreSecret(); truststore != nil {
		tlsArgs.TrustStoreSecret = solrv1beta1.SolrCloudSpecSolrTLSTrustStoreSecretArgs{
			Name: pulumi.String(truststore.GetName()),
			Key:  pulumi.String(truststore.GetKey()),
		}
	}
	if truststorePassword := tls.GetTruststorePasswordSecret(); truststorePassword != nil {
		tlsArgs.TrustStorePasswordSecret = solrv1beta1.SolrCloudSpecSolrTLSTrustStorePasswordSecretArgs{
			Name: pulumi.String(truststorePassword.GetName()),
			Key:  pulumi.String(truststorePassword.GetKey()),
		}
	}
	if tls.GetVerifyClientHostname() {
		tlsArgs.VerifyClientHostname = pulumi.Bool(true)
	}

	return tlsArgs
}

// buildSolrSecurity renders only when basic auth is enabled — the CRD's
// authenticationType value is capitalized ("Basic") while the spec's enum
// is lowercase. A declared-but-empty security block means security stays
// disabled and nothing renders.
func buildSolrSecurity(security *kubernetessolrv1alpha1.KubernetesSolrSecurity) solrv1beta1.SolrCloudSpecSolrSecurityPtrInput {
	if security.GetAuthenticationType() != "basic" {
		return nil
	}

	securityArgs := solrv1beta1.SolrCloudSpecSolrSecurityArgs{
		AuthenticationType: pulumi.String("Basic"),
	}

	if security.GetBasicAuthSecret().GetValue() != "" {
		securityArgs.BasicAuthSecret = pulumi.String(security.GetBasicAuthSecret().GetValue())
	}
	if security.GetProbesRequireAuth() {
		securityArgs.ProbesRequireAuth = pulumi.Bool(true)
	}
	if bootstrap := security.GetBootstrapSecurityJson(); bootstrap != nil {
		securityArgs.BootstrapSecurityJson = solrv1beta1.SolrCloudSpecSolrSecurityBootstrapSecurityJsonArgs{
			Name: pulumi.String(bootstrap.GetName()),
			Key:  pulumi.String(bootstrap.GetKey()),
		}
	}

	return securityArgs
}

func buildBackupRepositories(repositories []*kubernetessolrv1alpha1.KubernetesSolrBackupRepository) solrv1beta1.SolrCloudSpecBackupRepositoriesArray {
	repositoryArgs := solrv1beta1.SolrCloudSpecBackupRepositoriesArray{}
	for _, repository := range repositories {
		args := solrv1beta1.SolrCloudSpecBackupRepositoriesArgs{
			Name: pulumi.String(repository.GetName()),
		}

		if s3 := repository.GetS3(); s3 != nil {
			s3Args := solrv1beta1.SolrCloudSpecBackupRepositoriesS3Args{
				Region: pulumi.String(s3.GetRegion()),
				Bucket: pulumi.String(s3.GetBucket()),
			}
			if s3.GetBaseLocation() != "" {
				s3Args.BaseLocation = pulumi.String(s3.GetBaseLocation())
			}
			if s3.GetEndpoint() != "" {
				s3Args.Endpoint = pulumi.String(s3.GetEndpoint())
			}
			// Declared credentials ride secretKeyRefs; an empty credentials
			// block means the nodes' ambient identity (IRSA) — the keyless
			// path — and renders nothing.
			if credentials := s3.GetCredentials(); credentials != nil {
				credentialsArgs := solrv1beta1.SolrCloudSpecBackupRepositoriesS3CredentialsArgs{}
				hasAny := false
				if accessKey := credentials.GetAccessKeyIdSecret(); accessKey != nil {
					credentialsArgs.AccessKeyIdSecret = solrv1beta1.SolrCloudSpecBackupRepositoriesS3CredentialsAccessKeyIdSecretArgs{
						Name: pulumi.String(accessKey.GetName()),
						Key:  pulumi.String(accessKey.GetKey()),
					}
					hasAny = true
				}
				if secretKey := credentials.GetSecretAccessKeySecret(); secretKey != nil {
					credentialsArgs.SecretAccessKeySecret = solrv1beta1.SolrCloudSpecBackupRepositoriesS3CredentialsSecretAccessKeySecretArgs{
						Name: pulumi.String(secretKey.GetName()),
						Key:  pulumi.String(secretKey.GetKey()),
					}
					hasAny = true
				}
				if hasAny {
					s3Args.Credentials = credentialsArgs
				}
			}
			args.S3 = s3Args
		}

		if gcs := repository.GetGcs(); gcs != nil {
			gcsArgs := solrv1beta1.SolrCloudSpecBackupRepositoriesGcsArgs{
				Bucket: pulumi.String(gcs.GetBucket()),
			}
			if credential := gcs.GetGcsCredentialSecret(); credential != nil {
				gcsArgs.GcsCredentialSecret = solrv1beta1.SolrCloudSpecBackupRepositoriesGcsGcsCredentialSecretArgs{
					Name: pulumi.String(credential.GetName()),
					Key:  pulumi.String(credential.GetKey()),
				}
			}
			if gcs.GetBaseLocation() != "" {
				gcsArgs.BaseLocation = pulumi.String(gcs.GetBaseLocation())
			}
			args.Gcs = gcsArgs
		}

		if volume := repository.GetVolume(); volume != nil {
			volumeArgs := solrv1beta1.SolrCloudSpecBackupRepositoriesVolumeArgs{
				Source: solrv1beta1.SolrCloudSpecBackupRepositoriesVolumeSourceArgs{
					PersistentVolumeClaim: solrv1beta1.SolrCloudSpecBackupRepositoriesVolumeSourcePersistentVolumeClaimArgs{
						ClaimName: pulumi.String(volume.GetPvcClaimName()),
					},
				},
			}
			if volume.GetDirectory() != "" {
				volumeArgs.Directory = pulumi.String(volume.GetDirectory())
			}
			args.Volume = volumeArgs
		}

		repositoryArgs = append(repositoryArgs, args)
	}
	return repositoryArgs
}

// resourceMaps translates the shared ContainerResources message into the
// limits/requests quantity maps the CRD expects — nil sides stay nil so
// only declared sides render.
func resourceMaps(resources *kubernetesprovider.ContainerResources) (limits, requests pulumi.MapInput) {
	if resourceLimits := resources.GetLimits(); resourceLimits != nil {
		limits = pulumi.Map{
			"cpu":    pulumi.String(resourceLimits.GetCpu()),
			"memory": pulumi.String(resourceLimits.GetMemory()),
		}
	}
	if resourceRequests := resources.GetRequests(); resourceRequests != nil {
		requests = pulumi.Map{
			"cpu":    pulumi.String(resourceRequests.GetCpu()),
			"memory": pulumi.String(resourceRequests.GetMemory()),
		}
	}
	return limits, requests
}

// intOrString honors the CRD's int-or-string semantics for the managed
// update-strategy budgets: "2" renders as the number 2, "25%" stays a
// string — the exact same coercion the Terraform twin applies with
// try(tonumber(x), x).
func intOrString(value string) pulumi.Input {
	if number, err := strconv.Atoi(value); err == nil {
		return pulumi.Int(number)
	}
	return pulumi.String(value)
}

// chrootOrDefault applies the shared chroot default ("/") of both
// zookeeper arms.
func chrootOrDefault(chroot string) string {
	if chroot == "" {
		return "/"
	}
	return chroot
}
