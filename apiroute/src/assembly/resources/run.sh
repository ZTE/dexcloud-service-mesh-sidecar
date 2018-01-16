#!/bin/sh
DIRNAME=`dirname $0`
RUNHOME=`cd $DIRNAME/; pwd`


echo @RUNHOME@ $RUNHOME
cd $RUNHOME

# set CONSUL_IP in ha mode
if [ -n "${HA_MODE}" ]; then
        # export CONSUL_IP=`cat /root/consul.ip`
		export CONSUL_REGISTER_MODE="catalog"
fi

echo "\n\n### Starting apiroute-go";
$RUNHOME/apiroute




