package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

func WorkFlowRouter(r *gin.Engine) {

	{
		Appsv1.
			GET("/wf_list", controller.GetWFHandler).
			POST("/wf_create", controller.CreateWFHandler).
			DELETE("/wf_delete", controller.DeleteWFHandler)
	}
}
