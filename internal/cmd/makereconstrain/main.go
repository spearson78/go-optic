package main

import (
	"embed"
	"log"
	"os"
	"strings"
	"text/template"
	//"github.com/spearson78/go-optic"
)

//go:embed *.tmpl
var templateFS embed.FS

type TemplateData struct {
	PackageName           string
	Prefix                string
	RestrictedName        string
	AnyName               string
	RestrictedType        string
	AnyType               string
	Desc                  string
	TypeParamName         string
	ReconstrainRestricted string
	ReconstrainAny        string
	InRestrictedRET       string
	InRestrictedRW        string
	InRestrictedDIR       string
	InRestrictedERR       string
}

func (t *TemplateData) InOpticPrefixL() string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {
		if v == t.TypeParamName {
			sb.WriteString("CompositionTree[" + v + "L, " + v + "R]")
		} else {
			sb.WriteString(v)
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (t *TemplateData) InOptic3(order []string) string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {
		if v == t.TypeParamName {
			sb.WriteString("CompositionTree[CompositionTree[" + v + order[0] + ", " + v + order[1] + "], " + v + order[2] + "]")
		} else {
			sb.WriteString(v)
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (t *TemplateData) InOptic4(order []string) string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {
		if v == t.TypeParamName {
			sb.WriteString("CompositionTree[CompositionTree[" + v + order[0] + ", " + v + order[1] + "], CompositionTree[" + v + order[2] + ", " + v + order[3] + "]]")
		} else {
			sb.WriteString(v)
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (t *TemplateData) InConstraintsPrefixL() string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {
		if v == t.TypeParamName {
			sb.WriteString(v)
			sb.WriteString("L any")
			sb.WriteString(", ")
			sb.WriteString(v)
			sb.WriteString("R T")
			sb.WriteString(t.RestrictedType)
		} else {
			sb.WriteString(v)
			sb.WriteString(" any")
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	x := sb.String()
	log.Println(x)
	return x
}

func (t *TemplateData) InConstraints3() string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {
		if v == t.TypeParamName {
			sb.WriteString(v)
			sb.WriteString("1 any, ")
			sb.WriteString(v)
			sb.WriteString("2 any, ")
			sb.WriteString(v)
			sb.WriteString("3 any")
		} else {
			sb.WriteString(v)
			sb.WriteString(" any")
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	x := sb.String()
	log.Println(x)
	return x
}

func (t *TemplateData) InConstraints4() string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {
		if v == t.TypeParamName {
			sb.WriteString(v)
			sb.WriteString("1 any, ")
			sb.WriteString(v)
			sb.WriteString("2 any, ")
			sb.WriteString(v)
			sb.WriteString("3 any, ")
			sb.WriteString(v)
			sb.WriteString("4 any")
		} else {
			sb.WriteString(v)
			sb.WriteString(" any")
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	x := sb.String()
	log.Println(x)
	return x
}

func (t *TemplateData) ReconstrainPrefixL() string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {
		sb.WriteString(v)
		if v == t.TypeParamName {
			sb.WriteString("L")
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (t *TemplateData) InConstraintsPrefixR() string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {
		if v == t.TypeParamName {
			sb.WriteString(v)
			sb.WriteString("L T")
			sb.WriteString(t.RestrictedType)
			sb.WriteString(", ")
			sb.WriteString(v)
			sb.WriteString("R any")
		} else {
			sb.WriteString(v)
			sb.WriteString(" any")
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	x := sb.String()
	log.Println(x)
	return x
}

func (t *TemplateData) ReconstrainPrefixR() string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {
		sb.WriteString(v)
		if v == t.TypeParamName {
			sb.WriteString("R")
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (t *TemplateData) InConstraintsPrefixSwapL() string {
	return t.InConstraints2([]string{"L", "R"})
}

func (t *TemplateData) InConstraints2(order []string) string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {
		if v == t.TypeParamName {
			sb.WriteString(v)
			sb.WriteString(order[0])
			sb.WriteString(" any")
			sb.WriteString(", ")
			sb.WriteString(v)
			sb.WriteString(order[1])
			sb.WriteString(" any")
		} else {
			sb.WriteString(v)
			sb.WriteString(" any")
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	x := sb.String()
	log.Println(x)
	return x
}

func (t *TemplateData) ReconstrainPrefixSwap() string {
	return t.Reconstrain2([]string{"R", "L"})
}

func (t *TemplateData) Reconstrain2(order []string) string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {

		if v == t.TypeParamName {
			//CompositionTree[RERR, LERR]
			sb.WriteString("CompositionTree[")
			sb.WriteString(v)
			sb.WriteString(order[0])
			sb.WriteString(", ")
			sb.WriteString(v)
			sb.WriteString(order[1])
			sb.WriteString("]")

		} else {
			sb.WriteString(v)
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (t *TemplateData) Reconstrain4(order []string) string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {

		if v == t.TypeParamName {
			//CompositionTree[RERR, LERR]
			sb.WriteString("CompositionTree[CompositionTree[")
			sb.WriteString(v)
			sb.WriteString(order[0])
			sb.WriteString(", ")
			sb.WriteString(v)
			sb.WriteString(order[1])
			sb.WriteString("], CompositionTree[")
			sb.WriteString(v)
			sb.WriteString(order[2])
			sb.WriteString(", ")
			sb.WriteString(v)
			sb.WriteString(order[3])
			sb.WriteString("]]")

		} else {
			sb.WriteString(v)
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (t *TemplateData) Reconstrain3(order []string) string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {
		if v == t.TypeParamName {

			sb.WriteString("CompositionTree[CompositionTree[" + v + order[0] + ", " + v + order[1] + "], " + v + order[2] + "]")

		} else {
			sb.WriteString(v)
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (t *TemplateData) ReconstrainSwapL() string {
	return t.Reconstrain3([]string{"2", "1", "3"})
}

func (t *TemplateData) InOpticSwapL() string {
	return t.InOptic3([]string{"1", "2", "3"})
}

func (t *TemplateData) InOpticMerge() string {

	var sb strings.Builder

	for i, v := range []string{"RET", "RW", "DIR", "ERR"} {
		if v == t.TypeParamName {
			sb.WriteString("CompositionTree[" + v + ", " + v + "]")
		} else {
			sb.WriteString(v)
		}

		if i != 3 {
			sb.WriteString(", ")
		}
	}

	return sb.String()
}

func (t *TemplateData) InConstraintsMerge() string {
	return t.InConstraints2([]string{"1", "2"})
}

func (t *TemplateData) InOpticMergeL() string {
	return t.InOptic3([]string{"1", "1", "2"})
}

func (t *TemplateData) ReconstrainMergeL() string {
	return t.Reconstrain2([]string{"1", "2"})
}

func (t *TemplateData) InOpticTrans() string {
	return t.InOptic4([]string{"1", "2", "3", "4"})
}

func (t *TemplateData) InConstraintsTrans() string {
	return t.InConstraints4()
}

func (t *TemplateData) ReconstrainTrans() string {
	return t.Reconstrain4([]string{"1", "3", "2", "4"})
}

func (t *TemplateData) InConstraintsTransL() string {
	return t.InConstraints3()
}

func (t *TemplateData) InOpticTransL() string {
	return t.InOptic3([]string{"1", "2", "3"})
}

func (t *TemplateData) ReconstrainTransL() string {
	return t.Reconstrain3([]string{"1", "3", "2"})
}

func main() {

	tmpl, err := template.New("reconstrain.tmpl").ParseFS(templateFS, "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	fileName := "generated_test.go"
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	/*
		//RET
		templateData := TemplateData{
			PackageName:           "main_test",
			Prefix:                "Ret",
			RestrictedName:        "1",
			AnyName:               "M",
			RestrictedType:        "ReturnOne",
			AnyType:               "ReturnMany",
			Desc:                  "return",
			TypeParamName:         "RET",
			ReconstrainRestricted: "ReturnOne, RW, DIR, ERR",
			ReconstrainAny:        "ReturnMany, RW, DIR, ERR",
			InRestrictedRET:       "TReturnOne",
			InRestrictedRW:        "any",
			InRestrictedDIR:       "any",
			InRestrictedERR:       "any",
		}
	*/

	/*
		//DIR
		templateData := TemplateData{
			PackageName:           "main_test",
			Prefix:                "Dir",
			RestrictedName:        "Bd",
			AnyName:               "Ud",
			RestrictedType:        "BiDir",
			AnyType:               "UniDir",
			Desc:                  "direction",
			TypeParamName:         "DIR",
			ReconstrainRestricted: "RET, RW, BiDir, ERR",
			ReconstrainAny:        "RET, RW, UniDir, ERR",
			InRestrictedRET:       "TBiDir",
			InRestrictedRW:        "any",
			InRestrictedDIR:       "any",
			InRestrictedERR:       "any",
		}
	*/

	//ERR
	templateData := TemplateData{
		PackageName:           "main_test",
		Prefix:                "EErr",
		RestrictedName:        "Pure",
		AnyName:               "Err",
		RestrictedType:        "Pure",
		AnyType:               "Err",
		Desc:                  "error",
		TypeParamName:         "ERR",
		ReconstrainRestricted: "RET, RW, DIR, Pure",
		ReconstrainAny:        "RET, RW, DIR, Err",
		InRestrictedRET:       "TPure",
		InRestrictedRW:        "any",
		InRestrictedDIR:       "any",
		InRestrictedERR:       "any",
	}

	err = tmpl.Execute(file, &templateData)
	if err != nil {
		log.Fatal(err)
	}

}
