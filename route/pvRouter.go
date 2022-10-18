package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

func PvRouter(r *gin.Engine) {

	{
		v1.
			GET("/pv_list", controller.GetPvHandler).
			GET("/pv_num", controller.GetPvNumHandler).
			GET("/pv_details", controller.GetPVDetailsHandler).
			GET("/pv_status", controller.GetPvBindHandler).
			POST("/pv_create", controller.CreatePvHandler).
			DELETE("/pv_delete", controller.DeletePvHandler).
			PATCH("pv_update", controller.UpdatePvHandler)

	}
}
