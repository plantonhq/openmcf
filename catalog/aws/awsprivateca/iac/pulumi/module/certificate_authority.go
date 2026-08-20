package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/acmpca"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// certificateAuthority creates the CA, composes its activation,
// issues certificates, grants the ACM renewal permission, attaches
// the resource policy, and exports outputs.
//
// Lifecycle facts the render below depends on:
//   - a fresh CA sits in PENDING_CERTIFICATE until a certificate is
//     INSTALLED on it; the CA resource's own certificate attribute is
//     read at create (still empty), so the ca_certificate output
//     comes from the ACTIVATION path, never the CA attribute;
//   - a ROOT self-signs: its own CSR issued against ITSELF with the
//     RootCACertificate template, then installed (three provider
//     calls the raw provider makes users wire by hand);
//   - a SUBORDINATE's certificate is issued by the PARENT CA with a
//     path-length template; the issue call's signing algorithm is
//     this spec's signing_algorithm, so same-family hierarchies are
//     assumed (the norm - a cross-family hierarchy needs out-of-band
//     activation);
//   - issued certificates require the CA ACTIVE - they depend on the
//     activation resource; delete REVOKES them (the provider's
//     documented delete semantic);
//   - deleting the CA parks it restorable for
//     permanent_deletion_time_in_days; billing stops at delete.
func certificateAuthority(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.Spec

	// Template ARNs are partition-scoped
	// (arn:{partition}:acm-pca:::template/...).
	partition, err := aws.GetPartition(ctx, nil, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "resolve partition")
	}

	subject := spec.Subject
	subjectArgs := &acmpca.CertificateAuthorityCertificateAuthorityConfigurationSubjectArgs{}
	if subject.CommonName != "" {
		subjectArgs.CommonName = pulumi.String(subject.CommonName)
	}
	if subject.Organization != "" {
		subjectArgs.Organization = pulumi.String(subject.Organization)
	}
	if subject.OrganizationalUnit != "" {
		subjectArgs.OrganizationalUnit = pulumi.String(subject.OrganizationalUnit)
	}
	if subject.Country != "" {
		subjectArgs.Country = pulumi.String(subject.Country)
	}
	if subject.State != "" {
		subjectArgs.State = pulumi.String(subject.State)
	}
	if subject.Locality != "" {
		subjectArgs.Locality = pulumi.String(subject.Locality)
	}

	caArgs := &acmpca.CertificateAuthorityArgs{
		Type: pulumi.String(spec.Type),
		CertificateAuthorityConfiguration: &acmpca.CertificateAuthorityCertificateAuthorityConfigurationArgs{
			KeyAlgorithm:     pulumi.String(spec.KeyAlgorithm),
			SigningAlgorithm: pulumi.String(spec.SigningAlgorithm),
			Subject:          subjectArgs,
		},
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}
	if spec.UsageMode != "" {
		caArgs.UsageMode = pulumi.String(spec.UsageMode)
	}
	if spec.KeyStorageSecurityStandard != "" {
		caArgs.KeyStorageSecurityStandard = pulumi.String(spec.KeyStorageSecurityStandard)
	}
	if spec.PermanentDeletionTimeInDays > 0 {
		caArgs.PermanentDeletionTimeInDays = pulumi.Int(int(spec.PermanentDeletionTimeInDays))
	}
	if spec.Enabled != nil && !*spec.Enabled {
		caArgs.Enabled = pulumi.Bool(false)
	}
	if revocation := spec.Revocation; revocation != nil {
		revocationArgs := &acmpca.CertificateAuthorityRevocationConfigurationArgs{}
		if crl := revocation.Crl; crl != nil && crl.Enabled {
			crlArgs := &acmpca.CertificateAuthorityRevocationConfigurationCrlConfigurationArgs{
				Enabled:          pulumi.Bool(true),
				ExpirationInDays: pulumi.Int(int(crl.ExpirationInDays)),
				S3BucketName:     pulumi.String(crl.S3BucketName.GetValue()),
			}
			if crl.S3ObjectAcl != "" {
				crlArgs.S3ObjectAcl = pulumi.String(crl.S3ObjectAcl)
			}
			if crl.CustomCname != "" {
				crlArgs.CustomCname = pulumi.String(crl.CustomCname)
			}
			if crl.CustomPath != "" {
				crlArgs.CustomPath = pulumi.String(crl.CustomPath)
			}
			revocationArgs.CrlConfiguration = crlArgs
		}
		if ocsp := revocation.Ocsp; ocsp != nil && ocsp.Enabled {
			ocspArgs := &acmpca.CertificateAuthorityRevocationConfigurationOcspConfigurationArgs{
				Enabled: pulumi.Bool(true),
			}
			if ocsp.CustomCname != "" {
				ocspArgs.OcspCustomCname = pulumi.String(ocsp.CustomCname)
			}
			revocationArgs.OcspConfiguration = ocspArgs
		}
		caArgs.RevocationConfiguration = revocationArgs
	}

	createdCa, err := acmpca.NewCertificateAuthority(ctx, "certificate-authority", caArgs, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "create certificate authority")
	}

	// Activation: the composed CSR->issue->install dance.
	var activationDependency pulumi.Resource
	caCertificate := pulumi.String("").ToStringOutput()
	caCertificateChain := pulumi.String("").ToStringOutput()
	activationCertificateArn := pulumi.String("").ToStringOutput()

	switch {
	case spec.Type == "ROOT":
		validityType, validityValue := "YEARS", "10"
		if spec.RootCaValidity != nil {
			validityType, validityValue = spec.RootCaValidity.Type, spec.RootCaValidity.Value
		}
		createdRootCertificate, err := acmpca.NewCertificate(ctx, "root-ca-certificate", &acmpca.CertificateArgs{
			CertificateAuthorityArn:   createdCa.Arn,
			CertificateSigningRequest: createdCa.CertificateSigningRequest,
			SigningAlgorithm:          pulumi.String(spec.SigningAlgorithm),
			TemplateArn:               pulumi.Sprintf("arn:%s:acm-pca:::template/RootCACertificate/V1", partition.Partition),
			Validity: &acmpca.CertificateValidityArgs{
				Type:  pulumi.String(validityType),
				Value: pulumi.String(validityValue),
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "self-sign root certificate")
		}
		createdActivation, err := acmpca.NewCertificateAuthorityCertificate(ctx, "ca-activation",
			&acmpca.CertificateAuthorityCertificateArgs{
				CertificateAuthorityArn: createdCa.Arn,
				Certificate:             createdRootCertificate.Certificate,
			}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "install root certificate")
		}
		activationDependency = createdActivation
		caCertificate = createdRootCertificate.Certificate
		activationCertificateArn = createdRootCertificate.Arn
		// A root IS the trust anchor - AWS reports no chain for it.

	case spec.SubordinateActivation != nil:
		activation := spec.SubordinateActivation
		createdSubordinateCertificate, err := acmpca.NewCertificate(ctx, "subordinate-ca-certificate", &acmpca.CertificateArgs{
			CertificateAuthorityArn:   pulumi.String(activation.ParentCaArn.GetValue()),
			CertificateSigningRequest: createdCa.CertificateSigningRequest,
			SigningAlgorithm:          pulumi.String(spec.SigningAlgorithm),
			TemplateArn: pulumi.Sprintf("arn:%s:acm-pca:::template/SubordinateCACertificate_PathLen%d/V1",
				partition.Partition, activation.PathLength),
			Validity: &acmpca.CertificateValidityArgs{
				Type:  pulumi.String(activation.Validity.Type),
				Value: pulumi.String(activation.Validity.Value),
			},
		}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "issue subordinate certificate from parent")
		}
		createdActivation, err := acmpca.NewCertificateAuthorityCertificate(ctx, "ca-activation",
			&acmpca.CertificateAuthorityCertificateArgs{
				CertificateAuthorityArn: createdCa.Arn,
				Certificate:             createdSubordinateCertificate.Certificate,
				CertificateChain:        createdSubordinateCertificate.CertificateChain,
			}, pulumi.Provider(provider))
		if err != nil {
			return errors.Wrap(err, "install subordinate certificate")
		}
		activationDependency = createdActivation
		caCertificate = createdSubordinateCertificate.Certificate
		activationCertificateArn = createdSubordinateCertificate.Arn
		caCertificateChain = createdSubordinateCertificate.CertificateChain.ApplyT(func(chain *string) string {
			if chain == nil {
				return ""
			}
			return *chain
		}).(pulumi.StringOutput)

	default:
		// An unactivated subordinate: created, billed, waiting for its
		// certificate out of band (the ca_csr output is what the
		// external parent signs).
	}

	issuedCertificateArns := pulumi.StringMap{}
	for _, certificate := range spec.IssuedCertificates {
		certificateArgs := &acmpca.CertificateArgs{
			CertificateAuthorityArn:   createdCa.Arn,
			CertificateSigningRequest: pulumi.String(certificate.Csr),
			SigningAlgorithm:          pulumi.String(certificate.SigningAlgorithm),
			Validity: &acmpca.CertificateValidityArgs{
				Type:  pulumi.String(certificate.Validity.Type),
				Value: pulumi.String(certificate.Validity.Value),
			},
		}
		if certificate.TemplateArn != "" {
			certificateArgs.TemplateArn = pulumi.String(certificate.TemplateArn)
		}
		if certificate.ApiPassthrough != "" {
			certificateArgs.ApiPassthrough = pulumi.String(certificate.ApiPassthrough)
		}
		options := []pulumi.ResourceOption{pulumi.Provider(provider)}
		if activationDependency != nil {
			// Issuing needs the CA ACTIVE, which activation makes it.
			options = append(options, pulumi.DependsOn([]pulumi.Resource{activationDependency}))
		}
		createdCertificate, err := acmpca.NewCertificate(ctx,
			fmt.Sprintf("certificate-%s", certificate.Name),
			certificateArgs, options...)
		if err != nil {
			return errors.Wrapf(err, "issue certificate %s", certificate.Name)
		}
		issuedCertificateArns[certificate.Name] = createdCertificate.Arn
	}

	if spec.AcmRenewalPermission {
		if _, err := acmpca.NewPermission(ctx, "acm-renewal-permission", &acmpca.PermissionArgs{
			CertificateAuthorityArn: createdCa.Arn,
			Principal:               pulumi.String("acm.amazonaws.com"),
			// All three actions - AWS documents that ACM auto-renewal
			// requires the full set; a partial grant fails silently
			// at renewal time.
			Actions: pulumi.ToStringArray([]string{"IssueCertificate", "GetCertificate", "ListPermissions"}),
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "grant acm renewal permission")
		}
	}

	if spec.Policy != "" {
		if _, err := acmpca.NewPolicy(ctx, "resource-policy", &acmpca.PolicyArgs{
			ResourceArn: createdCa.Arn,
			Policy:      pulumi.String(spec.Policy),
		}, pulumi.Provider(provider)); err != nil {
			return errors.Wrap(err, "attach resource policy")
		}
	}

	ctx.Export(OpCertificateAuthorityArn, createdCa.Arn)
	ctx.Export(OpCertificateAuthorityId, createdCa.ID())
	ctx.Export(OpCaCertificate, caCertificate)
	ctx.Export(OpCaCertificateChain, caCertificateChain)
	ctx.Export(OpCaCsr, createdCa.CertificateSigningRequest)
	ctx.Export(OpIssuedCertificateArns, issuedCertificateArns)
	ctx.Export(OpActivationCertificateArn, activationCertificateArn)
	return nil
}
