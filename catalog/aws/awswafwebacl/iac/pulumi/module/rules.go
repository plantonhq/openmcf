package module

import (
	"encoding/json"

	"github.com/pkg/errors"
	awswafwebaclv1alpha1 "github.com/plantonhq/planton/catalog/aws/awswafwebacl/v1alpha1"
)

// This file serializes the spec's typed, recursive statement tree into the
// AWS WAFv2 API JSON (PascalCase keys) that the provider's rule_json argument
// expects. Both engines pass rules to the provider as this JSON — the typed
// HCL/SDK rule schemas add nothing here because AWS validates the JSON
// directly — so THIS mapping (and its Terraform twin in iac/tf/locals.tf) is
// the single behavioral surface that must stay in lockstep with the spec.
//
// Defaults applied while serializing (identical in both engines):
//   - rule visibility_config: metrics on, sampling on, metric_name = rule name
//   - rate_based aggregate_key_type: "IP"
//   - ip_set_reference forwarded-IP position: "FIRST"

// buildRulesJSON constructs the AWS WAFv2 API JSON representation of all rules.
func buildRulesJSON(spec *awswafwebaclv1alpha1.AwsWafWebAclSpec) (string, error) {
	var rules []map[string]interface{}

	for _, rule := range spec.Rules {
		ruleMap := map[string]interface{}{
			"Name": rule.Name,
			// GetPriority: the field is presence-typed so priority 0 is
			// expressible; requiredness is CEL-enforced, so it is never nil.
			"Priority": rule.GetPriority(),
		}

		statement, err := buildStatement(rule.Statement)
		if err != nil {
			return "", errors.Wrapf(err, "rule %q", rule.Name)
		}
		ruleMap["Statement"] = statement

		// Exactly one of action / override_action is present (CEL-enforced):
		// match rules carry an action, group rules carry an override action.
		if rule.Action != "" {
			ruleMap["Action"] = buildAction(rule.Action, rule.CustomResponse, rule.CustomRequestHeaders)
		}
		if rule.OverrideAction != "" {
			ruleMap["OverrideAction"] = buildOverrideAction(rule.OverrideAction)
		}

		ruleMap["VisibilityConfig"] = buildRuleVisibilityConfig(rule)

		if len(rule.RuleLabels) > 0 {
			var labels []map[string]interface{}
			for _, label := range rule.RuleLabels {
				labels = append(labels, map[string]interface{}{"Name": label})
			}
			ruleMap["RuleLabels"] = labels
		}

		// Per-rule immunity-time overrides for CAPTCHA / silent challenges.
		if rule.CaptchaConfig != nil {
			ruleMap["CaptchaConfig"] = buildImmunityConfig(rule.CaptchaConfig)
		}
		if rule.ChallengeConfig != nil {
			ruleMap["ChallengeConfig"] = buildImmunityConfig(rule.ChallengeConfig)
		}

		rules = append(rules, ruleMap)
	}

	bytes, err := json.Marshal(rules)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// buildStatement converts one statement-tree node to the AWS API JSON format.
// The tree is recursive through and/or/not and scope-down statements; the
// recursion terminates because the spec value is finite.
//
// PARITY-EXCEPTION: this Go recursion handles arbitrary nesting depth, while
// the Terraform twin (iac/tf/locals.tf) unrolls the tree to THREE levels below
// the root — HCL cannot recurse — and fails the plan loudly beyond that. The
// divergence is depth-only (identical JSON for any tree Terraform accepts) and
// is documented on both sides; see locals.tf for the full rationale.
func buildStatement(statement *awswafwebaclv1alpha1.AwsWafWebAclStatement) (map[string]interface{}, error) {
	if statement == nil {
		return nil, errors.New("statement is required")
	}

	switch stmt := statement.Statement.(type) {
	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_ManagedRuleGroup:
		return buildManagedRuleGroupStatement(stmt.ManagedRuleGroup)

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_RateBased:
		return buildRateBasedStatement(stmt.RateBased)

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_RuleGroupReference:
		return buildRuleGroupReferenceStatement(stmt.RuleGroupReference), nil

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_IpSetReference:
		return buildIpSetReferenceStatement(stmt.IpSetReference), nil

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_RegexPatternSetReference:
		return buildRegexPatternSetReferenceStatement(stmt.RegexPatternSetReference)

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_GeoMatch:
		return buildGeoMatchStatement(stmt.GeoMatch), nil

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_ByteMatch:
		return buildByteMatchStatement(stmt.ByteMatch)

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_SqliMatch:
		return buildSqliMatchStatement(stmt.SqliMatch)

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_XssMatch:
		return buildXssMatchStatement(stmt.XssMatch)

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_SizeConstraint:
		return buildSizeConstraintStatement(stmt.SizeConstraint)

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_RegexMatch:
		return buildRegexMatchStatement(stmt.RegexMatch)

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_LabelMatch:
		return map[string]interface{}{
			"LabelMatchStatement": map[string]interface{}{
				"Scope": stmt.LabelMatch.Scope,
				"Key":   stmt.LabelMatch.Key,
			},
		}, nil

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_AsnMatch:
		asn := map[string]interface{}{
			"AsnList": stmt.AsnMatch.AsnList,
		}
		if stmt.AsnMatch.ForwardedIpConfig != nil {
			asn["ForwardedIPConfig"] = buildForwardedIpConfig(stmt.AsnMatch.ForwardedIpConfig)
		}
		return map[string]interface{}{"AsnMatchStatement": asn}, nil

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_AndStatement:
		children, err := buildStatementList(stmt.AndStatement.Statements)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"AndStatement": map[string]interface{}{"Statements": children},
		}, nil

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_OrStatement:
		children, err := buildStatementList(stmt.OrStatement.Statements)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"OrStatement": map[string]interface{}{"Statements": children},
		}, nil

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_NotStatement:
		child, err := buildStatement(stmt.NotStatement.Statement)
		if err != nil {
			return nil, err
		}
		return map[string]interface{}{
			"NotStatement": map[string]interface{}{"Statement": child},
		}, nil

	case *awswafwebaclv1alpha1.AwsWafWebAclStatement_CustomStatement:
		// Escape hatch: the user writes PascalCase keys matching the AWS
		// WAFv2 API format, so the Struct passes through verbatim.
		return stmt.CustomStatement.AsMap(), nil

	default:
		return nil, errors.New("exactly one statement type must be set")
	}
}

