package kube

import (
	"context"
	"errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
)

// Connect make connection
func Connect() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := makeConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// GetPodsList Get list of pods
func ListPods(clientset *kubernetes.Clientset, namespace string) (*v1.PodList, error) {
	return clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}

// DeletePod Delete pod
func DeletePod(clientset *kubernetes.Clientset, namespace string, name string) error {
	return clientset.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{})
}

// Detect cluster configuration
// Using the cluster configuration if running inside a cluster,
// otherwise using the ~/.kube/config configuration file.
func makeConfig() (*rest.Config, error) {
	home, exists := os.LookupEnv("HOME")
	if !exists {
		home = "/root"
	}

	_, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST")

	configPath := filepath.Join(home, ".kube", "config")

	if _, err := os.Stat(configPath); err == nil {
		config, err := clientcmd.BuildConfigFromFlags("", configPath)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
	if inCluster {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		return config, nil
	}

	return nil, errors.New("K8S config not found")
}
