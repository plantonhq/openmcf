// Package capacityestimator evaluates a component's capacity derivation
// against one typed manifest. It is the execution engine of the
// capacity-derivation standard, the cluster-capacity twin of the cost
// estimator: workload bindings locate the manifest's ContainerResources
// blocks, instance counts, and per-instance volumes, and the engine sums
// them into the capacity footprint the workload reserves from its target
// cluster -- exact Kubernetes-quantity arithmetic, never a float. A
// manifest that states no reservation at all refuses honestly instead of
// emitting an empty footprint.
//
// The result is the same shape an authored estimate model preset carries
// (a capacity_footprint with configuration-scoped exclusions and notes),
// so the estimate generator's emission runs unchanged on evaluated
// output. Condition semantics for the conditional prose live in the cost
// estimator (one home): effective values, placeholder-as-unset, and
// reference presence behave identically across both standards.
package capacityestimator

import (
	"fmt"
	"math/big"
	"strings"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"

	capacityv1 "github.com/plantonhq/planton/finops/componentcapacityderivation/v1"
	costestimatev1 "github.com/plantonhq/planton/finops/componentcostestimate/v1"
	estimatemodelv1 "github.com/plantonhq/planton/finops/componentcostestimatemodel/v1"
	"github.com/plantonhq/planton/pkg/finops/costestimator"
	"github.com/plantonhq/planton/pkg/specpath"
)

// Evaluate runs a capacity derivation against one manifest and returns
// exactly one of: the evaluated preset-model shape (a capacity footprint
// with configuration-scoped exclusions and notes), a refusal, or an error
// for a malformed derivation (unresolvable path, unparseable quantity --
// the conformance gate catches path problems before CI ever runs the
// engine, so errors here mean the gate was bypassed or a preset carries a
// quantity outside the grammar).
func Evaluate(
	manifest proto.Message,
	spec *capacityv1.ComponentCapacityDerivationSpec,
) (*estimatemodelv1.PresetEstimateModel, *costestimator.Refusal, error) {
	specMsg, err := costestimator.ManifestSpec(manifest)
	if err != nil {
		return nil, nil, err
	}

	totals := &footprintTotals{}
	var segments []string
	for _, workload := range spec.GetWorkloads() {
		segment, err := accumulateWorkload(specMsg, workload, totals)
		if err != nil {
			return nil, nil, err
		}
		if segment != "" {
			segments = append(segments, segment)
		}
	}

	if totals.empty() {
		return nil, &costestimator.Refusal{Reason: "the manifest declares no resource requests, limits, or volumes -- the reservation is unknowable from this manifest"}, nil
	}

	exclusions, err := costestimator.MatchingTexts(specMsg, spec.GetExclusions())
	if err != nil {
		return nil, nil, err
	}
	notes, err := costestimator.MatchingTexts(specMsg, spec.GetNotes())
	if err != nil {
		return nil, nil, err
	}

	footprint, err := totals.render()
	if err != nil {
		return nil, nil, err
	}
	footprint.Basis = strings.Join(segments, "; ")

	return &estimatemodelv1.PresetEstimateModel{
		CapacityFootprint: footprint,
		Exclusions:        exclusions,
		Notes:             strings.Join(notes, " "),
	}, nil, nil
}

// footprintTotals accumulates the five capacity buckets exactly: CPU in
// millicores, memory and storage in bytes -- big.Int arithmetic, so no
// quantity ever rounds.
type footprintTotals struct {
	cpuRequests    big.Int
	memoryRequests big.Int
	cpuLimits      big.Int
	memoryLimits   big.Int
	storage        big.Int
}

func (t *footprintTotals) empty() bool {
	return t.cpuRequests.Sign() == 0 && t.memoryRequests.Sign() == 0 &&
		t.cpuLimits.Sign() == 0 && t.memoryLimits.Sign() == 0 && t.storage.Sign() == 0
}

func (t *footprintTotals) render() (*costestimatev1.CapacityFootprint, error) {
	return &costestimatev1.CapacityFootprint{
		CpuRequests:       renderCpu(&t.cpuRequests),
		MemoryRequests:    renderBytes(&t.memoryRequests),
		CpuLimits:         renderCpu(&t.cpuLimits),
		MemoryLimits:      renderBytes(&t.memoryLimits),
		PersistentStorage: renderBytes(&t.storage),
	}, nil
}

