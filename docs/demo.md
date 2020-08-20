# Demo deployment for 研揚
## itri-edge platform deployment
### prework

1. ssh key
2. `$ sudo apt-get install ansible`
3. update kernal

[Building the NVIDIA Kernel](https://docs.nvidia.com/jetson/l4t/index.html#page/Tegra%2520Linux%2520Driver%2520Package%2520Development%2520Guide%2Fkernel_custom.html%23wwpID0E0FD0HA)

[Kernel config for TX2/Xavier to enable docker extensions & Kubernetes](https://gist.github.com/stevennick/71ba2c71bc43ad665e1aab93d6cc6372)

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
程式client.py透過Tersorflow_Serving_IP提供服務的API，輸入一張測試圖片image1.jpg，並得到預測結果out_image2.json和out_image2.jpg。
Label文件參考https://github.com/fpaupier/tensorflow-serving_sidecar/blob/master/data/labels.pbtxt
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
$ sudo apt install protobuf-compiler
用protoc your/path/to/object_detection/protos/string_int_label_map.proto --python_out=.命令生成string_int_label_map_pb2.py文件。
```
<補充>如果遇到tensorflow無法在jetson nano安裝
```
$ sudo pip3 install --pre --extra-index-url https://developer.download.nvidia.com/compute/redist/jp/v43 tensorflow==1.15.2+nv20.2
```

## Cloud Kubeflow Install
### prework 
* 準備Cloud Kubernetes Cluster
* 準備NFS Server (共用儲存空間for clould and edge)
### Install NFS Persistent Volumes
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

### Install helm 3
```shell=
$ curl -L https://raw.githubusercontent.com/helm/helm/master/scripts/get-helm-3 | bash
```

### Install Dynamic Provisioner
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

### Install Kubeflow
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

### PipeLine Source Code
編譯後透過Kubeflow UI上傳至kubeflow
編譯方法請參照官網[kubeflow pipeline](https://www.kubeflow.org/docs/pipelines/sdk/build-component/)
```python=
#!/usr/bin/env python3

import kfp
from kfp import dsl

def deployed_model_op(url, key):
    return dsl.ContainerOp(
        name='Deployed - Model',
        image='google/cloud-sdk:279.0.0',
        command=['sh', '-c'],
        arguments=['curl -X GET -k -v http://$0:9000/v1/key/$1 | tee $2', url, key, '/tmp/results.txt'],
        file_outputs={
            'data': '/tmp/results.txt',
        }
    )

def echo2_op(text1):
    return dsl.ContainerOp(
        name='echo',
        image='library/bash:4.4.23',
        command=['sh', '-c'],
        arguments=['echo "Text 1: $0"', text1]
    )


@dsl.pipeline(
  name='Deployed pipeline',
  description='Trigger edge cluster to get model and prints the concatenated result.'
)
def download_and_join(
    url='',
    key='',
):

    deployed_task = deployed_model_op(url, key)

    echo_task = echo2_op(deployed_task.output)

if __name__ == '__main__':
    kfp.compiler.Compiler().compile(download_and_join, __file__ + '.yaml')
```


### Remove Kubeflow
```shell=
$ kfctl delete -V -f kfctl_k8s_istio.v1.0.2.yaml
```
### Cloud Infomation
| Master        | Worker        | NFS           | Kubeflow            |
| ------------- | ------------- | ------------- | ------------------- |
| 10.172.21.100 | 10.172.21.101 | 10.172.21.100 | 10.172.21.100:31380 |

### Edge Infomation
TX2
| node              | Serving Service  | Master            | Deploy        |
| ----------------- | ---------------- | ----------------- | ------------- |
| 192.168.1.101-104 | 10.12.50.5:30005 | 192.168.1.101-102 | 192.168.1.103 |

研揚
| node              | Serving Service  | Master            | Deploy       |
| ----------------- | ---------------- | ----------------- | ------------ |
| 192.168.1.101-106 | 10.12.50.6:30005 | 192.168.1.105-106 | 192.168.1.11 |

### 註冊yaml on Edge
8770a079a419b8eed7adf746eebd2
```shell=
curl -X POST -k -v POST http://192.168.1.151:9000/v1/deploymentTemplate -d  '{ "namespace": "default", "options": [{"key": "template_key", "value": "template_value" }], "deployment_template": {
    "apiVersion": "apps/v1",
    "kind": "Deployment",
    "metadata": {
        "labels": {
            "app": "faster-rcnn-resnet"
        },
        "name": "faster-rcnn-resnet",
        "namespace": "default"
    },
    "spec": {
        "selector": {
            "matchLabels": {
                "app": "faster-rcnn-resnet"
            }
        },
        "template": {
            "metadata": {
                "labels": {
                    "app": "faster-rcnn-resnet",
                    "version": "v1"
                }
            },
            "spec": {
                "containers": [
                    {
                        "args": [
                            "--rest_api_port=8501",
                            "--model_name=faster_rcnn_resnet",
                            "--model_base_path=/var/data/faster-rcnn-resnet101"
                        ],
                        "command": [
                            "/usr/bin/tensorflow_model_server"
                        ],
                        "image": "emacski/tensorflow-serving:latest-linux_arm64",
                        "imagePullPolicy": "IfNotPresent",
                        "name": "faster-rcnn-resnet",
                        "ports": [
                            {
                                "containerPort": 8501
                            }
                        ],
                        "volumeMounts": [
                            {
                                "mountPath": "/var/data",
                                "name": "nfs"
                            }
                        ]
                    }
                ],
                "volumes": [
                    {
                        "name": "nfs",
                        "nfs": {
                            "server": "10.172.21.100",
                            "path": "/var/data"
                        }
                    }
                ]
            }
        }
    }
}, "service_template": {
    "apiVersion": "v1",
    "kind": "Service",
    "metadata": {
        "labels": {
            "app": "faster-rcnn-resnet"
        },
        "name": "faster-rcnn-resnet-service",
        "namespace": "default"
    },
    "spec": {
        "ports": [
            {
                "name": "http-faster-rcnn-resnet",
                "port": 8501,
                "targetPort": 8501,
                "nodePort": 30005
            }
        ],
        "selector": {
            "app": "faster-rcnn-resnet"
        },
        "type": "NodePort"
    }
} }'
```

### Pipeline操作
登入kubeflow

![](https://i.imgur.com/ixmadDe.png)

選用Pipeline

![](https://i.imgur.com/QCn98be.png)

Create Demo2 Pipeline run

![](https://i.imgur.com/KCXeyo5.png)

填入url與key

![](https://i.imgur.com/8UA1dsQ.png)
