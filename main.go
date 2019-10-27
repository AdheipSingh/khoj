package main

import (
	"github.com/khoj/controller"
	"github.com/khoj/utils"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
)

const (
	AllNamespaces = ""
)

func main() {
	clientset := utils.GetKubeClientset()

	// create the deployment watcher
	deployListWatcher := cache.NewListWatchFromClient(clientset.AppsV1().RESTClient(), "deployments", AllNamespaces, fields.Everything())

	// create the workqueue
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())
	// initialize the controller object
	controller := controller.NewController(clientset, queue, deployListWatcher)

	stop := make(chan struct{})
	defer close(stop)
	go controller.Run(1, stop)

	select {}
}
