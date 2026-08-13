package module

import (
	"strings"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// optionGroup provisions the option group managed for inline
// spec.options -- engine features like Oracle TDE/OEM or SQL Server
// native backup, activated as a named option list (glue, so it stays
// inside this module). Engine name and major version derive from the
// spec. The provider's EC2-Classic db_security_group_memberships
// argument is deliberately unused -- security group access composes
// through vpc_security_group_memberships references. Returns nil when
// the spec brings no inline options.
func optionGroup(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) (*rds.OptionGroup, error) {
	spec := locals.AwsRdsInstance.Spec
	if len(spec.Options) == 0 {
		return nil, nil
	}

	options := rds.OptionGroupOptionArray{}
	for _, option := range spec.Options {
		optionArgs := &rds.OptionGroupOptionArgs{
			OptionName: pulumi.String(option.OptionName),
		}
		if option.Port != 0 {
			optionArgs.Port = pulumi.Int(int(option.Port))
		}
		if option.Version != "" {
			optionArgs.Version = pulumi.String(option.Version)
		}
		if len(option.VpcSecurityGroupMemberships) > 0 {
			memberships := pulumi.StringArray{}
			for _, membership := range option.VpcSecurityGroupMemberships {
				memberships = append(memberships, pulumi.String(membership.GetValue()))
			}
			optionArgs.VpcSecurityGroupMemberships = memberships
		}
		if len(option.OptionSettings) > 0 {
			settings := rds.OptionGroupOptionOptionSettingArray{}
			for _, setting := range option.OptionSettings {
				settings = append(settings, &rds.OptionGroupOptionOptionSettingArgs{
					Name:  pulumi.String(setting.Name),
					Value: pulumi.String(setting.Value),
				})
			}
			optionArgs.OptionSettings = settings
		}
		options = append(options, optionArgs)
	}

	createdOptionGroup, err := rds.NewOptionGroup(ctx, "option-group",
		&rds.OptionGroupArgs{
			Name:                   pulumi.String(locals.InstanceIdentifier),
			EngineName:             pulumi.String(spec.Engine),
			MajorEngineVersion:     pulumi.String(optionMajorEngineVersion(spec.Engine, spec.EngineVersion)),
			OptionGroupDescription: pulumi.String("Managed by Planton for " + locals.InstanceIdentifier),
			Options:                options,
			Tags:                   pulumi.ToStringMap(locals.AwsTags),
		}, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create option group")
	}
	return createdOptionGroup, nil
}

// optionMajorEngineVersion keeps AWS's raw segment convention for the
// option group's major engine version (sqlserver wants "16.00", not
// "16.0"; oracle wants the bare major). Mirrors the Terraform module's
// locals key-for-key.
func optionMajorEngineVersion(engine, engineVersion string) string {
	segments := strings.Split(engineVersion, ".")
	if strings.HasPrefix(engine, "oracle-") {
		return segments[0]
	}
	if len(segments) >= 2 {
		return segments[0] + "." + segments[1]
	}
	return segments[0]
}
