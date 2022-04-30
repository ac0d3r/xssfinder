package javascript

import (
	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
)

type varWalker struct {
	vars []string
}

var _ js.IVisitor = &varWalker{}

func newVarWalker() *varWalker {
	return &varWalker{
		vars: make([]string, 0),
	}
}

func (w *varWalker) Enter(n js.INode) js.IVisitor {
	switch n := n.(type) {
	case *js.VarDecl:
		for i := range n.List {
			v, ok := n.List[i].Binding.(*js.Var)
			if !ok {
				continue
			}
			w.vars = append(w.vars, string(v.Data))
		}
	}

	return w
}

func (w *varWalker) Exit(n js.INode) {}

func GetAllVariable(code string) ([]string, error) {
	ast, err := js.Parse(parse.NewInputString(code), js.Options{})
	if err != nil {
		return nil, err
	}

	walker := &varWalker{}
	js.Walk(walker, ast)

	return walker.vars, nil
}
