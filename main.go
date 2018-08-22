package main

import (
	"boink/handler"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/retry"
)

var clientset *kubernetes.Clientset

var namespace string

var pathToConfig string

var selectors string

var action string

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true})

	// Output to stdout instead of the default stderr, could also be a file.
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.DebugLevel)

}

func main() {
	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Kube config path for outside of cluster access",
			Destination: &pathToConfig,
		},

		cli.StringFlag{
			Name:        "namespace, n",
			Value:       "default",
			Usage:       "the namespace where the application will poll the service.",
			Destination: &namespace,
		},
		cli.StringFlag{
			Name:        "action, a",
			Value:       "none",
			Usage:       "The action to perform either (stop or start)",
			Destination: &action,
		},
		cli.StringFlag{
			Name:        "label, l",
			Value:       "default",
			Usage:       "The deployment selector based on the labels .",
			Destination: &selectors,
		},
	}

	app.Action = func(c *cli.Context) error {
		var err error
		clientset, err = getClient()
		if err != nil {
			logrus.Error(err)
			return err
		}
		return manageDeployments()

	}
	app.Run(os.Args)
}

func manageDeployments() error {
	deploymentClient := clientset.AppsV1().Deployments(namespace)
	var listOptions metav1.ListOptions

	if selectors != "" {
		listOptions = metav1.ListOptions{
			LabelSelector: selectors,
			Limit:         100,
		}

	} else {
		listOptions = metav1.ListOptions{}
	}
	deployments, err := deploymentClient.List(listOptions)
	if err != nil {
		logrus.Warnf("Failed to find deployments: %v", err)
		return err
	}
	for _, deployment := range deployments.Items {
		retryErr := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			handler.HandleDeployment(deployment, deploymentClient, action)
			return nil
		})
		if retryErr != nil {
			panic(fmt.Errorf("Update failed: %v", retryErr))
		}

	}
	return nil
}

func getClient() (*kubernetes.Clientset, error) {
	var config *rest.Config
	var err error
	if pathToConfig == "" {
		logrus.Info("Using in cluster config")
		config, err = rest.InClusterConfig()
		// in cluster access
	} else {
		logrus.Info("Using out of cluster config")
		config, err = clientcmd.BuildConfigFromFlags("", pathToConfig)
	}
	if err != nil {
		return nil, err
	}
	return kubernetes.NewForConfig(config)
}
