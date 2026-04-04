package tools

import (
	"easyChat/global"
	"encoding/json"
	"fmt"
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	dypnsapi20170525 "github.com/alibabacloud-go/dypnsapi-20170525/v3/client"
	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/credentials-go/credentials"
	credential "github.com/aliyun/credentials-go/credentials"
)

// 发送验证码

func CreateClient() (_result *dypnsapi20170525.Client, _err error) {
	credentialsConfig := new(credentials.Config).
		SetType("access_key").
		SetAccessKeyId(global.Config.AliAccess.AccessKey).
		SetAccessKeySecret(global.Config.AliAccess.SecretKey)
	credential, _err := credential.NewCredential(credentialsConfig)
	if _err != nil {
		return _result, _err
	}

	config := &openapi.Config{
		Credential: credential,
	}

	config.Endpoint = tea.String("dypnsapi.aliyuncs.com")
	_result = &dypnsapi20170525.Client{}
	_result, _err = dypnsapi20170525.NewClient(config)
	return _result, _err
}

func SendTelephoneCode(args []*string, code string) (_err error) {
	// 参数验证
	if len(args) < 4 {
		return fmt.Errorf("参数不足，需要4个参数: 签名、手机号、模板代码、验证码有效期(分钟)")
	}

	for i, arg := range args {
		if arg == nil {
			return fmt.Errorf("第%d个参数为空", i+1)
		}
	}

	client, _err := CreateClient()
	if _err != nil {
		return _err
	}

	sendSmsVerifyCodeRequest := &dypnsapi20170525.SendSmsVerifyCodeRequest{
		SignName:      tea.String(*args[0]),
		TemplateCode:  tea.String(*args[2]),
		PhoneNumber:   tea.String(*args[1]),
		TemplateParam: tea.String("{\"code\":\"" + code + "\",\"min\":\"" + *args[3] + "\"}"),
	}
	runtime := &util.RuntimeOptions{}
	tryErr := func() (_e error) {
		defer func() {
			if r := tea.Recover(recover()); r != nil {
				_e = r
			}
		}()
		resp, _err := client.SendSmsVerifyCodeWithOptions(sendSmsVerifyCodeRequest, runtime)
		if _err != nil {
			return _err
		}

		global.Logger.Infof("SendSmsVerifyCodeResponse: %v", resp)

		return nil
	}()

	if tryErr != nil {
		var error = &tea.SDKError{}
		if _t, ok := tryErr.(*tea.SDKError); ok {
			error = _t
		} else {
			error.Message = tea.String(tryErr.Error())
		}
		// 诊断地址
		var data interface{}
		d := json.NewDecoder(strings.NewReader(tea.StringValue(error.Data)))
		d.Decode(&data)
		if m, ok := data.(map[string]interface{}); ok {
			recommend, _ := m["Recommend"]
			global.Logger.Errorf("SendSmsVerifyCodeRecommend: %v,%v", recommend, error.Message)

		} else {
			global.Logger.Errorf("SendSmsVerifyCodeError: %v", error)
		}
	}
	return _err
}
