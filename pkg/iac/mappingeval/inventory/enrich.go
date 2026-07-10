package inventory

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/iac/mappingeval"
)

// Declared per-type enrichments: Cloud Control's models are complete for
// most types but not all, and the gaps are closed HERE, per type, with
// typed read-only SDK calls -- never ad hoc at consumption sites. This is
// the same architecture the platform's inventory capability uses (its
// first entry likewise fills S3 bucket regions), so a proposer developed
// against this scanner sees the same completed data the product's scan
// channel serves.
//
// Current entries and why each exists:
//
//   - AWS::S3::Bucket gains "Region" (GetBucketLocation): Cloud Control's
//     bucket listing is global and region-less, and region is both the
//     scan-scope filter and a required spec field.
//   - AWS::EC2::RouteTable gains "Routes" and "Associations"
//     (DescribeRouteTables): CloudFormation models routes and associations
//     as separate resources, so the route table's Cloud Control model
//     carries neither -- yet they are exactly the ownership and routing
//     facts grouping and spec reconstruction need.
//   - AWS::EC2::SubnetRouteTableAssociation gains "RouteTableId" and
//     "SubnetId" (from the same DescribeRouteTables read): its list
//     entries carry only the association id, and GetResource on it fails
//     for main-table associations -- the describe is the reliable source.
//   - AWS::EC2::InternetGateway gains "Attachments"
//     (DescribeInternetGateways): the Cloud Control model omits the VPC
//     attachment, which is the gateway's one relationship.
//   - AWS::SNS::Topic gains "Policy" (GetTopicAttributes): CloudFormation
//     models the access policy as a separate TopicPolicy resource, so the
//     topic's Cloud Control model omits it -- yet it is a real setting the
//     topic carries.
//
// INVARIANT: enriched values are stored in JSON-generic shapes (map/slice
// of `any`, float64 numbers) -- byte-identical to what decoding Cloud
// Control's own property documents produces -- so consumers and recorded
// fixtures see one shape regardless of whether a property came from Cloud
// Control or an enrichment. A Go-typed slice here would silently fail the
// consumers' JSON-generic type assertions on live scans while passing on
// fixtures.

const (
	s3BucketTypeName        = "AWS::S3::Bucket"
	s3RegionProperty        = "Region"
	routeTableTypeName      = "AWS::EC2::RouteTable"
	rtbAssociationTypeName  = "AWS::EC2::SubnetRouteTableAssociation"
	internetGatewayTypeName = "AWS::EC2::InternetGateway"
	snsTopicTypeName        = "AWS::SNS::Topic"
)

// applyEnrichments runs every declared enrichment whose type appears in
// the scan. Enrichment failures fail the scan loudly: unlike a single
// entry's degraded get, a missing enrichment silently changes what every
// consumer sees, and the entries exist precisely because the base model is
// insufficient.
func applyEnrichments(ctx context.Context, cfg aws.Config, region string, scan *mappingeval.Scan) error {
	if hasType(scan, s3BucketTypeName) {
		if err := enrichBucketRegions(ctx, cfg, scan); err != nil {
			return errors.Wrap(err, "enriching S3 bucket regions")
		}
	}
	if hasType(scan, routeTableTypeName) || hasType(scan, rtbAssociationTypeName) {
		if err := enrichRouteTables(ctx, cfg, region, scan); err != nil {
			return errors.Wrap(err, "enriching route tables/associations")
		}
	}
	if hasType(scan, internetGatewayTypeName) {
		if err := enrichInternetGateways(ctx, cfg, region, scan); err != nil {
			return errors.Wrap(err, "enriching internet gateways")
		}
	}
	if hasType(scan, snsTopicTypeName) {
		if err := enrichTopicPolicies(ctx, cfg, region, scan); err != nil {
			return errors.Wrap(err, "enriching SNS topic policies")
		}
	}
	return nil
}

func enrichBucketRegions(ctx context.Context, cfg aws.Config, scan *mappingeval.Scan) error {
	client := s3.NewFromConfig(cfg)
	for i := range scan.Resources {
		r := &scan.Resources[i]
		if r.TypeName != s3BucketTypeName {
			continue
		}
		out, err := client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{Bucket: aws.String(r.Identifier)})
		if err != nil {
			// Per-bucket degradation: a bucket we cannot locate (e.g. a
			// permission-scoped one) keeps its entry without a Region and
			// falls out of a region-scoped scan.
			continue
		}
		region := string(out.LocationConstraint)
		if region == "" {
			// The S3 API's legacy encoding: an empty constraint IS
			// us-east-1.
			region = "us-east-1"
		}
		r.Properties[s3RegionProperty] = region
	}
	return nil
}

