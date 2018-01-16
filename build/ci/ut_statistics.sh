#!/bin/bash

resultfile=${proCIShRepo}/result.txt
resultfilecsv=${proCIShRepo}/result.csv
allpkgcover=${proCIShRepo}/allpkgcover
jenkinsws_src=${proCIShRepo}
if [ -f "${resultfile}" ];then
   rm -rf ${resultfile}
fi
totalcount=0
hitcount=0
coverrate=0
datetime=`date "+%Y-%m-%d %H:%M:%S"`
result=`cat ${allpkgcover} |tail -n 5 |head -n 1 |awk '{print $2}'`
if [ "${result}" = "success" ];then 
   totalcount=`cat ${allpkgcover} |tail -n 4 |head -n 1 |awk '{print $4}'`
   utcount=`cat ${allpkgcover} |tail -n 2 |head -n 1 |awk '{print $4}'`
   coverrate=`cat ${allpkgcover} |tail -n 1 |awk '{print $3}'`
   cover=${coverrate%\%*}
   #coverf=$(awk -v x=${cover} -v y=100 'BEGIN {printf "%.2f\n",x/y}')
   #hitcount=$(awk -v x=${coverf} -v y=${totalcount} 'BEGIN {printf "%d\n",x*y}')
   hitcount=`cat ${allpkgcover} |tail -n 3 |head -n 1 |awk '{print $4}'`
fi
#echo "data type:ut" >> ${resultfile}
echo "result:${result}" >> ${resultfile}
echo "total count:${totalcount}" >> ${resultfile}
echo "hit count:${hitcount}" >> ${resultfile}
echo "cover rate:${coverrate}" >> ${resultfile}
echo "duration:${duration}" >> ${resultfile}
echo "date time:${datetime}" >> ${resultfile}

if [ -f "${resultfile}" ];then
   filedatetime=`date "+%Y%m%d_%H%M%S"`
   #cp ${resultfile} ${jenkinsws_src}/result_$BUILD_NUMBER.txt
   cp ${resultfile} ${jenkinsws_src}/result_${filedatetime}.txt
fi

if [ -f "${resultfilecsv}" ];then
  echo "${proname},${result},${datetime},${totalcount},${hitcount},${cover},${duration}" >> ${resultfilecsv}
else
  echo "App,Result,Date,TotalLine,HitLine,Coverage,Duration" >> ${resultfilecsv}
  echo "${proname},${result},${datetime},${totalcount},${hitcount},${cover},${duration}" >> ${resultfilecsv}  
fi

#ops all pro
#${proCIShRepo}/statisticsops.sh
if [ "${result}"x = "fail"x ]; then
    echo "################################ut test fail################################## "
    exit 1
else
    echo "################################ut test end################################ "
fi
