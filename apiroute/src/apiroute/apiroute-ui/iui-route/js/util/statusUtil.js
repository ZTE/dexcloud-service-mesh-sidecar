var statusUtil = {};

//轮询时间
statusUtil.statisticsPollTime=60000;
statusUtil.connectionsPollTime=5000;

//横坐标最大显示数
statusUtil.connectionsXAxisCount=12;
statusUtil.statisticsXAxisCount=60;

statusUtil.connection=true;

statusUtil.initChart= function(){
	statusUtil.init_statistics_requestChart();
	statusUtil.init_status_requestChart();
	statusUtil.init_connectionChart();

     if(statusUtil.connection==true){    
          statusUtil.getRealTimeChartData();

            setInterval(function () {
            statusUtil.getRealTimeChartData();
         }, statusUtil.connectionsPollTime);   
      }

}

statusUtil.getQueryString= function(name) { 
    var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i"); 
    var r = window.location.search.substr(1).match(reg); 
    if (r != null) 
        return unescape(r[2]);
    return null; 
} 




/*
responses-json

{
"ip" : [
{
stats_time:"2017-06-19T01:28:30Z",
"responses" :{ "total" : 20943825 }
}, 
{
stats_time:"2017-06-19T01:29:00Z",
"responses" :{ "total" : 20943825 }
}, 
{
stats_time:"2017-06-19T01:29:30Z",
"responses" :{ "total" : 20943825 }
}
]
}

*/

statusUtil.init_statistics_requestChart= function(){


	var times = [];
	var datas=[];
	var latestTime="";
	$.ajax({
        "type": 'get',
        "async": false,
        "timeout" : 5000, 
        "url":  apiBasePath+"/stats/apigateway/responses?latestNum=10",
        "dataType": "json",
        success: function (resp) { 
         var respObj = (resp==null)?{}:resp;  
     
         $.each(respObj, function(key, val) { 
           var initDatas=val;
           for(var i=0;i<initDatas.length;i++){
           
           	times.push(statusUtil.toLocaleTime(initDatas[i].stats_time));
           	datas.push(initDatas[i].responses.total);
           }

           if(initDatas.length>0){
           	latestTime=initDatas[initDatas.length-1].stats_time;
        	}
        });

        			                                       	
        },
         error: function(XMLHttpRequest, textStatus, errorThrown) {
           statusUtil.connection=false;
           routeUtil.notify('get Chart Data fails：',XMLHttpRequest.statusText,'danger'); 		
         }
    });


	var requestChart = echarts.init(document.getElementById('statisticsLineChartDiv'), 'macarons');

        var option = {

             tooltip : {
		        trigger: 'axis'
		    },
          	toolbox: {
		        feature: {
		            saveAsImage: {
                  name:$.i18n.prop('org_openo_msb_status_chart_statisticsLineChart'),
                   title : $.i18n.prop('org_openo_msb_status_chart_save_picture')
                }
		        }
		    },
            legend: {
                data:[$.i18n.prop('org_openo_msb_status_chart_minute_requests')]
            },             
		xAxis: 
        {
            type: 'category',
            boundaryGap: false,
            data: times,
            nameTextStyle:{
              color:'#8D8D8D'
            },
            name: $.i18n.prop('org_openo_msb_status_chart_time')

        },
		yAxis:
        {
            type: 'value',
            scale: true,
            name: $.i18n.prop('org_openo_msb_status_chart_requests'),
            nameTextStyle:{
              color:'#8D8D8D'
            },
            minInterval: 1,
            min: 0,
            boundaryGap: [0.1, 0.1]
        },
        series: [
	        {
	            name:$.i18n.prop('org_openo_msb_status_chart_minute_requests'),
	            type:'line',
	            stack: 'status',
              areaStyle: {normal: {
                color: {
                    type: 'linear',
                    x: 0,
                    y: 0,
                    x2: 0,
                    y2: 1,
                    colorStops: [{
                        offset: 0, color: '#88D8FF' // 0% 处的颜色
                    }, {
                        offset: 1, color: '#FFFFFF' // 100% 处的颜色
                    }],
                    globalCoord: false 
                }
              }},
	            data:datas
	        }	        
        ]
 	};

        // 使用刚指定的配置项和数据显示图表。
         var maxVaule=Math.max.apply(null, datas);
           if(maxVaule<5){
            option.yAxis.max=5;
           }
          
        requestChart.setOption(option);
        window.onresize = requestChart.resize;
 if(statusUtil.connection==true){
      setInterval(function () {

      $.ajax({
        "type": 'get',
        "async": true,
        "timeout" : 3000, 
        "url":  apiBasePath+"/stats/apigateway/responses",
        "dataType": "json",
        success: function (resp) { 
        var respObj = (resp==null)?{}:resp;  
       
        $.each(respObj, function(key, val) {            
        var latestData=val;    
         
         if(latestData.length>0){        
         	
         	var data = option.series[0].data;
          var latest_stats_time=latestData[0].stats_time;
    		
    		if(latestTime!=latest_stats_time){	 
    			  data.push(latestData[0].responses.total);    			  
				    option.xAxis.data.push(statusUtil.toLocaleTime(latest_stats_time));
      			if(data.length>=statusUtil.statisticsXAxisCount){
			    	  data.shift();			    	
			    	  option.xAxis.data.shift();
				    }


           var maxVaule=Math.max.apply(null, data);
           if(maxVaule<5){
            option.yAxis.max=5;
           }
           else{
            option.yAxis.max=null;
           }  
				  requestChart.setOption(option);

          latestTime=latest_stats_time;
  			}
		                                       	
      }
    });
   }
});

		   
  }, statusUtil.statisticsPollTime);  
 }
}

