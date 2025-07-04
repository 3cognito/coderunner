package main

import (
	"3cognito/coderunner/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		utils.SuccessResponse(ctx, 200, "Up", "d")
	})
	r.Run()
}
