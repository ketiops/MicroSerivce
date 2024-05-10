####################################
#### Change Version as you want ####
####################################
version=1.25.6-00

sudo apt-get update
sudo apt-get install apt-transport-https ca-certificates curl -y
curl -s https://packages.cloud.google.com/apt/doc/apt-key.gpg | apt-key add -
echo "deb http://apt.kubernetes.io/ kubernetes-xenial main" > /etc/apt/sources.list.d/kubernetes.list
sudo apt-get update
apt install -qy kubelet=$version kubectl=$version kubeadm=$version kubernetes-cni
sudo apt-mark hold kubelet kubeadm kubectl
