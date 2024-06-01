package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
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

var resourceInfoQueue = make(chan ResourceInfo, 100)

func fetchResources() {
    for {
        resourceInfo, err := getRandomResourceInfo()
        if err != nil {
            log.Printf("Error fetching resource info: %s\n", err)
        } else {
            resourceInfoQueue <- resourceInfo
        }
        time.Sleep(1 * time.Second)
    }
}

type KubernetesResource interface {
    List(ctx context.Context, namespace string, opts metav1.ListOptions) ([]ResourceInfo, error)
    Delete(ctx context.Context, namespace, name string) error
}

type PodResource struct {
    clientset *kubernetes.Clientset
}

type ResourceInfo struct {
    Name      string
    Namespace string
    Type      string
}

func (p *PodResource) List(ctx context.Context, namespace string, opts metav1.ListOptions) ([]ResourceInfo, error) {
    pods, err := p.clientset.CoreV1().Pods(namespace).List(ctx, opts)
    if err != nil {
        return nil, err
    }
    var results []ResourceInfo
    for _, pod := range pods.Items {
        if !isCriticalPod(pod) {
            results = append(results, ResourceInfo{Name: pod.Name, Namespace: pod.Namespace, Type: "pod"})
        }
    }
    return results, nil
}

func (p *PodResource) Delete(ctx context.Context, namespace, name string) error {
    return p.clientset.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
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

func getResourceHandler(resourceType string) (KubernetesResource, error) {
    switch resourceType {
    case "pods":
        return &PodResource{clientset: clientset}, nil
    // Add cases for other resource types here
    default:
        return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
    }
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

func getRandomResourceInfo() (ResourceInfo, error) {
    if len(gameConfig.ResourceTypes) == 0 {
        return ResourceInfo{}, fmt.Errorf("no resource types configured")
    }

    var allResources []ResourceInfo
    for _, resourceType := range gameConfig.ResourceTypes {
        handler, err := getResourceHandler(resourceType)
        if err != nil {
            log.Printf("Error getting handler for resource type %s: %s\n", resourceType, err)
            continue
        }

        filteredNamespaces, err := getAllNamespaces()
        if err != nil {
            log.Printf("Error fetching namespaces: %s\n", err)
            continue
        }

        for _, ns := range filteredNamespaces {
            resources, err := handler.List(context.TODO(), ns, metav1.ListOptions{})
            if err != nil {
                log.Printf("Error listing resources in namespace %s: %s\n", ns, err)
                continue
            }
            allResources = append(allResources, resources...)
        }
    }

    if len(allResources) == 0 {
        log.Println("No eligible resources found.")
        return ResourceInfo{}, fmt.Errorf("no eligible resources found")
    }

    randIndex := rand.Intn(len(allResources))
    return allResources[randIndex], nil
}


func isCriticalPod(pod v1.Pod) bool {
	_, isCritical := pod.Annotations["scheduler.alpha.kubernetes.io/critical-pod"]
	return isCritical
}

func deleteResource(resourceInfo ResourceInfo) {
    var err error
    switch resourceInfo.Type {
    case "pod":
        err = clientset.CoreV1().Pods(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
    // Add cases for other resource types here
    default:
        log.Printf("Unsupported resource type: %s", resourceInfo.Type)
        return
    }

    if err != nil {
        log.Printf("Error deleting %s %s in namespace %s: %s\n", resourceInfo.Type, resourceInfo.Name, resourceInfo.Namespace, err.Error())
    } else {
        log.Printf("%s deleted: %s in namespace %s\n", resourceInfo.Type, resourceInfo.Name, resourceInfo.Namespace)
    }
}
