# AWS Cert Manager Cert Examples

Below are several examples demonstrating how to define an AWS Cert Manager Cert component in
OpenMCF. After creating one of these YAML manifests, apply it with Terraform using the OpenMCF CLI:

```shell
openmcf tofu apply --manifest <yaml-path> --stack <stack-name>
```

---

## Basic Example (Single Domain)

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsCertManagerCert
metadata:
  name: single-cert-example
spec:
  primaryDomainName: example.com
  route53HostedZoneId: Z123456ABCXYZ
```

This manifest requests an SSL certificate for `example.com` and uses Route53 hosted zone `Z123456ABCXYZ`
for DNS validation. The validation method defaults to "DNS".

---

## Multiple Domains Example

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsCertManagerCert
metadata:
  name: multi-cert-example
spec:
  primaryDomainName: myapp.io
  alternateDomainNames:
    - api.myapp.io
    - blog.myapp.io
  route53HostedZoneId: Z987654XYZABC
  validationMethod: DNS
```

In this example, we issue a certificate for `myapp.io` as the primary domain and add two alternate domains:
`api.myapp.io` and `blog.myapp.io`. It automatically creates DNS records in the
Route53 hosted zone `Z987654XYZABC`.

---

## Wildcard Domain Example

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsCertManagerCert
metadata:
  name: wildcard-example
spec:
  primaryDomainName: '*.wildexample.com'
  route53HostedZoneId: Z999999EXAMPLE
```

This manifest requests a wildcard certificate for `*.wildexample.com`, allowing you to secure any subdomain under
`wildexample.com` without creating separate certificates.

---

## Email Validation Example

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsCertManagerCert
metadata:
  name: email-validation-example
spec:
  primaryDomainName: emailvalidated.com
  alternateDomainNames:
    - www.emailvalidated.com
  validationMethod: EMAIL
```

This example uses email validation instead of DNS validation. AWS will send validation emails to the domain's
administrative contact addresses (admin@, administrator@, hostmaster@, postmaster@, webmaster@).

---

## Complete Production Certificate

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsCertManagerCert
metadata:
  name: production-cert
spec:
  primaryDomainName: app.production.com
  alternateDomainNames:
    - api.production.com
    - admin.production.com
    - '*.api.production.com'
  route53HostedZoneId: Z1234567890ABCDEF
  validationMethod: DNS
```

A production-ready certificate with:
• Primary domain and multiple alternate domains including a wildcard subdomain.
• DNS validation for automatic certificate renewal.
• Route53 integration for DNS record management.

---

## Minimal Example

```yaml
apiVersion: aws.openmcf.org/v1
kind: AwsCertManagerCert
metadata:
  name: minimal-cert
spec:
  primaryDomainName: barebones.io
  route53HostedZoneId: Z1111111111111
```

Here, we only specify the required fields. The certificate will be created for `barebones.io`, with DNS records
automatically managed in the hosted zone. Alternate domain names are not specified.

---

## After Deploying

Once you've applied your manifest with OpenMCF tofu, you can confirm that the certificate is active via the AWS console or by
using the AWS CLI:

```shell
aws acm list-certificates
```

You should see your new certificate and any domain validation records created within Route53. The certificate will be
automatically renewed by AWS before expiration.

---

Enjoy secure connections with OpenMCF!
