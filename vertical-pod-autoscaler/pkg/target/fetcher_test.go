/*
Copyright 2019 The Kubernetes Authors.

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

package target

import (
	"context"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	vpa_types "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"
)

func TestVpaTargetSelectorFetcher_PodLabelSelector(t *testing.T) {
	fetcher := &vpaTargetSelectorFetcher{}
	selector := &metav1.LabelSelector{
		MatchLabels: map[string]string{"app": "test"},
	}
	vpa := &vpa_types.VerticalPodAutoscaler{
		Spec: vpa_types.VerticalPodAutoscalerSpec{
			PodLabelSelector: selector,
		},
	}
	ctx := context.TODO()
	result, err := fetcher.Fetch(ctx, vpa)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.String() != labels.Set(selector.MatchLabels).AsSelector().String() {
		t.Errorf("expected selector %v, got %v", labels.Set(selector.MatchLabels).AsSelector().String(), result.String())
	}
}

func TestVpaTargetSelectorFetcher_TargetRefNil(t *testing.T) {
	fetcher := &vpaTargetSelectorFetcher{}
	vpa := &vpa_types.VerticalPodAutoscaler{
		Spec: vpa_types.VerticalPodAutoscalerSpec{},
	}
	ctx := context.TODO()
	_, err := fetcher.Fetch(ctx, vpa)
	if err == nil {
		t.Fatalf("expected error when neither podLabelSelector nor targetRef is set")
	}
}
