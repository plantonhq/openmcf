package keycloak

import (
	"strings"
	"testing"

	"github.com/plantonhq/planton/operator/internal/resources"
)

func testLDAPFederation() *OwnedLDAPFederation {
	return &OwnedLDAPFederation{
		Servers:            []string{"ldaps://dc1.lab.example.internal:636", "ldaps://dc2.lab.example.internal:636"},
		BindDN:             "CN=svc-planton,OU=Service Accounts,DC=lab,DC=example,DC=internal",
		BindCredential:     "bind-password",
		UsersDN:            "OU=Staff,DC=lab,DC=example,DC=internal",
		GroupsDN:           "OU=Groups,DC=lab,DC=example,DC=internal",
		UserObjectClasses:  []string{"person", "organizationalPerson", "user"},
		UsernameAttribute:  "sAMAccountName",
		EmailAttribute:     "mail",
		FirstNameAttribute: "givenName",
		LastNameAttribute:  "sn",
		GroupNameAttribute: "cn",
		GroupMemberAttr:    "member",
		NestedGroups:       true,
		SyncPeriodMinutes:  60,
	}
}

// The component config's load-bearing invariants: multiple servers ride one
// space-joined connectionUrl, the sync window converts to seconds, the
// directory is never writable, and objectGUID is the immutable id.
func TestLDAPComponentConfig(t *testing.T) {
	cfg := testLDAPFederation().componentConfig()

	if got := cfg["connectionUrl"][0]; got != "ldaps://dc1.lab.example.internal:636 ldaps://dc2.lab.example.internal:636" {
		t.Errorf("connectionUrl = %q, want space-joined server list", got)
	}
	if got := cfg["fullSyncPeriod"][0]; got != "3600" {
		t.Errorf("fullSyncPeriod = %q, want 3600 (60 minutes in seconds)", got)
	}
	if got := cfg["editMode"][0]; got != "READ_ONLY" {
		t.Errorf("editMode = %q -- the directory is the truth, Planton never writes it", got)
	}
	if got := cfg["uuidLDAPAttribute"][0]; got != "objectGUID" {
		t.Errorf("uuidLDAPAttribute = %q, want objectGUID (the stable id)", got)
	}
	if _, ok := cfg["bindCredential"]; ok {
		t.Error("bindCredential must NOT be in the diffable owned config -- it is masked on read and rides create/rotation writes only")
	}
}

// The truststore posture follows the CA declaration: a private CA pins
// useTruststoreSpi=always; without one, ldapsOnly rides the JVM defaults.
func TestLDAPComponentConfig_Truststore(t *testing.T) {
	fed := testLDAPFederation()
	if got := fed.componentConfig()["useTruststoreSpi"][0]; got != "ldapsOnly" {
		t.Errorf("useTruststoreSpi = %q without a CA bundle, want ldapsOnly", got)
	}
	fed.UseTruststore = true
	if got := fed.componentConfig()["useTruststoreSpi"][0]; got != "always" {
		t.Errorf("useTruststoreSpi = %q with a CA bundle, want always", got)
	}
}

// The group mapper's two stable-id guarantees: directory groups mirror under
// the dedicated /directory subtree (drop-non-existing stays scoped there,
// structurally away from admin-created realm groups), and each mirrored
// group carries the directory objectGUID -- the id the platform's mappings
// key on. Nested resolution follows the manifest.
func TestLDAPGroupMapperConfig(t *testing.T) {
	fed := testLDAPFederation()
	cfg := fed.groupMapperConfig()

	if got := cfg["groups.path"][0]; got != resources.IdentityDirectoryGroupsPath {
		t.Errorf("groups.path = %q, want %s", got, resources.IdentityDirectoryGroupsPath)
	}
	if got := cfg["mapped.group.attributes"][0]; got != "objectGUID" {
		t.Errorf("mapped.group.attributes = %q, want objectGUID -- without it, group identity on the LDAP arm is display names, which rename and collide", got)
	}
	if got := cfg["user.roles.retrieve.strategy"][0]; got != "LOAD_GROUPS_BY_MEMBER_ATTRIBUTE_RECURSIVELY" {
		t.Errorf("retrieve strategy = %q with nestedGroups, want RECURSIVELY", got)
	}
	if got := cfg["preserve.group.inheritance"][0]; got != "false" {
		t.Errorf("preserve.group.inheritance = %q, want false ALWAYS -- a multi-parent directory group (normal in real ADs) fails a nested sync with GroupsMultipleParents", got)
	}

	fed.NestedGroups = false
	if got := fed.groupMapperConfig()["user.roles.retrieve.strategy"][0]; got != "LOAD_GROUPS_BY_MEMBER_ATTRIBUTE" {
		t.Errorf("retrieve strategy = %q without nestedGroups, want non-recursive", got)
	}
}

