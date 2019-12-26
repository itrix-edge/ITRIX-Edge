#!/bin/bash
REDISCLI="/usr/local/bin/redis-cli"
LOGFILE="/var/log/keepalived-redis-state.log"
echo "[backup]" >> $LOGFILE
date >> $LOGFILE
echo "Being slave...." >> $LOGFILE 2>&1
sleep 15 
echo "Run SLAVEOF cmd ..." >> $LOGFILE
$REDISCLI SLAVEOF 192.168.99.11 6379 >> $LOGFILE  2>&1
