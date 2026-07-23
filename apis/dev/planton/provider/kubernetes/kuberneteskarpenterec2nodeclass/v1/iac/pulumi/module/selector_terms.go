// Selector-term builders for the EC2NodeClass spec: AMI, subnet,
// security-group and capacity-reservation terms. Terms are ORed by
// Karpenter; fields within a term are ANDed. Every field is emitted only
// when present so mutually-exclusive arms (id vs tags vs alias, enforced by
// protovalidate CEL) render exactly one key.
package module

import (
	kuberneteskarpenterec2nodeclassv1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kuberneteskarpenterec2nodeclass/v1"
	karpenterv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/karpenter/kubernetes/karpenter/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func buildAmiSelectorTerms(terms []*kuberneteskarpenterec2nodeclassv1.KubernetesKarpenterEc2NodeClassAmiSelectorTerm) karpenterv1.EC2NodeClassSpecAmiSelectorTermsArray {
	arr := karpenterv1.EC2NodeClassSpecAmiSelectorTermsArray{}
	for _, term := range terms {
		termArgs := karpenterv1.EC2NodeClassSpecAmiSelectorTermsArgs{}
		if term.Alias != nil {
			termArgs.Alias = pulumi.String(term.GetAlias())
		}
		if term.Id != nil {
			termArgs.Id = pulumi.String(term.GetId())
		}
		if term.Name != nil {
			termArgs.Name = pulumi.String(term.GetName())
		}
		if term.Owner != nil {
			termArgs.Owner = pulumi.String(term.GetOwner())
		}
		if term.SsmParameter != nil {
			termArgs.SsmParameter = pulumi.String(term.GetSsmParameter())
		}
		if tags := term.GetTags(); len(tags) > 0 {
			termArgs.Tags = pulumi.ToStringMap(tags)
		}
		arr = append(arr, termArgs)
	}
	return arr
}

func buildSubnetSelectorTerms(terms []*kuberneteskarpenterec2nodeclassv1.KubernetesKarpenterEc2NodeClassSubnetSelectorTerm) karpenterv1.EC2NodeClassSpecSubnetSelectorTermsArray {
	arr := karpenterv1.EC2NodeClassSpecSubnetSelectorTermsArray{}
	for _, term := range terms {
		termArgs := karpenterv1.EC2NodeClassSpecSubnetSelectorTermsArgs{}
		if term.Id != nil {
			termArgs.Id = pulumi.String(term.GetId())
		}
		if tags := term.GetTags(); len(tags) > 0 {
			termArgs.Tags = pulumi.ToStringMap(tags)
		}
		arr = append(arr, termArgs)
	}
	return arr
}

func buildSecurityGroupSelectorTerms(terms []*kuberneteskarpenterec2nodeclassv1.KubernetesKarpenterEc2NodeClassSecurityGroupSelectorTerm) karpenterv1.EC2NodeClassSpecSecurityGroupSelectorTermsArray {
	arr := karpenterv1.EC2NodeClassSpecSecurityGroupSelectorTermsArray{}
	for _, term := range terms {
		termArgs := karpenterv1.EC2NodeClassSpecSecurityGroupSelectorTermsArgs{}
		if term.Id != nil {
			termArgs.Id = pulumi.String(term.GetId())
		}
		if term.Name != nil {
			termArgs.Name = pulumi.String(term.GetName())
		}
		if tags := term.GetTags(); len(tags) > 0 {
			termArgs.Tags = pulumi.ToStringMap(tags)
		}
		arr = append(arr, termArgs)
	}
	return arr
}

// buildCapacityReservationSelectorTerms maps the capacity-reservation terms.
// owner_id renders as the acronym-cased 'ownerID' key — the SDK's OwnerID
// arg carries that pulumi tag, matching the CRD.
func buildCapacityReservationSelectorTerms(terms []*kuberneteskarpenterec2nodeclassv1.KubernetesKarpenterEc2NodeClassCapacityReservationSelectorTerm) karpenterv1.EC2NodeClassSpecCapacityReservationSelectorTermsArray {
	arr := karpenterv1.EC2NodeClassSpecCapacityReservationSelectorTermsArray{}
	for _, term := range terms {
		termArgs := karpenterv1.EC2NodeClassSpecCapacityReservationSelectorTermsArgs{}
		if term.Id != nil {
			termArgs.Id = pulumi.String(term.GetId())
		}
		if term.InstanceMatchCriteria != nil {
			termArgs.InstanceMatchCriteria = pulumi.String(term.GetInstanceMatchCriteria())
		}
		if term.OwnerId != nil {
			termArgs.OwnerID = pulumi.String(term.GetOwnerId())
		}
		if tags := term.GetTags(); len(tags) > 0 {
			termArgs.Tags = pulumi.ToStringMap(tags)
		}
		arr = append(arr, termArgs)
	}
	return arr
}