// buildStatementList serializes the children of an and/or statement.
func buildStatementList(statements []*awswafwebaclv1alpha1.AwsWafWebAclStatement) ([]map[string]interface{}, error) {
	children := make([]map[string]interface{}, 0, len(statements))
	for _, statement := range statements {
		child, err := buildStatement(statement)
		if err != nil {
			return nil, err
		}
		children = append(children, child)
	}
	return children, nil
}

// buildManagedRuleGroupStatement constructs the AWS API JSON for a managed
// rule group reference, including the intelligent-threat configs (Bot
// Control, ATP, ACFP, anti-DDoS).
func buildManagedRuleGroupStatement(mg *awswafwebaclv1alpha1.AwsWafWebAclManagedRuleGroupStatement) (map[string]interface{}, error) {
	stmt := map[string]interface{}{
		"Name":       mg.Name,
		"VendorName": mg.VendorName,
	}

	if mg.Version != "" {
		stmt["Version"] = mg.Version
	}

	if len(mg.RuleActionOverrides) > 0 {
		stmt["RuleActionOverrides"] = buildRuleActionOverrides(mg.RuleActionOverrides)
	}

	if mg.ScopeDownStatement != nil {
		scopeDown, err := buildStatement(mg.ScopeDownStatement)
		if err != nil {
			return nil, errors.Wrap(err, "scope_down_statement")
		}
		stmt["ScopeDownStatement"] = scopeDown
	}

	if mg.ManagedRuleGroupConfigs != nil {
		stmt["ManagedRuleGroupConfigs"] = buildManagedRuleGroupConfigs(mg.ManagedRuleGroupConfigs)
	}

	return map[string]interface{}{"ManagedRuleGroupStatement": stmt}, nil
}

