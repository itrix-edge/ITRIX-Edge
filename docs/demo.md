Demo deployment for 研揚
------------
### prework
1. ssh key

2. sudo apt-get install ansible

### download from git
```
$ sudo su
$ cd /root
$ git clone https://github.com/itrix-edge/ITRIX-Edge.git
$ cd ITRIX-Edge
$ git checkout v0.1.1
$ git submodule update --init
```
### pre-setup
```
$ cd /root/ITRIX-Edge/pre-setup
$ vi inventory.sample

[CP]
192.168.1.101
192.168.1.102
192.168.1.103
192.168.1.104

[OOBM]
```
```
$ ansible-playbook -i inventory.sample playbook.yml
```
### install K8S by kubespray
```
$ cd /root/ITRIX-Edge/kubespray
$ ansible-playbook -i /root/ITRIX-Edge/kubespray/inventory/edge/hosts.yaml cluster.yml
```
檢查
```
$ kubectl get node
NAME    STATUS   ROLES    AGE   VERSION
node1   Ready    master   79m   v1.16.6
node2   Ready    master   78m   v1.16.6
node3   Ready    master   78m   v1.16.6
node4   Ready    <none>   76m   v1.16.6
```

### install matallb
https://github.com/itrix-edge/metallb/tree/v0.9.3-itri
```
$ cd /root/ITRIX-Edge/metallb/manifests

$ kubectl apply -f namespace.yaml
$ kubectl apply -f current-config.yaml
$ kubectl apply -f metallb.yaml
$ kubectl create secret generic -n metallb-system memberlist --from-literal=secretkey="$(openssl rand -base64 128)"
```
### install posgress
```
$ docker run -it --name postgresql-local -p 192.168.1.103:5432:5432/tcp -e POSTGRES_PASSWORD=postgres -e POSTGRES_USER=postgres -e POSTGRES_DB=postgresdb -d postgres:11.8
```
### install edge-client-agent
```
$ cd edge-client-agent
$ vi external-IP.yml

kind: Service
apiVersion: v1
metadata:
  name: postgres-external
spec:
  ports:
  - protocol: TCP
    port: 5432
    targetPort: 5432
------------
kind: Endpoints
apiVersion: v1
metadata:
  name: postgres-external
subsets:
  - addresses:
      - ip: 192.168.1.103
    ports:
      - port: 5432
```
```
$ kubectl apply -f external-IP.yml
$ kubectl apply -f edge-agent-all.yaml
```
```
$ curl http://<10.233.62.205>:9000/v1/migrate/hook
{"result":true}

$ curl http://<10.233.62.205>:9000/v1/migrate/deploymentTemplate
{"result":true}
```
