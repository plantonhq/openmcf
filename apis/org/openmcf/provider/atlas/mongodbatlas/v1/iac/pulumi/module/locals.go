package module

import (
	mongodbatlasv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/atlas/mongodbatlas/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type Locals struct {
	MongodbAtlas *mongodbatlasv1.MongodbAtlas
	// Computed values for cluster configuration
	ClusterName        string
	ProjectId          string
	ClusterType        string
	ElectableNodes     int32
	Priority           int32
	ReadOnlyNodes      int32
	CloudBackup        bool
	AutoScalingEnabled bool
	MongoDBVersion     string
	ProviderName       string
	InstanceSize       string
}

func initializeLocals(ctx *pulumi.Context, stackInput *mongodbatlasv1.MongodbAtlasStackInput) *Locals {
	locals := &Locals{}

	// Assign value for the locals variable to make it available across the project
	locals.MongodbAtlas = stackInput.Target

	// Extract cluster configuration for easy access
	spec := stackInput.Target.Spec
	clusterConfig := spec.ClusterConfig

	// Compute local values from spec
	locals.ClusterName = stackInput.Target.Metadata.Name
	locals.ProjectId = clusterConfig.ProjectId
	locals.ClusterType = clusterConfig.ClusterType
	locals.ElectableNodes = clusterConfig.ElectableNodes
	locals.Priority = clusterConfig.Priority
	locals.ReadOnlyNodes = clusterConfig.ReadOnlyNodes
	locals.CloudBackup = clusterConfig.CloudBackup
	locals.AutoScalingEnabled = clusterConfig.AutoScalingDiskGbEnabled
	locals.MongoDBVersion = clusterConfig.MongoDbMajorVersion
	locals.ProviderName = clusterConfig.ProviderName
	locals.InstanceSize = clusterConfig.ProviderInstanceSizeName

	return locals
}
