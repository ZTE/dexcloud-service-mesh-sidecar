#!/bin/sh
#
# Copyright 2016 ZTE, Inc. and others.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#


DIRNAME=`dirname $0`
RUNHOME=`cd $DIRNAME/; pwd`
echo @RUNHOME@ $RUNHOME

echo @JAVA_HOME@ $JAVA_HOME

if [ "x$RELEASE_VERSION" != "x" ]; then
    DOCKER_RELEASE_VERSION=$RELEASE_VERSION    
else
	DOCKER_RELEASE_VERSION=latest
fi

echo @RELEASE_VERSION@ $RELEASE_VERSION
#build
#mvn clean install 

RELEASE_BASE_DIR=$RUNHOME/release
STANDALONE_RELEASE_DIR=${RELEASE_BASE_DIR}/standalone/dexcloud-mesh-sidecar
DOCKER_RELEASE_DIR=${RELEASE_BASE_DIR}/docker/dexcloud-mesh-sidecar

rm -rf $RELEASE_BASE_DIR
mkdir  $STANDALONE_RELEASE_DIR -p
mkdir  $DOCKER_RELEASE_DIR -p

DOCKER_RUN_NAME=dexcloud_mesh_sidecar
DOCKER_IMAGE_NAME=dexcloud_mesh_sidecar


VERSION_DIR=$RUNHOME/distributions/sidecar/target/version/

cp -r $VERSION_DIR/* ${STANDALONE_RELEASE_DIR}
rm -rf ${STANDALONE_RELEASE_DIR}/blueprint


#build docker image
cd ${STANDALONE_RELEASE_DIR}
cp  $RUNHOME/build/ci/build_docker_image.sh ${STANDALONE_RELEASE_DIR}
chmod 777 build_docker_image.sh

#clear old version
docker rmi ${DOCKER_IMAGE_NAME}:${DOCKER_RELEASE_VERSION}

./build_docker_image.sh -n=${DOCKER_IMAGE_NAME} -v=${DOCKER_RELEASE_VERSION} -d=$DOCKER_RELEASE_DIR
rm build_docker_image.sh


#docker run
docker images
#docker run -d --net=host  --name ${DOCKER_RUN_NAME} ${DOCKER_IMAGE_NAME}:${DOCKER_RELEASE_VERSION}
#docker ps |grep ${DOCKER_RUN_NAME} 




