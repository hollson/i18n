# i18n 
**i18n**是一个Go语言版本的多语言程序包

- 支持Unicode[通用语言环境数据存储库(CLDR)](https://www.unicode.org/cldr/charts/28/supplemental/language_plural_rules.html)；
- 支持使用具有命名变量的字符串 `Template`语法 ；
-   支持多种消息文件格式, 如:  TOML、JSON、YAML；

# i18n包

i18n软件包提供了根据一组区域设置首选项查找消息的支持。 

```go
import "github.com/hollson/i18n"
```

创建一个捆绑包，以在您的应用程序的整个生命周期内使用。 

```go
bundle := i18n.NewBundle(language.English)
```

在初始化期间将翻译加载到您的捆绑软件中。 

```go
bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
bundle.LoadMessageFile("es.toml")
```

创建一个本地化程序以用于一组语言首选项。 

```go
func(w http.ResponseWriter, r *http.Request) {
    lang := r.FormValue("lang")
    accept := r.Header.Get("Accept-Language")
    localizer := i18n.NewLocalizer(bundle, lang, accept)
}
```

使用本地化程序查找消息。 

```go
localizer.Localize(&i18n.LocalizeConfig{
    DefaultMessage: &i18n.Message{
        ID: "PersonCats",
        One: "{{.Name}} has {{.Count}} cat.",
        Other: "{{.Name}} has {{.Count}} cats.",
    },
    TemplateData: map[string]interface{}{
        "Name": "Nick",
        "Count": 2,
    },
    PluralCount: 2,
}) // Nick has 2 cats.
```

# goi18n工具

goi18n命令管理i18n软件包使用的消息文件。 

```bash
go get -u github.com/hollson/i18n
goi18n -help
```

## 提取消息

用  `goi18n extract` 将Go源文件中的所有`i18n.Message`结构体文字提取到消息文件中进行翻译。 

```toml
# active.en.toml
[PersonCats]
description = "The number of cats a person has"
one = "{{.Name}} has {{.Count}} cat."
other = "{{.Name}} has {{.Count}} cats."
```

## 翻译消息

1.创建一个空的消息文件，如 `translate.es.toml`

2.执行`goi18n merge active.en.toml translate.es.toml` 填充  `translate.es.toml` 与要翻译的消息。 
```toml
# translate.es.toml
[HelloPerson]
hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
other = "Hello {{.Name}}"
 ```

3.`translate.es.toml` 已翻译，将其重命名为  `active.es.toml`. 
```toml
# active.es.toml
[HelloPerson]
hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
other = "Hola {{.Name}}"
```

4.加载`active.es.toml`文件。 

```go
bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
bundle.LoadMessageFile("active.es.toml")
```


> https://github.com/nicksnyder/go-i18n

