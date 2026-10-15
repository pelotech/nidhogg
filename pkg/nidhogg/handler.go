package nidhogg

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/uswitch/nidhogg/pkg/utils"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/prometheus/client_golang/prometheus"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	corev1helpers "k8s.io/component-helpers/scheduling/corev1"
	"k8s.io/component-helpers/scheduling/corev1/nodeaffinity"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/metrics"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	defaultTaintKeyPrefix      = "nidhogg.uswitch.com"
	taintOperationAdded        = "added"
	taintOperationRemoved      = "removed"
	readySinceAnnotationSuffix = "/ready-since"
)

var (
	taintOperations = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "taint_operations",
		Help: "Total number of added/removed taints operations",
	},
		[]string{
			"operation",
			"taint",
		},
	)
	taintOperationErrors = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "taint_operation_errors",
		Help: "Total number of errors during taint operations",
	},
		[]string{
			"operation",
		},
	)
)

func init() {
	metrics.Registry.MustRegister(
		taintOperations,
		taintOperationErrors,
	)
}

// Handler performs the main business logic of the Wave controller
type Handler struct {
	client.Client
	recorder record.EventRecorder
	config   HandlerConfig
}

// HandlerConfig contains the options for Nidhogg
type HandlerConfig struct {
	TaintNamePrefix            string      `json:"taintNamePrefix,omitempty" yaml:"taintNamePrefix,omitempty"`
	TaintEffect                string      `json:"taintEffect,omitempty" yaml:"taintEffect,omitempty"`
	TaintRemovalDelayInSeconds int         `json:"taintRemovalDelayInSeconds,omitempty" yaml:"taintRemovalDelayInSeconds,omitempty"`
	Daemonsets                 []Daemonset `json:"daemonsets" yaml:"daemonsets"`
	NodeSelector               []string    `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	DaemonsetSelectors         map[Daemonset]labels.Selector
}

func (hc *HandlerConfig) BuildSelectors() error {
	hc.DaemonsetSelectors = make(map[Daemonset]labels.Selector)
	globalSelector := labels.NewSelector()
	for _, rawSelector := range hc.NodeSelector {
		if selector, err := labels.Parse(rawSelector); err != nil {
			return fmt.Errorf("error parsing selector: %v", err)
		} else {
			requirements, _ := selector.Requirements()
			globalSelector = globalSelector.Add(requirements...)
		}
	}
	//Will initialize all daemonsets with the same selector, either representing the NodeSelector config or labels.Nothing if no config was provided for NodeSelector
	for _, daemonset := range hc.Daemonsets {
		hc.DaemonsetSelectors[daemonset] = globalSelector
	}
	return nil
}

// Daemonset contains the name and namespace of a Daemonset
type Daemonset struct {
	Name      string `json:"name" yaml:"name"`
	Namespace string `json:"namespace" yaml:"namespace"`
}

type taintChanges struct {
	taintsAdded   []string
	taintsRemoved []string
	requeueAfter  time.Duration
}

// NewHandler constructs a new instance of Handler
func NewHandler(c client.Client, r record.EventRecorder, conf HandlerConfig) *Handler {
	return &Handler{Client: c, recorder: r, config: conf}
}

// HandleNode works out what taints need to be applied to the nodeName
func (h *Handler) HandleNode(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	log := logf.Log.WithName("nidhogg")

	// Fetch the Node instance
	latestNode := &corev1.Node{}
	err := h.Get(ctx, request.NamespacedName, latestNode)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	updatedNode, taintChanges, err := h.calculateTaints(ctx, latestNode)
	if err != nil {
		taintOperationErrors.WithLabelValues("calculateTaints").Inc()
		return reconcile.Result{}, fmt.Errorf("error calculating taints for node %s: %v", request.Name, err)
	}

	taintLess := true
	for _, taint := range updatedNode.Spec.Taints {
		if strings.HasPrefix(taint.Key, h.getTaintNamePrefix()) {
			taintLess = false
		}
	}

	var readySinceKey = h.getTaintNamePrefix() + readySinceAnnotationSuffix
	var readySinceValue string
	if taintLess {
		readySinceValue = time.Now().UTC().Format(time.RFC3339)
		if updatedNode.Annotations == nil {
			updatedNode.Annotations = map[string]string{
				readySinceKey: readySinceValue,
			}
		} else if _, ok := updatedNode.Annotations[readySinceKey]; !ok {
			updatedNode.Annotations[readySinceKey] = readySinceValue
		} else {
			readySinceValue = updatedNode.Annotations[readySinceKey]
		}
	} else if updatedNode.Annotations != nil {
		readySinceValue = updatedNode.Annotations[readySinceKey]
	}

	if !reflect.DeepEqual(updatedNode, latestNode) {
		log.Info("Updating Node taints", "instance", updatedNode.Name, "taints added", taintChanges.taintsAdded, "taints removed", taintChanges.taintsRemoved, "taintLess", taintLess, "readySinceValue", readySinceValue)

		//err := h.Patch(ctx, updatedNode, client.StrategicMergeFrom(latestNode))
		err := h.Update(ctx, updatedNode)

		if err != nil {
			taintOperationErrors.WithLabelValues("nodeUpdate").Inc()
			return reconcile.Result{}, err
		} else {
			log.Info("Node taints updated.")
		}
		for _, taintAdded := range taintChanges.taintsAdded {
			taintOperations.WithLabelValues(taintOperationAdded, taintAdded).Inc()
		}
		for _, taintRemoved := range taintChanges.taintsRemoved {
			taintOperations.WithLabelValues(taintOperationRemoved, taintRemoved).Inc()
		}

		h.recorder.Eventf(updatedNode, corev1.EventTypeNormal, "TaintsChanged", "Taints added: %s, Taints removed: %s, TaintLess: %v, FirstTimeReady: %q", taintChanges.taintsAdded, taintChanges.taintsRemoved, taintLess, readySinceValue)
	}

	return reconcile.Result{RequeueAfter: taintChanges.requeueAfter}, nil
}

// nodeMatches reports whether the daemonset is expected to run on the node.
func (h *Handler) nodeMatches(ctx context.Context, daemonset Daemonset, node *corev1.Node) (bool, error) {
	if h.config.NodeSelector != nil {
		return h.config.DaemonsetSelectors[daemonset].Matches(labels.Set(node.Labels)), nil
	}
	ds := &appsv1.DaemonSet{}
	if err := h.Get(ctx, types.NamespacedName{Namespace: daemonset.Namespace, Name: daemonset.Name}, ds); err != nil {
		return false, fmt.Errorf("fetching daemonset %s/%s: %w", daemonset.Namespace, daemonset.Name, err)
	}
	pod := &corev1.Pod{Spec: ds.Spec.Template.Spec}
	matches, err := nodeaffinity.GetRequiredNodeAffinity(pod).Match(node)
	if err != nil || !matches {
		return matches, err
	}
	_, untolerated := corev1helpers.FindMatchingUntoleratedTaint(logf.Log, node.Spec.Taints, pod.Spec.Tolerations, h.taintBlocksDaemonset, false)
	return !untolerated, nil
}

// transientTaints come and go on their own or are tolerated by every daemonset pod, so they never mean the daemonset will not run on the node.
var transientTaints = map[string]bool{
	corev1.TaintNodeNotReady:                         true,
	corev1.TaintNodeUnreachable:                      true,
	corev1.TaintNodeUnschedulable:                    true,
	corev1.TaintNodeMemoryPressure:                   true,
	corev1.TaintNodeDiskPressure:                     true,
	corev1.TaintNodePIDPressure:                      true,
	corev1.TaintNodeNetworkUnavailable:               true,
	"node.cloudprovider.kubernetes.io/uninitialized": true,
}

// taintBlocksDaemonset reports whether a taint can keep the daemonset pod off the node for good.
func (h *Handler) taintBlocksDaemonset(taint *corev1.Taint) bool {
	if taint.Effect == corev1.TaintEffectPreferNoSchedule || transientTaints[taint.Key] {
		return false
	}
	return !strings.HasPrefix(taint.Key, h.getTaintNamePrefix())
}

func (h *Handler) calculateTaints(ctx context.Context, instance *corev1.Node) (*corev1.Node, taintChanges, error) {

	nodeCopy := instance.DeepCopy()

	var changes taintChanges

	taintsToRemove := make(map[string]struct{})
	for _, taint := range nodeCopy.Spec.Taints {
		// we could have some older taints from a different configuration file
		// storing them all to reconcile from a previous state
		if strings.HasPrefix(taint.Key, h.getTaintNamePrefix()) {
			taintsToRemove[taint.Key] = struct{}{}
		}
	}
	for _, daemonset := range h.config.Daemonsets {
		matches, err := h.nodeMatches(ctx, daemonset, instance)
		if err != nil {
			return nil, taintChanges{}, err
		}
		if matches {
			taint := fmt.Sprintf("%s/%s.%s", h.getTaintNamePrefix(), daemonset.Namespace, daemonset.Name)
			taintEffect := h.getTaintEffect()
			// Get Pod for nodeName
			pods, err := h.getDaemonsetPods(ctx, instance.Name, daemonset)
			if err != nil {
				return nil, taintChanges{}, fmt.Errorf("error fetching pods: %v", err)
			}

			_, tainted := taintsToRemove[taint]
			if len(pods) == 0 || !utils.AllTrue(pods, podReady) {
				// pod doesn't exist or is not ready
				if tainted {
					// we want to keep this already existing taint on it
					delete(taintsToRemove, taint)
				} else {
					// taint is not already present, adding it
					changes.taintsAdded = append(changes.taintsAdded, taint)
					nodeCopy.Spec.Taints = addTaint(nodeCopy.Spec.Taints, taint, taintEffect)
				}
			} else if wait := h.removalDelayRemaining(pods); tainted && wait > 0 {
				// pod is ready but not for long enough yet, keep the taint and come back later
				delete(taintsToRemove, taint)
				if changes.requeueAfter == 0 || wait < changes.requeueAfter {
					changes.requeueAfter = wait
				}
			}
		}
	}

	for taint := range taintsToRemove {
		nodeCopy.Spec.Taints = removeTaint(nodeCopy.Spec.Taints, taint)
		changes.taintsRemoved = append(changes.taintsRemoved, taint)
	}
	return nodeCopy, changes, nil
}

// removalDelayRemaining is how long until the newest pod has been ready for the configured delay.
func (h *Handler) removalDelayRemaining(pods []*corev1.Pod) time.Duration {
	var readySince time.Time
	for _, pod := range pods {
		for _, condition := range pod.Status.Conditions {
			if condition.Type == corev1.PodReady && condition.LastTransitionTime.After(readySince) {
				readySince = condition.LastTransitionTime.Time
			}
		}
	}
	return time.Duration(h.config.TaintRemovalDelayInSeconds)*time.Second - time.Since(readySince)
}

func (h *Handler) getTaintNamePrefix() string {
	if h.config.TaintNamePrefix != "" {
		return h.config.TaintNamePrefix
	}

	return defaultTaintKeyPrefix
}

func (h *Handler) getTaintEffect() corev1.TaintEffect {
	var effect corev1.TaintEffect

	switch h.config.TaintEffect {
	case "NoSchedule":
		effect = corev1.TaintEffectNoSchedule
	case "PreferNoSchedule":
		effect = corev1.TaintEffectPreferNoSchedule
	case "NoExecute":
		effect = corev1.TaintEffectNoExecute
	default:
		effect = corev1.TaintEffectNoSchedule
	}

	return effect
}

func (h *Handler) getDaemonsetPods(ctx context.Context, nodeName string, ds Daemonset) ([]*corev1.Pod, error) {
	opts := client.InNamespace(ds.Namespace)
	pods := &corev1.PodList{}
	err := h.List(ctx, pods, opts)
	if err != nil {
		return nil, err
	}

	matchingPods := make([]*corev1.Pod, 0)
	for i, pod := range pods.Items {
		for _, owner := range pod.OwnerReferences {
			if owner.Kind == "DaemonSet" && owner.Name == ds.Name && pod.Spec.NodeName == nodeName {
				matchingPods = append(matchingPods, &pods.Items[i])
			}
		}
	}

	return matchingPods, nil
}

func podReady(pod *corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

func addTaint(taints []corev1.Taint, taintName string, taintEffect corev1.TaintEffect) []corev1.Taint {
	return append(taints, corev1.Taint{Key: taintName, Effect: taintEffect})
}

func removeTaint(taints []corev1.Taint, taintName string) []corev1.Taint {
	var newTaints []corev1.Taint

	for _, taint := range taints {
		if taint.Key == taintName {
			continue
		}
		newTaints = append(newTaints, taint)
	}
	return newTaints
}
