package rest

import (
	"fmt"
	"github.com/panyam/relay/bindings"
	"os"
	"text/template"
)

/**
 * Responsible for generating the code for the client classes.
 */
type Generator struct {
	// where the templates are
	Bindings     map[string]*HttpBinding
	TypeSystem   bindings.TypeSystem
	TemplatesDir string

	// Parameters to determine Generated output
	Package           string
	ClientPackageName string
	ServiceName       string
	ClientPrefix      string
	ClientSuffix      string
	httpBindings      map[string]*HttpBinding
	ArgListMaker      func([]*bindings.Type, bool) string
	ServiceType       *bindings.RecordTypeData
	TransportRequest  string
	OpName            string
	OpType            *bindings.FunctionTypeData
}

func NewGenerator(bindings map[string]*HttpBinding, typeSys bindings.TypeSystem, templatesDir string) *Generator {
	if bindings == nil {
		bindings := make(map[string]*HttpBinding)
	}
	out := Generator{Bindings: bindings,
		TypeSystem:        typeSys,
		TemplatesDir:      templatesDir,
		ClientPackageName: "restclient",
		ClientSuffix:      "Client",
		TransportRequest:  "*http.Request",
		ArgListMaker:      ArgListMaker,
	}
	return &out
}

/**
 * Emits the class that acts as a client for the service.
 */
func (g *Generator) EmitClientClass(pkgName string, serviceName string) error {
	g.ServiceName = serviceName
	g.ServiceType = g.TypeSystem.GetType(pkgName, serviceName).TypeData.(*bindings.RecordTypeData)

	tmpl, err := template.New("client.gen").ParseFiles(g.TemplatesDir + "client.gen")
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(os.Stdout, g)
	if err != nil {
		panic(err)
	}
	return err
}

/**
 * For a given service operation, emits a method which:
 * 1. Has inputs the same as those of the underlying service operation,
 * 2. creates a transport level request
 * 3. Sends the transport level request
 * 4. Gets a response from the transport level and returns it
 */
func (g *Generator) EmitSendRequestMethod(opName string, opType *bindings.FunctionTypeData, argPrefix string) error {
	g.StartWritingMethod(opName, opType, "arg")
	if opType.NumInputs() > 0 {
		if opType.NumInputs() == 1 {
			g.EmitObjectWriterCall("arg0", opType.InputTypes[0])
		} else {
			g.StartWritingList()
			for index, param := range opType.InputTypes {
				g.StartWritingChild(index)
				g.EmitObjectWriterCall(fmt.Sprintf("arg%d", index), param)
				g.EndWritingChild(index)
			}
			g.EndWritingList()
		}
	}
	g.EndWritingMethod(opName, opType)
	return nil
}

/**
 * For a given service operation, emits a method:
 * 1. whose input is a http.Response object
 * 2. Which can be parsed into the output values as expected by the service
 * 	  operations's output signature
 */
func (g *Generator) EmitReadResponseMethod(opName string, opType *bindings.FunctionTypeData, argPrefix string) error {
	g.StartReadingMethod(opName, opType, "arg")
	if opType.NumOutputs() > 0 {
		if opType.NumOutputs() == 1 {
			g.EmitObjectReaderCall("arg0", opType.OutputTypes[0])
		} else {
			g.StartReadingList()
			for index, param := range opType.OutputTypes {
				g.StartReadingChild()
				g.EmitObjectReaderCall(fmt.Sprintf("arg%d", index), param)
				g.EndReadingChild()
			}
			g.EndReadingList()
		}
	}
	g.EndReadingMethod(opName, opType)
}

func ArgListMaker(paramTypes []*bindings.Type, withNames bool) string {
	out := ""
	for index, param := range paramTypes {
		if index > 0 {
			out += ", "
		}
		if withNames {
			out += fmt.Sprintf("arg%d ", index)
		}
		out += param.Signature()
	}
	return out
}

func (g *Generator) ClientName() string {
	return g.ClientPrefix + g.ServiceName + g.ClientSuffix
}
