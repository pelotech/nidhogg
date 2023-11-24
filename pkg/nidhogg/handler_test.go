package nidhogg

import (
	"context"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

const (
	nodeName     = "nodeName"
	nodeSelector = "nodeSelector"

	namespace     = "namespace"
	daemonsetName = "daemonsetName"

	taintNamePrefix = "pelo.tech"
	taintName       = taintNamePrefix + "/" + namespace + "." + daemonsetName
)

func TestCalculateTaintsWithReadyPod(t *testing.T) {
	ctx := context.TODO()
	node := buildNode()
	pod := buildPod(corev1.PodReady)
	cfg := buildNidhoggConfig()

	handler := buildHandler(pod, cfg)
	updatedNode, changes, err := handler.calculateTaints(ctx, &node)

	assert.NoError(t, err)
	assert.NotNil(t, updatedNode)
	assert.NotContains(t, updatedNode.Spec.Taints, taintName)
	assert.Empty(t, changes.taintsAdded)
	assert.Contains(t, changes.taintsRemoved, taintName)
}

func TestCalculateTaintsWithUnreadyPod(t *testing.T) {
	ctx := context.TODO()
	node := buildNode()
	pod := buildPod(corev1.PodScheduled)
	cfg := buildNidhoggConfig()

	handler := buildHandler(pod, cfg)
	updatedNode, changes, err := handler.calculateTaints(ctx, &node)

	assert.NoError(t, err)
	assert.NotNil(t, updatedNode)
	assert.Contains(t, updatedNode.Spec.Taints, buildActiveTaint())
	assert.Empty(t, changes.taintsRemoved)
	assert.Empty(t, changes.taintsAdded, taintName)
}

func buildHandler(pod corev1.Pod, config HandlerConfig) Handler {
	return Handler{
		Client: fake.NewClientBuilder().WithLists(&corev1.PodList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items:    []corev1.Pod{pod},
		}).Build(),
		recorder: record.NewFakeRecorder(0),
		config:   config,
	}
}

func buildNidhoggConfig() HandlerConfig {
	return HandlerConfig{
		TaintNamePrefix:            taintNamePrefix,
		TaintRemovalDelayInSeconds: 0,
		Daemonsets: []Daemonset{
			{
				Name:      daemonsetName,
				Namespace: namespace,
			},
		},
		NodeSelector: []string{nodeSelector},
	}
}

func buildPod(conditionType corev1.PodConditionType) corev1.Pod {
	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:            "pod",
			Namespace:       namespace,
			OwnerReferences: []metav1.OwnerReference{{Name: daemonsetName}},
		},
		Spec: corev1.PodSpec{
			NodeName: nodeName,
		},
		Status: corev1.PodStatus{
			Phase: corev1.PodRunning,
			Conditions: []corev1.PodCondition{
				{
					Type:   conditionType,
					Status: corev1.ConditionTrue,
				},
			},
		},
	}
}

func buildNode() corev1.Node {
	return corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: nodeName,
			Labels: map[string]string{
				nodeSelector: "true",
			},
		},
		Spec: corev1.NodeSpec{
			Taints: []corev1.Taint{
				buildActiveTaint(),
			},
		},
	}
}

func buildActiveTaint() corev1.Taint {
	return corev1.Taint{
		Key:    taintName,
		Effect: corev1.TaintEffectNoSchedule,
		Value:  string(corev1.ConditionTrue),
	}
}
