package module

import (
	"github.com/pkg/errors"
	gcpeventarcmessagebusv1alpha1 "github.com/plantonhq/planton/catalog/gcp/gcpeventarcmessagebus/v1alpha1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/eventarc"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/projects"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// messageBus provisions the Eventarc Advanced family: the bus first, then
// its satellites wired BY RESOURCE REFERENCE — api sources point their
// destination at the created bus's computed full name, and enrollments
// point at the created pipelines' computed full names. Reference wiring
// (never string assembly) keeps the rendered names correct in the
// ambient-project case and gives the engine the dependency order for
// free; the Terraform module wires the same attributes for byte-identical
// results.
func messageBus(ctx *pulumi.Context, locals *Locals, gcpProvider *gcp.Provider) error {
	spec := locals.GcpEventarcMessageBus.Spec

	// Enable the Eventarc API so a fresh project can host the family.
	// disable_on_destroy stays false (the provider default): tearing down
	// one bus must never disable Eventarc for everything else in the
	// project. Matches the Terraform module.
	serviceArgs := &projects.ServiceArgs{
		Service:                  pulumi.String("eventarc.googleapis.com"),
		DisableDependentServices: pulumi.BoolPtr(true),
	}
	if spec.ProjectId.GetValue() != "" {
		serviceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	createdProjectService, err := projects.NewService(ctx,
		"eventarcmessagebus-eventarc.googleapis.com", serviceArgs, pulumi.Provider(gcpProvider))
	if err != nil {
		return errors.Wrap(err, "failed to enable eventarc.googleapis.com api")
	}

	busArgs := &eventarc.MessageBusArgs{
		MessageBusId: pulumi.String(locals.MessageBusId),
		Location:     pulumi.String(spec.Location),
		Labels:       pulumi.ToStringMap(locals.GcpLabels),
	}
	// Honor the spec contract: an empty project_id falls back to the
	// provider's default project.
	if spec.ProjectId.GetValue() != "" {
		busArgs.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	if spec.DisplayName != "" {
		busArgs.DisplayName = pulumi.StringPtr(spec.DisplayName)
	}
	if len(spec.Annotations) > 0 {
		busArgs.Annotations = pulumi.ToStringMap(spec.Annotations)
	}
	if spec.CryptoKey.GetValue() != "" {
		busArgs.CryptoKeyName = pulumi.StringPtr(spec.CryptoKey.GetValue())
	}
	if spec.LogSeverity != "" {
		busArgs.LoggingConfig = &eventarc.MessageBusLoggingConfigArgs{
			LogSeverity: pulumi.StringPtr(spec.LogSeverity),
		}
	}
	// Unset defers to the provider default (DELETE). The same value is
	// wired to every satellite below — one spec lever, every resource.
	if spec.DeletionPolicy != "" {
		busArgs.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
	}

	createdBus, err := eventarc.NewMessageBus(ctx, "message-bus", busArgs,
		pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
	if err != nil {
		return errors.Wrap(err, "failed to create message bus")
	}

	// Google API sources — each auto-wired to THIS bus (the kind's
	// contract: a source feeding another bus belongs to that bus's kind
	// instance).
	for _, source := range spec.GoogleApiSources {
		sourceArgs := &eventarc.GoogleApiSourceArgs{
			GoogleApiSourceId: pulumi.String(source.SourceId),
			Location:          pulumi.String(spec.Location),
			Destination:       createdBus.Name,
			Labels:            pulumi.ToStringMap(satelliteLabels(locals.GcpLabels, source.Labels)),
		}
		if spec.ProjectId.GetValue() != "" {
			sourceArgs.Project = pulumi.String(spec.ProjectId.GetValue())
		}
		if source.DisplayName != "" {
			sourceArgs.DisplayName = pulumi.StringPtr(source.DisplayName)
		}
		if len(source.Annotations) > 0 {
			sourceArgs.Annotations = pulumi.ToStringMap(source.Annotations)
		}
		if source.CryptoKey.GetValue() != "" {
			sourceArgs.CryptoKeyName = pulumi.StringPtr(source.CryptoKey.GetValue())
		}
		if source.LogSeverity != "" {
			sourceArgs.LoggingConfig = &eventarc.GoogleApiSourceLoggingConfigArgs{
				LogSeverity: pulumi.StringPtr(source.LogSeverity),
			}
		}
		if spec.DeletionPolicy != "" {
			sourceArgs.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
		}
		if _, err := eventarc.NewGoogleApiSource(ctx, "api-source-"+source.SourceId, sourceArgs,
			pulumi.Provider(gcpProvider)); err != nil {
			return errors.Wrapf(err, "failed to create google api source %s", source.SourceId)
		}
	}

	// Pipelines — created before enrollments so enrollments can reference
	// their computed full names.
	createdPipelines := map[string]*eventarc.Pipeline{}
	for _, pipeline := range spec.Pipelines {
		pipelineArgs, err := expandPipeline(spec, locals, pipeline)
		if err != nil {
			return err
		}
		createdPipeline, err := eventarc.NewPipeline(ctx, "pipeline-"+pipeline.PipelineId, pipelineArgs,
			pulumi.Provider(gcpProvider), pulumi.DependsOn([]pulumi.Resource{createdProjectService}))
		if err != nil {
			return errors.Wrapf(err, "failed to create pipeline %s", pipeline.PipelineId)
		}
		createdPipelines[pipeline.PipelineId] = createdPipeline
	}

	// Enrollments — the routing table binding bus messages to pipelines.
	// The sibling-id contract is proto-CEL-enforced, so the lookup always
	// resolves.
	for _, enrollment := range spec.Enrollments {
		enrollmentArgs := &eventarc.EnrollmentArgs{
			EnrollmentId: pulumi.String(enrollment.EnrollmentId),
			Location:     pulumi.String(spec.Location),
			CelMatch:     pulumi.String(enrollment.CelMatch),
			MessageBus:   createdBus.Name,
			Destination:  createdPipelines[enrollment.Pipeline].Name,
			Labels:       pulumi.ToStringMap(satelliteLabels(locals.GcpLabels, enrollment.Labels)),
		}
		if spec.ProjectId.GetValue() != "" {
			enrollmentArgs.Project = pulumi.String(spec.ProjectId.GetValue())
		}
		if enrollment.DisplayName != "" {
			enrollmentArgs.DisplayName = pulumi.StringPtr(enrollment.DisplayName)
		}
		if len(enrollment.Annotations) > 0 {
			enrollmentArgs.Annotations = pulumi.ToStringMap(enrollment.Annotations)
		}
		if spec.DeletionPolicy != "" {
			enrollmentArgs.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
		}
		if _, err := eventarc.NewEnrollment(ctx, "enrollment-"+enrollment.EnrollmentId, enrollmentArgs,
			pulumi.Provider(gcpProvider)); err != nil {
			return errors.Wrapf(err, "failed to create enrollment %s", enrollment.EnrollmentId)
		}
	}

	ctx.Export(OpMessageBusName, createdBus.Name)

	return nil
}

// expandPipeline maps one spec pipeline onto the provider's pipeline
// resource. The provider supports exactly one destination per pipeline
// (its own schema note); the spec models exactly one, and the wiring
// renders the single-element destinations list the provider expects. The
// spec's pipeline-level output_payload_format also lives INSIDE that
// destination element (the provider's shape).
func expandPipeline(spec *gcpeventarcmessagebusv1alpha1.GcpEventarcMessageBusSpec, locals *Locals,
	pipeline *gcpeventarcmessagebusv1alpha1.GcpEventarcMessageBusPipeline) (*eventarc.PipelineArgs, error) {

	destination := pipeline.Destination
	destinationArgs := &eventarc.PipelineDestinationArgs{}

	if httpEndpoint := destination.HttpEndpoint; httpEndpoint != nil {
		httpArgs := &eventarc.PipelineDestinationHttpEndpointArgs{
			Uri: pulumi.String(httpEndpoint.Uri),
		}
		if httpEndpoint.MessageBindingTemplate != "" {
			httpArgs.MessageBindingTemplate = pulumi.StringPtr(httpEndpoint.MessageBindingTemplate)
		}
		destinationArgs.HttpEndpoint = httpArgs
		// Required for HTTP endpoints, forbidden otherwise (provider rule,
		// proto-CEL-enforced) — the spec carries it inside the arm; the
		// wiring restores the provider's sibling network_config shape.
		destinationArgs.NetworkConfig = &eventarc.PipelineDestinationNetworkConfigArgs{
			NetworkAttachment: pulumi.StringPtr(httpEndpoint.NetworkAttachment),
		}
	}
	if destination.Topic.GetValue() != "" {
		destinationArgs.Topic = pulumi.StringPtr(destination.Topic.GetValue())
	}
	if destination.Workflow.GetValue() != "" {
		destinationArgs.Workflow = pulumi.StringPtr(destination.Workflow.GetValue())
	}
	if destination.MessageBus != "" {
		destinationArgs.MessageBus = pulumi.StringPtr(destination.MessageBus)
	}

	if auth := pipeline.Authentication; auth != nil {
		authArgs := &eventarc.PipelineDestinationAuthenticationConfigArgs{}
		if oidc := auth.GoogleOidc; oidc != nil {
			oidcArgs := &eventarc.PipelineDestinationAuthenticationConfigGoogleOidcArgs{
				ServiceAccount: pulumi.String(oidc.ServiceAccount.GetValue()),
			}
			if oidc.Audience != "" {
				oidcArgs.Audience = pulumi.StringPtr(oidc.Audience)
			}
			authArgs.GoogleOidc = oidcArgs
		}
		if oauth := auth.OauthToken; oauth != nil {
			oauthArgs := &eventarc.PipelineDestinationAuthenticationConfigOauthTokenArgs{
				ServiceAccount: pulumi.String(oauth.ServiceAccount.GetValue()),
			}
			if oauth.Scope != "" {
				oauthArgs.Scope = pulumi.StringPtr(oauth.Scope)
			}
			authArgs.OauthToken = oauthArgs
		}
		destinationArgs.AuthenticationConfig = authArgs
	}

	if pipeline.OutputPayloadFormat != nil {
		destinationArgs.OutputPayloadFormat = expandOutputPayloadFormat(pipeline.OutputPayloadFormat)
	}

	args := &eventarc.PipelineArgs{
		PipelineId:   pulumi.String(pipeline.PipelineId),
		Location:     pulumi.String(spec.Location),
		Destinations: eventarc.PipelineDestinationArray{destinationArgs},
		Labels:       pulumi.ToStringMap(satelliteLabels(locals.GcpLabels, pipeline.Labels)),
	}
	if spec.ProjectId.GetValue() != "" {
		args.Project = pulumi.String(spec.ProjectId.GetValue())
	}
	if pipeline.DisplayName != "" {
		args.DisplayName = pulumi.StringPtr(pipeline.DisplayName)
	}
	if len(pipeline.Annotations) > 0 {
		args.Annotations = pulumi.ToStringMap(pipeline.Annotations)
	}
	if pipeline.CryptoKey.GetValue() != "" {
		args.CryptoKeyName = pulumi.StringPtr(pipeline.CryptoKey.GetValue())
	}
	if pipeline.LogSeverity != "" {
		args.LoggingConfig = &eventarc.PipelineLoggingConfigArgs{
			LogSeverity: pulumi.StringPtr(pipeline.LogSeverity),
		}
	}
	if pipeline.InputPayloadFormat != nil {
		args.InputPayloadFormat = expandInputPayloadFormat(pipeline.InputPayloadFormat)
	}
	// The API allows at most ONE mediation (transformation) per pipeline —
	// the single spec template renders as the single-element list.
	if pipeline.MediationTransformationTemplate != "" {
		args.Mediations = eventarc.PipelineMediationArray{
			&eventarc.PipelineMediationArgs{
				Transformation: &eventarc.PipelineMediationTransformationArgs{
					TransformationTemplate: pulumi.StringPtr(pipeline.MediationTransformationTemplate),
				},
			},
		}
	}
	if retry := pipeline.RetryPolicy; retry != nil {
		retryArgs := &eventarc.PipelineRetryPolicyArgs{}
		if retry.MaxAttempts != 0 {
			retryArgs.MaxAttempts = pulumi.IntPtr(int(retry.MaxAttempts))
		}
		if retry.MinRetryDelay != "" {
			retryArgs.MinRetryDelay = pulumi.StringPtr(retry.MinRetryDelay)
		}
		if retry.MaxRetryDelay != "" {
			retryArgs.MaxRetryDelay = pulumi.StringPtr(retry.MaxRetryDelay)
		}
		args.RetryPolicy = retryArgs
	}
	if spec.DeletionPolicy != "" {
		args.DeletionPolicy = pulumi.String(spec.DeletionPolicy)
	}

	return args, nil
}

// expandInputPayloadFormat maps the spec's exactly-one format form onto
// the provider's input_payload_format block. JSON is a presence marker (an
// empty block selects it).
func expandInputPayloadFormat(format *gcpeventarcmessagebusv1alpha1.GcpEventarcMessageBusPayloadFormat) *eventarc.PipelineInputPayloadFormatArgs {
	args := &eventarc.PipelineInputPayloadFormatArgs{}
	if avro := format.Avro; avro != nil {
		avroArgs := &eventarc.PipelineInputPayloadFormatAvroArgs{}
		if avro.SchemaDefinition != "" {
			avroArgs.SchemaDefinition = pulumi.StringPtr(avro.SchemaDefinition)
		}
		args.Avro = avroArgs
	}
	if format.Json {
		args.Json = &eventarc.PipelineInputPayloadFormatJsonArgs{}
	}
	if protobuf := format.Protobuf; protobuf != nil {
		protobufArgs := &eventarc.PipelineInputPayloadFormatProtobufArgs{}
		if protobuf.SchemaDefinition != "" {
			protobufArgs.SchemaDefinition = pulumi.StringPtr(protobuf.SchemaDefinition)
		}
		args.Protobuf = protobufArgs
	}
	return args
}

// expandOutputPayloadFormat is the destination-side mirror of
// expandInputPayloadFormat (the provider nests output format inside the
// destination element).
func expandOutputPayloadFormat(format *gcpeventarcmessagebusv1alpha1.GcpEventarcMessageBusPayloadFormat) *eventarc.PipelineDestinationOutputPayloadFormatArgs {
	args := &eventarc.PipelineDestinationOutputPayloadFormatArgs{}
	if avro := format.Avro; avro != nil {
		avroArgs := &eventarc.PipelineDestinationOutputPayloadFormatAvroArgs{}
		if avro.SchemaDefinition != "" {
			avroArgs.SchemaDefinition = pulumi.StringPtr(avro.SchemaDefinition)
		}
		args.Avro = avroArgs
	}
	if format.Json {
		args.Json = &eventarc.PipelineDestinationOutputPayloadFormatJsonArgs{}
	}
	if protobuf := format.Protobuf; protobuf != nil {
		protobufArgs := &eventarc.PipelineDestinationOutputPayloadFormatProtobufArgs{}
		if protobuf.SchemaDefinition != "" {
			protobufArgs.SchemaDefinition = pulumi.StringPtr(protobuf.SchemaDefinition)
		}
		args.Protobuf = protobufArgs
	}
	return args
}
