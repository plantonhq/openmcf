package module

import (
	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/redshiftserverless"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// namespace provisions the serverless namespace itself -- database name,
// admin credentials, data-encryption key, engine IAM roles, and audit
// log exports. Create-only in AWS: the namespace name and the first
// database's name. Everything else edits in place, including credential
// changes and KMS re-encryption (in place but long-running).
func namespace(
	ctx *pulumi.Context,
	locals *Locals,
	provider *aws.Provider,
) (*redshiftserverless.Namespace, error) {
	spec := locals.AwsRedshiftServerlessNamespace.Spec

	args := &redshiftserverless.NamespaceArgs{
		NamespaceName: pulumi.String(locals.NamespaceName),
		Tags:          pulumi.ToStringMap(locals.AwsTags),
	}

	// Empty keeps the AWS default first database ("dev"). Create-time
	// only -- additional databases are created with SQL, not here.
	if spec.DbName != "" {
		args.DbName = pulumi.String(spec.DbName)
	}

	// Empty keeps the AWS default admin ("admin"). A serverless
	// namespace does not hard-require admin credentials at all -- IAM
	// identities can use temporary credentials without one.
	if spec.AdminUsername != "" {
		args.AdminUsername = pulumi.String(spec.AdminUsername)
	}

	// The password contract (CEL enforces exactly one strategy): the
	// AWS-managed Secrets Manager secret (recommended -- no secret in
	// manifest or state) or a directly supplied password.
	// ManageAdminPassword is forwarded ONLY when true: an explicit false
	// conflicts with admin_user_password in the provider's ConflictsWith
	// machinery.
	if spec.ManageAdminPassword {
		args.ManageAdminPassword = pulumi.Bool(true)
		if spec.AdminPasswordSecretKmsKeyId.GetValue() != "" {
			args.AdminPasswordSecretKmsKeyId = pulumi.String(spec.AdminPasswordSecretKmsKeyId.GetValue())
		}
	} else if spec.AdminUserPassword != "" {
		args.AdminUserPassword = pulumi.String(spec.AdminUserPassword)
	}

	// Data encryption at rest. Empty keeps the AWS-owned Redshift
	// service key; switching keys later is an in-place but long-running
	// re-encryption.
	if spec.KmsKeyId.GetValue() != "" {
		args.KmsKeyId = pulumi.String(spec.KmsKeyId.GetValue())
	}

	// IAM roles the serverless engine assumes for COPY/UNLOAD/Spectrum.
	// The default role must also be in iam_roles (an AWS requirement the
	// error message makes obvious enough to leave to the API).
	if len(spec.IamRoles) > 0 {
		iamRoles := pulumi.StringArray{}
		for _, iamRole := range spec.IamRoles {
			iamRoles = append(iamRoles, pulumi.String(iamRole.GetValue()))
		}
		args.IamRoles = iamRoles
	}
	if spec.DefaultIamRoleArn.GetValue() != "" {
		args.DefaultIamRoleArn = pulumi.String(spec.DefaultIamRoleArn.GetValue())
	}

	// Audit log delivery to CloudWatch Logs. Empty exports nothing.
	if len(spec.LogExports) > 0 {
		logExports := pulumi.StringArray{}
		for _, logExport := range spec.LogExports {
			logExports = append(logExports, pulumi.String(logExport))
		}
		args.LogExports = logExports
	}

	createdNamespace, err := redshiftserverless.NewNamespace(ctx, "namespace", args, pulumi.Provider(provider))
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Redshift Serverless namespace")
	}
	return createdNamespace, nil
}
