/*
Copyright 2026.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package singleton enforces the one-operator-per-cluster invariant.
//
// The operator watches PlantonPlatform resources in EVERY namespace, but its
// leader-election lease lives in its OWN namespace -- so two operator
// installations in different namespaces both win their local elections and
// double-reconcile the same platforms, fighting over server-side-apply field
// ownership. Leader election protects replicas of one installation, never the
// cluster. This guard closes that gap at the only reliable moment: startup.
// The losing operator exits with a plain-language remedy instead of silently
// corrupting a running platform, and because the check re-runs on every
// restart, removing the other installation lets the waiting one recover on
// its own.
package singleton

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Operator Deployments are recognized by the labels the planton-operator
// Helm chart stamps on them (both the standalone chart and the umbrella
// chart's subchart render the same labels). A dev deployment applied through
// kustomize does not carry the chart's name label and is deliberately outside
// the guard -- the invariant protects real installations, not lab runs.
const (
	operatorNameLabel      = "app.kubernetes.io/name"
	operatorNameLabelValue = "planton-operator"
	controlPlaneLabel      = "control-plane"
	controlPlaneLabelValue = "controller-manager"
)

// serviceAccountNamespaceFile is mounted into every pod by Kubernetes; its
// absence means the process runs outside a cluster (e.g. `make run`), where
// the guard has nothing meaningful to check.
const serviceAccountNamespaceFile = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"

// OwnNamespace reports the namespace this process runs in, or ok=false when
// running outside a cluster.
func OwnNamespace() (namespace string, ok bool) {
	raw, err := os.ReadFile(serviceAccountNamespaceFile)
	if err != nil {
		return "", false
	}
	ns := strings.TrimSpace(string(raw))
	return ns, ns != ""
}

// Check returns an error naming every OTHER namespace that runs an active
// planton-operator installation. Same-namespace Deployments are always
// allowed (HA replicas and rolling upgrades of one installation), and a
// Deployment scaled to zero replicas is not managing anything, so it does
// not block startup either.
func Check(ctx context.Context, reader client.Reader, ownNamespace string) error {
	var deployments appsv1.DeploymentList
	if err := reader.List(ctx, &deployments, client.MatchingLabels{
		operatorNameLabel: operatorNameLabelValue,
		controlPlaneLabel: controlPlaneLabelValue,
	}); err != nil {
		return fmt.Errorf("listing operator Deployments to enforce one operator per cluster: %w", err)
	}

	foreign := map[string]struct{}{}
	for _, d := range deployments.Items {
		if d.Namespace == ownNamespace {
			continue
		}
		if d.Spec.Replicas != nil && *d.Spec.Replicas == 0 {
			continue
		}
		foreign[d.Namespace] = struct{}{}
	}
	if len(foreign) == 0 {
		return nil
	}

	namespaces := make([]string, 0, len(foreign))
	for ns := range foreign {
		namespaces = append(namespaces, ns)
	}
	sort.Strings(namespaces)

	return fmt.Errorf(
		"another planton-operator is already managing this cluster from namespace %q; "+
			"one operator per cluster -- running two would make them fight over the same platforms. "+
			"Uninstall one of the two operator releases (a platform is declared separately, with the planton chart "+
			"or a PlantonPlatform manifest, and needs no operator of its own). "+
			"This operator will keep restarting and will recover on its own once the conflict is gone",
		strings.Join(namespaces, `", "`))
}
