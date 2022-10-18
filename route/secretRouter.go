package route

import (
	"github.com/gin-gonic/gin"
	"k8sManager/controller"
)

func SecretRouter(r *gin.Engine) {
	{
		v1.
			GET("/secret_list", controller.GetSecretHandler).
			GET("/secret_num", controller.GetNsSecretNumHandler).
			GET("/secret_details", controller.GetSecretDetailsHandler).
			GET("/secret_data", controller.GetSecretDataHandler).
			POST("/secret_create", controller.CreateSecretHandler).
			DELETE("/secret_delete", controller.DeleteSecretHandler).
			PATCH("secret_update", controller.UpdateSecretHandler)

	}
}
