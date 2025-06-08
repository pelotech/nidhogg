package utils

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	NodeName     = "nodeName"
	NodeSelector = "nodeSelector"

	Namespace = "namespace"
	Daemonset = "daemonset"

	Daemonset1 = "daemonset1"
	Daemonset2 = "daemonset2"

	TaintNamePrefix = "pelo.tech"
	TaintName       = TaintNamePrefix + "/" + Namespace + "." + Daemonset
)

func BuildPod(podName string, daemonsetName string, conditionType corev1.PodConditionType) corev1.Pod {
	return corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:            podName,
			Namespace:       Namespace,
			OwnerReferences: []metav1.OwnerReference{{Name: daemonsetName}},
		},
		Spec: corev1.PodSpec{
			NodeName: NodeName,
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

func BuildDaemonset(daemonsetName string) appsv1.DaemonSet {
	return appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      daemonsetName,
			Namespace: Namespace,
		},
		Spec: appsv1.DaemonSetSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					NodeSelector: map[string]string{
						NodeSelector: "true",
					},
				},
			},
		},
	}
}

func BuildNode(namespace string, daemonsets []string) corev1.Node {
	return corev1.Node{
		ObjectMeta: metav1.ObjectMeta{
			Name: NodeName,
			Labels: map[string]string{
				NodeSelector: "true",
			},
		},
		Spec: corev1.NodeSpec{
			Taints: BuildActiveTaintsFromDaemonsets(namespace, daemonsets),
		},
	}
}

func BuildActiveTaintsFromDaemonsets(namespace string, daemonsets []string) []corev1.Taint {
	var taints []corev1.Taint
	for _, daemonset := range daemonsets {
		taints = append(taints, BuildActiveTaint(namespace, daemonset))
	}
	return taints
}

func BuildActiveTaint(namespace string, daemonset string) corev1.Taint {
	return corev1.Taint{
		Key:    BuildTaintName(namespace, daemonset),
		Effect: corev1.TaintEffectNoSchedule,
		Value:  string(corev1.ConditionTrue),
	}
}

func BuildTaintName(namespace string, daemonset string) string {
	return TaintNamePrefix + "/" + namespace + "." + daemonset
}
