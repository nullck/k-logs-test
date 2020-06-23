package kubernetes_pods

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func CreatePod(podName string, namespaceName string) (string, error) {
	var kubeconfig *string
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
						"/bin/sh", "-ec", "while :; do echo \"$POD_NAME: hello logs\"; sleep 5 ; done",
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
					},
				},
			},
		},
	}
	fmt.Println("Creating the pod ...")
	result, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		return "", err
	}
	fmt.Printf("Pod created %q.\n", result.GetObjectMeta().GetName())
	return result.GetObjectMeta().GetName(), nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
