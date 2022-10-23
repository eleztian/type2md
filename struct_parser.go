package main

import (
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"
	"go/types"
	"strings"

	"golang.org/x/tools/go/loader"
)

type Parser struct {
	Tag             string
	RootModPath     string
	program         *loader.Program
	loaderCfg       *loader.Config
	fileImportNamed map[*ast.File]map[string]string
	fileTypeDoc     map[*ast.File]map[string]string
}

func NewParser(modPath string) (*Parser, error) {
	cfg := &loader.Config{
		ParserMode: goparser.ParseComments,
	}
	cfg.Import(modPath)

	p, err := cfg.Load()
	if err != nil {
		return nil, err
	}
	res := &Parser{
		Tag:             "json",
		RootModPath:     "",
		program:         p,
		loaderCfg:       cfg,
		fileImportNamed: map[*ast.File]map[string]string{},
		fileTypeDoc:     map[*ast.File]map[string]string{},
	}

	return res, nil
}

func (p *Parser) getFileImportNamedMap(file *ast.File) map[string]string {
	namedImportMap := map[string]string{}
	for _, imp := range file.Imports {
		modPath := strings.Trim(imp.Path.Value, "\"")
		if imp.Name == nil {
			namedImportMap[p.program.Package(modPath).Pkg.Name()] = modPath
		} else {
			namedImportMap[imp.Name.Name] = modPath
		}
	}
	return namedImportMap
}

func (p *Parser) getFileTypeDocMap(file *ast.File) map[string]string {
	typeDocMap := map[string]string{}
	for _, imp := range file.Decls {
		genDecl, ok := imp.(*ast.GenDecl)
		if !ok {
			continue
		}
		if genDecl.Tok != token.TYPE {
			continue
		}
		var comment = ""
		if genDecl.Doc != nil {
			comment = TrimComment(genDecl.Doc.Text())
			if comment != "" {
				comment += "\n"
			}
		}
		for _, spec := range genDecl.Specs {
			if specType, ok := spec.(*ast.TypeSpec); ok {
				var subComment = GetDescribeFromComment(specType.Doc, specType.Comment)
				typeDocMap[specType.Name.Name] = comment + subComment
			}
		}
	}

	return typeDocMap
}

func (p *Parser) parseFileStruct(
	modPath string,
	file *ast.File,
	exParseMap map[string]struct{},
	objName string,
	objStructType *ast.StructType) map[string]StructInfo {

	res := make(map[string]StructInfo, 0)
	typeKey := TypeKey(modPath, objName)

	namedImportMap := p.fileImportNamed[file]
	typeDocMap := p.fileTypeDoc[file]

	structInfo := p.parseStruct(objStructType, typeDocMap[objName])
	res[typeKey] = structInfo
	for idx, field := range structInfo.Fields {
		if field.Reference == "" {
			continue
		}
		subModPath := ""
		if field.Reference == "." {
			subModPath = modPath
		} else {
			subModPath = namedImportMap[field.Reference]
		}
		structInfo.Fields[idx].Reference = subModPath

		key := TypeKey(subModPath, field.Type)
		if _, ok := exParseMap[key]; !ok {
			for typeName, fields := range p.Parse(subModPath, field.Type) {
				res[typeName] = fields
			}
			exParseMap[key] = struct{}{}
		}
	}

	return res
}

