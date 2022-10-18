package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

func CMRouter(r *gin.Engine) {
	{
		v1.
			GET("/config_list", controller.GetConfigmapHandler).
			GET("/config_num", controller.GetNsConfigmapNumHandler).
			GET("/config_details", controller.GetConfigmapDetailsHandler).
			GET("/config_data", controller.GetConfigmapDataHandler).
			POST("/config_create", controller.CreateConfigmapHandler).
			DELETE("/config_delete", controller.DeleteConfigmapHandler).
			PATCH("config_update", controller.UpdateConfigmapHandler)

	}
}
