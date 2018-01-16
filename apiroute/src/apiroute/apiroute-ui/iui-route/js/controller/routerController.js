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
var table;
var vm = avalon
		.define({
			$id : "routerController",
			namespace:"",
			route:{
				routeWay:"ip",
				routeHost:"",
				routeIP:"",
				routePort:"",
				routeSubDomain:"",
				iuiRootPath:"iui",
			    apiRootPath:"api"
			},				
			route_type:[],
			$publish_protocol:["Http","Https","TCP","UDP"],		
			$routeListUrl :routeBasePath+'/routelist',	
			$defaultPortsUrl :routeBasePath+'/conf/defaultports',				
			selectSearch:{
				name:"",
				index:1,
				context:""				
			},
			searchTypeName:[],	
			setSearchType:function(name,index){
				
				if(vm.selectSearch.context!=""){
				$('#msbTable').DataTable().column(vm.selectSearch.index).search(
        			"",true,true
    				).draw();
				vm.selectSearch.context="";
			  }

				vm.selectSearch.name=name;
				vm.selectSearch.index=index;
				
			},
			initMSBRoute:function(){
				vm.route.routeHost=window.location.host.split(":")[0]; //Default show port：80
				vm.route.routeIP=window.location.host.split(":")[0];
				vm.initIUIfori18n();

				var url= window.location.search.substr(1);
		 		var namespace=routeUtil.getQueryString(url,"namespace");
		 		
		 		vm.namespace=namespace==null?"default":namespace; //==""?"default":namespace;

	var t=$('#msbTable').DataTable( {
        "bProcessing": true, 
        ajax:{ 
         //"url":"./data/data.json", 
         "timeout": 30000,
         "url":vm.$routeListUrl+"?namespace="+vm.namespace,
          "dataSrc": "",
          "error":function(msg){
          	routeUtil.notify('get ServiceList fails:',"Rest URL "+msg.statusText,'danger');          
          	$("#msbTable_processing").hide();
          }
      }, 
        columns: [
            {"data": null},
            { "data": 'serviceName',
            "render": function ( data, type, full, meta ) {
            return routeUtil.showServiceName(full);
        	 }
            },
            { "data": 'serviceVersion' },
            { "data": 'routeType',
            "render": function ( data, type, full, meta ) {
	           return routeUtil.showProtocol(data);
          } 
         },
          
            { data: 'publishUrl',
             "render": function ( data, type, full, meta ) {
	           return routeUtil.showPublishUrl(full);
             }  
            },
            { data: 'publishPort' },
            { data: 'publishProtocol' },
            {
            "data": null,
            "render": function ( data, type, full, meta ) {
            return routeUtil.showButton(full);
          	}          
           }
        ],
        orderFixed: [[5, 'asc'],[3, 'asc']],
        rowGroup: {
            dataSrc: 'publishPort',
             startRender: function ( rows, group ) {

                return '<div>'+group 
                + '<span class="protocal">'+ rows.data()[0].publishProtocol+'</span>'
                //+' <span class="label label-info">'+rows.count()+'</span>'
                +'</div>';
            },
        },
        "oLanguage": vm.dataTableLanguage,
	    "dom": '<"top">rt<"bottom"lip><"clear">',
	    "sPaginationType": "bootstrap_extended",
	    "columnDefs": [ 		
			{
		      "targets": [0,2,7],
		      "searchable": false,
		    },
		    {
		      "targets": [0,1,2,3,4,5,6,7],
		      "bSortable": false,
		    }			
       	 ]
    } );
	t.columns( [ 5,6 ] ).visible( false);
    t.on( 'order.dt search.dt', function () {
        t.column(0, {search:'applied', order:'applied'}).nodes().each( function (cell, i) {
            cell.innerHTML = i+1;
        } );
    } ).draw();

	           
},
   searchService4keydown:function(event){
				if(event.keyCode == "13")      
  				{  
  					vm.searchService();
  				}
				
			},
	searchService:function(){
				var seachContext;
				if(vm.selectSearch.index==6){
					seachContext = '^' + vm.selectSearch.context +'$';
				}
				else{
					seachContext=vm.selectSearch.context;
				}

				 $('#msbTable').DataTable().column(vm.selectSearch.index).search(
        			seachContext,true,true
    				).draw();
			},		
			initIUIfori18n:function(){
				vm.selectSearch.name=$.i18n.prop("org_openo_msb_route_servicename");	
				vm.searchTypeName=[
			   {
			   	name:$.i18n.prop("org_openo_msb_route_servicename"),
			   	index:1
			   },
			   {
			   	name:$.i18n.prop("org_openo_msb_route_routetype"),
			   	index:3
			   },			   
			   {
			   	name:$.i18n.prop("org_openo_msb_route_publishPort"),
			   	index:5
			   },
			   {
			   	name:$.i18n.prop("org_openo_msb_route_publishProtocol"),
			   	index:6
			   },
			   {
			   	name:$.i18n.prop("org_openo_msb_route_publishUrl"),
			   	index:4
			   }

			  
			 ];

			 vm.route_type=["",$.i18n.prop("org_openo_msb_route_api"),$.i18n.prop("org_openo_msb_route_iui"),$.i18n.prop("org_openo_msb_route_custom")];


 				vm.dataTableLanguage={
                "sProcessing": "<img src='./img/loading-spinner-grey.gif'/><span>&nbsp;&nbsp;Loadding...</span>",   
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
                "sInfoFiltered": $.i18n.prop("org_openo_msb_route-table-infofilter"),
                "sLoadingRecords" :$.i18n.prop("org_openo_msb_route-table-loading") 
            };	

			}
			
		

	});


