package envpartition

import (
	"fmt"
	"sort"
	"strings"
)

// Resource is the engine's neutral input: one discovered account resource,
// reduced to exactly the surfaces the precedence ladder reads. Provider
// adapters (see the awsscan subpackage) build these; the engine itself
// knows nothing about any provider's property model.
type Resource struct {
	// TypeName and Identifier are the resource's scan coordinates (the
	// same currency claims and overrides use).
	TypeName   string
	Identifier string
	// Name is the human-authored name surface ("" when the account's
	// owners never named it). NEVER an opaque cloud identifier -- a random
	// id can embed a misleading token, so adapters must only populate
	// Name from surfaces a human wrote (a Name tag, a chosen queue/bucket/
	// role name).
	Name string
	// Tags are the resource's cloud tags.
	Tags map[string]string
	// Containers are the identifiers of the resources this one lives
	// inside (a subnet's VPC, a NAT gateway's subnet). Used only by
	// containment inheritance, and only for containers present in the
	// same partitioned set.
	Containers []string
}

// Tiers of the precedence ladder, recorded on every assignment so review
// surfaces can show WHY a resource landed where it did.
const (
	TierOverride    = "override"
	TierTag         = "tag"
	TierNameToken   = "name-token"
	TierContainment = "containment"
	TierFallback    = "fallback"
)

// Assignment is the engine's verdict for one resource. Environment is ""
// when the resource is honestly unassigned (no signal, or contradictory
// signals -- Conflicts says which).
type Assignment struct {
	TypeName    string   `json:"typeName"`
	Identifier  string   `json:"identifier"`
	Environment string   `json:"environment,omitempty"`
	Tier        string   `json:"tier,omitempty"`
	Signal      string   `json:"signal,omitempty"`
	Conflicts   []string `json:"conflicts,omitempty"`
}

// Result is the complete partition of one resource set: one assignment per
// input resource, in input order (determinism is part of the contract --
// see the package doc).
type Result struct {
	Assignments []Assignment `json:"assignments"`
}

// Environments returns the distinct assigned environment names, sorted.
func (r *Result) Environments() []string {
	seen := map[string]bool{}
	var envs []string
	for _, a := range r.Assignments {
		if a.Environment != "" && !seen[a.Environment] {
			seen[a.Environment] = true
			envs = append(envs, a.Environment)
		}
	}
	sort.Strings(envs)
	return envs
}

// ByRef returns the assignments keyed by scan coordinates.
func (r *Result) ByRef() map[ResourceRef]Assignment {
	byRef := make(map[ResourceRef]Assignment, len(r.Assignments))
	for _, a := range r.Assignments {
		byRef[ResourceRef{TypeName: a.TypeName, Identifier: a.Identifier}] = a
	}
	return byRef
}

// ResourceRef is a resource's scan coordinates.
type ResourceRef struct {
	TypeName   string `json:"typeName"`
	Identifier string `json:"identifier"`
}

// nameTokenSeparators split a human-authored name into whole tokens.
// Matching is whole-token exact -- never substring -- so "prodfix" or
// "reproduction" can never read as an environment signal.
const nameTokenSeparators = "-_./ "

// Partition assigns every resource to an environment per the rule's
// precedence ladder, or leaves it honestly unassigned. Deterministic: the
// same rule and resources always produce the same result, byte for byte.
func Partition(rule *Rule, resources []Resource) *Result {
	result := &Result{Assignments: make([]Assignment, len(resources))}

	// Tiers 1-3: each resource's OWN signals, strictly in ladder order.
	// A resource whose own signals contradict each other is conflicted and
	// stays unassigned -- and is excluded from inheritance and fallback,
	// because both would bury a contradiction the review gate must see.
	for i, resource := range resources {
		a := Assignment{TypeName: resource.TypeName, Identifier: resource.Identifier}

		if env, taught := rule.overrides[resourceKey{typeName: resource.TypeName, identifier: resource.Identifier}]; taught {
			a.Environment, a.Tier, a.Signal = env, TierOverride, "taught per-resource override"
			result.Assignments[i] = a
			continue
		}

		if env, signal, conflict := envFromAuthoritativeTags(rule, resource.Tags); conflict != "" {
			a.Conflicts = append(a.Conflicts, conflict)
			result.Assignments[i] = a
			continue
		} else if env != "" {
			a.Environment, a.Tier, a.Signal = env, TierTag, signal
			result.Assignments[i] = a
			continue
		}

		if env, signal, conflict := envFromNameTokens(rule, resource.Name); conflict != "" {
			a.Conflicts = append(a.Conflicts, conflict)
			result.Assignments[i] = a
			continue
		} else if env != "" {
			a.Environment, a.Tier, a.Signal = env, TierNameToken, signal
			result.Assignments[i] = a
			continue
		}

		result.Assignments[i] = a
	}

	if rule.containmentInheritance {
		inheritFromContainers(resources, result)
	}

	// Fallback: only for resources no tier reached and no conflict marked.
	if rule.fallbackEnvironment != "" {
		for i := range result.Assignments {
			a := &result.Assignments[i]
			if a.Environment == "" && len(a.Conflicts) == 0 {
				a.Environment, a.Tier = rule.fallbackEnvironment, TierFallback
				a.Signal = "rule fallback environment"
			}
		}
	}
	return result
}

