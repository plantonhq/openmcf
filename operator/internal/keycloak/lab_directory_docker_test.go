//go:build requires_docker

package keycloak

// The lab directory's Docker driver for this suite: boots the Samba AD
// container, seeds it from the ONE seed manifest
// (hack/lab-directory/seed.yaml -- the definition every identity
// suite shares; see its README), and extracts the lab CA for Keycloak's
// truststore. The loader lives here rather than in a shared package because
// it is Docker-tagged test infrastructure: gazelle skips tagged files, so
// nothing of this ever reaches Bazel or a shipped binary. Other consumers
// (the sync-engine suites, the e2e lab) bring their own loaders to the same
// manifest -- the manifest is the contract, not this code.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"sigs.k8s.io/yaml"

	"github.com/plantonhq/planton/operator/internal/resources"
)

// labSeedPath reaches the lab directory's seed manifest (hack/lab-directory,
// the single definition of the lab population) from this package.
const labSeedPath = "../../hack/lab-directory/seed.yaml"

const labImage = "smblds/smblds:latest"

// labSeed mirrors seed.yaml's shape.
type labSeed struct {
	Realm         string `json:"realm"`
	NetbiosDomain string `json:"netbiosDomain"`
	BaseDN        string `json:"baseDn"`
	AdminPassword string `json:"adminPassword"`
	UsersOu       string `json:"usersOu"`
	GroupsOu      string `json:"groupsOu"`
	UserPassword  string `json:"userPassword"`
	Users         []struct {
		Username  string `json:"username"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Email     string `json:"email"`
		Disabled  bool   `json:"disabled"`
	} `json:"users"`
	Groups []struct {
		Name         string   `json:"name"`
		Members      []string `json:"members"`
		MemberGroups []string `json:"memberGroups"`
	} `json:"groups"`
}

// labDirectory is the running lab AD container plus the facts a federation
// spec needs to point at it.
type labDirectory struct {
	containerID string
	caDir       string
	seed        *labSeed

	// host is the container's name on the suite network -- and the name the
	// Samba-minted TLS certificate is issued for, so LDAPS hostname
	// verification passes exactly when the URL uses it.
	host string
}

// startLabDirectory boots and seeds the lab AD on the suite's network.
func startLabDirectory(network, containerName string) (*labDirectory, error) {
	raw, err := os.ReadFile(labSeedPath)
	if err != nil {
		return nil, fmt.Errorf("reading the lab seed manifest (run from the keycloak package directory): %w", err)
	}
	seed := &labSeed{}
	if err := yaml.Unmarshal(raw, seed); err != nil {
		return nil, fmt.Errorf("parsing the lab seed manifest: %w", err)
	}

	// The cert's DNS name is {hostname}.{realm-domain}; the network alias
	// matches it so LDAPS hostname verification holds.
	shortHost := "lab-dc"
	fqdn := shortHost + "." + strings.ToLower(seed.Realm)

	out, err := exec.Command("docker", "run", "-d", "--rm",
		"--name", containerName,
		"--hostname", shortHost,
		"--network", network, "--network-alias", fqdn,
		"-e", "REALM="+seed.Realm,
		"-e", "DOMAIN="+seed.NetbiosDomain,
		"-e", "ADMINPASS="+seed.AdminPassword,
		labImage).Output()
	if err != nil {
		return nil, fmt.Errorf("starting the lab directory container: %w", err)
	}
	lab := &labDirectory{containerID: strings.TrimSpace(string(out)), seed: seed, host: fqdn}

	if err := lab.waitReady(90 * time.Second); err != nil {
		lab.stop()
		return nil, err
	}
	if err := lab.applySeed(); err != nil {
		lab.stop()
		return nil, err
	}
	if err := lab.extractCA(); err != nil {
		lab.stop()
		return nil, err
	}
	return lab, nil
}

func (l *labDirectory) waitReady(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		err := exec.Command("docker", "exec", l.containerID,
			"ldbsearch", "-H", "/var/lib/samba/private/sam.ldb",
			"-b", l.seed.BaseDN, "-s", "base", "dn").Run()
		if err == nil {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	logs, _ := exec.Command("docker", "logs", "--tail", "30", l.containerID).CombinedOutput()
	return fmt.Errorf("the lab directory never became ready; container logs:\n%s", logs)
}

// applySeed renders the whole population as ONE shell script and applies it
// in a single exec -- per-entry execs would take minutes. Generated bulk
// users (seed.yaml's `generated` block) are deliberately not seeded here:
// this suite proves federation mechanics, not blast-radius realism.
func (l *labDirectory) applySeed() error {
	var script strings.Builder
	script.WriteString("set -e\n")
	fmt.Fprintf(&script, "samba-tool ou add %q\n", l.seed.UsersOu)
	fmt.Fprintf(&script, "samba-tool ou add %q\n", l.seed.GroupsOu)

	// Entries are created with CN = username (a common real-AD convention),
	// with the person attributes set by a follow-up modify: samba-tool
	// derives the CN from "given surname" when those flags are passed, and
	// two people SHARING a display name (the alex.kim fixture) would then
	// collide at the RDN -- a genuine AD lesson this suite caught on its
	// first run.
	for _, user := range l.seed.Users {
		fmt.Fprintf(&script, "samba-tool user create %q %q --userou=%q",
			user.Username, l.seed.UserPassword, l.seed.UsersOu)
		if user.Email != "" {
			fmt.Fprintf(&script, " --mail-address=%q", user.Email)
		}
		script.WriteString("\n")

		userDN := fmt.Sprintf("CN=%s,%s,%s", user.Username, l.seed.UsersOu, l.seed.BaseDN)
		fmt.Fprintf(&script, "cat <<'LDIF' | ldbmodify -H /var/lib/samba/private/sam.ldb\n"+
			"dn: %s\nchangetype: modify\nreplace: givenName\ngivenName: %s\n-\n"+
			"replace: sn\nsn: %s\n-\nreplace: displayName\ndisplayName: %s %s\nLDIF\n",
			userDN, user.FirstName, user.LastName, user.FirstName, user.LastName)

		if user.Disabled {
			fmt.Fprintf(&script, "samba-tool user disable %q\n", user.Username)
		}
	}

	// Groups first, memberships after: a nested member group must exist
	// before it can be added to its parent.
	for _, group := range l.seed.Groups {
		fmt.Fprintf(&script, "samba-tool group add %q --groupou=%q\n", group.Name, l.seed.GroupsOu)
	}
	for _, group := range l.seed.Groups {
		members := append(append([]string{}, group.Members...), group.MemberGroups...)
		if len(members) == 0 {
			continue
		}
		fmt.Fprintf(&script, "samba-tool group addmembers %q %q\n", group.Name, strings.Join(members, ","))
	}

	if out, err := exec.Command("docker", "exec", l.containerID, "sh", "-c", script.String()).CombinedOutput(); err != nil {
		return fmt.Errorf("seeding the lab directory: %v\n%s", err, out)
	}
	return nil
}

// extractCA copies the lab's private CA to a host temp dir whose CONTENTS
// are later docker-cp'd into Keycloak's container at the production
// truststore path -- copied, never bind-mounted, because Docker Desktop's VM
// file sharing silently materializes an unshared host path as an empty
// directory.
func (l *labDirectory) extractCA() error {
	dir, err := os.MkdirTemp("", "lab-directory-ca-")
	if err != nil {
		return err
	}
	// The dir's mode travels with the docker cp: MkdirTemp's 0700 would
	// arrive root-owned and unreadable to the keycloak user, and Keycloak
	// REFUSES BOOT on an unreadable truststore path.
	if err := os.Chmod(dir, 0o755); err != nil {
		return err
	}
	l.caDir = dir
	target := filepath.Join(dir, resources.IdentityCABundleFileName)
	if out, err := exec.Command("docker", "cp",
		l.containerID+":/var/lib/samba/private/tls/ca.pem", target).CombinedOutput(); err != nil {
		return fmt.Errorf("extracting the lab CA: %v\n%s", err, out)
	}
	// World-readable: the Keycloak container reads it as its own user.
	return os.Chmod(target, 0o644)
}

func (l *labDirectory) stop() {
	_ = exec.Command("docker", "rm", "-f", l.containerID).Run()
	if l.caDir != "" {
		_ = os.RemoveAll(l.caDir)
	}
}

// execSamba runs one shell script inside the lab container -- the mutation
// seam the sync-engine probes use to change the directory mid-test (a person
// leaves a group, a person leaves the company, a new group appears). Every
// probe that mutates MUST restore the seeded state before returning: the lab
// is a suite-wide singleton and later tests assert seeded facts (the 16-user
// count, platform-eng's GUID cross-check).
func (l *labDirectory) execSamba(t *testing.T, script string) {
	t.Helper()
	if out, err := exec.Command("docker", "exec", l.containerID, "sh", "-c", "set -e\n"+script).CombinedOutput(); err != nil {
		t.Fatalf("mutating the lab directory: %v\n%s", err, out)
	}
}

func (l *labDirectory) removeGroupMember(t *testing.T, group, user string) {
	l.execSamba(t, fmt.Sprintf("samba-tool group removemembers %q %q", group, user))
}

func (l *labDirectory) addGroupMember(t *testing.T, group, user string) {
	l.execSamba(t, fmt.Sprintf("samba-tool group addmembers %q %q", group, user))
}

func (l *labDirectory) deleteUser(t *testing.T, username string) {
	l.execSamba(t, fmt.Sprintf("samba-tool user delete %q", username))
}

// createUser recreates a seeded person (restoration after a deletion probe).
// The recreated entry gets a NEW objectGUID -- exactly what a real
// rehire-after-offboarding looks like, and nothing in the suite pins a
// user's GUID.
func (l *labDirectory) createUser(t *testing.T, username, email string) {
	script := fmt.Sprintf("samba-tool user create %q %q --userou=%q", username, l.seed.UserPassword, l.seed.UsersOu)
	if email != "" {
		script += fmt.Sprintf(" --mail-address=%q", email)
	}
	l.execSamba(t, script)
}

func (l *labDirectory) addGroup(t *testing.T, name string) {
	l.execSamba(t, fmt.Sprintf("samba-tool group add %q --groupou=%q", name, l.seed.GroupsOu))
}

func (l *labDirectory) deleteGroup(t *testing.T, name string) {
	l.execSamba(t, fmt.Sprintf("samba-tool group delete %q", name))
}

// groupObjectGUID reads a group's objectGUID from the directory ITSELF
// (Samba's sam.ldb prints it in canonical dashed form) -- the ground-truth
// value the Keycloak mirror's objectGUID attribute is cross-checked against.
func (l *labDirectory) groupObjectGUID(t *testing.T, name string) string {
	t.Helper()
	out, err := exec.Command("docker", "exec", l.containerID,
		"ldbsearch", "-H", "/var/lib/samba/private/sam.ldb",
		"-b", l.seed.GroupsOu+","+l.seed.BaseDN,
		fmt.Sprintf("(cn=%s)", name), "objectGUID").Output()
	if err != nil {
		t.Fatalf("reading %s's objectGUID from the lab directory: %v", name, err)
	}
	for _, line := range strings.Split(string(out), "\n") {
		if guid, found := strings.CutPrefix(line, "objectGUID: "); found {
			return strings.TrimSpace(guid)
		}
	}
	t.Fatalf("no objectGUID in the directory's answer for %s:\n%s", name, out)
	return ""
}

// ldapFederation returns the desired LDAP federation state pointing at the
// lab -- what the identity component would translate from a manifest naming
// this directory.
func (l *labDirectory) ldapFederation() *OwnedLDAPFederation {
	return &OwnedLDAPFederation{
		Servers:            []string{"ldaps://" + l.host + ":636"},
		UseTruststore:      true,
		BindDN:             "CN=Administrator,CN=Users," + l.seed.BaseDN,
		BindCredential:     l.seed.AdminPassword,
		UsersDN:            l.seed.UsersOu + "," + l.seed.BaseDN,
		GroupsDN:           l.seed.GroupsOu + "," + l.seed.BaseDN,
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
