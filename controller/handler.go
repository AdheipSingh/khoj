package controller

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog"
)

func (c *Controller) syncHandler(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	deploymentName := obj.(*appsv1.Deployment).Name
	deploymentImage := obj.(*appsv1.Deployment).Spec.Template.Spec.Containers[0].Image

	if err != nil {
		klog.Errorf("Error fetching obj %s from store failed due to %v\n", key, err)
		return err
	}

	if !exists {
		fmt.Printf("Deployment %s does not exist anymore\n", key)
	} else {

		fmt.Printf("Deployment Name [%s] : Image Name [%s] \n", deploymentName, deploymentImage)
	}
	return nil

}

func (c *Controller) handleErr(err error, key interface{}) {
	if err == nil {

		c.queue.Forget(key)
		return
	}

	if c.queue.NumRequeues(key) < 5 {
		klog.Infof("Error syncing pod %v: %v", key, err)
		c.queue.AddRateLimited(key)
		return
	}

	c.queue.Forget(key)

	runtime.HandleError(err)
	klog.Infof("Dropping Deployemnt %q out of the queue: %v", key, err)
}