// envFromAuthoritativeTags reads tier 2: the declared tag keys, scanned in
// declared order (the deterministic order; the keys are expected to agree).
// One distinct environment wins; disagreement is a conflict, never a guess.
func envFromAuthoritativeTags(rule *Rule, tags map[string]string) (env, signal, conflict string) {
	var found []string // "key=value" per present authoritative tag, declared order
	distinct := map[string]bool{}
	for _, key := range rule.authoritativeTagKeys {
		value, present := tags[key]
		if !present || value == "" {
			continue
		}
		canonical := rule.normalizeEnvValue(value)
		found = append(found, fmt.Sprintf("%s=%s", key, value))
		if !distinct[canonical] {
			distinct[canonical] = true
			env = canonical
		}
	}
	if len(distinct) > 1 {
		return "", "", "authoritative tags disagree: " + strings.Join(found, ", ")
	}
	if len(distinct) == 0 {
		return "", "", ""
	}
	return env, "tag " + found[0], ""
}

// envFromNameTokens reads tier 3: whole tokens of the human-authored name,
// matched against the DECLARED vocabulary only. Multiple tokens naming the
// same environment agree; tokens naming different environments are a
// conflict ("prod-to-staging-sync" is nobody's to guess).
func envFromNameTokens(rule *Rule, name string) (env, signal, conflict string) {
	if name == "" {
		return "", "", ""
	}
	var matched []string // matching tokens, in order of appearance
	distinct := map[string]bool{}
	for _, token := range strings.FieldsFunc(strings.ToLower(name), func(r rune) bool {
		return strings.ContainsRune(nameTokenSeparators, r)
	}) {
		canonical, declared := rule.tokenToEnv[token]
		if !declared {
			continue
		}
		matched = append(matched, token)
		if !distinct[canonical] {
			distinct[canonical] = true
			env = canonical
		}
	}
	if len(distinct) > 1 {
		return "", "", fmt.Sprintf("name %q carries tokens of different environments: %s", name, strings.Join(matched, ", "))
	}
	if len(distinct) == 0 {
		return "", "", ""
	}
	return env, fmt.Sprintf("name token %q in %q", matched[0], name), ""
}

// inheritFromContainers runs tier 4 to a fixpoint: a resource with no
// signal of its own inherits when every container of it PRESENT in the
// set is assigned and they all agree. Containers absent from the set are
// ignored (nothing is known about them); a still-unassigned container
// defers the decision to a later round, so multi-hop chains (NAT gateway
// -> subnet -> VPC) resolve without any topological bookkeeping. The loop
// terminates because every round either assigns at least one resource or
// stops.
func inheritFromContainers(resources []Resource, result *Result) {
	envByIdentifier := map[string]string{}
	assignedCount := 0
	for _, a := range result.Assignments {
		if a.Environment != "" {
			envByIdentifier[a.Identifier] = a.Environment
			assignedCount++
		}
	}
	present := map[string]bool{}
	for _, r := range resources {
		present[r.Identifier] = true
	}

	for {
		progress := false
		for i, resource := range resources {
			a := &result.Assignments[i]
			if a.Environment != "" || len(a.Conflicts) > 0 {
				continue
			}
			env, container, decided := containerConsensus(resource, present, envByIdentifier)
			if !decided || env == "" {
				continue
			}
			a.Environment, a.Tier = env, TierContainment
			a.Signal = fmt.Sprintf("inherited from container %q", container)
			envByIdentifier[resource.Identifier] = env
			progress = true
		}
		if !progress {
			break
		}
	}

	// Final sweep, once the fixpoint settled: flag what inheritance could
	// not honestly resolve (containers disagree) and own-signal
	// assignments that contradict their container -- the resource keeps
	// its own signal (a closer, more explicit fact), but the review gate
	// must see the tension.
	for i, resource := range resources {
		a := &result.Assignments[i]
		if len(a.Conflicts) > 0 {
			continue
		}
		env, container, decided := containerConsensus(resource, present, envByIdentifier)
		if !decided {
			continue
		}
		switch {
		case a.Environment == "" && env == "":
			a.Conflicts = append(a.Conflicts, "containers disagree: "+container)
		case a.Environment != "" && env != "" && a.Tier != TierContainment && env != a.Environment:
			a.Conflicts = append(a.Conflicts, fmt.Sprintf("own signal says %q but container %q is in %q", a.Environment, container, env))
		}
	}
}

// containerConsensus inspects a resource's containers that are present in
// the set. decided is false while any present container is still
// unassigned (the answer may yet arrive) or when the resource has no
// present containers at all. When decided: one agreed environment returns
// (env, firstContainer); disagreement returns env == "" with the
// disagreement rendered in container.
func containerConsensus(resource Resource, present map[string]bool, envByIdentifier map[string]string) (env, container string, decided bool) {
	var firstContainer string
	var disagreement []string
	distinct := map[string]bool{}
	sawPresent := false
	for _, id := range resource.Containers {
		if !present[id] {
			continue
		}
		sawPresent = true
		containerEnv, assigned := envByIdentifier[id]
		if !assigned {
			return "", "", false
		}
		if firstContainer == "" {
			firstContainer = id
		}
		disagreement = append(disagreement, fmt.Sprintf("%s in %q", id, containerEnv))
		if !distinct[containerEnv] {
			distinct[containerEnv] = true
			env = containerEnv
		}
	}
	if !sawPresent {
		return "", "", false
	}
	if len(distinct) > 1 {
		return "", strings.Join(disagreement, ", "), true
	}
	return env, firstContainer, true
}
