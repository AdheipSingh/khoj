package controller

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog"
)

func (c *Controller) syncHandler(key string) error {
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		klog.Errorf("Error fetching obj %s from store failed due to %v\n", key, err)
		return err
	}
	if !exists {
		fmt.Printf("Deployment [%s]  does not exist \n", key)
	} else {

		fmt.Printf("Deployment name [%s] ||  Image name [%s] \n", obj.(*appsv1.Deployment).Name, obj.(*appsv1.Deployment).Spec.Template.Spec.Containers[0].Image)
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
