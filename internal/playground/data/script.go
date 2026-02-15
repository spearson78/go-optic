package data

import (
	"reflect"

	goscript "github.com/spearson78/go-script"
)

func init() {

	pkg := goscript.Package{
		Name: "github.com/spearson78/go-optic/internal/playground/data",
		Exported: map[string]reflect.Value{
			"O":           reflect.ValueOf(&o{}),
			"NewBlogPost": reflect.ValueOf(NewBlogPost),
			"NewComment":  reflect.ValueOf(NewComment),
			"OBlogPostOf": reflect.ValueOf(OBlogPostOf[any, any, any, any, any, any, any]),
		},
		ExportedTypes: map[string]reflect.Type{
			"Comment":  reflect.TypeFor[Comment](),
			"BlogPost": reflect.TypeFor[BlogPost](),
			"Rating":   reflect.TypeFor[Rating](),
		},
	}

	goscript.RegisterPackage(&pkg)

}