// buildManagedRuleGroupConfigs constructs the AWS API config LIST: each
// intelligent-threat rule set contributes one entry.
func buildManagedRuleGroupConfigs(configs *awswafwebaclv1alpha1.AwsWafWebAclManagedRuleGroupConfigs) []map[string]interface{} {
	var out []map[string]interface{}

	if bc := configs.BotControl; bc != nil {
		botControl := map[string]interface{}{
			"InspectionLevel": bc.InspectionLevel,
		}
		if bc.EnableMachineLearning != nil {
			botControl["EnableMachineLearning"] = bc.GetEnableMachineLearning()
		}
		out = append(out, map[string]interface{}{"AWSManagedRulesBotControlRuleSet": botControl})
	}

	if atp := configs.AccountTakeoverPrevention; atp != nil {
		atpConfig := map[string]interface{}{
			"LoginPath": atp.LoginPath,
		}
		if atp.EnableRegexInPath {
			atpConfig["EnableRegexInPath"] = true
		}
		if ri := atp.RequestInspection; ri != nil {
			atpConfig["RequestInspection"] = map[string]interface{}{
				"PayloadType":   ri.PayloadType,
				"UsernameField": map[string]interface{}{"Identifier": ri.UsernameField.Identifier},
				"PasswordField": map[string]interface{}{"Identifier": ri.PasswordField.Identifier},
			}
		}
		if atp.ResponseInspection != nil {
			atpConfig["ResponseInspection"] = buildResponseInspection(atp.ResponseInspection)
		}
		out = append(out, map[string]interface{}{"AWSManagedRulesATPRuleSet": atpConfig})
	}

	if acfp := configs.AccountCreationFraudPrevention; acfp != nil {
		acfpConfig := map[string]interface{}{
			"CreationPath":         acfp.CreationPath,
			"RegistrationPagePath": acfp.RegistrationPagePath,
		}
		if acfp.EnableRegexInPath {
			acfpConfig["EnableRegexInPath"] = true
		}
		if ri := acfp.RequestInspection; ri != nil {
			inspection := map[string]interface{}{
				"PayloadType": ri.PayloadType,
			}
			if ri.UsernameField != nil {
				inspection["UsernameField"] = map[string]interface{}{"Identifier": ri.UsernameField.Identifier}
			}
			if ri.PasswordField != nil {
				inspection["PasswordField"] = map[string]interface{}{"Identifier": ri.PasswordField.Identifier}
			}
			if ri.EmailField != nil {
				inspection["EmailField"] = map[string]interface{}{"Identifier": ri.EmailField.Identifier}
			}
			// AWS models the multi-field groups as lists of {Identifier}.
			if ri.PhoneNumberFields != nil {
				inspection["PhoneNumberFields"] = buildIdentifierList(ri.PhoneNumberFields.Identifiers)
			}
			if ri.AddressFields != nil {
				inspection["AddressFields"] = buildIdentifierList(ri.AddressFields.Identifiers)
			}
			acfpConfig["RequestInspection"] = inspection
		}
		if acfp.ResponseInspection != nil {
			acfpConfig["ResponseInspection"] = buildResponseInspection(acfp.ResponseInspection)
		}
		out = append(out, map[string]interface{}{"AWSManagedRulesACFPRuleSet": acfpConfig})
	}

	if dd := configs.AntiDdos; dd != nil {
		challenge := map[string]interface{}{
			"UsageOfAction": dd.ClientSideAction.UsageOfAction,
		}
		if dd.ClientSideAction.Sensitivity != "" {
			challenge["Sensitivity"] = dd.ClientSideAction.Sensitivity
		}
		if len(dd.ClientSideAction.ExemptUriRegularExpressions) > 0 {
			var regexes []map[string]interface{}
			for _, expression := range dd.ClientSideAction.ExemptUriRegularExpressions {
				regexes = append(regexes, map[string]interface{}{"RegexString": expression})
			}
			challenge["ExemptUriRegularExpressions"] = regexes
		}
		antiDdos := map[string]interface{}{
			"ClientSideActionConfig": map[string]interface{}{"Challenge": challenge},
		}
		if dd.SensitivityToBlock != "" {
			antiDdos["SensitivityToBlock"] = dd.SensitivityToBlock
		}
		out = append(out, map[string]interface{}{"AWSManagedRulesAntiDDoSRuleSet": antiDdos})
	}

	return out
}

// buildIdentifierList wraps plain identifiers into AWS's [{Identifier}] shape.
func buildIdentifierList(identifiers []string) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(identifiers))
	for _, identifier := range identifiers {
		out = append(out, map[string]interface{}{"Identifier": identifier})
	}
	return out
}

