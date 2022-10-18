package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

var v1 *gin.RouterGroup

func PodRouter(r *gin.Engine) {
	v1 = r.Group("/api/v1")
	{
		v1.
			GET("/pod_list", controller.GetPodHandler).
			GET("/pod_log", controller.GetPodLogHandler).
			GET("/pod_num", controller.GetNsPodNumHandler).
			GET("/pod_details", controller.GetPodDetailsHandler).
			GET("/pod_container_list", controller.GetContainerNameHandler).
			GET("/event_list", controller.GetEventListHandler).
			POST("/pod_create", controller.CreatePodHandler).
			DELETE("/pod_delete", controller.DeletePodHandler).
			PATCH("pod_update", controller.UpdatePodHandler)

	}
}