func (p *Parser) Parse(modPath string, typeName string) map[string]StructInfo {
	pktInfo := p.program.Package(modPath)
	if pktInfo == nil {
		return map[string]StructInfo{}
	}
	res := make(map[string]StructInfo, 0)
	exParseMap := make(map[string]struct{})

	for _, file := range pktInfo.Files {
		if _, ok := p.fileImportNamed[file]; !ok {
			p.fileImportNamed[file] = p.getFileImportNamedMap(file)
		}
		if _, ok := p.fileTypeDoc[file]; !ok {
			p.fileTypeDoc[file] = p.getFileTypeDocMap(file)
		}
		obj := file.Scope.Lookup(typeName)
		if obj == nil {
			continue
		}
		objTypeSpec := obj.Decl.(*ast.TypeSpec)
		var objType = objTypeSpec.Type
		if v, ok := objType.(*ast.StarExpr); ok {
			objType = v.X
		}
		switch ost := objType.(type) {
		case *ast.StructType:
			structMaps := p.parseFileStruct(modPath, file, exParseMap, obj.Name, ost)
			for name, info := range structMaps {
				res[name] = info
			}
		case *ast.Ident:
			if ost.Obj == nil { // 基本数据类型, 枚举值
				enums := p.getEnumTypeValues(&pktInfo.Info, typeName, file.Decls)
				res[TypeKey(modPath, typeName)] = StructInfo{
					Describe: p.fileTypeDoc[file][typeName],
					Enums: &EnumInfo{
						Type:  ost.Name,
						Names: enums,
					}}
			}
		case *ast.SelectorExpr:
			subModPath := p.fileImportNamed[file][ost.X.(*ast.Ident).Name]
			subTypeName := ost.Sel.Name
			structMaps := p.Parse(subModPath, subTypeName)
			oldKey := TypeKey(subModPath, subTypeName)
			info := structMaps[oldKey]
			info.Describe += fmt.Sprintf("alias `%s`", oldKey)

			structMaps[TypeKey(modPath, typeName)] = info
			delete(structMaps, oldKey)
			return structMaps
		}
	}
	return res
}

func (p *Parser) getEnumTypeValues(pktInfo *types.Info, name string, decls []ast.Decl) [][2]string {
	res := make([][2]string, 0)

EXIT:
	for _, obj := range decls {

		if genDecl, ok := obj.(*ast.GenDecl); ok {
			if genDecl.Tok != token.CONST || len(genDecl.Specs) == 0 {
				continue
			}
			firstSpec := genDecl.Specs[0]
			if valueSpec, ok := firstSpec.(*ast.ValueSpec); ok {
				t, ok := valueSpec.Type.(*ast.Ident)
				if !ok || t.Name != name {
					continue
				}
			} else {
				continue
			}

			for _, s := range genDecl.Specs {
				v := s.(*ast.ValueSpec) // safe because decl.Tok == token.CONST
				for _, name := range v.Names {
					c := pktInfo.ObjectOf(name).(*types.Const)
					res = append(res, [2]string{
						c.Val().ExactString(),
						GetDescribeFromComment(v.Doc, v.Comment),
					})
				}
			}
			break EXIT
		}
	}

	return res
}

func (p *Parser) parseTypeExpr(obj ast.Expr) []FieldInfo {
	if v, ok := obj.(*ast.StarExpr); ok {
		obj = v.X
	}
	var res []FieldInfo
	switch ot := obj.(type) {
	case *ast.SelectorExpr:
		res = []FieldInfo{{Type: ot.Sel.Name, Reference: ot.X.(*ast.Ident).Name, skipNum: 1}}
	case *ast.Ident:
		field := FieldInfo{
			Type: ot.Name,
		}
		if ot.Obj != nil {
			field.Reference = "."
		}
		res = append(res, field)
	case *ast.StructType:
		res = p.parseStruct(ot, "").Fields
	case *ast.MapType:
		prefix := fmt.Sprintf("{%s}.", ot.Key)
		res = p.parseTypeExpr(ot.Value)
		for idx := range res {
			res[idx].Name = prefix + res[idx].Name
		}
	case *ast.SliceExpr:
		prefix := "[]."
		res = p.parseTypeExpr(ot.X)
		for idx := range res {
			res[idx].Name = prefix + res[idx].Name
		}
	case *ast.ArrayType:
		prefix := "[]"
		t, ok := ot.Len.(*ast.BasicLit)
		if ok {
			prefix = fmt.Sprintf("[%s]", t.Value)
		}
		res = p.parseTypeExpr(ot.Elt)
		for idx := range res {
			res[idx].Name = prefix + res[idx].Name
		}
	}
	return res
}

