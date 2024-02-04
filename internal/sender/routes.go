package sender

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (ns *NotificationsSender) RegisterRoutes(r *gin.Engine) {
	group := r.Group("/telegram", AccessControl())

	group.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
		return
	})
	group.POST("/notifications", ns.TriggerNotifications)
}
