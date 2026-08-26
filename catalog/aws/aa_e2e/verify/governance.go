package verify

import (
	"context"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	cloudtrailtypes "github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	"github.com/aws/aws-sdk-go-v2/service/configservice"
	configtypes "github.com/aws/aws-sdk-go-v2/service/configservice/types"
	"github.com/aws/aws-sdk-go-v2/service/guardduty"
	guarddutytypes "github.com/aws/aws-sdk-go-v2/service/guardduty/types"
	"github.com/pkg/errors"
)

// The governance-family verifiers. Each asserts what its scenario
// WROTE (selector styles, recording scope, feature statuses, active
// lists), never mere presence -- and each kind's destroy contract is a
// REAL delete, with two singleton twists taught on the profiles:
//
//   - AwsConfigRecorder destroys in ORDER (stop recorder, then channel
//     + recorder + retention) -- verify-absent asserts all three
//     regional singletons are gone;
//   - AwsGuardDuty's detector delete cascades every satellite --
//     verify-absent asserts the detector id no longer resolves.

// cloudTrailVerifier verifies AwsCloudTrail via GetTrail +
// GetTrailStatus + GetInsightSelectors, keyed on trail_arn.
type cloudTrailVerifier struct{}

func (*cloudTrailVerifier) IDOutputKey() string { return "trail_arn" }

func (*cloudTrailVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	client := cloudtrail.NewFromConfig(cfg, func(o *cloudtrail.Options) {
		if region != "" {
			o.Region = region
		}
	})
	trail, err := client.GetTrail(ctx, &cloudtrail.GetTrailInput{Name: &id})
	if err != nil {
		return errors.Wrap(err, "GetTrail")
	}
	if trail.Trail == nil {
		return errors.Errorf("trail (%s) not found after deploy", id)
	}
	// The scenario's posture: multi-region, validated, CloudWatch
	// mirroring wired.
	if trail.Trail.IsMultiRegionTrail == nil || !*trail.Trail.IsMultiRegionTrail {
		return errors.Errorf("trail (%s) is not multi-region after deploy", id)
	}
	if trail.Trail.LogFileValidationEnabled == nil || !*trail.Trail.LogFileValidationEnabled {
		return errors.Errorf("trail (%s) has log-file validation off after deploy", id)
	}
	if trail.Trail.CloudWatchLogsLogGroupArn == nil || *trail.Trail.CloudWatchLogsLogGroupArn == "" {
		return errors.Errorf("trail (%s) has no CloudWatch Logs group after deploy", id)
	}
	// enable_logging defaults on -- the trail must be DELIVERING.
	status, err := client.GetTrailStatus(ctx, &cloudtrail.GetTrailStatusInput{Name: &id})
	if err != nil {
		return errors.Wrap(err, "GetTrailStatus")
	}
	if status.IsLogging == nil || !*status.IsLogging {
		return errors.Errorf("trail (%s) is not logging after deploy", id)
	}
	// Both Insights engines were declared.
	insights, err := client.GetInsightSelectors(ctx, &cloudtrail.GetInsightSelectorsInput{TrailName: &id})
	if err != nil {
		return errors.Wrap(err, "GetInsightSelectors")
	}
	if len(insights.InsightSelectors) != 2 {
		return errors.Errorf("trail (%s) reports %d insight selectors after deploy (expected 2)", id, len(insights.InsightSelectors))
	}
	return nil
}

func (*cloudTrailVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := cloudtrail.NewFromConfig(cfg, func(o *cloudtrail.Options) {
		if region != "" {
			o.Region = region
		}
	}).GetTrail(ctx, &cloudtrail.GetTrailInput{Name: &id})
	if err == nil {
		return errors.Errorf("trail (%s) still exists after destroy", id)
	}
	if strings.Contains(err.Error(), "TrailNotFound") {
		return nil
	}
	return errors.Wrap(err, "GetTrail")
}

// configRecorderVerifier verifies AwsConfigRecorder via the four
// Describe calls covering the recorder, its status, the delivery
// channel, and the retention window, keyed on recorder_name.
type configRecorderVerifier struct{}

