
| Role | IP |
| ------- | -------------- |
| Redis 1 | 100.67.191.111 |
| Redis 2 | 100.67.191.112 |
| Redis VIP | 100.67.191.110 |

Keepalived virtual router id : 11 <=修改一個不能相同的數字

**Redis 1為Master**  
**Redis 2為Slave**  

1. 於Redis 1和Redis 2下載 Redis 和 keepalived  
```bash=
wget http://download.redis.io/releases/redis-4.0.6.tar.gz
tar -zxvf redis-4.0.6.tar.gz
sudo yum -y install keepalived
```

2. Redis 1 and 2 編譯Redis  
```bash=
cd redis-4.0.6
yum -y update
yum -y install gcc
make MALLOC=libc
sudo make install
```

3. 設定Redis cluster  

[On Master]  
修改[redis.conf](https://github.com/mJace/HarborHA/blob/master/redis_src/redis_cnf/master/redis.conf)  
```
daemonize  yes
bind  0.0.0.0
pidfile /var/run/redis_6379.pid
logfile "/var/log/redis.log"
```

[On Slave]  
修改[redis.conf](https://github.com/mJace/HarborHA/blob/master/redis_src/redis_cnf/slave/redis.conf)，跟Master不一樣的是Slave多了一行```slaveof <masterIP> <port>```  
```
daemonize  yes
bind  0.0.0.0
pidfile /var/run/redis_6379.pid
logfile "/var/log/redis.log"
slaveof 100.67.191.111 6379
```

啟動Redis  
[On Master and Slave]  
```bash=
cd redis-4.0.6
sudo redis-server ./redis.conf
```

檢查Redis狀態  
[On Master and Slave]  
```bash=
redis-cli info replication
```


4. 設定Keepalived  
#### [On Master]  
edit [/etc/keepalived/keepalived.conf](https://github.com/mJace/HarborHA/blob/master/redis_src/keepalived_cnf/master/keepalived.conf)  
```
global_defs {
	router_id redis
}
 
vrrp_script chk_redis {
	script "/etc/keepalived/scripts/check_redis.sh"
	interval 4
	weight -5
	fall 3  
	rise 2
}
 
vrrp_instance VI_REDIS {
	state MASTER
	interface eth1
	virtual_router_id 11
	priority 100
	advert_int 1
	nopreempt
 
	authentication {
		auth_type PASS
		auth_pass 1111
	}
 
	virtual_ipaddress {
		192.168.99.10
	}
 
	track_script {
		chk_redis
	}
 
	notify_master /etc/keepalived/scripts/redis_master.sh
	notify_backup /etc/keepalived/scripts/redis_backup.sh
	notify_fault  /etc/keepalived/scripts/redis_fault.sh
	notify_stop   /etc/keepalived/scripts/redsi_stop.sh
}
```
**Edit script for redis HA**  
edit [/etc/keepalived/scripts/check_redis.sh](https://github.com/mJace/HarborHA/blob/master/redis_src/keepalived_cnf/master/scripts/check_redis.sh) for checking master reids  
```
#!/bin/bash  
CHECK=`/usr/local/bin/redis-cli PING`  
if [ "$CHECK" == "PONG" ] ;then  
      echo $CHECK  
      exit 0  
else   
      echo $CHECK  
      service keepalived stop 
      exit 1  
fi  
```

edit [/etc/keepalived/scripts/redis_backup.sh](https://github.com/mJace/HarborHA/blob/master/redis_src/keepalived_cnf/master/scripts/redis_backup.sh) for master back online and being slave of new master  
```
#!/bin/bash
REDISCLI="/usr/local/bin/redis-cli"
LOGFILE="/var/log/keepalived-redis-state.log"
echo "[backup]" >> $LOGFILE
date >> $LOGFILE
echo "Being slave...." >> $LOGFILE 2>&1
sleep 15
echo "Run SLAVEOF cmd ..." >> $LOGFILE
$REDISCLI SLAVEOF 192.168.99.12 6379 >> $LOGFILE  2>&1
```

edit [/etc/keepalived/scripts/redis_fault.sh](https://github.com/mJace/HarborHA/blob/master/redis_src/keepalived_cnf/master/scripts/redis_fault.sh)  
```
# !/bin/bash
LOGFILE=/var/log/keepalived-redis-state.log 
echo "[fault]" >> $LOGFILE 
date >> $LOGFILE
```

edit [/etc/keepalived/scripts/redis_master.sh](https://github.com/mJace/HarborHA/blob/master/redis_src/keepalived_cnf/master/scripts/redis_master.sh)  
```
#!/bin/bash
REDISCLI="/usr/local/bin/redis-cli"
LOGFILE="/var/log/keepalived-redis-state.log"
echo "[master]" >> $LOGFILE
date >> $LOGFILE
echo "Being master...." >> $LOGFILE 2>&1
echo "Run SLAVEOF cmd ..." >> $LOGFILE
$REDISCLI SLAVEOF 192.168.99.12 6379 >> $LOGFILE  2>&1
sleep 10 
echo "Run SLAVEOF NO ONE cmd ..." >> $LOGFILE
$REDISCLI SLAVEOF NO ONE >> $LOGFILE 2>&1
```

edit [/etc/keepalived/scripts/redis_stop.sh](https://github.com/mJace/HarborHA/blob/master/redis_src/keepalived_cnf/master/scripts/redis_stop.sh)  
```
# !/bin/bash
LOGFILE=/var/log/keepalived-redis-state.log 
echo "[stop]" >> $LOGFILE 
date >> $LOGFILE
```

#### [On Slave]  
edit **[/etc/keepalived/keeplived.conf](https://github.com/mJace/HarborHA/blob/master/redis_src/keepalived_cnf/slave/keepalived.conf)** for slave keepalived  
```
global_defs {
	router_id redis
}
 
vrrp_script chk_redis {
	script "/etc/keepalived/scripts/check_redis.sh"
	interval 4
	weight -5
	fall 3  
	rise 2
}
 
vrrp_instance VI_REDIS {
	state BACKUP
	interface eth1
	virtual_router_id 11
	priority 99
	advert_int 1
	nopreempt
 
	authentication {
		auth_type PASS
		auth_pass 1111
	}
 
	virtual_ipaddress {
		192.168.99.10
	}
 
	track_script {
		chk_redis
	}
 
	notify_master /etc/keepalived/scripts/redis_master.sh
	notify_backup /etc/keepalived/scripts/redis_backup.sh
	notify_fault  /etc/keepalived/scripts/redis_fault.sh
	notify_stop   /etc/keepalived/scripts/redsi_stop.sh
}
```

**Edit script for redis HA**  
[/etc/keepalived/scripts/check_redis.sh](https://github.com/mJace/HarborHA/blob/master/redis_src/keepalived_cnf/slave/scripts/check_redis.sh) for checking slave reids
```
#!/bin/bash  
CHECK=`/usr/local/bin/redis-cli PING`  
if [ "$CHECK" == "PONG" ] ;then  
      echo $CHECK  
      exit 0  
else   
      echo $CHECK  
      service keepalived stop
      exit 1  
fi
```


edit [/etc/keepalived/scripts/redis_backup.sh](https://github.com/mJace/HarborHA/blob/master/redis_src/keepalived_cnf/slave/scripts/redis_backup.sh) for master back online and being slave of new master  
```
#!/bin/bash
REDISCLI="/usr/local/bin/redis-cli"
LOGFILE="/var/log/keepalived-redis-state.log"
echo "[backup]" >> $LOGFILE
date >> $LOGFILE
echo "Being slave...." >> $LOGFILE 2>&1
sleep 15 
echo "Run SLAVEOF cmd ..." >> $LOGFILE
$REDISCLI SLAVEOF 192.168.99.11 6379 >> $LOGFILE  2>&1
```

edit [/etc/keepalived/scripts/redis_fault.sh](https://github.com/mJace/HarborHA/blob/master/redis_src/keepalived_cnf/slave/scripts/redis_fault.sh)  
```
# !/bin/bash
LOGFILE=/var/log/keepalived-redis-state.log 
echo "[fault]" >> $LOGFILE 
date >> $LOGFILE
```

edit [/etc/keepalived/scripts/redis_master.sh](https://github.com/mJace/HarborHA/blob/master/redis_src/keepalived_cnf/slave/scripts/redis_master.sh)  
```
#!/bin/bash
REDISCLI="/usr/local/bin/redis-cli"
LOGFILE="/var/log/keepalived-redis-state.log"
echo "[master]" >> $LOGFILE
date >> $LOGFILE
echo "Being master...." >> $LOGFILE 2>&1
echo "Run SLAVEOF cmd ..." >> $LOGFILE
$REDISCLI SLAVEOF 192.168.99.11 6379 >> $LOGFILE  2>&1
sleep 10 
echo "Run SLAVEOF NO ONE cmd ..." >> $LOGFILE
$REDISCLI SLAVEOF NO ONE >> $LOGFILE 2>&1
```

edit [/etc/keepalived/scripts/redis_stop.sh](https://github.com/mJace/HarborHA/blob/master/redis_src/keepalived_cnf/slave/scripts/redis_stop.sh)  
```
# !/bin/bash
LOGFILE=/var/log/keepalived-redis-state.log 
echo "[stop]" >> $LOGFILE 
date >> $LOGFILE
```

編輯完master/slave的scripts後  
記得用chmod a+x 把所有腳本設定為可執行  
並且在redis-server啟用的狀態下執行 check_redis.sh  
檢查腳本以及redis-server運作狀況  

**Start keepalived**  
**[On master and slave]**  
```bash=
sudo systemctl restart keepalived
```

#### Redis cluster測試  
1. Master/Slave寫入測試  

[On Master]  
```bash=
redis-cli set a hello
```  

[On Slave]  
```bash=
redis-cli get a
"hello"
```

2. VIP 寫入測試  
[On other machine]  

```bash=
redis-cli -h 100.67191.110 set b test
> OK
redis-cli -h 100.67191.110 get b
> "test"
```

3. HA 測試  

[On Redis 1]  
關閉redis server  
```bash
sudo pkill redis-server
```

測試針對VIP寫入與讀取，驗證master轉移到Redis 2  
[On other machine]  
```bash
redis-cli -h 100.67191.110 set c test
> OK
redis-cli -h 100.67191.110 get c
> "test"
```

[On Redis 2]  
```bash=
ip a
```
應該看到VIP 已經轉移到Redis 2上面  

回到Redis 1重新啟動Redis server和keepalived，驗證VIP 切換回Redis 1  
[On Redis 1]
```bash=
cd redis-4.0.6
sudo redis-server ./redis.conf
sudo systemctl restart keepalived
```
測試Redis 1是否能夠同步最新資訊  
[On Redis 1]
```bash=
redis-cli get c
> "test"
```
[On Redis 1]  
```bash=
ip a
```
應該看到VIP 已經轉移到Redis 1上面  

到此為止已經完成Redis HA環境的架設  
