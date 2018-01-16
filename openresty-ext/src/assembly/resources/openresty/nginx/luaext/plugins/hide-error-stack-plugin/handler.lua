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

local HideErrorStackHandler = BasePlugin:extend()

function HideErrorStackHandler:new()
	HideErrorStackHandler.super.new(self, "hideerrorstackplugin")
end

function HideErrorStackHandler:header_filter()
	HideErrorStackHandler.super.header_filter(self)
	if(ngx.var.upstream_status and ngx.var.upstream_status == "500") then
	    ngx.log(ngx.WARN, "upstream response 500 internal server error: ")
        ngx.status = ngx.HTTP_INTERNAL_SERVER_ERROR
        ngx.exit(ngx.status)
	end

end

return HideErrorStackHandler