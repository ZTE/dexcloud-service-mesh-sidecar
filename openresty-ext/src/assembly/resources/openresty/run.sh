#!/bin/sh
#
# Copyright (C) 2016 ZTE, Inc. and others. All rights reserved. (ZTE)
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#         http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

DIRNAME=`dirname $0`
CONFIG_DIR=`cd $DIRNAME/nginx/msb-enabled; pwd`
HTTP_CONF="$CONFIG_DIR/msb.conf"
HTTPS_CONF="$CONFIG_DIR/msbhttps.conf"
HOME=`cd $DIRNAME/nginx; pwd`
_NGINXCMD="$HOME/sbin/nginx"
_LUAINIT="$HOME/luaext/conf/msbinit.lua"
_NGINXCONF="$HOME/conf/nginx.conf"
NUM_OF_WORKERS_LARGE=8 #in case scale = large
NUM_OF_WORKERS_SMALL=4 #in case scale = small
NUM_OF_WORKERS_TINY=2 #in case scale = tiny
LUAJIT_HOME=`cd $DIRNAME/luajit; pwd`
echo =========== prepare the symbolic links  ========================================
ln -s -f $_NGINXCMD $DIRNAME/bin/openresty
#ln -s -f $LUAJIT_HOME/bin/luajit-2.1.0-beta2 $LUAJIT_HOME/bin/luajit
ln -s -f $LUAJIT_HOME/lib/libluajit-5.1.so.2.1.0 $LUAJIT_HOME/lib/libluajit-5.1.so.2
ln -s -f $LUAJIT_HOME/lib/libluajit-5.1.so.2.1.0 $LUAJIT_HOME/lib/libluajit-5.1.so
echo ================================================================================

echo =========== create symbolic link for libluajit-5.1.so.2  ========================================
LUAJIT_HOME=`cd $DIRNAME/luajit; pwd`
LUAJIT_FILENAME="$LUAJIT_HOME/lib/libluajit-5.1.so.2"
LN_TARGET_FILE='/lib/libluajit-5.1.so.2'
LN_TARGET_FILE64='/lib64/libluajit-5.1.so.2'
ln -s -f $LUAJIT_FILENAME $LN_TARGET_FILE
ln -s -f $LUAJIT_FILENAME $LN_TARGET_FILE64
echo ===============================================================================

echo =========== openresty config info  =============================================
echo HOME=$HOME
echo _NGINXCMD=$_NGINXCMD
echo ===============================================================================
cd $HOME; pwd

# clean the consul-template generated files
rm -fr stream-enabled/service.conf
rm -fr stream-enabled/port.lst
rm -fr sites-enabled/service.conf
rm -fr sites-enabled/server_http.conf
rm -fr sites-enabled/server_https.conf

if [ -n "${APIGATEWAY_MODE}" -a -n "${APIGATEWAY_REDIS_PORT}" ]; then
        sed -i 's/= 6379/= '${APIGATEWAY_REDIS_PORT}'/g' $_LUAINIT
fi

# substitute listen ports with env specified ones if any
if [ -n "$HTTP_OVERWRITE_PORT" -a -f "$HTTP_CONF" ]; then
   # validate env port
   echo $HTTP_OVERWRITE_PORT | grep -E ^[0-9]+$ >/dev/null
   if [ $? -eq 0 ]; then
      sed -i -E 's/(.*listen +)[0-9]+(.+)/\1'$HTTP_OVERWRITE_PORT'\2/' $HTTP_CONF
   else
      echo "http listen port:$HTTP_OVERWRITE_PORT is not valid"
   fi
fi

if [ -n "$HTTPS_OVERWRITE_PORT" -a -f "$HTTPS_CONF" ]; then
   # validate env port
   echo $HTTPS_OVERWRITE_PORT | grep -E ^[0-9]+$ >/dev/null
   if [ $? -eq 0 ]; then
      sed -i -E 's/(.*listen +)[0-9]+(.+)/\1'$HTTPS_OVERWRITE_PORT'\2/' $HTTPS_CONF
   else
      echo "https listen port:$HTTPS_OVERWRITE_PORT is not valid"
   fi
fi


echo @SCALE@ ${SCALE}
# set worker_processes number by scale ,default is 4
if [ "${SCALE}" == "small" ]; then
   echo "set worker_processes to $NUM_OF_WORKERS_SMALL in ${SCALE} scenario."
	sed -i -E 's/(.*worker_processes +)([0-9]+)(.+)/\1'$NUM_OF_WORKERS_SMALL'\3/' $_NGINXCONF	
fi	

if [ "${SCALE}" == "tiny" ]; then
   echo "set worker_processes to $NUM_OF_WORKERS_TINY in ${SCALE} scenario."
   sed -i -E 's/(.*worker_processes +)([0-9]+)(.+)/\1'$NUM_OF_WORKERS_TINY'\3/' $_NGINXCONF 
fi 

if [ "${SCALE}" == "large" ]; then
   echo "set worker_processes to $NUM_OF_WORKERS_LARGE in ${SCALE} scenario."
   sed -i -E 's/(.*worker_processes +)([0-9]+)(.+)/\1'$NUM_OF_WORKERS_LARGE'\3/' $_NGINXCONF 
fi

echo @WORK_DIR@ $HOME
echo @C_CMD@ $_NGINXCMD -p $HOME/
$_NGINXCMD -p $HOME/

