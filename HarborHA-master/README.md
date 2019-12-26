# Harbor HA 安裝手冊

## Harbor HA 架構
**!!! 此文件的IP以及環境配置皆以這張架構圖為準，若環境不同，請參考架構圖自行更換IP !!!**  
![](https://i.imgur.com/ZcNEr5A.png)



| Role     | Node 1 IP     | Node 2 IP     | Node 3 IP    |
| -------- | ------------- | ------------- | ------------ | 
| Harbor   | 100.67.191.10 | 100.67.191.11 | N\A          |
| HarborLB | 100.67.191.119| 100.67.191.118| N\A          |
| HarborVIP| 100.67.191.201|               |              |
||||
| MariaDB  | 100.67.191.121| 
||||
| Redis    | 100.67.191.111| 100.67.191.112| N\A          |
| RedisVIP | 100.67.191.110|               |              |
||||
|NFS Server| 100.67.191.8  |               |              |

Harbor 此架構主要由三大項目組成   
#### 1. Harbor Node和其LoadBalancer   
(Active/Active)  
藉由keepalived將流量平均分配到其下Harbor Node  

#### 2. MariaDB Cluster  
(Active/Active)  
由三個MariaDB 透過Galera 建立cluster，在藉由Keepalived將流量平均分配  

#### 3. Redis Cluster
(Active/StandBy)  
透過Redis本身的Master/Slave設定加上Keepalived的腳本設計，達到Redis的高可用架構  

## 軟體版本/需求  


| Software | Version  | Description |
| -------- | -------- | ----------- |
| Harbor   | v1.5.0   | https://storage.googleapis.com/harbor-releases/release-1.5.0/harbor-offline-installer-v1.5.0.tgz |
| Redis    | >= 4.0.6   |             |
| MariaDB  | >= 10.2.14 |             |
| Python   | > v2.7   |             |
| Docker   | > 1.10   |             |
| Docker compose| > 1.6.0 |         |

## Network Ports / Settings  

| Port   | Protocol  | Description |
| -----  | --------  | ----------- |
| 80     | http  | Harbor UI and API will accept requests on this port for http protocol |
| 443    | https | Harbor UI and API will accept requests on this port for https protocol |
| 4443   | https | Connections to the Docker Content Trust service for Harbor, only needed when Notary is enabled|
| 6379   |    | Default redis port |
| 3306   |    | Default MariaDB port |
| 5432   |    | Default Postgress port |

**Keepalived setting**  

|role|vid|
|----|---|
|Harbor| 75 |
|Redis| 11 |
|MariaDB| 21 |


---


## 環境建置  
### Redis HA Cluster  
[Link to redis cluster setup guide](https://github.com/momobin/HarborHA/blob/master/redis)  


---


### Maria DB Cluster  

[Link to MariaDB cluster setup guide](https://github.com/momobin/HarborHA/tree/master/mariaDB)  




---


### NFS 架設  

Harbor建議使用具有共享且同步功能的儲存空間例如 Swift, S3, azure, Ceph 或是NFS  
在這個範例內我們直接在100.67.191.8上面架設了一個nfs server.  
可參考  
[Setup NFS Mounts On centos7.5 LTS Servers For Client Computers To Access](https://www.phpini.com/linux/rhel-centos-7-install-nfs-server/)  




---


### Harbor Cluster  
[Link to harbor cluster](https://github.com/momobin/HarborHA/tree/master/Harbor)  


---

## HA 注意事項  

1. 確保防火牆不會阻擋我們使用到的port  
2. Node fail 之後重開 有可能iptable prerouting會消失
   建議設定成每次開機自動設定  
3. Harbor Node上面NFS的mount也建議設定為開機自動Mount  
4. Redis Node上面redis sever fail或是node reboot後都要同時啟動redis-server和keepalived  
5. 不同service 的keepalived virtual router id必須不同，否則keepalived會失效。

---


## Reference  
[Redis Keepalived HA](https://www.jianshu.com/p/711542e7d347)  
[MariaDB cluster](https://www.techrepublic.com/article/how-to-set-up-a-mariadb-galera-cluster-on-ubuntu-16-04/)  
[Harbor HA Guide](https://github.com/goharbor/harbor/blob/master/docs/high_availability_installation_guide.md)  
