package notify

import "testing"

func TestNotifiers(t *testing.T) {
	n, err := NewNotifierWithYaml("notifier.yaml")
	t.Log(err)
	t.Log(n.Notify("http://localhost:8080/", `- type dom`))
}
