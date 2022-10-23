package main

import (
	"reflect"
	"strings"
)

type TagInfo struct {
	Name    string
	Default string
	Require bool
	Inline  bool
	Enums   [][2]string
}

func ParseStructTag(tagType string, tagStr string) *TagInfo {
	res := &TagInfo{}
	tag := strings.TrimSpace(reflect.StructTag(tagStr).Get(tagType))
	if tag != "" {
		sp := strings.Split(tag, ",")
		res.Name = strings.TrimSpace(sp[0])
		for _, item := range sp[1:] {
			item = strings.TrimSpace(item)
			if item == "inline" {
				res.Inline = true
			}
		}
	}
	res.Default, _ = reflect.StructTag(tagStr).Lookup("default")
	_, res.Require = reflect.StructTag(tagStr).Lookup("require")
	enumsStr, _ := reflect.StructTag(tagStr).Lookup("enums")

	if enumsStr != "" {
		res.Enums = make([][2]string, 0)
		sp := strings.Split(enumsStr, ",")
		for _, k := range sp {
			var (
				value string
				desc  string
			)
			idx := strings.Index(k, ":")
			if idx >= 0 {
				value = k[:idx]
				desc = k[idx+1:]
			} else {
				value = k
			}
			res.Enums = append(res.Enums, [2]string{
				value, desc,
			})
		}
	}

	return res
}
