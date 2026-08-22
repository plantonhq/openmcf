package module

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/acm"
	"github.com/pulumi/pulumi-aws/sdk/v7/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// certManagerCert provisions an ACM certificate in one of the three creation
// modes -- requested (Amazon-issued), imported (bring-your-own material), or
// private (ACM-PCA) -- plus, for DNS-validated requested certificates with a
// managed Route53 zone, the validation records and the issuance waiter.
// Exactly one mode is populated (CEL enforces the exclusivity), so the arms
// below never mix.
func certManagerCert(ctx *pulumi.Context, locals *Locals, provider *aws.Provider) error {
	spec := locals.AwsCertManagerCert.Spec
	meta := locals.AwsCertManagerCert.Metadata

	isImported := spec.Imported != nil
	isPrivate := !isImported && spec.GetCertificateAuthorityArn().GetValue() != ""
	isRequested := !isImported && !isPrivate

	// Requested certificates validate via DNS unless EMAIL is chosen;
	// imported and private certificates never validate publicly.
	validationMethod := ""
	if isRequested {
		validationMethod = spec.ValidationMethod
		if validationMethod == "" {
			validationMethod = "DNS"
		}
	}

	args := &acm.CertificateArgs{
		Tags: pulumi.ToStringMap(locals.AwsTags),
	}

	switch {
	case isImported:
		// Imported arm: the PEM material. The private key is sensitive spec
		// input and never appears in outputs; re-importing new material
		// before expiry updates in place and keeps the ARN stable for
		// consumers. Domains are derived from the certificate body.
		args.CertificateBody = pulumi.String(spec.Imported.CertificateBody)
		args.PrivateKey = pulumi.String(spec.Imported.PrivateKey)
		if spec.Imported.CertificateChain != "" {
			args.CertificateChain = pulumi.String(spec.Imported.CertificateChain)
		}

	default:
		// Requested + private arms: ACM issues for these domains.
		args.DomainName = pulumi.String(spec.PrimaryDomainName)
		if len(spec.AlternateDomainNames) > 0 {
			args.SubjectAlternativeNames = pulumi.ToStringArray(spec.AlternateDomainNames)
		}
		if isPrivate {
			// Private-CA issuance is authorized by the CA itself -- passing a
			// validation method alongside the CA ARN is rejected by AWS.
			args.CertificateAuthorityArn = pulumi.String(spec.GetCertificateAuthorityArn().GetValue())
			// Managed early renewal -- a private-CA mechanism (ACM's
			// RenewCertificate is private-certificate-only; CEL ties the
			// field to the private arm). Durations under 60 days have no
			// effect.
			if spec.EarlyRenewalDuration != "" {
				args.EarlyRenewalDuration = pulumi.String(spec.EarlyRenewalDuration)
			}
		} else {
			args.ValidationMethod = pulumi.String(validationMethod)
		}
		// Create-time immutable key algorithm; empty keeps the ACM default
		// (RSA_2048).
		if spec.KeyAlgorithm != "" {
			args.KeyAlgorithm = pulumi.String(spec.KeyAlgorithm)
		}
		// Per-domain overrides of where the validation request is sent.
		if len(spec.ValidationOptions) > 0 {
			var validationOptions acm.CertificateValidationOptionArray
			for _, vo := range spec.ValidationOptions {
				validationOptions = append(validationOptions, &acm.CertificateValidationOptionArgs{
					DomainName:       pulumi.String(vo.DomainName),
					ValidationDomain: pulumi.String(vo.ValidationDomain),
				})
			}
			args.ValidationOptions = validationOptions
		}
		// Certificate Transparency logging and exportability. Empty-string
		// sentinels keep the ACM defaults (CT ENABLED, export DISABLED).
		if spec.Options != nil {
			optionArgs := &acm.CertificateOptionsArgs{}
			if spec.Options.CertificateTransparencyLoggingPreference != "" {
				optionArgs.CertificateTransparencyLoggingPreference = pulumi.String(spec.Options.CertificateTransparencyLoggingPreference)
			}
			if spec.Options.Export != "" {
				optionArgs.Export = pulumi.String(spec.Options.Export)
			}
			args.Options = optionArgs
		}
	}

	cert, err := acm.NewCertificate(ctx, meta.Name, args, pulumi.Provider(provider))
	if err != nil {
		return errors.Wrap(err, "failed to create ACM certificate")
	}

	// Managed DNS validation: one validation record per certificate domain,
	// in the referenced Route53 zone. The domains are config-known, so each
	// record is a first-class resource visible at preview; only its
	// name/type/value flow from the certificate's computed
	// domain_validation_options. A domain and its wildcard SAN share the
	// SAME validation CNAME -- AllowOverwrite turns the duplicate into an
	// idempotent UPSERT (and adopts a record left behind by a prior partial
	// apply instead of colliding on "InvalidChangeBatch ... already exists").
	// Matches the Terraform module's per-domain for_each key-for-key.
	route53ZoneId := spec.GetRoute53HostedZoneId().GetValue()
	managesValidationRecords := isRequested && validationMethod == "DNS" && route53ZoneId != ""

	if managesValidationRecords {
		domains := append([]string{spec.PrimaryDomainName}, spec.AlternateDomainNames...)

		var recordFqdns pulumi.StringArray
		for _, domain := range domains {
			domain := domain
			lookup := func(field func(acm.CertificateDomainValidationOption) *string) pulumi.StringOutput {
				return cert.DomainValidationOptions.ApplyT(func(dvos []acm.CertificateDomainValidationOption) string {
					for _, dvo := range dvos {
						if dvo.DomainName != nil && *dvo.DomainName == domain {
							if v := field(dvo); v != nil {
								return *v
							}
						}
					}
					return ""
				}).(pulumi.StringOutput)
			}

			record, err := route53.NewRecord(ctx, fmt.Sprintf("%s-validation-%s", meta.Name, domain), &route53.RecordArgs{
				AllowOverwrite: pulumi.Bool(true),
				ZoneId:         pulumi.String(route53ZoneId),
				Name:           lookup(func(dvo acm.CertificateDomainValidationOption) *string { return dvo.ResourceRecordName }),
				Type:           lookup(func(dvo acm.CertificateDomainValidationOption) *string { return dvo.ResourceRecordType }),
				Records: pulumi.StringArray{
					lookup(func(dvo acm.CertificateDomainValidationOption) *string { return dvo.ResourceRecordValue }),
				},
				Ttl: pulumi.Int(60),
			}, pulumi.Provider(provider))
			if err != nil {
				return errors.Wrapf(err, "failed to create validation record for %s", domain)
			}
			recordFqdns = append(recordFqdns, record.Fqdn)
		}

		// The issuance waiter -- a read-only resource that blocks until ACM
		// reports ISSUED (75-minute ceiling; DNS-validated issuance typically
		// lands in minutes once the records propagate). Absent
		// wait_for_validation resolves to the annotated default (true), the
		// same resolution the Terraform contract applies.
		if spec.WaitForValidation == nil || spec.GetWaitForValidation() {
			_, err = acm.NewCertificateValidation(ctx, meta.Name, &acm.CertificateValidationArgs{
				CertificateArn:        cert.Arn,
				ValidationRecordFqdns: recordFqdns,
			}, pulumi.Provider(provider))
			if err != nil {
				return errors.Wrap(err, "failed to create certificate validation waiter")
			}
		}
	}

	// The validation records external DNS needs when the module does not
	// manage them -- shaped key-for-key with the Terraform output.
	domainValidationRecords := cert.DomainValidationOptions.ApplyT(func(dvos []acm.CertificateDomainValidationOption) []map[string]string {
		records := make([]map[string]string, 0, len(dvos))
		for _, dvo := range dvos {
			record := map[string]string{"domain_name": "", "record_name": "", "record_type": "", "record_value": ""}
			if dvo.DomainName != nil {
				record["domain_name"] = *dvo.DomainName
			}
			if dvo.ResourceRecordName != nil {
				record["record_name"] = *dvo.ResourceRecordName
			}
			if dvo.ResourceRecordType != nil {
				record["record_type"] = *dvo.ResourceRecordType
			}
			if dvo.ResourceRecordValue != nil {
				record["record_value"] = *dvo.ResourceRecordValue
			}
			records = append(records, record)
		}
		return records
	})

	ctx.Export(OpCertArn, cert.Arn)
	ctx.Export(OpStatus, cert.Status)
	ctx.Export(OpDomainValidationRecords, domainValidationRecords)
	ctx.Export(OpNotBefore, cert.NotBefore)
	ctx.Export(OpNotAfter, cert.NotAfter)
	ctx.Export(OpCertificateType, cert.Type)
	return nil
}
