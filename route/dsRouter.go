package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

func DsRouter(r *gin.Engine) {
	appsV1 := r.Group("/api/appsv1")
	{
		appsV1.
			GET("/ds_list", controller.GetDsHandler).
			GET("/ds_num", controller.GetNsDsNumHandler).
			GET("/ds_details", controller.GetDsDetailsHandler).
			POST("/ds_create", controller.CreateDsHandler).
			DELETE("/ds_delete", controller.DeleteDsHandler).
			PATCH("/ds_update", controller.UpdateDsHandler).
			PATCH("/ds_restart", controller.RestartDs)

	}
}