// buildResponseInspection constructs the ATP/ACFP response-inspection JSON —
// how the rule set recognizes success vs failure in the application's own
// responses.
func buildResponseInspection(ri *awswafwebaclv1alpha1.AwsWafWebAclResponseInspection) map[string]interface{} {
	inspection := map[string]interface{}{}
	if sc := ri.StatusCode; sc != nil {
		inspection["StatusCode"] = map[string]interface{}{
			"SuccessCodes": sc.SuccessCodes,
			"FailureCodes": sc.FailureCodes,
		}
	}
	if h := ri.Header; h != nil {
		inspection["Header"] = map[string]interface{}{
			"Name":          h.Name,
			"SuccessValues": h.SuccessValues,
			"FailureValues": h.FailureValues,
		}
	}
	if j := ri.BodyJson; j != nil {
		inspection["Json"] = map[string]interface{}{
			"Identifier":    j.Identifier,
			"SuccessValues": j.SuccessValues,
			"FailureValues": j.FailureValues,
		}
	}
	if bc := ri.BodyContains; bc != nil {
		inspection["BodyContains"] = map[string]interface{}{
			"SuccessStrings": bc.SuccessStrings,
			"FailureStrings": bc.FailureStrings,
		}
	}
	return inspection
}

// buildRuleActionOverrides serializes per-rule action overrides for managed
// and referenced rule groups.
func buildRuleActionOverrides(overrides []*awswafwebaclv1alpha1.AwsWafWebAclRuleActionOverride) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(overrides))
	for _, override := range overrides {
		out = append(out, map[string]interface{}{
			"Name":        override.Name,
			"ActionToUse": buildSimpleAction(override.Action),
		})
	}
	return out
}

// buildRuleGroupReferenceStatement constructs the AWS API JSON for a
// customer-owned rule group reference.
func buildRuleGroupReferenceStatement(rg *awswafwebaclv1alpha1.AwsWafWebAclRuleGroupReferenceStatement) map[string]interface{} {
	stmt := map[string]interface{}{
		"ARN": rg.Arn,
	}
	if len(rg.RuleActionOverrides) > 0 {
		stmt["RuleActionOverrides"] = buildRuleActionOverrides(rg.RuleActionOverrides)
	}
	return map[string]interface{}{"RuleGroupReferenceStatement": stmt}
}

// buildRateBasedStatement constructs the AWS API JSON for a rate-based rule.
func buildRateBasedStatement(rb *awswafwebaclv1alpha1.AwsWafWebAclRateBasedStatement) (map[string]interface{}, error) {
	stmt := map[string]interface{}{
		"Limit": rb.Limit,
	}

	// Aggregate key type defaults to IP (both engines apply the same default).
	aggregateKeyType := "IP"
	if rb.AggregateKeyType != "" {
		aggregateKeyType = rb.AggregateKeyType
	}
	stmt["AggregateKeyType"] = aggregateKeyType

	if rb.EvaluationWindowSec > 0 {
		stmt["EvaluationWindowSec"] = rb.EvaluationWindowSec
	}

	if len(rb.CustomKeys) > 0 {
		stmt["CustomKeys"] = buildRateBasedCustomKeys(rb.CustomKeys)
	}

	if rb.ForwardedIpConfig != nil {
		stmt["ForwardedIPConfig"] = buildForwardedIpConfig(rb.ForwardedIpConfig)
	}

	if rb.ScopeDownStatement != nil {
		scopeDown, err := buildStatement(rb.ScopeDownStatement)
		if err != nil {
			return nil, errors.Wrap(err, "scope_down_statement")
		}
		stmt["ScopeDownStatement"] = scopeDown
	}

	return map[string]interface{}{"RateBasedStatement": stmt}, nil
}

