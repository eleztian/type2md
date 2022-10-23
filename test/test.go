package test

import "github.com/eleztian/type2md/test/ext"

//go:generate type2md -f ../docs/doc_config.md github.com/eleztian/type2md/test Config

// Config doc.
type Config struct {
	Pre     ext.Hook
	Post    *ext.Hook
	Servers map[string]struct {
		Host string `json:"host"`
		Port int    `json:"port" enums:"22,65522" require:""`
	} `json:"servers"` // server list
	Test1 `json:",inline"`
	Test  []string // sss
}

type Test1 struct {
	A string `json:"a"`
}
