package handler

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes/fake"
)

func createDeployment(replicas int32, annotation map[string]string) *appsv1.Deployment {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "demo-deployment",
			Annotations: annotation,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(replicas),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "demo",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "demo",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "web",
							Image: "nginx:1.12",
							Ports: []corev1.ContainerPort{
								{
									Name:          "http",
									Protocol:      corev1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	return deployment
}

func TestHandleDeploymentScaleDown(t *testing.T) {
	// Create the fake client.
	client := fake.NewSimpleClientset()

	//deploymentsClient := client.AppsV1().Deployments(corev1.NamespaceDefault)

	deploymentClient := client.AppsV1().Deployments("default")
	deployment := createDeployment(2, make(map[string]string))
	_, err := deploymentClient.Create(deployment)

	if err != nil {
		t.Error("failed test.")
	}

	HandleDeployment(*deployment, deploymentClient, "stop")

}

func TestHandleDeploymentScaleDownTwice(t *testing.T) {
	// Create the fake client.
	client := fake.NewSimpleClientset()

	//deploymentsClient := client.AppsV1().Deployments(corev1.NamespaceDefault)

	deploymentClient := client.AppsV1().Deployments("default")
	deployment := createDeployment(2, make(map[string]string))
	_, err := deploymentClient.Create(deployment)

	if err != nil {
		t.Error("failed test.")
	}

	HandleDeployment(*deployment, deploymentClient, "stop")
	deployments, err := deploymentClient.List(metav1.ListOptions{})
	for _, deployment := range deployments.Items {

		HandleDeployment(deployment, deploymentClient, "stop")
	}

}

func TestHandleDeploymentScaleUpWithNoReplicas(t *testing.T) {
	// Create the fake client.
	client := fake.NewSimpleClientset()

	//deploymentsClient := client.AppsV1().Deployments(corev1.NamespaceDefault)

	deploymentClient := client.AppsV1().Deployments("default")
	deployment := createDeployment(2, make(map[string]string))
	_, err := deploymentClient.Create(deployment)

	if err != nil {
		t.Error("failed test.")
	}

	HandleDeployment(*deployment, deploymentClient, "start")

}

func TestHandleDeploymentScaleUpWithTargetReplicas(t *testing.T) {

	// Create the fake client.
	client := fake.NewSimpleClientset()

	deploymentClient := client.AppsV1().Deployments("default")
	annotation := map[string]string{"bal.io/target-replicas": "5"}
	deployment := createDeployment(0, annotation)
	logrus.Infof("replicas: %d", *deployment.Spec.Replicas)
	_, err := deploymentClient.Create(deployment)

	if err != nil {
		t.Error("failed test.")
	}
	HandleDeployment(*deployment, deploymentClient, "start")

}