// buildRateBasedCustomKeys serializes the composite aggregation key entries.
func buildRateBasedCustomKeys(keys []*awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(keys))
	for _, key := range keys {
		switch k := key.Key.(type) {
		case *awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_Header:
			out = append(out, map[string]interface{}{"Header": map[string]interface{}{
				"Name":                k.Header.Name,
				"TextTransformations": buildTextTransformations(k.Header.TextTransformations),
			}})
		case *awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_Cookie:
			out = append(out, map[string]interface{}{"Cookie": map[string]interface{}{
				"Name":                k.Cookie.Name,
				"TextTransformations": buildTextTransformations(k.Cookie.TextTransformations),
			}})
		case *awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_QueryArgument:
			out = append(out, map[string]interface{}{"QueryArgument": map[string]interface{}{
				"Name":                k.QueryArgument.Name,
				"TextTransformations": buildTextTransformations(k.QueryArgument.TextTransformations),
			}})
		case *awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_QueryString:
			out = append(out, map[string]interface{}{"QueryString": map[string]interface{}{
				"TextTransformations": buildTextTransformations(k.QueryString.TextTransformations),
			}})
		case *awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_UriPath:
			out = append(out, map[string]interface{}{"UriPath": map[string]interface{}{
				"TextTransformations": buildTextTransformations(k.UriPath.TextTransformations),
			}})
		case *awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_HttpMethod:
			if k.HttpMethod {
				out = append(out, map[string]interface{}{"HTTPMethod": map[string]interface{}{}})
			}
		case *awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_Ip:
			if k.Ip {
				out = append(out, map[string]interface{}{"IP": map[string]interface{}{}})
			}
		case *awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_ForwardedIp:
			if k.ForwardedIp {
				out = append(out, map[string]interface{}{"ForwardedIP": map[string]interface{}{}})
			}
		case *awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_Asn:
			if k.Asn {
				out = append(out, map[string]interface{}{"ASN": map[string]interface{}{}})
			}
		case *awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_LabelNamespace:
			out = append(out, map[string]interface{}{"LabelNamespace": map[string]interface{}{
				"Namespace": k.LabelNamespace.Namespace,
			}})
		case *awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_Ja3Fingerprint:
			out = append(out, map[string]interface{}{"JA3Fingerprint": map[string]interface{}{
				"FallbackBehavior": k.Ja3Fingerprint.FallbackBehavior,
			}})
		case *awswafwebaclv1alpha1.AwsWafWebAclRateBasedCustomKey_Ja4Fingerprint:
			out = append(out, map[string]interface{}{"JA4Fingerprint": map[string]interface{}{
				"FallbackBehavior": k.Ja4Fingerprint.FallbackBehavior,
			}})
		}
	}
	return out
}

// buildGeoMatchStatement constructs the AWS API JSON for a geo match rule.
func buildGeoMatchStatement(gm *awswafwebaclv1alpha1.AwsWafWebAclGeoMatchStatement) map[string]interface{} {
	stmt := map[string]interface{}{
		"CountryCodes": gm.CountryCodes,
	}

	if gm.ForwardedIpConfig != nil {
		stmt["ForwardedIPConfig"] = buildForwardedIpConfig(gm.ForwardedIpConfig)
	}

	return map[string]interface{}{"GeoMatchStatement": stmt}
}

// buildIpSetReferenceStatement constructs the AWS API JSON for an IP set
// reference rule. The set ARN arrives resolved (the orchestrator resolves
// references before IaC runs).
func buildIpSetReferenceStatement(ip *awswafwebaclv1alpha1.AwsWafWebAclIpSetReferenceStatement) map[string]interface{} {
	stmt := map[string]interface{}{
		"ARN": ip.Arn.GetValue(),
	}

	if ip.ForwardedIpConfig != nil {
		fwdConfig := map[string]interface{}{
			"HeaderName":       ip.ForwardedIpConfig.HeaderName,
			"FallbackBehavior": ip.ForwardedIpConfig.FallbackBehavior,
		}
		// Position only exists on the IP-set variant of forwarded-IP config;
		// both engines default it to FIRST.
		position := "FIRST"
		if ip.ForwardedIpConfig.Position != "" {
			position = ip.ForwardedIpConfig.Position
		}
		fwdConfig["Position"] = position
		stmt["IPSetForwardedIPConfig"] = fwdConfig
	}

	return map[string]interface{}{"IPSetReferenceStatement": stmt}
}

// buildRegexPatternSetReferenceStatement constructs the AWS API JSON for a
// regex pattern set reference rule.
func buildRegexPatternSetReferenceStatement(rp *awswafwebaclv1alpha1.AwsWafWebAclRegexPatternSetReferenceStatement) (map[string]interface{}, error) {
	fieldToMatch, err := buildFieldToMatch(rp.FieldToMatch)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"RegexPatternSetReferenceStatement": map[string]interface{}{
			"ARN":                 rp.Arn.GetValue(),
			"FieldToMatch":        fieldToMatch,
			"TextTransformations": buildTextTransformations(rp.TextTransformations),
		},
	}, nil
}

// buildByteMatchStatement constructs the AWS API JSON for a byte match rule.
func buildByteMatchStatement(bm *awswafwebaclv1alpha1.AwsWafWebAclByteMatchStatement) (map[string]interface{}, error) {
	fieldToMatch, err := buildFieldToMatch(bm.FieldToMatch)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"ByteMatchStatement": map[string]interface{}{
			"SearchString":         bm.SearchString,
			"PositionalConstraint": bm.PositionalConstraint,
			"FieldToMatch":         fieldToMatch,
			"TextTransformations":  buildTextTransformations(bm.TextTransformations),
		},
	}, nil
}

