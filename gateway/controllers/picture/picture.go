package picture

import (
	"github.com/gin-gonic/gin"
)

func GetPicture(c *gin.Context) {
	var (
		data = make(map[string]interface{})
	)

	data = gin.H{}
	controllers.Response(c, common.WebOK, "", data)
}
