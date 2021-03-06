// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package i18n

import (
	"fmt"
	"text/template"

	"github.com/hollson/i18n/internal"
	"github.com/hollson/i18n/internal/plural"
)

// 消息的可执行模板
type MessageTemplate struct {
	*Message                                           // 消息
	PluralTemplates map[plural.Form]*internal.Template // 模板
}

// 创建消息的可执行模板
func NewMessageTemplate(m *Message) *MessageTemplate {
	pluralTemplates := map[plural.Form]*internal.Template{}
	setPluralTemplate(pluralTemplates, plural.Zero, m.Zero, m.LeftDelim, m.RightDelim)
	setPluralTemplate(pluralTemplates, plural.One, m.One, m.LeftDelim, m.RightDelim)
	setPluralTemplate(pluralTemplates, plural.Two, m.Two, m.LeftDelim, m.RightDelim)
	setPluralTemplate(pluralTemplates, plural.Few, m.Few, m.LeftDelim, m.RightDelim)
	setPluralTemplate(pluralTemplates, plural.Many, m.Many, m.LeftDelim, m.RightDelim)
	setPluralTemplate(pluralTemplates, plural.Other, m.Other, m.LeftDelim, m.RightDelim)
	if len(pluralTemplates) == 0 {
		return nil
	}
	return &MessageTemplate{
		Message:         m,
		PluralTemplates: pluralTemplates,
	}
}

func setPluralTemplate(pluralTemplates map[plural.Form]*internal.Template, pluralForm plural.Form, src, leftDelim, rightDelim string) {
	if src != "" {
		pluralTemplates[pluralForm] = &internal.Template{
			Src:        src,
			LeftDelim:  leftDelim,
			RightDelim: rightDelim,
		}
	}
}

type pluralFormNotFoundError struct {
	pluralForm plural.Form
	messageID  string
}

func (e pluralFormNotFoundError) Error() string {
	return fmt.Sprintf("message %q has no plural form %q", e.messageID, e.pluralForm)
}

// Execute executes the template for the plural form and template data.
func (mt *MessageTemplate) Execute(pluralForm plural.Form, data interface{}, funcs template.FuncMap) (string, error) {
	t := mt.PluralTemplates[pluralForm]
	if t == nil {
		return "", pluralFormNotFoundError{
			pluralForm: pluralForm,
			messageID:  mt.Message.ID,
		}
	}
	return t.Execute(funcs, data)
}
