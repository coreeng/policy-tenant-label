package main

import (
	"encoding/json"

	appsv1 "github.com/kubewarden/k8s-objects/api/apps/v1"
	batchv1 "github.com/kubewarden/k8s-objects/api/batch/v1"
	corev1 "github.com/kubewarden/k8s-objects/api/core/v1"
	apimachinery_pkg_apis_meta_v1 "github.com/kubewarden/k8s-objects/apimachinery/pkg/apis/meta/v1"
	kubewarden "github.com/kubewarden/policy-sdk-go"
	kubewarden_protocol "github.com/kubewarden/policy-sdk-go/protocol"
)

const httpBadRequestStatusCode = 400

func validate(payload []byte) ([]byte, error) {
	validationRequest := kubewarden_protocol.ValidationRequest{}
	err := json.Unmarshal(payload, &validationRequest)
	if err != nil {
		return kubewarden.RejectRequest(
			kubewarden.Message(err.Error()),
			kubewarden.Code(httpBadRequestStatusCode))
	}

	return mutateRequest(validationRequest)
}

func mutateRequest(validationRequest kubewarden_protocol.ValidationRequest) ([]byte, error) { //nolint:funlen
	switch validationRequest.Request.Kind.Kind {
	case "Deployment":
		deployment := appsv1.Deployment{}
		if err := json.Unmarshal(validationRequest.Request.Object, &deployment); err != nil {
			return nil, err
		}
		addTenantLabel(deployment.Metadata)
		return kubewarden.MutateRequest(deployment)
	case "ReplicaSet":
		replicaset := appsv1.ReplicaSet{}
		if err := json.Unmarshal(validationRequest.Request.Object, &replicaset); err != nil {
			return nil, err
		}
		addTenantLabel(replicaset.Metadata)
		return kubewarden.MutateRequest(replicaset)
	case "StatefulSet":
		statefulset := appsv1.StatefulSet{}
		if err := json.Unmarshal(validationRequest.Request.Object, &statefulset); err != nil {
			return nil, err
		}
		addTenantLabel(statefulset.Metadata)
		return kubewarden.MutateRequest(statefulset)
	case "DaemonSet":
		daemonset := appsv1.DaemonSet{}
		if err := json.Unmarshal(validationRequest.Request.Object, &daemonset); err != nil {
			return nil, err
		}
		addTenantLabel(daemonset.Metadata)
		return kubewarden.MutateRequest(daemonset)
	case "ReplicationController":
		replicationController := corev1.ReplicationController{}
		if err := json.Unmarshal(validationRequest.Request.Object, &replicationController); err != nil {
			return nil, err
		}
		addTenantLabel(replicationController.Metadata)
		return kubewarden.MutateRequest(replicationController)
	case "CronJob":
		cronjob := batchv1.CronJob{}
		if err := json.Unmarshal(validationRequest.Request.Object, &cronjob); err != nil {
			return nil, err
		}
		addTenantLabel(cronjob.Metadata)
		return kubewarden.MutateRequest(cronjob)
	case "Job":
		job := batchv1.Job{}
		if err := json.Unmarshal(validationRequest.Request.Object, &job); err != nil {
			return nil, err
		}
		addTenantLabel(job.Metadata)
		return kubewarden.MutateRequest(job)
	case "Pod":
		pod := corev1.Pod{}
		if err := json.Unmarshal(validationRequest.Request.Object, &pod); err != nil {
			return nil, err
		}
		addTenantLabel(pod.Metadata)
		return kubewarden.MutateRequest(pod)
	default:
		return kubewarden.RejectRequest(
			"Object should be one of these kinds: Deployment, ReplicaSet, StatefulSet, DaemonSet, ReplicationController, Job, CronJob, Pod", //nolint:lll
			kubewarden.NoCode,
		)
	}
}

func addTenantLabel(metadata *apimachinery_pkg_apis_meta_v1.ObjectMeta) {
	if metadata.Labels == nil {
		metadata.Labels = make(map[string]string)
	}

	logger.Debug("adding tenant label")

	metadata.Labels["tenant"] = metadata.Namespace // this will be improved once we agree on how to get the tenant name
}