// buildSqliMatchStatement constructs the AWS API JSON for a SQLi detection rule.
func buildSqliMatchStatement(sm *awswafwebaclv1alpha1.AwsWafWebAclSqliMatchStatement) (map[string]interface{}, error) {
	fieldToMatch, err := buildFieldToMatch(sm.FieldToMatch)
	if err != nil {
		return nil, err
	}
	stmt := map[string]interface{}{
		"FieldToMatch":        fieldToMatch,
		"TextTransformations": buildTextTransformations(sm.TextTransformations),
	}
	if sm.SensitivityLevel != "" {
		stmt["SensitivityLevel"] = sm.SensitivityLevel
	}
	return map[string]interface{}{"SqliMatchStatement": stmt}, nil
}

// buildXssMatchStatement constructs the AWS API JSON for an XSS detection rule.
func buildXssMatchStatement(xm *awswafwebaclv1alpha1.AwsWafWebAclXssMatchStatement) (map[string]interface{}, error) {
	fieldToMatch, err := buildFieldToMatch(xm.FieldToMatch)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"XssMatchStatement": map[string]interface{}{
			"FieldToMatch":        fieldToMatch,
			"TextTransformations": buildTextTransformations(xm.TextTransformations),
		},
	}, nil
}

// buildSizeConstraintStatement constructs the AWS API JSON for a size
// comparison rule.
func buildSizeConstraintStatement(sc *awswafwebaclv1alpha1.AwsWafWebAclSizeConstraintStatement) (map[string]interface{}, error) {
	fieldToMatch, err := buildFieldToMatch(sc.FieldToMatch)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"SizeConstraintStatement": map[string]interface{}{
			"ComparisonOperator":  sc.ComparisonOperator,
			"Size":                sc.Size,
			"FieldToMatch":        fieldToMatch,
			"TextTransformations": buildTextTransformations(sc.TextTransformations),
		},
	}, nil
}

// buildRegexMatchStatement constructs the AWS API JSON for an inline regex rule.
func buildRegexMatchStatement(rm *awswafwebaclv1alpha1.AwsWafWebAclRegexMatchStatement) (map[string]interface{}, error) {
	fieldToMatch, err := buildFieldToMatch(rm.FieldToMatch)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{
		"RegexMatchStatement": map[string]interface{}{
			"RegexString":         rm.RegexString,
			"FieldToMatch":        fieldToMatch,
			"TextTransformations": buildTextTransformations(rm.TextTransformations),
		},
	}, nil
}

// buildFieldToMatch converts the field selector oneof into the AWS API shape:
// empty-config components are empty objects; configured components carry
// their settings.
func buildFieldToMatch(ftm *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch) (map[string]interface{}, error) {
	if ftm == nil {
		return nil, errors.New("field_to_match is required")
	}

	switch f := ftm.Field.(type) {
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_UriPath:
		return map[string]interface{}{"UriPath": map[string]interface{}{}}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_QueryString:
		return map[string]interface{}{"QueryString": map[string]interface{}{}}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_Method:
		return map[string]interface{}{"Method": map[string]interface{}{}}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_AllQueryArguments:
		return map[string]interface{}{"AllQueryArguments": map[string]interface{}{}}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_SingleHeader:
		return map[string]interface{}{"SingleHeader": map[string]interface{}{"Name": f.SingleHeader.Name}}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_SingleQueryArgument:
		return map[string]interface{}{"SingleQueryArgument": map[string]interface{}{"Name": f.SingleQueryArgument.Name}}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_Body:
		body := map[string]interface{}{}
		if f.Body.OversizeHandling != "" {
			body["OversizeHandling"] = f.Body.OversizeHandling
		}
		return map[string]interface{}{"Body": body}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_JsonBody:
		jsonBody := map[string]interface{}{
			"MatchScope":   f.JsonBody.MatchScope,
			"MatchPattern": buildJsonMatchPattern(f.JsonBody.IncludedPaths),
		}
		if f.JsonBody.InvalidFallbackBehavior != "" {
			jsonBody["InvalidFallbackBehavior"] = f.JsonBody.InvalidFallbackBehavior
		}
		if f.JsonBody.OversizeHandling != "" {
			jsonBody["OversizeHandling"] = f.JsonBody.OversizeHandling
		}
		return map[string]interface{}{"JsonBody": jsonBody}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_Headers:
		return map[string]interface{}{"Headers": map[string]interface{}{
			"MatchPattern":     buildNamePattern(f.Headers.MatchPattern, "IncludedHeaders", "ExcludedHeaders"),
			"MatchScope":       f.Headers.MatchScope,
			"OversizeHandling": f.Headers.OversizeHandling,
		}}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_Cookies:
		return map[string]interface{}{"Cookies": map[string]interface{}{
			"MatchPattern":     buildNamePattern(f.Cookies.MatchPattern, "IncludedCookies", "ExcludedCookies"),
			"MatchScope":       f.Cookies.MatchScope,
			"OversizeHandling": f.Cookies.OversizeHandling,
		}}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_HeaderOrder:
		return map[string]interface{}{"HeaderOrder": map[string]interface{}{
			"OversizeHandling": f.HeaderOrder.OversizeHandling,
		}}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_Ja3Fingerprint:
		return map[string]interface{}{"JA3Fingerprint": map[string]interface{}{
			"FallbackBehavior": f.Ja3Fingerprint.FallbackBehavior,
		}}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_Ja4Fingerprint:
		return map[string]interface{}{"JA4Fingerprint": map[string]interface{}{
			"FallbackBehavior": f.Ja4Fingerprint.FallbackBehavior,
		}}, nil
	case *awswafwebaclv1alpha1.AwsWafWebAclFieldToMatch_UriFragment:
		fragment := map[string]interface{}{}
		if f.UriFragment.FallbackBehavior != "" {
			fragment["FallbackBehavior"] = f.UriFragment.FallbackBehavior
		}
		return map[string]interface{}{"UriFragment": fragment}, nil
	default:
		return nil, errors.New("exactly one field_to_match component must be set")
	}
}

