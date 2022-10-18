package route

import "github.com/gin-gonic/gin"

func RouterLoad(r *gin.Engine) {
	{

		PodRouter(r)
		DeployRouter(r)
		DsRouter(r)
		PvRouter(r)
		StsRouter(r)
		NsRouter(r)
		NodeRouter(r)
		PvcRouter(r)
		CMRouter(r)
		SecretRouter(r)
		ServiceRouter(r)
		IngressRouter(r)
		WorkFlowRouter(r)
	}
}
