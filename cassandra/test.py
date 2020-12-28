
# Connecting to Cassandra
--------------------------
from cassandra.cluster import Cluster

cluster = Cluster()#127.0.0.1
cluster = Cluster(['10.98.187.203', ''])

--------------------------
#創建keyspace
session.execute('CREATE KEYSPACE IF NOT EXISTS cycling WITH replication = { 'class' : 'NetworkTopologyStrategy', 'dc1' : '3' };')

session = cluster.connect('cycling')

#session.set_keyspace('users')

--------------------------
# 創建table
session.execute('CREATE TABLE IF NOT EXISTS cycling.table_1 (
   id int, 
   firstname text, 
   lastname text, 
   age int, 
   affiliation text,
   country text,
   registration date,
   PRIMARY KEY (id));')
# 刪除table
session.execute('DROP TABLE cycling.table_1;')

--------------------------
from cassandra.cluster import Cluster, ExecutionProfile, EXEC_PROFILE_DEFAULT
from cassandra.policies import WhiteListRoundRobinPolicy, DowngradingConsistencyRetryPolicy
from cassandra.query import tuple_factory

#執行配置文件
profile = ExecutionProfile(
    load_balancing_policy=WhiteListRoundRobinPolicy(['127.0.0.1']),
    retry_policy=DowngradingConsistencyRetryPolicy(),
    consistency_level=ConsistencyLevel.LOCAL_QUORUM,
    serial_consistency_level=ConsistencyLevel.LOCAL_SERIAL,
    request_timeout=15,
    row_factory=tuple_factory
)
cluster = Cluster(execution_profiles={EXEC_PROFILE_DEFAULT: profile})
session = cluster.connect()

print(session.execute("SELECT release_version FROM system.local").one())
--------------------------

查詢keyspaces/tables/columns狀態
# -*- encoding: utf-8 -*-
# 引入Cluster模塊
#from cassandra.cluster import Cluster
# 引入DCAwareRoundRobinPolicy模塊，可用來自定義驅動程序的行為
#from cassandra.policies import DCAwareRoundRobinPolicy
# 默認本機數據庫集群(IP127.0.0.1).
#cluster = Cluster()
# 連接並創建一個會話
#session = cluster.connect()
# 查詢keyspaces/tables/columns狀態
print(cluster.metadata.keyspaces)
print(‘----------‘)
print(cluster.metadata.keyspaces[‘test‘].tables)
print(‘----------‘)
print(cluster.metadata.keyspaces[‘test‘].tables[‘user‘])
print(‘----------‘)
print(cluster.metadata.keyspaces[‘test‘].tables[‘user‘].columns)
print(‘----------‘)
print(cluster.metadata.keyspaces[‘test‘].tables[‘user‘].columns[‘age‘])
print(‘----------‘)
# 關閉連接
cluster.shutdown()
# 查看是否關閉連接
print(cluster.is_shutdown)

-----------------------


插入和查詢表中的數據
# -*- encoding: utf-8 -*-
# 引入Cluster模塊
from cassandra.cluster import Cluster
# 引入DCAwareRoundRobinPolicy模塊，可用來自定義驅動程序的行為
from cassandra.policies import DCAwareRoundRobinPolicy


# 默認本機數據庫集群(IP127.0.0.1).
cluster = Cluster()
# 連接並創建一個會話
session = cluster.connect()
# table中插入數據
session.execute(‘insert into test.user (name, age, email) values (%s, %s, %s);‘, [‘aaa‘, 21, ‘222@21.com‘])
session.execute(‘insert into test.user (name, age, email) values (%s, %s, %s);‘, [‘bbb‘, 22, ‘bbb@22.com‘])
session.execute(‘insert into test.user (name, age, email) values (%s, %s, %s);‘, [‘ddd‘, 20, ‘ccc@20.com‘])



-------------------------
# table中查詢數據
rows = session.execute(‘select * from test.user;‘)
for row in rows:
    print(row)
# 關閉連接
cluster.shutdown()
# 查看是否關閉連接
print(cluster.is_shutdown)


--------------------------
連接遠程數據庫
# -*- encoding: utf-8 -*-
from cassandra import ConsistencyLevel
# 引入Cluster模塊
from cassandra.cluster import Cluster
# 引入DCAwareRoundRobinPolicy模塊，可用來自定義驅動程序的行為
# from cassandra.policies import DCAwareRoundRobinPolicy
from cassandra.auth import PlainTextAuthProvider
from cassandra.query import SimpleStatement
import pandas as pd


# 配置Cassandra集群的IP
contact_points = [‘1.1.1.1‘, ‘2.2.2.2‘, ‘3.3.3.3‘]
# 配置登陸Cassandra集群的賬號和密碼
auth_provider = PlainTextAuthProvider(username=‘XXX‘, password=‘XXX‘)
# 創建一個Cassandra的cluster
cluster = Cluster(contact_points=contact_points, auth_provider=auth_provider)
# 連接並創建一個會話
session = cluster.connect()
# 定義一條cql查詢語句
cql_str = ‘select * from keyspace.table limit 5;‘
simple_statement = SimpleStatement(cql_str,consistency_level=ConsistencyLevel.ONE)
# 對語句的執行設置超時時間為None
execute_result = session.execute(simple_statement, timeout=None)
# 獲取執行結果中的原始數據
result = execute_result._current_rows
# 把結果轉成DataFrame格式
result = pd.DataFrame(result)
# 把查詢結果寫入csv
result.to_csv(‘連接遠程數據庫.csv‘, mode=‘a‘, header=True)
# 關閉連接
cluster.shutdown()
---------------------------------


