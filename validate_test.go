package main

import (
	"encoding/json"
	"testing"

	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	metav1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"

	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

func TestMutate(t *testing.T) {
	pod := corev1.Pod{
		Metadata: &metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: "test-tenant",
		},
	}

	payload, err := buildValidationRequest(&pod, &struct{}{}, "Pod")
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	responsePayload, err := validate(payload)
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var response kubewarden_protocol.ValidationResponse
	if err = json.Unmarshal(responsePayload, &response); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if response.Accepted != true {
		t.Error("Unexpected rejection", response.Message)
	}

	if response.MutatedObject == nil {
		t.Error("Expected mutation")
	}

	mutatedRequestJSON, err := json.Marshal(response.MutatedObject.(map[string]interface{}))
	if err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	var mutatedPod corev1.Pod
	if err = json.Unmarshal(mutatedRequestJSON, &mutatedPod); err != nil {
		t.Errorf("Unexpected error: %+v", err)
	}

	if mutatedPod.Metadata.Labels["tenant"] != pod.Metadata.Namespace {
		t.Errorf("Missing tenant label")
	}
}

func buildValidationRequest(object, settings interface{}, kind string) ([]byte, error) {
	objectRaw, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	kubeAdmissionReq := kubewarden_protocol.KubernetesAdmissionRequest{
		Object: objectRaw,
		Kind: kubewarden_protocol.GroupVersionKind{
			Kind: kind,
		},
	}

	settingsRaw, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	validationRequest := kubewarden_protocol.ValidationRequest{
		Request:  kubeAdmissionReq,
		Settings: settingsRaw,
	}

	return json.Marshal(validationRequest)
}
