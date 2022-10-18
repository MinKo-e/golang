package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

func PvcRouter(r *gin.Engine) {
	{
		v1.
			GET("/pvc_list", controller.GetPvcHandler).
			GET("/pvc_num", controller.GetNsPvcNumHandler).
			GET("/pvc_details", controller.GetPvcDetailsHandler).
			POST("/pvc_create", controller.CreatePvcHandler).
			DELETE("/pvc_delete", controller.DeletePvcHandler).
			PATCH("pvc_update", controller.UpdatePvcHandler)

	}
}
