package html

import (
	"strings"
	"testing"
)

const h = `<!DOCTYPE html>

<head>
    <title>Document</title>
    <style>
        zznq::after
    </style>
</head>

<body>
    <tag-zznq>
        in-html test zznq
    </tag-zznq>

    <div>
        in-html test zznq-div
    </div>

    <script>
        console.log("zznq");
    </script>

    <h1 style="style-tzznq: 12px;">

    </h1>

    <!-- comment-zznq -->

</body>

</html>
`

func TestMarkLocation(t *testing.T) {
	t.Log(MarkReflexLocation("zznq", strings.NewReader(h)))
}
