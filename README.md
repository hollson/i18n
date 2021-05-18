go-i18n是一个Go [程序包](github.com/hollson/go-i18n#package-i18n)和一个 [命令](https://github.com/hollson/go-i18n#command-goi18n)，可以帮助您将Go程序翻译成多种语言。 

-   支持 

    复数字符串 

    所有200种以上语言的 

    Unicode通用语言环境数据存储库（CLDR）中 

    -   代码和测试是 [自动生成的 ](https://github.com/nicksnyder/go-i18n/tree/main/v2/internal/plural/codegen)从 [CLDR数据 ](http://cldr.unicode.org/index/downloads)。 

-   支持使用 具有命名变量的字符串 [文本/模板 ](http://golang.org/pkg/text/template/)语法的 。 

-   支持任何格式的消息文件（例如JSON，TOML，YAML）。 

## 套餐i18n

i18n软件包提供了根据一组区域设置首选项查找消息的支持。 

```
import "github.com/nicksnyder/go-i18n/v2/i18n"
```

创建一个捆绑包，以在您的应用程序的整个生命周期内使用。 

```
bundle := i18n.NewBundle(language.English)
```

在初始化期间将翻译加载到您的捆绑软件中。 

```
bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
bundle.LoadMessageFile("es.toml")
```

创建一个本地化程序以用于一组语言首选项。 

```
func(w http.ResponseWriter, r *http.Request) {
    lang := r.FormValue("lang")
    accept := r.Header.Get("Accept-Language")
    localizer := i18n.NewLocalizer(bundle, lang, accept)
}
```

使用本地化程序查找消息。 

```
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

## 命令goi18n

goi18n命令管理i18n软件包使用的消息文件。 

```bash
go get -u github.com/nicksnyder/go-i18n/v2/goi18n
goi18n -help
```

### 提取消息

用  `goi18n extract` 将Go源文件中的所有i18n.Message结构体文字提取到消息文件中进行翻译。 

```toml
# active.en.toml
[PersonCats]
description = "The number of cats a person has"
one = "{{.Name}} has {{.Count}} cat."
other = "{{.Name}} has {{.Count}} cats."
```

### 翻译新语言

1.  为您要添加的语言创建一个空的消息文件（例如，  `translate.es.toml`). 

2.  跑步  `goi18n merge active.en.toml translate.es.toml` 填充  `translate.es.toml` 与要翻译的消息。 

    ```toml
    # translate.es.toml
    [HelloPerson]
    hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
    other = "Hello {{.Name}}"
    ```

3.  后  `translate.es.toml` 已翻译，将其重命名为  `active.es.toml`. 

    ```toml
    # active.es.toml
    [HelloPerson]
    hash = "sha1-5b49bfdad81fedaeefb224b0ffc2acc58b09cff5"
    other = "Hola {{.Name}}"
    ```

4.  加载  `active.es.toml` 进入您的捆绑包。 

    ```toml
    bundle.RegisterUnmarshalFunc("toml", toml.Unmarshal)
    bundle.LoadMessageFile("active.es.toml")
    ```

### 翻译新讯息 

如果您在程序中添加了新消息： 

1.  跑步  `goi18n extract` 更新  `active.en.toml` 与新消息。 
2.  跑步  `goi18n merge active.*.toml` 生成更新  `translate.*.toml` 文件。 
3.  翻译  `translate.*.toml` 文件。 
4.  跑步  `goi18n merge active.*.toml translate.*.toml` 将转换后的消息合并到活动消息文件中。 

## 有关更多信息和示例： 

-   阅读 [文档 ](https://godoc.org/github.com/nicksnyder/go-i18n/v2)。 
-   查看 [代码示例 ](https://github.com/nicksnyder/go-i18n/blob/main/v2/i18n/example_test.go)和 [测试 ](https://github.com/nicksnyder/go-i18n/blob/main/v2/i18n/localizer_test.go)。 
-   看一个示例 [应用程序 ](https://github.com/nicksnyder/go-i18n/tree/main/v2/example)。 

## 执照 

go-i18n在MIT许可下可用。 有关 请参见 [LICENSE ](https://github.com/hollson/go-i18n/blob/main/LICENSE)更多信息， 文件。 