package keycloak

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"strconv"
	"strings"

	"github.com/plantonhq/planton/operator/internal/resources"
)

// The verification pass: every check the identity manifest's status reports
// is composed here, as a named verdict with a message a person can act on --
// a wrong bind credential is a verdict naming the remedy, never a stack
// trace. The identity server does all directory speaking (testLDAPConnection,
// user/group sync); the operator's only direct external fetch is the OIDC
// discovery document, which it needs anyway to write the broker's explicit
// endpoints.
//
// Verification runs on the manifest's cadence, not the reconcile's: the
// caller invokes it when the spec generation moved, a credential rotated, or
// federation state was repaired -- NEVER as a steady-state probe. A directory
// probed every 30 seconds is a security team's intrusion alert.

// Verdict is one check's outcome.
type Verdict string

const (
	VerdictPassed  Verdict = "Passed"
	VerdictFailed  Verdict = "Failed"
	VerdictUnknown Verdict = "Unknown"
)

// Check is one named verification verdict.
type Check struct {
	Name    string
	Verdict Verdict
	Message string
}

// VerifyInput carries everything one verification pass needs.
type VerifyInput struct {
	Realm      string
	ServerRoot string

	AdminUsername string
	AdminPassword string
	HTTPClient    *http.Client

	// Federation is the same desired state the convergence consumed; the
	// checks are arm-appropriate.
	Federation *OwnedFederation

	// SeededAdminEmail is the platform manifest's declared adminEmail (empty
	// when none): the seeded LOCAL realm user it creates collides with a
	// directory user holding the same email, and the collision is reported
	// here -- never auto-resolved (both objects are legitimate; deleting
	// either is a human call).
	SeededAdminEmail string
}

// Verify runs the arm-appropriate checks and returns the verdicts in the
// order they ran. An error is returned only when verification itself could
// not run (the admin API unreachable) -- a failing CHECK is a verdict, not
// an error.
func Verify(ctx context.Context, in VerifyInput) ([]Check, error) {
	if in.Federation == nil {
		return nil, nil
	}
	admin := NewAdminClient(in.HTTPClient, in.ServerRoot)
	if err := admin.Authenticate(ctx, in.AdminUsername, in.AdminPassword); err != nil {
		return nil, fmt.Errorf("authenticating to the admin API: %w", err)
	}

	var checks []Check
	switch {
	case in.Federation.LDAP != nil:
		checks = verifyLDAP(ctx, admin, in.Realm, in.Federation.LDAP)
	case in.Federation.Broker != nil:
		checks = verifyBroker(ctx, in.HTTPClient, in.Federation.Broker)
	}

	if in.SeededAdminEmail != "" {
		checks = append(checks, verifySeededAdminCollision(ctx, admin, in.Realm, in.SeededAdminEmail))
	}
	return checks, nil
}

// verifyLDAP probes the directory through the identity server: reachability,
// the bind credential, then real user and group syncs whose counts make the
// verdicts concrete ("42 users visible", not "OK").
func verifyLDAP(ctx context.Context, admin *AdminClient, realm string, ldap *OwnedLDAPFederation) []Check {
	probe := Representation{
		"connectionUrl":     strings.Join(ldap.Servers, " "),
		"authType":          "simple",
		"bindDn":            ldap.BindDN,
		"startTls":          strconv.FormatBool(ldap.StartTLS),
		"useTruststoreSpi":  ldap.componentConfig()["useTruststoreSpi"][0],
		"connectionTimeout": "10000",
	}

	checks := make([]Check, 0, 4)

	connProbe := cloneRepresentation(probe)
	connProbe["action"] = "testConnection"
	if err := admin.TestLDAPConnection(ctx, realm, connProbe); err != nil {
		checks = append(checks, Check{Name: "connection", Verdict: VerdictFailed, Message: fmt.Sprintf(
			"the directory could not be reached at %s -- check network reach from the cluster, the server address, and (for ldaps) that the certificate chains to a trusted CA (declare caBundleSecretRef for a private CA): %s",
			strings.Join(ldap.Servers, " "), keycloakReason(err))})
		// Without a connection the remaining probes can only restate the
		// same failure; stopping keeps the verdict list signal, not noise.
		return checks
	}
	checks = append(checks, Check{Name: "connection", Verdict: VerdictPassed, Message: fmt.Sprintf(
		"the directory answered at %s", strings.Join(ldap.Servers, " "))})

	bindProbe := cloneRepresentation(probe)
	bindProbe["action"] = "testAuthentication"
	bindProbe["bindCredential"] = ldap.BindCredential
	if err := admin.TestLDAPConnection(ctx, realm, bindProbe); err != nil {
		checks = append(checks, Check{Name: "bind", Verdict: VerdictFailed, Message: fmt.Sprintf(
			"the service account %s could not authenticate -- check the bind DN and the password in the referenced Secret: %s",
			ldap.BindDN, keycloakReason(err))})
		return checks
	}
	checks = append(checks, Check{Name: "bind", Verdict: VerdictPassed, Message: fmt.Sprintf(
		"authenticated as %s", ldap.BindDN)})

	checks = append(checks, verifyLDAPSyncs(ctx, admin, realm)...)
	return checks
}

