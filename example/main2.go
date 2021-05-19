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

var (
	error_server = i18n.Message{ID: "ERR_Server", Other: "server error", Description: "服务错误"}
	error_db     = i18n.Message{ID: "ERR_DB", Other: "Database error", Description: "数据库错误"}
	error_auth   = i18n.Message{ID: "ERR_Auth", Other: "auth error", Description: "授权验证失败"}
)

//go:generate  go run main2.go zh
func main() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 不需要加载active.en.toml，因为我们提供了默认的翻译。
	bundle.MustLoadMessageFile("active.zh.toml")

	lang := os.Args[1]                           // 接收一个lang参数
	localizer := i18n.NewLocalizer(bundle, lang) // 本地化转换器

	// response提供i18n.Message数据

	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &error_auth}))
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &error_server}))
	fmt.Println(localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &error_db}))
}