func (p *Parser) parseStruct(objStructType *ast.StructType, desc string) StructInfo {
	res := StructInfo{
		Describe: desc,
		Fields:   make([]FieldInfo, 0, len(objStructType.Fields.List)),
	}

	for _, f := range objStructType.Fields.List {
		if f.Names[0].Name[0] <= 'Z' && f.Names[0].Name[0] >= 'A' {
			res.Fields = append(res.Fields, p.parseStructField(f)...)
		}
	}

	return res
}

func (p *Parser) parseStructField(f *ast.Field) []FieldInfo {
	res := make([]FieldInfo, 0)

	tagStr := ""
	if f.Tag != nil {
		tagStr = strings.Trim(f.Tag.Value, "`")
	}
	tagInfo := ParseStructTag(p.Tag, tagStr)
	if tagInfo.Name == "-" {
		return res
	}

	baseField := FieldInfo{
		Name:    tagInfo.Name,
		Default: tagInfo.Default,
		Require: tagInfo.Require,
		Enums: EnumInfo{
			Names: tagInfo.Enums,
		},
	}
	if baseField.Name == "" {
		baseField.Name = f.Names[0].Name
	}

	baseField.Describe += GetDescribeFromComment(f.Doc, f.Comment)

	if s, ok := f.Type.(*ast.StarExpr); ok {
		f.Type = s.X
	}
	switch tt := f.Type.(type) {
	case *ast.Ident:
		baseField.Type = tt.Name
		if tt.Obj != nil {
			baseField.Reference = "."
		}
		if tagInfo.Inline && baseField.Reference == "." {
			for _, field := range p.parseTypeExpr(tt.Obj.Decl.(*ast.TypeSpec).Type) {
				res = append(res, field)
			}
		} else {
			res = append(res, baseField)
		}
	case *ast.StructType:
		for _, f := range p.parseStruct(tt, baseField.Describe).Fields {
			if tagInfo.Inline {
				res = append(res, f)
			} else {
				if f.skipNum != 0 {
					field := baseField.Copy()
					field.Type = f.Type
					if f.Name != "" {
						field.Name += "." + f.Name
					}
					res = append(res, field)
				} else {
					f.Name = baseField.Name + "." + f.Name
					res = append(res, f)
				}
			}
		}
	case *ast.MapType:
		baseField.Name += fmt.Sprintf(".{%s}", tt.Key)

		var subFields = p.parseTypeExpr(tt.Value)
		for _, f := range subFields {
			if f.skipNum != 0 {
				field := baseField.Copy()
				field.Type = f.Type
				field.Reference = f.Reference
				if f.Name != "" {
					field.Name += "." + f.Name
				}
				res = append(res, field)
			} else {
				f.Name = baseField.Name + "." + f.Name
				res = append(res, f)
			}
		}
	case *ast.SliceExpr:
		baseField.Name += ".[]"
		for _, f := range p.parseTypeExpr(tt.X) {
			if f.skipNum != 0 {
				field := baseField.Copy()
				field.Type = f.Type
				field.Reference = f.Reference
				if f.Name != "" {
					field.Name += "." + f.Name
				}
				res = append(res, field)
			} else {
				f.Name = baseField.Name + "." + f.Name
				res = append(res, f)
			}
		}
	case *ast.ArrayType:
		t, ok := tt.Len.(*ast.BasicLit)
		if ok {
			baseField.Name += fmt.Sprintf(".[%s]", t.Value)
		} else {
			baseField.Name += ".[]"
		}
		for _, f := range p.parseTypeExpr(tt.Elt) {
			if f.skipNum != 0 {
				field := baseField.Copy()
				field.Type = f.Type
				field.Reference = f.Reference
				if f.Name != "" {
					field.Name += "." + f.Name
				}
				res = append(res, field)
			} else {
				f.Name = baseField.Name + "." + f.Name
				res = append(res, f)
			}
		}
	case *ast.SelectorExpr:
		baseField.Type = tt.Sel.Name
		baseField.Reference = tt.X.(*ast.Ident).Name
		res = append(res, baseField)
	}

	return res
}

func TypeKey(modPath string, typeName string) string {
	return modPath + "." + typeName
}
