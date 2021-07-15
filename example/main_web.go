// Command example runs a sample webserver that uses go-i18n/v2/i18n.
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/hollson/i18n"
	"golang.org/x/text/language"
)

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