// The broker config's load-bearing invariants: syncMode FORCE (import-once
// would freeze a user's groups at first login, silently breaking offboarding
// on the brokered arm), signature validation on, and endpoint fields present
// ONLY when the caller ran discovery -- the steady-state pass never fetches
// the upstream's discovery document.
func TestBrokerConfig(t *testing.T) {
	broker := &OwnedOIDCBroker{
		IssuerURL:   "https://login.microsoftonline.com/tenant-id/v2.0",
		ClientID:    "client-id",
		Scopes:      []string{"openid", "profile", "email"},
		GroupsClaim: "groups",
		DisplayName: "Sign in with your organization",
	}

	cfg := broker.brokerConfig()
	if cfg["syncMode"] != "FORCE" {
		t.Errorf("syncMode = %q, want FORCE", cfg["syncMode"])
	}
	if cfg["validateSignature"] != "true" {
		t.Error("upstream tokens must be signature-validated")
	}
	if _, ok := cfg["authorizationUrl"]; ok {
		t.Error("endpoint fields must be absent without fresh discovery (owned-when-discovered)")
	}
	if _, ok := cfg["clientSecret"]; ok {
		t.Error("clientSecret must NOT be in the diffable owned config -- masked on read, rides create/rotation writes only")
	}

	broker.Endpoints = &OIDCEndpoints{
		AuthorizationURL: "https://login/authorize", TokenURL: "https://login/token", JWKSURL: "https://login/keys",
	}
	cfg = broker.brokerConfig()
	if cfg["authorizationUrl"] != "https://login/authorize" || cfg["tokenUrl"] != "https://login/token" || cfg["jwksUrl"] != "https://login/keys" {
		t.Errorf("endpoint fields not carried from discovery: %v", cfg)
	}
	if _, ok := cfg["userInfoUrl"]; ok {
		t.Error("a spec-optional endpoint absent from discovery must stay absent from config")
	}
}

// The broker's groups importer re-imports at EVERY sign-in and lands the
// upstream claim in the one attribute the protocol mapper reads back out.
func TestBrokerGroupsImporterConfig(t *testing.T) {
	broker := &OwnedOIDCBroker{GroupsClaim: "groups"}
	cfg := broker.groupsImporterConfig()
	if cfg["syncMode"] != "FORCE" {
		t.Errorf("importer syncMode = %q, want FORCE (import-once freezes groups at first login)", cfg["syncMode"])
	}
	if cfg["user.attribute"] != resources.IdentityDirectoryGroupsAttribute {
		t.Errorf("importer user.attribute = %q, want %s", cfg["user.attribute"], resources.IdentityDirectoryGroupsAttribute)
	}
}

// The broker's subject importer mirrors the groups importer exactly: FORCE
// re-import (a stale directory id after an upstream re-keying would break
// the account correlation offboarding runs on), landing in the one attribute
// the subject protocol mapper reads back out.
func TestBrokerSubjectImporterConfig(t *testing.T) {
	broker := &OwnedOIDCBroker{SubjectClaim: "oid"}
	cfg := broker.subjectImporterConfig()
	if cfg["syncMode"] != "FORCE" {
		t.Errorf("importer syncMode = %q, want FORCE (a stale directory id breaks offboarding correlation)", cfg["syncMode"])
	}
	if cfg["claim"] != "oid" {
		t.Errorf("importer claim = %q, want the manifest's subjectClaim", cfg["claim"])
	}
	if cfg["user.attribute"] != resources.IdentityDirectorySubjectAttribute {
		t.Errorf("importer user.attribute = %q, want %s", cfg["user.attribute"], resources.IdentityDirectorySubjectAttribute)
	}
}

