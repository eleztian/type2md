package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

var (
	dstModPath string
	dstTypes   []string
)
var (
	fileName    = flag.String("f", "", "file path")
	title       = flag.String("t", "", "file title")
	tagType     = flag.String("tag", "json", "struct tag name")
	showVersion = flag.Bool("v", false, "show version")
)

func main() {
	flag.Parse()

	if *showVersion {
		PrintVersion()
		os.Exit(0)
	}
	dstModPath = os.Args[len(os.Args)-2]
	dstTypes = strings.Split(os.Args[len(os.Args)-1], ",")

	rootMod, _, _ := getModPath()
	log.Println("Current Module:", rootMod)

	parser, err := NewParser(*tagType, dstModPath)
	if err != nil {
		log.Fatalf("new parser %v", err)
		return
	}

	for _, tp := range dstTypes {
		log.Printf("start generate %s.%s\n", dstModPath, tp)
		fs := parser.Parse(dstModPath, tp)

		md := Markdown{
			Title:          *title,
			MainStructName: dstModPath + "." + tp,
			ObjTitleFunc: func(modPath string, typeName string) string {
				modPath = strings.TrimPrefix(modPath, dstModPath)
				modPath = strings.TrimPrefix(modPath, rootMod)
				modPath = strings.TrimLeft(modPath, "./")
				if modPath == "" {
					return typeName
				} else {
					return fmt.Sprintf("%s.%s", modPath, typeName)
				}
			},
		}
		if md.Title == "" {
			md.Title = tp + " Doc"
		}

		data := md.Generate(fs)

		filename := *fileName
		if filename == "" {
			filename = tp + "_doc.md"
		}
		log.Printf("start to save to %s\n", filename)
		_ = os.WriteFile(filename, data, 0655)
	}
}

func getModPath() (string, string, error) {
	current, _ := os.Getwd()
	var path, _ = filepath.Abs(current)
	for {
		_, err := os.Stat(filepath.Join(path, "go.mod"))
		if err != nil {
			if !os.IsNotExist(err) {
				return "", "", err
			}
			path, _ = filepath.Abs(filepath.Join(path, "../"))
			if path == "" {
				return "", "", nil
			}
			continue
		}
		content, err := os.ReadFile(filepath.Join(path, "go.mod"))
		if err != nil {
			return "", "", err
		}
		f, err := modfile.Parse("go.mod", content, nil)
		if err != nil {
			return "", "", err
		}

		return f.Module.Mod.Path, strings.Replace(current, path, f.Module.Mod.Path, 1), nil
	}
}