// verifyLDAPSyncs triggers a real user sync and a real group sync on the
// provisioned component and reports their counts.
func verifyLDAPSyncs(ctx context.Context, admin *AdminClient, realm string) []Check {
	realmRep, err := admin.GetRealm(ctx, realm)
	if err != nil {
		return []Check{{Name: "usersSearch", Verdict: VerdictUnknown, Message: "could not read the realm to locate the federation component: " + err.Error()}}
	}
	realmID, _ := realmRep["id"].(string)

	components, err := admin.ListComponents(ctx, realm, realmID, userStorageProviderType)
	if err != nil {
		return []Check{{Name: "usersSearch", Verdict: VerdictUnknown, Message: "could not list federation components: " + err.Error()}}
	}
	var componentID string
	for _, c := range components {
		if name, _ := c["name"].(string); name == resources.IdentityLDAPComponentName {
			componentID, _ = c["id"].(string)
			break
		}
	}
	if componentID == "" {
		return []Check{{Name: "usersSearch", Verdict: VerdictUnknown, Message: "the federation component has not been provisioned yet; the next reconcile pass verifies it"}}
	}

	checks := make([]Check, 0, 2)

	userSync, err := admin.TriggerUserStorageSync(ctx, realm, componentID, "triggerFullSync")
	switch {
	case err != nil:
		checks = append(checks, Check{Name: "usersSearch", Verdict: VerdictFailed, Message: "searching the users DN failed -- check usersDn and the service account's read permissions: " + keycloakReason(err)})
	case userSync.Failed > 0:
		checks = append(checks, Check{Name: "usersSearch", Verdict: VerdictFailed, Message: fmt.Sprintf(
			"%d users synced (%d added, %d updated) but %d FAILED -- failures are usually email collisions with existing local users (the seeded admin among them) or entries missing required attributes; the identity server's log names each one",
			userSync.Added+userSync.Updated, userSync.Added, userSync.Updated, userSync.Failed)})
	default:
		checks = append(checks, Check{Name: "usersSearch", Verdict: VerdictPassed, Message: fmt.Sprintf(
			"%d directory users are visible to the identity server (%d added, %d updated this sync)",
			userSync.Added+userSync.Updated, userSync.Added, userSync.Updated)})
	}

	mappers, err := admin.ListComponents(ctx, realm, componentID, ldapStorageMapperType)
	if err != nil {
		checks = append(checks, Check{Name: "groupsSearch", Verdict: VerdictUnknown, Message: "could not list the federation mappers: " + err.Error()})
		return checks
	}
	var mapperID string
	for _, m := range mappers {
		if name, _ := m["name"].(string); name == ldapGroupMapperName {
			mapperID, _ = m["id"].(string)
			break
		}
	}
	if mapperID == "" {
		checks = append(checks, Check{Name: "groupsSearch", Verdict: VerdictUnknown, Message: "the group mapper has not been provisioned yet; the next reconcile pass verifies it"})
		return checks
	}

	groupSync, err := admin.TriggerLDAPMapperSync(ctx, realm, componentID, mapperID, "fedToKeycloak")
	switch {
	case err != nil:
		checks = append(checks, Check{Name: "groupsSearch", Verdict: VerdictFailed, Message: "searching the groups DN failed -- check groupsDn and the service account's read permissions: " + keycloakReason(err)})
	case groupSync.Failed > 0:
		checks = append(checks, Check{Name: "groupsSearch", Verdict: VerdictFailed, Message: fmt.Sprintf(
			"%d groups synced but %d failed; the identity server's log names each one", groupSync.Added+groupSync.Updated, groupSync.Failed)})
	default:
		checks = append(checks, Check{Name: "groupsSearch", Verdict: VerdictPassed, Message: fmt.Sprintf(
			"%d directory groups mirrored under %s (%d added, %d updated this sync)",
			groupSync.Added+groupSync.Updated, resources.IdentityDirectoryGroupsPath, groupSync.Added, groupSync.Updated)})
	}
	return checks
}

