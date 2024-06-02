package main

import (
	"context"
	"encoding/json"
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

// Resource types supported
// Pods
// ReplicaSets
// Deployments
// StatefulSets
// Services
// DaemonSets
// Secrets
// ConfigMaps
// Jobs
// CronJobs
// Ingresses

type PodResource struct {
    clientset *kubernetes.Clientset
}

type ReplicaSetResource struct {
    clientset *kubernetes.Clientset
}

type DeploymentResource struct {
    clientset *kubernetes.Clientset
}

type StatefulSetResource struct {
    clientset *kubernetes.Clientset
}

type ServiceResource struct {
    clientset *kubernetes.Clientset
}

type DaemonSetResource struct {
    clientset *kubernetes.Clientset
}

type SecretResource struct {
    clientset *kubernetes.Clientset
}

type ConfigMapResource struct {
    clientset *kubernetes.Clientset
}

type JobResource struct {
    clientset *kubernetes.Clientset
}

type CronJobResource struct {
    clientset *kubernetes.Clientset
}

type IngressResource struct {
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

func (r *ReplicaSetResource) List(ctx context.Context, namespace string, opts metav1.ListOptions) ([]ResourceInfo, error) {
    replicasets, err := r.clientset.AppsV1().ReplicaSets(namespace).List(ctx, opts)
    if err != nil {
        return nil, err
    }
    var results []ResourceInfo
    for _, rs := range replicasets.Items {
        results = append(results, ResourceInfo{Name: rs.Name, Namespace: rs.Namespace, Type: "replicaset"})
    }
    return results, nil
}

func (d *DeploymentResource) List(ctx context.Context, namespace string, opts metav1.ListOptions) ([]ResourceInfo, error) {
    deployments, err := d.clientset.AppsV1().Deployments(namespace).List(ctx, opts)
    if err != nil {
        return nil, err
    }
    var results []ResourceInfo
    for _, deployment := range deployments.Items {
        results = append(results, ResourceInfo{Name: deployment.Name, Namespace: deployment.Namespace, Type: "deployment"})
    }
    return results, nil
}

func (s *StatefulSetResource) List(ctx context.Context, namespace string, opts metav1.ListOptions) ([]ResourceInfo, error) {
    statefulSets, err := s.clientset.AppsV1().StatefulSets(namespace).List(ctx, opts)
    if err != nil {
        return nil, err
    }
    var results []ResourceInfo
    for _, ss := range statefulSets.Items {
        results = append(results, ResourceInfo{Name: ss.Name, Namespace: ss.Namespace, Type: "statefulset"})
    }
    return results, nil
}

func (s *ServiceResource) List(ctx context.Context, namespace string, opts metav1.ListOptions) ([]ResourceInfo, error) {
    services, err := s.clientset.CoreV1().Services(namespace).List(ctx, opts)
    if err != nil {
        return nil, err
    }
    var results []ResourceInfo
    for _, svc := range services.Items {
        results = append(results, ResourceInfo{Name: svc.Name, Namespace: svc.Namespace, Type: "service"})
    }
    return results, nil
}

func (d *DaemonSetResource) List(ctx context.Context, namespace string, opts metav1.ListOptions) ([]ResourceInfo, error) {
    daemonSets, err := d.clientset.AppsV1().DaemonSets(namespace).List(ctx, opts)
    if err != nil {
        return nil, err
    }
    var results []ResourceInfo
    for _, ds := range daemonSets.Items {
        results = append(results, ResourceInfo{Name: ds.Name, Namespace: ds.Namespace, Type: "daemonset"})
    }
    return results, nil
}

func (s *SecretResource) List(ctx context.Context, namespace string, opts metav1.ListOptions) ([]ResourceInfo, error) {
    secrets, err := s.clientset.CoreV1().Secrets(namespace).List(ctx, opts)
    if err != nil {
        return nil, err
    }
    var results []ResourceInfo
    for _, secret := range secrets.Items {
        results = append(results, ResourceInfo{Name: secret.Name, Namespace: secret.Namespace, Type: "secret"})
    }
    return results, nil
}

func (c *ConfigMapResource) List(ctx context.Context, namespace string, opts metav1.ListOptions) ([]ResourceInfo, error) {
    configMaps, err := c.clientset.CoreV1().ConfigMaps(namespace).List(ctx, opts)
    if err != nil {
        return nil, err
    }
    var results []ResourceInfo
    for _, cm := range configMaps.Items {
        results = append(results, ResourceInfo{Name: cm.Name, Namespace: cm.Namespace, Type: "configmap"})
    }
    return results, nil
}

func (j *JobResource) List(ctx context.Context, namespace string, opts metav1.ListOptions) ([]ResourceInfo, error) {
    jobs, err := j.clientset.BatchV1().Jobs(namespace).List(ctx, opts)
    if err != nil {
        return nil, err
    }
    var results []ResourceInfo
    for _, job := range jobs.Items {
        results = append(results, ResourceInfo{Name: job.Name, Namespace: job.Namespace, Type: "job"})
    }
    return results, nil
}

func (c *CronJobResource) List(ctx context.Context, namespace string, opts metav1.ListOptions) ([]ResourceInfo, error) {
    cronJobs, err := c.clientset.BatchV1beta1().CronJobs(namespace).List(ctx, opts)
    if err != nil {
        return nil, err
    }
    var results []ResourceInfo
    for _, cj := range cronJobs.Items {
        results = append(results, ResourceInfo{Name: cj.Name, Namespace: cj.Namespace, Type: "cronjob"})
    }
    return results, nil
}

func (i *IngressResource) List(ctx context.Context, namespace string, opts metav1.ListOptions) ([]ResourceInfo, error) {
    ingresses, err := i.clientset.NetworkingV1().Ingresses(namespace).List(ctx, opts)
    if err != nil {
        return nil, err
    }
    var results []ResourceInfo
    for _, ing := range ingresses.Items {
        results = append(results, ResourceInfo{Name: ing.Name, Namespace: ing.Namespace, Type: "ingress"})
    }
    return results, nil
}

func (p *PodResource) Delete(ctx context.Context, namespace, name string) error {
    return p.clientset.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (r *ReplicaSetResource) Delete(ctx context.Context, namespace, name string) error {
    return r.clientset.AppsV1().ReplicaSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (d *DeploymentResource) Delete(ctx context.Context, namespace, name string) error {
    return d.clientset.AppsV1().Deployments(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *StatefulSetResource) Delete(ctx context.Context, namespace, name string) error {
    return s.clientset.AppsV1().StatefulSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *ServiceResource) Delete(ctx context.Context, namespace, name string) error {
    return s.clientset.CoreV1().Services(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (d *DaemonSetResource) Delete(ctx context.Context, namespace, name string) error {
    return d.clientset.AppsV1().DaemonSets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (s *SecretResource) Delete(ctx context.Context, namespace, name string) error {
    return s.clientset.CoreV1().Secrets(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *ConfigMapResource) Delete(ctx context.Context, namespace, name string) error {
    return c.clientset.CoreV1().ConfigMaps(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (j *JobResource) Delete(ctx context.Context, namespace, name string) error {
    return j.clientset.BatchV1().Jobs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (c *CronJobResource) Delete(ctx context.Context, namespace, name string) error {
    return c.clientset.BatchV1beta1().CronJobs(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

func (i *IngressResource) Delete(ctx context.Context, namespace, name string) error {
    return i.clientset.NetworkingV1().Ingresses(namespace).Delete(ctx, name, metav1.DeleteOptions{})
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
    case "replicasets":
        return &ReplicaSetResource{clientset: clientset}, nil
    case "deployments":
        return &DeploymentResource{clientset: clientset}, nil
    case "statefulsets":
        return &StatefulSetResource{clientset: clientset}, nil
    case "services":
        return &ServiceResource{clientset: clientset}, nil
    case "daemonsets":
        return &DaemonSetResource{clientset: clientset}, nil
    case "secrets":
        return &SecretResource{clientset: clientset}, nil
    case "configmaps":
        return &ConfigMapResource{clientset: clientset}, nil
    case "jobs":
        return &JobResource{clientset: clientset}, nil
    case "cronjobs":
        return &CronJobResource{clientset: clientset}, nil
    case "ingresses":
        return &IngressResource{clientset: clientset}, nil
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
	var kubeconfig string
	if kc := os.Getenv("KUBECONFIG"); kc != "" {
		kubeconfig = kc
	} else {
		home := homedir.HomeDir()
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
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
    case "replicaset":
        err = clientset.AppsV1().ReplicaSets(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "deployment":
        err = clientset.AppsV1().Deployments(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
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
