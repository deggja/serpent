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

type KubeResourceInfo struct {
	Name      string
	Namespace string
	Kind      string
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

func getRandomResourceInfo(resourceType string) (KubeResourceInfo, error) {
	var resources []KubeResourceInfo

	switch resourceType {
	case "Pod":
		pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing pods: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, pod := range pods.Items {
			if pod.Namespace != "kube-system" && pod.Namespace != "kube-public" && !isCriticalPod(pod) {
				resources = append(resources, KubeResourceInfo{Name: pod.Name, Namespace: pod.Namespace, Kind: "Pod"})
			}
		}
	case "Deployment":
		deployments, err := clientset.AppsV1().Deployments("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing deployments: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, deployment := range deployments.Items {
			if deployment.Namespace != "kube-system" && deployment.Namespace != "kube-public" {
				resources = append(resources, KubeResourceInfo{Name: deployment.Name, Namespace: deployment.Namespace, Kind: "Deployment"})
			}
		}
	case "Service":
		services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing services: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, service := range services.Items {
			if service.Namespace != "kube-system" && service.Namespace != "kube-public" {
				resources = append(resources, KubeResourceInfo{Name: service.Name, Namespace: service.Namespace, Kind: "Service"})
			}
		}

	case "Job":
		jobs, err := clientset.BatchV1().Jobs("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing jobs: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, job := range jobs.Items {
			if job.Namespace != "kube-system" && job.Namespace != "kube-public" {
				resources = append(resources, KubeResourceInfo{Name: job.Name, Namespace: job.Namespace, Kind: "Job"})
			}
		}

	case "ConfigMap":
		configMaps, err := clientset.CoreV1().ConfigMaps("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing configmaps: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, configMap := range configMaps.Items {
			if configMap.Namespace != "kube-system" && configMap.Namespace != "kube-public" {
				resources = append(resources, KubeResourceInfo{Name: configMap.Name, Namespace: configMap.Namespace, Kind: "ConfigMap"})
			}
		}

	case "Secret":
		secrets, err := clientset.CoreV1().Secrets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing secrets: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, secret := range secrets.Items {
			if secret.Namespace != "kube-system" && secret.Namespace != "kube-public" {
				resources = append(resources, KubeResourceInfo{Name: secret.Name, Namespace: secret.Namespace, Kind: "Secret"})
			}
		}

	case "StatefulSet":
		statefulSets, err := clientset.AppsV1().StatefulSets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing statefulsets: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, statefulSet := range statefulSets.Items {
			if statefulSet.Namespace != "kube-system" && statefulSet.Namespace != "kube-public" {
				resources = append(resources, KubeResourceInfo{Name: statefulSet.Name, Namespace: statefulSet.Namespace, Kind: "StatefulSet"})
			}
		}

	case "DaemonSet":
		daemonSets, err := clientset.AppsV1().DaemonSets("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing daemonsets: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, ds := range daemonSets.Items {
			if ds.Namespace != "kube-system" && ds.Namespace != "kube-public" {
				resources = append(resources, KubeResourceInfo{Name: ds.Name, Namespace: ds.Namespace, Kind: "DaemonSet"})
			}
		}

	case "PersistentVolume":
		persistentVolumes, err := clientset.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing persistent volumes: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, pv := range persistentVolumes.Items {
			if pv.Namespace != "kube-system" && pv.Namespace != "kube-public" {
				resources = append(resources, KubeResourceInfo{Name: pv.Name, Kind: "PersistentVolume"})
			}
		}

	case "PersistentVolumeClaim":
		persistentVolumeClaims, err := clientset.CoreV1().PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing persistent volume claims: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, pvc := range persistentVolumeClaims.Items {
			if pvc.Namespace != "kube-system" && pvc.Namespace != "kube-public" {
				resources = append(resources, KubeResourceInfo{Name: pvc.Name, Namespace: pvc.Namespace, Kind: "PersistentVolumeClaim"})
			}
		}

	case "ServiceAccount":
		serviceAccounts, err := clientset.CoreV1().ServiceAccounts("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing service accounts: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, sa := range serviceAccounts.Items {
			if sa.Namespace != "kube-system" && sa.Namespace != "kube-public" {
				resources = append(resources, KubeResourceInfo{Name: sa.Name, Namespace: sa.Namespace, Kind: "ServiceAccount"})
			}
		}
	case "NetworkPolicy":
		networkPolicy, err := clientset.NetworkingV1().NetworkPolicies("").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing network policies: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, np := range networkPolicy.Items {
			if np.Namespace != "kube-system" && np.Namespace != "kube-public" {
				resources = append(resources, KubeResourceInfo{Name: np.Name, Namespace: np.Namespace, Kind: "NetworkPolicy"})
			}
		}
	case "Node":
		nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			log.Fatalf("Error listing nodes: %s\n", err.Error())
			return KubeResourceInfo{}, err
		}
		for _, node := range nodes.Items {
			resources = append(resources, KubeResourceInfo{Name: node.Name, Kind: "Node"})
		}
	}

	if len(resources) == 0 {
		log.Println("No eligible resources found.")
		return KubeResourceInfo{}, nil
	}

	randIndex := rand.Intn(len(resources))
	return resources[randIndex], nil
}

func isCriticalPod(pod v1.Pod) bool {
	_, isCritical := pod.Annotations["scheduler.alpha.kubernetes.io/critical-pod"]
	return isCritical
}

func deleteKubeResource(resourceInfo KubeResourceInfo) error {
	var err error
	switch resourceInfo.Kind {
	case "Pod":
		err = clientset.CoreV1().Pods(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "Deployment":
		err = clientset.AppsV1().Deployments(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "Service":
		err = clientset.CoreV1().Services(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "Job":
		err = clientset.BatchV1().Jobs(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "ConfigMap":
		err = clientset.CoreV1().ConfigMaps(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "Secret":
		err = clientset.CoreV1().Secrets(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "StatefulSet":
		err = clientset.AppsV1().StatefulSets(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "DaemonSet":
		err = clientset.AppsV1().DaemonSets(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "PersistentVolume":
		err = clientset.CoreV1().PersistentVolumes().Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "PersistentVolumeClaim":
		err = clientset.CoreV1().PersistentVolumeClaims(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "ServiceAccount":
		err = clientset.CoreV1().ServiceAccounts(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "NetworkPolicy":
		err = clientset.NetworkingV1().NetworkPolicies(resourceInfo.Namespace).Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	case "Node":
		err = clientset.CoreV1().Nodes().Delete(context.TODO(), resourceInfo.Name, metav1.DeleteOptions{})
	}

	if err != nil {
		log.Printf("Error deleting %s %s in namespace %s: %s\n", resourceInfo.Kind, resourceInfo.Name, resourceInfo.Namespace, err)
		return err
	}

	log.Printf("%s deleted: %s in namespace %s\n", resourceInfo.Kind, resourceInfo.Name, resourceInfo.Namespace)
	return nil
}