var statusLineChart;
var statusLineChartOption;
statusUtil.init_status_requestChart= function(){

statusLineChart = echarts.init(document.getElementById('statusLineChartDiv'), 'macarons');

statusLineChartOption = {
             color:[ "#A1DF6A", "#40C0FF", "#F17874"],           
             tooltip : {
		        trigger: 'axis'
		    },
          	toolbox: {
		        feature: {
		            saveAsImage: {
                  name:$.i18n.prop('org_openo_msb_status_chart_statusLineChart'),
                  title : $.i18n.prop('org_openo_msb_status_chart_save_picture')
                }
		        }
		    },
            legend: {
                data:[$.i18n.prop('org_openo_msb_status_chart_forward_waiting_response'),$.i18n.prop('org_openo_msb_status_chart_accept_preparing_forward'),$.i18n.prop('org_openo_msb_status_chart_receive_resp_not_return')],
                itemGap:20
            },             
		xAxis: 
        {
            type: 'category',
            boundaryGap: false,
            name: $.i18n.prop('org_openo_msb_status_chart_time'),
            nameTextStyle:{
              color:'#8D8D8D'
            },
            data: []
        },
		yAxis:
        {
            type: 'value',
            scale: true,
            name: $.i18n.prop('org_openo_msb_status_chart_requests'),
            nameTextStyle:{
              color:'#8D8D8D'
            },
            min: 0,
            minInterval: 1,
            boundaryGap: [0.1, 0.1]
        },
        series: [
	        {
	            name:$.i18n.prop('org_openo_msb_status_chart_forward_waiting_response'),
	            type:'line',
	            data:[]
	        },
	        {
	            name:$.i18n.prop('org_openo_msb_status_chart_accept_preparing_forward'),
	            type:'line',
	            data:[]
	        },
	        {
	            name:$.i18n.prop('org_openo_msb_status_chart_receive_resp_not_return'),
	            type:'line',
	            data:[]
	        }
        
        ]
 	};

        // 使用刚指定的配置项和数据显示图表。
        statusLineChart.setOption(statusLineChartOption);
        window.onresize = statusLineChart.resize;
}


var connectionChart;
var connectionChartOption;

statusUtil.init_connectionChart= function(){


	  connectionChart = echarts.init(document.getElementById('connectionBarChartDiv'), 'macarons');
     
       
     connectionChartOption = {
            color:[ "#A2E9FF", "#A1DF6A", "#53C6FF","#F67D79"],           
            tooltip : {
  		        trigger: 'axis',
  		        axisPointer : {           
  		            type : 'shadow'       
  		        }
		        },
          	toolbox: {
		        feature: {
		            saveAsImage: {
                  name:$.i18n.prop('org_openo_msb_status_chart_connectionBarChart'),
                   title : $.i18n.prop('org_openo_msb_status_chart_save_picture')
                }
		        }
		    },
            legend: {
                data:['Active','Waiting','Writing','Reading'],
                itemGap:20
            },             
		xAxis: 
        {
            type: 'category',
            boundaryGap: true,
            name: $.i18n.prop('org_openo_msb_status_chart_time'),
            nameTextStyle:{
              color:'#8D8D8D'
            },           
            data: []
        },
		yAxis:
        {
            type: 'value',
            scale: true,
            name: $.i18n.prop('org_openo_msb_status_chart_connects'),
            nameTextStyle:{
              color:'#8D8D8D'
            },
            min: 0,
            minInterval: 1,
            boundaryGap: [0.1, 0.1]
          
        },
        series: [
        {
            name:'Active',
            type:'bar',
            barWidth: '12',
            itemStyle:{
              normal:{
                barBorderRadius:[3, 3, 0, 0]
              }
            },
            data:[]
        },
        {
            name:'Waiting',
            type:'bar',
            barWidth: '12',            
            stack: 'connection',
            data:[]
        },
        {
            name:'Writing',
            type:'bar',
            barWidth: '12',
            itemStyle:{
              normal:{
                barBorderRadius:[3, 3, 0, 0]
              }
            },
            stack: 'connection',
            data:[]
        },
        {
            name:'Reading',
            type:'bar',
            barWidth: '12',
            itemStyle:{
              normal:{
                barBorderRadius:[3, 3, 0, 0]
              }
            },
            stack: 'connection',
            data:[]
        }
        ]
 	};

        // 使用刚指定的配置项和数据显示图表。
        connectionChart.setOption(connectionChartOption);
        window.onresize = connectionChart.resize;       
   
}

