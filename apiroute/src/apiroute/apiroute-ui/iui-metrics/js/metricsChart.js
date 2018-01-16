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
var metricsChart = {};


metricsChart.memoryPieMetrics = function(data){
 var memoryPieChart = echarts.init(document.getElementById('memoryPieChartDiv'), 'macarons');   

 var labelTop = {
    normal : {
        label : {
            show : true,
            position : 'center',
            formatter : '{b}',
            textStyle: {
                baseline : 'bottom',
                color: '#000000'
            }
        },
        labelLine : {
            show : false 
        }
    }
};
var labelFromatter = {
    normal : {
        label : {
            formatter : function (params){
                return (100 - params.value).toFixed(1) + '%'
            },
            textStyle: {
                baseline : 'top'
            }
        }
    },
}
var labelBottom = {
    normal : {
        color: '#eee',
        label : {
            show : true,
            position : 'center'
        },
        labelLine : {
            show : false
        }
    },
    emphasis: {
      shadowBlur:0,
      color: '#eee',
    }
};
var radius = [40, 55];
option = {
    legend: {
        x : 'center',
        y:"bottom",
        data:[
            'Code-Cache','Eden-Space','Perm-Gen','Survivor-Space','Tenured-Gen'
        ],
        itemGap:20
    },
    title : {
        text: $.i18n.prop('org_openo_msb_metrics_jvm_memory_radius'),
        subtext: '',
        x: 'center',
        textStyle:{
            fontWeight:"normal",
            fontSize:16
        }
    },   
    toolbox: {
        show : true,       
        feature : {
                    
            saveAsImage : {
            show : true,
            title : $.i18n.prop('org_openo_msb_metrics_chart_save_picture'),
            type : 'png'

            }
        }
    },
    series : [
        {
            type : 'pie',
            center : ['10%', '55%'],
            radius : radius,
            x: '0%', // for funnel
            itemStyle : labelFromatter, 
            label: {
                normal: {                                   
                    textStyle: {
                        color:'rgb(64,192,255)'
                    }
                }
            },         
            data : [
                {name:'other', value:100-data.CodeCache, itemStyle : labelBottom},
                {name:'Code-Cache', value:data.CodeCache,itemStyle : labelTop}
            ]
        },
        {
            type : 'pie',
            center : ['30%', '55%'],
            radius : radius,
            x:'20%', // for funnel
            itemStyle : labelFromatter,
            label: {
                normal: {                                   
                    textStyle: {
                        color:'rgb(150,219,89)'
                    }
                }
            },      
            data : [
                {name:'other', value:100-data.EdenSpace, itemStyle : labelBottom},
                {name:'Eden-Space', value:data.EdenSpace,itemStyle : labelTop}
            ]
        },
        {
            type : 'pie',
            center : ['50%', '55%'],
            radius : radius,
            x:'40%', // for funnel
            itemStyle : labelFromatter,
            label: {
                normal: {                                   
                    textStyle: {
                        color:'rgb(249,212,80)'
                    }
                }
            },    
            data : [
                {name:'other', value:100-data.PermGen, itemStyle : labelBottom},
                {name:'Perm-Gen', value:data.PermGen,itemStyle : labelTop}
            ]
        },
         {
            type : 'pie',
            center : ['70%', '55%'],
            radius : radius,
            x:'60%', // for funnel
            itemStyle : labelFromatter,
             label: {
                normal: {                                   
                    textStyle: {
                        color:'rgb(255,132,128)'
                    }
                }
            },   
            data : [
                {name:'other', value:100-data.SurvivorSpace, itemStyle : labelBottom},
                {name:'Survivor-Space', value:data.SurvivorSpace,itemStyle : labelTop}
            ]
        },
         {
            type : 'pie',
            center : ['90%', '55%'],
            radius : radius,
            x:'80%', // for funnel
            itemStyle : labelFromatter,
             label: {
                normal: {                                   
                    textStyle: {
                        color:'#b6a2de'
                    }
                }
            },   
            data : [
                {name:'other', value:100-data.TenuredGen, itemStyle : labelBottom},
                {name:'Tenured-Gen', value:data.TenuredGen,itemStyle : labelTop}
            ]
        }
    ]
};
                    
       
        // load data for echarts objects
         memoryPieChart.setOption(option); 
         window.onresize = memoryPieChart.resize;


}


