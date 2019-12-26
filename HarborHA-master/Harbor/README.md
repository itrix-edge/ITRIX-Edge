### Harbor Cluster  

#### Requirement  

**Software**  
Docker  > 1.10  
Docker compose > 1.6.0  
Python > v2.7

**Ports**  
80  
443  
4443  


| Role | IP |
| ------- | -------------- |
| Harbor 1 | 100.67.191.10 |
| Harbor 2 | 100.67.191.11 |
| Harbor LB1 | 100.67.191.119 |
| Harbor LB2 | 100.67.191.118 |
| HarborVIP | 100.67.191.201 |  
  
keepalived virtual router ID : 75
1. 下載Harbor v1.5.0 offline package  
[On Harbor 1,2]  
```bash=
#vi /etc/hostname <= 編輯所有的Harbor 資訊
cd ~
wget https://storage.googleapis.com/harbor-releases/release-1.5.0/harbor-offline-installer-v1.5.0.tgz
tar -zxvf harbor-offline-installer-v1.5.0.tgz
```

2. 安裝並設置NFS  
[On Harbor 1,2]  
```bash=
sudo yum -y install epel-release
sudo yum -y install nfs-utils
sudo yum -y install python-pip
yum -y install docker docker-registry
sudo pip install docker-compose

# Harbor預設的registry放在 /data 目前還不能透過config直接修改
# 所以我們直接把NFS mount到 /data 上面
sudo rm -r /data # 移除舊資料
sudo mkdir /data
cd /
sudo mount 100.67.191.8:/var/nfs/harbor /data
# 100.67.191.8:/var/nfs/harbor 是我們自行設定的nfs server路徑
# 請依情況更改
 
sudo chown 10000:10000 /data
# Harbor官方要求修改/data權限
vim /etc/fstab
```
架設完成後請分別於 Harbor 1,2測試讀寫檔案能否同步  


3. 修改Harbor.cfg 並安裝Harbor  
[On Harbor 1,2]  
Edit ~/harbor/harbor.cfg  

```
# hostname 改為 VIP
hostname = 100.67.191.201

#db_host改成mariaDB的VIP
db_host = 100.67.191.120

#The password for the root user of Harbor DB. Change this before any production use.
db_password = root123

#The port of Harbor database host
db_port = 3306

#The user name of Harbor database
db_user = root

##### End of Harbor DB configuration#######

#redis_url 改成redis的VIP
redis_url = 100.67.191.110:6379

sudo ./install.sh --ha
```

4.修改iptables並儲存
```
iptables -t nat -A PREROUTING -p tcp -d 100.67.165.202/16 --dport 80 -j REDIRECT
iptables -t nat -A PREROUTING -p tcp -d 100.67.165.202/16 --dport 443 -j REDIRECT
iptables -t nat -A PREROUTING -p tcp -d 100.67.165.202/16 --dport 4443 -j REDIRECT
iptables-save > /etc/iptables-rules
vim a+x /etc/rc.d/rc.local
chmod a+x /etc/rc.d/rc.local

```

