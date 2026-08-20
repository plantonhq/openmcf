package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/amp"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// scraper creates the AMP scraper and its logging configuration, and
// exports outputs.
//
// Lifecycle facts the render below depends on:
//   - the whole SOURCE replaces on change (AWS re-provisions the
//     collector); alias, destination, role configuration, and the
//     scrape configuration update in place;
//   - the spec guarantees exactly one source arm and exactly one
//     destination arm; VPC sources carry their own scrape
//     configuration (AWS's published default exists only for EKS);
//   - creates run long (AWS provisions collector infrastructure - the
//     provider waits up to 30 minutes) and deletes drain before
//     removal (up to 20);
//   - the logging configuration is created via update against the
//     scraper ID and replaces with it.
func scraper(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	// Unset scrape configuration resolves to AWS's published default
	// (EKS sources only, spec-guaranteed).
	scrapeConfiguration := spec.ScrapeConfiguration
	if scrapeConfiguration == "" {
		defaultConfiguration, err := amp.GetDefaultScraperConfiguration(ctx, nil, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "resolve default scraper configuration")
		}
		scrapeConfiguration = defaultConfiguration.Configuration
	}

	sourceArgs := &amp.ScraperSourceArgs{}
	if eks := spec.SourceEks; eks != nil {
		subnetIds := pulumi.StringArray{}
		for _, subnetId := range eks.SubnetIds {
			subnetIds = append(subnetIds, pulumi.String(subnetId.GetValue()))
		}
		eksArgs := &amp.ScraperSourceEksArgs{
			ClusterArn: pulumi.String(eks.ClusterArn.GetValue()),
			SubnetIds:  subnetIds,
		}
		if len(eks.SecurityGroupIds) > 0 {
			securityGroupIds := pulumi.StringArray{}
			for _, securityGroupId := range eks.SecurityGroupIds {
				securityGroupIds = append(securityGroupIds, pulumi.String(securityGroupId.GetValue()))
			}
			eksArgs.SecurityGroupIds = securityGroupIds
		}
		sourceArgs.Eks = eksArgs
	}
	if vpc := spec.SourceVpc; vpc != nil {
		subnetIds := pulumi.StringArray{}
		for _, subnetId := range vpc.SubnetIds {
			subnetIds = append(subnetIds, pulumi.String(subnetId.GetValue()))
		}
		securityGroupIds := pulumi.StringArray{}
		for _, securityGroupId := range vpc.SecurityGroupIds {
			securityGroupIds = append(securityGroupIds, pulumi.String(securityGroupId.GetValue()))
		}
		sourceArgs.Vpc = &amp.ScraperSourceVpcArgs{
			SubnetIds:        subnetIds,
			SecurityGroupIds: securityGroupIds,
		}
	}

	destinationArgs := &amp.ScraperDestinationArgs{}
	if spec.AmpWorkspaceArn.GetValue() != "" {
		destinationArgs.Amp = &amp.ScraperDestinationAmpArgs{
			WorkspaceArn: pulumi.String(spec.AmpWorkspaceArn.GetValue()),
		}
	}
	if spec.CloudwatchDatasetArn.GetValue() != "" {
		destinationArgs.Cloudwatch = &amp.ScraperDestinationCloudwatchArgs{
			DatasetArn: pulumi.String(spec.CloudwatchDatasetArn.GetValue()),
		}
	}

	args := &amp.ScraperArgs{
		ScrapeConfiguration: pulumi.String(scrapeConfiguration),
		Source:              sourceArgs,
		Destination:         destinationArgs,
		Tags:                pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.Alias != "" {
		args.Alias = pulumi.String(spec.Alias)
	}
	if roleConfiguration := spec.RoleConfiguration; roleConfiguration != nil {
		roleArgs := &amp.ScraperRoleConfigurationArgs{}
		if roleConfiguration.SourceRoleArn.GetValue() != "" {
			roleArgs.SourceRoleArn = pulumi.String(roleConfiguration.SourceRoleArn.GetValue())
		}
		if roleConfiguration.TargetRoleArn.GetValue() != "" {
			roleArgs.TargetRoleArn = pulumi.String(roleConfiguration.TargetRoleArn.GetValue())
		}
		args.RoleConfiguration = roleArgs
	}

	createdScraper, err := amp.NewScraper(ctx, "scraper", args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create scraper")
	}

	if logging := spec.Logging; logging != nil {
		loggingArgs := &amp.ScraperLoggingConfigurationArgs{
			ScraperId: createdScraper.ID(),
			LoggingDestination: &amp.ScraperLoggingConfigurationLoggingDestinationArgs{
				CloudwatchLogs: &amp.ScraperLoggingConfigurationLoggingDestinationCloudwatchLogsArgs{
					LogGroupArn: pulumi.String(wildcardLogGroupArn(logging.LogGroupArn.GetValue())),
				},
			},
		}
		if len(logging.Components) > 0 {
			loggingArgs.ScraperComponents = pulumi.ToStringArray(logging.Components)
		}
		if _, err := amp.NewScraperLoggingConfiguration(ctx, "scraper_logging_configuration", loggingArgs, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "create scraper logging configuration")
		}
	}

	ctx.Export(OpScraperId, createdScraper.ID())
	ctx.Export(OpScraperArn, createdScraper.Arn)
	ctx.Export(OpScraperRoleArn, createdScraper.RoleArn)
	return nil
}
