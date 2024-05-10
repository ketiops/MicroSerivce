# install necessary packages
sudo apt install haproxy -y

# setting to haproxy
cat << EOF | sudo tee -a /etc/haproxy/haproxy.cfg
frontend kubernetes-master-lb
bind 0.0.0.0:26443
option tcplog
mode tcp
default_backend kubernetes-master-nodes
backend kubernetes-master-nodes
mode tcp
balance roundrobin
option tcp-check
option tcplog
server sjjeon-m1 10.0.2.131:6443 check
server sjjeon-m2 10.0.2.132:6443 check
server sjjeon-m3 10.0.2.131:6443 check
EOF

# restart the server
sudo systemctl restart haproxy
sudo systemctl enable haproxy