/*
json:real-time
{
  "nginx_version" : "1.7.9",
  "ngx_lua_version" : "0.9.13",
  "worker_count":8,
  "connections" :
  { "active" : "19", "reading" : "0", "waiting" : "18", "writing" : "1" },
  "requests" :
  { "accept_preparing_forward" : "19", "forward_waiting_response" : "0", "receive_resp_not_return" : "18" }
}


 */


statusUtil.getRealTimeChartData=function(){
  $.ajax({
        "type": 'get',
        "async": true,
        "url":  apiBasePath+"/real-time/apigateway",
        "dataType": "json",
        success: function (resp) { 
   
         if(resp!=null){
          var currentTime=statusUtil.getCurrentTime();
            //requestsChartData
          var forward_data= statusLineChartOption.series[0].data;
          var accept_data = statusLineChartOption.series[1].data;
          var receive_data = statusLineChartOption.series[2].data;
      


          if(accept_data.length>=statusUtil.statisticsXAxisCount){
              accept_data.shift();
              forward_data.shift();
              receive_data.shift();
              statusLineChartOption.xAxis.data.shift();
          }
          
           
           accept_data.push(resp.requests.accept_preparing_forward);           
           forward_data.push(resp.requests.forward_waiting_response);
           receive_data.push(resp.requests.receive_resp_not_return);

            var allValue=accept_data.concat(forward_data).concat(receive_data);
           var maxVaule=Math.max.apply(null, allValue);
           if(maxVaule<5){
            statusLineChartOption.yAxis.max=5;
           }
           else{
            statusLineChartOption.yAxis.max=null;
           }
          
       
          statusLineChartOption.xAxis.data.push(currentTime);

          statusLineChart.setOption(statusLineChartOption);

            //connectionChartData
            var active_data =  connectionChartOption.series[0].data;
            var waiting_data = connectionChartOption.series[1].data;
            var writing_data = connectionChartOption.series[2].data;
            var reading_data = connectionChartOption.series[3].data;


          if(active_data.length>=statusUtil.connectionsXAxisCount){
              active_data.shift();
              waiting_data.shift();
              writing_data.shift();
              reading_data.shift();
              connectionChartOption.xAxis.data.shift();
          }
          
           active_data.push(resp.connections.active); 
           waiting_data.push(resp.connections.waiting);           
           writing_data.push(resp.connections.writing);
           reading_data.push(resp.connections.reading);

           var allValue=active_data.concat(waiting_data).concat(writing_data).concat(reading_data)
           var maxVaule=Math.max.apply(null, allValue)
           if(maxVaule<5){
            connectionChartOption.yAxis.max=5;
           }
           else{
            connectionChartOption.yAxis.max=null;
           }


          connectionChartOption.xAxis.data.push(currentTime);

          connectionChart.setOption(connectionChartOption);
                                            
        }
      }
    });
}


statusUtil.getCurrentTime=function(){
          var date = new Date();          
          var axisData = [statusUtil.addZero(date.getHours()),statusUtil.addZero(date.getMinutes()),statusUtil.addZero(date.getSeconds())].join(":");
          return axisData;
  }

statusUtil.addZero=function(s) {
        return s < 10 ? '0' + s: s;
    }

statusUtil.toLocaleTime=function(timestamp){
  return new Date(timestamp).Format("yyyy-MM-dd hh:mm:ss");
}

statusUtil.toLocaleHours=function(timestamp){
  return new Date(timestamp).Format("hh:mm:ss");
}

statusUtil.getLocalTime=function(nS) {  
          var date = new Date(parseInt(nS)*1000);    
          var fullDate= [statusUtil.addZero(date.getFullYear()),statusUtil.addZero(date.getMonth() + 1 ) ,statusUtil.addZero(date.getDate())].join("-");
              fullDate+=" ";   
              fullDate += [statusUtil.addZero(date.getHours()),statusUtil.addZero(date.getMinutes()),statusUtil.addZero(date.getSeconds())].join(":");
          return fullDate;

}



Date.prototype.Format = function (fmt) { //author: meizz 
    var o = {
        "M+": this.getMonth() + 1, //月份 
        "d+": this.getDate(), //日 
        "h+": this.getHours(), //小时 
        "m+": this.getMinutes(), //分 
        "s+": this.getSeconds(), //秒 
        "q+": Math.floor((this.getMonth() + 3) / 3), //季度 
        "S": this.getMilliseconds() //毫秒 
    };
    if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
    for (var k in o)
    if (new RegExp("(" + k + ")").test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
    return fmt;
}