// verifyBroker checks what CAN be known before anyone signs in: the issuer's
// discovery document (fetched by the operator, the same fetch that feeds the
// broker's explicit endpoints) and the Entra tenancy constraint. Client
// authentication is honestly Unknown -- there is no directory to browse on
// the brokered arm, and pretending otherwise is the failure mode this
// program exists to avoid.
func verifyBroker(ctx context.Context, httpClient *http.Client, broker *OwnedOIDCBroker) []Check {
	checks := make([]Check, 0, 3)

	if tenancy := entraTenancyCheck(broker.IssuerURL); tenancy != nil {
		checks = append(checks, *tenancy)
	}

	endpoints, advertisedIssuer, err := DiscoverOIDC(ctx, httpClient, broker.IssuerURL)
	switch {
	case err != nil:
		checks = append(checks, Check{Name: "issuer", Verdict: VerdictFailed, Message: fmt.Sprintf(
			"the issuer's discovery document could not be fetched from %s/.well-known/openid-configuration -- check the issuer URL and the cluster's egress: %v",
			strings.TrimSuffix(broker.IssuerURL, "/"), err)})
	case advertisedIssuer != broker.IssuerURL && advertisedIssuer != strings.TrimSuffix(broker.IssuerURL, "/"):
		checks = append(checks, Check{Name: "issuer", Verdict: VerdictFailed, Message: fmt.Sprintf(
			"the provider advertises issuer %q but the manifest declares %q -- token validation is strict about this exact string; use the advertised value",
			advertisedIssuer, broker.IssuerURL)})
	default:
		checks = append(checks, Check{Name: "issuer", Verdict: VerdictPassed, Message: fmt.Sprintf(
			"issuer discovered; sign-in will authorize at %s", endpoints.AuthorizationURL)})
	}

	checks = append(checks, Check{Name: "clientAuthentication", Verdict: VerdictUnknown, Message: "the client id and secret are verifiable only at a real sign-in (the upstream provider offers no credential probe); the first sign-in attempt is the proof"})
	return checks
}

// entraTenancyCheck names the Entra single-tenant constraint when the issuer
// is a Microsoft multiplexer endpoint -- /common and /organizations cannot
// satisfy strict issuer validation, and the failure they cause downstream
// names neither.
func entraTenancyCheck(issuerURL string) *Check {
	lower := strings.ToLower(issuerURL)
	if !strings.Contains(lower, "login.microsoftonline.com") {
		return nil
	}
	if strings.Contains(lower, "/common/") || strings.HasSuffix(lower, "/common") ||
		strings.Contains(lower, "/organizations/") || strings.Contains(lower, "/consumers/") {
		return &Check{Name: "issuerTenancy", Verdict: VerdictFailed, Message: "the issuer uses Entra's multi-tenant endpoint (/common, /organizations, or /consumers), which cannot satisfy strict issuer validation -- use the single-tenant form: https://login.microsoftonline.com/{tenant-id}/v2.0"}
	}
	return &Check{Name: "issuerTenancy", Verdict: VerdictPassed, Message: "the issuer is tenant-scoped"}
}

