#!/bin/sh
DIRNAME=`dirname $0`
RUNHOME=`cd $DIRNAME/; pwd`
echo @RUNHOME@ $RUNHOME


echo ================env info before runing===========================================
echo @OPENPALETTE_NAMESPACE@ $OPENPALETTE_NAMESPACE
echo @OPENPALETTE_MSB_IP@ $OPENPALETTE_MSB_IP
echo =================================================================================


if [ "x$OPENPALETTE_NAMESPACE" != "x" ] && [ "x$OPENPALETTE_MSB_IP" != "x" ]; then
      export SDCLIENT_IP=$OPENPALETTE_MSB_IP
      export ROUTE_WAY=ip
      export CUSTOM_FILTER_CONFIG=namespace:$OPENPALETTE_NAMESPACE
      export ROUTE_LABELS=visualRange:1
fi

echo ================env info after runing===========================================
echo @SDCLIENT_IP@ $SDCLIENT_IP
echo @ROUTE_WAY@ $ROUTE_WAY
echo @CUSTOM_FILTER_CONFIG@ $CUSTOM_FILTER_CONFIG
echo @ROUTE_LABELS@ $ROUTE_LABELS
echo @HTTP_OVERWRITE_PORT@ $HTTP_OVERWRITE_PORT
echo @HTTPS_OVERWRITE_PORT@ $HTTPS_OVERWRITE_PORT
echo @APIGATEWAY_APIROUTE_PORT@ $APIGATEWAY_APIROUTE_PORT
echo @APIGATEWAY_REDIS_PORT@ $APIGATEWAY_REDIS_PORT
echo @UPSTREAM_DNS_SERVERS@ $UPSTREAM_DNS_SERVERS
echo @SCALE@ $SCALE
echo =================================================================================

echo "### Starting redis";
cd $RUNHOME/redis
./run.sh 


echo "\n\n### Starting openresty";
cd $RUNHOME/openresty
./run.sh 

echo "\n\n### Starting apiroute"
cd $RUNHOME/apiroute
./run.sh