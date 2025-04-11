package nidhogg

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/uswitch/nidhogg/pkg/utils"
	"github.com/uswitch/nidhogg/pkg/daemonsets/k8s/apps/v1"
	"github.com/uswitch/nidhogg/pkg/daemonsets/datadoghq.com/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"

	"github.com/prometheus/client_golang/prometheus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
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
	TaintRemovalDelayInSeconds int         `json:"taintRemovalDelayInSeconds,omitempty" yaml:"taintRemovalDelayInSeconds,omitempty"`
	Daemonsets                 []DaemonsetConfig `json:"daemonsets" yaml:"daemonsets"`
	NodeSelector               []string    `json:"nodeSelector,omitempty" yaml:"nodeSelector,omitempty"`
	DaemonsetSelectors         map[DaemonsetConfig]labels.Selector
}

func (hc *HandlerConfig) BuildSelectors() error {
	hc.DaemonsetSelectors = make(map[DaemonsetConfig]labels.Selector)
	globalSelector := labels.Nothing()
	for _, rawSelector := range hc.NodeSelector {
		if selector, err := labels.Parse(rawSelector); err != nil {
			return fmt.Errorf("error parsing selector: %v", err)
		} else {
			requirements, _ := selector.Requirements()
			globalSelector = labels.NewSelector().Add(requirements...)
		}
	}
	//Will initialize all daemonsets with the same selector, either representing the NodeSelector config or labels.Nothing if no config was provided for NodeSelector
	for _, daemonset := range hc.Daemonsets {
		hc.DaemonsetSelectors[daemonset] = globalSelector
	}
	return nil
}

// DaemonsetConfig contains the namel namespace and api/version of a Daemonset
type DaemonsetConfig struct {
	Name       string `json:"name" yaml:"name"`
	Namespace  string `json:"namespace" yaml:"namespace"`
	ApiVersion string `json:"apiVersion" yaml:"apiVersion"`
}

type Daemonset interface {
	GetSelector(ctx context.Context)(labels.Selector, error)
	GetPods(ctx context.Context, nodeName string) ([]*corev1.Pod, error)
}


type taintChanges struct {
	taintsAdded   []string
	taintsRemoved []string
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
		return reconcile.Result{}, fmt.Errorf("error calculating taints for nodeName: %v", err)
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
		readySinceValue = time.Now().Format("2006-01-02T15:04:05Z")
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

		// this is a hack to make the event work on a non-namespaced object
		updatedNode.UID = types.UID(updatedNode.Name)

		h.recorder.Eventf(updatedNode, corev1.EventTypeNormal, "TaintsChanged", "Taints added: %s, Taints removed: %s, TaintLess: %v, FirstTimeReady: %q", taintChanges.taintsAdded, taintChanges.taintsRemoved, taintLess, readySinceValue)
	}

	return reconcile.Result{}, nil
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

		var ds Daemonset
		switch daemonset.ApiVersion {
		case "datadoghq.com/v1alpha1":
			ds = &v1alpha1.ExtendedDaemonset{
				Name: daemonset.Name,
				Namespace: daemonset.Namespace,
				List: h.List,
				Get: h.Get,
			}
		default:
			ds = &v1.Daemonset{
				Name: daemonset.Name,
				Namespace: daemonset.Namespace,
				List: h.List,
				Get: h.Get,
			}
		}

		//If NodeSelector was not provided upfront through config
		if h.config.NodeSelector == nil {
			//Will try to get selectors from daemonset directly
			selector, err := ds.GetSelector(ctx)
			if err != nil {
				logf.Log.Info(fmt.Sprintf("Could not fetch selector from daemonset %s in namespace %s", daemonset.Name, daemonset.Namespace))
			} else {
				//Override existing daemonset selector with the one freshly retrieved from the daemonset
				h.config.DaemonsetSelectors[daemonset] = selector
			}
		}

		//make sure daemonset selector matches node selector
		if h.config.DaemonsetSelectors[daemonset].Matches(labels.Set(instance.Labels)) {
			taint := fmt.Sprintf("%s/%s.%s", h.getTaintNamePrefix(), daemonset.Namespace, daemonset.Name)
			// Get Pod for nodeName
			pods, err := ds.GetPods(ctx, instance.Name)
			if err != nil {
				return nil, taintChanges{}, fmt.Errorf("error fetching pods: %v", err)
			}
			if len(pods) == 0 || (len(pods) > 0 && !utils.AllTrue(pods, func(pod *corev1.Pod) bool { return podReady(pod) })) {
				// pod doesn't exist or is not ready
				_, ok := taintsToRemove[taint]
				if ok {
					// we want to keep this already existing taint on it
					delete(taintsToRemove, taint)
				} else {
					// taint is not already present, adding it
					changes.taintsAdded = append(changes.taintsAdded, taint)
					nodeCopy.Spec.Taints = addTaint(nodeCopy.Spec.Taints, taint)
				}
			}
		}
	}

	for taint := range taintsToRemove {
		h.applyTaintRemovalDelay()
		nodeCopy.Spec.Taints = removeTaint(nodeCopy.Spec.Taints, taint)
		changes.taintsRemoved = append(changes.taintsRemoved, taint)
	}
	return nodeCopy, changes, nil
}

func (h *Handler) applyTaintRemovalDelay() {
	if h.config.TaintRemovalDelayInSeconds == 0 {
		return
	}
	logf.Log.Info("Daemonset is running, a delay has been set before removing taint.", "delay", h.config.TaintRemovalDelayInSeconds)
	time.Sleep(time.Duration(h.config.TaintRemovalDelayInSeconds) * time.Second)
}

func (h *Handler) getTaintNamePrefix() string {
	if h.config.TaintNamePrefix != "" {
		return h.config.TaintNamePrefix
	}

	return defaultTaintKeyPrefix
}

func podReady(pod *corev1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type == corev1.PodReady {
			return condition.Status == corev1.ConditionTrue
		}
	}
	return false
}

func addTaint(taints []corev1.Taint, taintName string) []corev1.Taint {
	return append(taints, corev1.Taint{Key: taintName, Effect: corev1.TaintEffectNoSchedule})
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
