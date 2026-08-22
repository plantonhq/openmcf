package module

import (
	"github.com/pkg/errors"
	cloudflarezerotrustdeviceposturerulev1alpha1 "github.com/plantonhq/planton/catalog/cloudflare/cloudflarezerotrustdeviceposturerule/v1alpha1"
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// postureRule creates the device posture rule: a health check WARP evaluates
// on enrolled devices, which Access and Gateway policies can then require. A
// plain CRUD resource (real create/update/delete; only the account forces
// replacement).
func postureRule(
	ctx *pulumi.Context,
	locals *Locals,
	cloudflareProvider *cloudflare.Provider,
) error {
	spec := locals.CloudflareZeroTrustDevicePostureRule.Spec

	args := &cloudflare.ZeroTrustDevicePostureRuleArgs{
		AccountId: pulumi.String(spec.AccountId),
		Name:      pulumi.StringPtr(spec.Name),
		Type:      pulumi.String(spec.Type),
	}

	if spec.Description != "" {
		args.Description = pulumi.String(spec.Description)
	}
	if spec.Expiration != "" {
		args.Expiration = pulumi.String(spec.Expiration)
	}
	if spec.Schedule != "" {
		args.Schedule = pulumi.String(spec.Schedule)
	}

	if len(spec.Match) > 0 {
		matches := cloudflare.ZeroTrustDevicePostureRuleMatchArray{}
		for _, row := range spec.Match {
			matches = append(matches, cloudflare.ZeroTrustDevicePostureRuleMatchArgs{
				Platform: pulumi.StringPtr(row.Platform),
			})
		}
		args.Matches = matches
	}

	if spec.Input != nil {
		args.Input = buildInput(spec.Input)
	}

	createdRule, err := cloudflare.NewZeroTrustDevicePostureRule(
		ctx,
		"posture_rule",
		args,
		pulumi.Provider(cloudflareProvider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create posture rule")
	}

	ctx.Export(OpRuleId, createdRule.ID())

	return nil
}

// buildInput maps the check parameters; unset fields are never sent, so each
// rule's payload holds exactly the fields its type reads.
func buildInput(input *cloudflarezerotrustdeviceposturerulev1alpha1.CloudflareZeroTrustDevicePostureRuleInput) cloudflare.ZeroTrustDevicePostureRuleInputTypePtrInput {
	inputArgs := cloudflare.ZeroTrustDevicePostureRuleInputTypeArgs{}

	if input.OperatingSystem != "" {
		inputArgs.OperatingSystem = pulumi.String(input.OperatingSystem)
	}
	if input.Path != "" {
		inputArgs.Path = pulumi.String(input.Path)
	}
	if input.Exists != nil {
		inputArgs.Exists = pulumi.BoolPtr(input.GetExists())
	}
	if input.Sha256 != "" {
		inputArgs.Sha256 = pulumi.String(input.Sha256)
	}
	if input.Thumbprint != "" {
		inputArgs.Thumbprint = pulumi.String(input.Thumbprint)
	}
	if input.Id != "" {
		inputArgs.Id = pulumi.String(input.Id)
	}
	if input.Domain != "" {
		inputArgs.Domain = pulumi.String(input.Domain)
	}
	if input.Operator != "" {
		inputArgs.Operator = pulumi.String(input.Operator)
	}
	if input.Version != "" {
		inputArgs.Version = pulumi.String(input.Version)
	}
	if input.OsDistroName != "" {
		inputArgs.OsDistroName = pulumi.String(input.OsDistroName)
	}
	if input.OsDistroRevision != "" {
		inputArgs.OsDistroRevision = pulumi.String(input.OsDistroRevision)
	}
	if input.OsVersionExtra != "" {
		inputArgs.OsVersionExtra = pulumi.String(input.OsVersionExtra)
	}
	if input.Enabled != nil {
		inputArgs.Enabled = pulumi.BoolPtr(input.GetEnabled())
	}
	if len(input.CheckDisks) > 0 {
		inputArgs.CheckDisks = pulumi.ToStringArray(input.CheckDisks)
	}
	if input.RequireAll != nil {
		inputArgs.RequireAll = pulumi.BoolPtr(input.GetRequireAll())
	}
	if input.CertificateId != "" {
		inputArgs.CertificateId = pulumi.String(input.CertificateId)
	}
	if input.Cn != "" {
		inputArgs.Cn = pulumi.String(input.Cn)
	}
	if input.CheckPrivateKey != nil {
		inputArgs.CheckPrivateKey = pulumi.BoolPtr(input.GetCheckPrivateKey())
	}
	if len(input.ExtendedKeyUsage) > 0 {
		inputArgs.ExtendedKeyUsages = pulumi.ToStringArray(input.ExtendedKeyUsage)
	}
	if input.Locations != nil {
		locations := cloudflare.ZeroTrustDevicePostureRuleInputLocationsArgs{}
		if len(input.Locations.Paths) > 0 {
			locations.Paths = pulumi.ToStringArray(input.Locations.Paths)
		}
		if len(input.Locations.TrustStores) > 0 {
			locations.TrustStores = pulumi.ToStringArray(input.Locations.TrustStores)
		}
		inputArgs.Locations = locations
	}
	if len(input.SubjectAlternativeNames) > 0 {
		inputArgs.SubjectAlternativeNames = pulumi.ToStringArray(input.SubjectAlternativeNames)
	}
	if input.UpdateWindowDays != nil {
		inputArgs.UpdateWindowDays = pulumi.Float64Ptr(float64(input.GetUpdateWindowDays()))
	}
	if input.ComplianceStatus != "" {
		inputArgs.ComplianceStatus = pulumi.String(input.ComplianceStatus)
	}
	if input.ConnectionId != "" {
		inputArgs.ConnectionId = pulumi.String(input.ConnectionId)
	}
	if input.LastSeen != "" {
		inputArgs.LastSeen = pulumi.String(input.LastSeen)
	}
	if input.Os != "" {
		inputArgs.Os = pulumi.String(input.Os)
	}
	if input.Overall != "" {
		inputArgs.Overall = pulumi.String(input.Overall)
	}
	if input.SensorConfig != "" {
		inputArgs.SensorConfig = pulumi.String(input.SensorConfig)
	}
	if input.State != "" {
		inputArgs.State = pulumi.String(input.State)
	}
	if input.VersionOperator != "" {
		inputArgs.VersionOperator = pulumi.String(input.VersionOperator)
	}
	if len(input.AuthState) > 0 {
		inputArgs.AuthStates = pulumi.ToStringArray(input.AuthState)
	}
	if input.CountOperator != "" {
		inputArgs.CountOperator = pulumi.String(input.CountOperator)
	}
	if input.IssueCount != "" {
		inputArgs.IssueCount = pulumi.String(input.IssueCount)
	}
	if input.EidLastSeen != "" {
		inputArgs.EidLastSeen = pulumi.String(input.EidLastSeen)
	}
	if input.RiskLevel != "" {
		inputArgs.RiskLevel = pulumi.String(input.RiskLevel)
	}
	if input.ScoreOperator != "" {
		inputArgs.ScoreOperator = pulumi.String(input.ScoreOperator)
	}
	if input.TotalScore != nil {
		inputArgs.TotalScore = pulumi.Float64Ptr(input.GetTotalScore())
	}
	if input.ActiveThreats != nil {
		inputArgs.ActiveThreats = pulumi.Float64Ptr(float64(input.GetActiveThreats()))
	}
	if input.Infected != nil {
		inputArgs.Infected = pulumi.BoolPtr(input.GetInfected())
	}
	if input.IsActive != nil {
		inputArgs.IsActive = pulumi.BoolPtr(input.GetIsActive())
	}
	if input.NetworkStatus != "" {
		inputArgs.NetworkStatus = pulumi.String(input.NetworkStatus)
	}
	if input.OperationalState != "" {
		inputArgs.OperationalState = pulumi.String(input.OperationalState)
	}
	if input.Score != nil {
		inputArgs.Score = pulumi.Float64Ptr(input.GetScore())
	}

	return inputArgs
}
