package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

var Networkv1 *gin.RouterGroup

func IngressRouter(r *gin.Engine) {
	Networkv1 = r.Group("/api/networkv1")
	{
		Networkv1.
			GET("/ingress_list", controller.GetIngressHandler).
			GET("/ingress_num", controller.GetNsIngressNumHandler).
			GET("/ingress_details", controller.GetIngressDetailsHandler).
			POST("/ingress_create", controller.CreateIngressHandler).
			DELETE("/ingress_delete", controller.DeleteIngressHandler).
			PATCH("ingress_update", controller.UpdateIngressHandler)

	}
}