func (*configRecorderVerifier) IDOutputKey() string { return "recorder_name" }

func (*configRecorderVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	client := configservice.NewFromConfig(cfg, func(o *configservice.Options) {
		if region != "" {
			o.Region = region
		}
	})
	recorders, err := client.DescribeConfigurationRecorders(ctx, &configservice.DescribeConfigurationRecordersInput{
		ConfigurationRecorderNames: []string{id},
	})
	if err != nil {
		return errors.Wrap(err, "DescribeConfigurationRecorders")
	}
	if len(recorders.ConfigurationRecorders) != 1 {
		return errors.Errorf("configuration recorder (%s) not found after deploy", id)
	}
	// The scenario's scoped-inclusion posture -- never all-supported.
	group := recorders.ConfigurationRecorders[0].RecordingGroup
	if group == nil || group.AllSupported || len(group.ResourceTypes) == 0 {
		return errors.Errorf("configuration recorder (%s) is not the scoped inclusion posture after deploy", id)
	}
	// recording_enabled defaults on -- the recorder must be RUNNING.
	status, err := client.DescribeConfigurationRecorderStatus(ctx, &configservice.DescribeConfigurationRecorderStatusInput{
		ConfigurationRecorderNames: []string{id},
	})
	if err != nil {
		return errors.Wrap(err, "DescribeConfigurationRecorderStatus")
	}
	if len(status.ConfigurationRecordersStatus) != 1 || !status.ConfigurationRecordersStatus[0].Recording {
		return errors.Errorf("configuration recorder (%s) is not recording after deploy", id)
	}
	channels, err := client.DescribeDeliveryChannels(ctx, &configservice.DescribeDeliveryChannelsInput{})
	if err != nil {
		return errors.Wrap(err, "DescribeDeliveryChannels")
	}
	if len(channels.DeliveryChannels) == 0 {
		return errors.Errorf("delivery channel for recorder (%s) not found after deploy", id)
	}
	// Retention is deliberately NOT asserted here: it is an optional arm
	// (the modules deploy aws_config_retention_configuration only when
	// spec.retention_period_in_days is set) and the kind's outputs carry
	// no retention key this verifier could condition on -- an
	// unconditional presence assert fails every legitimate
	// retention-less manifest, which is exactly how fixture deployments
	// of this kind compose it (live-caught). The retention posture is
	// proven by the scenario that declares it: its import round-trip
	// re-imports the retention singleton at a no-op plan, and
	// VerifyAbsent below asserts it gone after destroy.
	return nil
}

func (*configRecorderVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	client := configservice.NewFromConfig(cfg, func(o *configservice.Options) {
		if region != "" {
			o.Region = region
		}
	})
	// All three regional singletons must be gone (the stop-then-delete
	// teardown ordering the modules encode).
	recorders, err := client.DescribeConfigurationRecorders(ctx, &configservice.DescribeConfigurationRecordersInput{})
	if err != nil {
		return errors.Wrap(err, "DescribeConfigurationRecorders")
	}
	if len(recorders.ConfigurationRecorders) != 0 {
		return errors.Errorf("configuration recorder (%s) still exists after destroy", id)
	}
	channels, err := client.DescribeDeliveryChannels(ctx, &configservice.DescribeDeliveryChannelsInput{})
	if err != nil {
		return errors.Wrap(err, "DescribeDeliveryChannels")
	}
	if len(channels.DeliveryChannels) != 0 {
		return errors.Errorf("delivery channel (%s) still exists after destroy", id)
	}
	retention, err := client.DescribeRetentionConfigurations(ctx, &configservice.DescribeRetentionConfigurationsInput{})
	if err != nil {
		return errors.Wrap(err, "DescribeRetentionConfigurations")
	}
	if len(retention.RetentionConfigurations) != 0 {
		return errors.Errorf("retention configuration (%s) still exists after destroy", id)
	}
	return nil
}

// configRuleVerifier verifies AwsConfigRule via DescribeConfigRules,
// keyed on rule_name (the account-scoped scenarios; organization
// rules are the recorded single-account deferral).
type configRuleVerifier struct{}

