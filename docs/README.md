![build](https://github.com/pelotech/nidhogg/actions/workflows/build.yaml/badge.svg)
![release](https://github.com/pelotech/nidhogg/actions/workflows/release.yaml/badge.svg)
![publish-chart](https://github.com/pelotech/nidhogg/actions/workflows/publish-chart.yaml/badge.svg)

# Nidhogg

Nidhogg is a controller that taints nodes based on whether a Pod from a specific Daemonset is running on them.

Sometimes you have a Daemonset that is so important that you don't want other pods to run on your node until that Daemonset is up and running on the node. Nidhogg solves this problem by tainting the node until your Daemonset pod is ready, preventing pods that don't tolerate the taint from scheduling there.

Nidhogg annotate the node when all the required taints are removed: `nidhogg.uswitch.com/ready-since: 2006-01-02T15:04:05Z`

Nidhogg was built using [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder)

## Usage

Nidhogg requires a yaml/json config file to tell it what Daemonsets to watch and what nodes to act on.

| Attribute name | Required/Optional | Description |
| :--- | :--- | :--- |
| `daemonsets` | Required | Array of Daemonsets to watch, each containing two required fields `name` and `namespace` and optional `apiVersion` set to `datadoghq.com/v1alpha1` to manage [DataDog](https://github.com/DataDog/extendeddaemonset/)|
| `nodeSelector` | Optional | Map of keys/values corresponding to node labels, will default to get selectors from daemonsets directly if not provided |
| `taintNamePrefix` | Optional | Prefix of the taint name, defaults to `nidhogg.uswitch.com` if not specified |
| `taintRemovalDelayInSeconds` | Optional | Delay to apply before removing taint on the node when ready, defaults to 0 if not specified |

Nodes are tainted with a taint that follows the format of `taintNamePrefix/namespace.name:NoSchedule`

Example:

YAML:
```yaml
daemonsets:
  - name: kiam
    namespace: kube-system
  - name: datadog-agent
    namespace: datadog-operator-system
    apiVersion: datadoghq.com/v1alpha1
nodeSelector:
  - "node-role.kubernetes.io/node"
  - "!node-role.kubernetes.io/master"
  - "aws.amazon.com/ec2.asg.name in (standard, special)"
taintNamePrefix: "nidhogg.uswitch.com"
taintRemovalDelayInSeconds: 10
```
JSON:
```json
{
  "daemonsets": [
    {
      "name": "kiam",
      "namespace": "kube-system"
    },
    {
      "name": "datadog-agent",
      "namespace": "datadog-operator-system",
      "apiVersion": "datadoghq.com/v1alpha1"
    }
  ],
  "nodeSelector": [
    "node-role.kubernetes.io/node",
    "!node-role.kubernetes.io/master",
    "aws.amazon.com/ec2.asg.name in (standard, special)"
  ],
  "taintNamePrefix": "nidhogg.uswitch.com",
  "taintRemovalDelayInSeconds": 10
}
```
This example will select any nodes in AWS ASGs named "standard" or "special" that have the label `node-role.kubernetes.io/node` present, and no nodes with label `node-role.kubernetes.io/master`

If the matching nodes do not have a running and ready pod from the `kiam` daemonset in the `kube-system` namespace, it will add a taint of `nidhogg.uswitch.com/kube-system.kiam:NoSchedule` until there is a ready kiam pod on the node.

Whenever the pod becomes ready, a delay of 10s will be applied before removing the taint.

If you want pods to be able to run on the nidhogg tainted nodes you can add a toleration:

```yaml
spec:
  tolerations:
  - key: nidhogg.uswitch.com/kube-system.kiam
    operator: "Exists"
    effect: NoSchedule
```

## Deploying
Docker images can be found at https://ghcr.io/pelotech/nidhogg

### helm

We publish a helm chart to https://ghcr.io/pelotech/charts/nidhogg

### kustomize

[Kustomize](https://github.com/kubernetes-sigs/kustomize) manifests can be found  [here](/kustomize) to quickly deploy this to a cluster.

## Flags
```
-config-file string
    Path to config file (default "config.json")
-kubeconfig string
    Paths to a kubeconfig. Only required if out-of-cluster.
-leader-configmap string
    Name of configmap to use for leader election
-leader-election
    enable leader election
-leader-namespace string
    Namespace where leader configmap located
-master string
    The address of the Kubernetes API server. Overrides any value in kubeconfig. Only required if out-of-cluster.
-metrics-addr string
    The address the metric endpoint binds to. (default ":8080")
-kube-api-qps float
    QPS rate for throttling requests sent to the Kubernetes API server (default 20)
-kube-api-burst int
    Maximum burst for throttling requests sent to the Kubernetes API server (default 30)
-disable-compression bool
    Disable response compression for k8s restAPI in client-go (default true)
```

