package test

import (
	credential "github.com/bytedance/douyin-openapi-credential-go/client"
	openApiSdkClient "github.com/bytedance/douyin-openapi-sdk-go/client"
	"testing"
)

func Test_GetCode2Session(t *testing.T) {

	//初始化
	core.BaseInit()
	defer core.CloseDB()

}

func TestAppsV2Token(t *testing.T) {
	// 初始化SDK client
	opt := new(credential.Config).
		SetClientKey("tt6dc61b7bf2cd0d0002").                       // 改成自己的app_id
		SetClientSecret("ca0092473225d541f494284f6b11a6ea45a7fd96") // 改成自己的secret
	sdkClient, err := openApiSdkClient.NewClient(opt)
	if err != nil {
		t.Log("sdk init err:", err)
		return
	}

	/* 构建请求参数，该代码示例中只给出部分参数，请用户根据需要自行构建参数值
	   	token:
	   	   1.若用户自行维护token,将用户维护的token赋值给该参数即可
	          2.SDK包中有获取token的函数，请根据接口path在《OpenAPI SDK 总览》文档中查找获取token函数的名字
	            在使用过程中，请注意token互刷问题
	       header:
	          sdk中默认填充content-type请求头，若不需要填充除content-type之外的请求头，删除该参数即可
	*/
	sdkRequest := &openApiSdkClient.AppsV2TokenRequest{}
	sdkRequest.SetAppid("tt6dc61b7bf2cd0d0002")
	sdkRequest.SetGrantType("client_credential")
	sdkRequest.SetSecret("ca0092473225d541f494284f6b11a6ea45a7fd96")
	// sdk调用
	sdkResponse, err := sdkClient.AppsV2Token(sdkRequest)
	if err != nil {
		t.Log("sdk call err:", err)
		return
	}
	//0801121847445035792f485169705149384261556b2b58646f413d3d
	t.Log(sdkResponse)
}

func TestAppsJscode2session(t *testing.T) {
	// 初始化SDK client

	opt := new(credential.Config).
		SetClientKey("tt6dc61b7bf2cd0d0002").                       // 改成自己的app_id
		SetClientSecret("ca0092473225d541f494284f6b11a6ea45a7fd96") // 改成自己的secret
	sdkClient, err := openApiSdkClient.NewClient(opt)
	if err != nil {
		t.Log("sdk init err:", err)
		return
	}

	/* 构建请求参数，该代码示例中只给出部分参数，请用户根据需要自行构建参数值
	   	token:
	   	   1.若用户自行维护token,将用户维护的token赋值给该参数即可
	          2.SDK包中有获取token的函数，请根据接口path在《OpenAPI SDK 总览》文档中查找获取token函数的名字
	            在使用过程中，请注意token互刷问题
	       header:
	          sdk中默认填充content-type请求头，若不需要填充除content-type之外的请求头，删除该参数即可
	*/
	sdkRequest := &openApiSdkClient.AppsJscode2sessionRequest{}
	sdkRequest.SetAppid("tt6dc61b7bf2cd0d0002")
	sdkRequest.SetSecret("ca0092473225d541f494284f6b11a6ea45a7fd96")
	sdkRequest.SetAnonymousCode("UrMISfLeTrSqZUrw_b0ib9P6tPjNdbJpt1tfYi9c4y_ZET6Gdd3C7vVAH-L_UGHlPB3UR4TqCdR8dRe2OKVOeD9e1feMOYZ8FL9Ivg")
	//sdkRequest.SetCode("tOPmWTgyFFx8YCAN4WQOxQmZDmCw6MHZuFCciFSp_8Qf-hsjRkrBGhvvUJkrOvYWo2ZXGVeSoYl125Zka9Xc7gTPTtBWbOoBhYLW2ZUjWZH-dMSnuGGqYwG5dPo")
	// sdk调用
	sdkResponse, err := sdkClient.AppsJscode2session(sdkRequest)
	if err != nil {
		t.Log("sdk call err:", err)
		return
	}
	t.Log(sdkResponse)
}
