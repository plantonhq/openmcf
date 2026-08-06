package module

import (
	"github.com/pkg/errors"
	kubernetescertificatev1alpha1 "github.com/plantonhq/planton/apis/dev/planton/provider/kubernetes/kubernetescertificate/v1alpha1"
	certmanagerv1 "github.com/plantonhq/planton/pkg/kubernetes/kubernetestypes/certmanager/kubernetes/cert_manager/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// buildCertificateSpec maps the proto spec onto the typed crd2pulumi
// CertificateSpecArgs. Proto enum vocabularies are lowercase; the CRD wants
// its exact casing ("RSA", "PKCS1", "Always", "DER", "LegacyRC2", ...) — the
// translation tables here are the single place that mapping lives for the
// Pulumi engine; the Terraform module's locals carry the identical tables.
func buildCertificateSpec(locals *Locals) (certmanagerv1.CertificateSpecArgs, error) {
	spec := locals.Spec

	out := certmanagerv1.CertificateSpecArgs{
		SecretName: pulumi.String(spec.GetSecretName()),
	}

	// ---- names -----------------------------------------------------------
	if len(spec.GetDnsNames()) > 0 {
		out.DnsNames = pulumi.ToStringArray(spec.GetDnsNames())
	}
	if len(spec.GetIpAddresses()) > 0 {
		out.IpAddresses = pulumi.ToStringArray(spec.GetIpAddresses())
	}
	if len(spec.GetUris()) > 0 {
		out.Uris = pulumi.ToStringArray(spec.GetUris())
	}
	if len(spec.GetEmailAddresses()) > 0 {
		out.EmailAddresses = pulumi.ToStringArray(spec.GetEmailAddresses())
	}
	if spec.GetCommonName() != "" {
		out.CommonName = pulumi.StringPtr(spec.GetCommonName())
	}
	if spec.GetLiteralSubject() != "" {
		out.LiteralSubject = pulumi.StringPtr(spec.GetLiteralSubject())
	}
	if subject := spec.GetSubject(); subject != nil {
		subjectArgs := &certmanagerv1.CertificateSpecSubjectArgs{}
		if len(subject.GetOrganizations()) > 0 {
			subjectArgs.Organizations = pulumi.ToStringArray(subject.GetOrganizations())
		}
		if len(subject.GetOrganizationalUnits()) > 0 {
			subjectArgs.OrganizationalUnits = pulumi.ToStringArray(subject.GetOrganizationalUnits())
		}
		if len(subject.GetCountries()) > 0 {
			subjectArgs.Countries = pulumi.ToStringArray(subject.GetCountries())
		}
		if len(subject.GetProvinces()) > 0 {
			subjectArgs.Provinces = pulumi.ToStringArray(subject.GetProvinces())
		}
		if len(subject.GetLocalities()) > 0 {
			subjectArgs.Localities = pulumi.ToStringArray(subject.GetLocalities())
		}
		if len(subject.GetStreetAddresses()) > 0 {
			subjectArgs.StreetAddresses = pulumi.ToStringArray(subject.GetStreetAddresses())
		}
		if len(subject.GetPostalCodes()) > 0 {
			subjectArgs.PostalCodes = pulumi.ToStringArray(subject.GetPostalCodes())
		}
		if subject.GetSerialNumber() != "" {
			subjectArgs.SerialNumber = pulumi.StringPtr(subject.GetSerialNumber())
		}
		out.Subject = subjectArgs
	}
	if len(spec.GetOtherNames()) > 0 {
		otherNames := make(certmanagerv1.CertificateSpecOtherNamesArray, 0, len(spec.GetOtherNames()))
		for _, otherName := range spec.GetOtherNames() {
			otherNames = append(otherNames, certmanagerv1.CertificateSpecOtherNamesArgs{
				Oid:       pulumi.StringPtr(otherName.GetOid()),
				Utf8Value: pulumi.StringPtr(otherName.GetUtf8Value()),
			})
		}
		out.OtherNames = otherNames
	}

	// ---- issuer ----------------------------------------------------------
	issuerRef, err := buildIssuerRef(spec.GetIssuerRef())
	if err != nil {
		return out, err
	}
	out.IssuerRef = issuerRef

	// ---- lifetime and renewal ---------------------------------------------
	if spec.Duration != nil {
		out.Duration = pulumi.StringPtr(spec.GetDuration())
	}
	if spec.GetRenewBefore() != "" {
		out.RenewBefore = pulumi.StringPtr(spec.GetRenewBefore())
	}
	if spec.RenewBeforePercentage != nil {
		out.RenewBeforePercentage = pulumi.IntPtr(int(spec.GetRenewBeforePercentage()))
	}
	if spec.RevisionHistoryLimit != nil {
		out.RevisionHistoryLimit = pulumi.IntPtr(int(spec.GetRevisionHistoryLimit()))
	}

	// ---- key material -------------------------------------------------------
	if pk := spec.GetPrivateKey(); pk != nil {
		privateKey := &certmanagerv1.CertificateSpecPrivateKeyArgs{}
		if pk.Algorithm != nil {
			switch pk.GetAlgorithm() {
			case kubernetescertificatev1alpha1.KubernetesCertificatePrivateKey_rsa:
				privateKey.Algorithm = pulumi.StringPtr("RSA")
			case kubernetescertificatev1alpha1.KubernetesCertificatePrivateKey_ecdsa:
				privateKey.Algorithm = pulumi.StringPtr("ECDSA")
			case kubernetescertificatev1alpha1.KubernetesCertificatePrivateKey_ed25519:
				privateKey.Algorithm = pulumi.StringPtr("Ed25519")
			}
		}
		if pk.Size != nil {
			privateKey.Size = pulumi.IntPtr(int(pk.GetSize()))
		}
		if pk.Encoding != nil {
			switch pk.GetEncoding() {
			case kubernetescertificatev1alpha1.KubernetesCertificatePrivateKey_pkcs1:
				privateKey.Encoding = pulumi.StringPtr("PKCS1")
			case kubernetescertificatev1alpha1.KubernetesCertificatePrivateKey_pkcs8:
				privateKey.Encoding = pulumi.StringPtr("PKCS8")
			}
		}
		if pk.RotationPolicy != nil {
			switch pk.GetRotationPolicy() {
			case kubernetescertificatev1alpha1.KubernetesCertificatePrivateKey_always:
				privateKey.RotationPolicy = pulumi.StringPtr("Always")
			case kubernetescertificatev1alpha1.KubernetesCertificatePrivateKey_never:
				privateKey.RotationPolicy = pulumi.StringPtr("Never")
			}
		}
		out.PrivateKey = privateKey
	}

	// ---- usages ----------------------------------------------------------
	if len(spec.GetUsages()) > 0 {
		out.Usages = pulumi.ToStringArray(spec.GetUsages())
	}
	if spec.GetEncodeUsagesInRequest() {
		out.EncodeUsagesInRequest = pulumi.BoolPtr(true)
	}
	if spec.GetIsCa() {
		out.IsCA = pulumi.BoolPtr(true)
	}
	if spec.GetSignatureAlgorithm() != "" {
		out.SignatureAlgorithm = pulumi.StringPtr(spec.GetSignatureAlgorithm())
	}

	// ---- outputs beyond PEM ----------------------------------------------
	if keystores := spec.GetKeystores(); keystores != nil {
		keystoresArgs := &certmanagerv1.CertificateSpecKeystoresArgs{}
		if jks := keystores.GetJks(); jks != nil {
			keystoresArgs.Jks = &certmanagerv1.CertificateSpecKeystoresJksArgs{
				Create:   pulumi.BoolPtr(jks.GetCreate()),
				Alias:    pulumi.StringPtr(jks.GetAlias()),
				Password: pulumi.ToSecret(pulumi.String(jks.GetPassword())).(pulumi.StringOutput),
			}
		}
		if pkcs12 := keystores.GetPkcs12(); pkcs12 != nil {
			pkcs12Args := &certmanagerv1.CertificateSpecKeystoresPkcs12Args{
				Create:   pulumi.BoolPtr(pkcs12.GetCreate()),
				Password: pulumi.ToSecret(pulumi.String(pkcs12.GetPassword())).(pulumi.StringOutput),
			}
			switch pkcs12.GetProfile() {
			case "legacy_rc2":
				pkcs12Args.Profile = pulumi.StringPtr("LegacyRC2")
			case "legacy_des":
				pkcs12Args.Profile = pulumi.StringPtr("LegacyDES")
			case "modern2023":
				pkcs12Args.Profile = pulumi.StringPtr("Modern2023")
			}
			keystoresArgs.Pkcs12 = pkcs12Args
		}
		out.Keystores = keystoresArgs
	}
	if len(spec.GetAdditionalOutputFormats()) > 0 {
		formats := make(certmanagerv1.CertificateSpecAdditionalOutputFormatsArray, 0, len(spec.GetAdditionalOutputFormats()))
		for _, format := range spec.GetAdditionalOutputFormats() {
			crdType := "DER"
			if format.GetType() == "combined_pem" {
				crdType = "CombinedPEM"
			}
			formats = append(formats, certmanagerv1.CertificateSpecAdditionalOutputFormatsArgs{
				Type: pulumi.StringPtr(crdType),
			})
		}
		out.AdditionalOutputFormats = formats
	}

	// ---- CA name constraints ------------------------------------------------
	if nameConstraints := spec.GetNameConstraints(); nameConstraints != nil {
		nameConstraintsArgs := &certmanagerv1.CertificateSpecNameConstraintsArgs{
			Critical: pulumi.BoolPtr(nameConstraints.GetCritical()),
		}
		if permitted := nameConstraints.GetPermitted(); permitted != nil {
			nameConstraintsArgs.Permitted = &certmanagerv1.CertificateSpecNameConstraintsPermittedArgs{
				DnsDomains:     pulumi.ToStringArray(permitted.GetDnsDomains()),
				IpRanges:       pulumi.ToStringArray(permitted.GetIpRanges()),
				EmailAddresses: pulumi.ToStringArray(permitted.GetEmailAddresses()),
				UriDomains:     pulumi.ToStringArray(permitted.GetUriDomains()),
			}
		}
		if excluded := nameConstraints.GetExcluded(); excluded != nil {
			nameConstraintsArgs.Excluded = &certmanagerv1.CertificateSpecNameConstraintsExcludedArgs{
				DnsDomains:     pulumi.ToStringArray(excluded.GetDnsDomains()),
				IpRanges:       pulumi.ToStringArray(excluded.GetIpRanges()),
				EmailAddresses: pulumi.ToStringArray(excluded.GetEmailAddresses()),
				UriDomains:     pulumi.ToStringArray(excluded.GetUriDomains()),
			}
		}
		out.NameConstraints = nameConstraintsArgs
	}

	// ---- secret template ------------------------------------------------------
	if secretTemplate := spec.GetSecretTemplate(); secretTemplate != nil {
		templateArgs := &certmanagerv1.CertificateSpecSecretTemplateArgs{}
		if len(secretTemplate.GetLabels()) > 0 {
			templateArgs.Labels = pulumi.ToStringMap(secretTemplate.GetLabels())
		}
		if len(secretTemplate.GetAnnotations()) > 0 {
			templateArgs.Annotations = pulumi.ToStringMap(secretTemplate.GetAnnotations())
		}
		out.SecretTemplate = templateArgs
	}

	return out, nil
}

// buildIssuerRef maps the issuer-selection oneof onto the CRD's issuerRef.
// Planton-managed issuers carry group cert-manager.io implicitly; the
// external arm names a third-party issuer kind explicitly.
func buildIssuerRef(issuerRef *kubernetescertificatev1alpha1.KubernetesCertificateIssuerRef) (*certmanagerv1.CertificateSpecIssuerRefArgs, error) {
	switch {
	case issuerRef.GetClusterIssuer() != nil:
		return &certmanagerv1.CertificateSpecIssuerRefArgs{
			Kind: pulumi.StringPtr("ClusterIssuer"),
			Name: pulumi.StringPtr(issuerRef.GetClusterIssuer().GetName().GetValue()),
		}, nil
	case issuerRef.GetIssuer() != nil:
		return &certmanagerv1.CertificateSpecIssuerRefArgs{
			Kind: pulumi.StringPtr("Issuer"),
			Name: pulumi.StringPtr(issuerRef.GetIssuer().GetName().GetValue()),
		}, nil
	case issuerRef.GetExternal() != nil:
		external := issuerRef.GetExternal()
		return &certmanagerv1.CertificateSpecIssuerRefArgs{
			Group: pulumi.StringPtr(external.GetGroup()),
			Kind:  pulumi.StringPtr(external.GetKind()),
			Name:  pulumi.StringPtr(external.GetName()),
		}, nil
	default:
		return nil, errors.New("issuer_ref must select cluster_issuer, issuer, or external")
	}
}
