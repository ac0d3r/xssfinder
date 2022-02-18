package hookconv

import (
	"bytes"
	"io"
	"testing"

	"github.com/tdewolff/parse/v2"
	"github.com/tdewolff/parse/v2/js"
)

func TestJsParser(t *testing.T) {
	code := `
	// a=={"a":"1.2", b:2.2};
	a.b;
	a[b];
	a["b"];
	// a += "1";
	b = a + "1";
	var c = a + b;
	const name1 = "zznQ"; 
	let name;
	var b = decodeURI(location.hash.split("#")[1]);
	document.write("Hello2 " + b + "!");`

	ast, err := js.Parse(parse.NewInputString(code), js.Options{})

	if err != nil && err != io.EOF {
		t.Fatal(err)
	}
	for i := range ast.List {
		t.Logf("%#v\n", ast.List[i])
		t.Logf(ast.List[i].String())
	}
}

// testing walker https://github.com/tdewolff/parse/blob/master/js/walk_test.go

type walker struct{}

func (w *walker) Enter(n js.INode) js.IVisitor {
	switch n := n.(type) {
	case *js.Var:
		if bytes.Equal(n.Data, []byte("x")) {
			n.Data = []byte("obj")
		}
	}

	return w
}

func (w *walker) Exit(n js.INode) {}

func TestWalk(t *testing.T) {
	code := `
	if (true) {
		for (i = 0; i < 1; i++) {
			x.y = i
		}
	}`

	ast, err := js.Parse(parse.NewInputString(code), js.Options{})
	if err != nil {
		t.Fatal(err)
	}

	js.Walk(&walker{}, ast)

	t.Log(ast.JS())
}

func TestConvExpression(t *testing.T) {
	code := `
	var a;
        a += "1";
        var b = decodeURI(location.hash.split("#")[1]);
        document.write("Hello2 " + b + "!");

        typeof name;
        typeof "123";
        typeof("123");
        if (typeof "123" === "string"){
            console.log("123")
        }
        if (typeof a === "string"){
            console.log("bbbbb");
        }
        var bb = new String("123123");
        console.log(bb);
        
        function hello(){
            typeof name;
            typeof "123";
            typeof("123");
        }

        const sum = new Function('a', 'b', 'return a + b');
	`

	t.Log(HookConv(code))
}
