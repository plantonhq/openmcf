package module

import (
	"sort"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/bedrock"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// guardrail creates the Bedrock guardrail's mutable draft definition plus
// one published version per spec.versions entry, and exports outputs.
func guardrail(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	args := &bedrock.GuardrailArgs{
		// Create-time naming basis; doubles as the Name tag. metadata.name
		// on both engines.
		Name: pulumi.String(locals.GuardrailName),
		// Required by AWS for every guardrail: what the caller sees when
		// the guardrail intervenes on input/output.
		BlockedInputMessaging:   pulumi.String(spec.BlockedInputMessaging),
		BlockedOutputsMessaging: pulumi.String(spec.BlockedOutputsMessaging),
		Tags:                    pulumi.ToStringMap(locals.AwsTags),
	}

	// Description is sent only when set: the provider attribute is
	// Optional+Computed, so sending "" would fight AWS's normalization.
	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}

	// Customer-managed key when referenced; Bedrock-managed key otherwise.
	if spec.KmsKeyArn.GetValue() != "" {
		args.KmsKeyArn = pulumi.String(spec.KmsKeyArn.GetValue())
	}

	// -------------------------------------------------------------------
	// Content filters
	// -------------------------------------------------------------------
	if spec.ContentPolicy != nil {
		policy := &bedrock.GuardrailContentPolicyConfigArgs{}
		// The safeguard tier rides an optional single-entry list. Omitted
		// entirely when unset -- the provider treats the absent list as
		// "AWS default" (CLASSIC) and pins whatever AWS materializes.
		if spec.ContentPolicy.Tier != "" {
			policy.TierConfigs = bedrock.GuardrailContentPolicyConfigTierConfigArray{
				&bedrock.GuardrailContentPolicyConfigTierConfigArgs{
					TierName: pulumi.String(spec.ContentPolicy.Tier),
				},
			}
		}
		var filters bedrock.GuardrailContentPolicyConfigFiltersConfigArray
		for _, f := range spec.ContentPolicy.Filters {
			filter := &bedrock.GuardrailContentPolicyConfigFiltersConfigArgs{
				Type:           pulumi.String(f.Type),
				InputStrength:  pulumi.String(f.InputStrength),
				OutputStrength: pulumi.String(f.OutputStrength),
			}
			// Action/enabled arms are send-when-set: AWS defaults actions
			// to BLOCK and enabled to true; explicit values (including
			// false) are always transmitted so disablement is expressible.
			if f.InputAction != "" {
				filter.InputAction = pulumi.String(f.InputAction)
			}
			if f.OutputAction != "" {
				filter.OutputAction = pulumi.String(f.OutputAction)
			}
			if f.InputEnabled != nil {
				filter.InputEnabled = pulumi.Bool(*f.InputEnabled)
			}
			if f.OutputEnabled != nil {
				filter.OutputEnabled = pulumi.Bool(*f.OutputEnabled)
			}
			if len(f.InputModalities) > 0 {
				filter.InputModalities = pulumi.ToStringArray(f.InputModalities)
			}
			if len(f.OutputModalities) > 0 {
				filter.OutputModalities = pulumi.ToStringArray(f.OutputModalities)
			}
			filters = append(filters, filter)
		}
		policy.FiltersConfigs = filters
		args.ContentPolicyConfig = policy
	}

	// -------------------------------------------------------------------
	// Denied topics
	// -------------------------------------------------------------------
	if spec.TopicPolicy != nil {
		policy := &bedrock.GuardrailTopicPolicyConfigArgs{}
		if spec.TopicPolicy.Tier != "" {
			policy.TierConfigs = bedrock.GuardrailTopicPolicyConfigTierConfigArray{
				&bedrock.GuardrailTopicPolicyConfigTierConfigArgs{
					TierName: pulumi.String(spec.TopicPolicy.Tier),
				},
			}
		}
		var topics bedrock.GuardrailTopicPolicyConfigTopicsConfigArray
		for _, t := range spec.TopicPolicy.Topics {
			topic := &bedrock.GuardrailTopicPolicyConfigTopicsConfigArgs{
				Name:       pulumi.String(t.Name),
				Definition: pulumi.String(t.Definition),
				// DENY is the only topic type AWS defines -- the modules
				// own the constant so the spec never asks for a one-value
				// field.
				Type: pulumi.String("DENY"),
			}
			if len(t.Examples) > 0 {
				topic.Examples = pulumi.ToStringArray(t.Examples)
			}
			topics = append(topics, topic)
		}
		policy.TopicsConfigs = topics
		args.TopicPolicyConfig = policy
	}

	// -------------------------------------------------------------------
	// Word filters
	// -------------------------------------------------------------------
	if spec.WordPolicy != nil {
		policy := &bedrock.GuardrailWordPolicyConfigArgs{}
		// The AWS-managed profanity list -- PROFANITY is the only managed
		// list type AWS defines; presence of spec.profanity_filter enables
		// it and the modules own the type constant.
		if pf := spec.WordPolicy.ProfanityFilter; pf != nil {
			managed := &bedrock.GuardrailWordPolicyConfigManagedWordListsConfigArgs{
				Type: pulumi.String("PROFANITY"),
			}
			if pf.InputAction != "" {
				managed.InputAction = pulumi.String(pf.InputAction)
			}
			if pf.OutputAction != "" {
				managed.OutputAction = pulumi.String(pf.OutputAction)
			}
			if pf.InputEnabled != nil {
				managed.InputEnabled = pulumi.Bool(*pf.InputEnabled)
			}
			if pf.OutputEnabled != nil {
				managed.OutputEnabled = pulumi.Bool(*pf.OutputEnabled)
			}
			policy.ManagedWordListsConfigs = bedrock.GuardrailWordPolicyConfigManagedWordListsConfigArray{managed}
		}
		var words bedrock.GuardrailWordPolicyConfigWordsConfigArray
		for _, w := range spec.WordPolicy.CustomWords {
			word := &bedrock.GuardrailWordPolicyConfigWordsConfigArgs{
				Text: pulumi.String(w.Text),
			}
			if w.InputAction != "" {
				word.InputAction = pulumi.String(w.InputAction)
			}
			if w.OutputAction != "" {
				word.OutputAction = pulumi.String(w.OutputAction)
			}
			if w.InputEnabled != nil {
				word.InputEnabled = pulumi.Bool(*w.InputEnabled)
			}
			if w.OutputEnabled != nil {
				word.OutputEnabled = pulumi.Bool(*w.OutputEnabled)
			}
			words = append(words, word)
		}
		if len(words) > 0 {
			policy.WordsConfigs = words
		}
		args.WordPolicyConfig = policy
	}

	// -------------------------------------------------------------------
	// Sensitive information (PII + regexes)
	// -------------------------------------------------------------------
	if spec.SensitiveInformationPolicy != nil {
		policy := &bedrock.GuardrailSensitiveInformationPolicyConfigArgs{}
		var piiEntities bedrock.GuardrailSensitiveInformationPolicyConfigPiiEntitiesConfigArray
		for _, e := range spec.SensitiveInformationPolicy.PiiEntities {
			entity := &bedrock.GuardrailSensitiveInformationPolicyConfigPiiEntitiesConfigArgs{
				Type:   pulumi.String(e.Type),
				Action: pulumi.String(e.Action),
			}
			// Per-direction overrides are Optional+Computed at the
			// provider: AWS materializes them from `action` when omitted,
			// and once set they never revert to AWS-derived (taught on the
			// spec fields).
			if e.InputAction != "" {
				entity.InputAction = pulumi.String(e.InputAction)
			}
			if e.OutputAction != "" {
				entity.OutputAction = pulumi.String(e.OutputAction)
			}
			if e.InputEnabled != nil {
				entity.InputEnabled = pulumi.Bool(*e.InputEnabled)
			}
			if e.OutputEnabled != nil {
				entity.OutputEnabled = pulumi.Bool(*e.OutputEnabled)
			}
			piiEntities = append(piiEntities, entity)
		}
		if len(piiEntities) > 0 {
			policy.PiiEntitiesConfigs = piiEntities
		}
		var regexes bedrock.GuardrailSensitiveInformationPolicyConfigRegexesConfigArray
		for _, r := range spec.SensitiveInformationPolicy.Regexes {
			regex := &bedrock.GuardrailSensitiveInformationPolicyConfigRegexesConfigArgs{
				Name:    pulumi.String(r.Name),
				Pattern: pulumi.String(r.Pattern),
				Action:  pulumi.String(r.Action),
			}
			if r.Description != "" {
				regex.Description = pulumi.String(r.Description)
			}
			if r.InputAction != "" {
				regex.InputAction = pulumi.String(r.InputAction)
			}
			if r.OutputAction != "" {
				regex.OutputAction = pulumi.String(r.OutputAction)
			}
			if r.InputEnabled != nil {
				regex.InputEnabled = pulumi.Bool(*r.InputEnabled)
			}
			if r.OutputEnabled != nil {
				regex.OutputEnabled = pulumi.Bool(*r.OutputEnabled)
			}
			regexes = append(regexes, regex)
		}
		if len(regexes) > 0 {
			policy.RegexesConfigs = regexes
		}
		args.SensitiveInformationPolicyConfig = policy
	}

	// -------------------------------------------------------------------
	// Contextual grounding
	// -------------------------------------------------------------------
	if spec.ContextualGroundingPolicy != nil {
		var filters bedrock.GuardrailContextualGroundingPolicyConfigFiltersConfigArray
		for _, f := range spec.ContextualGroundingPolicy.Filters {
			filters = append(filters, &bedrock.GuardrailContextualGroundingPolicyConfigFiltersConfigArgs{
				Type:      pulumi.String(f.Type),
				Threshold: pulumi.Float64(f.Threshold),
			})
		}
		args.ContextualGroundingPolicyConfig = &bedrock.GuardrailContextualGroundingPolicyConfigArgs{
			FiltersConfigs: filters,
		}
	}

	// -------------------------------------------------------------------
	// Cross-region inference profile
	// -------------------------------------------------------------------
	if spec.CrossRegionProfileArn != "" {
		args.CrossRegionConfig = &bedrock.GuardrailCrossRegionConfigArgs{
			GuardrailProfileIdentifier: pulumi.String(spec.CrossRegionProfileArn),
		}
	}

	createdGuardrail, err := bedrock.NewGuardrail(ctx, locals.GuardrailName, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create guardrail")
	}

	ctx.Export(OpGuardrailId, createdGuardrail.GuardrailId)
	ctx.Export(OpGuardrailArn, createdGuardrail.GuardrailArn)
	ctx.Export(OpDraftVersion, createdGuardrail.Version)

	// Published versions -- one immutable numbered version per
	// spec.versions entry, keyed by the entry's stable name. Passing the
	// guardrail's ARN output (plus the explicit dependency) orders every
	// publish AFTER the draft update in the same deploy, so a spec edit
	// plus a new entry captures the edited definition. Iteration is
	// name-sorted for deterministic previews.
	versionNumbers := pulumi.StringMap{}
	sortedVersions := make([]string, 0, len(spec.Versions))
	versionsByName := map[string]int{}
	for i, v := range spec.Versions {
		sortedVersions = append(sortedVersions, v.Name)
		versionsByName[v.Name] = i
	}
	sort.Strings(sortedVersions)
	for _, name := range sortedVersions {
		v := spec.Versions[versionsByName[name]]
		versionArgs := &bedrock.GuardrailVersionArgs{
			GuardrailArn: createdGuardrail.GuardrailArn,
			// Keep the published version in AWS when the entry (or the
			// whole guardrail) is removed from management.
			SkipDestroy: pulumi.Bool(v.KeepOnDelete),
		}
		if v.Description != "" {
			versionArgs.Description = pulumi.String(v.Description)
		}
		createdVersion, err := bedrock.NewGuardrailVersion(ctx, "version-"+v.Name, versionArgs,
			pulumi.Provider(provider), pulumi.DependsOn([]pulumi.Resource{createdGuardrail}))
		if err != nil {
			return errors.Wrapf(err, "publish guardrail version %q", v.Name)
		}
		versionNumbers[v.Name] = createdVersion.Version
	}
	ctx.Export(OpVersionNumbers, versionNumbers)

	return nil
}
