package main

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

func main() {
	x := CreatePod()
	fmt.Println(x)
}

func CreatePod() string {
	var kubeconfig *string
	if home := homeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	namespace := "default"
	podName := "test-logs"
	podsClient := clientset.CoreV1().Pods(namespace)
	// https://github.com/kubernetes/client-go/blob/kubernetes-1.18.0/examples/create-update-delete-deployment/main.go#L64
	// https://godoc.org/k8s.io/api/core/v1#Pod
	//https://godoc.org/k8s.io/api/core/v1#PodSpec
	// https://github.com/kubernetes/client-go/blob/kubernetes-1.18.0/kubernetes/typed/core/v1/pod.go#L116
	pod := &apiv1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: podName,
		},
		Spec: apiv1.PodSpec{
			Containers: []apiv1.Container{
				{
					Name:  "test",
					Image: "busybox",
					Ports: []apiv1.ContainerPort{
						{
							ContainerPort: 80,
						},
					},
					Command: []string{
						"/bin/sh", "-ec", "while :; do echo '$POD_ID: hello logs'; sleep 5 ; done",
					},
				},
			},
		},
	}
	fmt.Println("Creating a Pod ...")
	result, err := podsClient.Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Pod created %q.\n", result.GetObjectMeta().GetName())
	return result.GetObjectMeta().GetName()
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}
