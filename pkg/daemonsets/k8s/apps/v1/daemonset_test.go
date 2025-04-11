package v1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"

	"github.com/uswitch/nidhogg/pkg/test"
)

const (
	nodeName  = utils.NodeName
	namespace = utils.Namespace
	daemonset = utils.Daemonset
)

func TestGetPodsReturnsUniquePods(t *testing.T) {
	ctx := context.TODO()

	pod1 := utils.BuildPod("pod1", daemonset, corev1.PodReady)
	pod2 := utils.BuildPod("pod2", daemonset, corev1.PodReady)


	ds := Daemonset{
		Name: daemonset,
		Namespace: namespace,
		List: func(ctx Context, list ObjectList, opts ...ListOption) error{
			podList, _ := list.(*corev1.PodList)
			podList.Items = append(podList.Items, pod1)
			podList.Items = append(podList.Items, pod2)

			return nil
		},
		Get: nil,
	}

	pods, err := ds.GetPods(ctx, nodeName)

	assert.NoError(t, err)
	assert.NotNil(t, pods)
	assert.Equal(t, 2, len(pods))
	assert.Equal(t, pod1.Name, pods[0].Name)
	assert.Equal(t, pod2.Name, pods[1].Name)
}
