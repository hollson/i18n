// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/hollson/i18n"
	"golang.org/x/text/language"
)

type ApiResult struct {
	Code int
	Msg  *i18n.Message // 消息内容可以是Template模板
}

const (
	CODE_ERR_UNKNOWN = 1001
	CODE_ERR_SERVER  = 1003
	CODE_ERR_DB      = 1004
	CODE_ERR_AUTH    = 1005
	ERR_ROOM_CLOSED  = 1006
)

var (
	ERROR_SERVER      = i18n.Message{ID: "ERR_Server", Other: "server error", Desc: "服务错误"}
	ERROR_DB          = i18n.Message{ID: "ERR_DB", Other: "database error", Desc: "数据库错误"}
	ERROR_AUTH        = i18n.Message{ID: "ERR_Auth", Other: "auth error", Desc: "授权验证失败"}
	ERROR_ROOM_CLOSED = i18n.Message{ID: "ERR_ROOM_CLOSED", Other: "{{.RoomId}} room closed", Desc: "游戏房间已关闭"}
)

//go:generate  go run main2.go zh
//go:generate  go run main2.go en
func main() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 不需要加载active.en.toml，因为我们提供了默认的翻译。
	bundle.MustLoadMessageFile("active.zh.toml")

	lang := os.Args[1]                           // 接收一个lang参数
	localizer := i18n.NewLocalizer(bundle, lang) // 本地化转换器

	// 封装的API-Response
	response := func(res ApiResult) {
		if res.Msg == nil {
			fmt.Printf("code=%d,msg=%s\n", res.Code, `""`)
			return
		}

		fmt.Printf("code=%d,msg=%s\n", res.Code, localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: res.Msg,
			TemplateData: map[string]string{
				"RoomId": "1008",
			},
		}))
	}

	// 模拟Http-Handler
	response(ApiResult{Code: CODE_ERR_SERVER, Msg: &ERROR_SERVER})
	response(ApiResult{Code: CODE_ERR_DB, Msg: &ERROR_DB})
	response(ApiResult{Code: CODE_ERR_AUTH, Msg: &ERROR_AUTH})
	response(ApiResult{Code: ERR_ROOM_CLOSED, Msg: &ERROR_ROOM_CLOSED})
	response(ApiResult{Code: CODE_ERR_UNKNOWN})
}
