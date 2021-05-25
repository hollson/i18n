// Copyright 2021 Hollson. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"golang.org/x/text/language"
)

type languageTag language.Tag

func (lt *languageTag) Set(value string) error {
	t, err := language.Parse(value)
	if err != nil {
		return err
	}
	*lt = languageTag(t)
	return nil
}

func (lt languageTag) String() string {
	return lt.Tag().String()
}

func (lt *languageTag) Tag() language.Tag {
	tag := language.Tag(*lt)
	if tag.IsRoot() {
		return language.English
	}
	return tag
}
