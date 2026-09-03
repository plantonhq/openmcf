package component

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	v1 "github.com/plantonhq/planton/operator/api/v1"
)

// reconcileIdentityProviderBindings resolves which PlantonIdentityProviders
// in the namespace bind to this platform and records the resolution on each
// resource's status -- boundPlatform plus the Bound condition -- so an empty
// platformRef's resolution is always visible, never implicit.
//
// The binding rules (each with its reason):
//   - An explicit platformRef binds to exactly the named platform; a ref
//     naming another platform is that platform's business and is skipped
//     here entirely.
//   - An empty ref resolves to the namespace's ONLY platform. Two platforms
//     in the namespace with no ref is a named status error -- never a guess,
//     because guessing wrong here silently federates a company's directory
//     into the wrong Planton.
//   - At most ONE resource binds per platform (the v1 constraint). Several
//     candidates resolve deterministically to the oldest (then
//     lexicographically first) and the others get a condition naming the
//     winner -- deterministic, so two reconcile passes can never flip-flop.
//
// Federation provisioning consumes the winning binding, returned here (nil
// when nothing binds to this platform).
func (id *Identity) reconcileIdentityProviderBindings(ctx context.Context, c client.Client, planton *v1.PlantonPlatform) (*v1.PlantonIdentityProvider, error) {
	var idps v1.PlantonIdentityProviderList
	if err := c.List(ctx, &idps, client.InNamespace(planton.Namespace)); err != nil {
		return nil, fmt.Errorf("listing PlantonIdentityProviders: %w", err)
	}
	if len(idps.Items) == 0 {
		return nil, nil
	}

	var platforms v1.PlantonPlatformList
	if err := c.List(ctx, &platforms, client.InNamespace(planton.Namespace)); err != nil {
		return nil, fmt.Errorf("listing PlantonPlatforms: %w", err)
	}
	platformNames := make([]string, 0, len(platforms.Items))
	for i := range platforms.Items {
		platformNames = append(platformNames, platforms.Items[i].Name)
	}
	sort.Strings(platformNames)

	// Resolve every resource first, then pick the per-platform winner, so
	// the at-most-one constraint is decided over the full candidate set.
	var candidates []*v1.PlantonIdentityProvider
	var ambiguous []*v1.PlantonIdentityProvider
	for i := range idps.Items {
		idp := &idps.Items[i]
		switch {
		case idp.Spec.PlatformRef != nil && idp.Spec.PlatformRef.Name != "":
			if idp.Spec.PlatformRef.Name == planton.Name {
				candidates = append(candidates, idp)
			}
			// A ref to another platform is resolved by that platform's pass.
		case len(platforms.Items) == 1:
			candidates = append(candidates, idp)
		default:
			ambiguous = append(ambiguous, idp)
		}
	}

	sort.Slice(candidates, func(a, b int) bool {
		at, bt := candidates[a].CreationTimestamp, candidates[b].CreationTimestamp
		if !at.Equal(&bt) {
			return at.Before(&bt)
		}
		return candidates[a].Name < candidates[b].Name
	})

	var winner *v1.PlantonIdentityProvider
	for i, idp := range candidates {
		if i == 0 {
			winner = idp
			id.setBindingStatus(ctx, c, idp, planton.Name, metav1.Condition{
				Type: v1.ConditionBound, Status: metav1.ConditionTrue,
				Reason:  "Resolved",
				Message: fmt.Sprintf("Bound to PlantonPlatform %s", planton.Name),
			})
			continue
		}
		id.setBindingStatus(ctx, c, idp, "", metav1.Condition{
			Type: v1.ConditionBound, Status: metav1.ConditionFalse,
			Reason: "PlatformAlreadyBound",
			Message: fmt.Sprintf(
				"PlantonPlatform %s is already bound to PlantonIdentityProvider %s -- at most one identity provider binds per platform; delete or repoint one of them",
				planton.Name, candidates[0].Name),
		})
	}

	for _, idp := range ambiguous {
		id.setBindingStatus(ctx, c, idp, "", metav1.Condition{
			Type: v1.ConditionBound, Status: metav1.ConditionFalse,
			Reason: "AmbiguousPlatform",
			Message: fmt.Sprintf(
				"several PlantonPlatforms exist in this namespace (%s) -- set spec.platformRef.name to the one this identity config belongs to",
				strings.Join(platformNames, ", ")),
		})
	}

	return winner, nil
}

// setBindingStatus writes boundPlatform + the Bound condition, touching the
// API only when something actually changed -- the reconcile loop runs every
// 30s and an unconditional write would churn the resource forever.
func (id *Identity) setBindingStatus(ctx context.Context, c client.Client, idp *v1.PlantonIdentityProvider, boundPlatform string, condition metav1.Condition) {
	log := logf.FromContext(ctx).WithValues("component", id.Name())

	changed := meta.SetStatusCondition(&idp.Status.Conditions, condition)
	if idp.Status.BoundPlatform != boundPlatform {
		idp.Status.BoundPlatform = boundPlatform
		changed = true
	}
	if !changed {
		return
	}
	// A conflict here is not a component failure: the next pass re-resolves
	// from fresh reads, so log-and-continue is the whole recovery.
	if err := c.Status().Update(ctx, idp); err != nil {
		log.Error(err, "Failed to update PlantonIdentityProvider binding status", "name", idp.Name)
		return
	}
	log.Info("Resolved PlantonIdentityProvider binding",
		"name", idp.Name, "boundPlatform", boundPlatform, "reason", condition.Reason)
}
