eval $(minikube docker-env)
docker build -t client:0.1 -f client/Dockerfile .
docker build -t service1:0.2 -f services/service1/Dockerfile .
docker build -t service1-migrations:0.1 -f services/service1/migrations/Dockerfile .
docker build -t service2:0.1 -f services/service2/Dockerfile .
