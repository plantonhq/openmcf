package module

import (
	cloudflareworkerv1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflareworker/v1alpha1"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
	cloudfl "github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// fkValue unwraps a StringValueOrRef to the literal the provider wants.
// Empty when the field is unset — callers skip the corresponding arg.
func fkValue(r *foreignkeyv1.StringValueOrRef) string {
	if r == nil {
		return ""
	}
	return r.GetValue()
}

// buildBindings flattens the spec's grouped, type-specific binding lists into the
// provider's single discriminated bindings array. Sensitive values (secret_text,
// secret_key material) arrive resolved via StringValueOrRef.GetValue().
func buildBindings(spec *cloudflareworkerv1alpha1.CloudflareWorkerSpec) cloudfl.WorkersScriptBindingArray {
	var bindings cloudfl.WorkersScriptBindingArray

	for k, v := range spec.Vars {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(k),
			Type: pulumi.String("plain_text"),
			Text: pulumi.String(v),
		})
	}
	for _, b := range spec.Secrets {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("secret_text"),
			Text: pulumi.String(b.Value.GetValue()),
		})
	}
	for _, b := range spec.KvNamespaces {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name:        pulumi.String(b.Name),
			Type:        pulumi.String("kv_namespace"),
			NamespaceId: pulumi.String(b.NamespaceId.GetValue()),
		})
	}
	for _, b := range spec.R2Buckets {
		args := cloudfl.WorkersScriptBindingArgs{
			Name:       pulumi.String(b.Name),
			Type:       pulumi.String("r2_bucket"),
			BucketName: pulumi.String(b.BucketName.GetValue()),
		}
		if b.Jurisdiction != "" {
			args.Jurisdiction = pulumi.String(b.Jurisdiction)
		}
		bindings = append(bindings, args)
	}
	for _, b := range spec.D1Databases {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("d1"),
			Id:   pulumi.String(b.DatabaseId.GetValue()),
		})
	}
	for _, b := range spec.HyperdriveConfigs {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("hyperdrive"),
			Id:   pulumi.String(b.ConfigId.GetValue()),
		})
	}
	for _, b := range spec.Services {
		args := cloudfl.WorkersScriptBindingArgs{
			Name:    pulumi.String(b.Name),
			Type:    pulumi.String("service"),
			Service: pulumi.String(b.Service.GetValue()),
		}
		if b.Environment != "" {
			args.Environment = pulumi.String(b.Environment)
		}
		if b.Entrypoint != "" {
			args.Entrypoint = pulumi.String(b.Entrypoint)
		}
		bindings = append(bindings, args)
	}
	for _, b := range spec.Queues {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name:      pulumi.String(b.Name),
			Type:      pulumi.String("queue"),
			QueueName: pulumi.String(b.QueueName.GetValue()),
		})
	}
	for _, b := range spec.DurableObjects {
		args := cloudfl.WorkersScriptBindingArgs{
			Name:      pulumi.String(b.Name),
			Type:      pulumi.String("durable_object_namespace"),
			ClassName: pulumi.String(b.ClassName),
		}
		if v := fkValue(b.ScriptName); v != "" {
			args.ScriptName = pulumi.String(v)
		}
		if b.Environment != "" {
			args.Environment = pulumi.String(b.Environment)
		}
		if b.NamespaceId != "" {
			args.NamespaceId = pulumi.String(b.NamespaceId)
		}
		if b.DispatchNamespace != "" {
			args.DispatchNamespace = pulumi.String(b.DispatchNamespace)
		}
		bindings = append(bindings, args)
	}
	for _, b := range spec.AnalyticsEngineDatasets {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name:    pulumi.String(b.Name),
			Type:    pulumi.String("analytics_engine"),
			Dataset: pulumi.String(b.Dataset),
		})
	}
	for _, b := range spec.VectorizeIndexes {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name:      pulumi.String(b.Name),
			Type:      pulumi.String("vectorize"),
			IndexName: pulumi.String(b.IndexName),
		})
	}
	for _, b := range spec.Ai {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("ai"),
		})
	}
	for _, b := range spec.VersionMetadata {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("version_metadata"),
		})
	}
	for _, b := range spec.MtlsCertificates {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name:          pulumi.String(b.Name),
			Type:          pulumi.String("mtls_certificate"),
			CertificateId: pulumi.String(b.CertificateId),
		})
	}
	for _, b := range spec.DispatchNamespaces {
		args := cloudfl.WorkersScriptBindingArgs{
			Name:      pulumi.String(b.Name),
			Type:      pulumi.String("dispatch_namespace"),
			Namespace: pulumi.String(b.Namespace),
		}
		if o := b.Outbound; o != nil {
			out := &cloudfl.WorkersScriptBindingOutboundArgs{}
			if len(o.Params) > 0 {
				out.Params = pulumi.ToStringArray(o.Params)
			}
			if w := o.Worker; w != nil {
				wa := &cloudfl.WorkersScriptBindingOutboundWorkerArgs{}
				if v := fkValue(w.Service); v != "" {
					wa.Service = pulumi.String(v)
				}
				if w.Environment != "" {
					wa.Environment = pulumi.String(w.Environment)
				}
				out.Worker = wa
			}
			args.Outbound = out
		}
		bindings = append(bindings, args)
	}
	for _, b := range spec.RateLimits {
		args := cloudfl.WorkersScriptBindingArgs{
			Name:      pulumi.String(b.Name),
			Type:      pulumi.String("ratelimit"),
			Namespace: pulumi.String(b.Namespace),
		}
		if s := b.Simple; s != nil {
			// PARITY-EXCEPTION: pulumi-cloudflare SDK v6.17.0's Simple args
			// have Limit + Period only — mitigation_timeout is tofu-only.
			args.Simple = &cloudfl.WorkersScriptBindingSimpleArgs{
				Limit:  pulumi.Float64(s.Limit),
				Period: pulumi.Int(int(s.Period)),
			}
		}
		bindings = append(bindings, args)
	}
	for _, b := range spec.SendEmail {
		args := cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("send_email"),
		}
		if b.DestinationAddress != "" {
			args.DestinationAddress = pulumi.String(b.DestinationAddress)
		}
		if len(b.AllowedDestinationAddresses) > 0 {
			args.AllowedDestinationAddresses = pulumi.ToStringArray(b.AllowedDestinationAddresses)
		}
		if len(b.AllowedSenderAddresses) > 0 {
			args.AllowedSenderAddresses = pulumi.ToStringArray(b.AllowedSenderAddresses)
		}
		bindings = append(bindings, args)
	}
	for _, b := range spec.SecretsStoreSecrets {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name:       pulumi.String(b.Name),
			Type:       pulumi.String("secrets_store_secret"),
			StoreId:    pulumi.String(b.StoreId),
			SecretName: pulumi.String(b.SecretName),
		})
	}
	for _, b := range spec.SecretKeys {
		args := cloudfl.WorkersScriptBindingArgs{
			Name:      pulumi.String(b.Name),
			Type:      pulumi.String("secret_key"),
			Algorithm: pulumi.String(b.Algorithm),
			Format:    pulumi.String(b.Format),
		}
		if len(b.Usages) > 0 {
			args.Usages = pulumi.ToStringArray(b.Usages)
		}
		if b.KeyBase64 != "" {
			args.KeyBase64 = pulumi.String(b.KeyBase64)
		}
		if b.KeyJwk != "" {
			args.KeyJwk = pulumi.String(b.KeyJwk)
		}
		bindings = append(bindings, args)
	}
	for _, b := range spec.Workflows {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name:         pulumi.String(b.Name),
			Type:         pulumi.String("workflow"),
			WorkflowName: pulumi.String(b.WorkflowName),
		})
	}
	for _, b := range spec.Pipelines {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name:     pulumi.String(b.Name),
			Type:     pulumi.String("pipelines"),
			Pipeline: pulumi.String(b.Pipeline),
		})
	}
	for _, b := range spec.JsonBindings {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("json"),
			Json: pulumi.String(b.Json),
		})
	}
	for _, b := range spec.InheritBindings {
		args := cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("inherit"),
		}
		if b.OldName != "" {
			args.OldName = pulumi.String(b.OldName)
		}
		if b.VersionId != "" {
			args.VersionId = pulumi.String(b.VersionId)
		}
		bindings = append(bindings, args)
	}
	for _, b := range spec.DataBlobs {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("data_blob"),
			Part: pulumi.String(b.Part),
		})
	}
	for _, b := range spec.TextBlobs {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("text_blob"),
			Part: pulumi.String(b.Part),
		})
	}
	for _, b := range spec.Browsers {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("browser"),
		})
	}
	for _, b := range spec.AiSearch {
		args := cloudfl.WorkersScriptBindingArgs{
			Name:         pulumi.String(b.Name),
			Type:         pulumi.String("ai_search"),
			InstanceName: pulumi.String(b.InstanceName),
		}
		if b.Namespace != "" {
			args.Namespace = pulumi.String(b.Namespace)
		}
		if b.AppId != "" {
			args.AppId = pulumi.String(b.AppId)
		}
		bindings = append(bindings, args)
	}
	for _, b := range spec.AiSearchNamespaces {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name:      pulumi.String(b.Name),
			Type:      pulumi.String("ai_search_namespace"),
			Namespace: pulumi.String(b.Namespace),
		})
	}
	for _, b := range spec.Images {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("images"),
		})
	}
	for _, b := range spec.Media {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("media"),
		})
	}
	for _, b := range spec.WasmModules {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("wasm_module"),
			Part: pulumi.String(b.Part),
		})
	}
	for _, b := range spec.VpcServices {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name:      pulumi.String(b.Name),
			Type:      pulumi.String("vpc_service"),
			ServiceId: pulumi.String(b.ServiceId),
		})
	}
	for _, b := range spec.VpcNetworks {
		args := cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(b.Name),
			Type: pulumi.String("vpc_network"),
		}
		if b.NetworkId != "" {
			args.NetworkId = pulumi.String(b.NetworkId)
		}
		if v := fkValue(b.TunnelId); v != "" {
			args.TunnelId = pulumi.String(v)
		}
		bindings = append(bindings, args)
	}
	for _, b := range spec.TailConsumerBindings {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name:    pulumi.String(b.Name),
			Type:    pulumi.String("tail_consumer"),
			Service: pulumi.String(fkValue(b.Service)),
		})
	}
	// Assets binding (env.<NAME>) for full-stack workers that read assets in code.
	if a := spec.Assets; a != nil && a.BindingName != "" {
		bindings = append(bindings, cloudfl.WorkersScriptBindingArgs{
			Name: pulumi.String(a.BindingName),
			Type: pulumi.String("assets"),
		})
	}

	return bindings
}
