cli to clone Kubernetes deployment from one cluster to another

Usage:

`cluster_clone clone <sourceKubeConfig>  <destKubeConfig>`


Lang/tools used:
- Go lang
- Cobra CLI

Working:
1) It uses `kubectl get deployment -o json --kubeconfig=kubeConfigFile> source_resources.json"` command to get all deployment from given source kubeConfig writes resources to a file

2) Next to fetch kubectl.kubernetes.io/last-applied-configuration and write to a json file

3) Then `kubectl apply -f source_resources.json --kubeconfig=kubeConfigFile` is run to apply resources from json to destination