func (*configRuleVerifier) IDOutputKey() string { return "rule_name" }

func (*configRuleVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := configservice.NewFromConfig(cfg, func(o *configservice.Options) {
		if region != "" {
			o.Region = region
		}
	}).DescribeConfigRules(ctx, &configservice.DescribeConfigRulesInput{
		ConfigRuleNames: []string{id},
	})
	if err != nil {
		return errors.Wrap(err, "DescribeConfigRules")
	}
	if len(out.ConfigRules) != 1 {
		return errors.Errorf("config rule (%s) not found after deploy", id)
	}
	rule := out.ConfigRules[0]
	// Assert the source the scenario wrote (managed identifier or the
	// Guard runtime) and that the rule is active.
	if rule.Source == nil {
		return errors.Errorf("config rule (%s) reports no source after deploy", id)
	}
	if rule.Source.Owner == configtypes.OwnerAws && (rule.Source.SourceIdentifier == nil || *rule.Source.SourceIdentifier == "") {
		return errors.Errorf("managed config rule (%s) reports no source identifier after deploy", id)
	}
	if rule.Source.Owner == configtypes.OwnerCustomPolicy && rule.Source.CustomPolicyDetails == nil {
		return errors.Errorf("custom-policy config rule (%s) reports no policy details after deploy", id)
	}
	if rule.ConfigRuleState == configtypes.ConfigRuleStateDeleting || rule.ConfigRuleState == configtypes.ConfigRuleStateDeletingResults {
		return errors.Errorf("config rule (%s) is in state %s after deploy", id, rule.ConfigRuleState)
	}
	return nil
}

func (*configRuleVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	out, err := configservice.NewFromConfig(cfg, func(o *configservice.Options) {
		if region != "" {
			o.Region = region
		}
	}).DescribeConfigRules(ctx, &configservice.DescribeConfigRulesInput{
		ConfigRuleNames: []string{id},
	})
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchConfigRule") {
			return nil
		}
		return errors.Wrap(err, "DescribeConfigRules")
	}
	// Deletion is asynchronous; a rule still draining in DELETING is
	// the soft-deleted state, not an orphan.
	for _, rule := range out.ConfigRules {
		if rule.ConfigRuleState != configtypes.ConfigRuleStateDeleting && rule.ConfigRuleState != configtypes.ConfigRuleStateDeletingResults {
			return errors.Errorf("config rule (%s) still exists after destroy (state %s)", id, rule.ConfigRuleState)
		}
	}
	return nil
}

// guardDutyVerifier verifies AwsGuardDuty via GetDetector plus the
// satellite List calls, keyed on detector_id. The full-surface
// scenario declares every single-account arm, so each is asserted
// firmly (a scenario-keyed variant carries the day if a minimal
// scenario ever lands).
type guardDutyVerifier struct{}

func (*guardDutyVerifier) IDOutputKey() string { return "detector_id" }

func (*guardDutyVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	client := guardduty.NewFromConfig(cfg, func(o *guardduty.Options) {
		if region != "" {
			o.Region = region
		}
	})
	detector, err := client.GetDetector(ctx, &guardduty.GetDetectorInput{DetectorId: &id})
	if err != nil {
		return errors.Wrap(err, "GetDetector")
	}
	if detector.Status != "ENABLED" {
		return errors.Errorf("detector (%s) is %s after deploy (expected ENABLED)", id, detector.Status)
	}
	// The declared protection plans echo on the detector.
	if len(detector.Features) == 0 {
		return errors.Errorf("detector (%s) reports no features after deploy", id)
	}
	filters, err := client.ListFilters(ctx, &guardduty.ListFiltersInput{DetectorId: &id})
	if err != nil {
		return errors.Wrap(err, "ListFilters")
	}
	if len(filters.FilterNames) == 0 {
		return errors.Errorf("detector (%s) has no finding filters after deploy", id)
	}
	ipSets, err := client.ListIPSets(ctx, &guardduty.ListIPSetsInput{DetectorId: &id})
	if err != nil {
		return errors.Wrap(err, "ListIPSets")
	}
	if len(ipSets.IpSetIds) == 0 {
		return errors.Errorf("detector (%s) has no trusted IP sets after deploy", id)
	}
	tiSets, err := client.ListThreatIntelSets(ctx, &guardduty.ListThreatIntelSetsInput{DetectorId: &id})
	if err != nil {
		return errors.Wrap(err, "ListThreatIntelSets")
	}
	if len(tiSets.ThreatIntelSetIds) == 0 {
		return errors.Errorf("detector (%s) has no threat intel sets after deploy", id)
	}
	destinations, err := client.ListPublishingDestinations(ctx, &guardduty.ListPublishingDestinationsInput{DetectorId: &id})
	if err != nil {
		return errors.Wrap(err, "ListPublishingDestinations")
	}
	if len(destinations.Destinations) == 0 {
		return errors.Errorf("detector (%s) has no publishing destination after deploy", id)
	}
	return nil
}

