package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	//"github.com/spearson78/go-optic"
)

//go:embed *.tmpl
var templateFS embed.FS

type ColOp struct {
	BaseName       string
	DocName        string
	Parameters     string
	TypeParameters string
	IxParameters   string
	Constructor    string
	ConstructorP   string
	IConstructor   string
	IConstructorP  string
	ReadWrite      bool

	Indexed bool
	BiDir   bool
	Err     string
	Pure    bool
}

func (c ColOp) ErrRecon() string {
	if c.Pure {
		return "EPureR"
	} else {
		return "EErrMerge"
	}
}

func (c ColOp) OErr() string {
	if c.Pure {
		return "Pure"
	} else {
		return c.Err
	}
}

func (c ColOp) RW() string {
	if c.ReadWrite {
		return "ReadWrite"
	} else {
		return "ReadOnly"
	}
}

type TemplateData struct {
	ImportOptics           bool
	AdditionalImports      []string
	PackageName            string
	CollectionName         string
	Parameters             string
	CollectionConstructor  string
	collectionConstructorP string
	CollectionType         string
	CollectionTypeP        string
	TypeParameters         string
	TypeParametersP        string
	IndexType              string
	FocusType              string
	FocusTypeP             string
	ErrType                string
	Unordered              bool
	NoColOf                bool
	Operations             []ColOp
	EqPredicate            string
}

func (TemplateData) AsListPrefix(v string) string {
	if v == "" {
		return ""
	}

	return v + ", "
}

func (TemplateData) AsCombinedList(v string, p string) string {
	if v == "" {
		if p == "" {
			return ""
		} else {
			return p
		}
	} else {
		if p == "" {
			return v
		} else {
			return v + ", " + p
		}
	}
}

func (TemplateData) AsCombinedListT(v string, p string) string {
	if v == "" {
		if p == "" {
			return ""
		} else {
			return "[" + p + "]"
		}
	} else {
		if p == "" {
			return "[" + v + "]"
		} else {
			return "[" + v + ", " + p + "]"
		}
	}
}

func (TemplateData) AsList(v string) string {
	if v == "" {
		return ""
	}

	return "[" + v + "]"
}

func (t TemplateData) Poly() bool {
	return false
	//return t.TypeParametersP != ""
}

func (t TemplateData) CollectionConstructorP() string {
	if t.collectionConstructorP == "" {
		return t.CollectionConstructor
	} else {
		return t.collectionConstructorP
	}
}

func (c ColOp) Dir() string {
	if c.BiDir {
		return "BiDir"
	} else {
		return "UniDir"
	}
}

func (t TemplateData) ParameterNames() string {
	var params strings.Builder
	for i, paramDef := range strings.Split(t.Parameters, ",") {
		nameIndex := strings.IndexAny(paramDef, " \t")
		if nameIndex != -1 {

			if i != 0 {
				params.WriteString(", ")
			}

			paramName := paramDef[0:nameIndex]
			params.WriteString(paramName)
		}
	}

	return params.String()
}

func (t TemplateData) TypeParameterNames() string {
	var params strings.Builder
	for i, paramDef := range strings.Split(t.TypeParameters, ",") {
		nameIndex := strings.IndexAny(paramDef, " \t")
		if nameIndex != -1 {

			if i != 0 {
				params.WriteString(", ")
			}

			paramName := paramDef[0:nameIndex]
			params.WriteString(paramName)
		}
	}

	return params.String()
}

type multiFlag []string

func (i *multiFlag) String() string {
	return "multiFlag"
}

func (i *multiFlag) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var imports multiFlag

