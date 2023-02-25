package cmd

import (
	"boink/handler"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/retry"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop deployments",
	Long:  `This command stop kubernetes deployments`,
	Run: func(cmd *cobra.Command, args []string) {
		Clientset, _ = getClient()
		deploymentClient := Clientset.AppsV1().Deployments(Namespace)
		deployments, err := getDeployments()
		if err != nil {
			panic(err)
		}
		if len(deployments.Items) == 0 {
			logrus.Error("no deployment found")
		}
		for _, deployment := range deployments.Items {
			retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
				handler.HandleDeployment(deployment, deploymentClient, "stop")
				return nil
			})
			if retryErr != nil {
				panic(fmt.Errorf("Update failed: %v", retryErr))
			}

		}
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)

}
