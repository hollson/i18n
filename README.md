# i18n

**i18n**是一个Go语言版本的[多语言](./lang.md)(国际化)包，并提供了`go18in`命令行工具。

- 支持Unicode[通用语言环境数据存储库(CLDR)](https://www.unicode.org/cldr/charts/28/supplemental/language_plural_rules.html)；
- 支持使用具有命名变量的字符串 `Template`语法 ；
- 支持多种消息文件格式, 如:  TOML、JSON、YAML；


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

# 附： 语言代码表

|代码|名称|
|--|--|
|af | 南非语|
| af-ZA| 南非语 |
| ar | 阿拉伯语 |
| ar-AE| 阿拉伯语(阿联酋) |
| ar-BH| 阿拉伯语(巴林) |
| ar-DZ| 阿拉伯语(阿尔及利亚) |
| ar-EG| 阿拉伯语(埃及) |
| ar-IQ| 阿拉伯语(伊拉克) |
| ar-JO| 阿拉伯语(约旦) |
| ar-KW| 阿拉伯语(科威特) |
| ar-LB| 阿拉伯语(黎巴嫩) |
| ar-LY| 阿拉伯语(利比亚) |
| ar-MA| 阿拉伯语(摩洛哥) |
| ar-OM| 阿拉伯语(阿曼) |
| ar-QA| 阿拉伯语(卡塔尔) |
| ar-SA| 阿拉伯语(沙特阿拉伯) |
| ar-SY| 阿拉伯语(叙利亚) |
| ar-TN| 阿拉伯语(突尼斯) |
| ar-YE| 阿拉伯语(也门) |
| az | 阿塞拜疆语 |
| az-AZ| 阿塞拜疆语(拉丁文) |
| az-AZ| 阿塞拜疆语(西里尔文) |
| be | 比利时语 |
| be-BY| 比利时语 |
| bg | 保加利亚语 |
| bg-BG| 保加利亚语 |
| bs-BA| 波斯尼亚语(拉丁文，波斯尼亚和黑塞哥维那) |
| ca | 加泰隆语 |
| ca-ES| 加泰隆语 |
| cs | 捷克语 |
| cs-CZ| 捷克语 |
| cy | 威尔士语 |
| cy-GB| 威尔士语 |
| da | 丹麦语 |
| da-DK| 丹麦语 |
| de | 德语 |
| de-AT| 德语(奥地利) |
| de-CH| 德语(瑞士) |
| de-DE| 德语(德国) |
| de-LI| 德语(列支敦士登) |
| de-LU| 德语(卢森堡) |
| dv | 第维埃语 |
| dv-MV| 第维埃语 |
| el | 希腊语 |
| el-GR| 希腊语 |
| en | 英语 |
| en-AU| 英语(澳大利亚) |
| en-BZ| 英语(伯利兹) |
| en-CA| 英语(加拿大) |
| en-CB| 英语(加勒比海) |
| en-GB| 英语(英国) |
| en-IE| 英语(爱尔兰) |
| en-JM| 英语(牙买加) |
| en-NZ| 英语(新西兰) |
| en-PH| 英语(菲律宾) |
| en-TT| 英语(特立尼达) |
| en-US| 英语(美国) |
| en-ZA| 英语(南非) |
| en-ZW| 英语(津巴布韦) |
| eo | 世界语 |
| es | 西班牙语 |
| es-AR| 西班牙语(阿根廷) |
| es-BO| 西班牙语(玻利维亚) |
| es-CL| 西班牙语(智利) |
| es-CO| 西班牙语(哥伦比亚) |
| es-CR| 西班牙语(哥斯达黎加) |
| es-DO| 西班牙语(多米尼加共和国) |
| es-EC| 西班牙语(厄瓜多尔) |
| es-ES| 西班牙语(传统) |
| es-ES| 西班牙语(国际) |
| es-GT| 西班牙语(危地马拉) |
| es-HN| 西班牙语(洪都拉斯) |
| es-MX| 西班牙语(墨西哥) |
| es-NI| 西班牙语(尼加拉瓜) |
| es-PA| 西班牙语(巴拿马) |
| es-PE| 西班牙语(秘鲁) |
| es-PR| 西班牙语(波多黎各(美)) |
| es-PY| 西班牙语(巴拉圭) |
| es-SV| 西班牙语(萨尔瓦多) |
| es-UY| 西班牙语(乌拉圭) |
| es-VE| 西班牙语(委内瑞拉) |
| et | 爱沙尼亚语 |
| et-EE| 爱沙尼亚语 |
| eu | 巴士克语 |
| eu-ES| 巴士克语 |
| fa | 法斯语 |
| fa-IR| 法斯语 |
| fi | 芬兰语 |
| fi-FI| 芬兰语 |
| fo | 法罗语 |
| fo-FO| 法罗语 |
| fr | 法语 |
| fr-BE| 法语(比利时) |
| fr-CA| 法语(加拿大) |
| fr-CH| 法语(瑞士) |
| fr-FR| 法语(法国) |
| fr-LU| 法语(卢森堡) |
| fr-MC| 法语(摩纳哥) |
| gl | 加里西亚语 |
| gl-ES| 加里西亚语 |
| gu | 古吉拉特语 |
| gu-IN| 古吉拉特语 |
| he | 希伯来语 |
| he-IL| 希伯来语 |
| hi | 印地语 |
| hi-IN| 印地语 |
| hr | 克罗地亚语 |
| hr-BA| 克罗地亚语(波斯尼亚和黑塞哥维那) |
| hr-HR| 克罗地亚语 |
| hu | 匈牙利语 |
| hu-HU| 匈牙利语 |
| hy | 亚美尼亚语 |
| hy-AM| 亚美尼亚语 |
| id | 印度尼西亚语 |
| id-ID| 印度尼西亚语 |
| is | 冰岛语 |
| is-IS| 冰岛语 |
| it | 意大利语 |
| it-CH| 意大利语(瑞士) |
| it-IT| 意大利语(意大利) |
| ja | 日语 |
| ja-JP| 日语 |
| ka | 格鲁吉亚语 |
| ka-GE| 格鲁吉亚语 |
| kk | 哈萨克语 |
| kk-KZ| 哈萨克语 |
| kn | 卡纳拉语 |
| kn-IN| 卡纳拉语 |
| ko | 朝鲜语 |
| ko-KR| 朝鲜语 |
| kok| 孔卡尼语 |
| kok-IN | 孔卡尼语 |
| ky | 吉尔吉斯语 |
| ky-KG| 吉尔吉斯语(西里尔文) |
| lt | 立陶宛语 |
| lt-LT| 立陶宛语 |
| lv | 拉脱维亚语 |
| lv-LV| 拉脱维亚语 |
| mi | 毛利语 |
| mi-NZ| 毛利语 |
| mk | 马其顿语 |
| mk-MK| 马其顿语(FYROM)|
| mn | 蒙古语 |
| mn-MN| 蒙古语(西里尔文) |
| mr | 马拉地语 |
| mr-IN| 马拉地语 |
| ms | 马来语 |
| ms-BN| 马来语(文莱达鲁萨兰) |
| ms-MY| 马来语(马来西亚) |
| mt | 马耳他语 |
| mt-MT| 马耳他语 |
| nb | 挪威语(伯克梅尔) |
| nb-NO| 挪威语(伯克梅尔)(挪威) |
| nl | 荷兰语 |
| nl-BE| 荷兰语(比利时) |
| nl-NL| 荷兰语(荷兰) |
| nn-NO| 挪威语(尼诺斯克)(挪威) |
| ns | 北梭托语 |
| ns-ZA| 北梭托语 |
| pa | 旁遮普语 |
| pa-IN| 旁遮普语 |
| pl | 波兰语 |
| pl-PL| 波兰语 |
| pt | 葡萄牙语 |
| pt-BR| 葡萄牙语(巴西) |
| pt-PT| 葡萄牙语(葡萄牙) |
| qu | 克丘亚语 |
| qu-BO| 克丘亚语(玻利维亚) |
| qu-EC| 克丘亚语(厄瓜多尔) |
| qu-PE| 克丘亚语(秘鲁) |
| ro | 罗马尼亚语 |
| ro-RO| 罗马尼亚语 |
| ru | 俄语 |
| ru-RU| 俄语 |
| sa | 梵文 |
| sa-IN| 梵文 |
| se | 北萨摩斯语 |
| se-FI| 北萨摩斯语(芬兰) |
| se-FI| 斯科特萨摩斯语(芬兰) |
| se-FI| 伊那里萨摩斯语(芬兰) |
| se-NO| 北萨摩斯语(挪威) |
| se-NO| 律勒欧萨摩斯语(挪威) |
| se-NO| 南萨摩斯语(挪威) |
| se-SE| 北萨摩斯语(瑞典) |
| se-SE| 律勒欧萨摩斯语(瑞典) |
| se-SE| 南萨摩斯语(瑞典) |
| sk | 斯洛伐克语 |
| sk-SK| 斯洛伐克语 |
| sl | 斯洛文尼亚语 |
| sl-SI| 斯洛文尼亚语 |
| sq | 阿尔巴尼亚语 |
| sq-AL| 阿尔巴尼亚语 |
| sr-BA| 塞尔维亚语(拉丁文，波斯尼亚和黑塞哥维那) |
| sr-BA| 塞尔维亚语(西里尔文，波斯尼亚和黑塞哥维那) |
| sr-SP| 塞尔维亚(拉丁) |
| sr-SP| 塞尔维亚(西里尔文) |
| sv | 瑞典语 |
| sv-FI| 瑞典语(芬兰) |
| sv-SE| 瑞典语 |
| sw | 斯瓦希里语 |
| sw-KE| 斯瓦希里语 |
| syr| 叙利亚语 |
| syr-SY | 叙利亚语 |
| ta | 泰米尔语 |
| ta-IN| 泰米尔语 |
| te | 泰卢固语 |
| te-IN| 泰卢固语 |
| th | 泰语 |
| th-TH| 泰语 |
| tl | 塔加路语 |
| tl-PH| 塔加路语(菲律宾) |
| tn | 茨瓦纳语 |
| tn-ZA| 茨瓦纳语 |
| tr | 土耳其语 |
| tr-TR| 土耳其语 |
| ts | 宗加语 |
| tt | 鞑靼语 |
| tt-RU| 鞑靼语 |
| uk | 乌克兰语 |
| uk-UA| 乌克兰语 |
| ur | 乌都语 |
| ur-PK| 乌都语 |
| uz | 乌兹别克语 |
| uz-UZ| 乌兹别克语(拉丁文) |
| uz-UZ| 乌兹别克语(西里尔文) |
| vi | 越南语 |
| vi-VN| 越南语 |
| xh | 班图语 |
| xh-ZA| 班图语 |
| zh | 中文 |
| zh-CN| 中文(简体) |
| zh-HK| 中文(香港) |
| zh-MO| 中文(澳门) |
| zh-SG| 中文(新加坡) |
| zh-TW| 中文(繁体) |
| zu | 祖鲁语 |
| zu-ZA| 祖鲁语 |



<br/>

>   参考链接： https://github.com/nicksnyder/go-i18n