metricsChart.memoryBarMetrics = function(data){
 var memoryBarChart = echarts.init(document.getElementById('memoryBarChartDiv')); 
var option = {
    title : {
        text: $.i18n.prop('org_openo_msb_metrics_jvm_memory_bar'),
        x:'center',
        textStyle:{
            fontWeight:"normal",
            fontSize:16
        }
    },
    tooltip : {
        trigger: 'axis',
         axisPointer : {            
            type : 'shadow'        
        }
    },
    legend: {
         data:[
            $.i18n.prop('org_openo_msb_metrics_jvm_memory_bar_init'),$.i18n.prop('org_openo_msb_metrics_jvm_memory_bar_used'),$.i18n.prop('org_openo_msb_metrics_jvm_memory_bar_total')
        ],
        x:'left'
    },
    toolbox: {
        show : true,
         iconStyle:{
            normal:{
                borderColor:'#00abff'
            }
        },
        feature : {
            
            saveAsImage : {
            show : true,
            title : $.i18n.prop('org_openo_msb_metrics_chart_save_picture'),
            type : 'png'

            }
        }
    },
    yAxis : [
        {
            type : 'category',
            name: $.i18n.prop('org_openo_msb_metrics_chart_jvm_unit'),  
            axisLine:{
                lineStyle:{
                    color:'#CCCCCC'
                }
            },
            nameTextStyle:{
              color:'#8D8D8D'
            },
            axisLabel:{
                 textStyle:{
                    color:'#4D5761'
                }
           },
            data : [$.i18n.prop('org_openo_msb_metrics_jvm_memory_bar_heap'),$.i18n.prop('org_openo_msb_metrics_jvm_memory_bar_non-heap'),$.i18n.prop('org_openo_msb_metrics_jvm_memory_bar_total-heap')]
            
        }
    
    ],
    xAxis : [
        {
            type : 'value',                   
            axisLine:{
                lineStyle:{
                    color:'#CCCCCC'
                 }
            },
            axisLabel:{
             textStyle:{
                color:'#4D5761'
            }
            //formatter:'{value} M'
           }
        }
    ],
    series : [
        {
            name:$.i18n.prop('org_openo_msb_metrics_jvm_memory_bar_init'),
            type:'bar',
            stack:'barGroup',
            barWidth: '20',
            itemStyle: {normal: {color:'#B7ECFC', label:{show:false}}},
            data:data.init
        },
        
        {
            name:$.i18n.prop('org_openo_msb_metrics_jvm_memory_bar_used'),
            type:'bar',
            stack:'barGroup',
            itemStyle: {normal: {color:'#66CAFC', label:{show:false,formatter:function(p){return p.value > 0 ? (p.value +'\n'):'';}}}},
            data:data.used
        },

        {
            name:$.i18n.prop('org_openo_msb_metrics_jvm_memory_bar_total'),
            type:'bar',
            stack:'barGroup',            
            itemStyle: {normal: {color:'#00ABFF',  barBorderRadius:[0, 3, 3, 0],label:{show:true,position: 'insideRight',formatter:function(p){return p.value > 0 ? (p.value +'\n'):'';}}}},
            data:data.max
        },
        
        
    ]
};


                    

 memoryBarChart.setOption(option); 
 window.onresize = memoryBarChart.resize;

}




metricsChart.threadsMetrics = function(data){

 var threadsChart = echarts.init(document.getElementById('threadsChartDiv')); 


 var option = {
    /*title : {
        text: $.i18n.prop('org_openo_msb_metrics_jvm_thread_chart'),
        subtext: '',
        x:'center'
    },*/
    tooltip : {
        trigger: 'item',
        formatter: "{b}{a}: <br/> {c} ({d}%)"
    },
    color:["#a1df6a","#f7da64","#53c6ff","#f67d79"],//Runnable,Timed waiting,Waiting,Blocked
    legend: {
        //orient : 'vertical',
        //x : 'left',
        data:['Blocked','Waiting','Timed waiting','Runnable'],
        itemGap:20,
        top:10
    },
    toolbox: {
        show : true,
        iconStyle:{
            normal:{
                borderColor:'#00abff'
            }
        },
        feature : {
                     
            saveAsImage : {
            show : true,
            title : $.i18n.prop('org_openo_msb_metrics_chart_save_picture'),
            type : 'png'

            }
        }
    },
    calculable : true,
    series : [
        {
            name:$.i18n.prop('org_openo_msb_metrics_thread'),
            type:'pie',
             radius: ['40%', '60%'],
            center: ['50%', '60%'],
            data:data
        }
    ]
};


 threadsChart.setOption(option); 
 window.onresize = threadsChart.resize;

                    

}


metricsChart.restMetrics = function(restMetrics_data,restMetrics_name){



  var labelFromatter=function (value){

       if(value.length>12) return value.substring(0,11)+"\n"+value.substring(11);
       else return value; 
    }


var restChart = echarts.init(document.getElementById('restChartDiv'), 'macarons'); 
var option = {
    title : {
        text: '',
        subtext: ''
    },
    
    tooltip : {
        trigger: 'axis'
    },
    legend: {
        data:restMetrics_name,
        //orient:'vertical',
        x:'left',
        itemGap:5
    },
    toolbox: {
        show : true,
        feature : {           
          
            saveAsImage : {
            show : true,
            title : $.i18n.prop('org_openo_msb_metrics_chart_save_picture'),
            name:$.i18n.prop('org_openo_msb_metrics_rest_title'),
            type : 'png'
            }
        }
    },

    xAxis : [
        {
            type : 'category',
            data : [$.i18n.prop('org_openo_msb_metrics_rest_bar_count')]
           
        }
    ],
    yAxis : [
        {
            type : 'value'
           
        }
    ],
    series : restMetrics_data
};
                    
 restChart.setOption(option); 
  window.onresize = restChart.resize;
}


metricsChart.requestsMetrics = function(requestsMetrics_data,requestsMetrics_name){
var requestsChart = echarts.init(document.getElementById('requestsChartDiv'), 'macarons'); 

option = {
    title : {
        text: '',
        subtext: ''
    },    
    tooltip : {
        trigger: 'axis'
    },
    legend: {
        data:requestsMetrics_name,
        itemGap:20

    },
    toolbox: {
        show : true,
        feature : {          
           
            saveAsImage : {
            show : true,
            title : $.i18n.prop('org_openo_msb_metrics_chart_save_picture'),
            name:$.i18n.prop('org_openo_msb_metrics_requests_title'),
            type : 'png'

            }
        }
    },
    calculable : true,
    yAxis : [
        {
             type : 'value'
        }
    ],
    xAxis : [
         {
            type : 'category',
            data : [$.i18n.prop('org_openo_msb_metrics_rest_bar_count')],
            boundaryGap : [0, 0.01]         
           
        }
    ],
    series : requestsMetrics_data
};


 requestsChart.setOption(option); 
   window.onresize = requestsChart.resize;
}



  
