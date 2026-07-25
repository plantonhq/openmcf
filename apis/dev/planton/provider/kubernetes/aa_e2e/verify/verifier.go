package verify

import (
	"context"
	"strings"
)

// ResourceVerifier knows how to verify a specific Kubernetes resource type.
type ResourceVerifier interface {
	VerifyExists(ctx context.Context, kubeconfig string) error
	VerifyAbsent(ctx context.Context, kubeconfig string) error
}

// operatorKinds lists manifest kind values (lowercased) for operator/controller
// components. Operators install CRD controllers that watch resources but typically
// do not expose a Kubernetes Service. Verification checks namespace + running
// pods only (no service requirement).
var operatorKinds = map[string]bool{
	"kubernetesstrimzikafkaoperator": true,
	// Tier 4 operators with configurable namespace
	"kubernetesgharunnerscalesetcontroller": true,
}

// helmTier2Kinds lists manifest kind values (lowercased) for Helm-based
// Kubernetes components that deploy applications with Services.
// These must match the CloudResourceKind enum names from cloud_resource_kind.proto
// (case-insensitive via lowercasing).
var helmTier2Kinds = map[string]bool{
	// Tier 2 Helm applications
	"kubernetesnats":     true,
	"kubernetesopenbao":  true,
	"kubernetesopenfga":  true,
	"kubernetesjenkins":  true,
	"kubernetestemporal": true,
	"kubernetesargocd":   true,
	"kubernetesharbor":   true,
	"kuberneteslocust":   true,
	"kuberneteskeycloak": true,
	// Data-tier Helm applications with dedicated behavioral verifiers
	"kubernetesseaweedfs": true,
	"kubernetesqdrant":    true,
}

// crdInstallKinds maps manifest kind values (lowercased) to their expected CRD
// names for components that only install cluster-scoped CRDs without deploying
// any pods or services.
var crdInstallKinds = map[string][]string{
	// The standard channel serves all of these from the v1.6 release onward
	// (TCPRoute and UDPRoute graduated to GA; ListenerSet is standard since
	// v1.5) — assert the full set the catalog's projection kinds deploy onto.
	"kubernetesgatewayapicrds": {
		"gatewayclasses.gateway.networking.k8s.io",
		"gateways.gateway.networking.k8s.io",
		"listenersets.gateway.networking.k8s.io",
		"httproutes.gateway.networking.k8s.io",
		"grpcroutes.gateway.networking.k8s.io",
		"tcproutes.gateway.networking.k8s.io",
		"udproutes.gateway.networking.k8s.io",
		"tlsroutes.gateway.networking.k8s.io",
		"referencegrants.gateway.networking.k8s.io",
		"backendtlspolicies.gateway.networking.k8s.io",
	},
	// KubernetesIstioBaseCrds installs the istio/base CRD bundle (no istiod). Verify the
	// CRDs backing the seven typed Istio components are present.
	"kubernetesistiobasecrds": {
		"destinationrules.networking.istio.io",
		"serviceentries.networking.istio.io",
		"envoyfilters.networking.istio.io",
		"peerauthentications.security.istio.io",
		"requestauthentications.security.istio.io",
		"authorizationpolicies.security.istio.io",
		"telemetries.telemetry.istio.io",
	},
}

// gatewayApiCustomResource describes how to verify a Gateway API custom resource
// created by one of the Gateway API deployment components. These components do
// not run pods; verification confirms the CR itself exists after apply and is
// gone after destroy. The CRDs are installed by the KubernetesGatewayApiCrds
// registry prerequisite before the component applies.
type gatewayApiCustomResource struct {
	// resource is the fully-qualified kubectl resource (plural.group), which is
	// stable across the served apiVersion.
	resource string
	// clusterScoped is true for cluster-scoped kinds (GatewayClass), which must
	// be queried without a namespace.
	clusterScoped bool
}

// gatewayApiKinds maps manifest kind values (lowercased) to their Gateway API
// custom-resource verification descriptor.
var gatewayApiKinds = map[string]gatewayApiCustomResource{
	"kubernetesgatewayclass":     {resource: "gatewayclasses.gateway.networking.k8s.io", clusterScoped: true},
	"kubernetesgateway":          {resource: "gateways.gateway.networking.k8s.io"},
	"kuberneteslistenerset":      {resource: "listenersets.gateway.networking.k8s.io"},
	"kuberneteshttproute":        {resource: "httproutes.gateway.networking.k8s.io"},
	"kubernetesgrpcroute":        {resource: "grpcroutes.gateway.networking.k8s.io"},
	"kubernetestcproute":         {resource: "tcproutes.gateway.networking.k8s.io"},
	"kubernetesudproute":         {resource: "udproutes.gateway.networking.k8s.io"},
	"kubernetestlsroute":         {resource: "tlsroutes.gateway.networking.k8s.io"},
	"kubernetesreferencegrant":   {resource: "referencegrants.gateway.networking.k8s.io"},
	"kubernetesbackendtlspolicy": {resource: "backendtlspolicies.gateway.networking.k8s.io"},
}

// istioApiKinds maps manifest kind values (lowercased) to the fully-qualified
// kubectl resource (plural.group) for the typed Istio API deployment components
// (853-859). Like the Gateway API kinds, these components do not run pods;
// verification confirms the CR itself exists after apply and is gone after
// destroy. The Istio CRDs are installed by the KubernetesIstioBaseCrds registry
// prerequisite before the component applies. All seven Istio kinds are
// namespaced, so no clusterScoped flag is needed.
var istioApiKinds = map[string]string{
	"kubernetespeerauthentication":    "peerauthentications.security.istio.io",
	"kubernetesrequestauthentication": "requestauthentications.security.istio.io",
	"kubernetesauthorizationpolicy":   "authorizationpolicies.security.istio.io",
	"kubernetesserviceentry":          "serviceentries.networking.istio.io",
	"kubernetesdestinationrule":       "destinationrules.networking.istio.io",
	"kubernetesenvoyfilter":           "envoyfilters.networking.istio.io",
	"kubernetestelemetry":             "telemetries.telemetry.istio.io",
}

