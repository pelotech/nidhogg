/*

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package node

import (
	"context"
	"github.com/uswitch/nidhogg/pkg/nidhogg"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new Node Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager, cfg nidhogg.HandlerConfig) error {
	return add(mgr, newReconciler(mgr, cfg))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager, cfg nidhogg.HandlerConfig) reconcile.Reconciler {
	eventRecorder := mgr.GetEventRecorderFor("nidhogg")
	reconcilerHandler := nidhogg.NewHandler(mgr.GetClient(), eventRecorder, cfg)
	return &ReconcileNode{reconcilerHandler, mgr.GetScheme()}
}

var _ handler.EventHandler = &nodeEnqueue{}

type nodeEnqueue struct{}

// Update implements the interface
func (e *nodeEnqueue) Update(_ context.Context, _ event.UpdateEvent, _ workqueue.RateLimitingInterface) {
}

// Delete implements the interface
func (e *nodeEnqueue) Delete(_ context.Context, _ event.DeleteEvent, _ workqueue.RateLimitingInterface) {
}

// Generic implements the interface
func (e *nodeEnqueue) Generic(_ context.Context, _ event.GenericEvent, _ workqueue.RateLimitingInterface) {
}

// Create adds the node to the queue, the node is created as NotReady and without daemonset pods
func (e *nodeEnqueue) Create(_ context.Context, evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	if evt.Object == nil {
		return
	}
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name: evt.Object.GetName(),
	}})
}

var _ handler.EventHandler = &podEnqueue{}

type podEnqueue struct{}

// Generic implements the interface
func (e *podEnqueue) Generic(_ context.Context, _ event.GenericEvent, _ workqueue.RateLimitingInterface) {
}

// canAddToQueue check if the Pod is associated to a node and is a daemonset pod
func (e *podEnqueue) canAddToQueue(pod *corev1.Pod) bool {
	if pod.Spec.NodeName == "" {
		return false
	}
	owner := v1.GetControllerOf(pod)
	if owner == nil {
		return false
	}
	return owner.Kind == "DaemonSet"
}

// Create adds the node of the daemonset pod to the queue
func (e *podEnqueue) Create(_ context.Context, evt event.CreateEvent, q workqueue.RateLimitingInterface) {
	pod, ok := evt.Object.(*corev1.Pod)
	if !ok {
		return
	}
	if !e.canAddToQueue(pod) {
		return
	}
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name: pod.Spec.NodeName,
	}})

}

// Update adds the node of the updated daemonset pod to the queue
func (e *podEnqueue) Update(_ context.Context, evt event.UpdateEvent, q workqueue.RateLimitingInterface) {
	pod, ok := evt.ObjectNew.(*corev1.Pod)
	if !ok {
		return
	}
	if !e.canAddToQueue(pod) {
		return
	}
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name: pod.Spec.NodeName,
	}})
}

// Delete adds the node of the deleted daemonset pod to the queue
func (e *podEnqueue) Delete(_ context.Context, evt event.DeleteEvent, q workqueue.RateLimitingInterface) {
	pod, ok := evt.Object.(*corev1.Pod)
	if !ok {
		return
	}
	if !e.canAddToQueue(pod) {
		return
	}
	q.Add(reconcile.Request{NamespacedName: types.NamespacedName{
		Name: pod.Spec.NodeName,
	}})
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("node-controller", mgr, controller.Options{
		Reconciler:              r,
		MaxConcurrentReconciles: 1,
	})
	if err != nil {
		return err
	}

	// Watch for changes to Node
	err = c.Watch(source.Kind(mgr.GetCache(), &corev1.Node{}), &nodeEnqueue{})
	if err != nil {
		return err
	}

	err = c.Watch(source.Kind(mgr.GetCache(), &corev1.Pod{}), &podEnqueue{})
	if err != nil {
		return err
	}

	return nil
}

// ReconcileNode reconciles a Node object
var _ reconcile.Reconciler = &ReconcileNode{}

type ReconcileNode struct {
	handler *nidhogg.Handler
	scheme  *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Node object and makes changes based on the state read
// and what is in the Node.Spec
// +kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=nodes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch
// +kubebuilder:rbac:groups=,resources=events,verbs=create;update;patch
func (r *ReconcileNode) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	return r.handler.HandleNode(ctx, request)
}
