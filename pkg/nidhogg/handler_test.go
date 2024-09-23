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
)

const (
	nodeName     = "nodeName"
	nodeSelector = "nodeSelector"

	namespace = "namespace"
	daemonset = "daemonset"

	daemonset1 = "daemonset1"
	daemonset2 = "daemonset2"

	taintNamePrefix = "pelo.tech"
	taintName       = taintNamePrefix + "/" + namespace + "." + daemonset
)

func TestCalculateTaintsWithReadyPod(t *testing.T) {
	ctx := context.TODO()
	node := buildNode(namespace, []string{daemonset})
	pod := buildPod("pod", daemonset, corev1.PodReady)
	cfg := buildNidhoggConfig(namespace, []string{daemonset})
	cfg.BuildSelectors()

	handler := buildHandler([]corev1.Pod{pod}, nil, cfg)
	updatedNode, changes, err := handler.calculateTaints(ctx, &node)

	assert.NoError(t, err)
	assert.NotNil(t, updatedNode)
	assert.NotContains(t, updatedNode.Spec.Taints, taintName)
	assert.Empty(t, changes.taintsAdded)
	assert.Contains(t, changes.taintsRemoved, taintName)
}

func TestCalculateTaintsWithReadyPodAndWithoutNodeSelector(t *testing.T) {
	ctx := context.TODO()
	node := buildNode(namespace, []string{daemonset})
	pod := buildPod("pod", daemonset, corev1.PodReady)
	cfg := buildNidhoggConfigWithoutNodeSelector(namespace, []string{daemonset})
	cfg.BuildSelectors()

	handler := buildHandler([]corev1.Pod{pod}, []appsv1.DaemonSet{buildDaemonset(daemonset)}, cfg)
	updatedNode, changes, err := handler.calculateTaints(ctx, &node)

	assert.NoError(t, err)
	assert.NotNil(t, updatedNode)
	assert.NotContains(t, updatedNode.Spec.Taints, taintName)
	assert.Empty(t, changes.taintsAdded)
	assert.Contains(t, changes.taintsRemoved, taintName)
}

func TestCalculateTaintsWithMultipleDaemonsets(t *testing.T) {
	ctx := context.TODO()
	node := buildNode(namespace, []string{daemonset1, daemonset2})
	pod1 := buildPod("pod1", daemonset1, corev1.PodReady)
	pod2 := buildPod("pod2", daemonset2, corev1.PodScheduled)
	cfg := buildNidhoggConfig(namespace, []string{daemonset1, daemonset2})
	cfg.BuildSelectors()

	handler := buildHandler([]corev1.Pod{pod1, pod2}, nil, cfg)
	updatedNode, changes, err := handler.calculateTaints(ctx, &node)

	assert.NoError(t, err)
	assert.NotNil(t, updatedNode)
	assert.NotContains(t, updatedNode.Spec.Taints, buildActiveTaint(namespace, daemonset1))
	assert.Contains(t, updatedNode.Spec.Taints, buildActiveTaint(namespace, daemonset2))
	assert.Contains(t, changes.taintsRemoved, buildTaintName(namespace, daemonset1))
	assert.Empty(t, changes.taintsAdded)
}

func TestCalculateTaintsWithUnreadyPod(t *testing.T) {
	ctx := context.TODO()
	node := buildNode(namespace, []string{daemonset})
	pod := buildPod("pod", daemonset, corev1.PodScheduled)
	cfg := buildNidhoggConfig(namespace, []string{daemonset})
	cfg.BuildSelectors()

	handler := buildHandler([]corev1.Pod{pod}, nil, cfg)
	updatedNode, changes, err := handler.calculateTaints(ctx, &node)

	assert.NoError(t, err)
	assert.NotNil(t, updatedNode)
	assert.Contains(t, updatedNode.Spec.Taints, buildActiveTaint(namespace, daemonset))
	assert.Empty(t, changes.taintsRemoved)
	assert.Empty(t, changes.taintsAdded, taintName)
}

