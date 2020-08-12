# Demo deployment for 研揚
## itri-edge platform deployment
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
## inference demo
下載物件辨識模型範例包
```
$ git clone https://github.com/fpaupier/tensorflow-serving_sidecar.git
$ cd tensorflow-serving_sidecar
```
該程式透過Tersorflow-Serving提供服務的API，輸入一張測試圖片image1.jpg，並得到預測結果out_image2.json和out_image2.jpg。
Label文件 : https://github.com/fpaupier/tensorflow-serving_sidecar/blob/master/data/labels.pbtxt
```
python3 client.py --server_url "http://<tersorflow_serving_IP>:30005/v1/models/faster_rcnn_resnet:predict" --image_path "/home/nvidia/tensorflow-serving_sidecar/object_detection/test_images/image1.jpg" --output_json "/home/nvidia/tensorflow-serving_sidecar/object_detection/test_images/out_image2.json" --save_output_image "TRUE" --label_map "/home/nvidia/tensorflow-serving_sidecar/data/labels.pbtxt"
```
<補充>如果遇到執行python3缺少lib
```
$ sudo su
$ apt-get install python3-pip
$ pip3 install --upgrade setuptools
$ pip3 install --upgrade pip
$ pip3 install numpy
$ pip3 install matplotlib
```
<補充>如果遇到tensorflow無法在jetson nano安裝
```
$ sudo pip3 install --pre --extra-index-url https://developer.download.nvidia.com/compute/redist/jp/v43 tensorflow==1.15.2+nv20.2
$ sudo apt install protobuf-compiler
用protoc your/path/to/object_detection/protos/string_int_label_map.proto --python_out=.命令生成string_int_label_map_pb2.py文件。
```
