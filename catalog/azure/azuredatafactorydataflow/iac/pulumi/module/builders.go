package module

import (
	azuredatafactorydataflowv1alpha1 "github.com/plantonhq/planton/catalog/azure/azuredatafactorydataflow/v1alpha1"
	"github.com/pulumi/pulumi-azure/sdk/v6/go/azure/datafactory"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// The classic SDK generates a parallel nested-type set per resource
// (DataFlow* and FlowletDataFlow*), so each spec block gets twin
// builders in lockstep -- both write the exact same ARM wire shape.
// Rejected-row routing is set on sinks only, mirroring the spec (the
// provider silently drops it on sources).

func buildDataFlowSources(sources []*azuredatafactorydataflowv1alpha1.AzureDataFactoryDataFlowSource) datafactory.DataFlowSourceArray {
	result := datafactory.DataFlowSourceArray{}
	for _, source := range sources {
		args := datafactory.DataFlowSourceArgs{
			Name: pulumi.String(source.Name),
		}
		if source.Description != "" {
			args.Description = pulumi.String(source.Description)
		}
		if source.Dataset != nil {
			datasetArgs := datafactory.DataFlowSourceDatasetArgs{
				Name: pulumi.String(source.Dataset.Name.GetValue()),
			}
			if len(source.Dataset.Parameters) > 0 {
				datasetArgs.Parameters = pulumi.ToStringMap(source.Dataset.Parameters)
			}
			args.Dataset = datasetArgs
		}
		if source.Flowlet != nil {
			flowletArgs := datafactory.DataFlowSourceFlowletArgs{
				Name: pulumi.String(source.Flowlet.Name.GetValue()),
			}
			if len(source.Flowlet.Parameters) > 0 {
				flowletArgs.Parameters = pulumi.ToStringMap(source.Flowlet.Parameters)
			}
			if source.Flowlet.DatasetParameters != "" {
				flowletArgs.DatasetParameters = pulumi.String(source.Flowlet.DatasetParameters)
			}
			args.Flowlet = flowletArgs
		}
		if source.LinkedService != nil {
			linkedServiceArgs := datafactory.DataFlowSourceLinkedServiceArgs{
				Name: pulumi.String(source.LinkedService.Name.GetValue()),
			}
			if len(source.LinkedService.Parameters) > 0 {
				linkedServiceArgs.Parameters = pulumi.ToStringMap(source.LinkedService.Parameters)
			}
			args.LinkedService = linkedServiceArgs
		}
		if source.SchemaLinkedService != nil {
			schemaArgs := datafactory.DataFlowSourceSchemaLinkedServiceArgs{
				Name: pulumi.String(source.SchemaLinkedService.Name.GetValue()),
			}
			if len(source.SchemaLinkedService.Parameters) > 0 {
				schemaArgs.Parameters = pulumi.ToStringMap(source.SchemaLinkedService.Parameters)
			}
			args.SchemaLinkedService = schemaArgs
		}
		result = append(result, args)
	}
	return result
}

func buildDataFlowSinks(sinks []*azuredatafactorydataflowv1alpha1.AzureDataFactoryDataFlowSink) datafactory.DataFlowSinkArray {
	result := datafactory.DataFlowSinkArray{}
	for _, sink := range sinks {
		args := datafactory.DataFlowSinkArgs{
			Name: pulumi.String(sink.Name),
		}
		if sink.Description != "" {
			args.Description = pulumi.String(sink.Description)
		}
		if sink.Dataset != nil {
			datasetArgs := datafactory.DataFlowSinkDatasetArgs{
				Name: pulumi.String(sink.Dataset.Name.GetValue()),
			}
			if len(sink.Dataset.Parameters) > 0 {
				datasetArgs.Parameters = pulumi.ToStringMap(sink.Dataset.Parameters)
			}
			args.Dataset = datasetArgs
		}
		if sink.Flowlet != nil {
			flowletArgs := datafactory.DataFlowSinkFlowletArgs{
				Name: pulumi.String(sink.Flowlet.Name.GetValue()),
			}
			if len(sink.Flowlet.Parameters) > 0 {
				flowletArgs.Parameters = pulumi.ToStringMap(sink.Flowlet.Parameters)
			}
			if sink.Flowlet.DatasetParameters != "" {
				flowletArgs.DatasetParameters = pulumi.String(sink.Flowlet.DatasetParameters)
			}
			args.Flowlet = flowletArgs
		}
		if sink.LinkedService != nil {
			linkedServiceArgs := datafactory.DataFlowSinkLinkedServiceArgs{
				Name: pulumi.String(sink.LinkedService.Name.GetValue()),
			}
			if len(sink.LinkedService.Parameters) > 0 {
				linkedServiceArgs.Parameters = pulumi.ToStringMap(sink.LinkedService.Parameters)
			}
			args.LinkedService = linkedServiceArgs
		}
		if sink.SchemaLinkedService != nil {
			schemaArgs := datafactory.DataFlowSinkSchemaLinkedServiceArgs{
				Name: pulumi.String(sink.SchemaLinkedService.Name.GetValue()),
			}
			if len(sink.SchemaLinkedService.Parameters) > 0 {
				schemaArgs.Parameters = pulumi.ToStringMap(sink.SchemaLinkedService.Parameters)
			}
			args.SchemaLinkedService = schemaArgs
		}
		if sink.RejectedLinkedService != nil {
			rejectedArgs := datafactory.DataFlowSinkRejectedLinkedServiceArgs{
				Name: pulumi.String(sink.RejectedLinkedService.Name.GetValue()),
			}
			if len(sink.RejectedLinkedService.Parameters) > 0 {
				rejectedArgs.Parameters = pulumi.ToStringMap(sink.RejectedLinkedService.Parameters)
			}
			args.RejectedLinkedService = rejectedArgs
		}
		result = append(result, args)
	}
	return result
}

func buildDataFlowTransformations(transformations []*azuredatafactorydataflowv1alpha1.AzureDataFactoryDataFlowTransformation) datafactory.DataFlowTransformationArray {
	result := datafactory.DataFlowTransformationArray{}
	for _, transformation := range transformations {
		args := datafactory.DataFlowTransformationArgs{
			Name: pulumi.String(transformation.Name),
		}
		if transformation.Description != "" {
			args.Description = pulumi.String(transformation.Description)
		}
		if transformation.Dataset != nil {
			datasetArgs := datafactory.DataFlowTransformationDatasetArgs{
				Name: pulumi.String(transformation.Dataset.Name.GetValue()),
			}
			if len(transformation.Dataset.Parameters) > 0 {
				datasetArgs.Parameters = pulumi.ToStringMap(transformation.Dataset.Parameters)
			}
			args.Dataset = datasetArgs
		}
		if transformation.Flowlet != nil {
			flowletArgs := datafactory.DataFlowTransformationFlowletArgs{
				Name: pulumi.String(transformation.Flowlet.Name.GetValue()),
			}
			if len(transformation.Flowlet.Parameters) > 0 {
				flowletArgs.Parameters = pulumi.ToStringMap(transformation.Flowlet.Parameters)
			}
			if transformation.Flowlet.DatasetParameters != "" {
				flowletArgs.DatasetParameters = pulumi.String(transformation.Flowlet.DatasetParameters)
			}
			args.Flowlet = flowletArgs
		}
		if transformation.LinkedService != nil {
			linkedServiceArgs := datafactory.DataFlowTransformationLinkedServiceArgs{
				Name: pulumi.String(transformation.LinkedService.Name.GetValue()),
			}
			if len(transformation.LinkedService.Parameters) > 0 {
				linkedServiceArgs.Parameters = pulumi.ToStringMap(transformation.LinkedService.Parameters)
			}
			args.LinkedService = linkedServiceArgs
		}
		result = append(result, args)
	}
	return result
}

func buildFlowletSources(sources []*azuredatafactorydataflowv1alpha1.AzureDataFactoryDataFlowSource) datafactory.FlowletDataFlowSourceArray {
	result := datafactory.FlowletDataFlowSourceArray{}
	for _, source := range sources {
		args := datafactory.FlowletDataFlowSourceArgs{
			Name: pulumi.String(source.Name),
		}
		if source.Description != "" {
			args.Description = pulumi.String(source.Description)
		}
		if source.Dataset != nil {
			datasetArgs := datafactory.FlowletDataFlowSourceDatasetArgs{
				Name: pulumi.String(source.Dataset.Name.GetValue()),
			}
			if len(source.Dataset.Parameters) > 0 {
				datasetArgs.Parameters = pulumi.ToStringMap(source.Dataset.Parameters)
			}
			args.Dataset = datasetArgs
		}
		if source.Flowlet != nil {
			flowletArgs := datafactory.FlowletDataFlowSourceFlowletArgs{
				Name: pulumi.String(source.Flowlet.Name.GetValue()),
			}
			if len(source.Flowlet.Parameters) > 0 {
				flowletArgs.Parameters = pulumi.ToStringMap(source.Flowlet.Parameters)
			}
			if source.Flowlet.DatasetParameters != "" {
				flowletArgs.DatasetParameters = pulumi.String(source.Flowlet.DatasetParameters)
			}
			args.Flowlet = flowletArgs
		}
		if source.LinkedService != nil {
			linkedServiceArgs := datafactory.FlowletDataFlowSourceLinkedServiceArgs{
				Name: pulumi.String(source.LinkedService.Name.GetValue()),
			}
			if len(source.LinkedService.Parameters) > 0 {
				linkedServiceArgs.Parameters = pulumi.ToStringMap(source.LinkedService.Parameters)
			}
			args.LinkedService = linkedServiceArgs
		}
		if source.SchemaLinkedService != nil {
			schemaArgs := datafactory.FlowletDataFlowSourceSchemaLinkedServiceArgs{
				Name: pulumi.String(source.SchemaLinkedService.Name.GetValue()),
			}
			if len(source.SchemaLinkedService.Parameters) > 0 {
				schemaArgs.Parameters = pulumi.ToStringMap(source.SchemaLinkedService.Parameters)
			}
			args.SchemaLinkedService = schemaArgs
		}
		result = append(result, args)
	}
	return result
}

func buildFlowletSinks(sinks []*azuredatafactorydataflowv1alpha1.AzureDataFactoryDataFlowSink) datafactory.FlowletDataFlowSinkArray {
	result := datafactory.FlowletDataFlowSinkArray{}
	for _, sink := range sinks {
		args := datafactory.FlowletDataFlowSinkArgs{
			Name: pulumi.String(sink.Name),
		}
		if sink.Description != "" {
			args.Description = pulumi.String(sink.Description)
		}
		if sink.Dataset != nil {
			datasetArgs := datafactory.FlowletDataFlowSinkDatasetArgs{
				Name: pulumi.String(sink.Dataset.Name.GetValue()),
			}
			if len(sink.Dataset.Parameters) > 0 {
				datasetArgs.Parameters = pulumi.ToStringMap(sink.Dataset.Parameters)
			}
			args.Dataset = datasetArgs
		}
		if sink.Flowlet != nil {
			flowletArgs := datafactory.FlowletDataFlowSinkFlowletArgs{
				Name: pulumi.String(sink.Flowlet.Name.GetValue()),
			}
			if len(sink.Flowlet.Parameters) > 0 {
				flowletArgs.Parameters = pulumi.ToStringMap(sink.Flowlet.Parameters)
			}
			if sink.Flowlet.DatasetParameters != "" {
				flowletArgs.DatasetParameters = pulumi.String(sink.Flowlet.DatasetParameters)
			}
			args.Flowlet = flowletArgs
		}
		if sink.LinkedService != nil {
			linkedServiceArgs := datafactory.FlowletDataFlowSinkLinkedServiceArgs{
				Name: pulumi.String(sink.LinkedService.Name.GetValue()),
			}
			if len(sink.LinkedService.Parameters) > 0 {
				linkedServiceArgs.Parameters = pulumi.ToStringMap(sink.LinkedService.Parameters)
			}
			args.LinkedService = linkedServiceArgs
		}
		if sink.SchemaLinkedService != nil {
			schemaArgs := datafactory.FlowletDataFlowSinkSchemaLinkedServiceArgs{
				Name: pulumi.String(sink.SchemaLinkedService.Name.GetValue()),
			}
			if len(sink.SchemaLinkedService.Parameters) > 0 {
				schemaArgs.Parameters = pulumi.ToStringMap(sink.SchemaLinkedService.Parameters)
			}
			args.SchemaLinkedService = schemaArgs
		}
		if sink.RejectedLinkedService != nil {
			rejectedArgs := datafactory.FlowletDataFlowSinkRejectedLinkedServiceArgs{
				Name: pulumi.String(sink.RejectedLinkedService.Name.GetValue()),
			}
			if len(sink.RejectedLinkedService.Parameters) > 0 {
				rejectedArgs.Parameters = pulumi.ToStringMap(sink.RejectedLinkedService.Parameters)
			}
			args.RejectedLinkedService = rejectedArgs
		}
		result = append(result, args)
	}
	return result
}

func buildFlowletTransformations(transformations []*azuredatafactorydataflowv1alpha1.AzureDataFactoryDataFlowTransformation) datafactory.FlowletDataFlowTransformationArray {
	result := datafactory.FlowletDataFlowTransformationArray{}
	for _, transformation := range transformations {
		args := datafactory.FlowletDataFlowTransformationArgs{
			Name: pulumi.String(transformation.Name),
		}
		if transformation.Description != "" {
			args.Description = pulumi.String(transformation.Description)
		}
		if transformation.Dataset != nil {
			datasetArgs := datafactory.FlowletDataFlowTransformationDatasetArgs{
				Name: pulumi.String(transformation.Dataset.Name.GetValue()),
			}
			if len(transformation.Dataset.Parameters) > 0 {
				datasetArgs.Parameters = pulumi.ToStringMap(transformation.Dataset.Parameters)
			}
			args.Dataset = datasetArgs
		}
		if transformation.Flowlet != nil {
			flowletArgs := datafactory.FlowletDataFlowTransformationFlowletArgs{
				Name: pulumi.String(transformation.Flowlet.Name.GetValue()),
			}
			if len(transformation.Flowlet.Parameters) > 0 {
				flowletArgs.Parameters = pulumi.ToStringMap(transformation.Flowlet.Parameters)
			}
			if transformation.Flowlet.DatasetParameters != "" {
				flowletArgs.DatasetParameters = pulumi.String(transformation.Flowlet.DatasetParameters)
			}
			args.Flowlet = flowletArgs
		}
		if transformation.LinkedService != nil {
			linkedServiceArgs := datafactory.FlowletDataFlowTransformationLinkedServiceArgs{
				Name: pulumi.String(transformation.LinkedService.Name.GetValue()),
			}
			if len(transformation.LinkedService.Parameters) > 0 {
				linkedServiceArgs.Parameters = pulumi.ToStringMap(transformation.LinkedService.Parameters)
			}
			args.LinkedService = linkedServiceArgs
		}
		result = append(result, args)
	}
	return result
}
