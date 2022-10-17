package route

import "github.com/gin-gonic/gin"

var R *gin.Engine

func InitRouter() {
	R = gin.Default()
	RouterLoad(R)
}
