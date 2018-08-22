# go-k8s
This is a sample code to interact with Kubernetes using client-go.

Two ways to run the application:
1. out of the cluster - This is a normal executable program to interact with kubernetes cluster.  You can use minikube.

2.  In cluster - TODO 


## What the application does

The application polls all the services in the `<your namespace>` namespace and gets all the kubernetes services in that namespace.  After which the application checks if the service contains the annotation `sample.com/job-orchestrator` is set to `true`.

This is a preclude to the possibility of performing actions when services with such annotations exist.

### Tools:
1. command line library using `github.com/urfave:v1.18.0`.  This library makes it easy to create a command line application in `Go`.

2.  Uses client-go to interact with kubernetes cluster.

3.  `github.com/Sirupsen/logrus` - Standard logging mechanism for `Go`.

Check the `Gopkg.toml` to see all the dependencies.

## To build and run 
1. Make sure you have `dep` installed in your system
2. Go to the working directory `$GOPATH/src/go-k8s/`.
3. Update all the dependencies using `dep ensure` to update all the dependencies.
4. Do `go build`
5. Finally to run the application `./go-k8s --config $KUBECONFIG --namespace <yournamespace>`.  This allows the application to retrieve the kubernetes configuration.

** If you need cluster admin, make sure you have the necessary administrator rights on the cluster.



