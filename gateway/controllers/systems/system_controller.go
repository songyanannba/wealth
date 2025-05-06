// Package systems 系统查询
package systems

import (
	"fmt"
	"gate_way/common"
	"gate_way/controllers"
	"gate_way/servers/websocket"
	"runtime"

	"github.com/gin-gonic/gin"
)

// Status 查询系统状态
func Status(c *gin.Context) {
	isDebug := c.Query("isDebug")
	fmt.Println("http_request 查询系统状态", isDebug)
	data := make(map[string]interface{})
	numGoroutine := runtime.NumGoroutine()
	numCPU := runtime.NumCPU()

	// goroutine 的数量
	data["numGoroutine"] = numGoroutine
	data["numCPU"] = numCPU

	// ClientManager 信息
	data["managerInfo"] = websocket.GetManagerInfo(isDebug)
	controllers.Response(c, common.OK, "", data)
}
