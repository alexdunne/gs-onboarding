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

- Install depedencies:
    - [Docker](https://docs.docker.com/get-docker/) 
    - [Minikube](https://minikube.sigs.k8s.io/docs/start/)
    - kubectl `brew install kubectl`
    - kubesec `brew install shyiko/kubesec/kubesec`
    - gpg `brew install gnupg`
- Create a `GPG` key
    - `gpg --full-generate-key`
    - Answers the questions
    - List your `GPG` keys
        - `gpg --list-secret-keys --keyid-format=long`
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
- Deploy the RabbitMQ Cluster
    - `kubectl apply -f ./k8s/base/rabbitmq/rabbitmq_cluster.yaml`
- Acquire the RabbitMQ credentials and update `./k8s/shared/secrets.yaml` (see below)
    - `kubectl get secret onboarding-rabbitmq-default-user -o jsonpath='{.data.username}'`
    - `kubectl get secret onboarding-rabbitmq-default-user -o jsonpath='{.data.password}'`
- Update the `./k8s/overlays/develop/secrets.yaml` file with your secrets
- Encrypt the secrets
    - `kubesec encrypt --key=pgp:<your_GPG_key> ./k8s/overlays/develop/secrets.yaml -o ./k8s/overlays/develop/secrets.enc.yaml`
    - Whenever you edit the `./k8s/overlays/develop/secrets.yaml` file redo this step
- Apply the secrets
    `kubesec decrypt ./k8s/overlays/develop/secrets.enc.yaml | kubectl apply -f -`
- Apply the deployments and services
    - `kustomize build k8s/overlays/develop | kubectl apply -f -`
- Create a tunnel to access the gateway service
    - `minikube tunnel`
- Access the gateway endpoint
    - `localhost:8000/all`
    - `localhost:8000/stories`
    - `localhost:8000/jobs`