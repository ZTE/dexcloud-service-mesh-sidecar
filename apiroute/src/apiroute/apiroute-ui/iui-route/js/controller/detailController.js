var vm = avalon
		.define({
			$id : "detailController",
			$routeInstanceUrl :routeBasePath+'/routes/internal',
			$serviceInstanceUrl :routeBasePath+'/routes',
			route:{				
				routeIP:""
			},
			routeInfo : {				
				serviceName : "",
				serviceVersion : "",
				routeName:"",
				routeVersion:"",
				routeType:"",
				publishUrl:"",
				status:"",
				namespace:"",
				updateTimestamp:"",
				host:"",
				url:"",
				metricsUrl:"",
				apiJson:"",
				lb_policy:"",
				publish_port:"",
				publish_protocol:"",
				enable_ssl:"",
				enable_refer_match:"",
				connect_timeout:"",
				send_timeout:"",
				read_timeout:"",
				proxy_timeout:"",
				proxy_responses:"",
				enable_tls:"",
				nodes: []
			},
			viewHttpRoute:function(){

			 vm.route.routeIP=window.location.host.split(":")[0];

			 var url= window.location.search.substr(1);
		 	 var publishPort=routeUtil.getQueryString(url,"publishPort");
		 	 var routeType=routeUtil.getQueryString(url,"routetype");
		 	 var routeName=routeUtil.getQueryString(url,"routename");
		 	 var routeVersion=routeUtil.getQueryString(url,"routeversion");	

		 	

   			$.ajax({
	                "type": 'get',
	                 "url":  vm.$routeInstanceUrl+"?publishport="+publishPort+"&routetype="+routeType+"&routename="+routeName+"&routeversion="+routeVersion,
	                //"url":  "./data/apiRouteInfo.json",
	                timeout: 8000,
	                "dataType": "json",
	                success: function (responseRouteInfo) {  
	                     
				vm.routeInfo.routeName=routeName;
				vm.routeInfo.routeVersion= routeVersion;
				vm.routeInfo.serviceName=responseRouteInfo.metadata.serviceName;
			    vm.routeInfo.serviceVersion= responseRouteInfo.metadata.serviceVersion;
				vm.routeInfo.routeType= routeType;
				vm.routeInfo.status= responseRouteInfo.status;
				vm.routeInfo.namespace= responseRouteInfo.metadata.namespace==""?"default":responseRouteInfo.metadata.namespace;
				vm.routeInfo.updateTimestamp= responseRouteInfo.metadata.updateTimestamp;
				vm.routeInfo.host= responseRouteInfo.spec.host;
				vm.routeInfo.url= responseRouteInfo.spec.url==""?"/":responseRouteInfo.spec.url;
				vm.routeInfo.metricsUrl= responseRouteInfo.spec.metricsUrl;
				vm.routeInfo.apiJson= responseRouteInfo.spec.apijson;
				vm.routeInfo.lb_policy= responseRouteInfo.spec.lb_policy==""?$.i18n.prop("org_openo_msb_route_detail_round-robin"):responseRouteInfo.spec.lb_policy;
				vm.routeInfo.publish_protocol= responseRouteInfo.spec.publish_protocol;
				vm.routeInfo.publish_port= publishPort;
				vm.routeInfo.enable_ssl= responseRouteInfo.spec.enable_ssl;
				vm.routeInfo.enable_refer_match= responseRouteInfo.spec.enable_refer_match;
				vm.routeInfo.connect_timeout= responseRouteInfo.spec.connect_timeout==""?"--":responseRouteInfo.spec.connect_timeout+" s";
				vm.routeInfo.send_timeout= responseRouteInfo.spec.send_timeout==""?"--":responseRouteInfo.spec.send_timeout+" s";
				vm.routeInfo.read_timeout= responseRouteInfo.spec.read_timeout==""?"--":responseRouteInfo.spec.read_timeout+" s";
				vm.routeInfo.nodes= responseRouteInfo.spec.nodes;
				vm.routeInfo.publishUrl=routeUtil.generatePublishUrl(vm.routeInfo);


	                                    	
	       },
             error: function(XMLHttpRequest, textStatus, errorThrown) {
				  // bootbox.alert("get apiRouteInfo  fails："+XMLHttpRequest.responseText);  
				  routeUtil.notify('get RouteInfo fails:',XMLHttpRequest.statusText,'danger');                      
                 
             }
		});		
   				
   			    
	},
	viewTcpudpRoute:function(){

			 vm.route.routeIP=window.location.host.split(":")[0];

			 var url= window.location.search.substr(1);
		 	 var namespace=routeUtil.getQueryString(url,"namespace");
		 	
		 	 vm.routeInfo.serviceName=routeUtil.getQueryString(url,"serviceName");
			 vm.routeInfo.serviceVersion= routeUtil.getQueryString(url,"serviceVersion");

			 var version=vm.routeInfo.serviceVersion==""?"null":vm.routeInfo.serviceVersion;

   			$.ajax({
	                "type": 'get',
	                 "url":  vm.$serviceInstanceUrl+"/servicename/"+vm.routeInfo.serviceName+"/version/"+version+"?namespace="+namespace,
	                //"url":  "./data/tcpudpRouteInfo.json",
	                timeout: 8000,
	                "dataType": "json",
	                success: function (responseRouteList) {  
	                     
            	var responseRouteInfo=responseRouteList[0];
				vm.routeInfo.routeType= responseRouteInfo.spec.publish_protocol;
				vm.routeInfo.status= responseRouteInfo.status;
				vm.routeInfo.namespace= responseRouteInfo.metadata.namespace==""?"default":responseRouteInfo.metadata.namespace;
				vm.routeInfo.updateTimestamp= responseRouteInfo.metadata.updateTimestamp;
				vm.routeInfo.host= responseRouteInfo.spec.host;
				vm.routeInfo.url= responseRouteInfo.spec.url==""?"/":responseRouteInfo.spec.url;
				vm.routeInfo.lb_policy= responseRouteInfo.spec.lb_policy==""?$.i18n.prop("org_openo_msb_route_detail_round-robin"):responseRouteInfo.spec.lb_policy;
				vm.routeInfo.publish_protocol= responseRouteInfo.spec.publish_protocol;
				vm.routeInfo.publish_port= responseRouteInfo.spec.publish_port;
				vm.routeInfo.enable_tls= responseRouteInfo.spec.enable_tls;
				vm.routeInfo.connect_timeout= responseRouteInfo.spec.connect_timeout==""?"--":responseRouteInfo.spec.connect_timeout+" s";
				vm.routeInfo.proxy_responses= responseRouteInfo.spec.proxy_responses==""?"--":responseRouteInfo.spec.proxy_responses+" s";
				vm.routeInfo.proxy_timeout= responseRouteInfo.spec.proxy_timeout==""?"--":responseRouteInfo.spec.proxy_timeout+" s";
				vm.routeInfo.nodes= responseRouteInfo.spec.nodes;
				vm.routeInfo.publishUrl=routeUtil.generatePublishUrl(vm.routeInfo);
	       },
             error: function(XMLHttpRequest, textStatus, errorThrown) {			
				  routeUtil.notify('get RouteInfo fails:',XMLHttpRequest.statusText,'danger');                      
                 
             }
		});		
   				
   			    
	},
	
			showApiDoc:function(){

				if($('#apidocFrame').attr("src")!=null) return;
				

				var version=vm.routeInfo.routeVersion==""?"":"/"+vm.routeInfo.routeVersion;
				var sourceUrl= "/apijson/"+vm.routeInfo.routeName+version;	


				var url="../api-doc/index.html?publish_protocol="+vm.routeInfo.publish_protocol+"&publish_port="+vm.routeInfo.publish_port+"&url="+sourceUrl+"&api=/api/"+vm.routeInfo.routeName+version;

				 $('#apidocFrame').attr("src",url);  
	
			},
			showMetrics:function(){
					
				if($('#metricsFrame').attr("src")!=null) return;

				var version=vm.routeInfo.routeVersion==""?"":"/"+vm.routeInfo.routeVersion;	
				var sourceUrl= "/admin/"+vm.routeInfo.routeName+version;

				var url="../iui-metrics/index.html?publish_protocol="+vm.routeInfo.publish_protocol+"&publish_port="+vm.routeInfo.publish_port+"&url="+sourceUrl;

				 $('#metricsFrame').attr("src",url);   
			},

	});