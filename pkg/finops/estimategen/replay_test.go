package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	costestimatev1 "github.com/plantonhq/planton/finops/componentcostestimate/v1"
	"github.com/plantonhq/planton/pkg/finops/costderivation"
	"github.com/plantonhq/planton/pkg/protobufyaml"
)

// TestPresetReplayHonesty pins the honesty semantics of derivation-driven
// estimates -- the behaviors that make a replayed estimate trustworthy and
// that a careless derivation edit could silently lose. Every pin is
// price-refresh-proof: quantities, preset sets, and line presence are
// functions of the presets and the rules, never of the volatile dollar
// values.
//
//   - A preset whose configuration the derivation cannot price (the
//     launch-template node group) is ABSENT -- never a guessed number.
//   - A scale-to-zero Spot pool prices to the literal "0.00" with no
//     line items and its market-rate story in the exclusions.
//   - Same-meter same-price lines merge (the warm pool's floor and
//     pooled root volumes are ONE line whose quantity is the sum).
//   - A read replica prices under STATED assumptions (the basis says
//     "assumed"), and conditional lines exist exactly for the
//     configurations that incur them (no public-IPv4 line on an
//     internal load balancer).
func TestPresetReplayHonesty(t *testing.T) {
	root := repoRoot(t)
	if _, err := os.Stat(filepath.Join(root, "catalog")); err != nil {
		t.Skip("catalog source tree not present (bazel sandbox); runs under go test and the lint.catalog-data lane")
	}
	if derived, err := costderivation.Discover(root); err != nil || len(derived) == 0 {
		t.Skip("no cost derivations authored yet")
	}

	summary, err := Generate(root)
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	if len(summary.Problems) > 0 {
		t.Fatalf("coherence problems: %v", summary.Problems)
	}

	t.Run("awseksnodegroup", func(t *testing.T) {
		presets := generatedPresets(t, summary, "awseksnodegroup")

		if _, exists := presets["03-launch-template"]; exists {
			t.Error("the launch-template preset is priced -- its rate is delegated to the referenced AwsLaunchTemplate and must be refused")
		}

		spot, ok := presets["02-spot-cost-optimized"]
		if !ok {
			t.Fatal("the spot preset is missing -- scale-to-zero must price to an honest 0.00, not disappear")
		}
		if len(spot.GetLineItems()) != 0 || spot.GetTotalListCost() != "0.00" {
			t.Errorf("spot preset: want zero lines and total 0.00, got %d lines and %q",
				len(spot.GetLineItems()), spot.GetTotalListCost())
		}

		warm, ok := presets["04-warm-pool"]
		if !ok {
			t.Fatal("the warm-pool preset is missing")
		}
		volumes := linesForMeter(warm, "node root volumes")
		if len(volumes) != 1 {
			t.Fatalf("warm pool root volumes: want ONE merged line, got %d", len(volumes))
		}
		if volumes[0].GetPricingQuantity() != "400" {
			t.Errorf("warm pool root volumes quantity: got %q, want 400 (2 floor + 2 pooled nodes x 100 GiB)",
				volumes[0].GetPricingQuantity())
		}
	})

	t.Run("awsrdsinstance", func(t *testing.T) {
		presets := generatedPresets(t, summary, "awsrdsinstance")
		replica, ok := presets["03-read-replica"]
		if !ok {
			t.Fatal("the read-replica preset is missing -- it prices under stated assumptions")
		}
		for _, line := range replica.GetLineItems() {
			if line.GetSkuMeter() == "instance hours" && !strings.Contains(line.GetQuantityBasis(), "assumed") {
				t.Errorf("replica instance line basis does not state its engine assumption: %q", line.GetQuantityBasis())
			}
		}
		storage := linesForMeter(replica, "provisioned storage")
		if len(storage) != 1 || storage[0].GetPricingQuantity() != "50" {
			t.Errorf("replica storage: want one line of 50 GB-months (the stated assumption), got %+v", storage)
		}
	})

	t.Run("awsalb", func(t *testing.T) {
		presets := generatedPresets(t, summary, "awsalb")
		if internal, ok := presets["02-internal-hardened"]; !ok {
			t.Fatal("the internal preset is missing")
		} else if lines := linesForMeter(internal, "public IPv4 address usage"); len(lines) != 0 {
			t.Error("an internal load balancer wears a public IPv4 line -- the conditional rule leaked")
		}
		if facing, ok := presets["01-internet-facing"]; !ok {
			t.Fatal("the internet-facing preset is missing")
		} else if lines := linesForMeter(facing, "public IPv4 address usage"); len(lines) != 1 || lines[0].GetPricingQuantity() != "1460" {
			t.Errorf("internet-facing IPv4: want one line of 1460 IP-hours (2 subnets x 730), got %+v", lines)
		}
	})

	t.Run("awsvpc", func(t *testing.T) {
		presets := generatedPresets(t, summary, "awsvpc")
		if len(presets) != 3 {
			t.Fatalf("awsvpc: want all 3 presets priced, got %d", len(presets))
		}
		for key, preset := range presets {
			if len(preset.GetLineItems()) != 0 || preset.GetTotalListCost() != "0.00" {
				t.Errorf("%s: a VPC is free -- want zero lines and total 0.00, got %d lines and %q",
					key, len(preset.GetLineItems()), preset.GetTotalListCost())
			}
			if len(preset.GetExclusions()) == 0 {
				t.Errorf("%s: a zero estimate without exclusions hides where the network's real spend lives", key)
			}
		}
	})
}

// generatedPresets parses one component's generated estimate document out
// of the in-memory summary, keyed by preset stem.
func generatedPresets(t *testing.T, summary *Summary, component string) map[string]*costestimatev1.PresetEstimate {
	t.Helper()
	content, ok := summary.Files[filepath.Join("catalog/_pricing/estimates", component+".yaml")]
	if !ok {
		t.Fatalf("no generated estimate for %s", component)
	}
	estimate := &costestimatev1.ComponentCostEstimate{}
	if err := protobufyaml.LoadYamlBytes([]byte(content), estimate); err != nil {
		t.Fatalf("parsing generated estimate for %s: %v", component, err)
	}
	presets := map[string]*costestimatev1.PresetEstimate{}
	for _, preset := range estimate.GetSpec().GetPresets() {
		presets[preset.GetPreset()] = preset
	}
	return presets
}

// linesForMeter returns a preset's line items carrying one meter.
func linesForMeter(preset *costestimatev1.PresetEstimate, meter string) []*costestimatev1.LineItem {
	var lines []*costestimatev1.LineItem
	for _, line := range preset.GetLineItems() {
		if line.GetSkuMeter() == meter {
			lines = append(lines, line)
		}
	}
	return lines
}
