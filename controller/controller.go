package controller

import (
	"fmt"
	"time"

	"github.com/khoj/utils"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"

	"k8s.io/klog"
)

// Controller Struct
type Controller struct {
	clientset *kubernetes.Clientset
	indexer   cache.Indexer
	queue     workqueue.RateLimitingInterface
	informer  cache.Controller
}

// NewController init
func NewController(clientset *kubernetes.Clientset, queue workqueue.RateLimitingInterface, listWatcher *cache.ListWatch) *Controller {

	indexer, informer := cache.NewIndexerInformer(listWatcher, &appsv1.Deployment{}, 0, cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(obj)
			var event Event

			if err == nil {
				event.key = key
				event.eventType = "create"
				event.deploymentName = obj.(*appsv1.Deployment).Name
				event.deploymentImage = obj.(*appsv1.Deployment).Spec.Template.Spec.Containers[0].Image
				queue.Add(key)
			}
		},
		UpdateFunc: func(oldObj interface{}, newObj interface{}) {
			key, err := cache.MetaNamespaceKeyFunc(newObj)
			var event Event

			if err == nil {
				event.key = key
				event.eventType = "update"
				event.deploymentName = oldObj.(*appsv1.Deployment).Name
				event.deploymentImage = oldObj.(*appsv1.Deployment).Spec.Template.Spec.Containers[0].Image
				queue.Add(key)
			}
		},
		DeleteFunc: func(obj interface{}) {

			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			var event Event
			if err == nil {
				event.key = key
				event.eventType = "delete"
				event.deploymentName = obj.(*appsv1.Deployment).Name
				event.deploymentImage = obj.(*appsv1.Deployment).Spec.Template.Spec.Containers[0].Image
				queue.Add(key)
			}

		},
	}, cache.Indexers{})

	return &Controller{
		clientset: clientset,
		informer:  informer,
		indexer:   indexer,
		queue:     queue,
	}

}

func (c *Controller) Run(threadiness int, stopCh chan struct{}) {
	defer runtime.HandleCrash()

	defer c.queue.ShutDown()
	klog.Info("Starting Deployment controller")

	go c.informer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh, c.informer.HasSynced) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
		return
	}

	for i := 0; i < threadiness; i++ {
		go wait.Until(c.runWorker, time.Second, stopCh)
	}

	<-stopCh
	klog.Info("Stopping Deployment controller")
}

func (c *Controller) runWorker() {
	for c.processNextItem() {
	}
}

func (c *Controller) processNextItem() bool {
	// Get() returns an item with the highest priority.

	key, ok := c.queue.Get()
	if ok {
		return false
	}
	// Every item returned by Get() needs a Done(item) called.
	// Basically tells the queue we are done processing with the key, so
	// know two pods with the same key are never processed in parallel.
	defer c.queue.Done(key)

	err := c.syncHandler(utils.TypeConv(key))
	c.handleErr(err, key)
	return true

}
