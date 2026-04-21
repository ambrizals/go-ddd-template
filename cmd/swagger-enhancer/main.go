package main

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	data, err := os.ReadFile("docs/swagger.json")
	if err != nil {
		fmt.Println("Reading swagger.json:", err)
		os.Exit(1)
	}

	var swagger map[string]interface{}
	if err := json.Unmarshal(data, &swagger); err != nil {
		fmt.Println("Parsing swagger.json:", err)
		os.Exit(1)
	}

	definitions := swagger["definitions"].(map[string]interface{})
	enhanceSchemas(definitions)
	swagger["definitions"] = definitions

	output, err := os.Create("docs/swagger.json")
	if err != nil {
		fmt.Println("Writing swagger.json:", err)
		os.Exit(1)
	}
	defer output.Close()

	encoder := json.NewEncoder(output)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(swagger); err != nil {
		fmt.Println("Encoding swagger.json:", err)
		os.Exit(1)
	}

	jsonToYaml("docs/swagger.json", "docs/swagger.yaml")

	fmt.Println("swagger.json and swagger.yaml enhanced successfully")
}

func jsonToYaml(jsonPath, yamlPath string) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		fmt.Println("Reading JSON:", err)
		return
	}

	var v interface{}
	if err := json.Unmarshal(data, &v); err != nil {
		fmt.Println("Unmarshaling JSON:", err)
		return
	}

	var sb strings.Builder
	writeYAML(&sb, v, 0)
	os.WriteFile(yamlPath, []byte(sb.String()), 0644)
}

func writeYAML(sb *strings.Builder, v interface{}, indent int) {
	prefix := strings.Repeat("  ", indent)

	switch val := v.(type) {
	case map[string]interface{}:
		first := true
		for k, v := range val {
			if !first {
				sb.WriteString("\n")
			}
			first = false
			sb.WriteString(prefix)
			sb.WriteString(k)
			sb.WriteString(":")
			if v == nil {
				continue
			}
			switch v.(type) {
			case map[string]interface{}, []interface{}:
				sb.WriteString("\n")
				writeYAML(sb, v, indent+1)
			default:
				sb.WriteString(" ")
				writeYAML(sb, v, indent+1)
			}
		}
	case []interface{}:
		for i, item := range val {
			if i > 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(prefix)
			sb.WriteString("- ")
			switch item.(type) {
			case map[string]interface{}:
				writeYAML(sb, item, indent+1)
			default:
				writeYAML(sb, item, indent+1)
			}
		}
	case string:
		sb.WriteString(fmt.Sprintf("%q", val))
	case float64:
		sb.WriteString(fmt.Sprintf("%v", val))
	case bool:
		sb.WriteString(fmt.Sprintf("%v", val))
	case nil:
		sb.WriteString("null")
	default:
		sb.WriteString(fmt.Sprintf("%v", val))
	}
}

func enhanceSchemas(definitions map[string]interface{}) {
	schemasToCheck := []struct {
		name       string
		goStruct   string
		packageDir string
	}{
		{"dto.UserResponseDTO", "UserResponse", "user"},
		{"dto.DeactivateUserResponseDTO", "DeactivateUserResponseData", "user"},
		{"update_user.UpdateUserInput", "UpdateUserInput", "user"},
	}

	for _, schema := range schemasToCheck {
		if _, exists := definitions[schema.name]; !exists {
			continue
		}

		if isEmptyObject(definitions[schema.name]) {
			goSchema := extractFromGoFile(schema.goStruct, schema.packageDir)
			if goSchema != nil {
				definitions[schema.name] = goSchema
			}
		}
	}
}

func isEmptyObject(v interface{}) bool {
	m, ok := v.(map[string]interface{})
	if !ok {
		return false
	}
	props, hasProps := m["properties"]
	if !hasProps {
		return true
	}
	propsMap, ok := props.(map[string]interface{})
	return ok && len(propsMap) == 0
}

func extractFromGoFile(structName, packageDir string) map[string]interface{} {
	dtoDir := filepath.Join("internal", "modules", packageDir, "dto")

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dtoDir, nil, parser.ParseComments)
	if err != nil {
		fmt.Printf("Parsing %s: %v\n", dtoDir, err)
		return nil
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}

				for _, spec := range genDecl.Specs {
					typeSpec, ok := spec.(*ast.TypeSpec)
					if !ok || typeSpec.Name.Name != structName {
						continue
					}

					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						return nil
					}

					return buildSchemaFromStruct(structType)
				}
			}
		}
	}

	return nil
}

func buildSchemaFromStruct(structType *ast.StructType) map[string]interface{} {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": map[string]interface{}{},
	}

	required := []string{}

	for _, field := range structType.Fields.List {
		if len(field.Names) == 0 {
			continue
		}

		fieldName := field.Names[0].Name
		jsonTag := getJSONTag(field)
		if jsonTag == "" {
			jsonTag = strings.ToLower(fieldName[:1]) + fieldName[1:]
		}

		fieldType := getFieldType(field.Type)
		prop := map[string]interface{}{
			"type": fieldType,
		}

		if fieldType == "string" && isTimeField(fieldName) {
			prop["example"] = "2023-10-27T10:00:00Z"
		}

		if fieldType == "integer" && strings.HasSuffix(fieldName, "ID") {
			prop["example"] = 1
		}

		if fieldType == "boolean" && strings.Contains(strings.ToLower(fieldName), "activated") {
			prop["example"] = true
		}

		schema["properties"].(map[string]interface{})[jsonTag] = prop
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

func getJSONTag(field *ast.Field) string {
	if field.Tag != nil {
		tag := strings.Trim(field.Tag.Value, "`")
		
		start := strings.Index(tag, `json:"`)
		if start < 0 {
			return ""
		}
		
		start += 6 // len of `json:"`
		end := start
		for end < len(tag) && tag[end] != '"' {
			end++
		}
		
		result := tag[start:end]
		if commaIdx := strings.Index(result, ","); commaIdx > 0 {
			return result[:commaIdx]
		}
		return result
	}
	return ""
}

func getFieldType(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		switch t.Name {
		case "string":
			return "string"
		case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
			return "integer"
		case "float32", "float64":
			return "number"
		case "bool":
			return "boolean"
		default:
			return "string"
		}
	case *ast.ArrayType:
		return "array"
	case *ast.MapType:
		return "object"
	default:
		return "string"
	}
}

func isTimeField(fieldName string) bool {
	lowerName := strings.ToLower(fieldName)
	return strings.Contains(lowerName, "at") || strings.Contains(lowerName, "date") || strings.Contains(lowerName, "time")
}