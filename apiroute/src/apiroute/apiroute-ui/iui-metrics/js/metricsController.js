/*
 * Copyright (C) 2016 ZTE, Inc. and others. All rights reserved. (ZTE)
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

var vm = avalon
		.define({
			$id : "metricsController",
			routeTargetTitle:$.i18n.prop("org_openo_msb_route_content_title"),		
			$metricsUrl : '/admin/metrics',		
			metricsLoading:false,		
			metricsArray :  [],	
			threadNum:"",
			routeLoading:false,	
			loadingTip:"",
			jvmTime:"",
			isShowJVM:false,
			isShowRequestsChart:false,
			restArray:[],
			metersArray:[],		
        	 $dataTableLanguage: {
                "sProcessing": "<img src='../img/loading-spinner-grey.gif'/><span>&nbsp;&nbsp;Loadding...</span>",   
                "sLengthMenu": $.i18n.prop("org_openo_msb_route-table-sLengthMenu"),
                "sZeroRecords": $.i18n.prop("org_openo_msb_route-table-sZeroRecords"),
                "sInfo": "<span class='seperator'>  </span>" + $.i18n.prop("org_openo_msb_route-table-sInfo"),
                "sInfoEmpty": $.i18n.prop("org_openo_msb_route-table-sInfoEmpty"),
                "sGroupActions": $.i18n.prop("org_openo_msb_route-table-sGroupActions"),
                "sAjaxRequestGeneralError": $.i18n.prop("org_openo_msb_route-table-sAjaxRequestGeneralError"),
                "sEmptyTable": $.i18n.prop("org_openo_msb_route-table-sEmptyTable"),
                "oPaginate": {
                    "sPrevious": $.i18n.prop("org_openo_msb_route-table-sPrevious"),
                    "sNext": $.i18n.prop("org_openo_msb_route-table-sNext"),
                    "sPage": $.i18n.prop("org_openo_msb_route-table-sPage"),
                    "sPageOf": $.i18n.prop("org_openo_msb_route-table-sPageOf")
                },
                "sSearch": $.i18n.prop("org_openo_msb_route-table-search"),
                "sInfoFiltered": $.i18n.prop("org_openo_msb_route-table-infofilter") 
            },	
            initMeters:function(testRestJson){
            
 				 //Initialize the Meter access list in detail
	          $.each(testRestJson,function(name,value) {
	          	
	           	var obj=value;
	           	obj.name=name;
				vm.metersArray.push(obj);
				});

	         	$('#metersTable').DataTable({
			      responsive: true,
				  destroy: true,
				  "dom": '<"top">frt<"bottom"lip><"clear">',
				  "oLanguage": vm.$dataTableLanguage
				});

				if(vm.metersArray.length>0){
	         	var unit=vm.metersArray[0].units;
				$('#org_openo_msb_metrics_meters_table_m1').text($('#org_openo_msb_metrics_meters_table_m1').text()+"("+unit+")");
				$('#org_openo_msb_metrics_meters_table_m5').text($('#org_openo_msb_metrics_meters_table_m5').text()+"("+unit+")");
				$('#org_openo_msb_metrics_meters_table_m15').text($('#org_openo_msb_metrics_meters_table_m15').text()+"("+unit+")");

				}
            },
            initTimersMetrics:function(testRestJson){

 				//Initialize the Rest interface traffic map
 				var restMetrics_data=[];
 				var restMetrics_name=[];	
 				var restArray=[];

 				 //Initialize the HTTP traffic
	           if(testRestJson["io.dropwizard.jetty.MutableServletContextHandler.get-requests"]!=null){
	           		vm.isShowRequestsChart=true;
	           		
		          
		            var requestsMetrics_data=[];
 					var requestsMetrics_name=["get","post","put","delete","other"];

 					 for(var i=0;i<requestsMetrics_name.length;i++){
 				//   	restMetrics_data.restName.push(restArray[i].name);
					// restMetrics_data.restCount.push(restArray[i].count);
					
					var data=new Object;
					data.name=requestsMetrics_name[i];
					switch(requestsMetrics_name[i]){
						case "get":data.data=[testRestJson["io.dropwizard.jetty.MutableServletContextHandler.get-requests"].count];break;
						case "post":data.data=[testRestJson["io.dropwizard.jetty.MutableServletContextHandler.post-requests"].count];break;
						case "put":data.data=[testRestJson["io.dropwizard.jetty.MutableServletContextHandler.put-requests"].count];break;
						case "delete":data.data=[testRestJson["io.dropwizard.jetty.MutableServletContextHandler.delete-requests"].count];break;
						case "other":data.data=[testRestJson["io.dropwizard.jetty.MutableServletContextHandler.other-requests"].count];break;
					

					}
					
					data.type='bar';
					data.itemStyle={ normal: {label : {show: true, position: 'top'},barBorderRadius:[3, 3, 0, 0]}};
					data.barWidth='25';
					data.barGap='200%';


		
					requestsMetrics_data.push(data);
					
 				  }


		       
		           metricsChart.requestsMetrics(requestsMetrics_data,requestsMetrics_name);
	       		}
	       		else{
	       			$("#restChartPanel").css("width","100%");	
	       		}

 				  $.each(testRestJson,function(name,value) {
 				  	if(name.indexOf("dropwizard") > 0 || name.indexOf("jetty") > 0)return true;
 				
		           	var nameArray=name.split(".");

		           	var rest=new Object();
		           	rest.name=nameArray[nameArray.length-1];
		           	rest.count=value.count;
		           	restArray.push(rest);
				});

 				   restArray.sort(function(a,b){
            		return b.count-a.count});
 				

 				  var restMaxNum=restArray.length>10?10:restArray.length;

 				  var barWidth,barGap;
 				  if(vm.isShowRequestsChart==true){
 				  	barWidth='20';
 				  	barGap='100%';
 				  }
 				  else{
 				  	barWidth='30';
 				  	barGap='150%';
 				  }
 				  
 				  for(var i=0;i<restMaxNum;i++){
 				//   	restMetrics_data.restName.push(restArray[i].name);
					// restMetrics_data.restCount.push(restArray[i].count);
					restMetrics_name.push(restArray[i].name);

					var data=new Object;
					data.name=restArray[i].name;
					data.data=[restArray[i].count];
					data.type='bar';
					data.itemStyle={ normal: {label : {show: true, position: 'top'},barBorderRadius:[3, 3, 0, 0]}};
					data.barWidth=barWidth;
					data.barGap=barGap;


		
					restMetrics_data.push(data);
					
 				  }
 				  

	         


	          

	       		  metricsChart.restMetrics(restMetrics_data,restMetrics_name); 

	           //Initialize the HTTP access list in detail
	          $.each(testRestJson,function(name,value) {
	          	if(name.indexOf("org.eclipse.jetty.server.HttpConnectionFactory") == 0) return true;
	           	var obj=value;
	           	obj.name=name;
				vm.restArray.push(obj);
				});

	         	$('#restTable').DataTable({
			      responsive: true,
				  destroy: true,
				  "dom": '<"top">frt<"bottom"lip><"clear">',
				  "oLanguage": vm.$dataTableLanguage
				});

	         	if(vm.restArray.length>0){
	         		var unit=vm.restArray[0].rate_units;
					$('#org_openo_msb_metrics_http_table_m1').text($('#org_openo_msb_metrics_http_table_m1').text()+"("+unit+")");
					$('#org_openo_msb_metrics_http_table_m5').text($('#org_openo_msb_metrics_http_table_m5').text()+"("+unit+")");
				}
            },	
            initGaugesMetrics:function(gaugesJson){
            	 //jvm Time
	           var jvmTime=gaugesJson["jvm.attribute.uptime"].value;

	           vm.jvmTime=metricsUtil.formatSeconds(jvmTime);


	           //Initialize the JVM memory usage
	           var Eden_Space_usage;
	           	if(gaugesJson["jvm.memory.pools.Eden-Space.usage"]==null){
	           		if(gaugesJson["jvm.memory.pools.PS-Eden-Space.usage"]==null)
	           		{
	           			Eden_Space_usage=0;
	           		}
	           		else{
	           			Eden_Space_usage=gaugesJson["jvm.memory.pools.PS-Eden-Space.usage"].value;
	           		}
	           	}
	           	else{
	           		Eden_Space_usage=gaugesJson["jvm.memory.pools.Eden-Space.usage"].value;
	           	}




	           	var Perm_Gen_usage;
	           	if(gaugesJson["jvm.memory.pools.Perm-Gen.usage"]==null){
	           		if(gaugesJson["jvm.memory.pools.PS-Perm-Gen.usage"]==null)
	           		{
	           			Perm_Gen_usage=0;
	           		}
	           		else{
	           			Perm_Gen_usage=gaugesJson["jvm.memory.pools.PS-Perm-Gen.usage"].value;
	           		}
	           	}
	           	else{
	           		Perm_Gen_usage=gaugesJson["jvm.memory.pools.Perm-Gen.usage"].value;
	           	}


	           	var Survivor_Space_usage;
	           	if(gaugesJson["jvm.memory.pools.Survivor-Space.usage"]==null){
	           		if(gaugesJson["jvm.memory.pools.PS-Survivor-Space.usage"]==null)
	           		{
	           			Survivor_Space_usage=0;
	           		}
	           		else{
	           			Survivor_Space_usage=gaugesJson["jvm.memory.pools.PS-Survivor-Space.usage"].value;
	           		}
	           	}
	           	else{
	           		Survivor_Space_usage=gaugesJson["jvm.memory.pools.Survivor-Space.usage"].value;
	           	}


	           	var Tenured_Gen_usage;
	           	if(gaugesJson["jvm.memory.pools.Tenured-Gen.usage"]==null){
	           		if(gaugesJson["jvm.memory.pools.PS-Old-Gen.usage"]==null)
	           		{
	           			Tenured_Gen_usage=0;
	           		}
	           		else{
	           			Tenured_Gen_usage=gaugesJson["jvm.memory.pools.PS-Old-Gen.usage"].value;
	           		}
	           	}
	           	else{
	           		Tenured_Gen_usage=gaugesJson["jvm.memory.pools.Tenured-Gen.usage"].value;
	           	}	


	           var memoryPieMetrics_data={
	           	CodeCache:(gaugesJson["jvm.memory.pools.Code-Cache.usage"].value*100).toFixed(1),
	           	EdenSpace:(Eden_Space_usage*100).toFixed(1),
	           	PermGen:(Perm_Gen_usage*100).toFixed(1),
	           	SurvivorSpace:(Survivor_Space_usage*100).toFixed(1),
	           	TenuredGen:(Tenured_Gen_usage*100).toFixed(1)
	           	};
	           metricsChart.memoryPieMetrics(memoryPieMetrics_data);

	           // initialize the JVM memory map
	           var heap_init=Math.round(gaugesJson["jvm.memory.heap.init"].value/1000000);
	           var non_heap_init=Math.round(gaugesJson["jvm.memory.non-heap.init"].value/1000000);

	           var heap_used=Math.round(gaugesJson["jvm.memory.heap.used"].value/1000000);
	           var non_heap_used=Math.round(gaugesJson["jvm.memory.non-heap.used"].value/1000000);

	           var heap_max=Math.round(gaugesJson["jvm.memory.heap.max"].value/1000000);
	           var non_heap_max=Math.round(gaugesJson["jvm.memory.non-heap.max"].value/1000000);

	           var memoryBarMetrics_data={
	           	init:[
	           		heap_init,
	           		non_heap_init,
	           		non_heap_init+heap_init
	           		],
	           	used:[
	           		heap_used,
	           		non_heap_used,
	           		non_heap_used+heap_used
	           	   	],
	           	max:[
	           		heap_max,
	           		non_heap_max,
	           		non_heap_max+heap_max
	           		]
	           };
	           metricsChart.memoryBarMetrics(memoryBarMetrics_data);


	             //Initializes the thread profile
	           var threadsMetrics_data= [{value:gaugesJson["jvm.threads.runnable.count"].value, name:'Runnable'},
                {value:gaugesJson["jvm.threads.timed_waiting.count"].value, name:'Timed waiting'},
                {value:gaugesJson["jvm.threads.waiting.count"].value, name:'Waiting'},
                {value:gaugesJson["jvm.threads.blocked.count"].value, name:'Blocked'}];
                vm.threadNum=gaugesJson["jvm.threads.count"].value;
 				metricsChart.threadsMetrics(threadsMetrics_data);
            },
			initMetrics : function() {
 			
		 var fullUrl= window.location.search.substr(1);
		 var publish_protocol=metricsUtil.getQueryString(fullUrl,"publish_protocol").replace(/<[^>]+>/g,"");
		 var publish_port=metricsUtil.getQueryString(fullUrl,"publish_port").replace(/<[^>]+>/g,"");
		 var url=metricsUtil.getQueryString(fullUrl,"url").replace(/<[^>]+>/g,"");
		 var ip=window.location.host.split(":")[0];
		 var metricsUrl=publish_protocol+"://"+ip+":"+publish_port+url;
		 
		//metricsUrl="./metrics.json";
		    	

		vm.routeLoading=true;	
		vm.loadingTip="<div class='loadingImg'></div>  "+$.i18n.prop("org_openo_msb_metrics_loading");

		 $.ajax({
            "type": 'get',
            "url": metricsUrl,
            timeout: 8000,
            "dataType": "json",
            success: function (resp) { 
            	   vm.routeLoading=false;	
	               var restJson = resp;  
	               
	                if(restJson.meters!=null){
	               	vm.initMeters(restJson.meters);
	               }
	               
	               if(restJson.timers!=null){
	               	vm.initTimersMetrics(restJson.timers);
	               }

	               

	               if(restJson.gauges!=null){
	               	vm.isShowJVM=true;
	               	vm.initGaugesMetrics(restJson.gauges);

	               }

	              

	           
				 },
		        error: function(XMLHttpRequest, textStatus, errorThrown) {
		        	if(XMLHttpRequest.status==502){
		        		vm.loadingTip=XMLHttpRequest.responseText;
		        	}
		        	else if(XMLHttpRequest.status==404){
		        		vm.loadingTip="<img  src='../iui-route/img/icon_card_no_data.png' class='card-image'/>"+$.i18n.prop("org_openo_msb_metrics_loading_404");
		        	}
		        	else{
				    	vm.loadingTip="<img  src='../iui-route/img/icon_card_no_data.png' class='card-image'/>"+$.i18n.prop("org_openo_msb_metrics_loading_fail")+XMLHttpRequest.statusText+"<center>"+metricsUrl+"</center>";             
		           }
		          }
		       });

			}			
			

	});
avalon.scan();
vm.initMetrics();

