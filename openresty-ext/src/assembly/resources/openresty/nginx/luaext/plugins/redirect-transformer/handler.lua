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
local str_util  =  require('lib.utils.str_util')
local str_low = string.lower
local log = log_util.log
local str_startswith = str_util.startswith
local RedirectTransformerPluginHandler = BasePlugin:extend()

function RedirectTransformerPluginHandler:new()
	RedirectTransformerPluginHandler.super.new(self, "redirect-transformer-plugin")
end

function RedirectTransformerPluginHandler:header_filter()
	RedirectTransformerPluginHandler.super.header_filter(self)
	local originloc = ngx.header.Location
	local newloc
	if(originloc) then
		log("origin location:",originloc)
		local m, err = ngx.re.match(originloc, "^(\\w+)://([^/:]*)(?::(\\d+))?(/[^\\?]*).*", "o")
		local scheme,host,port,uri
		if m then
			scheme = m[1]
			host = m[2]
			port = m[3]
			uri = m[4]
		else
			return --It is not normal to enter this branch. This match result just let redirect transformer ignore this request(do nothing)
		end
		
		--If the port number is omitted, use the default port according to the scheme
		if port==false then
			local scheme = str_low(scheme)
			if scheme == "http" then
				port = "80"
			elseif scheme == "https" then
				port = "443"
			end
		end
		local ngx_var = ngx.var
		local req_host = ngx_var.host
		local req_port = ngx_var.server_port
		local req_scheme = ngx_var.scheme
		local last_peer = ngx.ctx.last_peer
		local backend_ip
		local backend_port
		if last_peer then
			backend_ip = last_peer.ip
			backend_port = last_peer.port
		end

		if not (host == req_host or host == backend_ip)  then return end
		if not (port == req_port or port == backend_port)  then return end
		--replace scheme,host,port
		newloc = ngx.re.sub(originloc, "^(https|http)://([^/]+)", req_scheme.."://"..req_host..":"..req_port, "oi")
		--check url
		local svc_pub_url = ngx.ctx.svc_pub_url
		local svc_url = ngx.ctx.svc_url
		--replace only if redirect location do not starts with pub_url and starts with svc_url
		if((not str_startswith(uri,svc_pub_url)) and (svc_url == "" or str_startswith(uri,svc_url))) then
			--replace $svc_url with $svc_pub_url
			newloc = ngx.re.sub(newloc, "^(https|http)://([^/]+)"..svc_url, "$1".."://".."$2"..svc_pub_url, "oi")
		end
		ngx.header["Location"] = newloc
		log("redirect-transformer output:","replace the redirect address to :"..newloc)
		ngx.log(ngx.WARN, "redirect-transformer replace the redirect address to:"..newloc, " origin location:",originloc)
	end
end

return RedirectTransformerPluginHandler