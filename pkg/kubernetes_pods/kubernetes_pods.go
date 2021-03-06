package kubernetes_pods

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strconv"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig *string

type Pod struct {
	PodName       string
	NamespaceName string
}

func (p Pod) genkubeconfig() {
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = flag.String("kubeconfig", os.Getenv("KUBECONFIG"), "(optional) absolute path to the kubeconfig file")
	} else {
		if home := homeDir(); home != "" {
			kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
		}
	}
	flag.Parse()
}

func (p Pod) CreatePod(logsHits int) (string, error) {
	p.genkubeconfig()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if os.Getenv("INCLUSTER") != "" {
		log.Printf("Incluster configuration enabled")
		config, err = rest.InClusterConfig()
	}

	if err != nil {
		panic(err.Error())
	}

	// creating the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	podsClient := clientset.CoreV1().Pods(p.NamespaceName)
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: p.PodName,
			Labels: map[string]string{
				"k_logs": "true",
			},
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  p.PodName,
					Image: "busybox",
					Ports: []apiv1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
					Command: []string{
						"/bin/sh", "-ec", "for i in `seq 1 ${LOGS_HITS}`; do echo \"$POD_NAME: `date +\"%Y-%m-%dT%T\"`\"; done", "exit 0",
					},
					Env: []apiv1.EnvVar{
						{
							Name: "POD_NAME",
							ValueFrom: &apiv1.EnvVarSource{
								FieldRef: &apiv1.ObjectFieldSelector{
									FieldPath: "metadata.name",
								},
							},
						},
						{
							Name:  "LOGS_HITS",
							Value: strconv.Itoa(logsHits),
						},
					},
				},
			},
			RestartPolicy: "Never",
		},
	}

	log.Printf("pod \"%s\" is being created ...", p.PodName)
	result, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	log.Printf("pod created \"%q\"", result.GetObjectMeta().GetName())
	return result.GetObjectMeta().GetName(), nil
}

func (p Pod) DeletePod(podName string) (string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// creating the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	podsClient := clientset.CoreV1().Pods(p.NamespaceName)
	deletePolicy := metav1.DeletePropagationForeground

	log.Printf("pod \"%s\" is being deleted ...", podName)
	if err := podsClient.Delete(context.TODO(), podName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		log.Printf("Error deleting the pod \"%s\"", podName)
		return "", err
	}
	log.Printf("pod deleted \"%s\"", podName)
	return podName, nil
}

func (p Pod) Cleaner() {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// creating the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	listOptions := metav1.ListOptions{
		LabelSelector: "k_logs",
	}
	pods, err := clientset.CoreV1().Pods(p.NamespaceName).List(context.TODO(), listOptions)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("starting the cleaner process\n")
	// the p.PodName is already being deleted by the DeletePod func, so let's avoid repeating the same process again
	for _, i := range pods.Items {
		if p.PodName != i.ObjectMeta.Name {
			p.DeletePod(i.ObjectMeta.Name)
		}
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
