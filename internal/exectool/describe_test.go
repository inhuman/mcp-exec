package exectool

import (
	"strings"
	"testing"
)

func TestDescription_ListsEnvironment(t *testing.T) {
	d := Description(testCfg())
	for _, want := range []string{"Python", "PyYAML", "Jinja2", "no network", "timeout"} {
		if !strings.Contains(d, want) {
			t.Errorf("description missing %q: %s", want, d)
		}
	}
}
