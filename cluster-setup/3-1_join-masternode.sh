## Master node's ip address
## please edit mater ip address
## set keepalived's ip address
masterIP=10.0.2.130

## Joining node's ip address
## please edit node ip address
nodeIP=10.0.2.131

## k8s node name
## please edit node name
nodename=sjjeon-m1

## Master's k8s token & hash
TOKEN=$(kubeadm token list | grep 'authentication' | cut -f 1 -d " ")
HASH=$(openssl x509 -pubkey -in /etc/kubernetes/pki/ca.crt | openssl rsa -pubin -outform der 2>/dev/null | openssl dgst -sha256 -hex | sed 's/^.* //')

## Validation
## please edit hostname
if [ -z "$TOKEN" ]; then
	new_token=$(kubeadm token create)
	cert_key=$(kubeadm init phase upload-certs --upload-certs | tail -1)
	ssh root@$nodeIP "swapoff -a; kubeadm join $masterIP:6443 --token $new_token --discovery-token-ca-cert-hash sha256:$master_hash --control-plane --certificate-key $cert_key; mkdir -p $HOME/.kube; cp -i /etc/kubernetes/admin.conf $HOME/.kube/config; chown $(id -u):$(id -g) $HOME/.kube/config; exit"
else
	token=$TOKEN
	master_hash=$HASH
	cert_key=$(kubeadm init phase upload-certs --upload-certs | tail -1)
	ssh root@$nodeIP "swapoff -a; kubeadm join $masterIP:6443 --token $token --discovery-token-ca-cert-hash sha256:$master_hash --control-plane --certificate-key $cert_key; mkdir -p $HOME/.kube; cp -i /etc/kubernetes/admin.conf $HOME/.kube/config; chown $(id -u):$(id -g) $HOME/.kube/config; exit"
fi

## Checking node status
nodestatus=$(kubectl get nodes | grep $nodename | sed 's/   /,/g' | cut -d "," -f2)
while [ $nodestatus = "NotReady" ];
do
	nodestatus=$(kubectl get nodes | grep $nodename | sed 's/   /,/g' | cut -d "," -f2)
	nodestatus=$nodestatus
	echo -n .;
	sleep 1;
done

echo -e "\n"
