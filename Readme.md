## Pod Crash Notification Controller
A Kubernetes Controller that watches over pods in all namespaces and generates a notification to the slack channel when any of the pods crashes.

### Getting Started
Before deploying, first update the `<slack webhook url>` in `deploy/deployment.yaml`

#### Run the following command to deploy the resources.
``kubectl create -f deploy/deployment.yaml``

#### Run the following command to create a test pod
``kubectl create -f deploy/example-pod.yaml``

### Usages

Following Args can be passed to the controller deployment - 

```
--namespace string       namespace to be watched by the controller (default "pod-crash")
--slack-webhook string   slack webhook URL to post notifications
```