// buildJsonMatchPattern maps included paths to the AWS shape: empty means
// inspect ALL elements.
func buildJsonMatchPattern(includedPaths []string) map[string]interface{} {
	if len(includedPaths) > 0 {
		return map[string]interface{}{"IncludedPaths": includedPaths}
	}
	return map[string]interface{}{"All": map[string]interface{}{}}
}

// buildNamePattern maps the headers/cookies name selector to the AWS shape.
// The include/exclude key names differ between the two components
// (IncludedHeaders vs IncludedCookies), so the caller passes them in.
func buildNamePattern(pattern *awswafwebaclv1alpha1.AwsWafWebAclNamePattern, includedKey, excludedKey string) map[string]interface{} {
	switch {
	case len(pattern.IncludedNames) > 0:
		return map[string]interface{}{includedKey: pattern.IncludedNames}
	case len(pattern.ExcludedNames) > 0:
		return map[string]interface{}{excludedKey: pattern.ExcludedNames}
	default:
		return map[string]interface{}{"All": map[string]interface{}{}}
	}
}

// buildTextTransformations serializes the normalization steps applied before
// matching.
func buildTextTransformations(transformations []*awswafwebaclv1alpha1.AwsWafWebAclTextTransformation) []map[string]interface{} {
	out := make([]map[string]interface{}, 0, len(transformations))
	for _, transformation := range transformations {
		out = append(out, map[string]interface{}{
			"Priority": transformation.Priority,
			"Type":     transformation.Type,
		})
	}
	return out
}

// buildForwardedIpConfig constructs the AWS API JSON for forwarded IP config
// used by geo_match, asn_match, and rate_based statements. (The IP-set
// variant with Position is built inline by buildIpSetReferenceStatement.)
func buildForwardedIpConfig(fwd *awswafwebaclv1alpha1.AwsWafWebAclForwardedIpConfig) map[string]interface{} {
	return map[string]interface{}{
		"HeaderName":       fwd.HeaderName,
		"FallbackBehavior": fwd.FallbackBehavior,
	}
}

// buildImmunityConfig constructs a per-rule or web-ACL-level CAPTCHA/challenge
// immunity-time block.
func buildImmunityConfig(config *awswafwebaclv1alpha1.AwsWafWebAclImmunityTimeConfig) map[string]interface{} {
	return map[string]interface{}{
		"ImmunityTimeProperty": map[string]interface{}{
			"ImmunityTime": config.ImmunityTimeSec,
		},
	}
}

