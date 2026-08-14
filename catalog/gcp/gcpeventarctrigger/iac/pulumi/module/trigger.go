package module

import (
	"github.com/pkg/errors"
	gcpeventarctriggerv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpeventarctrigger/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/eventarc"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/organizations"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// trigger provisions the Eventarc trigger and its count-gated companions:
// a partner CHANNEL when the spec arms partner_channel (the trigger is
// wired to it — the channel's one-time activation_token is exported for
// the partner handshake), and the per-project-per-location
// googleChannelConfig SINGLETON when google_channel_crypto_key is set
// (adopt-and-patch semantics: the singleton always exists in GCP; its
// delete is a state-only no-op in the provider).
//
// The first trigger in a project provisions Eventarc's service agent —
// the first delivery can lag a few minutes behind the apply (P4SA
// propagation).
func trigger(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpEventarcTrigger.Spec

	// Enable the Eventarc API so a fresh project can host the trigger.
	// disable_on_destroy stays false (the provider default): tearing down
	// one trigger must never disable Eventarc for everything else in the
	// project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("eventarc.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"eventarctrigger-eventarc.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable eventarc.googleapis.com api")
	}

	dependencies := []pulumi.Resource{createdProjectService}

	args := &eventarc.TriggerArgs{
		Name:     pulumi.String(locals.TriggerName),
		Location: pulumi.String(spec.Location),
		Labels:   pulumi.ToStringMap(locals.GcpLabels),
	}

	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project.
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}

	criteria := eventarc.TriggerMatchingCriteriaArray{}
	for _, criterion := range spec.MatchingCriteria {
		criterionArgs := &eventarc.TriggerMatchingCriteriaArgs{
			Attribute: pulumi.String(criterion.Attribute),
			Value:     pulumi.String(criterion.Value),
		}
		if criterion.Operator != "" {
			criterionArgs.Operator = pulumi.StringPtr(criterion.Operator)
		}
		criteria = append(criteria, criterionArgs)
	}
	args.MatchingCriterias = criteria

	args.Destination = expandDestination(spec.Destination, spec.Location)

	if spec.ServiceAccount.GetValue() != "" {
		args.ServiceAccount = pulumi.String(spec.ServiceAccount.GetValue())
	}
	if spec.TransportPubsubTopic.GetValue() != "" {
		args.Transport = &eventarc.TriggerTransportArgs{
			Pubsub: &eventarc.TriggerTransportPubsubArgs{
				Topic: pulumi.StringPtr(spec.TransportPubsubTopic.GetValue()),
			},
		}
	}
	if spec.EventDataContentType != "" {
		args.EventDataContentType = pulumi.StringPtr(spec.EventDataContentType)
	}
	if spec.RetryMaxAttempts != 0 {
		args.RetryPolicy = &eventarc.TriggerRetryPolicyArgs{
			MaxAttempts: pulumi.IntPtr(int(spec.RetryMaxAttempts)),
		}
	}
	// Unset defers to the provider default (DELETE). The same value is
	// wired to the companions below — one spec lever, every resource.
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
	}

	// Partner channel companion: created first, then the trigger is wired
	// to its full resource name (assembled from the channel's own computed
	// project attribute so the ambient-project case renders correctly).
	if spec.PartnerChannel != nil {
		channelArgs := &eventarc.ChannelArgs{
			Name:               pulumi.String(locals.ChannelName),
			Location:           pulumi.String(spec.Location),
			ThirdPartyProvider: pulumi.StringPtr(spec.PartnerChannel.ThirdPartyProvider),
			Labels:             pulumi.ToStringMap(locals.GcpLabels),
		}
		if spec.ProjectId.GetValue() != "" {
			channelArgs.Project = pulumi.String(spec.ProjectId.GetValue())
		}
		if spec.PartnerChannel.CryptoKey.GetValue() != "" {
			channelArgs.CryptoKeyName = pulumi.StringPtr(spec.PartnerChannel.CryptoKey.GetValue())
		}
		if spec.DeletionPolicy != "" {
			channelArgs.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
		}
		createdChannel, err := eventarc.NewChannel(ctx, "partner-channel", channelArgs,
			pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
		if err != nil {
			return errors.Wrap(err, "failed to create partner channel")
		}

		args.Channel = pulumi.Sprintf("projects/%s/locations/%s/channels/%s",
			createdChannel.Project, spec.Location, createdChannel.Name)
		dependencies = append(dependencies, createdChannel)

		// The one-time token the SaaS partner needs to complete the
		// handshake — sensitive, exported like a credential.
		ctx.Export(OpPartnerChannelActivationToken, pulumi.ToSecret(createdChannel.ActivationToken))
	} else {
		// Deterministic output surface: the key exists on every apply.
		ctx.Export(OpPartnerChannelActivationToken, pulumi.ToSecret(pulumi.String("")))
	}

	// googleChannelConfig singleton companion: the resource's name arg IS
	// the fixed singleton path, so the ambient-project case resolves the
	// provider's default project from client config (the Terraform
	// module's count-gated google_client_config data source, mirrored).
	if spec.GoogleChannelCryptoKey.GetValue() != "" {
		var project pulumi.StringInput
		if spec.ProjectId.GetValue() != "" {
			project = pulumi.String(spec.ProjectId.GetValue())
		} else {
			clientConfig, err := organizations.GetClientConfig(ctx, pulumi.Provider(gcpProvider))
			if err != nil {
				return errors.Wrap(err, "failed to read provider client config for the default project")
			}
			if clientConfig.Project == "" {
				return errors.New("project_id is empty and the provider has no default project configured")
			}
			project = pulumi.String(clientConfig.Project)
		}

		configArgs := &eventarc.GoogleChannelConfigArgs{
			Name:          pulumi.Sprintf("projects/%s/locations/%s/googleChannelConfig", project, spec.Location),
			Location:      pulumi.String(spec.Location),
			CryptoKeyName: pulumi.StringPtr(spec.GoogleChannelCryptoKey.GetValue()),
			Project:       project,
		}
		createdConfig, err := eventarc.NewGoogleChannelConfig(ctx, "google-channel-config", configArgs,
			pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
		if err != nil {
			return errors.Wrap(err, "failed to configure google channel config")
		}
		dependencies = append(dependencies, createdConfig)
	}

	createdTrigger, err := eventarc.NewTrigger(ctx, "trigger", args,
		pulumi.Provider(gcpProvider), pulumi.DependsOn(dependencies))
	if err != nil {
		return errors.Wrap(err, "failed to create trigger")
	}

	// The module-derived name, not the provider's `name` echo: the
	// provider stores the SHORT name at create but reads back the FULL
	// resource name on refresh (live-caught flip on the Terraform side by
	// the idempotency gate; same latent class here). The name is
	// immutable and the create succeeded with exactly this value.
	// PARITY: matches the Terraform module's local.trigger_name output.
	ctx.Export(OpTriggerName, pulumi.String(locals.TriggerName))
	// The resource ID is the full trigger resource name
	// (projects/{p}/locations/{l}/triggers/{name}) with the ambient
	// project resolved — the trigger's canonical API handle.
	ctx.Export(OpTriggerId, createdTrigger.ID().ToStringOutput())

	return nil
}

// expandDestination maps the spec's exactly-one destination arm onto the
// provider's destination block. The provider validates arm exclusivity
// only server-side; the spec's CEL wall guarantees exactly one arm here.
// triggerLocation feeds the Cloud Run region default (see below).
func expandDestination(destination *gcpeventarctriggerv1alpha1.GcpEventarcTriggerDestination, triggerLocation string) *eventarc.TriggerDestinationArgs {
	destinationArgs := &eventarc.TriggerDestinationArgs{}

	if cloudRun := destination.CloudRunService; cloudRun != nil {
		cloudRunArgs := &eventarc.TriggerDestinationCloudRunServiceArgs{
			Service: pulumi.String(cloudRun.Service.GetValue()),
		}
		// The API REQUIRES the region on every create (live-verified:
		// 400 "cloud_run.region is empty" — it never infers it from the
		// trigger's location), so an empty spec region defaults to the
		// trigger's own location. A "global" trigger cannot self-default;
		// the spec CEL forces an explicit region there.
		if cloudRun.Region != "" {
			cloudRunArgs.Region = pulumi.StringPtr(cloudRun.Region)
		} else {
			cloudRunArgs.Region = pulumi.StringPtr(triggerLocation)
		}
		if cloudRun.Path != "" {
			cloudRunArgs.Path = pulumi.StringPtr(cloudRun.Path)
		}
		destinationArgs.CloudRunService = cloudRunArgs
	}

	if gke := destination.Gke; gke != nil {
		gkeArgs := &eventarc.TriggerDestinationGkeArgs{
			Cluster:   pulumi.String(gke.Cluster.GetValue()),
			Location:  pulumi.String(gke.Location),
			Namespace: pulumi.String(gke.Namespace),
			Service:   pulumi.String(gke.Service),
		}
		if gke.Path != "" {
			gkeArgs.Path = pulumi.StringPtr(gke.Path)
		}
		destinationArgs.Gke = gkeArgs
	}

	if destination.Workflow.GetValue() != "" {
		destinationArgs.Workflow = pulumi.StringPtr(destination.Workflow.GetValue())
	}

	if httpEndpoint := destination.HttpEndpoint; httpEndpoint != nil {
		destinationArgs.HttpEndpoint = &eventarc.TriggerDestinationHttpEndpointArgs{
			Uri: pulumi.String(httpEndpoint.Uri),
		}
		// The provider models the attachment as a sibling network_config
		// block permitted only with HTTP endpoints — the spec carries it
		// inside the arm; the wiring restores the provider shape.
		destinationArgs.NetworkConfig = &eventarc.TriggerDestinationNetworkConfigArgs{
			NetworkAttachment: pulumi.String(httpEndpoint.NetworkAttachment),
		}
	}

	return destinationArgs
}
