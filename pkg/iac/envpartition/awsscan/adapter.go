// Package awsscan adapts AWS Cloud Control scan records -- a CloudFormation
// type name plus the type's JSON property document -- into the partition
// engine's neutral Resource values. All AWS property-model knowledge for
// partitioning lives here, declared once, so every consumer of the engine
// (the eval harness's reference proposer, the import journey's scan stage)
// reads the same signals the same way.
package awsscan

import (
	"github.com/plantonhq/planton/pkg/iac/envpartition"
)

// nameProperties are the property keys that carry a HUMAN-AUTHORED name,
// in the fixed order they are consulted after the Name tag. Only these may
// populate Resource.Name: opaque cloud identifiers are never used as a
// name surface, because a random id can embed a misleading token (a real
// account grew a security group id containing "prodfix").
var nameProperties = []string{
	"BucketName",
	"QueueName",
	"TopicName",
	"TableName",
	"RepositoryName",
	"RoleName",
	"GroupName",
}

// containerProperties are the property keys that carry the identifier of a
// CONTAINING resource (the surface containment inheritance walks). A
// resource carrying VpcId lives in that VPC; a NAT gateway's SubnetId puts
// it in that subnet.
var containerProperties = []string{
	"VpcId",
	"SubnetId",
}

// Adapt converts one scan record into the engine's input.
func Adapt(typeName, identifier string, properties map[string]any) envpartition.Resource {
	resource := envpartition.Resource{
		TypeName:   typeName,
		Identifier: identifier,
		Tags:       tagMap(properties),
	}
	resource.Name = resource.Tags["Name"]
	if resource.Name == "" {
		for _, key := range nameProperties {
			if name, _ := properties[key].(string); name != "" {
				resource.Name = name
				break
			}
		}
	}
	for _, key := range containerProperties {
		if id, _ := properties[key].(string); id != "" {
			resource.Containers = append(resource.Containers, id)
		}
	}
	// An internet gateway carries its VPC in the enriched Attachments list
	// rather than a top-level VpcId.
	for _, raw := range attachments(properties) {
		if vpcID, _ := raw["VpcId"].(string); vpcID != "" {
			resource.Containers = append(resource.Containers, vpcID)
		}
	}
	return resource
}

// tagMap flattens the Cloud Control Tags shape ([{Key, Value}]) into a map.
func tagMap(properties map[string]any) map[string]string {
	raw, _ := properties["Tags"].([]any)
	if len(raw) == 0 {
		return nil
	}
	tags := map[string]string{}
	for _, entry := range raw {
		tag, _ := entry.(map[string]any)
		if tag == nil {
			continue
		}
		key, _ := tag["Key"].(string)
		value, _ := tag["Value"].(string)
		if key != "" {
			tags[key] = value
		}
	}
	return tags
}

func attachments(properties map[string]any) []map[string]any {
	raw, _ := properties["Attachments"].([]any)
	var result []map[string]any
	for _, entry := range raw {
		if attachment, _ := entry.(map[string]any); attachment != nil {
			result = append(result, attachment)
		}
	}
	return result
}
