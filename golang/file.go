package golang

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	code "github.com/trwk76/gocode"
)

type (
	File struct {
		doc   Comment
		pkg   PackageDecl
		imp   ImportDecls
		decls Decls
	}
)

func NewFile(doc Comment, pkgName string) File {
	return File{
		doc: doc,
		pkg: PackageDecl(pkgName),
	}
}

func (f *File) Add(decl Decl) {
	if decl == nil {
		return
	}

	f.decls = append(f.decls, decl)
}

func (f *File) NamedType(id ID, pkg string, args GenArgs) NamedType {
	return NamedType{
		ID:   id,
		Pkg:  f.pkgAlias(pkg),
		Args: args,
	}
}

func (f *File) TypeOf(v any) Type {
	return f.Type(reflect.TypeOf(v))
}

func (f *File) Type(t reflect.Type) Type {
	if t.Name() != "" {
		return f.NamedType(ID(t.Name()), t.PkgPath(), nil)
	}

	switch t.Kind() {
	case reflect.Array:
		return SliceType{
			Len:  IntExpr(t.Len()),
			Item: f.Type(t.Elem()),
		}
	case reflect.Bool:
		return Bool
	case reflect.Float32:
		return Float32
	case reflect.Float64:
		return Float64
	case reflect.Int:
		return Int
	case reflect.Int16:
		return Int16
	case reflect.Int32:
		return Int32
	case reflect.Int64:
		return Int64
	case reflect.Int8:
		return Int8
	case reflect.Map:
		return MapType{
			Key: f.Type(t.Key()),
			Val: f.Type(t.Elem()),
		}
	case reflect.Pointer:
		return PtrType{
			Target: f.Type(t.Elem()),
		}
	case reflect.Slice:
		return SliceType{
			Item: f.Type(t.Elem()),
		}
	case reflect.String:
		return String
	case reflect.Uint:
		return Uint
	case reflect.Uint16:
		return Uint16
	case reflect.Uint32:
		return Uint32
	case reflect.Uint64:
		return Uint64
	case reflect.Uint8:
		return Uint8
	case reflect.Func:
	case reflect.Interface:
	case reflect.Struct:
		flds := make(Fields, 0)

		for i := 0; i < t.NumField(); i++ {
			fld := t.Field(i)

			flds = append(flds, Field{
				ID:   ID(fld.Name),
				Type: f.Type(fld.Type),
				Tags: parseTags(string(fld.Tag)),
			})
		}

		return StructType{Fields: flds}
	}

	panic(fmt.Errorf("type %v is not supported", t))
}

func (f *File) Symbol(id ID, pkg string) SymbolExpr {
	return SymbolExpr{
		ID:  id,
		Pkg: f.pkgAlias(pkg),
	}
}

func (f *File) Render(w *code.Writer) {
	decls := Decls{Comment("THIS FILE WAS AUTOMATICALLY GENERATED; DO NOT EDIT")}

	if len(f.doc) > 0 {
		decls = append(decls, f.doc)
	}

	decls = append(decls, f.pkg)

	if len(f.imp) > 0 {
		decls = append(decls, f.imp)
	}

	decls = append(decls, f.decls...)

	for idx, itm := range decls {
		if idx > 0 {
			w.Newline()
		}

		itm.write(w)
		w.Newline()
	}
}

func (f *File) pkgAlias(path string) string {
	if path == "" {
		return ""
	}

	for _, itm := range f.imp {
		if itm.Path == path {
			return itm.Alias
		}
	}

	alias := fmt.Sprintf("pkg%d", len(f.imp)+1)

	f.imp = append(f.imp, ImportDecl{
		Path:  path,
		Alias: alias,
	})

	return alias
}

func parseTags(str string) Tags {
	if str == "" {
		return nil
	}

	items := strings.Split(str, " ")
	tags := make(Tags, 0)

	for _, itm := range items {
		pos := strings.IndexByte(itm, ':')
		if pos > 0 {
			name := itm[:pos]
			val, err := strconv.Unquote(itm[pos+1:])
			if err == nil {
				tags = append(tags, Tag{
					Name: name,
					Val:  val,
				})
			}
		}
	}

	return tags
}
