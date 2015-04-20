package bindings

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
			if ok && typeSpec.Name.Name == declName {
				return gendecl
			}
		}
	}
	return nil
}

/**
 * Parses a file and returns a map of types indexed by name.
 */
func ParseFile(parsedFile *ast.File, typeSystem ITypeSystem) error {
	for _, decl := range parsedFile.Decls {
		gendecl, ok := decl.(*ast.GenDecl)
		if ok && len(gendecl.Specs) > 0 {
			typeSpec, ok := gendecl.Specs[0].(*ast.TypeSpec)
			if ok {
				NodeToType(typeSpec, parsedFile.Name.Name, typeSystem)
			}
		}
	}
	return nil
}

/**
 * Convert a node to a type.
 */
func NodeToType(node ast.Node, pkg string, typeSystem ITypeSystem) *Type {
	switch typeExpr := node.(type) {
	case *ast.StarExpr:
		// we have a reference type
		out := &Type{TypeClass: ReferenceType}
		out.TypeData = &ReferenceTypeData{TargetType: NodeToType(typeExpr.X, pkg, typeSystem)}
		return out
	case *ast.FuncType:
		{
			out := &Type{TypeClass: FunctionType}

			// create a function type
			functionType := &FunctionTypeData{}
			out.TypeData = functionType
			for _, param := range typeExpr.Params.List {
				paramType := NodeToType(param.Type, pkg, typeSystem)
				functionType.InputTypes = append(functionType.InputTypes, paramType)
			}
			if typeExpr.Results != nil && typeExpr.Results.List != nil {
				for _, result := range typeExpr.Results.List {
					resultType := NodeToType(result.Type, pkg, typeSystem)
					functionType.OutputTypes = append(functionType.OutputTypes, resultType)
				}
			}
			return out
		}
	case *ast.MapType:
		typeData := &MapTypeData{}
		typeData.KeyType = NodeToType(typeExpr.Key, pkg, typeSystem)
		typeData.ValueType = NodeToType(typeExpr.Value, pkg, typeSystem)
		return &Type{TypeClass: MapType, TypeData: typeData}
	case *ast.ArrayType:
		return &Type{TypeClass: ArrayType,
			TypeData: &ArrayTypeData{TargetType: NodeToType(typeExpr.Elt, pkg, typeSystem)}}
	case *ast.Ident:
		t := typeSystem.GetType(pkg, typeExpr.Name)
		if t == nil {
			t = &Type{TypeClass: LazyType, TypeData: typeExpr.Name}
			typeSystem.AddType(pkg, typeExpr.Name, t)
		}
		return t
	case *ast.SelectorExpr:
		pkgName := typeExpr.X.(*ast.Ident).Name
		t := typeSystem.GetType(pkgName, typeExpr.Sel.Name)
		if t == nil {
			t = &Type{TypeClass: LazyType, TypeData: typeExpr.Sel.Name}
			typeSystem.AddType(pkgName, typeExpr.Sel.Name, t)
		}
		return t
	case *ast.TypeSpec:
		out := &Type{}
		recordData := &RecordTypeData{Name: typeExpr.Name.Name}
		currT := typeSystem.GetType(pkg, recordData.Name)
		if currT == nil {
			typeSystem.AddType(pkg, recordData.Name, out)
		} else {
			out = currT
			if currT.TypeClass != LazyType {
				// what if it already exists?
				log.Println("ERROR: Redefinition of type: ", recordData.Name, currT.TypeClass)
			}
		}
		out.TypeClass = RecordType
		out.TypeData = recordData

		switch typeExpr := typeExpr.Type.(type) {
		case *ast.InterfaceType:
			{
				log.Println("TypeSpec: ", typeExpr.Methods)
				fieldList := typeExpr.Methods.List
				for index, field := range fieldList {
					log.Println("Processing method: ", index, field.Names[0], field.Type, reflect.TypeOf(field.Type))
					fieldType := NodeToType(field.Type, pkg, typeSystem)
					for _, fieldName := range field.Names {
						field := &Field{Name: fieldName.Name, Type: fieldType}
						recordData.Fields = append(recordData.Fields, field)
					}
				}
			}
		case *ast.StructType:
			{
				log.Println("Struct TypeSpec: ", typeExpr.Fields)
				fieldList := typeExpr.Fields.List
				for index, field := range fieldList {
					fieldType := NodeToType(field.Type, pkg, typeSystem)
					log.Println("Processing field: ", index, field.Names, field.Type, reflect.TypeOf(field.Type))
					if len(field.Names) == 0 {
						recordData.Bases = append(recordData.Bases, fieldType)
					} else {
						for _, fieldName := range field.Names {
							field := &Field{Name: fieldName.Name, Type: fieldType}
							recordData.Fields = append(recordData.Fields, field)
						}
					}
				}
			}
		}
		return out
	}
	log.Println("Damn - the wrong type: ", node, reflect.TypeOf(node))
	return nil
}
