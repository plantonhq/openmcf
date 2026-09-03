//go:build e2e

// Package chartlifecycle proves, on a real Kind cluster, the lifecycle
// helm/planton-operator promises for the definitions it owns: a fresh install
// creates and owns them and the manager starts against them; an uninstall
// keeps them and a reinstall adopts them; crds.keep=false removes them; and an
// install whose definitions predate this ownership (Helm's install-once crds/
// directory, from the operator chart or from the chart that once bundled the
// operator) is stopped with a message whose commands, run verbatim, complete
// the adoption so the upgrade carries the new schema.
//
// Helm runs through its SDK, the code path the Planton CLI and desktop drive;
// cluster facts are read with kubectl, the way the operator's own e2e suite
// does. The published charts are installed with the operator scaled to zero:
// what is under proof is ownership and schema, not a platform boot.
package chartlifecycle

import (
	"fmt"
	"os/exec"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/plantonhq/planton/operator/test/utils"
)

// managerImage is the working-tree manager image built and loaded into Kind
// for the fresh-install scenario.
const managerImage = "example.com/planton-operator:chart-lifecycle"

func TestChartLifecycle(t *testing.T) {
	RegisterFailHandler(Fail)
	_, _ = fmt.Fprintf(GinkgoWriter, "Starting planton-operator chart lifecycle suite\n")
	RunSpecs(t, "chart lifecycle suite")
}

var _ = BeforeSuite(func() {
	SetDefaultEventuallyTimeout(3 * time.Minute)
	SetDefaultEventuallyPollingInterval(2 * time.Second)

	By("building the manager image from the working tree")
	cmd := exec.Command("make", "docker-build", fmt.Sprintf("IMG=%s", managerImage))
	_, err := utils.Run(cmd)
	Expect(err).NotTo(HaveOccurred(), "Failed to build the manager image")

	By("loading the manager image into Kind")
	Expect(utils.LoadImageToKindClusterWithName(managerImage)).To(Succeed())
})
