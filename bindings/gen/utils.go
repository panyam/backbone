package gen

import (
	"go/ast"
	"log"
	"reflect"
)

/**
 * Finds a GenDecl node in a parsed file.
 */
func FindDecl(parsedFile *ast.File, declName string) *ast.GenDecl {
	for _, decl := range parsedFile.Decls {
		gendecl, ok := decl.(*ast.GenDecl)
		if ok && len(gendecl.Specs) > 0 {
			typeSpec := gendecl.Specs[0].(*ast.TypeSpec)
			if typeSpec.Name.Name == declName {
				return gendecl
			}
		}
	}
	return nil
}

/**
 * Convert a node to a type.
 */
func NodeToType(node ast.Node) *Type {
	switch typeExpr := node.(type) {
	case *ast.StarExpr:
		// we have a reference type
		out := &Type{TypeClass: ReferenceType}
		out.TypeData = ReferenceTypeData{TargetType: NodeToType(typeExpr.X)}
		return out
	case *ast.FuncType:
		{
			out := &Type{TypeClass: FunctionType}

			// create a function type
			functionType := FunctionTypeData{}
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
	case *ast.MapType:
		typeData := &MapTypeData{}
		typeData.KeyType = NodeToType(typeExpr.Key)
		typeData.ValueType = NodeToType(typeExpr.Value)
		return &Type{TypeClass: MapType, TypeData: typeData}
	case *ast.ArrayType:
		return &Type{TypeClass: ArrayType, TypeData: ArrayTypeData{TargetType: NodeToType(typeExpr.Elt)}}
	case *ast.Ident:
		return &Type{TypeClass: LazyType, TypeData: typeExpr.Name}
	case *ast.GenDecl:
		typeSpec := typeExpr.Specs[0].(*ast.TypeSpec)
		out := &Type{TypeClass: RecordType}
		recordData := &RecordTypeData{Name: typeSpec.Name.Name}
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