// GetVerifierFromManifest creates the appropriate verifier by parsing the manifest.
func GetVerifierFromManifest(manifestPath string) (ResourceVerifier, error) {
	info, err := ParseManifestInfo(manifestPath)
	if err != nil {
		return nil, err
	}

	component := strings.ToLower(info.Kind)

	switch component {
	case "kubernetesnamespace":
		return &NamespaceVerifier{Name: info.Name}, nil

	case "kubernetesdeployment":
		return &WorkloadVerifier{
			Namespace: info.Namespace,
			Kind:      "deployment",
			Name:      info.Name,
		}, nil

	case "kubernetesstatefulset":
		return &WorkloadVerifier{
			Namespace: info.Namespace,
			Kind:      "statefulset",
			Name:      info.Name,
		}, nil

	case "kubernetessecret":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "secret",
			Name:      info.Name,
		}, nil

	case "kubernetesconfigmap":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "configmap",
			Name:      info.Name,
		}, nil

	case "kubernetesserviceaccount":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "serviceaccount",
			Name:      info.Name,
		}, nil

	// The RBAC grant deploys role/binding objects whose names derive from the
	// grant's shape, so its verifier re-reads the manifest itself.
	case "kubernetesrbac":
		return &RbacVerifier{ManifestPath: manifestPath}, nil

	// Existence is the every-lane contract; the aws-load-balancer scenario
	// (real-cluster profile) additionally asserts the cloud populated a real
	// LB address — the never-wait-at-deploy posture makes the verifier the
	// only honest place for that proof.
	case "kubernetesservice":
		if strings.Contains(manifestPath, "aws-load-balancer") {
			return &ServiceLbAddressVerifier{
				Namespace: info.Namespace,
				Name:      info.Name,
			}, nil
		}
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "service",
			Name:      info.Name,
		}, nil

	// An Ingress is a valid API object with or without a controller running in
	// the cluster, so existence (not load-balancer status) is the correct
	// verification on every lane.
	case "kubernetesingress":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "ingress",
			Name:      info.Name,
		}, nil

	// A NetworkPolicy has no runtime status of its own (enforcement lives in
	// the CNI); on the default cluster the object's existence is the
	// verifiable contract. Scenarios that run with an ENFORCING CNI (the
	// KubernetesCilium prerequisite, on the cilium-cni cluster profile) get
	// the behavioral verifier — traffic actually blocked while the policy
	// exists and flowing again after destroy.
	case "kubernetesnetworkpolicy":
		if manifestHasPrerequisite(manifestPath, "KubernetesCilium") {
			return &NetworkPolicyDenyBehavioralVerifier{
				Namespace:        info.Namespace,
				PolicyName:       info.Name,
				ClientDeployment: "e2e-netpol-client",
				BackendURL:       "http://e2e-netpol-backend." + info.Namespace + ".svc.cluster.local/",
			}, nil
		}
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "networkpolicy",
			Name:      info.Name,
		}, nil

	// Cilium: the dataplane install proof — agent rolled out, operator
	// available, and every node Ready (on the CNI-less profile the
	// NotReady→Ready transition is the proof Cilium became the CNI).
	case "kubernetescilium":
		return &CiliumInstallVerifier{
			Namespace: info.Namespace,
		}, nil

	// KEDA: install proof up to the external-metrics APIService; the
	// behavioral-scaling scenario (recognized by its scale-target fixture)
	// additionally proves a real cron-driven scale-up.
	case "kuberneteskeda":
		return &KedaInstallVerifier{
			Namespace:  info.Namespace,
			Behavioral: manifestHasPrerequisiteSuffix(manifestPath, "fixture-scale-target.yaml"),
			// The aws-sqs-irsa scenario (real-cluster profile) proves the
			// cloud pod-identity hop: real SQS depth drives a real scale.
			BehavioralSqs: strings.Contains(manifestPath, "aws-sqs-irsa"),
		}, nil

	// Cluster Autoscaler: install proof down to the reconcile loop's own
	// heartbeat (the status ConfigMap). The kind lane runs the KWOK
	// simulation arm; the aws-asg-scaling scenario (real-cluster profile)
	// additionally proves a real node-group scale-up/scale-down.
	case "kubernetesclusterautoscaler":
		return &ClusterAutoscalerInstallVerifier{
			Namespace:  info.Namespace,
			Behavioral: strings.Contains(manifestPath, "aws-asg-scaling"),
		}, nil

	// Velero: install proof down to the BackupStorageLocation handshake;
	// the behavioral scenario (recognized by name) additionally proves a
	// full backup → namespace-loss → restore cycle against the MinIO
	// fixture, with verifier-owned Backup/Restore CRs.
	case "kubernetesvelero":
		return &VeleroInstallVerifier{
			Namespace:  info.Namespace,
			Behavioral: strings.Contains(manifestPath, "behavioral-backup-restore"),
			// The aws-s3-csi scenario (real-cluster profile) puts a VOLUME
			// in the DR blast radius: CSI snapshot → real S3 → restore.
			CsiVolume: strings.Contains(manifestPath, "aws-s3-csi"),
		}, nil

	// Karpenter: the controller cannot start off AWS (region/credentials/
	// instance-type discovery are fatal at startup), so these run only on
	// the batched EKS real-cluster lanes — the profiles are deferred and
	// the kind entrypoints skip them. The controller install is verified
	// like other operators (namespace + running pods; Karpenter exposes a
	// metrics Service, so the generic Helm shape fits).
	case "kuberneteskarpenter":
		return &HelmComponentVerifier{
			Namespace:     info.Namespace,
			ComponentName: "karpenter",
		}, nil

	// Karpenter fleet CRs are cluster-scoped objects reconciled by the
	// controller; the object's presence is the verifiable contract (node
	// provisioning itself is the EKS lane's behavioral assertion).
	case "kuberneteskarpenternodepool":
		return &KarpenterNodePoolVerifier{
			Name:       info.Name,
			Behavioral: strings.Contains(manifestPath, "behavioral-node-launch"),
		}, nil
	case "kuberneteskarpenterec2nodeclass":
		return &ResourceExistenceVerifier{
			Kind: "ec2nodeclasses.karpenter.k8s.aws",
			Name: info.Name,
		}, nil

	// A claim under a WaitForFirstConsumer StorageClass is correctly Pending
	// until a pod consumes it, so existence (not Bound) is the verifiable
	// contract; the composed consumer scenario proves real binding.
	// PVC: existence on every lane; the data-source scenarios (real-cluster
	// profile, snapshot-capable CSI) prove clone/snapshot-restore
	// PROVISIONING — marker data must travel from the source volume to the
	// claim under test.
	case "kubernetespersistentvolumeclaim":
		if strings.Contains(manifestPath, "snapshot-restore") || strings.Contains(manifestPath, "clone") {
			return &PvcDataSourceVerifier{
				Name:     info.Name,
				Snapshot: strings.Contains(manifestPath, "snapshot-restore"),
			}, nil
		}
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "pvc",
			Name:      info.Name,
		}, nil

	// Cluster-scoped: kubectl ignores the empty namespace argument.
	case "kubernetesstorageclass":
		return &ResourceExistenceVerifier{
			Kind: "storageclass",
			Name: info.Name,
		}, nil

	// The quota object is the verifiable contract; the companion LimitRange
	// shares its name and is asserted by the with-limit-defaults scenario's
	// import round-trip (both resources ride the same state).
	case "kubernetesresourcequota":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "resourcequota",
			Name:      info.Name,
		}, nil

	// Cluster-scoped; preemption BEHAVIOR needs scheduling pressure and
	// rides the real-cluster lanes — the object is the verifiable contract.
	// PriorityClass: existence on every lane; the preemption scenario
	// (real-cluster profile) proves the class actually evicts lower
	// priorities under genuine scheduling pressure.
	case "kubernetespriorityclass":
		if strings.Contains(manifestPath, "preemption") {
			return &PriorityClassPreemptionVerifier{Name: info.Name}, nil
		}
		return &ResourceExistenceVerifier{
			Kind: "priorityclass",
			Name: info.Name,
		}, nil

	case "kubernetespoddisruptionbudget":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "pdb",
			Name:      info.Name,
		}, nil

	// Scaling BEHAVIOR needs a metrics source. Scenarios that install
	// metrics-server (declared in the e2e-prerequisites annotation) get the
	// behavioral verifier — ScalingActive + an actual scale-up; scenarios
	// without it keep the object as the verifiable contract.
	case "kuberneteshorizontalpodautoscaler":
		if manifestHasPrerequisite(manifestPath, "KubernetesMetricsServer") {
			return &HpaBehavioralVerifier{
				Namespace:   info.Namespace,
				Name:        info.Name,
				MinReplicas: manifestSpecInt(manifestPath, "minReplicas", 1),
			}, nil
		}
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "hpa",
			Name:      info.Name,
		}, nil

	case "kubernetescronjob":
		return &ResourceExistenceVerifier{
			Namespace: info.Namespace,
			Kind:      "cronjob",
			Name:      info.Name,
		}, nil

	case "kubernetesjob":
		return &JobVerifier{
			Namespace: info.Namespace,
			Name:      info.Name,
		}, nil

	case "kubernetesdaemonset":
		return &WorkloadVerifier{
			Namespace: info.Namespace,
			Kind:      "daemonset",
			Name:      info.Name,
		}, nil

	// CloudNativePG operator: Deployment Available + CRDs Established;
	// when the manifest enables the Barman Cloud plugin, the plugin
	// deployment and its ObjectStore CRD join the contract.
	case "kubernetescloudnativepgoperator":
		return &CnpgOperatorInstallVerifier{
			Namespace:     info.Namespace,
			PluginEnabled: manifestBarmanPluginEnabled(manifestPath),
		}, nil

	// Percona operator installs: the operator Deployment Available plus
	// the CRDs the database kinds render against, Established.
	case "kubernetesperconamongooperator":
		return &PsmdbOperatorInstallVerifier{
			Namespace:   info.Namespace,
			ReleaseName: info.Name,
		}, nil
	case "kubernetesperconamysqloperator":
		return &PxcOperatorInstallVerifier{
			Namespace:    info.Namespace,
			ReleaseName:  info.Name,
			WatchWidened: pxcWatchWidened(manifestSpecMap(manifestPath)),
		}, nil

	// A Percona-operator-managed MongoDB cluster: the resource in state
	// ready + the replica-set Service. The behavioral-failover scenario
	// (recognized by name) proves data durability through a live primary
	// loss and election; the with-backup scenario drives a real PBM
	// backup to completion.
	case "kubernetesmongodb":
		spec := manifestSpecMap(manifestPath)
		rsName, rsSize := mongodbFirstReplset(spec)
		return &PsmdbClusterVerifier{
			Namespace:     info.Namespace,
			ClusterName:   info.Name,
			ReplsetName:   rsName,
			Size:          rsSize,
			Behavioral:    strings.Contains(manifestPath, "behavioral-failover"),
			BackupProof:   strings.Contains(manifestPath, "with-backup"),
			BackupStorage: mongodbFirstBackupStorage(spec),
		}, nil

	// A Strimzi-operator-managed KRaft Kafka cluster: the Kafka resource
	// Ready + every node pod + the bootstrap Service. The
	// behavioral-durability scenario (recognized by name) proves
	// replicated durability through a live broker loss (produce with
	// acks=all → kill a broker → consume during the outage → full
	// recovery).
	case "kuberneteskafka":
		spec := manifestSpecMap(manifestPath)
		firstPool, totalNodes := kafkaFirstPool(spec)
		return &KafkaClusterVerifier{
			Namespace:     info.Namespace,
			ClusterName:   info.Name,
			FirstPoolName: firstPool,
			TotalNodes:    totalNodes,
			BootstrapPort: kafkaFirstInternalPort(spec),
			Durability:    strings.Contains(manifestPath, "behavioral-durability"),
		}, nil

	// A declared Kafka topic: the KafkaTopic resource Ready (the topic
	// operator reconciled it) AND kafka-topics --describe on the target
	// cluster reporting the declared partition count — the reconcile
	// proof runs on every lane.
	case "kuberneteskafkatopic":
		spec := manifestSpecMap(manifestPath)
		topicName, _ := spec["topicName"].(string)
		if topicName == "" {
			// Scenario manifests are written in the proto's snake_case form.
			topicName, _ = spec["topic_name"].(string)
		}
		if topicName == "" {
			topicName = info.Name
		}
		return &KafkaTopicVerifier{
			Namespace:   info.Namespace,
			CrName:      info.Name,
			TopicName:   topicName,
			ClusterName: kafkaClusterRef(spec),
			Partitions:  int(manifestSpecInt(manifestPath, "partitions", 0)),
		}, nil

	// A declared Kafka user: the KafkaUser resource Ready + the
	// operator-generated credentials Secret. The behavioral-auth
	// scenario (recognized by name) additionally produces and consumes
	// AS THE USER through the cluster's scram listener — authentication
	// and ACL authorization proven on the wire.
	case "kuberneteskafkauser":
		spec := manifestSpecMap(manifestPath)
		aclTopic, aclGroup := kafkaUserAcl(spec)
		return &KafkaUserVerifier{
			Namespace:   info.Namespace,
			Username:    info.Name,
			ClusterName: kafkaClusterRef(spec),
			AuthProof:   strings.Contains(manifestPath, "behavioral-auth"),
			ScramPort:   9094,
			AclTopic:    aclTopic,
			AclGroup:    aclGroup,
		}, nil

	// A Strimzi-operator-managed Kafka Connect cluster: the KafkaConnect
	// resource Ready + the Connect REST API Service.
	case "kuberneteskafkaconnect":
		return &KafkaConnectVerifier{
			Namespace:   info.Namespace,
			ConnectName: info.Name,
		}, nil

	// A declared connector on a Connect cluster: the KafkaConnector
	// resource Ready. The behavioral-dataflow scenario (recognized by
	// name) additionally consumes the connector's output topic on the
	// Kafka fixture cluster and asserts real source content arrived —
	// the pipe proven end-to-end.
	case "kuberneteskafkaconnector":
		return &KafkaConnectorVerifier{
			Namespace:        info.Namespace,
			ConnectorName:    info.Name,
			DataFlow:         strings.Contains(manifestPath, "behavioral-dataflow"),
			KafkaClusterName: kafkaEcosystemFixtureCluster,
			SourceTopic:      connectorDataFlowSourceTopic,
			MirroredTopic:    connectorDataFlowMirroredTopic,
		}, nil

	// A MirrorMaker 2 deployment: the KafkaMirrorMaker2 resource Ready.
	// The behavioral-migration scenario (recognized by name) produces on
	// the SOURCE fixture cluster and consumes from the TARGET fixture
	// cluster — records crossing clusters is the migration proof.
	case "kuberneteskafkamirrormaker2":
		spec := manifestSpecMap(manifestPath)
		return &KafkaMirrorMaker2Verifier{
			Namespace:       info.Namespace,
			MirrorMakerName: info.Name,
			Migration:       strings.Contains(manifestPath, "behavioral-migration"),
			SourceCluster:   kafkaEcosystemSourceCluster,
			TargetCluster:   kafkaEcosystemFixtureCluster,
			SourceAlias:     mirrorMaker2FirstSourceAlias(spec),
		}, nil

	// A Karapace schema registry: registry Deployment Available + Service
	// + a LIVE register/fetch round-trip through the Confluent-compatible
	// API on every lane (a registry that cannot persist schemas is not a
	// registry).
	case "kuberneteskarapace":
		spec := manifestSpecMap(manifestPath)
		return &KarapaceVerifier{
			Namespace:    info.Namespace,
			RegistryName: info.Name,
			Port:         int(manifestSpecInt(manifestPath, "port", 8081)),
			RestProxy:    karapaceRestProxyEnabled(spec),
		}, nil

	// A kafbat UI console: Deployment Available + Service + the console's
	// own API answering. The behavioral-observe scenario additionally
	// asserts the declared cluster is reported ONLINE; with-auth
	// scenarios assert anonymous access is refused.
	case "kuberneteskafkaui":
		spec := manifestSpecMap(manifestPath)
		clusterOnline := ""
		if strings.Contains(manifestPath, "behavioral-observe") {
			clusterOnline = kafkaUiFirstClusterName(spec)
		}
		return &KafkaUiVerifier{
			Namespace:     info.Namespace,
			ServiceName:   info.Name,
			Port:          int(manifestSpecInt(manifestPath, "servicePort", 80)),
			ClusterOnline: clusterOnline,
			Authenticated: strings.Contains(manifestPath, "with-auth"),
		}, nil

	// The Altinity ClickHouse operator: the operator Deployment
	// Available + the four chart-owned CRDs Established (kept on
	// destroy by Helm's crds/ semantics).
	case "kubernetesaltinityoperator":
		return &AltinityOperatorInstallVerifier{
			Namespace:   info.Namespace,
			ReleaseName: info.Name,
		}, nil

	// The RabbitMQ Cluster Operator: Deployment Available + the
	// RabbitmqCluster CRD Established + both admission webhook
	// configurations present (their serving cert is cert-manager-issued
	// — the operator's registry prerequisite). Every name is fixed by
	// the release manifest, including the rabbitmq-system namespace; on
	// destroy the CRD is asserted GONE (it is a document of the applied
	// manifest, not a kept CRD).
	case "kubernetesrabbitmqoperator":
		return &RabbitMqOperatorInstallVerifier{}, nil

	// An operator-managed RabbitMQ cluster: ClusterAvailable +
	// AllReplicasReady conditions, the naming contract's Services and
	// default-user Secret, plus a LIVE message round-trip through the
	// management API on every lane. The behavioral-durability scenario
	// (recognized by name) proves a quorum queue's marker survives a
	// live broker loss.
	case "kubernetesrabbitmq":
		spec := manifestSpecMap(manifestPath)
		replicas := 1
		if n, ok := specInt(spec["replicas"]); ok && n > 0 {
			replicas = n
		}
		return &RabbitMqClusterVerifier{
			Namespace:   info.Namespace,
			ClusterName: info.Name,
			Replicas:    replicas,
			Durability:  strings.Contains(manifestPath, "behavioral-durability"),
		}, nil

	// An operator-managed ClickHouse cluster: status Completed with
	// every declared host reconciled, the managed Keeper when the
	// coordination calls for one, plus a LIVE SQL round-trip on every
	// lane. The behavioral-durability scenario (recognized by name)
	// proves replica-synced rows survive a live replica loss.
	case "kubernetesclickhouse":
		spec := manifestSpecMap(manifestPath)
		clusterName, totalHosts, username, password, managedKeeper := clickhouseScenarioShape(spec)
		return &ClickHouseInstallationVerifier{
			Namespace:     info.Namespace,
			ChiName:       info.Name,
			ClusterName:   clusterName,
			TotalHosts:    totalHosts,
			Username:      username,
			Password:      password,
			ManagedKeeper: managedKeeper,
			Durability:    strings.Contains(manifestPath, "behavioral-durability"),
		}, nil

	// The OpenSearch Kubernetes Operator: controller-manager Available +
	// the module-owned CRDs Established (kept on destroy by design).
	case "kubernetesopensearchoperator":
		return &OpenSearchOperatorInstallVerifier{
			Namespace:   info.Namespace,
			ReleaseName: info.Name,
		}, nil

	// An operator-managed OpenSearch cluster: phase RUNNING with every
	// declared node available, plus a LIVE index/search round-trip on
	// every lane. The behavioral-durability scenario (recognized by
	// name) proves replicated data survives a live node loss.
	case "kubernetesopensearch":
		spec := manifestSpecMap(manifestPath)
		firstPool, totalNodes := opensearchFirstPool(spec)
		return &OpenSearchClusterVerifier{
			Namespace:     info.Namespace,
			ClusterName:   info.Name,
			TotalNodes:    totalNodes,
			FirstPoolName: firstPool,
			HttpPort:      int(manifestSpecInt(manifestPath, "http_port", 9200)),
			Durability:    strings.Contains(manifestPath, "behavioral-durability"),
			Dashboards:    opensearchDashboardsEnabled(spec),
		}, nil

	// The Apache Solr Operator: operator Available + the module-owned
	// CRDs (incl. the bundled zookeeper-operator's ZookeeperCluster)
	// Established; kept on destroy by design.
	case "kubernetessolroperator":
		spec := manifestSpecMap(manifestPath)
		return &SolrOperatorInstallVerifier{
			Namespace:         info.Namespace,
			ReleaseName:       info.Name,
			ZookeeperOperator: solrZookeeperOperatorInstalled(spec),
		}, nil

	// An operator-managed SolrCloud: every node ready + the common
	// Service. The behavioral-collection scenario (recognized by name)
	// creates a collection, indexes a document and queries it back —
	// the cluster proven working, not just running.
	case "kubernetessolr":
		spec := manifestSpecMap(manifestPath)
		return &SolrCloudVerifier{
			Namespace:   info.Namespace,
			ClusterName: info.Name,
			Replicas:    int(manifestSpecInt(manifestPath, "replicas", 3)),
			Security:    solrSecurityEnabled(spec),
			Behavioral:  strings.Contains(manifestPath, "behavioral-collection"),
		}, nil

	// A Neo4j server: the StatefulSet ready + the main Service, plus a
	// LIVE Cypher write/read round-trip whenever credentials are
	// A SeaweedFS object store: master + filer StatefulSets ready (the
	// dedicated S3 gateway Deployment too when deployed), the `-s3`
	// Service present, declared buckets created by the install hook, and
	// a live S3 put/get round-trip with the chart-materialized
	// credentials. The behavioral-durability scenario (recognized by
	// name) proves object bytes survive a volume-server pod loss through
	// the PVC.
	case "kubernetesseaweedfs":
		spec := manifestSpecMap(manifestPath)
		return &SeaweedFsVerifier{
			Namespace:             info.Namespace,
			Name:                  info.Name,
			DedicatedS3:           seaweedfsS3Dedicated(spec),
			CredentialsSecretName: seaweedfsCredentialsSecret(spec, info.Name),
			Buckets:               seaweedfsBuckets(spec),
			AdminEnabled:          seaweedfsAdminEnabled(spec),
			Durability:            strings.Contains(manifestPath, "behavioral-durability"),
		}, nil

	// A kube-prometheus-stack: operator Deployment available, the
	// operator-reconciled Prometheus StatefulSet at its declared count,
	// Alertmanager/bundled-Grafana when enabled, and a LIVE metric-flow
	// proof (a healthy kube-state-metrics target + a PromQL answer). The
	// behavioral-alerting scenario (recognized by name) proves the
	// pipeline end to end via the always-firing Watchdog alert in BOTH
	// Prometheus and Alertmanager. Destroy asserts the crds-subchart
	// keep posture (the monitoring CRDs must SURVIVE uninstall).
	case "kuberneteskubeprometheusstack":
		spec := manifestSpecMap(manifestPath)
		return &KubePrometheusStackVerifier{
			Namespace:            info.Namespace,
			Name:                 info.Name,
			PrometheusReplicas:   kpsReplicas(spec, "prometheus"),
			AlertmanagerEnabled:  kpsHalfEnabled(spec, "alertmanager"),
			AlertmanagerReplicas: kpsReplicas(spec, "alertmanager"),
			GrafanaEnabled:       kpsHalfEnabled(spec, "grafana"),
			Alerting:             strings.Contains(manifestPath, "behavioral-alerting"),
		}, nil

	// A standalone Grafana: Deployment available, /api/health reporting
	// a working database, an AUTHENTICATED API round-trip as the admin
	// credentials (read from the Secret — the credential-wiring proof),
	// and every declared datasource provisioned. The
	// behavioral-persistence scenario (recognized by name) proves a
	// UI-authored dashboard survives a pod REPLACEMENT through the PVC.
	case "kubernetesgrafana":
		spec := manifestSpecMap(manifestPath)
		return &GrafanaVerifier{
			Namespace:        info.Namespace,
			Name:             info.Name,
			AdminSecretName:  grafanaAdminSecretName(spec, info.Name),
			AdminUserKey:     grafanaAdminSecretKey(spec, "user_key"),
			AdminPasswordKey: grafanaAdminSecretKey(spec, "password_key"),
			Datasources:      grafanaDatasourceNames(spec),
			Persistence:      strings.Contains(manifestPath, "behavioral-persistence"),
		}, nil

	// A Grafana Loki log store: the monolithic StatefulSet ready, the
	// gateway Service present, and a LIVE push→LogQL round-trip (a line
	// pushed through the gateway and returned by a LogQL query). With
	// tenants declared the proof authenticates AS the first tenant —
	// the gateway guards every push/query route and derives the tenant
	// from the authenticated user, so the round-trip then proves the
	// whole multi-tenant path. The behavioral-durability scenario
	// (recognized by name) proves logs survive a pod loss through the
	// PVC.
	case "kubernetesloki":
		spec := manifestSpecMap(manifestPath)
		v := &LokiVerifier{
			Namespace:      info.Namespace,
			Name:           info.Name,
			GatewayEnabled: lokiGatewayEnabled(spec),
			Durability:     strings.Contains(manifestPath, "behavioral-durability"),
		}
		if tenant := lokiFirstTenantUser(spec); tenant != "" {
			v.TenantUser = tenant
			v.TenantPassword = lokiE2eTenantPassword
		}
		return v, nil

	// A Grafana Tempo trace store: the StatefulSet ready, the Service
	// present, and a LIVE OTLP-push→trace-by-ID round-trip. The
	// behavioral-persistence scenario (recognized by name) proves traces
	// survive a pod loss through the PVC.
	case "kubernetestempo":
		return &TempoVerifier{
			Namespace:   info.Namespace,
			Name:        info.Name,
			Persistence: strings.Contains(manifestPath, "behavioral-persistence"),
		}, nil

	// A SigNoz observability platform: server StatefulSet + collector
	// rolled out, the bundled ClickHouseInstallation reconciled (bundled
	// arm), health OK, and the product-grade round-trip — first-admin
	// REGISTRATION through the product API, a session, an OTLP span
	// push, and the trace retrieved by ID through the authenticated
	// query API. The behavioral-state scenario (recognized by name)
	// re-authenticates and re-queries AFTER a server pod replacement —
	// users survive on the SQLite PVC, telemetry in ClickHouse.
	case "kubernetessignoz":
		spec := manifestSpecMap(manifestPath)
		return &SignozVerifier{
			Namespace:         info.Namespace,
			Name:              info.Name,
			BundledClickHouse: signozBundledClickHouse(spec),
			StateProof:        strings.Contains(manifestPath, "behavioral-state"),
		}, nil

	// A Qdrant vector database: replicas ready, the main Service
	// present, and a live vector round-trip (collection create → upsert
	// → similarity search asserting the nearest neighbour), carrying the
	// declared API key when auth is on. The behavioral-persistence
	// scenario (recognized by name) proves vectors survive pod loss
	// through the PVC.
	case "kubernetesqdrant":
		spec := manifestSpecMap(manifestPath)
		return &QdrantVerifier{
			Namespace:        info.Namespace,
			Name:             info.Name,
			Replicas:         int(manifestSpecInt(manifestPath, "replicas", 1)),
			ApiKeySecretName: qdrantApiKeySecretName(spec, info.Name),
			Persistence:      strings.Contains(manifestPath, "behavioral-persistence"),
		}, nil

	// declared. The behavioral-persistence scenario (recognized by
	// name) proves the data survives pod loss through the PVC.
	case "kubernetesneo4j":
		spec := manifestSpecMap(manifestPath)
		return &Neo4jVerifier{
			Namespace:      info.Namespace,
			Name:           info.Name,
			AuthSecretName: neo4jAuthSecretName(spec, info.Name),
			Persistence:    strings.Contains(manifestPath, "behavioral-persistence"),
		}, nil

	// A Percona-operator-managed MySQL (XtraDB) cluster: the resource in
	// state ready + the proxy's write Service. The behavioral-durability
	// scenario proves Galera's synchronous replication through a live
	// node loss; the with-backup scenario drives a real XtraBackup to
	// completion.
	case "kubernetesmysql":
		spec := manifestSpecMap(manifestPath)
		return &PxcClusterVerifier{
			Namespace:     info.Namespace,
			ClusterName:   info.Name,
			Size:          manifestSpecInt(manifestPath, "instances", 3),
			ProxyService:  mysqlProxyService(info.Name, spec),
			Behavioral:    strings.Contains(manifestPath, "behavioral-durability"),
			BackupProof:   strings.Contains(manifestPath, "with-backup"),
			BackupStorage: mysqlFirstBackupStorage(spec),
		}, nil

	// A Valkey instance from the official chart: workload ready + the
	// write Service. The behavioral-persistence scenario proves
	// durability through a pod loss; the replication scenario proves
	// propagation through the read Service.
	case "kubernetesvalkey":
		spec := manifestSpecMap(manifestPath)
		replication, replicas := valkeyReplication(spec)
		_, hasAuth := spec["auth"]
		return &ValkeyVerifier{
			Namespace:        info.Namespace,
			Name:             info.Name,
			Replication:      replication,
			Replicas:         replicas,
			AuthDeclared:     hasAuth,
			PersistenceProof: strings.Contains(manifestPath, "behavioral-persistence"),
			ReplicationProof: strings.Contains(manifestPath, "with-replication"),
			ServicePort:      valkeyServicePort(spec),
		}, nil

	// A CloudNativePG-managed PostgreSQL cluster: Ready condition + every
	// declared instance ready + the -rw Service. The behavioral-failover
	// scenario (recognized by name) additionally proves data durability
	// through a live primary loss and promotion.
	case "kubernetespostgres":
		return &CnpgClusterVerifier{
			Namespace:   info.Namespace,
			ClusterName: info.Name,
			Instances:   manifestSpecInt(manifestPath, "instances", 1),
			Behavioral:  strings.Contains(manifestPath, "behavioral-failover"),
			BackupProof: strings.Contains(manifestPath, "with-backup"),
		}, nil

	// cert-manager installation: the three component Deployments must be
	// Available and the core CRDs Established — the preconditions every
	// Issuer/Certificate apply depends on.
	case "kubernetescertmanager":
		return &CertManagerInstallVerifier{
			Namespace:     info.Namespace,
			ComponentName: info.Name,
		}, nil

	// Issuers with in-cluster backends (self-signed, CA) reach Ready with
	// no external dependency, so Ready is the verifiable contract on the
	// kind cluster. ClusterIssuer is cluster-scoped (no namespace).
	case "kubernetesclusterissuer":
		return &IssuerVerifier{
			Kind: "clusterissuer",
			Name: info.Name,
		}, nil

	case "kubernetesissuer":
		return &IssuerVerifier{
			Kind:      "issuer",
			Name:      info.Name,
			Namespace: info.Namespace,
		}, nil

	// Real issuance proof: Ready condition + signed material present in the
	// target TLS Secret.
	case "kubernetescertificate":
		secretName, err := manifestSpecString(manifestPath, "secretName")
		if err != nil {
			return nil, err
		}
		return &CertificateVerifier{
			Name:       info.Name,
			Namespace:  info.Namespace,
			SecretName: secretName,
		}, nil

	// Istio control-plane installation: istiod Available + the Istio CRDs
	// Established; ambient scenarios additionally require the istio-cni and
	// ztunnel DaemonSets fully ready. The istiod name carries the revision
	// suffix when the scenario names one.
	case "kubernetesistio":
		istiodName := "istiod"
		if revision, _ := manifestSpecString(manifestPath, "revision"); revision != "" {
			istiodName = "istiod-" + revision
		}
		dataplaneMode, _ := manifestSpecString(manifestPath, "dataplaneMode")
		return &IstioInstallVerifier{
			Namespace:  info.Namespace,
			IstiodName: istiodName,
			Ambient:    dataplaneMode == "ambient",
		}, nil

	// ExternalDNS installation: the controller Deployment must be Available
	// (fullname pinned to metadata.name — multi-instance per cluster is
	// first-class, so the release's own Deployment is the contract, not the
	// namespace).
	case "kubernetesexternaldns":
		// The aws-route53-irsa scenario (real-cluster profile) proves REAL
		// record writes: verifier-owned source, records asserted in the
		// hosted zone via the AWS API, then removed with the source.
		if strings.Contains(manifestPath, "aws-route53-irsa") {
			return &ExternalDnsRoute53Verifier{
				Namespace:     info.Namespace,
				ComponentName: info.Name,
			}, nil
		}
		return &ExternalDnsInstallVerifier{
			Namespace:     info.Namespace,
			ComponentName: info.Name,
		}, nil

	// External Secrets Operator installation: the three component
	// Deployments must be Available and the core CRDs Established — the
	// preconditions every SecretStore/ExternalSecret apply depends on.
	case "kubernetesexternalsecretsoperator":
		return &ExternalSecretsOperatorInstallVerifier{
			Namespace:     info.Namespace,
			ComponentName: info.Name,
		}, nil

	// Stores with in-cluster backends (the fake provider) reach Ready with
	// no external dependency, so Ready is the verifiable contract on the
	// kind cluster. ClusterSecretStore is cluster-scoped (no namespace).
	case "kubernetesclustersecretstore":
		return &SecretStoreVerifier{
			Kind: "clustersecretstore",
			Name: info.Name,
			// The aws-sm-irsa scenario (real-cluster profile) proves a real
			// Secrets Manager read through the keyless IRSA identity.
			BehavioralAwsSm: strings.Contains(manifestPath, "aws-sm-irsa"),
		}, nil

	case "kubernetessecretstore":
		return &SecretStoreVerifier{
			Kind:      "secretstore",
			Name:      info.Name,
			Namespace: info.Namespace,
		}, nil

	// Real sync proof: Ready condition + synced data present in the
	// materialized Secret.
	case "kubernetesexternalsecret":
		secretName, err := externalSecretTargetName(manifestPath)
		if err != nil {
			return nil, err
		}
		return &ExternalSecretSyncVerifier{
			Name:       info.Name,
			Namespace:  info.Namespace,
			SecretName: secretName,
		}, nil

	// ingress-nginx installation: controller Deployment Available + the
	// instance's IngressClass present. Fullname pinned to metadata.name —
	// multiple controllers per cluster are first-class, so the release's
	// own objects are the contract, not the namespace.
	case "kubernetesingressnginx":
		className := manifestNestedSpecString(manifestPath, "ingressClass", "name")
		if className == "" {
			className = "nginx"
		}
		return &IngressNginxInstallVerifier{
			Namespace:        info.Namespace,
			ComponentName:    info.Name,
			IngressClassName: className,
			// The aws-nlb scenario (real-cluster profile) validates the
			// documented AWS annotation recipe down to a provisioned LB.
			LbAddress: strings.Contains(manifestPath, "aws-nlb"),
		}, nil

	// metrics-server installation, verified to METRIC FLOW (Deployment
	// Available + APIService Available + kubectl top returning values).
	case "kubernetesmetricsserver":
		return &MetricsServerInstallVerifier{
			Namespace: info.Namespace,
		}, nil

	case "kubernetesmanifest":
		return &ConfigGroupVerifier{
			Namespace:    info.Namespace,
			ManifestPath: manifestPath,
		}, nil

	// A real Helm release of an arbitrary chart: the Helm verifier's
	// namespace + running-pods + service assertions hold for any chart that
	// deploys a workload (the e2e chart, podinfo, deploys one pod and one
	// service).
	case "kuberneteshelmrelease":
		return &HelmComponentVerifier{
			Namespace:     info.Namespace,
			ComponentName: info.Name,
		}, nil

	// Fixed-namespace components: the proto spec has no namespace field because
	// the upstream tooling (Tekton, etc.) uses hardcoded namespaces. The manifest
	// YAML cannot carry a namespace hint because protojson.Unmarshal rejects
	// unknown fields. Namespace is therefore embedded here.
	case "kubernetestekton":
		return &HelmComponentVerifier{
			Namespace:     "tekton-pipelines",
			ComponentName: info.Name,
		}, nil

	case "kubernetestektonoperator":
		return &OperatorComponentVerifier{
			Namespace:     "tekton-operator",
			ComponentName: info.Name,
		}, nil

	default:
		if crdNames, ok := crdInstallKinds[component]; ok {
			return &CRDInstallVerifier{
				ComponentName: info.Name,
				CRDNames:      crdNames,
			}, nil
		}
		// HTTPRoute: scenarios that install a real Gateway API implementation
		// (KubernetesIstio prerequisite) get the behavioral verifier — the
		// route Accepted, its Gateway Programmed, and a live request routed
		// through the auto-provisioned gateway. Scenarios without it keep the
		// object as the verifiable contract.
		if component == "kuberneteshttproute" && manifestHasPrerequisite(manifestPath, "KubernetesIstio") {
			hostname, _ := manifestSpecFirstString(manifestPath, "hostnames")
			return &GatewayRoutingBehavioralVerifier{
				Namespace:   info.Namespace,
				RouteName:   info.Name,
				GatewayName: "e2e-istio-gw",
				Hostname:    hostname,
			}, nil
		}
		// Gateway: the aws-lb-address scenario (real-cluster profile, LB-
		// default mesh fixture) proves Programmed + a real cloud address in
		// .status.addresses — the half the kind lanes pin away.
		if component == "kubernetesgateway" && strings.Contains(manifestPath, "aws-lb-address") {
			return &GatewayLbAddressVerifier{
				Namespace: info.Namespace,
				Name:      info.Name,
			}, nil
		}
		if gw, ok := gatewayApiKinds[component]; ok {
			namespace := info.Namespace
			if gw.clusterScoped {
				namespace = ""
			}
			return &ResourceExistenceVerifier{
				Namespace: namespace,
				Kind:      gw.resource,
				Name:      info.Name,
			}, nil
		}
		// AuthorizationPolicy: scenarios that install a real mesh
		// (KubernetesIstio prerequisite) get the behavioral verifier — a
		// meshed client's request must actually be DENIED (403) and succeed
		// again once the policy is destroyed. Scenarios without it keep the
		// object as the verifiable contract.
		if component == "kubernetesauthorizationpolicy" && manifestHasPrerequisite(manifestPath, "KubernetesIstio") {
			return &AuthzDenyBehavioralVerifier{
				Namespace:        info.Namespace,
				PolicyName:       info.Name,
				ClientDeployment: "e2e-authz-client",
				BackendURL:       "http://e2e-authz-backend." + info.Namespace + ".svc.cluster.local/",
			}, nil
		}
		if resource, ok := istioApiKinds[component]; ok {
			return &ResourceExistenceVerifier{
				Namespace: info.Namespace,
				Kind:      resource,
				Name:      info.Name,
			}, nil
		}
		if operatorKinds[component] {
			return &OperatorComponentVerifier{
				Namespace:     info.Namespace,
				ComponentName: info.Name,
			}, nil
		}
		if helmTier2Kinds[component] {
			return &HelmComponentVerifier{
				Namespace:     info.Namespace,
				ComponentName: info.Name,
			}, nil
		}
		return &GenericVerifier{Component: component}, nil
	}
}
