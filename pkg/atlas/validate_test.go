package atlas

import (
	"fmt"
	"os"
	"testing"
)

func TestValidateReport(t *testing.T) {
	a := load(t)
	issues := a.Validate()
	counts := map[string]int{}
	for _, i := range issues {
		counts[i.Severity+"/"+i.Kind]++
	}
	for _, i := range issues {
		if i.Severity == "error" {
			t.Errorf("data error: %s %q: %s", i.Kind, i.Subject, i.Message)
		}
	}
	if os.Getenv("ATLAS_REPORT") != "" {
		for k, n := range counts {
			fmt.Println(n, k)
		}
		for _, i := range issues {
			fmt.Printf("%s %s %q %s\n", i.Severity, i.Kind, i.Subject, i.Message)
		}
	}
}
