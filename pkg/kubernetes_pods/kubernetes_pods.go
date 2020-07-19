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
	"k8s.io/client-go/tools/clientcmd"
)

var kubeconfig *string

func CreatePod(podName, namespaceName string, logsHits int) (string, error) {
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// creating the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	podsClient := clientset.CoreV1().Pods(namespaceName)

	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  podName,
					Image: "busybox",
					Ports: []apiv1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
					Command: []string{
						"/bin/sh", "-ec", "for i in `seq 1 ${LOGS_HITS}`; do echo \"$POD_NAME: `date +\"%Y-%m-%dT%T\"`\"; done",
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
		},
	}

	log.Printf("pod \"%s\" is being created ...", podName)
	result, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	log.Printf("pod created \"%q\"", result.GetObjectMeta().GetName())
	return result.GetObjectMeta().GetName(), nil
}

func DeletePod(podName, namespaceName string) (string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// creating the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	podsClient := clientset.CoreV1().Pods(namespaceName)
	deletePolicy := metav1.DeletePropagationForeground

	log.Printf("pod \"%s\" is being deleted ...", podName)
	if err := podsClient.Delete(context.TODO(), podName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return "", err
	}
	log.Printf("pod deleted \"%s\"", podName)
	return podName, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