// The directory-subject mapper is the structural directory-born mark: every
// federated user carries the claim (LDAP_ID exists on every LDAP-federated
// user; the subject importer fills the attribute on every brokered sign-in),
// local realm users never do. Same mapper type on both arms; only the source
// attribute differs; one shared name so the sweep owns it across arm
// switches.
func TestSubjectProtocolMapper_ArmSpecific(t *testing.T) {
	ldapMapper := subjectProtocolMapper(&OwnedFederation{LDAP: testLDAPFederation()})
	if ldapMapper == nil || ldapMapper.Config["user.attribute"] != "LDAP_ID" {
		t.Fatalf("LDAP arm subject mapper = %+v, want the LDAP_ID source attribute", ldapMapper)
	}
	if ldapMapper.Config["claim.name"] != resources.IdentityDirectorySubjectClaim {
		t.Errorf("LDAP arm claim = %q, want the canonical %s", ldapMapper.Config["claim.name"], resources.IdentityDirectorySubjectClaim)
	}
	if ldapMapper.Config["multivalued"] != "false" {
		t.Error("the directory subject is ONE id, never a list")
	}
	if ldapMapper.Config["userinfo.token.claim"] != "true" {
		t.Error("the claim must ride userinfo -- the control plane's provisioning reads it there")
	}

	brokerMapper := subjectProtocolMapper(&OwnedFederation{Broker: &OwnedOIDCBroker{SubjectClaim: "sub"}})
	if brokerMapper == nil || brokerMapper.Config["user.attribute"] != resources.IdentityDirectorySubjectAttribute {
		t.Fatalf("broker arm subject mapper = %+v, want the imported %s attribute", brokerMapper, resources.IdentityDirectorySubjectAttribute)
	}
	if brokerMapper.Name != ldapMapper.Name {
		t.Error("both arms must share ONE mapper name -- it is how the sweep owns the mapper across arm switches")
	}

	if got := subjectProtocolMapper(&OwnedFederation{}); got != nil {
		t.Errorf("no federation desired -> no mapper desired, got %+v", got)
	}
}

// ONE canonical claim, arm-specific emitting mapper: the LDAP arm's groups
// live as realm groups (group-membership mapper), the brokered arm's in a
// user attribute (user-attribute mapper). One mapper type cannot serve both;
// the NAME is shared so the sweep owns it regardless of arm.
func TestGroupsProtocolMapper_ArmSpecific(t *testing.T) {
	ldapMapper := groupsProtocolMapper(&OwnedFederation{LDAP: testLDAPFederation()})
	if ldapMapper == nil || ldapMapper.ProtocolMapper != "oidc-group-membership-mapper" {
		t.Fatalf("LDAP arm mapper = %+v, want oidc-group-membership-mapper", ldapMapper)
	}
	if ldapMapper.Config["claim.name"] != resources.IdentityGroupsClaim {
		t.Errorf("LDAP arm claim = %q, want the canonical %s", ldapMapper.Config["claim.name"], resources.IdentityGroupsClaim)
	}

	brokerMapper := groupsProtocolMapper(&OwnedFederation{Broker: &OwnedOIDCBroker{GroupsClaim: "groups"}})
	if brokerMapper == nil || brokerMapper.ProtocolMapper != "oidc-usermodel-attribute-mapper" {
		t.Fatalf("broker arm mapper = %+v, want oidc-usermodel-attribute-mapper", brokerMapper)
	}
	if brokerMapper.Config["claim.name"] != resources.IdentityGroupsClaim {
		t.Errorf("broker arm claim = %q, want the canonical %s", brokerMapper.Config["claim.name"], resources.IdentityGroupsClaim)
	}
	if brokerMapper.Name != ldapMapper.Name {
		t.Error("both arms must share ONE mapper name -- it is how the sweep owns the mapper across arm switches")
	}

	if got := groupsProtocolMapper(&OwnedFederation{}); got != nil {
		t.Errorf("no federation desired -> no mapper desired, got %+v", got)
	}
}

// The Entra tenancy check names the multi-tenant constraint precisely and
// stays silent for non-Microsoft issuers (nothing to say is a verdict too).
func TestEntraTenancyCheck(t *testing.T) {
	if check := entraTenancyCheck("https://login.microsoftonline.com/common/v2.0"); check == nil || check.Verdict != VerdictFailed {
		t.Fatalf("/common must fail the tenancy check, got %+v", check)
	} else if !strings.Contains(check.Message, "single-tenant") {
		t.Errorf("the failure must name the remedy (single-tenant form): %q", check.Message)
	}
	if check := entraTenancyCheck("https://login.microsoftonline.com/3f2504e0-tenant/v2.0"); check == nil || check.Verdict != VerdictPassed {
		t.Errorf("a tenant-scoped issuer must pass, got %+v", check)
	}
	if check := entraTenancyCheck("https://keycloak.lab.example.internal/idp/realms/entra-sim"); check != nil {
		t.Errorf("non-Microsoft issuers get no tenancy check, got %+v", check)
	}
}
