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

package singleton

import (
	"context"
	"strings"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func operatorDeployment(namespace, name string, replicas int32) *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
			Labels: map[string]string{
				operatorNameLabel: operatorNameLabelValue,
				controlPlaneLabel: controlPlaneLabelValue,
			},
		},
		Spec: appsv1.DeploymentSpec{Replicas: &replicas},
	}
}

func newReader(t *testing.T, objs ...client.Object) client.Reader {
	t.Helper()
	scheme := runtime.NewScheme()
	if err := clientgoscheme.AddToScheme(scheme); err != nil {
		t.Fatalf("building scheme: %v", err)
	}
	return fake.NewClientBuilder().WithScheme(scheme).WithObjects(objs...).Build()
}

func TestCheckAllowsOwnNamespace(t *testing.T) {
	reader := newReader(t, operatorDeployment("planton", "planton-operator", 1))
	if err := Check(context.Background(), reader, "planton"); err != nil {
		t.Fatalf("own-namespace Deployment must not block startup: %v", err)
	}
}

func TestCheckAllowsEmptyCluster(t *testing.T) {
	reader := newReader(t)
	if err := Check(context.Background(), reader, "planton"); err != nil {
		t.Fatalf("empty cluster must not block startup: %v", err)
	}
}

func TestCheckIgnoresUnrelatedDeployments(t *testing.T) {
	unrelated := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "other",
			Name:      "some-app",
			Labels:    map[string]string{"app.kubernetes.io/name": "some-app"},
		},
	}
	reader := newReader(t, unrelated)
	if err := Check(context.Background(), reader, "planton"); err != nil {
		t.Fatalf("unrelated Deployments must not block startup: %v", err)
	}
}

func TestCheckIgnoresScaledDownOperator(t *testing.T) {
	// An operator scaled to zero is not reconciling anything (the lab pattern
	// for pausing an install) -- blocking on it would be a false positive.
	reader := newReader(t, operatorDeployment("old-home", "planton-operator", 0))
	if err := Check(context.Background(), reader, "planton"); err != nil {
		t.Fatalf("scaled-to-zero operator must not block startup: %v", err)
	}
}

func TestCheckRefusesSecondOperator(t *testing.T) {
	reader := newReader(t,
		operatorDeployment("planton-operator-system", "planton-operator", 1),
		operatorDeployment("planton", "planton-planton-operator", 1),
	)
	err := Check(context.Background(), reader, "planton")
	if err == nil {
		t.Fatal("a second active operator in another namespace must block startup")
	}
	// The message is the product surface here: it must name the other
	// namespace, the way out, and that a platform needs no operator of its own.
	for _, want := range []string{
		"planton-operator-system",
		"one operator per cluster",
		"Uninstall one of the two operator releases",
		"needs no operator of its own",
		"recover on its own",
	} {
		if !strings.Contains(err.Error(), want) {
			t.Errorf("error message missing %q; got: %s", want, err.Error())
		}
	}
	if strings.Contains(err.Error(), `"planton"`) {
		t.Errorf("error message must not name the operator's own namespace as a conflict; got: %s", err.Error())
	}
}