5. 設定Harbor LoadBalancer  
[On Harbor LB 1,2]  
安裝keepalived  
```bash=
sudo apt-get install keepalived
```
  
  
[On Harbor LB 1,2]  
Edit [/etc/keepalived/keepalived.conf](https://github.com/mJace/HarborHA/blob/master/Harbor/keepalived_cnf/keepalived.conf)  

```
global_defs {
  router_id haborlb
}
vrrp_sync_groups VG1 {
  group {
    VI_1
  }
}
#Please change "ens160" to the interface name on you loadbalancer hosts.
#In some case it will be eth0, ens16xxx etc.
vrrp_instance VI_1 {
  interface enp0s8

  track_interface {
    enp0s8
  }

  state BACKUP
  virtual_router_id 75
  priority 100
  nopreempt
  virtual_ipaddress {
    100.67.191.201/16
  }
  advert_int 1
  authentication {
    auth_type PASS
    auth_pass d0cker
  }

}

virtual_server 100.67.191.201 80 {
  delay_loop 15
  lb_algo rr
  lb_kind DR
  protocol TCP
  nat_mask 255.255.0.0
  persistence_timeout 10

  real_server 100.67.191.10 80 {
    weight 10
    MISC_CHECK {
        misc_path "/usr/local/bin/check.sh 100.67.191.10"
        misc_timeout 5
    }
  }

  real_server 100.67.191.11 80 {
    weight 10
    MISC_CHECK {
        misc_path "/usr/local/bin/check.sh 100.67.191.11"
        misc_timeout 5
    }
  }
}

virtual_server 100.67.191.201 443 {
  delay_loop 15
  lb_algo rr
  lb_kind DR
  protocol TCP
  nat_mask 255.255.0.0
  persistence_timeout 10

  real_server 100.67.191.10 443 {
    weight 10
    MISC_CHECK {
       misc_path "/usr/local/bin/check.sh 100.67.191.10"
       misc_timeout 5
    }
  }

  real_server 100.67.191.11 443 {
    weight 10
    MISC_CHECK {
       misc_path "/usr/local/bin/check.sh 100.67.191.11"
       misc_timeout 5
    }
  }
}

virtual_server 100.67.191.201 4443 {
  delay_loop 15
  lb_algo rr
  lb_kind DR
  protocol TCP
  nat_mask 255.255.0.0
  persistence_timeout 10

  real_server 100.67.191.10 4443 {
    weight 10
    MISC_CHECK {
        misc_path "/usr/local/bin/check.sh 100.67.191.10"
        misc_timeout 5
    }
  }

  real_server 100.67.191.11 4443 {
    weight 10
    MISC_CHECK {
        misc_path "/usr/local/bin/check.sh 100.67.191.11"
        misc_timeout 5
    }
  }
}

```

Save the server health check script to /usr/local/bin/check.sh  
```bash
wget https://raw.githubusercontent.com/goharbor/harbor/release-1.5.0/make/ha/sample/active_active/check.sh
sudo mv check.sh /usr/local/bin
sudo chmod +x /usr/local/bin/check.sh
```


Enable ip forward  
```bash=
# add the follow two lines to /etc/sysctl.conf

net.ipv4.ip_forward = 1
net.ipv4.ip_nonlocal_bind = 1

#Run the follow command to apply the change.
sysctl -p
```

重啟Keepalived使其套用新設定  
```bash=
systemctl restart keepalived
```

6. 設定Harbor Node1,2 iptable  
[On Harbor 1,2]  
```bash=
iptables -t nat -A PREROUTING -p tcp -d 100.67.191.201 --dport 80 -j REDIRECT
iptables -t nat -A PREROUTING -p tcp -d 100.67.191.201 --dport 443 -j REDIRECT
iptables -t nat -A PREROUTING -p tcp -d 100.67.191.201 --dport 4443 -j REDIRECT
iptables-save > /etc/iptables-rules
```

7. 安裝Harbor  
[On Harbor 1,2]  
```bash=
cd ~/harbor
sudo ./install.sh --ha
# 過程中會跳出詢問/data是否是綁訂到Shared Storage上
# 這時請輸入yes
# 安裝完成後可用docker ps 觀察container狀況
watch docker ps
```

#### 測試Harbor HA cluster  
1. 網頁登入   
在任一網路可通的機器上用瀏覽器打開VIP 100.67.191.201  
帳號: ```admin```  
密碼: ```Harbor12345```  
應可正常登入  

2. Docker client登入  
在任一可連至Harbor cluster的機器上  
由於我們Harbor沒有啟用https所以要先設定insecure registry  
Edit /etc/docker/daemon.json  
```
{"insecure-registries":["100.67.191.201"]}
```
```bash=
sudo systemctl restart docker
```

Login to Harbor  
```bash=
sudo docker login 100.67.191.201
# admin
# Harbor12345
```


3. Docker client push  
假設使用者已經在Harbor網頁建立專案 Test1，以push ubuntu:16.04為例  
```bash=
docker tag ubuntu:16.04 100.67.191.201/test1/ubuntu:16.04
docker push 100.67.191.201/test1/ubuntu:16.04
```


4. Docker client pull  
從Harbor上的專案test1，pull ubuntu 16.04  

```bash=
docker pull harbor.domain.com/test1/ubuntu:16.04
```
