### Maria DB Cluster  

| Role | IP |
| ------- | -------------- |
| MariaDB  | 100.67.191.12 |
  

安裝MariaDB  
[On MariaDB]  
```
#vim /etc/yum.repos.d/mariadb.repo
[mariadb-main]
name = MariaDB Server
baseurl = https://downloads.mariadb.com/MariaDB/mariadb-10.3.9/yum/rhel/$releasever/$basearch
gpgkey = https://downloads.mariadb.com/MariaDB/RPM-GPG-KEY-MariaDB
gpgcheck = 1
enabled = 1

[mariadb-tools]
name = MariaDB Tools
baseurl = https://downloads.mariadb.com/Tools/rhel/$releasever/$basearch
gpgkey = https://downloads.mariadb.com/Tools/MariaDB-Enterprise-GPG-KEY
gpgcheck = 1
enabled = 1

#yum -y install MariaDB-server galera MariaDB-client
#systemctl start mariadb
#systemctl enable mariadb
#systemctl status mariadb
#mysql --version
```
安裝完成後 執行secure MariaDB安裝 且允許root remote登入並設定密碼  
```bash=
sudo mysql_secure_installation
```

#### 將個別DB載入 Harbor的schema  
嘗試過在MariaDB cluster架設起來之後在載入schema，但是一直會失敗  
測試到現在發現可行的方法就是在建立galera cluster之前就個別載入schema  

[On MariaDB]  
```bash=
wget https://raw.githubusercontent.com/goharbor/harbor/release-1.5.0/make/photon/db/registry.sql
mysql -u root -p < registry.sql
```
#### 設定使用者權限、允許遠端登入

```bash=
#mysql -u root -p
MariaDB [(none)]> SELECT User, Host FROM mysql.user;
+------+-----------+
| User | Host      |
+------+-----------+
| root | 127.0.0.1 |
| root | ::1       |
| root | localhost |
+------+-----------+
3 rows in set (0.000 sec)

MariaDB [(none)]> CREATE USER 'root'@'100.67.165.101' IDENTIFIED BY 'password';
Query OK, 0 rows affected (0.001 sec)

MariaDB [(none)]> GRANT ALL PRIVILEGES ON *.* TO 'root'@'100.67.165.101' WITH GRANT OPTION;

MariaDB [(none)]> FLUSH PRIVILEGES;

```

**Starting the Galera Cluster**  

#### 測試MariaDB功能
在任一點上執行  
```bash=
mysql -u root -p
> CREATE DATABASE tr_test;
>show databases;
```

#### 遠端連線測試在Harbor
```bash=
# mysql -u root -p -h 100.67.165.105
```
