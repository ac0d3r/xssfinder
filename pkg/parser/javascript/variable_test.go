package javascript

import "testing"

func TestGetAllVariable(t *testing.T) {
	t.Log(GetAllVariable(`
	var a = 1;
	var aa,bb,cc;
	let b = '1';
	let d = taest();
	
	const c = '';
	console.log(a,b,c);
	`))
}
