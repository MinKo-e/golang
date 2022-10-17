package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

var Appsv1 *gin.RouterGroup

func DeployRouter(r *gin.Engine) {
	Appsv1 = r.Group("/api/appsv1")
	{
		Appsv1.
			GET("/deploy_list", controller.GetDeployHandler).
			GET("/deploy_num", controller.GetNsDeployNumHandler).
			GET("/deploy_details", controller.GetDeployDetailsHandler).
			POST("/deploy_create", controller.CreateDeployHandler).
			DELETE("/deploy_delete", controller.DeleteDeployHandler).
			PATCH("/deploy_update", controller.UpdateDeployHandler).
			PATCH("/deploy_restart", controller.RestartDeploy)

	}
}
