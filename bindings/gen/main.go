package main

import (
	"flag"
	"github.com/panyam/relay/bindings"
	"github.com/panyam/relay/bindings/rest"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
)

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

	parsedFile := parsedFiles[flag.Args()[0]]
	typeSystem := bindings.NewTypeSystem()
	bindings.ParseFile(parsedFile, typeSystem)

	generator := rest.NewGenerator(nil, typeSystem, "../rest/templates/")

	generator.EmitClientClass(parsedFile.Name.Name, serviceName)

	// Generate code for each of the service methods
	for _, field := range generator.ServiceType.Fields {
		switch optype := field.Type.TypeData.(type) {
		case *bindings.FunctionTypeData:
			generator.EmitSendRequestMethod(os.Stdout, field.Name, optype, "arg")
		}
	}
}
