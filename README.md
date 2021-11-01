# gs-onboarding

[Architecture Diagram](https://lucid.app/lucidchart/invitations/accept/inv_b0ee37ba-6710-41c4-be4e-b4992c2d6e0d)

## Services

### Consumer

The consumer service periodically fetches the top stories from hacker news and stores all non dead nor deleted items in the database

### API

The API service is a gRPC server that offers a interface to fetched the stored hacker news stories

### Gateway

The gateway service is main entry point for third parties to access all other systems. Currently, it is responsible for proxying requests to the API service

## Local Kubernetes setup

- Install [Docker](https://docs.docker.com/get-docker/) and [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- Start Minikube 
    - `minikube start`
- Change the docker daemon your terminal points at 
    - `eval $(minikube -p minikube docker-env)`
    - This only changes your current terminal session. By closing the terminal, you will go back to using your own system’s docker daemon.
- Build all of the docker images:
    - `docker build -t onboarding-api --target api .`
    - `docker build -t onboarding-consumer --target consumer .`
    - `docker build -t onboarding-gateway --target gateway .`
    - `docker build -t onboarding-migrator --target migrator .`
- Install the RabbitMQ operator in your cluster 
    - `kubectl apply -f "https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml"`
    - `kubectl rabbitmq install-cluster-operator`
- Acquire the RabbitMQ credentials and update `./k8s/shared/secrets.yaml`
    - `kubectl get secret onboarding-rabbitmq-default-user -o jsonpath='{.data.username}'`
    - `kubectl get secret onboarding-rabbitmq-default-user -o jsonpath='{.data.password}'`
- Apply the shared secrets
    - `kubectl apply -f ./k8s/shared/secrets.yaml`
- Apply the infrastructure deployments and services
    - `kustomize build k8s/overlays/develop | kubectl apply -f -`
- Create a tunnel to access the gateway service
    - `minikube tunnel`
- Access the gateway endpoint
    - `localhost:8000/all`
    - `localhost:8000/stories`
    - `localhost:8000/jobs`