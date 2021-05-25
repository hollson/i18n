package main

import (
	"crypto/sha1"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/hollson/i18n"
	"github.com/hollson/i18n/internal"
	"github.com/hollson/i18n/internal/plural"

	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v2"
)

func usageMerge() {
	fmt.Fprintf(os.Stderr, `合并消息文件:

    合并多语言消息文件,文件名必须具有受支持格式的后缀(例如“ .json”),并包含RFC 5646定义的有效语言标签(例如“ en-us”,“ fr”,“ zh-hant”等)

Usage: i18n_cli merge [Option]... <Param>...

Option:
    -source
      翻译来自该语言的消息, 如: en(默认),en-US,zh-Hant-CN
    -out
      文件输出路径
    -format
      输出消息的文件格式,仅支持: toml(默认), json, yaml

Example: 
    i18n_cli merge active.en.toml active.zh.toml

`)
}

type mergeCommand struct {
	msgFiles []string
	source   languageTag
	out      string
	format   string
}

func (mc *mergeCommand) name() string {
	return "merge"
}

func (mc *mergeCommand) parse(args []string) error {
	flags := flag.NewFlagSet("merge", flag.ExitOnError)
	flags.Usage = usageMerge

	flags.Var(&mc.source, "source", "en")
	flags.StringVar(&mc.out, "out", ".", "")
	flags.StringVar(&mc.format, "format", "toml", "")
	if err := flags.Parse(args); err != nil {
		return err
	}

	mc.msgFiles = flags.Args()
	return nil
}

func (mc *mergeCommand) execute() error {
	if len(mc.msgFiles) < 1 {
		usageMerge()
		return nil
	}
	inFiles := make(map[string][]byte)
	for _, path := range mc.msgFiles {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}
		inFiles[path] = content
	}
	ops, err := merge(inFiles, mc.source.Tag(), mc.out, mc.format)
	if err != nil {
		return err
	}
	for path, content := range ops.writeFiles {
		if err := ioutil.WriteFile(path, content, 0666); err != nil {
			return err
		}
	}
	for _, path := range ops.deleteFiles {
		// Ignore error since it isn't guaranteed to exist.
		os.Remove(path)
	}
	return nil
}

type fileSystemOp struct {
	writeFiles  map[string][]byte
	deleteFiles []string
}

