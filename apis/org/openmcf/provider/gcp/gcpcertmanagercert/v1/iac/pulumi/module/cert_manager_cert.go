package module

import (
	"fmt"

	"github.com/pkg/errors"
	gcpcertmanagercertv1 "github.com/plantonhq/openmcf/apis/org/openmcf/provider/gcp/gcpcertmanagercert/v1"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/certificatemanager"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/compute"
	"github.com/pulumi/pulumi-gcp/sdk/v9/go/gcp/dns"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// certManagerCert provisions a GCP certificate with DNS validation
// and creates corresponding DNS records in Cloud DNS.
// It supports both Certificate Manager certificates and Load Balancer SSL certificates.
func certManagerCert(ctx *pulumi.Context, locals *Locals, provider *gcp.Provider) error {
	spec := locals.GcpCertManagerCert.Spec

	// Determine which type of certificate to create
	// Default to MANAGED if not specified
	certType := gcpcertmanagercertv1.CertificateType_MANAGED
	if spec.CertificateType != nil {
		certType = *spec.CertificateType
	}

	// Collect all domains (primary + alternates)
	allDomains := []string{spec.PrimaryDomainName}
	allDomains = append(allDomains, spec.AlternateDomainNames...)

	switch certType {
	case gcpcertmanagercertv1.CertificateType_MANAGED:
		return createManagedCertificate(ctx, locals, provider, allDomains, spec)
	case gcpcertmanagercertv1.CertificateType_LOAD_BALANCER:
		return createLoadBalancerCertificate(ctx, locals, provider, allDomains, spec)
	default:
		return errors.Errorf("unsupported certificate type: %v", certType)
	}
}

// createManagedCertificate creates a Certificate Manager certificate with DNS authorization
func createManagedCertificate(ctx *pulumi.Context, locals *Locals, provider *gcp.Provider,
	allDomains []string, spec *gcpcertmanagercertv1.GcpCertManagerCertSpec) error {

	meta := locals.GcpCertManagerCert.Metadata

	// Create DNS authorizations for each domain
	var dnsAuthorizations []*certificatemanager.DnsAuthorization
	for i, domain := range allDomains {
		dnsAuth, err := certificatemanager.NewDnsAuthorization(ctx,
			fmt.Sprintf("%s-dns-auth-%d", meta.Name, i),
			&certificatemanager.DnsAuthorizationArgs{
				Name:        pulumi.String(fmt.Sprintf("%s-dns-auth-%d", meta.Name, i)),
				Description: pulumi.Sprintf("DNS authorization for %s", domain),
				Domain:      pulumi.String(domain),
				Project:     pulumi.String(spec.GcpProjectId),
				Labels:      pulumi.ToStringMap(locals.GcpLabels),
			},
			pulumi.Provider(provider),
		)
		if err != nil {
			return errors.Wrapf(err, "failed to create DNS authorization for domain %s", domain)
		}
		dnsAuthorizations = append(dnsAuthorizations, dnsAuth)

		// Create the DNS record for validation
		// The DNS authorization provides the record name and data
		_, err = dns.NewRecordSet(ctx,
			fmt.Sprintf("%s-validation-record-%d", meta.Name, i),
			&dns.RecordSetArgs{
				Name: dnsAuth.DnsResourceRecords.ApplyT(func(records []certificatemanager.DnsAuthorizationDnsResourceRecord) string {
					if len(records) > 0 && records[0].Name != nil {
						return *records[0].Name
					}
					return ""
				}).(pulumi.StringOutput),
				Type: dnsAuth.DnsResourceRecords.ApplyT(func(records []certificatemanager.DnsAuthorizationDnsResourceRecord) string {
					if len(records) > 0 && records[0].Type != nil {
						return *records[0].Type
					}
					return ""
				}).(pulumi.StringOutput),
				Ttl: pulumi.Int(300),
				Rrdatas: dnsAuth.DnsResourceRecords.ApplyT(func(records []certificatemanager.DnsAuthorizationDnsResourceRecord) []string {
					if len(records) > 0 && records[0].Data != nil {
						return []string{*records[0].Data}
					}
					return []string{}
				}).(pulumi.StringArrayOutput),
				ManagedZone: pulumi.String(spec.CloudDnsZoneId.GetValue()),
				Project:     pulumi.String(spec.GcpProjectId),
			},
			pulumi.Provider(provider),
		)
		if err != nil {
			return errors.Wrapf(err, "failed to create DNS validation record for domain %s", domain)
		}
	}

	// Create the Certificate Manager certificate
	cert, err := certificatemanager.NewCertificate(ctx,
		meta.Name+"-cert",
		&certificatemanager.CertificateArgs{
			Name:        pulumi.String(meta.Name),
			Description: pulumi.Sprintf("SSL certificate for %s", spec.PrimaryDomainName),
			Project:     pulumi.String(spec.GcpProjectId),
			Labels:      pulumi.ToStringMap(locals.GcpLabels),
			Managed: &certificatemanager.CertificateManagedArgs{
				Domains: pulumi.ToStringArray(allDomains),
				DnsAuthorizations: pulumi.StringArray(
					func() []pulumi.StringInput {
						var authIds []pulumi.StringInput
						for _, auth := range dnsAuthorizations {
							authIds = append(authIds, auth.ID())
						}
						return authIds
					}(),
				),
			},
		},
		pulumi.Provider(provider),
		pulumi.DependsOn(func() []pulumi.Resource {
			resources := make([]pulumi.Resource, len(dnsAuthorizations))
			for i, auth := range dnsAuthorizations {
				resources[i] = auth
			}
			return resources
		}()),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create Certificate Manager certificate")
	}

	// Export outputs
	ctx.Export(OpCertificateId, cert.ID())
	ctx.Export(OpCertificateName, cert.Name)
	ctx.Export(OpCertificateDomainName, pulumi.String(spec.PrimaryDomainName))
	ctx.Export(OpCertificateStatus, pulumi.String("PROVISIONING"))

	return nil
}

// createLoadBalancerCertificate creates a Google-managed SSL certificate for load balancers
func createLoadBalancerCertificate(ctx *pulumi.Context, locals *Locals, provider *gcp.Provider,
	allDomains []string, spec *gcpcertmanagercertv1.GcpCertManagerCertSpec) error {

	meta := locals.GcpCertManagerCert.Metadata

	// Note: Google-managed SSL certificates for load balancers handle DNS validation automatically
	// when the domain is pointed to the load balancer. We don't need to create DNS records manually.
	cert, err := compute.NewManagedSslCertificate(ctx,
		meta.Name+"-ssl-cert",
		&compute.ManagedSslCertificateArgs{
			Name:        pulumi.String(meta.Name),
			Description: pulumi.Sprintf("SSL certificate for %s", spec.PrimaryDomainName),
			Project:     pulumi.String(spec.GcpProjectId),
			Managed: &compute.ManagedSslCertificateManagedArgs{
				Domains: pulumi.ToStringArray(allDomains),
			},
		},
		pulumi.Provider(provider),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create Google-managed SSL certificate")
	}

	// Export outputs
	ctx.Export(OpCertificateId, cert.ID())
	ctx.Export(OpCertificateName, cert.Name)
	ctx.Export(OpCertificateDomainName, pulumi.String(spec.PrimaryDomainName))
	ctx.Export(OpCertificateStatus, pulumi.String("PROVISIONING"))

	return nil
}
