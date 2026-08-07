package module

import (
	"encoding/json"
	"reflect"
	"testing"

	awswafwebaclv1alpha1 "github.com/plantonhq/planton/catalog/aws/awswafwebacl/v1alpha1"
	foreignkeyv1 "github.com/plantonhq/planton/shared/foreignkey/v1"
)

// TestBuildRulesJSON_MatchesWafApiShape locks the statement-tree serializer to
// the AWS WAFv2 API JSON shape. The Terraform module renders the SAME shape
// from the same spec (iac/tf/locals.tf) -- the expected document below is the
// cross-engine contract for a spec exercising a managed group with a NOT
// scope-down, CUSTOM_KEYS rate aggregation, and an OR of byte-match + SQLi.
func TestBuildRulesJSON_MatchesWafApiShape(t *testing.T) {
	spec := &awswafwebaclv1alpha1.AwsWafWebAclSpec{
		Region: "us-west-2",
		Scope:  "REGIONAL",
		DefaultAction: &awswafwebaclv1alpha1.AwsWafWebAclDefaultAction{
			Type: "allow",
		},
		Rules: []*awswafwebaclv1alpha1.AwsWafWebAclRule{
			{
				Name:           "aws-managed-common",
				Priority:       1,
				OverrideAction: "none",
				Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement{
					Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement_ManagedRuleGroup{
						ManagedRuleGroup: &awswafwebaclv1alpha1.AwsWafWebAclManagedRuleGroupStatement{
							Name:       "AWSManagedRulesCommonRuleSet",
							VendorName: "AWS",
							ScopeDownStatement: &awswafwebaclv1alpha1.AwsWafWebAclStatement{
								Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement_NotStatement{
									NotStatement: &awswafwebaclv1alpha1.AwsWafWebAclNotStatement{
										Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement{
											Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement_ByteMatch{
												ByteMatch: &awswafwebaclv1alpha1.AwsWafWebAclByteMatchStatement{
													SearchString:         "/health",
													PositionalConstraint: "STARTS_WITH",
													FieldToMatch: &awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch{
														Field: &awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_UriPath{UriPath: true},
													},
													TextTransformations: []*awswafwebaclv1alpha1.AwsWafWebAclTextTransformation{
														{Priority: 0, Type: "NONE"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				Name:     "rate-limit-per-api-key",
				Priority: 2,
				Action:   "block",
				Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement{
					Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement_RateBased{
						RateBased: &awswafwebaclv1alpha1.AwsWafWebAclRateBasedStatement{
							Limit:            2000,
							AggregateKeyType: "CUSTOM_KEYS",
							CustomKeys: []*awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey{
								{Key: &awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_Header{
									Header: &awswafwebaclv1alpha1.AwsWafWebAclKeyWithTransformations{
										Name: "x-api-key",
										TextTransformations: []*awswafwebaclv1alpha1.AwsWafWebAclTextTransformation{
											{Priority: 0, Type: "NONE"},
										},
									},
								}},
								{Key: &awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_Ip{Ip: true}},
							},
						},
					},
				},
			},
			{
				Name:     "block-bad-paths",
				Priority: 3,
				Action:   "block",
				CustomResponse: &awswafwebaclv1alpha1.AwsWafWebAclCustomResponse{
					ResponseCode: 403,
				},
				Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement{
					Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement_OrStatement{
						OrStatement: &awswafwebaclv1alpha1.AwsWafWebAclOrStatement{
							Statements: []*awswafwebaclv1alpha1.AwsWafWebAclStatement{
								{Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement_ByteMatch{
									ByteMatch: &awswafwebaclv1alpha1.AwsWafWebAclByteMatchStatement{
										SearchString:         "/wp-admin",
										PositionalConstraint: "STARTS_WITH",
										FieldToMatch: &awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch{
											Field: &awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_UriPath{UriPath: true},
										},
										TextTransformations: []*awswafwebaclv1alpha1.AwsWafWebAclTextTransformation{
											{Priority: 0, Type: "LOWERCASE"},
										},
									},
								}},
								{Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement_SqliMatch{
									SqliMatch: &awswafwebaclv1alpha1.AwsWafWebAclSqliMatchStatement{
										FieldToMatch: &awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch{
											Field: &awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_Body{
												Body: &awswafwebaclv1alpha1.AwsWafWebAclBodyMatch{OversizeHandling: "MATCH"},
											},
										},
										TextTransformations: []*awswafwebaclv1alpha1.AwsWafWebAclTextTransformation{
											{Priority: 0, Type: "URL_DECODE"},
										},
									},
								}},
							},
						},
					},
				},
			},
		},
	}

	got, err := buildRulesJSON(spec)
	if err != nil {
		t.Fatalf("buildRulesJSON: %v", err)
	}

	// The expected document, exactly as the Terraform twin renders it from
	// the same manifest (verified via the module's rule_json local).
	expected := `[
	  {
	    "Name": "aws-managed-common",
	    "Priority": 1,
	    "OverrideAction": {"None": {}},
	    "Statement": {
	      "ManagedRuleGroupStatement": {
	        "Name": "AWSManagedRulesCommonRuleSet",
	        "VendorName": "AWS",
	        "ScopeDownStatement": {
	          "NotStatement": {
	            "Statement": {
	              "ByteMatchStatement": {
	                "SearchString": "/health",
	                "PositionalConstraint": "STARTS_WITH",
	                "FieldToMatch": {"UriPath": {}},
	                "TextTransformations": [{"Priority": 0, "Type": "NONE"}]
	              }
	            }
	          }
	        }
	      }
	    },
	    "VisibilityConfig": {"CloudWatchMetricsEnabled": true, "SampledRequestsEnabled": true, "MetricName": "aws-managed-common"}
	  },
	  {
	    "Name": "rate-limit-per-api-key",
	    "Priority": 2,
	    "Action": {"Block": {}},
	    "Statement": {
	      "RateBasedStatement": {
	        "Limit": 2000,
	        "AggregateKeyType": "CUSTOM_KEYS",
	        "CustomKeys": [
	          {"Header": {"Name": "x-api-key", "TextTransformations": [{"Priority": 0, "Type": "NONE"}]}},
	          {"IP": {}}
	        ]
	      }
	    },
	    "VisibilityConfig": {"CloudWatchMetricsEnabled": true, "SampledRequestsEnabled": true, "MetricName": "rate-limit-per-api-key"}
	  },
	  {
	    "Name": "block-bad-paths",
	    "Priority": 3,
	    "Action": {"Block": {"CustomResponse": {"ResponseCode": 403}}},
	    "Statement": {
	      "OrStatement": {
	        "Statements": [
	          {
	            "ByteMatchStatement": {
	              "SearchString": "/wp-admin",
	              "PositionalConstraint": "STARTS_WITH",
	              "FieldToMatch": {"UriPath": {}},
	              "TextTransformations": [{"Priority": 0, "Type": "LOWERCASE"}]
	            }
	          },
	          {
	            "SqliMatchStatement": {
	              "FieldToMatch": {"Body": {"OversizeHandling": "MATCH"}},
	              "TextTransformations": [{"Priority": 0, "Type": "URL_DECODE"}]
	            }
	          }
	        ]
	      }
	    },
	    "VisibilityConfig": {"CloudWatchMetricsEnabled": true, "SampledRequestsEnabled": true, "MetricName": "block-bad-paths"}
	  }
	]`

	var gotValue, expectedValue interface{}
	if err := json.Unmarshal([]byte(got), &gotValue); err != nil {
		t.Fatalf("unmarshal generated JSON: %v", err)
	}
	if err := json.Unmarshal([]byte(expected), &expectedValue); err != nil {
		t.Fatalf("unmarshal expected JSON: %v", err)
	}
	if !reflect.DeepEqual(gotValue, expectedValue) {
		t.Errorf("rules JSON diverges from the WAF API contract.\ngot:      %s\nexpected: %s", got, expected)
	}
}

// TestBuildRulesJSON_RefArnsResolve asserts reference statements serialize the
// resolved literal value of their StringValueOrRef ARNs.
func TestBuildRulesJSON_RefArnsResolve(t *testing.T) {
	spec := &awswafwebaclv1alpha1.AwsWafWebAclSpec{
		Rules: []*awswafwebaclv1alpha1.AwsWafWebAclRule{
			{
				Name:     "office-allow",
				Priority: 1,
				Action:   "allow",
				Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement{
					Statement: &awswafwebaclv1alpha1.AwsWafWebAclStatement_IpSetReference{
						IpSetReference: &awswafwebaclv1alpha1.AwsWafWebAclIpSetReferenceStatement{
							Arn: &foreignkeyv1.StringValueOrRef{
								LiteralOrRef: &foreignkeyv1.StringValueOrRef_Value{Value: "arn:aws:wafv2:us-west-2:1:regional/ipset/office/x"},
							},
						},
					},
				},
			},
		},
	}

	got, err := buildRulesJSON(spec)
	if err != nil {
		t.Fatalf("buildRulesJSON: %v", err)
	}

	var rules []map[string]interface{}
	if err := json.Unmarshal([]byte(got), &rules); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	arn := rules[0]["Statement"].(map[string]interface{})["IPSetReferenceStatement"].(map[string]interface{})["ARN"]
	if arn != "arn:aws:wafv2:us-west-2:1:regional/ipset/office/x" {
		t.Errorf("IP set ARN not serialized from the resolved reference, got %v", arn)
	}
}
