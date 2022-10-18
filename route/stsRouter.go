package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

func StsRouter(r *gin.Engine) {

	{
		Appsv1.
			GET("/sts_list", controller.GetStsHandler).
			GET("/sts_num", controller.GetNsStsNumHandler).
			GET("/sts_details", controller.GetStsDetailsHandler).
			POST("/sts_create", controller.CreateStsHandler).
			DELETE("/sts_delete", controller.DeleteStsHandler).
			PATCH("/sts_update", controller.UpdateStsHandler).
			PATCH("/sts_restart", controller.RestartSts)

	}
}
