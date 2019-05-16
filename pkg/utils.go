package pkg

import (
	"fmt"
	gotypes "go/types"

	"github.com/dave/jennifer/jen"
	"github.com/vetcher/go-astra/types"
	"golang.org/x/tools/go/packages"
)

// FindInterface is a function which returns an interface struct from given file by given name.
// Result of functions is nil if interface with given name is not found in given file structure.
func FindInterface(f *types.File, name string) *types.Interface {
	for _, iface := range f.Interfaces {
		if iface.Name == name {
			return &iface
		}
	}

	return nil
}

// Constructs jennifer statement for given astra type.
func TypeQual(currentPackage string, targetPackage string, t types.Type) *jen.Statement {
	stringRepr := t.String()

	ti := types.TypeImport(t)

	stmt := &jen.Statement{}

	tn := types.TypeName(t)
	if tn == nil {
		return nil
	}

	var op string
	opStartsFrom := len(stringRepr) - len(*tn)
	if ti != nil {
		opStartsFrom -= len(ti.Name) + 1
	}

	if opStartsFrom >= 0 {
		op = stringRepr[:opStartsFrom]
	}

	if op != "" {
		stmt = stmt.Op(op)
	}

	var pkg string
	if ti == nil {
		if !types.IsBuiltinTypeString(*tn) && currentPackage != targetPackage {
			pkg = currentPackage
		}
	} else {
		pkg = ti.Package
	}

	if pkg != "" {
		return stmt.Qual(pkg, *tn)
	}

	return stmt.Id(*tn)
}

// Checks if a given type is a pointer.
func IsPointerType(t types.Type) bool {
	_, ok := t.(types.TPointer)
	return ok
}

// Checks if a given type is a pointer.
func IsChanType(t types.Type) bool {
	_, ok := t.(types.TChan)
	return ok
}

func isNillable(t gotypes.Type) bool {
	switch t := t.(type) {
	case *gotypes.Pointer, *gotypes.Array, *gotypes.Map, *gotypes.Interface, *gotypes.Signature, *gotypes.Chan, *gotypes.Slice:
		return true
	case *gotypes.Named:
		return isNillable(t.Underlying())
	}
	return false
}

// Checks if a given type may be nil.
func IsNillableType(t types.Type) bool {
	// Check if type is imported.
	// If it is imported, we need to know if it is interface.
	ti := types.TypeImport(t)
	if ti != nil {
		xType, err := packages.Load(&packages.Config{
			Mode: packages.LoadTypes,
		}, ti.Package)
		if err == nil {
			tn := types.TypeName(t)

			if tn != nil && len(xType) > 0 {
				loadedType := xType[0].Types.Scope().Lookup(*tn).Type()
				return isNillable(loadedType)
			}
		}
	}

	// Check by
	return IsPointerType(t) || IsChanType(t) || types.IsInterface(t) || types.IsArray(t) || types.IsMap(t)
}

// Checks if a given type is a standard error type.
func IsErrorType(t types.Type) bool {
	return t.String() == "error"
}

// Adds default autogenerated package comment for a file.
func AddDefaultPackageComment(f *jen.File, generatorName string, version string) *jen.File {
	f.PackageComment(fmt.Sprintf("File generated by gen. DO NOT EDIT.\nGenerator plugin %s %s.", generatorName, version))
	return f
}
