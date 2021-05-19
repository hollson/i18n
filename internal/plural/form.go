package plural

// 语言的复数形式，可参考：
// http://cldr.unicode.org/index/cldr-spec/plural-rules
type Form string

// All defined plural forms.
const (
	Invalid Form = ""
	Zero    Form = "zero"
	One     Form = "one"
	Two     Form = "two"
	Few     Form = "few"
	Many    Form = "many"
	Other   Form = "other"
)
