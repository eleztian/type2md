package test

import "github.com/eleztian/type2md/test/ext"

//go:generate type2md -f ../docs/doc_config.md github.com/eleztian/type2md/test Config

// Config doc.
type Config struct {
	Pre     ext.Hook
	Post    *ext.Hook
	Servers map[string]struct {
		Host string `json:"host,omitempty"`
		Port int    `json:"port" enums:"22,65522" require:"false"`
	} `json:"servers"` // server list
	Test1 `json:",inline"`
	Test  []string // sss
	Test2 map[string]map[int]*Test1
	Test3 [][2]string `json:"test3"`
	C     []interface{}
}

type Test1 struct {
	A string `json:"a"`
}
