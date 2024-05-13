#bin/bash

kubectl apply -f kafka1.yaml
kubectl apply -f kafka2.yaml
kubectl apply -f kafka3.yaml
kubectl apply -f kafka4.yaml
kubectl apply -f kafka5.yaml
kubectl apply -f kafka6.yaml
kubectl apply -f mysql.yaml
kubectl apply -f grafana.yaml

kubectl get all -o wide -n ketiops-kafka
