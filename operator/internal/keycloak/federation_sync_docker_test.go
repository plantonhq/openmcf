//go:build requires_docker

package keycloak

// The directory sync engine's read premises, probed against the REAL lab
// directory with the EXACT credential the control plane holds (the
// planton-users service account). The periodic sync's whole design stands on
// what these probes prove: whether a mirrored group's member list reflects
// directory truth, on what cadence, and what Keycloak does with a person the
// directory no longer has. Every directory mutation below RESTORES the
// seeded state before returning -- the lab is a suite-wide singleton and
// sibling tests assert seeded facts (the 16-user count, platform-eng's GUID
// cross-check).

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/plantonhq/planton/operator/internal/resources"
)

// memberUsernames lists a mirrored group's members through the admin API as
// the given client, keyed by username -> the member's representation.
func memberUsernames(t *testing.T, client *AdminClient, realm, groupID string) map[string]Representation {
	t.Helper()
	var members []Representation
	if err := client.do(context.Background(), http.MethodGet,
		client.adminPath(realm, "/groups/"+groupID+"/members?first=0&max=200"),
		nil, http.StatusOK, &members); err != nil {
		t.Fatalf("reading group %s members: %v", groupID, err)
	}
	byName := map[string]Representation{}
	for _, m := range members {
		if username, _ := m["username"].(string); username != "" {
			byName[username] = m
		}
	}
	return byName
}

// directoryChildren lists the /directory mirror's children (name -> group
// representation) as the given client.
func directoryChildren(t *testing.T, client *AdminClient, realm string) map[string]Representation {
	t.Helper()
	ctx := context.Background()
	parent, found, err := client.GetGroupByPath(ctx, realm, resources.IdentityDirectoryGroupsPath)
	if err != nil || !found {
		t.Fatalf("reading the %s parent (found=%v err=%v)", resources.IdentityDirectoryGroupsPath, found, err)
	}
	parentID, _ := parent["id"].(string)
	var children []Representation
	if err := client.do(ctx, http.MethodGet,
		client.adminPath(realm, "/groups/"+parentID+"/children?briefRepresentation=false&max=100"),
		nil, http.StatusOK, &children); err != nil {
		t.Fatalf("listing the mirror's children: %v", err)
	}
	byName := map[string]Representation{}
	for _, child := range children {
		if name, _ := child["name"].(string); name != "" {
			byName[name] = child
		}
	}
	return byName
}

// triggerFullUserSync runs the federation component's full user sync -- the
// same sync Keycloak's own fullSyncPeriod timer runs on a live install (the
// test triggers it directly because waiting out a timer proves nothing more).
func triggerFullUserSync(t *testing.T, admin *AdminClient, realm string) *SyncResult {
	t.Helper()
	ctx := context.Background()
	realmRep, err := admin.GetRealm(ctx, realm)
	if err != nil {
		t.Fatal(err)
	}
	realmID, _ := realmRep["id"].(string)
	components, err := admin.ListComponents(ctx, realm, realmID, userStorageProviderType)
	if err != nil {
		t.Fatal(err)
	}
	var componentID string
	for _, c := range components {
		if c["name"] == resources.IdentityLDAPComponentName {
			componentID, _ = c["id"].(string)
		}
	}
	if componentID == "" {
		t.Fatal("no LDAP federation component to sync")
	}
	result, err := admin.TriggerUserStorageSync(ctx, realm, componentID, "triggerFullSync")
	if err != nil {
		t.Fatalf("triggering the full user sync: %v", err)
	}
	return result
}

// triggerGroupMapperSync runs the group mapper's fedToKeycloak sync -- the
// pass that refreshes the /directory mirror's group OBJECTS.
func triggerGroupMapperSync(t *testing.T, admin *AdminClient, realm string) {
	t.Helper()
	ctx := context.Background()
	realmRep, err := admin.GetRealm(ctx, realm)
	if err != nil {
		t.Fatal(err)
	}
	realmID, _ := realmRep["id"].(string)
	components, err := admin.ListComponents(ctx, realm, realmID, userStorageProviderType)
	if err != nil {
		t.Fatal(err)
	}
	var componentID string
	for _, c := range components {
		if c["name"] == resources.IdentityLDAPComponentName {
			componentID, _ = c["id"].(string)
		}
	}
	mappers, err := admin.ListComponents(ctx, realm, componentID, ldapStorageMapperType)
	if err != nil {
		t.Fatal(err)
	}
	var mapperID string
	for _, m := range mappers {
		if m["name"] == ldapGroupMapperName {
			mapperID, _ = m["id"].(string)
		}
	}
	if mapperID == "" {
		t.Fatal("no group mapper to sync")
	}
	if _, err := admin.TriggerLDAPMapperSync(ctx, realm, componentID, mapperID, "fedToKeycloak"); err != nil {
		t.Fatalf("triggering the group mapper sync: %v", err)
	}
}

