package infrachart

import (
	"io"
	"testing"

	"github.com/nikolalohinski/gonja/v2/logging"
	"github.com/sirupsen/logrus"
)

// Linking this package must leave the host process's logger alone. A
// binary that imports it (every Planton CLI does) relies on logrus for its
// own messages, and silencing the standard logger here once made every
// fatal exit in those binaries print nothing at all.
func TestPackageInitLeavesTheProcessLoggerAlone(t *testing.T) {
	if logrus.StandardLogger().Out == io.Discard {
		t.Fatal("the standard logrus output is io.Discard after importing infrachart; gonja must be silenced through its own logging switch, never through logrus.SetOutput")
	}
	if logging.Enabled() {
		t.Fatal("gonja logging is enabled; the render library must stay silent in the host process")
	}
}
