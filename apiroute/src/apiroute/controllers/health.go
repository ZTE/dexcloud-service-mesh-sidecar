package controllers

import (
	"apiroute/redis"
	"github.com/astaxie/beego"
)

type HealthController struct {
	beego.Controller
}

// @Title GetHealthStatus
// @Description Get health status
// @Success 200 {string} Healthy
// @Failure 500 Unhealthy
// @router / [get]
func (c *HealthController) GetHealthStatus() {
	c.EnableRender = false
	if ok, err := redis.IsHealthy(); !ok {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.Output.Body([]byte(err.Error()))
		return
	}
}
