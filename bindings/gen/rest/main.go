package main

import (
	"flag"
	"fmt"
	"github.com/panyam/relay/bindings/gen"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"text/template"
)

func Hello(a int, b int, c int) {
}

func ArgListMaker(paramTypes []*gen.Type, withNames bool, isInput bool) string {
	numParams := len(paramTypes)
	out := ""
	if !isInput && numParams > 1 {
		out = "("
	}
	for index, param := range paramTypes {
		if index > 0 {
			out += ", "
		}
		if withNames {
			out += fmt.Sprintf("arg%d ", index)
		}
		out += param.Signature()
	}
	if !isInput && numParams > 1 {
		out += ")"
	}
	return out
}

type HttpTemplateParams struct {
	Package           string
	ClientPackageName string
	ServiceName       string
	ClientPrefix      string
	ClientSuffix      string
	ArgListMaker      func([]*gen.Type, bool, bool) string
	ServiceType       *gen.RecordTypeData
}

func NewHttpTemplateParams(packageName string, serviceName string, typeSystem gen.ITypeSystem) *HttpTemplateParams {
	out := &HttpTemplateParams{
		ClientPackageName: "client",
		ClientSuffix:      "Client",
		ArgListMaker:      ArgListMaker,
		ServiceName:       serviceName,
	}
	log.Println("File Package: ", packageName)
	out.ServiceType = typeSystem.GetType(packageName, serviceName).TypeData.(*gen.RecordTypeData)
	return out
}

func (h *HttpTemplateParams) ClientName() string {
	return h.ClientPrefix + h.ServiceName + h.ClientSuffix
}

func main() {
	var serviceName, operation string
	flag.StringVar(&serviceName, "service", "", "The service whose methods are to be extracted and for whome binding code is to be generated")
	flag.StringVar(&operation, "operation", "", "The operation within the service to be generated code for.  If this is empty or not provided then ALL operations in the service will code generated for them")

	flag.Parse()

	log.Println("Args: ", flag.Args())
	if serviceName == "" {
		log.Println("Service required")
	}

	parsedFiles := make(map[string]*ast.File)
	for _, srcFile := range flag.Args() {
		fset := token.NewFileSet() // positions are relative to fset
		parsedFile, err := parser.ParseFile(fset, srcFile, nil, parser.ParseComments|parser.AllErrors)
		if err != nil {
			log.Println("Parsing error: ", err)
			return
		}
		parsedFiles[srcFile] = parsedFile
	}

	/*
		fset := token.NewFileSet() // positions are relative to fset
		pkg, err := ast.NewPackage(fset, parsedFiles, types.GcImporter, types.Universe)
		if err != nil {
			log.Println("Package creation err: ", err)
			return
		}
		parsedFile := ast.MergePackageFiles(pkg, 0)
	*/

	parsedFile := parsedFiles[flag.Args()[0]]
	log.Println("Import Spec: ", parsedFile.Imports)
	log.Println("Unresolved: ", parsedFile.Unresolved)
	typeSystem := gen.NewTypeSystem()
	gen.ParseFile(parsedFile, typeSystem)

	// now take the service and generate it
	params := NewHttpTemplateParams(parsedFile.Name.Name, serviceName, typeSystem)

	// GenerateOperationMethods(params)

	// write the request/response serializers and deserializers
	GenerateOperationIOMethods(params)
}

func GenerateOperationMethods(params *HttpTemplateParams) {
	tmpl, err := template.New("service.gen").ParseFiles("templates/service.gen")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, params)
	if err != nil {
		panic(err)
	}
}

func GenerateOperationIOMethods(params *HttpTemplateParams) {
	tmpl, err := template.New("io.gen").ParseFiles("templates/io.gen")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, params)
	if err != nil {
		panic(err)
	}
}
