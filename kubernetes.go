package main

import (
	"context"
	"flag"
	"log"
	"math/rand"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type PodInfo struct {
	Name      string
	Namespace string
}

var clientset *kubernetes.Clientset

func initKubeClient() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Printf("Error building kubeconfig: %s\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("Error building kube config: %s\n", err.Error())
		}
	}

	clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating clientset: %s\n", err.Error())
	}
}

func getRandomPodInfo() PodInfo {
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		log.Fatalf("Error listing pods: %s\n", err.Error())
		return PodInfo{}
	}

	var nonCriticalPods []PodInfo
	for _, pod := range pods.Items {
		if pod.Namespace != "kube-system" && !isCriticalPod(pod) {
			nonCriticalPods = append(nonCriticalPods, PodInfo{Name: pod.Name, Namespace: pod.Namespace})
		}
	}

	if len(nonCriticalPods) == 0 {
		log.Println("No eligible pods found to delete.")
		return PodInfo{}
	}

	// Randomly select a pod from the non-critical list
	randIndex := rand.Intn(len(nonCriticalPods))
	return nonCriticalPods[randIndex]
}

func isCriticalPod(pod v1.Pod) bool {
	_, isCritical := pod.Annotations["scheduler.alpha.kubernetes.io/critical-pod"]
	return isCritical
}

func deletePod(podInfo PodInfo) {
	err := clientset.CoreV1().Pods(podInfo.Namespace).Delete(context.TODO(), podInfo.Name, metav1.DeleteOptions{})
	if err != nil {
		log.Printf("Error deleting pod %s in namespace %s: %s\n", podInfo.Name, podInfo.Namespace, err.Error())
	} else {
		log.Printf("Pod deleted: %s in namespace %s\n", podInfo.Name, podInfo.Namespace)
	}
}
