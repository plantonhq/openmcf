package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/sagemaker"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// userProfiles creates the folded aws_sagemaker_user_profile satellites — one
// per spec.user_profiles entry, keyed by profile name so adding or removing
// one person never disturbs the others. Returns the profiles' ARNs keyed by
// name for the user_profile_arns output, plus the created resources so the
// spaces (whose ownership references profiles by NAME, with no implicit
// engine edge) can depend on them.
func userProfiles(ctx *pulumi.Context, locals *Locals, createdDomain *sagemaker.Domain, provider *aws.Provider) (pulumi.StringMap, []pulumi.Resource, error) {
	profileArns := pulumi.StringMap{}
	var createdProfiles []pulumi.Resource

	for _, profile := range locals.Spec.UserProfiles {
		args := &sagemaker.UserProfileArgs{
			DomainId:        createdDomain.ID(),
			UserProfileName: pulumi.String(profile.UserProfileName),
			Tags:            pulumi.ToStringMap(locals.AwsTags),
		}

		// SSO linkage (SSO-auth domains only; the pair travels together,
		// CEL-enforced).
		if profile.SingleSignOnUserIdentifier != "" {
			args.SingleSignOnUserIdentifier = pulumi.String(profile.SingleSignOnUserIdentifier)
			args.SingleSignOnUserValue = pulumi.String(profile.SingleSignOnUserValue)
		}

		// A profile's user_settings is the SAME settings tree as the
		// domain's default_user_settings, rendered by the profile-typed
		// twin builders in profile_user_settings.go.
		if profile.UserSettings != nil {
			args.UserSettings = buildProfileUserSettings(profile.UserSettings)
		}

		createdProfile, err := sagemaker.NewUserProfile(ctx,
			fmt.Sprintf("user-profile-%s", profile.UserProfileName),
			args, pulumi.Provider(provider), pulumi.Parent(createdDomain))
		if err != nil {
			return nil, nil, errors.Wrapf(err, "user profile %s", profile.UserProfileName)
		}

		profileArns[profile.UserProfileName] = createdProfile.Arn
		createdProfiles = append(createdProfiles, createdProfile)
	}

	return profileArns, createdProfiles, nil
}