func merge(msgFiles map[string][]byte, sourceLanguageTag language.Tag, out, outputFormat string) (*fileSystemOp, error) {
	unmerged := make(map[language.Tag][]map[string]*i18n.MessageTemplate)
	sourceMessageTemplates := make(map[string]*i18n.MessageTemplate)
	unmarshalFuncs := map[string]i18n.UnmarshalFunc{
		"json": json.Unmarshal,
		"toml": toml.Unmarshal,
		"yaml": yaml.Unmarshal,
	}
	for path, content := range msgFiles {
		mf, err := i18n.ParseMessageFileBytes(content, path, unmarshalFuncs)
		if err != nil {
			return nil, fmt.Errorf("failed to load message file %s: %s", path, err)
		}
		templates := map[string]*i18n.MessageTemplate{}
		for _, m := range mf.Messages {
			template := i18n.NewMessageTemplate(m)
			if template == nil {
				continue
			}
			templates[m.ID] = template
		}
		if mf.Tag == sourceLanguageTag {
			for _, template := range templates {
				if sourceMessageTemplates[template.ID] != nil {
					return nil, fmt.Errorf("multiple source translations for id %q", template.ID)
				}
				template.Hash = hash(template)
				sourceMessageTemplates[template.ID] = template
			}
		}
		unmerged[mf.Tag] = append(unmerged[mf.Tag], templates)
	}

	if len(sourceMessageTemplates) == 0 {
		return nil, fmt.Errorf("no messages found for source locale %s", sourceLanguageTag)
	}

	pluralRules := plural.DefaultRules()
	all := make(map[language.Tag]map[string]*i18n.MessageTemplate)
	all[sourceLanguageTag] = sourceMessageTemplates
	for _, srcTemplate := range sourceMessageTemplates {
		for dstLangTag, messageTemplates := range unmerged {
			if dstLangTag == sourceLanguageTag {
				continue
			}
			pluralRule := pluralRules.Rule(dstLangTag)
			if pluralRule == nil {
				// Non-standard languages not supported because
				// we don't know if translations are complete or not.
				continue
			}
			if all[dstLangTag] == nil {
				all[dstLangTag] = make(map[string]*i18n.MessageTemplate)
			}
			dstMessageTemplate := all[dstLangTag][srcTemplate.ID]
			if dstMessageTemplate == nil {
				dstMessageTemplate = &i18n.MessageTemplate{
					Message: &i18n.Message{
						ID:   srcTemplate.ID,
						Desc: srcTemplate.Desc,
						Hash: srcTemplate.Hash,
					},
					PluralTemplates: make(map[plural.Form]*internal.Template),
				}
				all[dstLangTag][srcTemplate.ID] = dstMessageTemplate
			}

			// Check all unmerged message templates for this message id.
			for _, messageTemplates := range messageTemplates {
				unmergedTemplate := messageTemplates[srcTemplate.ID]
				if unmergedTemplate == nil {
					continue
				}
				// Ignore empty hashes for v1 backward compatibility.
				if unmergedTemplate.Hash != "" && unmergedTemplate.Hash != srcTemplate.Hash {
					// This was translated from different content so discard.
					continue
				}

				// Merge in the translated messages.
				for pluralForm := range pluralRule.PluralForms {
					dt := unmergedTemplate.PluralTemplates[pluralForm]
					if dt != nil && dt.Src != "" {
						dstMessageTemplate.PluralTemplates[pluralForm] = dt
					}
				}
			}
		}
	}

	translate := make(map[language.Tag]map[string]*i18n.MessageTemplate)
	active := make(map[language.Tag]map[string]*i18n.MessageTemplate)
	for langTag, messageTemplates := range all {
		active[langTag] = make(map[string]*i18n.MessageTemplate)
		if langTag == sourceLanguageTag {
			active[langTag] = messageTemplates
			continue
		}
		pluralRule := pluralRules.Rule(langTag)
		if pluralRule == nil {
			// Non-standard languages not supported because
			// we don't know if translations are complete or not.
			continue
		}
		for _, messageTemplate := range messageTemplates {
			srcMessageTemplate := sourceMessageTemplates[messageTemplate.ID]
			activeMessageTemplate, translateMessageTemplate := activeDst(srcMessageTemplate, messageTemplate, pluralRule)
			if translateMessageTemplate != nil {
				if translate[langTag] == nil {
					translate[langTag] = make(map[string]*i18n.MessageTemplate)
				}
				translate[langTag][messageTemplate.ID] = translateMessageTemplate
			}
			if activeMessageTemplate != nil {
				active[langTag][messageTemplate.ID] = activeMessageTemplate
			}
		}
	}

	writeFiles := make(map[string][]byte, len(translate)+len(active))
	for langTag, messageTemplates := range translate {
		path, content, err := writeFile(out, "translate", langTag, outputFormat, messageTemplates, false)
		if err != nil {
			return nil, err
		}
		writeFiles[path] = content
	}
	deleteFiles := []string{}
	for langTag, messageTemplates := range active {
		path, content, err := writeFile(out, "active", langTag, outputFormat, messageTemplates, langTag == sourceLanguageTag)
		if err != nil {
			return nil, err
		}
		if len(content) > 0 {
			writeFiles[path] = content
		} else {
			deleteFiles = append(deleteFiles, path)
		}
	}
	return &fileSystemOp{writeFiles: writeFiles, deleteFiles: deleteFiles}, nil
}

// activeDst returns the active part of the dst and whether dst is a complete translation of src.
func activeDst(src, dst *i18n.MessageTemplate, pluralRule *plural.Rule) (active *i18n.MessageTemplate, translateMessageTemplate *i18n.MessageTemplate) {
	pluralForms := pluralRule.PluralForms
	if len(src.PluralTemplates) == 1 {
		pluralForms = map[plural.Form]struct{}{
			plural.Other: {},
		}
	}
	for pluralForm := range pluralForms {
		dt := dst.PluralTemplates[pluralForm]
		if dt == nil || dt.Src == "" {
			if translateMessageTemplate == nil {
				translateMessageTemplate = &i18n.MessageTemplate{
					Message: &i18n.Message{
						ID:   src.ID,
						Desc: src.Desc,
						Hash: src.Hash,
					},
					PluralTemplates: make(map[plural.Form]*internal.Template),
				}
			}
			translateMessageTemplate.PluralTemplates[pluralForm] = src.PluralTemplates[plural.Other]
			continue
		}
		if active == nil {
			active = &i18n.MessageTemplate{
				Message: &i18n.Message{
					ID:   src.ID,
					Desc: src.Desc,
					Hash: src.Hash,
				},
				PluralTemplates: make(map[plural.Form]*internal.Template),
			}
		}
		active.PluralTemplates[pluralForm] = dt
	}
	return
}

func hash(t *i18n.MessageTemplate) string {
	h := sha1.New()
	_, _ = io.WriteString(h, t.Desc)
	_, _ = io.WriteString(h, t.PluralTemplates[plural.Other].Src)
	return fmt.Sprintf("sha1-%x", h.Sum(nil))
}
