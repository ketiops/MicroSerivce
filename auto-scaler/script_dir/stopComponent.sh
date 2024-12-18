"publishComponent.sh" 21L, 635C                                                                                               8,23          All
#/bin/bash

echo "Microservice cpu limit service publish"
kubectl delete -f ./cpu_limit/cpu-limit.yaml -n hpa
echo "Done!"

echo "Microservice Age Limit Service Publish"
kubectl delete -f ./cpu_age/age-limit.yaml -n hpa
echo "Done!"

echo "Microservice Request per Second Publish"
kubectl delete -f ./rps/rps-deploy.yaml
kubectl delete -f ./rps/rps-svc.yaml
kubectl delete -f ./rps/rps-hpa.yaml
echo "Done!"

echo "Microservice Memory Usage Based on Scale-OUT Service Publish"
kubectl delete -f ./mem-usage-latest/mem-deploy.yaml
kubectl delete -f ./mem-usage-latest/mem-svc.yaml
kubectl delete -f ./mem-usage-latest/mem-hpa.yaml
echo "Done!"
