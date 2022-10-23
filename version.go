package main

import (
	"fmt"
	"os"
)

var (
	Version   = "v1.0.0"
	CommitID  = ""
	BuildTime = ""
)

func PrintVersion() {
	fmt.Printf(`%s
---
Parse the source code through the ast syntax tree, 
extract the specified structure definition and 
convert it into a markdown file.
----
Version  : %s
CommitID : %s
BuildTime: %s
Author   : MoreSec CPF 中间件团队
`, os.Args[0], Version, CommitID, BuildTime)
}
