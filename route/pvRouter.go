package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

func PvRouter(r *gin.Engine) {

	{
		v1.
			GET("/pv_list", controller.GetpvHandler).
			GET("/pv_log", controller.GetpvLogHandler).
			GET("/pv_num", controller.GetNspvNumHandler).
			GET("/pv_details", controller.GetpvDetailsHandler).
			POST("/pv_create", controller.CreatepvHandler).
			DELETE("/pv_delete", controller.DeletepvHandler).
			PATCH("pv_update", controller.UpdatepvHandler)

	}
}
