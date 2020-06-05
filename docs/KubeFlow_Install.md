# KubeFlow Install
## 前置作業
* 準備好Kubernetes Cluster
* 安裝好NFS Server

## Install NFS Persistent Volumes
Install NFS Server
```shell=
$ sudo apt install nfs-common
$ sudo apt install nfs-kernel-server
$ sudo mkdir /var/nfsdata
```
Than you need to configure /etc/exports to share that directory:
```shell=
/nfsdata 192.168.0.0/16(rw,no_root_squash,no_subtree_check)
```
NFS Client
Each node of the cluster must be able to establish a connection to the NFS server. To enable this, install the following NFS client library on each node:
```shell=
sudo apt install nfs-common
```

## Install helm 3
```shell=
$ curl -L https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
```

## Install Dynamic Provisioner
Git clone helm chart source
```shell=
$ git clone https://github.com/helm/charts.git
```
Helm install Dynamic Provisioner
```shell=
$ cd charts/stable/
$ helm install   --generate-name  --set nfs.server=10.172.21.100   --set nfs.path=/var/nfsdata   --set storageClass.name=nfs   --set storageClass.defaultClass=true nfs-client-provisioner
$ kubectl get storageclass -n kubeflow
```
![](https://i.imgur.com/B3ybdDQ.png)

## Install Kubeflow
```shell=
$ wget https://github.com/kubeflow/kfctl/releases/download/v1.0.2/kfctl_v1.0.2-0-ga476281_linux.tar.gz
$ tar -xvf kfctl_v1.0.2-0-ga476281_linux.tar.gz
$ sudo mv ./kfctl /usr/local/bin/kfctl
$ mkdir kf-test
$ cd kf-test
$ kfctl build -V -f https://raw.githubusercontent.com/kubeflow/manifests/v1.0-branch/kfdef/kfctl_k8s_istio.v1.0.2.yaml
$ kfctl apply -V -f kfctl_k8s_istio.v1.0.2.yaml
$ kubectl -n kubeflow get all
```
![](https://i.imgur.com/F0cntgO.png)


## Remove Kubeflow
```shell=
$ kfctl delete -V -f kfctl_k8s_istio.v1.0.2.yaml
```
