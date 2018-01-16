#!/bin/bash


RUNHOME=$(pwd)
echo "================ENV=================="
echo "#WORKSPACE# $WORKSPACE"
go version
echo "#RUNHOME# $RUNHOME"
echo "================ENV==================="

#check env
WORK_DIR=$RUNHOME/../work
mkdir -p $WORK_DIR

#clone build export gocoverage
GO_COVERAGE_DIR=$WORK_DIR/gocoveragedir

mkdir -p $GO_COVERAGE_DIR/src
cd $GO_COVERAGE_DIR/src

git clone http://10.89.168.136/10065132/gocoverage.git
export GOPATH=$GO_COVERAGE_DIR
echo "#GOPATH#: $GOPATH"
go install gocoverage/
export PATH=$PATH:$GO_COVERAGE_DIR/bin
echo "==============build gocoverage finish============"


#vmHost=10.63.240.167
#scp root@$vmHost:/root/ci/consul $WORK_DIR
#export PATH=$PATH:$WORK_DIR
#echo "==============prepare consul finish============"
echo "PATH : $PATH"

#execute gocoverage
GO_DIR=$RUNHOME/../../apiroute


export GOPATH=$GO_DIR
cd $GO_DIR/src

start=$(date +%s) 
gocoverage apiroute/ apiroute/vendor/ > allpkgcover
end=$(date +%s) 
duration=$(( $end - $start ))  
export duration


echo "============== gocoverage result============"
echo "duration: $duration"
tail allpkgcover
cp allpkgcover $WORKSPACE
echo "============== build result.txt============"


rm -rf $WORKSPACE/result*

cd $RUNHOME
export proCIShRepo=$WORKSPACE
export proname=sdclient-go-utest

chmod +x ut_statistics.sh && ./ut_statistics.sh  
echo "============== ut_statistics finish============"