// accumulateWorkload adds one workload's reservations into the totals and
// returns its basis segment ("" when the workload contributes nothing --
// zero instances, or neither resources nor volumes resolved).
func accumulateWorkload(specMsg protoreflect.Message, workload *capacityv1.WorkloadBinding, totals *footprintTotals) (string, error) {
	count, err := instanceCount(specMsg, workload)
	if err != nil {
		return "", err
	}
	if count.Sign() == 0 {
		return "", nil
	}

	resources, resourcesBasis, err := workloadResources(specMsg, workload, count, totals)
	if err != nil {
		return "", err
	}

	volumesBasis, contributed, err := workloadVolumes(specMsg, workload, count, totals)
	if err != nil {
		return "", err
	}

	if !resources && !contributed {
		return "", nil
	}

	noun := workload.GetLabel()
	if count.Cmp(big.NewInt(1)) != 0 {
		noun += "s"
	}
	segment := ""
	if resources {
		segment = fmt.Sprintf("%s %s x %s", count.String(), noun, resourcesBasis)
	}
	if contributed {
		if segment != "" {
			segment += " + " + volumesBasis
		} else {
			segment = fmt.Sprintf("%s %s: %s", count.String(), noun, volumesBasis)
		}
	}
	return segment, nil
}

// instanceCount resolves a workload's multiplier: a literal constant, or
// a numeric spec field with its declared default applying to absent AND
// explicitly-zero values (the FieldValue contract the cost derivations
// share).
func instanceCount(specMsg protoreflect.Message, workload *capacityv1.WorkloadBinding) (*big.Int, error) {
	switch count := workload.GetInstances().GetCount().(type) {
	case *capacityv1.InstanceCount_Constant:
		value, ok := new(big.Int).SetString(count.Constant, 10)
		if !ok {
			return nil, fmt.Errorf("workload %q: instance constant %q is not a whole number", workload.GetLabel(), count.Constant)
		}
		return value, nil
	case *capacityv1.InstanceCount_FieldValue:
		resolved, err := specpath.Resolve(specMsg, count.FieldValue.GetFieldPath())
		if err != nil {
			return nil, fmt.Errorf("workload %q instances: %w", workload.GetLabel(), err)
		}
		value := big.NewInt(0)
		if resolved.Value.IsValid() {
			switch resolved.Field.Kind() {
			case protoreflect.Int32Kind, protoreflect.Sint32Kind, protoreflect.Sfixed32Kind,
				protoreflect.Int64Kind, protoreflect.Sint64Kind, protoreflect.Sfixed64Kind:
				value = big.NewInt(resolved.Value.Int())
			case protoreflect.Uint32Kind, protoreflect.Fixed32Kind,
				protoreflect.Uint64Kind, protoreflect.Fixed64Kind:
				value = new(big.Int).SetUint64(resolved.Value.Uint())
			default:
				return nil, fmt.Errorf("workload %q instances: field %q is not an integer", workload.GetLabel(), count.FieldValue.GetFieldPath())
			}
		}
		if value.Sign() == 0 && count.FieldValue.GetDefaultWhenUnset() != "" {
			defaulted, ok := new(big.Int).SetString(count.FieldValue.GetDefaultWhenUnset(), 10)
			if !ok {
				return nil, fmt.Errorf("workload %q: instance default %q is not a whole number", workload.GetLabel(), count.FieldValue.GetDefaultWhenUnset())
			}
			value = defaulted
		}
		return value, nil
	default:
		return nil, fmt.Errorf("workload %q declares no instance count", workload.GetLabel())
	}
}

// workloadResources reads the workload's ContainerResources block (when
// bound and present), adds count x each quantity into the totals, and
// renders the per-instance basis fragment, e.g.
// "(1 CPU / 2Gi memory requests, 2 CPU / 4Gi limits)".
func workloadResources(specMsg protoreflect.Message, workload *capacityv1.WorkloadBinding, count *big.Int, totals *footprintTotals) (bool, string, error) {
	path := workload.GetResourcesPath()
	if path == "" {
		return false, "", nil
	}
	resolved, err := specpath.Resolve(specMsg, path)
	if err != nil {
		return false, "", fmt.Errorf("workload %q resources: %w", workload.GetLabel(), err)
	}
	if !resolved.Present || !resolved.Value.IsValid() {
		return false, "", nil
	}
	resources := resolved.Value.Message()

	requestsCpu, requestsMemory := cpuMemory(resources, "requests")
	limitsCpu, limitsMemory := cpuMemory(resources, "limits")

	contributed := false
	add := func(total *big.Int, quantity string, milli bool) error {
		if quantity == "" {
			return nil
		}
		value, err := parseQuantity(quantity, milli)
		if err != nil {
			return fmt.Errorf("workload %q resources: %w", workload.GetLabel(), err)
		}
		total.Add(total, value.Mul(value, count))
		contributed = true
		return nil
	}
	if err := add(&totals.cpuRequests, requestsCpu, true); err != nil {
		return false, "", err
	}
	if err := add(&totals.memoryRequests, requestsMemory, false); err != nil {
		return false, "", err
	}
	if err := add(&totals.cpuLimits, limitsCpu, true); err != nil {
		return false, "", err
	}
	if err := add(&totals.memoryLimits, limitsMemory, false); err != nil {
		return false, "", err
	}
	if !contributed {
		return false, "", nil
	}

	var parts []string
	if fragment := cpuMemoryFragment(requestsCpu, requestsMemory, "memory "); fragment != "" {
		parts = append(parts, fragment+"requests")
	}
	if fragment := cpuMemoryFragment(limitsCpu, limitsMemory, ""); fragment != "" {
		parts = append(parts, fragment+"limits")
	}
	return true, "(" + strings.Join(parts, ", ") + ")", nil
}

