package cert

import (
	"os"
	"testing"
)

func TestGenCA(t *testing.T) {
	cert, key, err := GenCA()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(cert))
	t.Log(string(key))

	f, err := os.Create("xssfinder.ca.cert")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f.Write(cert)

	f1, err := os.Create("xssfinder.ca.key")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	f1.Write(key)
}
