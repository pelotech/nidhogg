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
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/uswitch/nidhogg/pkg/nidhogg"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var c client.Client

var expectedRequest = reconcile.Request{NamespacedName: types.NamespacedName{Name: "foo"}}

const timeout = time.Second * 5

func TestNodeEnqueueUpdateOnTaintOrLabelChange(t *testing.T) {
	g := gomega.NewWithT(t)
	q := workqueue.NewTypedRateLimitingQueue(workqueue.DefaultTypedControllerRateLimiter[reconcile.Request]())
	defer q.ShutDown()
	e := &nodeEnqueue{}
	old := &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "foo", Labels: map[string]string{"a": "b"}}}

	statusOnly := old.DeepCopy()
	statusOnly.Status.Phase = corev1.NodeRunning
	e.Update(context.TODO(), event.TypedUpdateEvent[*corev1.Node]{ObjectOld: old, ObjectNew: statusOnly}, q)
	g.Expect(q.Len()).To(gomega.Equal(0))

	tainted := old.DeepCopy()
	tainted.Spec.Taints = []corev1.Taint{{Key: "x", Effect: corev1.TaintEffectNoSchedule}}
	e.Update(context.TODO(), event.TypedUpdateEvent[*corev1.Node]{ObjectOld: old, ObjectNew: tainted}, q)
	g.Expect(q.Len()).To(gomega.Equal(1))
}

func TestReconcile(t *testing.T) {
	g := gomega.NewWithT(t)
	instance := &corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "foo"}}

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.
	mgr, err := manager.New(cfg, manager.Options{})
	g.Expect(err).NotTo(gomega.HaveOccurred())
	c = mgr.GetClient()

	handlerConfig := nidhogg.HandlerConfig{}
	_ = handlerConfig.BuildSelectors()

	recFn, requests := SetupTestReconcile(newReconciler(mgr, handlerConfig))
	g.Expect(add(mgr, recFn)).NotTo(gomega.HaveOccurred())

	_, cancel, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		cancel()
		mgrStopped.Wait()
	}()

	// Create the Node object and expect the Reconcile
	err = c.Create(context.TODO(), instance)
	// The instance object may not be a valid object because it might be missing some required fields.
	// Please modify the instance object by adding required fields and then remove the following if statement.
	if apierrors.IsInvalid(err) {
		t.Logf("failed to create object, got an invalid object error: %v", err)
		return
	}
	g.Expect(err).NotTo(gomega.HaveOccurred())
	defer c.Delete(context.TODO(), instance)
	g.Eventually(requests, timeout).Should(gomega.Receive(gomega.Equal(expectedRequest)))

}
