package main

import (
	"fmt"
	"html/template"
	"strings"
)

func generateHelpers(d APIDescription) error {
	helpers := strings.Builder{}
	helpers.WriteString(`
// THIS FILE IS AUTOGENERATED. DO NOT EDIT.
// Regen by running 'go generate' in the repo root.

package gen

`)

	for _, tgMethodName := range orderedMethods(d) {
		tgMethod := d.Methods[tgMethodName]

		helper, err := generateHelperDef(d, tgMethod)
		if err != nil {
			return fmt.Errorf("failed to generate helpersfor %s: %w", tgMethodName, err)
		}
		if helper == "" {
			continue
		}
		helpers.WriteString(helper)
	}

	return writeGenToFile(helpers, "gen/gen_helpers.go")
}

func generateHelperDef(d APIDescription, tgMethod MethodDescription) (string, error) {
	helperDef := strings.Builder{}

	hasFromChat := false
	for _, x := range tgMethod.Fields {
		if x.Name == "from_chat_id" {
			hasFromChat = true
			break
		}
	}

	for _, typeName := range orderedTgTypes(d) {
		tgType := d.Types[typeName]

		repl := strings.ReplaceAll(tgMethod.Name, typeName, "")
		if repl == tgMethod.Name {
			continue
		}

		fields := map[string]string{}
		for _, f := range tgMethod.Fields {
			if f.Name == titleToSnake(typeName)+"_id" || f.Name == "id" {
				idField := "id"
				if typeName == tgTypeMessage {
					idField = "message_id"
				} else if typeName == tgTypeFile {
					idField = "file_id"
				}
				fields[titleToSnake(typeName)+"_id"] = idField
			}
		}

		if len(fields) == 0 {
			continue
		}

		for _, f := range tgType.Fields {
			if f.Name == "reply_to_message" {
				// this subfield just causes confusion; we always want the message_id
				continue
			}

			for _, mf := range tgMethod.Fields {
				prefType, err := f.getPreferredType()
				if err != nil {
					return "", fmt.Errorf("failed to get preferred type for field %s of %s: %w", mf.Name, tgMethod.Name, err)
				}
				if isTgType(d, prefType) && f.Name+"_id" == mf.Name {
					repl = strings.ReplaceAll(repl, prefType, "")
					if hasFromChat && mf.Name == "chat_id" {
						fields["from_chat_id"] = f.Name + ".Id"
					} else {
						fields[mf.Name] = f.Name + ".Id" // Note: maybe not just assume ID field exists?
					}
				}
			}
		}

		repl = strings.Title(repl)
		ret, err := tgMethod.GetReturnType(d)
		if err != nil {
			return "", fmt.Errorf("failed to get return type for %s: %w", tgMethod.Name, err)
		}

		optsContent := strings.Builder{}
		var funcCallArgList []string
		funcDefArgList := []string{"b Bot"}
		hasOpts := false
		for _, mf := range tgMethod.Fields {
			hasOpts = hasOpts || !mf.Required
			prefType, err := mf.getPreferredType()
			if err != nil {
				return "", fmt.Errorf("failed to get preferred type for field %s of %s: %w", mf.Name, tgMethod.Name, err)
			}

			if fName, ok := fields[mf.Name]; ok {
				if !mf.Required {
					def := getDefaultReturnVal(prefType)
					optsContent.WriteString("\n	if opts." + snakeToTitle(mf.Name) + " == " + def + " {")
					optsContent.WriteString("\n		opts." + snakeToTitle(mf.Name) + " = v." + snakeToTitle(fName))
					optsContent.WriteString("\n	}")

					continue
				}
				funcCallArgList = append(funcCallArgList, "v."+snakeToTitle(fName))
				continue
			}

			if !mf.Required {
				continue
			}

			funcDefArgList = append(funcDefArgList, snakeToCamel(mf.Name)+" "+prefType)
			funcCallArgList = append(funcCallArgList, snakeToCamel(mf.Name))

		}

		if hasOpts {
			funcDefArgList = append(funcDefArgList, "opts "+tgMethod.optsName())
			funcCallArgList = append(funcCallArgList, "opts")
		}

		funcDefArgs := strings.Join(funcDefArgList, ", ")
		funcCallArgs := strings.Join(funcCallArgList, ", ")

		helperDef.WriteString("\n// Helper method for Bot." + strings.Title(tgMethod.Name))
		err = helperFuncTmpl.Execute(&helperDef, helperFuncData{
			TypeName:     typeName,
			HelperName:   repl,
			ReturnType:   ret,
			FuncDefArgs:  funcDefArgs,
			Contents:     optsContent.String(),
			MethodName:   strings.Title(tgMethod.Name),
			FuncCallArgs: funcCallArgs,
		})
		if err != nil {
			return "", fmt.Errorf("failed to execute template to generate %s helper method on %s: %w", repl, typeName, err)
		}
	}

	return helperDef.String(), nil
}

var helperFuncTmpl = template.Must(template.New("helperFunc").Parse(helperFunc))

type helperFuncData struct {
	TypeName     string
	HelperName   string
	ReturnType   string
	FuncDefArgs  string
	Contents     string
	MethodName   string
	FuncCallArgs string
}

const helperFunc = `
func (v {{.TypeName}}) {{.HelperName}}({{.FuncDefArgs}}) ({{.ReturnType}}, error) {
	{{- .Contents}}
	return b.{{.MethodName}}({{.FuncCallArgs}})
}
`
