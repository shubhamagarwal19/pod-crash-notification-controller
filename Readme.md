## Pod Crash Notification Controller
A Kubernetes Controller that watches over pods in all namespaces and generates a notification to the slack channel when any of the pods crashes.

### Getting Started
Before deploying, first update the `<slack webhook url>` in `deploy/deployment.yaml`

#### Run the following command to deploy the resources.
``kubectl create -f deploy/deployment.yaml``

#### Run the following command to create a test pod
``kubectl create -f deploy/example-pod.yaml``