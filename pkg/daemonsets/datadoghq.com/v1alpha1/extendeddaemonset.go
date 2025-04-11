package v1alpha1

import (
	"context"
	"fmt"
	"slices"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
    ddv1alpha1 "github.com/DataDog/extendeddaemonset/api/v1alpha1"
)

type ObjectKey = types.NamespacedName
type ObjectList = client.ObjectList
type Object = client.Object
type GetOption = client.GetOption
type ListOption = client.ListOption
type Context = context.Context

type ExtendedDaemonset struct {
	Name string
	Namespace string
	List func(ctx Context, list ObjectList, opts ...ListOption) error
	Get func(ctx Context, key ObjectKey, obj Object, opts ...GetOption) error
}

func (d *ExtendedDaemonset) GetSelector(ctx Context)(labels.Selector, error) {
	ds := &ddv1alpha1.ExtendedDaemonSet{}
	err := d.Get(ctx, types.NamespacedName{Namespace: d.Namespace, Name: d.Name}, ds)
	if err != nil {
		logf.Log.Info(fmt.Sprintf("Could not fetch extendeddaemonset %s from namespace %s", d.Name, d.Namespace))
		return nil, err
	}
	selector := labels.SelectorFromSet(ds.Spec.Template.Spec.NodeSelector)

	return selector, nil
}

func (d *ExtendedDaemonset) getOwnerNames(ctx Context) (*[]string, error) {
	ds := &ddv1alpha1.ExtendedDaemonSet{}
	err := d.Get(ctx, types.NamespacedName{Namespace: d.Namespace, Name: d.Name}, ds)
	if err != nil {
		logf.Log.Info(fmt.Sprintf("Could not fetch extendeddaemonset %s from namespace %s", d.Name, d.Namespace))
		return nil, err
	}

	// For extendeddaemonset there can be two active owners at the same time:
	// active and canary
	owner := make([]string,0)

	owner = append(owner, ds.Status.ActiveReplicaSet)
	if ds.Status.Canary != nil {
		owner = append(owner, ds.Status.Canary.ReplicaSet)
	}

	return &owner, nil
}

func is_owned(podOwnerReferences []metav1.OwnerReference,  replicasetNames *[]string) bool {
	for _, ref := range podOwnerReferences {
		if slices.Contains(*replicasetNames, ref.Name) {
			return true
		}
	}
	return false
}

func (d *ExtendedDaemonset) GetPods(ctx context.Context, nodeName string) ([]*corev1.Pod, error) {
	opts := client.InNamespace(d.Namespace)
	pods := &corev1.PodList{}
	err := d.List(ctx, pods, opts)
	if err != nil {
		return nil, err
	}

	owners, err := d.getOwnerNames(ctx)
	if err != nil {
		logf.Log.Error(err, "Failed to retrieve daemonset owner name")
		return nil, err
	}

	matchingPods := make([]*corev1.Pod, 0)
	for i, pod := range pods.Items {
		if is_owned(pod.OwnerReferences, owners) && pod.Spec.NodeName == nodeName {
			matchingPods = append(matchingPods, &pods.Items[i])
		}
	}
	return matchingPods, nil
}
