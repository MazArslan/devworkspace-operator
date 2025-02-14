// Copyright (c) 2019-2023 Red Hat, Inc.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package kubernetes

import (
	"fmt"

	dw "github.com/devfile/api/v2/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/devworkspace-operator/pkg/common"
	"github.com/devfile/devworkspace-operator/pkg/dwerrors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func HasKubelikeComponent(workspace *common.DevWorkspaceWithConfig) bool {
	for _, component := range workspace.Spec.Template.Components {
		if component.Kubernetes != nil || component.Openshift != nil {
			return true
		}
	}
	return false
}

func filterForKubeLikeComponents(components []dw.Component) ([]dw.Component, error) {
	var k8sLikeComponents []dw.Component
	for _, component := range components {
		k8sLikeComponent, err := getK8sLikeComponent(component)
		if err != nil {
			continue
		}
		if k8sLikeComponent.Uri != "" {
			return nil, &dwerrors.FailError{Message: fmt.Sprintf("kubernetes/openshift components that define a URI are unsupported (component %s)", component.Name)}
		}
		if k8sLikeComponent.Inlined == "" {
			continue
		}
		if k8sLikeComponent.GetDeployByDefault() {
			k8sLikeComponents = append(k8sLikeComponents, component)
		}
	}
	return k8sLikeComponents, nil
}

// getK8sLikeComponent returns the K8sLikeComponent from a DevWorkspace component,
// allowing Kubernetes and OpenShift components to be processed in the same way.
// Returns error if component is not a kube-like component.
func getK8sLikeComponent(component dw.Component) (*dw.K8sLikeComponent, error) {
	switch {
	case component.Kubernetes != nil:
		return &component.Kubernetes.K8sLikeComponent, nil
	case component.Openshift != nil:
		return &component.Openshift.K8sLikeComponent, nil
	default:
		return nil, fmt.Errorf("not a kube-like component")
	}
}

func checkOwnerrefs(ownerrefs, subset []metav1.OwnerReference) error {
	for _, checkOwnerref := range subset {
		found := false
		for _, ownerref := range ownerrefs {
			if ownerref.Name == checkOwnerref.Name && ownerref.UID == checkOwnerref.UID {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("ownerref not found")
		}
	}
	return nil
}
