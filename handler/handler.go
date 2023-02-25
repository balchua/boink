package handler

import (
	"context"
	"strconv"

	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/client-go/kubernetes/typed/apps/v1"
)

const targetReplicasAnnotation string = "bal.io/target-replicas"

/*
HandleDeployment - Verifies if the deployment contains the "intraday-enabled"annotation,

	it will scale the deployment to zero to stop
	and scale to targetReplica to start.
*/
func HandleDeployment(deployment appsv1.Deployment, deploymentClient v1.DeploymentInterface, action string) {
	if action == "stop" {
		scaleToZero(deployment, deploymentClient)
	} else {
		scaleUp(deployment, deploymentClient)
	}

}

func scaleToZero(deployment appsv1.Deployment, deploymentClient v1.DeploymentInterface) {
	if *deployment.Spec.Replicas > int32(0) {
		logrus.Infof("Deployment (%s) has the annotation scaling to zero", deployment.ObjectMeta.Name)
		replicas := deployment.Spec.Replicas
		deployment.ObjectMeta.Annotations[targetReplicasAnnotation] = strconv.Itoa(int(*replicas))
		deployment.Spec.Replicas = int32Ptr(0)
		deploymentClient.Update(context.Background(), &deployment, metav1.UpdateOptions{})
	} else {
		logrus.Infof("Deployment (%s) is already scaled to (0)", deployment.ObjectMeta.Name)
	}

}

func scaleUp(deployment appsv1.Deployment, deploymentClient v1.DeploymentInterface) error {
	a := deployment.ObjectMeta.GetAnnotations()
	var replicas = 1
	var err error
	if *deployment.Spec.Replicas == int32(0) {
		if a[targetReplicasAnnotation] != "" {
			logrus.Infof("Deployment (%s) will be scaled up to (%s)", deployment.ObjectMeta.Name, a[targetReplicasAnnotation])
			replicas, err = strconv.Atoi(a[targetReplicasAnnotation])
			if err != nil {
				logrus.Error("Unable to convert replicas to number")
				panic(err)
			}
		}
		deployment.Spec.Replicas = int32Ptr(int32(replicas))
		_, err = deploymentClient.Update(context.Background(), &deployment, metav1.UpdateOptions{})
		if err != nil {
			return err
		}

	} else {
		logrus.Infof("Deployment (%s) is already scaled up to (%d)", deployment.ObjectMeta.Name, *deployment.Spec.Replicas)
	}
	return nil
}

func int32Ptr(i int32) *int32 {
	return &i
}