# -*- coding: utf-8 -*-
import requests
import json


# 使用 GET 方式下載華城資料網頁
r = requests.get('http://113.196.140.150:9995/TPE1')
# 檢查狀態碼是否 OK
if r.status_code == requests.codes.ok:
        print("requests http://113.196.140.150:9995/TPE1    OK ")
# 輸出網頁 HTML 原始碼
TPE1 = json.loads(r.text)['TPE1']


class Cassandra_json():
        rows_to_insert = []
        bt_num = 28

        def __init__(self, request_data):
                self.Time = request_data['Time'].replace("/","-")
                self.BmsVoltage = request_data['BmsVoltage']
                self.BmsCurrent = request_data['BmsCurrent']
                self.BmsSOC = request_data['BmsSOC']
                self.BmsSOH = request_data['BmsSOH']
                self.BmsRelayStatus = request_data['BmsRelayStatus']
                self.BmsStatus = request_data['BmsStatus']
                self.batteryInfos = request_data['batteryInfos']
                self.v1 = None
                self.v2 = None
                self.v3 = None
                self.v4 = None
                self.v5 = None
                self.v6 = None
                self.v7 = None
                self.ch1_T = None
                self.ch2_T = None


        def __battery_value__(self, b_request_data):
                # BMU_1 = b_request_data['BMU_1']
                self.v1 = b_request_data['v1']
                self.v2 = b_request_data['v2']
                self.v3 = b_request_data['v3']
                self.v4 = b_request_data['v4']
                self.v5 = b_request_data['v5']
                self.v6 = b_request_data['v6']
                self.v7 = b_request_data['v7']
                self.ch1_T = b_request_data['ch1_T']
                self.ch2_T = b_request_data['ch2_T']


        def __BMU_json__(self, b_id):
                self.__battery_value__(self.batteryInfos[b_id])
                jsonStr = {"Time": self.Time, "BmsVoltage": self.BmsVoltage, "BmsCurrent": self.BmsCurrent, "BmsSOC": self.BmsSOC, "BmsSOH":self.BmsSOH,
                                        "BmsRelayStatus": self.BmsRelayStatus, "BmsStatus": self.BmsStatus, "BMU_ID": str(b_id + 1),
                                        "v1": self.v1, "v2": self.v2, "v3": self.v3, "v4": self.v4, "v5": self.v5, "v6": self.v6, "v7": self.v7,
                                        "ch1_T": self.ch1_T, "ch2_T": self.ch2_T}
                return jsonStr


        def get_json(self):
                for i in range(0, BigQ_json.bt_num):
                        BigQ_json.rows_to_insert.append(self.__BMU_json__(i))
                return BigQ_json.rows_to_insert

    
def main():
        Cassandra_object = Cassandra_json(TPE1[0])
        rows_to_insert = Cassandra_object.get_json()
        # print(rows_to_insert)

		
if __name__ == '__main__':
        main()


------------------------	
#為了可以持續頻繁接收每3~7秒一次更新的28筆資料，每個insert session包成一個API request。	
#華城對Alex的server IP丟資料，對Alex的程式來說，把資料轉發到Cassanfra的程式包成一支function。



from flask import Flask
from uwsgidecorators import postfork
from cassandra.cluster import Cluster

session = None
prepared = None

@postfork
def connect():
    global session, prepared
    session = Cluster().connect()
    prepared = session.prepare("SELECT release_version FROM system.local WHERE key=?")

app = Flask(__name__)

@app.route('/')
def server_version():
    row = session.execute(prepared, ('local',))[0]
    return row.release_version	
---------------------
>>> future = session.execute_async("SELECT * FROM system.local", trace=True)
>>> result = future.result()
>>> trace = future.get_query_trace()
>>> for e in trace.events:
>>>     print e.source_elapsed, e.description

0:00:00.000077 Parsing select * from system.local
0:00:00.000153 Preparing statement
0:00:00.000309 Computing ranges to query
0:00:00.000368 Submitting range requests on 1 ranges with a concurrency of 1 (279.77142 rows per range expected)
0:00:00.000422 Submitted 1 concurrent range requests covering 1 ranges
0:00:00.000480 Executing seq scan across 1 sstables for (min(-9223372036854775808), min(-9223372036854775808))
0:00:00.000669 Read 1 live and 0 tombstone cells
0:00:00.000755 Scanned 1 rows and matched 1
---------------------
>>> prepared = session.prepare("SELECT * FROM example.t WHERE key=?")
>>> bound = prepared.bind((1,))
>>> replicas = cluster.metadata.get_replicas(bound.keyspace, bound.routing_key)
>>> for h in replicas:
>>>   print h.address
127.0.0.1
127.0.0.2
---------------------
import json
   data=['Rakesh',{'marks':(50,60,70)}]
   s=json.dumps(data)
json.loads(s)
-------------------
