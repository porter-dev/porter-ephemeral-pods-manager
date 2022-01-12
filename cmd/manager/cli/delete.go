package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

var deleteCmd = &cobra.Command{
	Use: "delete",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := GetClientset()
		if err != nil {
			panic(err.Error())
		}

		err = deleteEphemeralPods(clientset)
		if err != nil {
			panic(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteEphemeralPods(clientset *kubernetes.Clientset) error {
	namespaces, err := clientset.CoreV1().Namespaces().List(context.Background(), v1.ListOptions{})
	if err != nil {
		return err
	}

	for _, namespace := range namespaces.Items {
		pods, err := clientset.CoreV1().Pods(namespace.Name).List(context.Background(), v1.ListOptions{
			LabelSelector: "porter/ephemeral-pod",
		})
		if err != nil {
			return err
		}

		for _, pod := range pods.Items {
			if time.Since(pod.CreationTimestamp.Time) >= (6 * time.Hour) {
				err = clientset.CoreV1().Pods(namespace.Name).Delete(
					context.Background(), pod.Name, v1.DeleteOptions{},
				)
				if err != nil {
					fmt.Printf("error deleting ephemeral pod: %s\n", pod.Name)
				}
			}
		}
	}

	return nil
}
