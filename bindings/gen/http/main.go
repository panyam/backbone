package main

import (
	"flag"
	"github.com/panyam/relay/bindings/gen"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"text/template"
)

type HttpTemplateParams struct {
	Package       string
	ServiceName   string
	HandlerPrefix string
	HandlerSuffix string
	ServiceType   *gen.RecordTypeData
}

func (h *HttpTemplateParams) HandlerName() string {
	return h.HandlerPrefix + h.ServiceName + h.HandlerSuffix
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
	packageName := parsedFile.Name.Name
	log.Println("File Package: ", parsedFile, packageName)
	log.Println("Import Spec: ", parsedFile.Imports)
	log.Println("Unresolved: ", parsedFile.Unresolved)
	typeSystem := gen.NewTypeSystem()
	gen.ParseFile(parsedFile, typeSystem)

	// now take the service and generate it
	params := &HttpTemplateParams{
		ServiceType:   typeSystem.GetType(packageName, serviceName).TypeData.(*gen.RecordTypeData),
		ServiceName:   serviceName,
		HandlerSuffix: "Handler",
	}
	tmpl, err := template.New("service.gen").ParseFiles("templates/service.gen")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, params)
	if err != nil {
		panic(err)
	}
}
