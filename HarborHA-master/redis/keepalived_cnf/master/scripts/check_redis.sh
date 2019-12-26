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
