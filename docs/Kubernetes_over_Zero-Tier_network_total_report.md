# Kubernetes over Zero-Tier Network Report
## Install and Setting Kubernetes
### Architecture
* The node default network interface is enp4s0f0. This two node can't connect by the interface enp4s0fs.
    * Master node ip is 100.67.170.1
    * GPU node ip is 10.174.14.99
* After join Zero-Tier network. We could see that there will be an additional network interface. And then the node will get virtual ip. This two node can connect by the virtual network.
    * Master node virtual ip is 192.168.196.37
    * GPU node virtual ip is 192.168.196.70

![](https://i.imgur.com/7ostIjW.png)
![](https://i.imgur.com/UC1HImt.png)
![](https://i.imgur.com/dP2tjEi.png)

### Install K8S over Zero-Tier Network bu kubeadm
When node join Zero-Tier network. We could see that there will be an additional network interface. And then we will use this interface to create Kubernetes cluster.
![](https://i.imgur.com/1tHbzvb.png)
* Use Zero-Tier network interface to init Kubernetes by kubeadm
```shell=
$ kubeadm init --apiserver-advertise-address 192.168.196.37 --pod-network-cidr=10.244.0.0/16
```
* Use kubectl to check kubernetes cluster
```shell=
$ kubectl get node
$ kubectl get pod -n kube-system
```
![](https://i.imgur.com/CUfTm5u.png)

* Add gup node by kubeadm
```shell=
$ kubeadm join 192.168.196.37:6443 --token zixk2c.hpz10zsu9q09fsur     --discovery-token-ca-cert-hash sha256:7f6a5499de455d39670ce6eb9151e5613b005d7a5fa0852513b506e9590e4e85
```
![](https://i.imgur.com/CLznzkZ.png)

* Running calico cni
```shell=
$ kubectl create -f https://docs.projectcalico.org/v3.11/manifests/calico.yaml
```
* Check Kubernetes cluster status
```shell=
$ kubectl get node
$ kubectl get pod -n kube-system
```
![](https://i.imgur.com/hojTJ7Q.png)

### Testing Kubernetes Function
* Create nginx pod and svc
```shell=
$ kubectl run nginx --image nginx --port 80
$ kubectl expose deployment nginx --port 50000 --target-port 80
```
* Pod & SVC create `OK`
![](https://i.imgur.com/nCeuYiz.png)

* Pod & SVC status `running`
![](https://i.imgur.com/mQtnyUY.png)

#### Worker node pod status
* Pod staus `running`
![](https://i.imgur.com/ErAT6BY.png)

* Get service over cluster IP `OK`
![](https://i.imgur.com/EPf7tSr.png)

#### Master expose service by nodeport
* Use nodeport to expose service `OK`
![](https://i.imgur.com/bBwGYyx.png)

* Use kubectl exec command to connect pod `fail`
![](https://i.imgur.com/wEp8xpM.png)

#### Virtual Network K8S issues
* Can't use kubectl exec to connect pod.
* Reason: Host default route still is enp4s0f0. When we use kubectl exec command, that network route path still use the default route. Because this two node can't use the default route to connection, so the kubectl exec command will fail.
 
![](https://i.imgur.com/bobzrif.png)

### Conclusion
Using Zero-Tier network to create Kubernetes cluser has some benefit below.
* Create Kubernetes everywhere: Every node can join to Zero-Tier network, that can get a vitrual ip. And all the node can connect over by this vitrual ip. So no matter how far away the nodes are, kubernetes cluster can be created. Even it doesn't need to be in the same network segment.
* Use Kubernetes API over Zero-Tier network: Kubernetes API can be used through Zero-Tier network.
* Zero-Tier support massive node: The Zero-Tier network supports a large number of nodes to join, and has a good and convenient management interface.
* Resource management of different regional nodes: Because all nodes can join the same Kubernetes cluster, the kubernetes master can control all resources and manage resources at will.

Using Zero-Tier network to create Kubernetes cluser has some weak point.
* kubectl exec command can't use: Because the node default route still need to use public interface. This will cause nodes to be unable to communicate with each other through the default route. So the kubectl command could not use.
* Need to use system resources: The Zero-Tier Kit is required for all nodes. And all node need to join to the same group.
