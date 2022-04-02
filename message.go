// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package i18n

import (
	"fmt"
	"strings"
)

// (多语言)国际化的消息。
//  LDML，即Unicode Locale Data Markup Language(UNICODE语言环境数据标记语言);
//  CLDR，即Unicode Common Locale Data Repository(通用语言环境数据库),可参考：
//   http://cldr.unicode.org/index
//   http://cldr.unicode.org/index/cldr-spec/plural-rules
//   https://unicode-org.github.io/cldr-staging/charts/37/supplemental/language_plural_rules.html
type Message struct {
	// 消息标识(必须唯一)
	ID string

	// 唯一标识了其翻译的消息内容。
	Hash string

	// 消息描述
	Desc string

	// Go模板的左分隔符
	LeftDelim string

	// Go模板的右分隔符
	RightDelim string

	// CLDR复数形式“Zero”的消息内容。
	Zero string

	// CLDR复数形式“One”的消息内容
	One string

	// CLDR复数形式“Two”的消息内容
	Two string

	// 很少,CLDR复数形式“Few”的消息内容
	Few string

	// 许多,CLDR复数形式“Many”的消息内容
	Many string

	// CLDR复数形式“Other”的消息内容
	Other string
}

func (m *Message) String() string {
	return fmt.Sprintf("%+v", *m)
}

// NewMessage parses data and returns a new message.
func NewMessage(data interface{}) (*Message, error) {
	m := &Message{}
	if err := m.unmarshalInterface(data); err != nil {
		return nil, err
	}
	return m, nil
}

// MustNewMessage is similar to NewMessage except it panics if an error happens.
func MustNewMessage(data interface{}) *Message {
	m, err := NewMessage(data)
	if err != nil {
		panic(any(err))
	}
	return m
}

// unmarshalInterface unmarshals a message from data.
func (m *Message) unmarshalInterface(v interface{}) error {
	strdata, err := stringMap(v)
	if err != nil {
		return err
	}
	for k, v := range strdata {
		switch strings.ToLower(k) {
		case "id":
			m.ID = v
		case "description":
			m.Desc = v
		case "hash":
			m.Hash = v
		case "leftdelim":
			m.LeftDelim = v
		case "rightdelim":
			m.RightDelim = v
		case "zero":
			m.Zero = v
		case "one":
			m.One = v
		case "two":
			m.Two = v
		case "few":
			m.Few = v
		case "many":
			m.Many = v
		case "other":
			m.Other = v
		}
	}
	return nil
}

type keyTypeErr struct {
	key interface{}
}

func (err *keyTypeErr) Error() string {
	return fmt.Sprintf("expected key to be a string but got %#v", err.key)
}

type valueTypeErr struct {
	value interface{}
}

func (err *valueTypeErr) Error() string {
	return fmt.Sprintf("unsupported type %#v", err.value)
}

func stringMap(v interface{}) (map[string]string, error) {
	switch value := v.(type) {
	case string:
		return map[string]string{
			"other": value,
		}, nil
	case map[string]string:
		return value, nil
	case map[string]interface{}:
		strdata := make(map[string]string, len(value))
		for k, v := range value {
			err := stringSubmap(k, v, strdata)
			if err != nil {
				return nil, err
			}
		}
		return strdata, nil
	case map[interface{}]interface{}:
		strdata := make(map[string]string, len(value))
		for k, v := range value {
			kstr, ok := k.(string)
			if !ok {
				return nil, &keyTypeErr{key: k}
			}
			err := stringSubmap(kstr, v, strdata)
			if err != nil {
				return nil, err
			}
		}
		return strdata, nil
	default:
		return nil, &valueTypeErr{value: value}
	}
}

func stringSubmap(k string, v interface{}, strdata map[string]string) error {
	if k == "translation" {
		switch vt := v.(type) {
		case string:
			strdata["other"] = vt
		default:
			v1Message, err := stringMap(v)
			if err != nil {
				return err
			}
			for kk, vv := range v1Message {
				strdata[kk] = vv
			}
		}
		return nil
	}

	switch vt := v.(type) {
	case string:
		strdata[k] = vt
		return nil
	case nil:
		return nil
	default:
		return fmt.Errorf("expected value for key %q be a string but got %#v", k, v)
	}
}

// isMessage tells whether the given data is a message, or a map containing
// nested messages.
// A map is assumed to be a message if it contains any of the "reserved" keys:
// "id", "description", "hash", "leftdelim", "rightdelim", "zero", "one", "two", "few", "many", "other"
// with a string value.
// e.g.,
// - {"message": {"description": "world"}} is a message
// - {"message": {"description": "world", "foo": "bar"}} is a message ("foo" key is ignored)
// - {"notmessage": {"description": {"hello": "world"}}} is not
// - {"notmessage": {"foo": "bar"}} is not
func isMessage(v interface{}) bool {
	reservedKeys := []string{"id", "description", "hash", "leftdelim", "rightdelim", "zero", "one", "two", "few", "many", "other"}
	switch data := v.(type) {
	case string:
		return true
	case map[string]interface{}:
		for _, key := range reservedKeys {
			val, ok := data[key]
			if !ok {
				continue
			}
			_, ok = val.(string)
			if !ok {
				continue
			}
			// v is a message if it contains a "reserved" key holding a string value
			return true
		}
	case map[interface{}]interface{}:
		for _, key := range reservedKeys {
			val, ok := data[key]
			if !ok {
				continue
			}
			_, ok = val.(string)
			if !ok {
				continue
			}
			// v is a message if it contains a "reserved" key holding a string value
			return true
		}
	}
	return false
}
