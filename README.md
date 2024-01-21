# KUBERNETES FAILED POD CLEANER

Finds all pods with Failed status in all namespaces and removes them.
All the necessary information and the reason why it switched to this status are entered into the log `Status.Message`

## Output
```json
{"level":"info","msg":"pod cleaner is running.","time":"2024-01-21T16:05:50Z"}
{"level":"info","msg":"removed","pod_message":"The node was low on resource: memory. Container sidekiq was using 1001776Ki, which exceeds its request of 500M. ","pod_name":"sidekiq-65f5b978cf-z6xh7","pod_namespace":"sidekiq","pod_node":"ip-10-20-20-18.compute.internal","pod_reason":"Evicted","pod_status":"Failed","time":"2024-01-21T17:05:50Z"}
{"level":"info","msg":"removed","pod_message":"The node was low on resource: memory. Container app was using 1094660Ki, which exceeds its request of 1G. Container watcher was using 4680Ki, which exceeds its request of 0. ","pod_name":"server-deployment-56db4cddc4-s2mk4","pod_namespace":"app","pod_node":"ip-10-20-21-177.compute.internal","pod_reason":"Evicted","pod_status":"Failed","time":"2024-01-21T17:05:50Z"}
...
```

## Images
GitHub registry
```shell
docker pull ghcr.io/yousysadmin/failed-pod-cleaner:latest
```
DockerHub registry
```shell
docker pull yousysadmin/failed-pod-cleaner:latest
```

## Tags
`edge` - latest master branch build  
`latest` - latest release  
`vx.x.x` - latest release by version  

## Usage

### Local
```shell
docker run --rm -v $HOME/.kube/config:/root/.kube/config yousysadmin/failed-pod-cleaner:latest
```

### Cluster
```shell
# Simple install in the `default` namespace
kubectl apply -f https://github.com/YouSysAdmin/failed-pod-cleaner/examples/dist.yaml
# Or
git clone https://github.com/YouSysAdmin/failed-pod-cleaner.git
kubectl apply -k failed-pod-cleaner/examples/kustomize
```

## Build
```shell
# Build local image with tag failed-pod-cleaner:local
docker buildx bake

# Build multi-platform image
docker buildx bake image-all
```