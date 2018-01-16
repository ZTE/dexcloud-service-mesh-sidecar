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
var routeUtil = {};

routeUtil.growl=function(title,message,type){
      $.growl({
		icon: "fa fa-envelope-o fa-lg",
		title: "&nbsp;&nbsp;"+$.i18n.prop('org_openo_msb_route_property_ttl')+title,
		message: message+"&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;"
			},{
				type: type
			});
}

routeUtil.cutString=function(str){
  var newStr;
  if(str.length>22){
     newStr=str.substring(0,20)+"...";
  }
  else{
    newStr=str;
  }

  return newStr;
}





routeUtil.showPublishUrl=function(router){
   var publishPort=router.publishPort==""?"":":"+router.publishPort;
 

  return router.publishProtocol.toLowerCase()+"://"+vm.route.routeIP+publishPort+router.publishUrl;
}



routeUtil.showButton=function(router){
  var buttonHtml;
  
  if(routeUtil.isHttpTypeProtocol(router.publishProtocol)){
      buttonHtml="<a  class='btn btn-default btn-s' target='_blank' href='routeDetail.html?publishport="+router.publishPort+"&routetype="+router.routeType+"&routename="+router.routeName+"&routeversion="+router.routeVersion+"' ><i class='fa fa-file-text-o' ></i> "+$.i18n.prop("org_openo_msb_route_box_btn_view")+"</a>";
  }
  else{
      buttonHtml="<a  class='btn btn-default btn-s' target='_blank' href='tcpudpDetail.html?namespace="+router.namespace+"&serviceName="+router.serviceName+"&serviceVersion="+router.serviceVersion+"' ><i class='fa fa-file-text-o' ></i> "+$.i18n.prop("org_openo_msb_route_box_btn_view")+"</a>";
  }

  return buttonHtml;
}

routeUtil.showServiceName=function(router){
  var serviceNameHtml;
  if(routeUtil.isHttpTypeProtocol(router.publishProtocol)){
    serviceNameHtml="<a  target='_blank' href='routeDetail.html?publishport="+router.publishPort+"&routetype="+router.routeType+"&routename="+router.routeName+"&routeversion="+router.routeVersion+"' >"+router.serviceName+"</a>";
  }
  else{
       serviceNameHtml="<a  target='_blank' href='tcpudpDetail.html?namespace="+router.namespace+"&serviceName="+router.serviceName+"&serviceVersion="+router.serviceVersion+"' >"+router.serviceName+"</a>";

  }

  return serviceNameHtml;
}

routeUtil.isHttpTypeProtocol=function(publishProtocol){
  var protocol=publishProtocol.toLowerCase();
  if(protocol=="http" || protocol=="https"){
    return true;
  }

  return false;
}






          


routeUtil.showProtocol=function(protocol){
            if(protocol=="api"){  
                 return  '<span class="label-protocol label-api">'+$.i18n.prop("org_openo_msb_route_api")+'</span>';
              }
              else if(protocol=="iui"){  
                 return  '<span class="label-protocol label-iui">'+$.i18n.prop("org_openo_msb_route_iui")+'</span>';
              }
              else if(protocol=="custom"){  
                 return  '<span class="label-protocol label-tcp">'+$.i18n.prop("org_openo_msb_route_custom")+'</span>';
              }
              else{
                return  '';
              }
}

routeUtil.getQueryString=function(url,name){
 var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)", "i"); 
    var r = url.match(reg); 
    if (r != null) return unescape(r[2]); 
    return null; 
      
}

routeUtil.notify=function(title,message,type){
  var delay;
  if(type=="danger"){
    delay=10000;
  }
  else{
    delay=5000;
  }
  $.notify({
            title: title,
            message: message
          },{
            type: type,
            placement: {
              from: "top",
              align: "center"
            },
            delay: delay,
            offset: {
              x: 50,
              y: 50
            }
          }); 
}

routeUtil.generatePublishUrl=function(routeInfo){
  var publishUrl;
  var publishProtocol=routeInfo.publish_protocol.toLowerCase();
  var routeType=routeInfo.routeType.toLowerCase();

  if(publishProtocol=="tcp" || publishProtocol=="udp"){
    publishUrl="";
  }
  else{
    if(routeType=="api"){
       var version=routeInfo.routeVersion==""?"":"/"+routeInfo.routeVersion;
       publishUrl="/api/"+routeInfo.routeName+version;
    }
    else if(routeType=="iui"){
       publishUrl="/iui/"+routeInfo.routeName;
    }
     else if(routeType=="custom"){
       publishUrl=routeInfo.routeName;
       //  var reg_customName_match=/^(\/.*?)$/;
       //  if(!reg_customName_match.test(serviceName)) serviceName="/"+serviceName;
       // vm.targetFullServiceUrl=vm.targetServiceUrl+serviceName;
    }
  }

  var publishPort=routeInfo.publish_port==""?"":":"+routeInfo.publish_port;
 

  return publishProtocol+"://"+vm.route.routeIP+publishPort+publishUrl;

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

routeUtil.toLocaleTime=function(timestamp){
  return new Date(timestamp).Format("yyyy-MM-dd hh:mm:ss");
}




routeUtil.generate_apiUrl=function(type){
  var publishProtocol=vm.routeInfo.publish_protocol.toLowerCase();
  var publishPort=vm.routeInfo.publish_port==""?"":":"+vm.routeInfo.publish_port;
 var url= publishProtocol+"://"+vm.route.routeIP+publishPort+"/"+type+"/"+vm.routeInfo.routeName+"/"+vm.routeInfo.routeVersion;
 window.open(url);

}

routeUtil.showEnableSet=function(value){
  if(value==true){
     return '<span class="label-protocol label-iui">'+$.i18n.prop("org_openo_msb_route_detail_enable")+'</span>';
  }
  else if(value==false){
    return '<span class="label-protocol label-tcp">'+$.i18n.prop("org_openo_msb_route_detail_disable")+'</span>';
  }
  else{
    return "";
  }
}

routeUtil.showEnableSet4String=function(value){
  if(value=="true"){
     return '<span class="label-protocol label-iui">'+$.i18n.prop("org_openo_msb_route_detail_enable")+'</span>';
  }
  else if(value=="false"){
    return '<span class="label-protocol label-tcp">'+$.i18n.prop("org_openo_msb_route_detail_disable")+'</span>';
  }
  else{
    return "--";
  }
}

routeUtil.closePage=function(){
          
      window.close(); 
 }

 routeUtil.showProtocolName=function(protocol){
            if(protocol=="api"){  
                 return  $.i18n.prop("org_openo_msb_route_api");
              }
              else if(protocol=="iui"){  
                 return  $.i18n.prop("org_openo_msb_route_iui");
              }
              else if(protocol=="custom"){  
                 return $.i18n.prop("org_openo_msb_route_custom");
              }
              else{
                return  '';
              }
}









