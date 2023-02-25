# Boink
Boink is a simple Go client application that can handle stopping and starting Kubernetes `Deployments`.
It works by selecting `Deployments` based on labels.  It can also remember the previous known replicas, unlike a standard `kubectl scale` command where you need to specify the replicas manually.

This tool can be helpful when you have certain applications which needs to be stopped during certain period of time.
Pair this tool with kubernetes `CronJob` to automatically stop or start a `Deployment`


## How to run the application

1. Outside of the cluster - This is a normal executable program to interact with kubernetes cluster.  You can use minikube.

    Commands:
    - start - To start deployments
    - stop - To stop deployments

    Arguments:
    - --config - The location of $KUBECONFIG
    - --namespace - The namespace to use.
    - --label - Specify the selectors.
    

    To Stop:

       boink stop --config $KUBECONFIG --namespace test --label app=nginx

    To Start:
        
       boink start --config $KUBECONFIG --namespace test --label app=nginx

2.  In cluster

    Make sure you have the right permission in the cluster.  See samples in `manifest/` folder.
    
    1. Create the `ServiceAccount`

    ```
    apiVersion: v1
    kind: ServiceAccount
    metadata:
      name: boink
      namespace: test
    ```

    2. Create the `ClusterRole`

    ```
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
    name: boink
    rules:
    - apiGroups: ["extensions","apps"]
    resources: ["deployments"]
    verbs: ["get", "list", "watch", "create", "update", "patch", "delete"]

    ```

    3.  Bind the Cluster role with the service account

    ```
    apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
        name: boink        
    roleRef:
      apiGroup: rbac.authorization.k8s.io
      kind: ClusterRole
      name: boink
    subjects:
    - kind: ServiceAccount
      name: boink
      namespace: test
    ```

    4.  Create the `CronJob`, Stop `Deployment` with label `app=nginx` every  minute.

    ```
    apiVersion: batch/v1beta1
    kind: CronJob
    metadata:
      name: nginx-starter
      namespace: test
    spec:
      #this is in UTC  will run at 6:46 AM everyday
      schedule: "46 6 * * *"
      startingDeadlineSeconds: 10
      concurrencyPolicy: Forbid
      jobTemplate:
        spec:      
          template:
            spec:
              serviceAccountName: boink
              containers:
              - name: boink
                image: boink:1.0
                command: ["/boink"]
                args: ["--namespace","test", "--label", "app=nginx", "--action" , "start"]
              restartPolicy: OnFailure
    ```



### Tools:
1. command line library using `cobra`.  This library makes it easy to create a command line application in `Go`.

2.  Uses client-go to interact with kubernetes cluster.

3.  `github.com/sirupsen/logrus` - Standard logging mechanism for `Go`.

Check the `go.mod` to see all the dependencies.

## To build and run 
1. Make sure you enable `$GO111MODULE` to `on`
2. Go to the working directory `$GOPATH/src/boink/`.
3. Do `go build`
4. Do `go test ./... -cover` to run unit test with code coverage.
5. Finally to run the application `boink --config $KUBECONFIG --namespace test --label app=nginx stop`.  


If you are using skaffold, there is `skaffold.yaml` included at the root of the project.  Simply do a `skaffold dev` from the `$GOPATH/src/boink` and you are good to go.

