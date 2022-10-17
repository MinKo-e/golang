package route

import "github.com/gin-gonic/gin"

func RouterLoad(r *gin.Engine) {
	PodRouter(r)
	DeployRouter(r)
	DsRouter(r)
}