func TestCalculateTaintsWithUnreadyPodAndWithoutNodeSelector(t *testing.T) {
	ctx := context.TODO()
	node := buildNode(namespace, []string{daemonset})
	pod := buildPod("pod", daemonset, corev1.PodScheduled)
	cfg := buildNidhoggConfigWithoutNodeSelector(namespace, []string{daemonset})
	cfg.BuildSelectors()

	handler := buildHandler([]corev1.Pod{pod}, []appsv1.DaemonSet{buildDaemonset(daemonset)}, cfg)
	updatedNode, changes, err := handler.calculateTaints(ctx, &node)

	assert.NoError(t, err)
	assert.NotNil(t, updatedNode)
	assert.Contains(t, updatedNode.Spec.Taints, buildActiveTaint(namespace, daemonset))
	assert.Empty(t, changes.taintsRemoved)
	assert.Empty(t, changes.taintsAdded, taintName)
}

func TestGetDaemonsetPodsReturnsUniquePods(t *testing.T) {
	ctx := context.TODO()
	pod1 := buildPod("pod1", daemonset, corev1.PodReady)
	pod2 := buildPod("pod2", daemonset, corev1.PodReady)
	cfg := buildNidhoggConfig(namespace, []string{daemonset})

	handler := buildHandler([]corev1.Pod{pod1, pod2}, nil, cfg)
	daemonset := Daemonset{Name: daemonset, Namespace: namespace}
	pods, err := handler.getDaemonsetPods(ctx, nodeName, daemonset)

	assert.NoError(t, err)
	assert.NotNil(t, pods)
	assert.Equal(t, len(pods), 2)
	assert.Equal(t, pods[0].Name, pod1.Name)
	assert.Equal(t, pods[1].Name, pod2.Name)
}

func buildHandler(pods []corev1.Pod, daemonsets []appsv1.DaemonSet, config HandlerConfig) Handler {
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

func buildDaemonsets(namespace string, daemonsetNames []string) []Daemonset {
	var daemonsets []Daemonset
	for _, daemonsetName := range daemonsetNames {
		daemonsets = append(daemonsets, Daemonset{Name: daemonsetName, Namespace: namespace})
	}
	return daemonsets
}

func buildNidhoggConfig(namespace string, daemonsets []string) HandlerConfig {
	return HandlerConfig{
		TaintNamePrefix:            taintNamePrefix,
		TaintRemovalDelayInSeconds: 0,
		Daemonsets:                 buildDaemonsets(namespace, daemonsets),
		NodeSelector:               []string{nodeSelector},
	}
}

func buildNidhoggConfigWithoutNodeSelector(namespace string, daemonsets []string) HandlerConfig {
	return HandlerConfig{
		TaintNamePrefix:            taintNamePrefix,
		TaintRemovalDelayInSeconds: 0,
		Daemonsets:                 buildDaemonsets(namespace, daemonsets),
	}
}

func buildPod(podName string, daemonsetName string, conditionType corev1.PodConditionType) corev1.Pod {
	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:            podName,
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

func buildDaemonset(daemonsetName string) appsv1.DaemonSet {
	return appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      daemonsetName,
			Namespace: namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{
						nodeSelector: "true",
					},
				},
			},
		},
	}
}

func buildNode(namespace string, daemonsets []string) corev1.Node {
	return corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: nodeName,
			Labels: map[string]string{
				nodeSelector: "true",
			},
		},
		Spec: corev1.NodeSpec{
			Taints: buildActiveTaintsFromDaemonsets(namespace, daemonsets),
		},
	}
}

func buildActiveTaintsFromDaemonsets(namespace string, daemonsets []string) []corev1.Taint {
	var taints []corev1.Taint
	for _, daemonset := range daemonsets {
		taints = append(taints, buildActiveTaint(namespace, daemonset))
	}
	return taints
}

func buildActiveTaint(namespace string, daemonset string) corev1.Taint {
	return corev1.Taint{
		Key:    buildTaintName(namespace, daemonset),
		Effect: corev1.TaintEffectNoSchedule,
		Value:  string(corev1.ConditionTrue),
	}
}

func buildTaintName(namespace string, daemonset string) string {
	return taintNamePrefix + "/" + namespace + "." + daemonset
}
