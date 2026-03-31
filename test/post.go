package test

import "github.com/eleztian/type2md/test/ext"

// Pre hook config.
type Pre struct {
	Name     string            `json:"name" require:"" default:"example"` // hook name
	Commands []string          `json:"commands"`                          // command list
	Envs     map[string]string `json:"envs"`                              // env key map
	Mode     ext.Mode          `json:"mode" default:"1"`                  // run mode
}
