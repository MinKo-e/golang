package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

func ServiceRouter(r *gin.Engine) {

	{
		v1.
			GET("/svc_list", controller.GetSvcHandler).
			GET("/svc_num", controller.GetNsSvcNumHandler).
			GET("/svc_details", controller.GetSvcDetailsHandler).
			POST("/svc_create", controller.CreateSvcHandler).
			DELETE("/svc_delete", controller.DeleteSvcHandler).
			PATCH("svc_update", controller.UpdateSvcHandler)

	}
}
