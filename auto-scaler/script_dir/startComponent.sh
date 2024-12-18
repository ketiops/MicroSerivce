#/bin/bash

echo "Microservice cpu limit service publish"
kubectl apply -f /root/replicas/deploy/cpu-limit.yaml -n hpa
sleep 1
echo "Done!"

echo "Microservice Age Limit Service Publish"
kubectl apply -f /root/replicas/deploy/age-limit.yaml
sleep 1
echo "Done!"

echo "Microservice Request per Second Publish"
kubectl apply -f /root/replicas/deploy/rps-deploy.yaml
kubectl apply -f /root/replicas/deploy/rps-svc.yaml
kubectl apply -f /root/replicas/deploy/rps-hpa.yaml
sleep 1
echo "Done!"

echo "Microservice Memory Usage Based on Scale-OUT Service Publish"
kubectl apply -f /root/replicas/deploy/mem-deploy.yaml
kubectl apply -f /root/replicas/deploy/mem-svc.yaml
kubectl apply -f /root/replicas/deploy/mem-hpa.yaml
sleep 1
echo "Done!"
