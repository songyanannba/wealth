package user

import (
	"gateway/common"
	"gateway/controllers"
	"gateway/global"
	"gateway/helper"
	"gateway/models"
	"gateway/models/table"
	"gateway/utils/cache"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

func Register(c *gin.Context) {
	request := &models.Login{}
	data := make(map[string]interface{})
	err := c.ShouldBindJSON(&request)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, "", data)
		return
	}

	uname := request.UserName
	//如果用户名称存在  需要换一个
	info, err := table.GetBUserByUsername(uname)
	if err != nil {
		controllers.Response(c, common.ServerError, "", data)
		return
	}

	if info.ID > 0 && info.UserName == request.UserName {
		//重名字
		controllers.Response(c, common.UserNameRepeat, "", data)
		return
	}

	city := "china" //todo
	password := helper.MD5V([]byte(request.Password))
	bUser := table.BUser{
		Uuid:     uuid.New().String(),
		UserName: request.UserName,
		Nickname: request.UserName,
		Password: password,
		Amount:   0,
		Status:   1,
		City:     city,
	}
	err = table.CreateBUser(&bUser)
	if err != nil {
		controllers.Response(c, common.ServerError, "", data)
		return
	}
	controllers.Response(c, common.OK, "", gin.H{})
}

func Login(c *gin.Context) {
	request := &models.Login{}
	data := make(map[string]interface{})
	err := c.ShouldBindJSON(&request)
	if err != nil {
		controllers.Response(c, common.ParameterIllegal, "", data)
		return
	}

	global.GVA_LOG.Infof("Login:{%v}", *request)
	//获取用户名和密码
	//是否注册过，注册过返回用户信息

	userInfo, err := table.GetBUserByUsername(request.UserName)
	if err != nil {
		controllers.Response(c, common.ServerError, "", data)
		return
	}
	if userInfo.ID <= 0 {
		controllers.Response(c, common.NotRegister, "", data)
		return
	}

	password := helper.MD5V([]byte(request.Password))
	info, err := table.GetBUserByUsernameAndPassword(request.UserName, password)
	if err != nil {
		controllers.Response(c, common.ServerError, "", data)
		return
	}
	//没有注册过返回错误提示 让用户进行注册
	if info.ID <= 0 {
		controllers.Response(c, common.PasswordErr, "", data)
		return
	}

	//24小时
	token, _ := helper.GenerateJWT(info.Uuid, "", 24)
	//存储用户的登陆信息到redis
	err = cache.SetGateWayUserWeb(info.Uuid, &models.UserWeb{
		UserID:    info.Uuid,
		Nickname:  request.UserName,
		LoginTime: time.Now().Unix(),
		Token:     token,
	})
	if err != nil {
		global.GVA_LOG.Error("SetGateWayUserWeb err :", zap.Error(err))
		controllers.Response(c, common.ModelStoreError, "", data)
		return
	}

	data = gin.H{
		"token": token,
	}
	global.GVA_LOG.Infof("Login:%v", data)
	controllers.Response(c, common.OK, "", data)
}

func GUserInfo(c *gin.Context) {
	data := make(map[string]interface{})
	userUUid, err := helper.GetUserID(c)
	if err != nil {
		controllers.Response(c, common.UnauthorizedUserID, "", data)
		return
	}

	userInfo, err := table.GetBUserByUUid(userUUid)
	if err != nil {
		controllers.Response(c, common.ServerError, "", data)
		return
	}

	userInfoResp := models.UserInfoResp{
		UserName: userInfo.UserName,
		Amount:   userInfo.Amount,
		City:     userInfo.City,
	}

	data = gin.H{
		"user_info": userInfoResp,
	}
	global.GVA_LOG.Infof("Login:%v", data)
	controllers.Response(c, common.OK, "", data)
}
