package main

import (
	"os"
	"strings"
	"testing"
)

func TestFixtureImageStartsHTTPServiceAsNonRoot(t *testing.T) {
	data, err := os.ReadFile("Dockerfile")
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{"USER 65532:65532", "ENTRYPOINT [\"/fixture-service\"]", "EXPOSE 8080"} {
		if !strings.Contains(string(data), required) {
			t.Fatalf("fixture image does not provide HTTP startup contract: missing %s", required)
		}
	}
}
