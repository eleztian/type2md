package ext

// Hook hook config.
type Hook struct {
	Name     string            `json:"name" require:"" default:"example"` // hook name
	Commands []string          `json:"commands"`                          // command list
	Envs     map[string]string `json:"envs"`                              // env key map
	Mode     Mode              `json:"mode" default:"1"`                  // run mode
}

// Mode mode define.
type Mode int

const (
	Mode_Q Mode = iota + 1 // mode q
	Mode_A                 // mode a
)
