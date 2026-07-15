// Package inventory is the mapping-eval harness's read-only account
// scanner: the blind side's only window into the cloud. It mirrors the
// shape of Planton's inventory capability -- Cloud Control list/get keyed
// by CloudFormation type names, each resource reported as (type name,
// primary identifier, property document), with declared per-type
// enrichments closing Cloud Control's model gaps through typed SDK reads.
//
// Read-only is STRUCTURAL: this package only ever constructs Cloud
// Control's read calls and the enrichment registry's describe/get reads.
// There is no code path that can mutate an account, which is the same
// guarantee the platform's inventory surface makes -- a proposer fed by
// this scanner physically cannot touch what it examines.
package inventory

import (
	"context"
	"strings"

	"encoding/json"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudcontrol"
	"github.com/pkg/errors"
	"github.com/plantonhq/planton/pkg/iac/mappingeval"
)

// Scanner lists and reads account resources through Cloud Control.
type Scanner struct {
	client *cloudcontrol.Client
	cfg    aws.Config
}

// NewScanner builds a scanner over an ambient-credential AWS config.
func NewScanner(cfg aws.Config) *Scanner {
	return &Scanner{client: cloudcontrol.NewFromConfig(cfg), cfg: cfg}
}

// Scan lists every declared type in the region and returns the full
// property model per resource: Cloud Control's list responses carry
// type-dependent property subsets, so each entry is completed by a
// GetResource read, then by the declared enrichments. Per-resource read
// failures degrade that entry honestly (list-level properties remain); a
// type that fails to LIST at all fails the scan -- the type allowlist is
// declared, so an unlistable entry is a declaration bug, not a runtime
// condition to paper over.
func (s *Scanner) Scan(ctx context.Context, region string, typeNames []string) (*mappingeval.Scan, error) {
	scan := &mappingeval.Scan{Region: region}
	for _, typeName := range typeNames {
		resources, err := s.listType(ctx, region, typeName)
		if err != nil {
			return nil, errors.Wrapf(err, "listing %s in %s", typeName, region)
		}
		scan.Resources = append(scan.Resources, resources...)
	}
	if err := applyEnrichments(ctx, s.cfg, region, scan); err != nil {
		return nil, err
	}
	// S3 bucket listings are global; a single-region scan honestly means
	// "what exists in this region", so buckets whose enriched Region
	// disagrees leave the scan.
	scan.Resources = filterBucketsToRegion(scan.Resources, region)
	return scan, nil
}

// listType pages through ListResources for one type and completes each
// entry with its full model.
func (s *Scanner) listType(ctx context.Context, region string, typeName string) ([]mappingeval.ScannedResource, error) {
	var resources []mappingeval.ScannedResource
	var nextToken *string
	for {
		out, err := s.client.ListResources(ctx, &cloudcontrol.ListResourcesInput{
			TypeName:  aws.String(typeName),
			NextToken: nextToken,
		}, func(o *cloudcontrol.Options) { o.Region = region })
		if err != nil {
			return nil, err
		}
		for _, description := range out.ResourceDescriptions {
			resource := mappingeval.ScannedResource{
				TypeName:   typeName,
				Identifier: aws.ToString(description.Identifier),
				Properties: decodeProperties(aws.ToString(description.Properties)),
			}
			// The list response's property subset varies by type; the get
			// returns the full model. A get failure keeps the list-level
			// entry -- existence must never be hidden by a read hiccup.
			full, err := s.client.GetResource(ctx, &cloudcontrol.GetResourceInput{
				TypeName:   aws.String(typeName),
				Identifier: description.Identifier,
			}, func(o *cloudcontrol.Options) { o.Region = region })
			if err == nil && full.ResourceDescription != nil {
				if props := decodeProperties(aws.ToString(full.ResourceDescription.Properties)); len(props) > 0 {
					resource.Properties = props
				}
			}
			resources = append(resources, resource)
		}
		if out.NextToken == nil || aws.ToString(out.NextToken) == "" {
			return resources, nil
		}
		nextToken = out.NextToken
	}
}

// decodeProperties parses Cloud Control's JSON property document; a
// malformed document yields an empty map rather than hiding the resource.
func decodeProperties(doc string) map[string]any {
	if doc == "" {
		return map[string]any{}
	}
	var props map[string]any
	if err := json.Unmarshal([]byte(doc), &props); err != nil {
		return map[string]any{}
	}
	return props
}

// filterBucketsToRegion drops S3 buckets outside the scanned region.
func filterBucketsToRegion(resources []mappingeval.ScannedResource, region string) []mappingeval.ScannedResource {
	filtered := resources[:0]
	for _, r := range resources {
		if r.TypeName == s3BucketTypeName {
			bucketRegion, _ := r.Properties[s3RegionProperty].(string)
			if !strings.EqualFold(bucketRegion, region) {
				continue
			}
		}
		filtered = append(filtered, r)
	}
	return filtered
}
