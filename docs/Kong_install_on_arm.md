# Kong Install on Arm
### Install by YAML manfifests
* Download Kong YAML and change image
```shell=
$ wget https://bit.ly/kong-ingress-dbless
$ vim kong-ingress-dbless

chage image
image: kong:2.0 ----> image: kong:2.0.1-ubuntu
image: kong-docker-kubernetes-ingress-controller.bintray.io/kong-ingress-controller:0.8.1 ---->  image: skilledboy/kong-ingress-controller:0.8.0
```
* Deploy Kong via kubectl
```shell=
$ kubectl apply -f kong-ingress-dbless
```
![](https://i.imgur.com/IzOCNyA.png)

### Install by YAML manfifests
* Git clone Kong helm chart
```shell=
$ git clone https://github.com/Kong/charts.git
```
* Change helm chart values.yaml image version
```shell=
$ vim values.yaml

change image
image:                           image:   
  repository: kong     ---->       repository: kong 
  tag: "2.0"                       tag: "kong:2.0.1-ubuntu"
  
image:                                                                                                       image:
    repository: kong-docker-kubernetes-ingress-controller.bintray.io/kong-ingress-controller    ---->           repository: skilledboy/kong-ingress-controller:0.8.0
    tag: 0.8.0                                                                                                  tag: 0.8.0  
```
### Deploy Kong via helm
* Helm 2
```shell
$ helm install charts/kong
```
* Helm 3
```shell=
$ helm install charts/kong --generate-name --set ingressController.installCRDs=false
```