// The sync engine's four pre-registered probes, one journey:
//
//	A -- membership truth: the planton-users credential reads a mirrored
//	     group's members, and the list reflects a membership change in the
//	     directory (recording WHETHER it needed a user sync first -- the
//	     finding that decides the offboarding guarantee's phrasing).
//	B -- vanished user: what Keycloak does with an imported person the
//	     directory deleted (removal on full sync, or a lingering local row).
//	C -- disabled user: how a disabled directory account appears in a
//	     members read (the enabled flag's visibility).
//	D -- mirror freshness: whether a group CREATED after provisioning
//	     reaches the /directory mirror on the periodic user sync alone, or
//	     only on the verification-time group mapper sync.
func TestFederation_DirectorySyncReadPremises(t *testing.T) {
	admin := authedAdmin(t)
	createRealm(t, admin, "syncread", nil)
	fed := &OwnedFederation{LDAP: testLab.ldapFederation()}
	fed.LDAP.RotateCredential = true
	mustConverge(t, federationInput("syncread", fed))

	// Verification populates the mirror and imports the users -- the same
	// pass a live install runs at provisioning.
	ctx := context.Background()
	if _, err := Verify(ctx, verifyInput("syncread", fed)); err != nil {
		t.Fatalf("verification could not run: %v", err)
	}

	// The control plane's own credential -- every read below must work with
	// planton-users' two roles (manage-users, view-users), never the
	// bootstrap admin.
	users := NewAdminClient(&http.Client{Timeout: 10 * time.Second}, testServerRoot)
	users.token = clientCredentialsToken(t, "syncread", resources.IdentityUsersClientID, "users-secret-value")

	mirror := directoryChildren(t, users, "syncread")
	platformEng, ok := mirror["platform-eng"]
	if !ok {
		t.Fatalf("platform-eng missing from the mirror: %v", mirror)
	}
	platformEngID, _ := platformEng["id"].(string)
	developers, ok := mirror["developers"]
	if !ok {
		t.Fatalf("developers missing from the mirror: %v", mirror)
	}
	developersID, _ := developers["id"].(string)

	// ---- Probe A: membership truth through the members read -------------

	members := memberUsernames(t, users, "syncread", platformEngID)
	if len(members) == 0 {
		t.Fatal("PAUSE CONDITION: the members read answers nothing -- the group-first sync has no read path")
	}
	ada, ok := members["ada.lovelace"]
	if !ok {
		t.Fatalf("ada.lovelace missing from platform-eng's members: %v", memberKeys(members))
	}
	t.Logf("FINDING A0: platform-eng members read answers %d members with planton-users", len(members))

	// The fallback read: one user's groups.
	adaID, _ := ada["id"].(string)
	var adaGroups []Representation
	if err := users.do(ctx, http.MethodGet,
		users.adminPath("syncread", "/users/"+adaID+"/groups?max=100"),
		nil, http.StatusOK, &adaGroups); err != nil {
		t.Fatalf("reading ada.lovelace's groups (the fallback read): %v", err)
	}
	foundPlatformEng := false
	for _, g := range adaGroups {
		if g["name"] == "platform-eng" {
			foundPlatformEng = true
		}
	}
	if !foundPlatformEng {
		t.Errorf("the per-user groups read must include platform-eng: %v", adaGroups)
	}

	// A membership change in the DIRECTORY: alan.turing leaves developers.
	// Read immediately (does the members endpoint query LDAP live?), then
	// after a full user sync (the periodic timer's pass). Restore either way.
	testLab.removeGroupMember(t, "developers", "alan.turing")
	defer testLab.addGroupMember(t, "developers", "alan.turing")

	// PINNED (suite-discovered): the members read reflects the directory
	// LIVE -- Keycloak's READ_ONLY group mapper queries LDAP at read time,
	// no user sync needed. This is what makes the product's sync cadence
	// the WHOLE offboarding window ("within 15 minutes", not "your
	// directory sync window plus 15"). If a Keycloak upgrade ever changes
	// this, the guarantee's phrasing breaks -- this assertion lands red
	// before the docs can lie.
	liveMembers := memberUsernames(t, users, "syncread", developersID)
	if _, still := liveMembers["alan.turing"]; still {
		t.Error("FINDING A1 drifted: the members read no longer reflects the directory live -- " +
			"the offboarding guarantee's phrasing (product cadence only) no longer holds")
	} else {
		t.Log("FINDING A1: the members read reflects the directory LIVE (no sync needed)")
	}

	syncResult := triggerFullUserSync(t, admin, "syncread")
	t.Logf("FINDING A2: full user sync result: added=%d updated=%d removed=%d failed=%d",
		syncResult.Added, syncResult.Updated, syncResult.Removed, syncResult.Failed)
	afterSync := memberUsernames(t, users, "syncread", developersID)
	if _, still := afterSync["alan.turing"]; still {
		t.Fatal("PAUSE CONDITION: a directory membership removal is invisible to the members read even after a full user sync -- no membership truth exists")
	}
	if _, adaStill := afterSync["ada.lovelace"]; !adaStill {
		t.Error("ada.lovelace must survive the sync in developers (only alan.turing left)")
	}
	t.Log("FINDING A3: after the full user sync the members read reflects the removal -- membership truth holds")

	// ---- Probe B: the vanished user --------------------------------------

	testLab.deleteUser(t, "lee.lonely")
	defer func() {
		testLab.createUser(t, "lee.lonely", "lee.lonely@lab.example.internal")
		testLab.addGroupMember(t, "everyone", "lee.lonely")
	}()

	syncResult = triggerFullUserSync(t, admin, "syncread")
	t.Logf("FINDING B1: full sync after a directory deletion: added=%d updated=%d removed=%d failed=%d",
		syncResult.Added, syncResult.Updated, syncResult.Removed, syncResult.Failed)

	var lonely []Representation
	if err := admin.do(ctx, http.MethodGet,
		admin.adminPath("syncread", "/users?exact=true&username=lee.lonely"),
		nil, http.StatusOK, &lonely); err != nil {
		t.Fatal(err)
	}
	// PINNED (suite-discovered): Keycloak REMOVES the imported local row
	// when the directory no longer has the person -- a vanished user leaves
	// no ghost the members read could still name.
	if len(lonely) != 0 {
		enabled, _ := lonely[0]["enabled"].(bool)
		t.Errorf("FINDING B2 drifted: the imported user survives the full sync (enabled=%v) -- "+
			"offboarding can no longer rely on Keycloak's vanished-user removal", enabled)
	} else {
		t.Log("FINDING B2: Keycloak REMOVED the imported user on the full sync (vanished directory user leaves no local row)")
	}
	everyoneID, _ := mirror["everyone"]["id"].(string)
	everyoneMembers := memberUsernames(t, users, "syncread", everyoneID)
	if _, still := everyoneMembers["lee.lonely"]; still {
		t.Error("a deleted directory user must vanish from the members read after a full sync -- the offboarding premise")
	} else {
		t.Log("FINDING B3: the deleted user is gone from the members read -- offboarding-as-drift-repair holds")
	}

	// ---- Probe C: the disabled user ---------------------------------------

	// PINNED (suite-discovered): a disabled directory account stays in the
	// members read carrying enabled=false -- the reconciler can SEE the
	// disabled state and deliberately keeps the membership (a disabled
	// person cannot sign in and their session dies with its 15-minute
	// token; group membership is the directory's fact, and stripping it
	// would misreport what the directory says).
	platformEngMembers := memberUsernames(t, users, "syncread", platformEngID)
	dara, present := platformEngMembers["dara.disabled"]
	if !present {
		t.Error("FINDING C drifted: the disabled fixture vanished from the members read -- " +
			"the keep-disabled-memberships law has nothing to key on")
	} else if enabled, _ := dara["enabled"].(bool); enabled {
		t.Error("FINDING C drifted: the disabled fixture reads enabled=true -- the flag is no longer visible")
	} else {
		t.Log("FINDING C: the disabled fixture appears in the members read (enabled=false)")
	}

	// ---- Probe D: mirror freshness for a NEW group ------------------------

	testLab.addGroup(t, "sync-probe-fresh")
	defer testLab.deleteGroup(t, "sync-probe-fresh")
	testLab.addGroupMember(t, "sync-probe-fresh", "grace.hopper")

	// PINNED (suite-discovered): the periodic USER sync alone brings a new
	// group into the mirror -- Keycloak's own fullSyncPeriod timer keeps
	// the group catalog fresh at the manifest's stated window, and no
	// operator- or product-side mapper-sync scheduling is needed. If this
	// drifts, groups created after provisioning become invisible to
	// mappings until a manifest edit -- a freshness hole the platform must
	// then close deliberately.
	triggerFullUserSync(t, admin, "syncread")
	afterUserSync := directoryChildren(t, users, "syncread")
	if _, present := afterUserSync["sync-probe-fresh"]; present {
		t.Log("FINDING D1: the periodic USER sync alone brings a new group into the mirror")
	} else {
		t.Error("FINDING D1 drifted: the periodic user sync no longer refreshes the mirror's group objects -- " +
			"new directory groups are invisible until a manifest edit; diagnostic below names the fallback")
		triggerGroupMapperSync(t, admin, "syncread")
		afterMapperSync := directoryChildren(t, users, "syncread")
		if _, present := afterMapperSync["sync-probe-fresh"]; !present {
			t.Error("a new directory group must reach the mirror at least on the group mapper sync")
		} else {
			t.Log("DIAGNOSTIC: the group mapper sync still brings the new group in -- align a mapper sync to the manifest cadence")
		}
	}
}

func memberKeys(members map[string]Representation) []string {
	keys := make([]string, 0, len(members))
	for k := range members {
		keys = append(keys, k)
	}
	return keys
}
