package v1alpha1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
    ddv1alpha1 "github.com/DataDog/extendeddaemonset/api/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/uswitch/nidhogg/pkg/test"
)

const (
	nodeName  = utils.NodeName
	namespace = utils.Namespace
	daemonset = utils.Daemonset
)



func TestGetPodsReturnsUniquePods(t *testing.T) {
	ctx := context.TODO()

	// DD Extended Daemonset owner is DD Extended Daemonset Replica set
	// Additionally, during rollout, there are two owners: active and canary,
	// hence we need to check for both
	canaryOwner := "CANARY"
	activeOwner := "ACTIVE"

	pod1 := utils.BuildPod("pod1", daemonset, corev1.PodReady)
	pod1.OwnerReferences = append(pod1.OwnerReferences, metav1.OwnerReference{Name: activeOwner})
	pod2 := utils.BuildPod("pod2", daemonset, corev1.PodReady)
	pod2.OwnerReferences = append(pod2.OwnerReferences, metav1.OwnerReference{Name: canaryOwner})
	// Shouldn't match owner
	pod3 := utils.BuildPod("pod3", daemonset, corev1.PodReady)

	ds := ExtendedDaemonset{
		Name: daemonset,
		Namespace: namespace,
		List: func(ctx Context, list ObjectList, opts ...ListOption) error{
			podList, _ := list.(*corev1.PodList)
			podList.Items = append(podList.Items, pod1)
			podList.Items = append(podList.Items, pod2)
			podList.Items = append(podList.Items, pod3)

			return nil
		},
		Get: func(ctx Context, key ObjectKey, obj Object, opts ...GetOption) error{
			extendedDaemonset, _ := obj.(*ddv1alpha1.ExtendedDaemonSet)
			extendedDaemonset.Status.ActiveReplicaSet = activeOwner
			canary := ddv1alpha1.ExtendedDaemonSetStatusCanary{ ReplicaSet: canaryOwner}
			extendedDaemonset.Status.Canary = &canary
			return nil
		},
	}

	pods, err := ds.GetPods(ctx, nodeName)

	assert.NoError(t, err)
	assert.NotNil(t, pods)
	assert.Equal(t, 2, len(pods))
	assert.Equal(t, pod1.Name, pods[0].Name)
	assert.Equal(t, pod2.Name, pods[1].Name)
}
