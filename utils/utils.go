package utils

import (
	"log"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// GetKubeClientset shall return a clientset
func GetKubeClientset() *kubernetes.Clientset {
	var conf *rest.Config

	conf, err := clientcmd.BuildConfigFromFlags("", os.Getenv("HOME")+"/.kube/config")
	if err != nil {
		log.Printf("error in getting Kubeconfig: %v", err)
	}

	cs, err := kubernetes.NewForConfig(conf)
	if err != nil {
		log.Printf("error in getting clientset from Kubeconfig: %v", err)
	}

	return cs
}

// TypeConv shall return a string
func TypeConv(i interface{}) string {
	a := i.(string)
	return a
}