// verifySeededAdminCollision reports (never resolves) the seeded-admin email
// collision: the platform manifest's adminEmail seeds a LOCAL realm user, and
// a directory user with the same email can then never materialize (realm
// emails are unique). Both objects are legitimate; the remedy is a human
// choice, named in the message.
func verifySeededAdminCollision(ctx context.Context, admin *AdminClient, realm, seededEmail string) Check {
	users, err := admin.FindUsersByEmail(ctx, realm, seededEmail)
	if err != nil {
		return Check{Name: "seededAdminCollision", Verdict: VerdictUnknown, Message: "could not search realm users: " + err.Error()}
	}
	for _, user := range users {
		if link, _ := user["federationLink"].(string); link == "" {
			// The seeded local admin exists alongside declared federation.
			return Check{Name: "seededAdminCollision", Verdict: VerdictFailed, Message: fmt.Sprintf(
				"a local admin user holds %s, so a directory account with the same email can never sign in (realm emails are unique). Remedy: list the person under bootstrap.admins WITHOUT adminEmail in the PlantonPlatform manifest (bootstrap grants match by email at first federated sign-in), then delete the local user in the identity server's admin console once the federated sign-in works",
				seededEmail)}
		}
	}
	return Check{Name: "seededAdminCollision", Verdict: VerdictPassed, Message: fmt.Sprintf(
		"no local user holds %s; the declared admin can arrive through the directory", seededEmail)}
}

// DiscoverOIDC fetches the issuer's discovery document and returns the
// endpoints the broker config needs plus the ADVERTISED issuer for strict
// comparison. Exported because the identity component runs discovery on the
// verification cadence and passes the result into convergence.
func DiscoverOIDC(ctx context.Context, httpClient *http.Client, issuerURL string) (*OIDCEndpoints, string, error) {
	wellKnown := strings.TrimSuffix(issuerURL, "/") + "/.well-known/openid-configuration"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, wellKnown, nil)
	if err != nil {
		return nil, "", fmt.Errorf("building discovery request: %w", err)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("discovery returned %d", resp.StatusCode)
	}

	var doc struct {
		Issuer                string `json:"issuer"`
		AuthorizationEndpoint string `json:"authorization_endpoint"`
		TokenEndpoint         string `json:"token_endpoint"`
		JWKSURI               string `json:"jwks_uri"`
		UserinfoEndpoint      string `json:"userinfo_endpoint"`
		EndSessionEndpoint    string `json:"end_session_endpoint"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return nil, "", fmt.Errorf("decoding discovery document: %w", err)
	}
	if doc.AuthorizationEndpoint == "" || doc.TokenEndpoint == "" || doc.JWKSURI == "" {
		return nil, "", errors.New("discovery document is missing required endpoints (authorization_endpoint, token_endpoint, jwks_uri)")
	}
	return &OIDCEndpoints{
		AuthorizationURL: doc.AuthorizationEndpoint,
		TokenURL:         doc.TokenEndpoint,
		JWKSURL:          doc.JWKSURI,
		UserInfoURL:      doc.UserinfoEndpoint,
		LogoutURL:        doc.EndSessionEndpoint,
	}, doc.Issuer, nil
}

// keycloakReason extracts the identity server's own explanation from an admin
// API error -- Keycloak wraps it as {"errorMessage": "..."} -- falling back
// to the raw error text.
func keycloakReason(err error) string {
	var statusErr *unexpectedStatusError
	if errors.As(err, &statusErr) {
		var body struct {
			ErrorMessage string `json:"errorMessage"`
			Error        string `json:"error"`
		}
		if jsonErr := json.Unmarshal([]byte(statusErr.ResponseBody()), &body); jsonErr == nil {
			if body.ErrorMessage != "" {
				return body.ErrorMessage
			}
			if body.Error != "" {
				return body.Error
			}
		}
		if statusErr.ResponseBody() != "" {
			return statusErr.ResponseBody()
		}
	}
	return err.Error()
}

// cloneRepresentation shallow-copies a representation so probe variants never
// share state.
func cloneRepresentation(rep Representation) Representation {
	return maps.Clone(rep)
}