func (*guardDutyVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := guardduty.NewFromConfig(cfg, func(o *guardduty.Options) {
		if region != "" {
			o.Region = region
		}
	}).GetDetector(ctx, &guardduty.GetDetectorInput{DetectorId: &id})
	if err == nil {
		return errors.Errorf("detector (%s) still exists after destroy", id)
	}
	// A deleted detector id stops resolving: the API answers
	// BadRequest ("The request is rejected because the input detectorId
	// is not owned by the current account") rather than NotFound.
	if strings.Contains(err.Error(), "BadRequest") || strings.Contains(err.Error(), "not owned") || strings.Contains(err.Error(), "NotFound") {
		return nil
	}
	return errors.Wrap(err, "GetDetector")
}

// eventDataStoreVerifier verifies AwsCloudTrailEventDataStore via
// GetEventDataStore, keyed on event_data_store_arn. The scenario's
// posture: single-region, 7-day retention, termination protection OFF
// (the single-step-destroy posture), one narrow selector, ingestion
// suspended.
type eventDataStoreVerifier struct{}

func (*eventDataStoreVerifier) IDOutputKey() string { return "event_data_store_arn" }

func (*eventDataStoreVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	store, err := cloudtrail.NewFromConfig(cfg, func(o *cloudtrail.Options) {
		if region != "" {
			o.Region = region
		}
	}).GetEventDataStore(ctx, &cloudtrail.GetEventDataStoreInput{EventDataStore: &id})
	if err != nil {
		return errors.Wrap(err, "GetEventDataStore")
	}
	// "suspend" lands as STOPPED_INGESTION once AWS processes it;
	// ENABLED is the moments-earlier echo of the same healthy store.
	if store.Status != cloudtrailtypes.EventDataStoreStatusStoppedIngestion && store.Status != cloudtrailtypes.EventDataStoreStatusEnabled {
		return errors.Errorf("event data store (%s) is %s after deploy", id, store.Status)
	}
	if store.RetentionPeriod == nil || *store.RetentionPeriod != 7 {
		return errors.Errorf("event data store (%s) does not carry the declared 7-day retention after deploy", id)
	}
	if store.TerminationProtectionEnabled == nil || *store.TerminationProtectionEnabled {
		return errors.Errorf("event data store (%s) has termination protection ON after deploy (the scenario declares it off)", id)
	}
	if store.MultiRegionEnabled == nil || *store.MultiRegionEnabled {
		return errors.Errorf("event data store (%s) is multi-region after deploy (the scenario declares single-region)", id)
	}
	if len(store.AdvancedEventSelectors) != 1 {
		return errors.Errorf("event data store (%s) reports %d selectors after deploy (expected the 1 declared)", id, len(store.AdvancedEventSelectors))
	}
	return nil
}

