#!/bin/bash
DIRNAME=`dirname $0`
RUNHOME=`cd $DIRNAME/; pwd`
LOG_DIR=${RUNHOME}/apiroute-works/logs

RETRY=$1
CURL_RETRY_COUNT=1
SLEEP_SECONDS_IN_MONITOR=2

if [ x"$HTTP_OVERWRITE_PORT" == x"" ]; then
	zenap_msb_apigateway_port=80	
else
	zenap_msb_apigateway_port=$HTTP_OVERWRITE_PORT
fi

#####################################
echo @RUNHOME@ $RUNHOME
echo @LOG_DIR@ $LOG_DIR
echo @zenap_msb_apigateway_port@ $zenap_msb_apigateway_port
#####################################
	
print_log() {   
    LOG_DATE=`date "+%Y-%m-%d %H:%M:%S"`
    LOG_FILENAME="$LOG_DIR/apigateway-monitor-`date "+%Y-%m-%d"`.log"    
    ######[$$]:进程id
    ###### $n 传递给脚本或函数的参数, $1为函数名称，默认为FUNCNAME[0] $2为日志输出内容 
    ###### $? 函数的返回值
    ###### FUNCNAME[@]: 函数调用栈，
    ###### {} 第一个空格后在写脚本
    { echo -e  $LOG_DATE "[$$]===${FUNCNAME[@]}=== $2 "; } >> $LOG_FILENAME
}

apigateway_monitor() {
    print_log ${FUNCNAME[0]} "Starting health check. ==$1==";
	if [ "$1" == "with-retry" ]; then
		CURL_RETRY_COUNT=5	
	fi
	
	for i in `seq $CURL_RETRY_COUNT`
	do
		success=0
		response=$(curl -m 2 -sL 127.0.0.1:$zenap_msb_apigateway_port/api/route/v1/health -w %{http_code} -o /dev/null)
		rc=$?
		if [ $rc -eq 0 ];then
			print_log "${FUNCNAME[0]}" "Info: Curl to health check endpoint succeeded.";
			success=1
			break
		else
			print_log "${FUNCNAME[0]}" "Warn: Curl to health check endpoint failed $i time(s).";
			if [ "$1" == "with-retry" ]; then
				sleep $SLEEP_SECONDS_IN_MONITOR
			fi			
		fi
	done
	
	if [ $success -eq 0 ];then
        print_log "${FUNCNAME[0]}" "Failed: health check Endpoint is unreachable.";
		return 1
    fi
    
    if [ $response -eq 200 ]; then
        print_log "${FUNCNAME[0]}" "health check succeed.";            
        return 0
    fi
	
    print_log "${FUNCNAME[0]}" "Failed:health check return non 200 but is : $response ";             
    return 1
}  

if [ ! -d "$LOG_DIR" ]; then
	echo there is no $LOG_DIR but will be created, You can ignore this message!
	mkdir -p "$LOG_DIR"
fi

apigateway_monitor $RETRY
