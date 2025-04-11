package nidhogg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/uswitch/nidhogg/pkg/test"
)

const (
	nodeSelector = utils.NodeSelector

	namespace = utils.Namespace
	daemonset = utils.Daemonset

	daemonset1 = utils.Daemonset1
	daemonset2 = utils.Daemonset2

	taintNamePrefix = utils.TaintNamePrefix
	taintName       = utils.TaintName
)

func TestCalculateTaintsWithReadyPod(t *testing.T) {
	ctx := context.TODO()
	node := utils.BuildNode(namespace, []string{daemonset})
	pod := utils.BuildPod("pod", daemonset, corev1.PodReady)
	cfg := BuildNidhoggConfig(namespace, []string{daemonset})
	cfg.BuildSelectors()

	handler := BuildHandler([]corev1.Pod{pod}, nil, cfg)
	updatedNode, changes, err := handler.calculateTaints(ctx, &node)

	assert.NoError(t, err)
	assert.NotNil(t, updatedNode)
	assert.NotContains(t, updatedNode.Spec.Taints, taintName)
	assert.Empty(t, changes.taintsAdded)
	assert.Contains(t, changes.taintsRemoved, taintName)
}

func TestCalculateTaintsWithReadyPodAndWithoutNodeSelector(t *testing.T) {
	ctx := context.TODO()
	node := utils.BuildNode(namespace, []string{daemonset})
	pod := utils.BuildPod("pod", daemonset, corev1.PodReady)
	cfg := BuildNidhoggConfigWithoutNodeSelector(namespace, []string{daemonset})
	cfg.BuildSelectors()

	handler := BuildHandler([]corev1.Pod{pod}, []appsv1.DaemonSet{utils.BuildDaemonset(daemonset)}, cfg)
	updatedNode, changes, err := handler.calculateTaints(ctx, &node)

	assert.NoError(t, err)
	assert.NotNil(t, updatedNode)
	assert.NotContains(t, updatedNode.Spec.Taints, taintName)
	assert.Empty(t, changes.taintsAdded)
	assert.Contains(t, changes.taintsRemoved, taintName)
}

func TestCalculateTaintsWithMultipleDaemonsets(t *testing.T) {
	ctx := context.TODO()
	node := utils.BuildNode(namespace, []string{daemonset1, daemonset2})
	pod1 := utils.BuildPod("pod1", daemonset1, corev1.PodReady)
	pod2 := utils.BuildPod("pod2", daemonset2, corev1.PodScheduled)
	cfg := BuildNidhoggConfig(namespace, []string{daemonset1, daemonset2})
	cfg.BuildSelectors()

	handler := BuildHandler([]corev1.Pod{pod1, pod2}, nil, cfg)
	updatedNode, changes, err := handler.calculateTaints(ctx, &node)

	assert.NoError(t, err)
	assert.NotNil(t, updatedNode)
	assert.NotContains(t, updatedNode.Spec.Taints, utils.BuildActiveTaint(namespace, daemonset1))
	assert.Contains(t, updatedNode.Spec.Taints, utils.BuildActiveTaint(namespace, daemonset2))
	assert.Contains(t, changes.taintsRemoved, utils.BuildTaintName(namespace, daemonset1))
	assert.Empty(t, changes.taintsAdded)
}

func TestCalculateTaintsWithUnreadyPod(t *testing.T) {
	ctx := context.TODO()
	node := utils.BuildNode(namespace, []string{daemonset})
	pod := utils.BuildPod("pod", daemonset, corev1.PodScheduled)
	cfg := BuildNidhoggConfig(namespace, []string{daemonset})
	cfg.BuildSelectors()

	handler := BuildHandler([]corev1.Pod{pod}, nil, cfg)
	updatedNode, changes, err := handler.calculateTaints(ctx, &node)

	assert.NoError(t, err)
	assert.NotNil(t, updatedNode)
	assert.Contains(t, updatedNode.Spec.Taints, utils.BuildActiveTaint(namespace, daemonset))
	assert.Empty(t, changes.taintsRemoved)
	assert.Empty(t, changes.taintsAdded, taintName)
}

func TestCalculateTaintsWithUnreadyPodAndWithoutNodeSelector(t *testing.T) {
	ctx := context.TODO()
	node := utils.BuildNode(namespace, []string{daemonset})
	pod := utils.BuildPod("pod", daemonset, corev1.PodScheduled)
	cfg := BuildNidhoggConfigWithoutNodeSelector(namespace, []string{daemonset})
	cfg.BuildSelectors()

	handler := BuildHandler([]corev1.Pod{pod}, []appsv1.DaemonSet{utils.BuildDaemonset(daemonset)}, cfg)
	updatedNode, changes, err := handler.calculateTaints(ctx, &node)

	assert.NoError(t, err)
	assert.NotNil(t, updatedNode)
	assert.Contains(t, updatedNode.Spec.Taints, utils.BuildActiveTaint(namespace, daemonset))
	assert.Empty(t, changes.taintsRemoved)
	assert.Empty(t, changes.taintsAdded, taintName)
}

func BuildHandler(pods []corev1.Pod, daemonsets []appsv1.DaemonSet, config HandlerConfig) Handler {
	return Handler{
		Client: fake.NewClientBuilder().WithLists(&corev1.PodList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items:    pods,
		}, &appsv1.DaemonSetList{
			TypeMeta: metav1.TypeMeta{},
			ListMeta: metav1.ListMeta{},
			Items:    daemonsets,
		}).Build(),
		recorder: record.NewFakeRecorder(0),
		config:   config,
	}
}

func BuildDaemonsets(namespace string, daemonsetNames []string) []DaemonsetConfig {
	var daemonsets []DaemonsetConfig
	for _, daemonsetName := range daemonsetNames {
		daemonsets = append(daemonsets, DaemonsetConfig{Name: daemonsetName, Namespace: namespace})
	}
	return daemonsets
}

func BuildNidhoggConfig(namespace string, daemonsets []string) HandlerConfig {
	return HandlerConfig{
		TaintNamePrefix:            taintNamePrefix,
		TaintRemovalDelayInSeconds: 0,
		Daemonsets:                 BuildDaemonsets(namespace, daemonsets),
		NodeSelector:               []string{utils.NodeSelector},
	}
}

func BuildNidhoggConfigWithoutNodeSelector(namespace string, daemonsets []string) HandlerConfig {
	return HandlerConfig{
		TaintNamePrefix:            taintNamePrefix,
		TaintRemovalDelayInSeconds: 0,
		Daemonsets:                 BuildDaemonsets(namespace, daemonsets),
	}
}
