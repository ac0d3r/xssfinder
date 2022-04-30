package html

import (
	"strings"
	"testing"
)

func TestGetParams(t *testing.T) {
	doc := `
	<input name="xxsfinder"></input>
	<input></input>
	<script></script>
	
	<script> 
		const ccca = '';
		console.log(ccca);
		
		var a = 1;
		var aa,bb,cc;
		let b = '1';
		let d = taest();
		
		const c = '';
		console.log(a,b,c);
	</script>
	`
	t.Log(GetParams(strings.NewReader(doc)))
}