// buildAction constructs the AWS API JSON for a rule action (allow, block,
// count, captcha, challenge) with optional custom response/headers.
func buildAction(
	actionType string,
	customResponse *awswafwebaclv1alpha1.AwsWafWebAclCustomResponse,
	customHeaders []*awswafwebaclv1alpha1.AwsWafWebAclCustomHeader,
) map[string]interface{} {
	actionContent := map[string]interface{}{}

	switch actionType {
	case "block":
		blockContent := map[string]interface{}{}
		if customResponse != nil {
			cr := map[string]interface{}{
				"ResponseCode": customResponse.ResponseCode,
			}
			if customResponse.CustomResponseBodyKey != "" {
				cr["CustomResponseBodyKey"] = customResponse.CustomResponseBodyKey
			}
			if len(customResponse.ResponseHeaders) > 0 {
				var headers []map[string]interface{}
				for _, h := range customResponse.ResponseHeaders {
					headers = append(headers, map[string]interface{}{
						"Name":  h.Name,
						"Value": h.Value,
					})
				}
				cr["ResponseHeaders"] = headers
			}
			blockContent["CustomResponse"] = cr
		}
		actionContent["Block"] = blockContent

	case "allow":
		allowContent := map[string]interface{}{}
		if len(customHeaders) > 0 {
			allowContent["CustomRequestHandling"] = buildCustomRequestHandling(customHeaders)
		}
		actionContent["Allow"] = allowContent

	case "count":
		countContent := map[string]interface{}{}
		if len(customHeaders) > 0 {
			countContent["CustomRequestHandling"] = buildCustomRequestHandling(customHeaders)
		}
		actionContent["Count"] = countContent

	case "captcha":
		captchaContent := map[string]interface{}{}
		if len(customHeaders) > 0 {
			captchaContent["CustomRequestHandling"] = buildCustomRequestHandling(customHeaders)
		}
		actionContent["Captcha"] = captchaContent

	case "challenge":
		challengeContent := map[string]interface{}{}
		if len(customHeaders) > 0 {
			challengeContent["CustomRequestHandling"] = buildCustomRequestHandling(customHeaders)
		}
		actionContent["Challenge"] = challengeContent
	}

	return actionContent
}

// buildSimpleAction constructs a simple action JSON object (no custom response/headers).
// Used for rule action overrides within managed rule groups.
func buildSimpleAction(actionType string) map[string]interface{} {
	switch actionType {
	case "block":
		return map[string]interface{}{"Block": map[string]interface{}{}}
	case "allow":
		return map[string]interface{}{"Allow": map[string]interface{}{}}
	case "count":
		return map[string]interface{}{"Count": map[string]interface{}{}}
	case "captcha":
		return map[string]interface{}{"Captcha": map[string]interface{}{}}
	case "challenge":
		return map[string]interface{}{"Challenge": map[string]interface{}{}}
	default:
		return map[string]interface{}{"Count": map[string]interface{}{}}
	}
}

// buildOverrideAction constructs the AWS API JSON for an override action
// (used with managed and referenced rule group rules).
func buildOverrideAction(overrideType string) map[string]interface{} {
	if overrideType == "count" {
		return map[string]interface{}{"Count": map[string]interface{}{}}
	}
	// "none" means use the rule group's own actions.
	return map[string]interface{}{"None": map[string]interface{}{}}
}

// buildCustomRequestHandling constructs the CustomRequestHandling JSON object.
func buildCustomRequestHandling(headers []*awswafwebaclv1alpha1.AwsWafWebAclCustomHeader) map[string]interface{} {
	var insertHeaders []map[string]interface{}
	for _, h := range headers {
		insertHeaders = append(insertHeaders, map[string]interface{}{
			"Name":  h.Name,
			"Value": h.Value,
		})
	}
	return map[string]interface{}{"InsertHeaders": insertHeaders}
}

// buildRuleVisibilityConfig constructs the visibility config for a single rule,
// applying smart defaults when the user omits it (metrics on, sampling on,
// metric name = rule name — identical in both engines).
func buildRuleVisibilityConfig(rule *awswafwebaclv1alpha1.AwsWafWebAclRule) map[string]interface{} {
	metricsEnabled := true
	sampledEnabled := true
	metricName := rule.Name

	if rule.VisibilityConfig != nil {
		metricsEnabled = rule.VisibilityConfig.CloudwatchMetricsEnabled
		sampledEnabled = rule.VisibilityConfig.SampledRequestsEnabled
		if rule.VisibilityConfig.MetricName != "" {
			metricName = rule.VisibilityConfig.MetricName
		}
	}

	return map[string]interface{}{
		"CloudWatchMetricsEnabled": metricsEnabled,
		"SampledRequestsEnabled":   sampledEnabled,
		"MetricName":               metricName,
	}
}
