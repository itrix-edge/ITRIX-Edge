ITRIX-Edge: Software Define Edge
================================

**This project page needs reconstruction. Please check here for the latest release/version.**

ITRIX-Edge leverages containers and kubernetes to edge software can easily connect to major cloud providers.
The main focus of this project is to make deployment and maintanence easily.
With the build in connection to cloud provider, users can easily connnect the edge to cloud.


# Spec
Package	version
Kubernetes	v1.12.3
Docker	v18.06.1-ce
Azure Arc
prometheus
Grafana
 
 
# Deploy Node Requirement
Python3
pip3

# Master/Minion Node Requirement
Package	version
Supported Os	Ubuntu 18.04 LTS Server


## Usage
  

 
## Install requirement

```=shell
cd 
sudo pip3 install -r requirements.txt
```
Edit hosts.ini in 

Edit /extraVars.yml

## Regular Meeting note
https://docs.google.com/document/d/1wQb8q7dXOevTFSIFiWSf9xacT_8qqiqOgxSLDL-Gn3E/edit#

## Helm deployment
```=config 
helm_enabled: true

 
## Enable basic auth
# kube_basic_auth: true
## User defined api password
# kube_api_pwd: xk8suser

## Change default NodePort range 
# kube_apiserver_node_port_range: "9000-32767"
```

## Deploy
```=shell
su -
./ 
CLI
 Installer
```

## Usage  
    

## Examples


## Options
   
