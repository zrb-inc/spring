package parser

import (
	"go/ast"
	"go/token"
	"regexp"

	"github.com/zrb-inc/spring/internals/definition"
)

var reAnnotation = regexp.MustCompile(`@.+[\((.*)\)]?`)
var reAnnotationName = regexp.MustCompile(`@(\w*)`)
var reAnnotationPayload = regexp.MustCompile(`\((.*)\)`)

func GetRegexValue(comment string, re *regexp.Regexp) string {
	matches := re.FindStringSubmatch(comment)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

func GetRegexValueSlice(comment string, re *regexp.Regexp) []string {
	return re.FindStringSubmatch(comment)
}

func GetAnnotationBody(comment string) []string {
	return GetRegexValueSlice(comment, reAnnotation)
}

func GetAnnotationName(annotationString string) string {
	return GetRegexValue(annotationString, reAnnotationName)
}

func GetAnnotationPayload(annotationString string) string {
	return GetRegexValue(annotationString, reAnnotationPayload)
}

type Collector struct {
	Ast *ast.File

	Imports []*ast.GenDecl

	Annotations []*AnnotationHolder
}

func (c *Collector) GetAllAnnotationsString() []string {
	strs := []string{}
	for _, v := range c.Annotations {
		strs = append(strs, v.Annotation.Name)
	}
	return strs
}

func (c *Collector) pushHolder(h *AnnotationHolder) {
	c.Annotations = append(c.Annotations, h)
}

func (c *Collector) pushImport(h *ast.GenDecl) {
	c.Imports = append(c.Imports, h)
}

type AnnotationHolder struct {
	Annotation *definition.Annotation
	RelateNode interface{}
	RootNode   interface{}
}

func NewCollector(a *ast.File) *Collector {
	c := &Collector{
		Ast: a,
	}
	c.Process()
	return c
}

func (c *Collector) Process() {
	for _, dec := range c.Ast.Decls {
		switch t := dec.(type) {
		case *ast.GenDecl:
			c.parseGenDecl(t)
		case *ast.FuncDecl:
		}
	}
}

func (c *Collector) parseFuncDecl(fun *ast.FuncDecl) {
	fs := fun.Doc
	if fs == nil {
		return
	}
	for _, ci := range fs.List {
		comment := ci.Text
		annotaionRaws := GetAnnotationBody(comment)
		for _, anno := range annotaionRaws {
			name := GetAnnotationName(anno)
			payload := GetAnnotationPayload(anno)
			annotation := definition.NewAnnotationStruct(name, payload)
			c.pushHolder(&AnnotationHolder{
				Annotation: annotation,
				RelateNode: fun,
				RootNode:   fun,
			})
		}
	}

}

//parse gen (import and type struct) type
func (c *Collector) parseGenDecl(gen *ast.GenDecl) {

	if gen.Tok == token.IMPORT {
		c.pushImport(gen)
		return
	}

	doc := gen.Doc
	if doc != nil {
		for _, cm := range doc.List {
			comment := cm.Text
			annotaionRaws := GetAnnotationBody(comment)
			for _, anno := range annotaionRaws {
				name := GetAnnotationName(anno)
				payload := GetAnnotationPayload(anno)
				annotation := definition.NewAnnotationStruct(name, payload)
				c.pushHolder(&AnnotationHolder{
					Annotation: annotation,
					RelateNode: gen,
					RootNode:   gen,
				})
			}
		}
	}

	//filed comment
	for _, spec := range gen.Specs {
		switch t := spec.(type) {
		case *ast.TypeSpec:
			switch typeT := t.Type.(type) {
			case *ast.StructType:
				for _, f := range typeT.Fields.List {
					fcomment := f.Doc
					if fcomment == nil {
						continue
					}
					fcommentList := fcomment.List
					for _, fcl := range fcommentList {
						fclcText := fcl.Text
						annotaionRaws := GetAnnotationBody(fclcText)
						for _, anno := range annotaionRaws {
							name := GetAnnotationName(anno)
							payload := GetAnnotationPayload(anno)
							annotation := definition.NewAnnotationField(name, payload)
							c.pushHolder(&AnnotationHolder{
								Annotation: annotation,
								RelateNode: f,
								RootNode:   gen,
							})
						}
					}

				}
			}
		}
	}

}
