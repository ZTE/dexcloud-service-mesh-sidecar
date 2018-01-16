#!/bin/bash

#check env
go version
RUNHOME=$(pwd)
echo "#RUNHOME# $RUNHOME"

#go vet
cd $RUNHOME/../../apiroute/src/apiroute
go tool vet $(find . -name "*.go" | grep -v vendor | uniq)

echo "========go vet finish==========="
#golint
cd $RUNHOME/../../apiroute/src/apiroute
for pkg in $(go list ./... |grep -v /vendor/) ; do \
realpkg=${pkg#*_}
golint $realpkg
done

echo "======== golint finish=========="






