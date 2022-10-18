package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

func NsRouter(r *gin.Engine) {

	{
		v1.
			GET("/ns_list", controller.GetNsHandler).
			GET("/ns_num", controller.GetNsNumHandler).
			GET("/ns_details", controller.GetNsDetailsHandler).
			POST("/ns_create", controller.CreateNsHandler).
			DELETE("/ns_delete", controller.DeleteNsHandler).
			PATCH("ns_update", controller.UpdateNsHandler)

	}
}
