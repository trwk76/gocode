package stdhttp

import (
	"net/http"
	"reflect"

	code "github.com/trwk76/gocode"
	g "github.com/trwk76/gocode/go"
	"github.com/trwk76/gocode/web/api"
)

func NewGenerator(mapUnit *g.Unit, modelUnit *g.Unit, opIDXform code.IDTransformer, opPath OperationPathFunc, opWrapper OperationWrapFunc, typeConv TypeConverter) Generator {
	if opIDXform == nil {
		opIDXform = func(id string) string { return id }
	}

	if opPath == nil {
		opPath = defaultOpPath
	}

	if opWrapper == nil {
		opWrapper = defaultOperationWrapper
	}

	if typeConv == nil {
		typeConv = (*DefaultTypeConverter)(nil)
	}

	return Generator{
		mapUnit:   mapUnit,
		mdlUnit:   modelUnit,
		opIDXform: opIDXform,
		opPath:    opPath,
		opWrap:    opWrapper,
		tcnv:      reflect.TypeOf(typeConv).Elem(),
	}
}

func (gen *Generator) Initialize(baseURL string) {
	gen.baseURL = baseURL
}

func (gen *Generator) Finalize() {
	mux := g.SymbolFor[http.ServeMux](gen.mapUnit)

	gen.mapUnit.Decls = append(
		gen.mapUnit.Decls,
		g.FuncDecls{
			g.FuncDecl{
				ID: g.ID("Map"),
				Params: g.Params{
					{
						ID:   g.ID("m"),
						Type: g.PtrType{Item: mux},
					},
				},
				Body: gen.mapStmts,
			},
		},
	)
}

type (
	Generator struct {
		baseURL   string
		mapUnit   *g.Unit
		mapStmts  g.BlockStmt
		mdlUnit   *g.Unit
		mdlTypes  g.TypeDecls
		opIDXform code.IDTransformer
		opPath    OperationPathFunc
		opWrap    OperationWrapFunc
		tcnv      reflect.Type
	}
)

var (
	varMux         g.Symbol = g.Symbol{ID: g.ID("mux")}
	funcHandleFunc g.ID     = g.ID("HandleFunc")
)

var (
	_ api.Generator = (*Generator)(nil)
)
