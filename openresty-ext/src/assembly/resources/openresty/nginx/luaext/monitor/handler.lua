--[[

    Copyright (C) 2017 ZTE, Inc. and others. All rights reserved. (ZTE)

    Licensed under the Apache License, Version 2.0 (the "License");
    you may not use this file except in compliance with the License.
    You may obtain a copy of the License at

            http://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing, software
    distributed under the License is distributed on an "AS IS" BASIS,
    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
    See the License for the specific language governing permissions and
    limitations under the License.

--]]

local _M = {}
_M._VERSION = '1.0.0'

local cjson_safe = require "cjson.safe"

local tbl_concat = table.concat
local version_num = ngx.config.ngx_lua_version
local ngx_lua_version = tbl_concat({math.floor(version_num/1000000)%1000,math.floor(version_num/1000)%1000,version_num%1000},".")
local nginx_version = ngx.var.nginx_version
local worker_count = ngx.worker.count()

local function _collect_conn_status()
	local conn_status = {}
	conn_status["active"] = ngx.var.connections_active
	conn_status["reading"] = ngx.var.connections_reading
	conn_status["writing"] = ngx.var.connections_writing
	conn_status["waiting"] = ngx.var.connections_waiting
	return conn_status
end


local function _collect_real_time(index)
	local metrics = {}
	if(index == "connections") then
		metrics["connections"] = _collect_conn_status()
	elseif(index == "requests") then
		metrics["requests"] = stats.collect_req_status()
	else
		metrics["nginx_version"] = nginx_version
		metrics["ngx_lua_version"] = ngx_lua_version
		metrics["worker_count"] = worker_count
		metrics["connections"] = _collect_conn_status()
		metrics["requests"] = stats.collect_req_status()
	end

	local value, err = cjson_safe.encode(metrics)
	if err then
		return nil,"Collect real-time connection status failed! Error:" .. err
	end
	return value,""
end

local function _collect_stats_apigateway(index)
	--if input latestNum is empty or illegal, then set to 1
	local latest_num = tonumber(ngx.var.arg_latestNum) or 1
	if latest_num<=0 then
		latest_num =1
	end
	return stats.get_reqnum_stats(latest_num)
end

---------------------------------------------------------------
--collect the monitor result
--      query_type: real-time or stats
--      metrics_obj: apigateway or services
--      index: latency requests responses
---------------------------------------------------------------
local function _collect(query_type,metrics_obj,index,service_name,service_version)
	if(query_type == "real-time") then
		if(metrics_obj == "apigateway") then
			return _collect_real_time(index)
		else
			return nil,"Now the real-time metrics of the service is not supported!"
		end
	elseif(query_type == "stats") then
		if(metrics_obj == "apigateway") then
			return _collect_stats_apigateway(index)
		elseif(metrics_obj == "services") then
			return nil,"the stat metrics of services is not supported in this version!"
		else
			return nil,"not found!"
		end
	else
		return nil,"not found!"
	end
end

function _M.do_get()
	local uri = ngx.var.uri
	--local m, err = ngx.re.match(uri, "^/admin/microservices/v1/([^/]+)/([^/]+)(?:/([^/]+))?$", "o")
	local m, err = ngx.re.match(uri, "^/admin/microservices/v1/([^/]+)/(apigateway|services)(?:/([^/]+)/version/([^/]+))?(?:/([^/]+))?$", "o")
	if m then
		local result,err = _collect(m[1],m[2],m[5],m[3],m[4])
		if result then
			ngx.header.content_type = "application/json;charset=utf-8"
			ngx.print(result)
		else
			ngx.status = ngx.HTTP_NOT_FOUND
			ngx.print(err)
			return ngx.exit(ngx.status)
		end
	else
		ngx.status = ngx.HTTP_NOT_FOUND
		return ngx.exit(ngx.status)
	end
end

return _M