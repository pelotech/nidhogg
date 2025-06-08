package v1

import (
	"context"
	"fmt"


	"k8s.io/apimachinery/pkg/labels"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ObjectKey = types.NamespacedName
type ObjectList = client.ObjectList
type Object = client.Object
type GetOption = client.GetOption
type ListOption = client.ListOption
type Context = context.Context

type Daemonset struct {
	Name string
	Namespace string
	List func(ctx Context, list ObjectList, opts ...ListOption) error
	Get func(ctx Context, key ObjectKey, obj Object, opts ...GetOption) error
}

func (d *Daemonset) GetSelector(ctx Context)(labels.Selector, error) {
	ds := &appsv1.DaemonSet{}

	err := d.Get(ctx, types.NamespacedName{Namespace: d.Namespace, Name: d.Name}, ds)
	if err != nil {
		logf.Log.Info(fmt.Sprintf("Could not fetch daemonset %s from namespace %s", d.Name, d.Namespace))
		return nil, err
	}
	selector := labels.SelectorFromSet(ds.Spec.Template.Spec.NodeSelector)
	return selector, nil
}

func (d *Daemonset) GetPods(ctx Context, nodeName string) ([]*corev1.Pod, error) {
	opts := client.InNamespace(d.Namespace)
	pods := &corev1.PodList{}

	err := d.List(ctx, pods, opts)
	if err != nil {
		return nil, err
	}

	matchingPods := make([]*corev1.Pod, 0)
	for i, pod := range pods.Items {
		for _, owner := range pod.OwnerReferences {
			if owner.Name == d.Name && pod.Spec.NodeName == nodeName {
				matchingPods = append(matchingPods, &pods.Items[i])
			}
		}
	}
	return matchingPods, nil
}