func (*eventDataStoreVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	store, err := cloudtrail.NewFromConfig(cfg, func(o *cloudtrail.Options) {
		if region != "" {
			o.Region = region
		}
	}).GetEventDataStore(ctx, &cloudtrail.GetEventDataStoreInput{EventDataStore: &id})
	if err != nil {
		if strings.Contains(err.Error(), "EventDataStoreNotFound") {
			return nil
		}
		return errors.Wrap(err, "GetEventDataStore")
	}
	// A deleted store stays describable in PENDING_DELETION for 7 days
	// (its name reserved) before the purge -- the soft-deleted state,
	// not an orphan.
	if store.Status == cloudtrailtypes.EventDataStoreStatusPendingDeletion {
		return nil
	}
	return errors.Errorf("event data store (%s) still exists after destroy (state %s)", id, store.Status)
}

// configAggregatorVerifier verifies AwsConfigAggregator via
// DescribeConfigurationAggregators + DescribeAggregationAuthorizations,
// keyed on aggregator_name. The scenario declares BOTH arms, so both
// are asserted.
type configAggregatorVerifier struct{}

func (*configAggregatorVerifier) IDOutputKey() string { return "aggregator_name" }

func (*configAggregatorVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	client := configservice.NewFromConfig(cfg, func(o *configservice.Options) {
		if region != "" {
			o.Region = region
		}
	})
	aggregators, err := client.DescribeConfigurationAggregators(ctx, &configservice.DescribeConfigurationAggregatorsInput{
		ConfigurationAggregatorNames: []string{id},
	})
	if err != nil {
		return errors.Wrap(err, "DescribeConfigurationAggregators")
	}
	if len(aggregators.ConfigurationAggregators) != 1 {
		return errors.Errorf("configuration aggregator (%s) not found after deploy", id)
	}
	// The scenario's account-source posture (the org source is the
	// recorded single-account deferral).
	if len(aggregators.ConfigurationAggregators[0].AccountAggregationSources) == 0 {
		return errors.Errorf("configuration aggregator (%s) reports no account aggregation source after deploy", id)
	}
	grants, err := client.DescribeAggregationAuthorizations(ctx, &configservice.DescribeAggregationAuthorizationsInput{})
	if err != nil {
		return errors.Wrap(err, "DescribeAggregationAuthorizations")
	}
	if len(grants.AggregationAuthorizations) == 0 {
		return errors.Errorf("aggregation authorization for (%s) not found after deploy", id)
	}
	return nil
}

func (*configAggregatorVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	client := configservice.NewFromConfig(cfg, func(o *configservice.Options) {
		if region != "" {
			o.Region = region
		}
	})
	aggregators, err := client.DescribeConfigurationAggregators(ctx, &configservice.DescribeConfigurationAggregatorsInput{
		ConfigurationAggregatorNames: []string{id},
	})
	if err == nil && len(aggregators.ConfigurationAggregators) != 0 {
		return errors.Errorf("configuration aggregator (%s) still exists after destroy", id)
	}
	if err != nil && !strings.Contains(err.Error(), "NoSuchConfigurationAggregator") {
		return errors.Wrap(err, "DescribeConfigurationAggregators")
	}
	grants, err := client.DescribeAggregationAuthorizations(ctx, &configservice.DescribeAggregationAuthorizationsInput{})
	if err != nil {
		return errors.Wrap(err, "DescribeAggregationAuthorizations")
	}
	if len(grants.AggregationAuthorizations) != 0 {
		return errors.Errorf("aggregation authorization (%s) still exists after destroy", id)
	}
	return nil
}

// conformancePackVerifier verifies AwsConfigConformancePack via
// DescribeConformancePacks + DescribeConformancePackStatus, keyed on
// pack_name (the account-scoped scenario; organization packs are the
// recorded single-account deferral).
type conformancePackVerifier struct{}

func (*conformancePackVerifier) IDOutputKey() string { return "pack_name" }

