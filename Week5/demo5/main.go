package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "[可选] kubeconfig 绝对路径")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		fmt.Printf("error %s", err.Error())
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Printf("error %s", err.Error())
	}

	timeout := int64(60)
	watcher, err := clientSet.CoreV1().Pods("default").Watch(context.TODO(), metav1.ListOptions{TimeoutSeconds: &timeout})
	for event := range watcher.ResultChan() {
		item := event.Object.(*corev1.pod)

		switch event.Type {
		case watch.Added:
			processPod(item.GetName())
		case watch.Modified:
		case watch.Deleted:
		}
	}
}

func processPod(name string) {
	fmt.Printf("Pod %s\n", name)
}
