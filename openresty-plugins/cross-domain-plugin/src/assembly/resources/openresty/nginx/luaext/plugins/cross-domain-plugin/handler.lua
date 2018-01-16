--[[

    Copyright (C) 2016 ZTE, Inc. and others. All rights reserved. (ZTE)

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

local BasePlugin = require "plugins.base_plugin"
local msbConf   =  require('conf.msbinit')
local log_util  =  require('lib.utils.log_util')
local log = log_util.log

local CrossDomainPluginHandler = BasePlugin:extend()

function CrossDomainPluginHandler:new()
	CrossDomainPluginHandler.super.new(self, "cross-domain-plugin")
end

function CrossDomainPluginHandler:header_filter()
	CrossDomainPluginHandler.super.header_filter(self)
	local origin =  ngx.var.http_origin
	if(not origin or origin=="") then return end
	--local router_subdomain = msbConf.routerConf.subdomain
	--local m, err = ngx.re.match(origin, "https?://(.+)\\."..router_subdomain,"o")
	--if m then
		ngx.header["Access_Control_Allow_Origin"] = origin
		ngx.header["Access_Control_Allow_Credentials"] = "true"
		ngx.header["Access_Control_Allow_Methods"] = "GET,POST,DELETE,PUT"
		log("added Access_Control_Allow_Methods",true)
		local access_control_request_headers = ngx.var.http_access_control_request_headers
		log("access_control_request_headers",access_control_request_headers)
		if(access_control_request_headers and access_control_request_headers ~= "") then
		    ngx.header["Access_Control_Allow_Headers"] = access_control_request_headers
			log("added Access_Control_Allow_Headers",true)
		end
	--end
end

return CrossDomainPluginHandler