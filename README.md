# i18n

**i18n**是一个Go语言版本的多语言(国际化)包，并提供了`go18in`命令行工具。

- 支持Unicode[通用语言环境数据存储库(CLDR)](https://www.unicode.org/cldr/charts/28/supplemental/language_plural_rules.html)；
- 支持使用具有命名变量的字符串 `Template`语法 ；
-   支持多种消息文件格式, 如:  TOML、JSON、YAML；


<br/>

## 安装goi18n工具

```sh
$ go get -u github.com/hollson/i18n
$ goi18n -help
```
```text
$ goi18n
这是一个「github.com/hollson/i18」多语言库的客户端工具.
Usage:
    i18n_cli <Command> [Option] <Param>...

Command:
    extract     从go源码提取「i18n.Message」,即预翻译的消息(不包含测试文件)
    merge       合并翻译文件
```

<br/>

##  i18n包使用示例

创建示例代码：
```go
var tpl = template.Must(template.New("").Parse(`
<!DOCTYPE html>
<html>
<body>
	<h1>{{.Title}}</h1>
	{{range .Paragraphs}}<p>{{.}}</p>{{end}}
</body>
</html>
`))

// 测试：
//go:generate curl "http://localhost:8080/?name=Nick&unread=2"
//go:generate curl -X GET "http://localhost:8080/?name=Nick&unread=20&lang=en"
//go:generate curl -X GET "http://localhost:8080/?name=Nick&unread=30&lang=zh"
func main() {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)

	// 不需要加载active.en.toml，因为我们提供了默认的翻译。
	bundle.MustLoadMessageFile("active.zh.toml")
	// bundle.MustLoadMessageFile("active.es.toml")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		lang := r.FormValue("lang")
		// Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2
		accept := r.Header.Get("Accept-Language")
		localizer := i18n.NewLocalizer(bundle, lang, accept)
		fmt.Printf("语言列表：lang=%s;\taccept=%s; \n", lang, accept)

		name := r.FormValue("name")
		if name == "" {
			name = "Unknown"
		}
		unread, _ := strconv.ParseInt(r.FormValue("unread"), 10, 64)

		helloPerson := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Greeting",
				Other: "Hello {{.Name}}",
			},
			TemplateData: map[string]string{
				"Name": name,
			},
		})

		unreadEmails := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "UnreadEmails",
				Desc:  "The number of unread emails I have",
				One:   "I have {{.PluralCount}} unread email.",
				Other: "I have {{.PluralCount}} unread emails.",
			},
			PluralCount: unread,
		})

		unreadSms := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:          "UnreadSMS",
				Desc: "The number of unread sms",
				One:         "{{.Name}} has one {{.UnreadSms}} unread sms.",
				Other:       "{{.Name}} has {{.UnreadSms}} unread sms.",
			},
			PluralCount: 5,
			TemplateData: map[string]interface{}{
				"Name":             name,
				"UnreadSms": unread,
			},
		})

		err := tpl.Execute(w, map[string]interface{}{
			"Title": helloPerson,
			"Paragraphs": []string{
				unreadEmails,
				unreadSms,
			},
		})

		if err != nil {
			panic(err)
		}
	})

	fmt.Println("Listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
```

<br/>

### 提取词汇

执行 `goi18n extract`命令，可提取Go源文件中的所有`i18n.Message`对象(即提取要翻译的词汇)，这时候会生成一个默认的词汇文件`./active.en.toml`。 

```toml
# active.en.toml
[UnreadEmails]
one = "I have {{.PluralCount}} unread email."
other = "I have {{.PluralCount}} unread emails."

[UnreadSMS]
one = "{{.Name}} has one {{.UnreadSms}} unread sms."
other = "{{.Name}} has {{.UnreadSms}} unread sms."
```

<br/>

### 合并消息

假如我们需要提供一个中文的词汇库，此时需要创建一个`zh`空文件： `touch translate.zh.toml`

执行`goi18n merge active.en.toml translate.zh.toml` ，将要翻译的词汇拷贝到 `translate.zh.toml`文件中。

注意，如果`translate.zh.toml`已存在，则merge命令会将新增的词汇合并到`touch translate.zh.toml`中。

<br/>

### 翻译并渲染

人工翻译汉化包文件

```toml
# active.zh.toml
[UnreadEmails]
hash = "sha1-55687b25cf8ac24dbc9a2e091d4d7f14bc85d90d"
other = "您有 {{.PluralCount}} 份未读邮件."

[UnreadSMS]
hash = "sha1-f5aca1f705a50a9d4bd59e20fc07c8bc1218615f"
other = "{{.Name}} 有 {{.UnreadSms}} 条未读短信 ."
```

运行示例程序，并测试

```sh
$ curl -X GET "http://localhost:8080/?name=Nick&unread=20"
```

**默认方式输出：**

```html
<!DOCTYPE html>
<html>
<body>
    <h1>Hello Nick</h1>
    <p>I have 20 unread emails.</p><p>Nick has 20 unread sms.</p>
</body>
</html>
```
```sh
$  curl -X GET "http://localhost:8080/?name=Nick&unread=30&lang=zh"
```

**汉化方式输出：** 

```html
<!DOCTYPE html>
<html>
<body>
    <h1>你好，Nick</h1>
    <p>您有 30 份未读邮件.</p><p>Nick 有 30 条未读短信 .</p>
</body>
</html>
```



<br/>

>   参考链接： https://github.com/nicksnyder/go-i18n


