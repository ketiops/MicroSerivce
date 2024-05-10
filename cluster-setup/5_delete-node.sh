## Specify delete node
nodename=sjjeon-m1

kubectl drain $nodename
kubectl delete node $nodename
