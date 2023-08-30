package clients

import (
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func GetKubernetesClient(env string) (*kubernetes.Clientset, error) {
	if env == "dev" {
		kubeconfigPath := "kubeconfig.yaml"
		if kubeconfigPath == "" {
			return nil, fmt.Errorf("KUBECONFIG environment variable not set")
		}

		config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			return nil, err
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}

		return clientset, nil
	} else if env == "prod" {
		config, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

		clientset, err := kubernetes.NewForConfig(config)
		if err != nil {
			return nil, err
		}

		return clientset, nil
	} else {
		return nil, fmt.Errorf("Invalid environment: %s", env)
	}
}

func GetKubeconfigPath(env string) (string, error) {
	var kubeconfigPath string

	switch env {
	case "dev":
		kubeconfigPath = "./config/kubeconfig-dev.yaml" // Set your development kubeconfig path
	case "prod":
		kubeconfigPath = "./config/kubeconfig-prod.yaml" // Set your production kubeconfig path
	default:
		return "", fmt.Errorf("Invalid environment: %s", env)
	}

	return kubeconfigPath, nil
}

// func GetKubernetesClient() (*kubernetes.Clientset, error) {
// 	kubeconfigPath := "kubeconfig.yaml"
// 	if kubeconfigPath == "" {
// 		return nil, fmt.Errorf("KUBECONFIG environment variable not set")
// 	}

// 	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
// 	if err != nil {
// 		return nil, err
// 	}

// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return clientset, nil
// }

// func GetInClusterKubernetesClient() (*kubernetes.Clientset, error) {
// 	config, err := rest.InClusterConfig()
// 	if err != nil {
// 		return nil, err
// 	}

// 	clientset, err := kubernetes.NewForConfig(config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return clientset, nil
// }