// cpuMemory reads one CpuMemory block ("requests" or "limits") off a
// ContainerResources message, returning empty strings for absent halves.
func cpuMemory(resources protoreflect.Message, block string) (string, string) {
	field := resources.Descriptor().Fields().ByName(protoreflect.Name(block))
	if field == nil || !resources.Has(field) {
		return "", ""
	}
	pair := resources.Get(field).Message()
	read := func(name string) string {
		inner := pair.Descriptor().Fields().ByName(protoreflect.Name(name))
		if inner == nil {
			return ""
		}
		return pair.Get(inner).String()
	}
	return read("cpu"), read("memory")
}

// cpuMemoryFragment renders one block's per-instance basis half, e.g.
// "1 CPU / 2Gi memory " (the caller appends "requests"/"limits").
func cpuMemoryFragment(cpu, memory, memoryLabel string) string {
	var parts []string
	if cpu != "" {
		parts = append(parts, cpu+" CPU")
	}
	if memory != "" {
		parts = append(parts, memory+" "+memoryLabel)
	}
	if len(parts) == 0 {
		return ""
	}
	joined := strings.Join(parts, " / ")
	if !strings.HasSuffix(joined, " ") {
		joined += " "
	}
	return joined
}

// workloadVolumes adds count x each resolved volume size into the storage
// total and renders the volumes basis fragment, e.g.
// "3 x (100Gi data + 20Gi WAL) volumes".
func workloadVolumes(specMsg protoreflect.Message, workload *capacityv1.WorkloadBinding, count *big.Int, totals *footprintTotals) (string, bool, error) {
	type resolvedVolume struct{ size, label string }
	var volumes []resolvedVolume
	for _, volume := range workload.GetVolumes() {
		resolved, err := specpath.Resolve(specMsg, volume.GetSizePath())
		if err != nil {
			return "", false, fmt.Errorf("workload %q volume %q: %w", workload.GetLabel(), volume.GetLabel(), err)
		}
		if !resolved.Present || !resolved.Value.IsValid() {
			continue
		}
		size := resolved.Value.String()
		if size == "" || isPlaceholder(size) {
			continue
		}
		bytes, err := parseQuantity(size, false)
		if err != nil {
			return "", false, fmt.Errorf("workload %q volume %q: %w", workload.GetLabel(), volume.GetLabel(), err)
		}
		totals.storage.Add(&totals.storage, bytes.Mul(bytes, count))
		volumes = append(volumes, resolvedVolume{size: size, label: volume.GetLabel()})
	}
	if len(volumes) == 0 {
		return "", false, nil
	}

	if len(volumes) == 1 {
		noun := "volume"
		if count.Cmp(big.NewInt(1)) != 0 {
			noun = "volumes"
		}
		return fmt.Sprintf("%s x %s %s %s", count.String(), volumes[0].size, volumes[0].label, noun), true, nil
	}
	var terms []string
	for _, volume := range volumes {
		terms = append(terms, volume.size+" "+volume.label)
	}
	return fmt.Sprintf("%s x (%s) volumes", count.String(), strings.Join(terms, " + ")), true, nil
}

// isPlaceholder recognizes the catalog's template placeholders ("<size>")
// so they never masquerade as real quantities -- the same rule the cost
// estimator applies.
func isPlaceholder(text string) bool {
	return len(text) > 2 && strings.HasPrefix(text, "<") && strings.HasSuffix(text, ">")
}