func enrichRouteTables(ctx context.Context, cfg aws.Config, region string, scan *mappingeval.Scan) error {
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) { o.Region = region })
	type associationFacts struct{ routeTableID, subnetID string }
	// []any per the package invariant: JSON-generic shapes only.
	routesByTable := map[string][]any{}
	associationsByTable := map[string][]any{}
	factsByAssociation := map[string]associationFacts{}

	paginator := ec2.NewDescribeRouteTablesPaginator(client, &ec2.DescribeRouteTablesInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, table := range page.RouteTables {
			tableID := aws.ToString(table.RouteTableId)
			for _, route := range table.Routes {
				entry := map[string]any{}
				setIfPresent(entry, "DestinationCidrBlock", route.DestinationCidrBlock)
				setIfPresent(entry, "DestinationIpv6CidrBlock", route.DestinationIpv6CidrBlock)
				setIfPresent(entry, "DestinationPrefixListId", route.DestinationPrefixListId)
				setIfPresent(entry, "GatewayId", route.GatewayId)
				setIfPresent(entry, "NatGatewayId", route.NatGatewayId)
				setIfPresent(entry, "TransitGatewayId", route.TransitGatewayId)
				setIfPresent(entry, "VpcPeeringConnectionId", route.VpcPeeringConnectionId)
				setIfPresent(entry, "NetworkInterfaceId", route.NetworkInterfaceId)
				setIfPresent(entry, "EgressOnlyInternetGatewayId", route.EgressOnlyInternetGatewayId)
				routesByTable[tableID] = append(routesByTable[tableID], entry)
			}
			for _, association := range table.Associations {
				associationID := aws.ToString(association.RouteTableAssociationId)
				entry := map[string]any{"RouteTableAssociationId": associationID}
				setIfPresent(entry, "SubnetId", association.SubnetId)
				if association.Main != nil {
					entry["Main"] = *association.Main
				}
				associationsByTable[tableID] = append(associationsByTable[tableID], entry)
				factsByAssociation[associationID] = associationFacts{
					routeTableID: tableID,
					subnetID:     aws.ToString(association.SubnetId),
				}
			}
		}
	}

	for i := range scan.Resources {
		r := &scan.Resources[i]
		switch r.TypeName {
		case routeTableTypeName:
			if routes, ok := routesByTable[r.Identifier]; ok {
				r.Properties["Routes"] = routes
			}
			if associations, ok := associationsByTable[r.Identifier]; ok {
				r.Properties["Associations"] = associations
			}
		case rtbAssociationTypeName:
			if facts, ok := factsByAssociation[r.Identifier]; ok {
				r.Properties["RouteTableId"] = facts.routeTableID
				if facts.subnetID != "" {
					r.Properties["SubnetId"] = facts.subnetID
				}
			}
		}
	}
	return nil
}

func enrichInternetGateways(ctx context.Context, cfg aws.Config, region string, scan *mappingeval.Scan) error {
	client := ec2.NewFromConfig(cfg, func(o *ec2.Options) { o.Region = region })
	// []any per the package invariant: JSON-generic shapes only.
	attachmentsByGateway := map[string][]any{}

	paginator := ec2.NewDescribeInternetGatewaysPaginator(client, &ec2.DescribeInternetGatewaysInput{})
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}
		for _, gateway := range page.InternetGateways {
			gatewayID := aws.ToString(gateway.InternetGatewayId)
			for _, attachment := range gateway.Attachments {
				entry := map[string]any{}
				setIfPresent(entry, "VpcId", attachment.VpcId)
				attachmentsByGateway[gatewayID] = append(attachmentsByGateway[gatewayID], entry)
			}
		}
	}

	for i := range scan.Resources {
		r := &scan.Resources[i]
		if r.TypeName != internetGatewayTypeName {
			continue
		}
		if attachments, ok := attachmentsByGateway[r.Identifier]; ok {
			r.Properties["Attachments"] = attachments
		}
	}
	return nil
}

// enrichTopicPolicies fills each topic's access policy from
// GetTopicAttributes, parsed into the same JSON-generic document shape a
// policy inside a Cloud Control property document would have.
func enrichTopicPolicies(ctx context.Context, cfg aws.Config, region string, scan *mappingeval.Scan) error {
	client := sns.NewFromConfig(cfg, func(o *sns.Options) { o.Region = region })
	for i := range scan.Resources {
		r := &scan.Resources[i]
		if r.TypeName != snsTopicTypeName {
			continue
		}
		out, err := client.GetTopicAttributes(ctx, &sns.GetTopicAttributesInput{TopicArn: aws.String(r.Identifier)})
		if err != nil {
			// Per-topic degradation: an unreadable topic keeps its entry
			// without a policy.
			continue
		}
		policyJSON := out.Attributes["Policy"]
		if policyJSON == "" {
			continue
		}
		var policy map[string]any
		if err := json.Unmarshal([]byte(policyJSON), &policy); err != nil {
			continue
		}
		r.Properties["Policy"] = policy
	}
	return nil
}

func hasType(scan *mappingeval.Scan, typeName string) bool {
	for _, r := range scan.Resources {
		if r.TypeName == typeName {
			return true
		}
	}
	return false
}

func setIfPresent(entry map[string]any, key string, value *string) {
	if value != nil && *value != "" {
		entry[key] = *value
	}
}
