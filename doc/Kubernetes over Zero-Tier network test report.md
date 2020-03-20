# Kubernetes over Zero-Tier Network

## Architecture
We use two hosts that are on different network segments and cannot be connected to each other. After adding these two hosts to the Zero-Tier network, the hosts will be assigned a Zero-Tier virtual network IP. We will use Zero-Tier virtual network IP for Kubernetes deployment.
![](https://i.imgur.com/o3L66vF.png)

## Install K8S on Zero-Tier Network
We use kubeadm to deplyed a kubernetes cluster over Zero-Tier virtual network. And we could see the master node can create well.
![](https://i.imgur.com/at5Z29A.png)

## Add Node Fail
When we use kubeadm to add node, that will get an error form kubeadm. The reason is the node can not use port 6443 to connecting. 
![](https://i.imgur.com/LEBShx9.png)

![](https://i.imgur.com/FNs5jzp.png)

![](https://i.imgur.com/az5bG79.png)

## Add Node Fail (root cause)
The reason of the node can not use port 6443 to connecting, because the firewalld.service still running and block port 6443.
![](https://i.imgur.com/IZakdmq.png)

![](https://i.imgur.com/7n9elcR.png)

## Add K8S node over Zero-Tier Network
After stop firewalld.service, We could add worker node to this cluster.
![](https://i.imgur.com/YWNLazK.png)

## K8S Status on Zero-Tier Network
The cluser seems that the service functions are working properly. And all the functions and components are running. We can even create pods normally.
![](https://i.imgur.com/3aPtnSN.png)

![](https://i.imgur.com/oL2dEGU.png)

## Virtual Network K8S issues 
Cannot connect to the container through the service.
![](https://i.imgur.com/ie90PTl.png)