func (*conformancePackVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	client := configservice.NewFromConfig(cfg, func(o *configservice.Options) {
		if region != "" {
			o.Region = region
		}
	})
	packs, err := client.DescribeConformancePacks(ctx, &configservice.DescribeConformancePacksInput{
		ConformancePackNames: []string{id},
	})
	if err != nil {
		return errors.Wrap(err, "DescribeConformancePacks")
	}
	if len(packs.ConformancePackDetails) != 1 {
		return errors.Errorf("conformance pack (%s) not found after deploy", id)
	}
	// CREATE_COMPLETE is the provider's own create waiter target -- it
	// means the pack's rules materialized.
	status, err := client.DescribeConformancePackStatus(ctx, &configservice.DescribeConformancePackStatusInput{
		ConformancePackNames: []string{id},
	})
	if err != nil {
		return errors.Wrap(err, "DescribeConformancePackStatus")
	}
	if len(status.ConformancePackStatusDetails) != 1 || status.ConformancePackStatusDetails[0].ConformancePackState != configtypes.ConformancePackStateCreateComplete {
		return errors.Errorf("conformance pack (%s) is not CREATE_COMPLETE after deploy", id)
	}
	return nil
}

func (*conformancePackVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	packs, err := configservice.NewFromConfig(cfg, func(o *configservice.Options) {
		if region != "" {
			o.Region = region
		}
	}).DescribeConformancePacks(ctx, &configservice.DescribeConformancePacksInput{
		ConformancePackNames: []string{id},
	})
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchConformancePack") {
			return nil
		}
		return errors.Wrap(err, "DescribeConformancePacks")
	}
	if len(packs.ConformancePackDetails) != 0 {
		return errors.Errorf("conformance pack (%s) still exists after destroy", id)
	}
	return nil
}

// malwareProtectionPlanVerifier verifies AwsGuardDutyMalwareProtectionPlan
// via GetMalwareProtectionPlan, keyed on malware_protection_plan_id.
type malwareProtectionPlanVerifier struct{}

func (*malwareProtectionPlanVerifier) IDOutputKey() string { return "malware_protection_plan_id" }

func (*malwareProtectionPlanVerifier) VerifyExists(ctx context.Context, cfg aws.Config, id, region string) error {
	plan, err := guardduty.NewFromConfig(cfg, func(o *guardduty.Options) {
		if region != "" {
			o.Region = region
		}
	}).GetMalwareProtectionPlan(ctx, &guardduty.GetMalwareProtectionPlanInput{MalwareProtectionPlanId: &id})
	if err != nil {
		return errors.Wrap(err, "GetMalwareProtectionPlan")
	}
	// ERROR is the failed posture; ACTIVE is steady state and WARNING
	// the transient EventBridge-wiring echo AWS documents.
	if plan.Status == guarddutytypes.MalwareProtectionPlanStatusError {
		return errors.Errorf("malware protection plan (%s) is ERROR after deploy", id)
	}
	if plan.ProtectedResource == nil || plan.ProtectedResource.S3Bucket == nil || plan.ProtectedResource.S3Bucket.BucketName == nil || *plan.ProtectedResource.S3Bucket.BucketName == "" {
		return errors.Errorf("malware protection plan (%s) reports no protected bucket after deploy", id)
	}
	// The scenario declares prefix scoping and verdict tagging.
	if len(plan.ProtectedResource.S3Bucket.ObjectPrefixes) == 0 {
		return errors.Errorf("malware protection plan (%s) reports no object prefixes after deploy (the scenario declares one)", id)
	}
	if plan.Actions == nil || plan.Actions.Tagging == nil || plan.Actions.Tagging.Status != guarddutytypes.MalwareProtectionPlanTaggingActionStatusEnabled {
		return errors.Errorf("malware protection plan (%s) does not report verdict tagging ENABLED after deploy", id)
	}
	return nil
}

func (*malwareProtectionPlanVerifier) VerifyAbsent(ctx context.Context, cfg aws.Config, id, region string) error {
	_, err := guardduty.NewFromConfig(cfg, func(o *guardduty.Options) {
		if region != "" {
			o.Region = region
		}
	}).GetMalwareProtectionPlan(ctx, &guardduty.GetMalwareProtectionPlanInput{MalwareProtectionPlanId: &id})
	if err == nil {
		return errors.Errorf("malware protection plan (%s) still exists after destroy", id)
	}
	if strings.Contains(err.Error(), "ResourceNotFound") || strings.Contains(err.Error(), "NotFound") {
		return nil
	}
	return errors.Wrap(err, "GetMalwareProtectionPlan")
}
