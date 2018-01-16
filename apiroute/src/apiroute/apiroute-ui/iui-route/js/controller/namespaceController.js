var vm = avalon
		.define({
			$id : "namespaceController",
			$namespaceUrl:routeBasePath+'/namespaces',
			namespaceArray :  [],											
			initNamespaceList:function(){
				//vm.namespaceArray=["default","wudithwudithwudithwudithwudith","ns","xwj"];

				var namespaceArray=[];
				$.ajax({
	                "type": 'get',
	                "url":  vm.$namespaceUrl,
	                //"url":  "./data/namespace.json",
	                "timeout": 10000,
	                "dataType": "json",
	                success: function (resp) {  
	                      namespaceArray = (resp==null)?[]:resp;  
	                      namespaceArray.sort(function(a,b){
	                      	if(a=="default") return -1;
	                      	return a>b?1:-1
	                      }); 
          	
	                },
	                 error: function(XMLHttpRequest, textStatus, errorThrown) {			
						  routeUtil.notify('get namespaceList fails:',XMLHttpRequest.statusText,'danger'); 
						 
	                 },
	                  complete:function(){
	                  	  if(namespaceArray.length==0){
	                  	  	window.location.href="serviceList.html";
	                  	  }
	                      else if(namespaceArray.length==1){
	                  		window.location.href="serviceList.html?namespace="+namespaceArray[0];
	                     }	                 	
	                 	vm.namespaceArray=namespaceArray;            	
	                  }
				});
			 
			},
			viewServiceList:function(namespace){

				window.location.href="serviceList.html?namespace="+namespace;

			}
	});		