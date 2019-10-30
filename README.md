# khoj
- A kubernetes controller watching for images running in deployments.

### Getting Started
- Assume you have a working kubeconfig present on $HOME/.kube/config
- It shall watch images running on all namespaces in cluster and if a new image is added it prints out the image.
```
- go run main.go
```