func main() {

	flag.Var(&imports, "import", "additional imports")
	eqPredicate := flag.String("eq", "", "custom EqPredicate for EqCol")
	unordered := flag.Bool("unordered", false, "suppress generation of operations that rely on an underling order in the collection.")
	nocolof := flag.Bool("nocolof", false, "suppress generation of ColOf functions.")
	flag.Parse()

	tmpl, err := template.New("colops.tmpl").ParseFS(templateFS, "*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	if flag.NArg() != 13 {
		log.Fatal("Expected 13 params : makecolops <OutputFile> <PackageName> <CollectionName> <TypeParameters> <TypeParametersP> <Parameters> <CollectionConstructor> <CollectionConstructorP> <IndexType> <CollectionType> <CollectionTypeP> <FocusType> <FocusTypeP>")
	}

	fileName := flag.Arg(0)
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	templateData := TemplateData{
		ImportOptics:           flag.Arg(1) != "optic",
		AdditionalImports:      imports,
		PackageName:            flag.Arg(1),
		CollectionName:         flag.Arg(2),
		TypeParameters:         flag.Arg(3),
		TypeParametersP:        flag.Arg(4),
		Parameters:             flag.Arg(5),
		CollectionConstructor:  flag.Arg(6),
		collectionConstructorP: flag.Arg(7),
		IndexType:              flag.Arg(8),
		CollectionType:         flag.Arg(9),
		CollectionTypeP:        flag.Arg(10),
		FocusType:              flag.Arg(11),
		FocusTypeP:             flag.Arg(12),
		ErrType:                "Pure",
		Unordered:              *unordered,
		NoColOf:                *nocolof,
		EqPredicate:            *eqPredicate,
	}

	templateData.Operations = []ColOp{
		{
			BaseName:       "Filtered",
			DocName:        "Filtered",
			TypeParameters: "ERR any",
			Err:            "ERR",
			Parameters:     fmt.Sprintf("pred Predicate[%v, ERR]", templateData.FocusType),
			IxParameters:   fmt.Sprintf("pred PredicateI[%v, %v, ERR]", templateData.IndexType, templateData.FocusType),
			Constructor:    fmt.Sprintf("FilteredCol[%v](pred)", templateData.IndexType),
			ConstructorP:   fmt.Sprintf("FilteredCol[%v](pred)", templateData.IndexType),
			IConstructor:   fmt.Sprintf("FilteredColI(pred,IxMatchComparable[%v]())", templateData.IndexType),
			IConstructorP:  "FilteredColI(pred)",

			Indexed:   true,
			BiDir:     false,
			ReadWrite: true,
		},
		{
			BaseName:       "Append",
			DocName:        "Append",
			TypeParameters: "ERR any",
			Err:            "ERR",
			Parameters:     fmt.Sprintf("toAppend Collection[%v, %v, ERR]", templateData.IndexType, templateData.FocusType),
			IxParameters:   fmt.Sprintf("toAppend Collection[%v, %v, ERR]", templateData.IndexType, templateData.FocusType),
			Constructor:    "AppendCol(toAppend)",
			ConstructorP:   "AppendCol(toAppend)",

			Indexed: false,
			BiDir:   false,
			Pure:    false,
		},
	}

	if !templateData.Unordered {
		templateData.Operations = append(templateData.Operations, ColOp{
			BaseName:       "Prepend",
			DocName:        "Prepend",
			TypeParameters: "ERR any",
			Err:            "ERR",
			Parameters:     fmt.Sprintf("toPrepend Collection[%v, %v, ERR]", templateData.IndexType, templateData.FocusType),
			IxParameters:   fmt.Sprintf("toPrepend Collection[%v, %v, ERR]", templateData.IndexType, templateData.FocusType),
			Constructor:    "PrependCol(toPrepend)",
			ConstructorP:   "PrependCol(toPrepend)",

			Indexed: false,
			BiDir:   false,
			Pure:    false,
		})

	}

	if !templateData.Unordered {
		templateData.Operations = append(templateData.Operations,
			ColOp{
				BaseName:     "Sub",
				DocName:      "SubCol",
				Err:          "Pure",
				Parameters:   "start int,length int",
				IxParameters: "start int,length int",
				Constructor:  fmt.Sprintf("SubCol[%v, %v](start,length)", templateData.IndexType, templateData.FocusType),
				ConstructorP: fmt.Sprintf("SubColP[%v, %v, %v](start,length)", templateData.IndexType, templateData.FocusType, templateData.FocusTypeP),

				Indexed:   false,
				BiDir:     false,
				ReadWrite: true,
			},
			ColOp{

				BaseName:     "Reversed",
				DocName:      "Reversed",
				Err:          "Pure",
				Parameters:   "",
				IxParameters: "",
				Constructor:  fmt.Sprintf("ReversedCol[%v, %v]()", templateData.IndexType, templateData.FocusType),
				ConstructorP: fmt.Sprintf("ReversedCol[%v, %v]()", templateData.IndexType, templateData.FocusType),

				Indexed:   false,
				BiDir:     true,
				ReadWrite: true,
			},
			ColOp{
				BaseName:       "Ordered",
				DocName:        "Ordered",
				TypeParameters: "ERR any",
				Err:            "ERR",
				Parameters:     fmt.Sprintf("orderBy OrderByPredicate[%v, ERR]", templateData.FocusType),
				IxParameters:   fmt.Sprintf("orderBy OrderByPredicateI[%v, %v, ERR]", templateData.IndexType, templateData.FocusType),
				Constructor:    fmt.Sprintf("OrderedCol[%v](orderBy)", templateData.IndexType),
				ConstructorP:   fmt.Sprintf("OrderedCol[%v](orderBy)", templateData.IndexType),
				IConstructor:   fmt.Sprintf("OrderedColI(orderBy,IxMatchComparable[%v]())", templateData.IndexType),
				IConstructorP:  fmt.Sprintf("OrderedColI(orderBy,IxMatchComparable[%v]())", templateData.IndexType),

				Indexed:   true,
				BiDir:     false,
				ReadWrite: true,
			},
		)
	}

	err = tmpl.Execute(file, templateData)
	if err != nil {
		log.Fatal(err)
	}

}
