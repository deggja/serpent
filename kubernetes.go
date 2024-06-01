package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var podInfoQueue = make(chan PodInfo, 100)

func fetchPods() {
    for {
        podInfo := getRandomPodInfo()
        if podInfo.Name != "" && podInfo.Namespace != "" {
            podInfoQueue <- podInfo
        }
        time.Sleep(1 * time.Second)
    }
}

type Config struct {
    ResourceTypes []string `json:"resource_types"`
    Namespaces    NamespacesConfig `json:"namespaces"`
}

type NamespacesConfig struct {
    Include []string `json:"include"`
    Exclude []string `json:"exclude"`
}

var defaultConfig = Config{
    ResourceTypes: []string{"pods"},
    Namespaces: NamespacesConfig{
        Include: []string{},
        Exclude: []string{"kube-system"},
    },
}

var gameConfig Config

func setDefaultConfig() {
    gameConfig = defaultConfig
}

func loadConfigFromFile(filename string) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    err = json.Unmarshal(data, &gameConfig)
    if err != nil {
        return err
    }
    return nil
}

func getAllNamespaces() ([]string, error) {
    allNamespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
    if err != nil {
        return nil, err
    }

    var filteredNamespaces []string
    for _, ns := range allNamespaces.Items {
        if (len(gameConfig.Namespaces.Include) == 0 || contains(gameConfig.Namespaces.Include, ns.Name)) &&
           !contains(gameConfig.Namespaces.Exclude, ns.Name) {
            filteredNamespaces = append(filteredNamespaces, ns.Name)
        }
    }
    return filteredNamespaces, nil
}


func contains(slice []string, str string) bool {
    for _, v := range slice {
        if v == str {
            return true
        }
    }
    return false
}

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
    // Retrieve all namespaces and filter them based on gameConfig
    filteredNamespaces, err := getAllNamespaces()
    if err != nil {
        log.Printf("Error fetching namespaces: %s\n", err)
        return PodInfo{}
    }

    var nonCriticalPods []PodInfo
    for _, ns := range filteredNamespaces {
        pods, err := clientset.CoreV1().Pods(ns).List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            log.Printf("Error listing pods in namespace %s: %s\n", ns, err)
            continue
        }

        for _, pod := range pods.Items {
            if !isCriticalPod(pod) {
                nonCriticalPods = append(nonCriticalPods, PodInfo{Name: pod.Name, Namespace: pod.Namespace})
            }
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
