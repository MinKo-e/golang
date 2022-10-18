package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

func NodeRouter(r *gin.Engine) {

	{
		v1.
			GET("/node_list", controller.GetNodeHandler).
			GET("/node_role", controller.GetNodeRoleHandler).
			GET("/node_details", controller.GetPVDetailsHandler).
			PATCH("node_update", controller.UpdateNodeHandler)

	}
}
