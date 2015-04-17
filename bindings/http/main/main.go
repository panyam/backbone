package main

import (
	"flag"
	"github.com/panyam/relay/bindings/gen"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"reflect"
)

func NodeToType(node ast.Node) *gen.Type {
	switch typeExpr := node.(type) {
	case *ast.StarExpr:
		// we have a reference type
		out := &gen.Type{TypeClass: gen.ReferenceType}
		out.TypeData = gen.ReferenceTypeData{TargetType: NodeToType(typeExpr.X)}
		return out
	case *ast.FuncType:
		{
			out := &gen.Type{TypeClass: gen.FunctionType}

			// create a function type
			functionType := gen.FunctionTypeData{}
			out.TypeData = functionType
			for _, param := range typeExpr.Params.List {
				paramType := NodeToType(param.Type)
				functionType.InputTypes = append(functionType.InputTypes, paramType)
			}
			if typeExpr.Results != nil && typeExpr.Results.List != nil {
				for _, result := range typeExpr.Results.List {
					resultType := NodeToType(result.Type)
					functionType.OutputTypes = append(functionType.OutputTypes, resultType)
				}
			}
			return out
		}
	case *ast.ArrayType:
		return &gen.Type{TypeClass: gen.ArrayType, TypeData: gen.ArrayTypeData{TargetType: NodeToType(typeExpr.Elt)}}
	case *ast.Ident:
		return &gen.Type{TypeClass: gen.LazyType, TypeData: typeExpr.Name}
	case *ast.GenDecl:
		typeSpec := typeExpr.Specs[0].(*ast.TypeSpec)
		out := &gen.Type{TypeClass: gen.RecordType}
		out.Name = typeSpec.Name.Name
		recordData := &gen.RecordTypeData{}
		out.TypeData = recordData

		switch typeExpr := typeSpec.Type.(type) {
		case *ast.InterfaceType:
			{
				log.Println("TypeSpec: ", typeExpr.Methods)
				fieldList := typeExpr.Methods.List
				for index, field := range fieldList {
					log.Println("Processing method: ", index, field.Names[0], field.Type, reflect.TypeOf(field.Type))
					fieldType := NodeToType(field.Type)
					for _, fieldName := range field.Names {
						recordData.FieldNames = append(recordData.FieldNames, fieldName.Name)
						recordData.FieldTypes = append(recordData.FieldTypes, fieldType)
					}
				}
			}
		case *ast.StructType:
			{
			}
		}
		return out
	}
	log.Println("Damn - the wrong type: ", node, reflect.TypeOf(node))
	return nil
}

func FindDecl(parsedFile *ast.File, declName string, declType string) *ast.GenDecl {
	for _, decl := range parsedFile.Decls {
		gendecl, ok := decl.(*ast.GenDecl)
		if ok && len(gendecl.Specs) > 0 {
			typeSpec := gendecl.Specs[0].(*ast.TypeSpec)
			if typeSpec.Name.Name == declName {
				if declType == "" || declType == "asdf" {
					return gendecl
				}
			}
		}
	}
	return nil
}

func main() {
	var srcFilename, serviceName, operation string
	flag.StringVar(&srcFilename, "src", "", "Input file where the service interface resides.")
	flag.StringVar(&serviceName, "service", "", "The service whose methods are to be extracted and for whome binding code is to be generated")
	flag.StringVar(&operation, "operation", "", "The operation within the service to be generated code for.  If this is empty or not provided then ALL operations in the service will code generated for them")

	flag.Parse()

	if srcFilename == "" {
		log.Println("Filename required")
		return
	}

	if serviceName == "" {
		log.Println("Service required")
	}

	fset := token.NewFileSet() // positions are relative to fset
	parsedFile, err := parser.ParseFile(fset, srcFilename, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		log.Println("Parsing error: ", err)
		return
	}

	serviceDecl := FindDecl(parsedFile, serviceName, "")
	nt := NodeToType(serviceDecl)
	log.Println("ServiceType: ", nt)
}
