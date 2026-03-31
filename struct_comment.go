package main

import (
	"go/ast"
	"strings"
)

func TrimComment(src string) string {
	res := strings.Trim(src, "\n ")
	if len(res) != 0 && res[len(res)-1] != '.' {
		res += "."
	}
	return res
}

func GetDescribeFromComment(doc *ast.CommentGroup, comment *ast.CommentGroup) string {
	res := ""
	if doc != nil {
		res += TrimComment(doc.Text())
	}
	if comment != nil {
		res += TrimComment(comment.Text())
	}
	// skip comments starting with plus sign, known as "markers" and used for example by kubebuilder
	if strings.HasPrefix(res, "+") {
		res = ""
	}

	return res
}
