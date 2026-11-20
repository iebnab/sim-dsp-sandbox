package fixture

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func cleanGate(t *testing.T) string {
	t.Helper()
	data, err := os.ReadFile(".github/workflows/publish.yml")
	if err != nil {
		t.Fatal(err)
	}
	text := string(data)
	start := strings.Index(text, "# Clean image evidence gate.")
	end := strings.Index(text, "# End clean image evidence gate.")
	if start < 0 || end <= start {
		t.Fatal("clean image publish lacks a fail-closed zero vulnerability/secret evidence gate")
	}
	return text[start:end]
}

func TestCleanPublishRequiresCompleteZeroEvidence(t *testing.T) {
	gate := cleanGate(t)
	for _, test := range []struct{ sbom, secrets string; passes bool }{
		{`{"bomFormat":"CycloneDX","components":[{"name":"fixture"}]}`, `[]`, true},
		{`{"bomFormat":"CycloneDX","components":[{"name":"fixture"}],"vulnerabilities":[{"id":"CVE-1"}]}`, `[]`, false},
		{`{"bomFormat":"CycloneDX","components":[{"name":"fixture"}]}`, `[{"rule_id":"secret"}]`, false},
		{`{}`, `[]`, false},
		{`not-json`, `[]`, false},
	} {
		dir := t.TempDir()
		if err := os.WriteFile(dir+"/sbom.json", []byte(test.sbom), 0600); err != nil { t.Fatal(err) }
		if err := os.WriteFile(dir+"/secrets.json", []byte(test.secrets), 0600); err != nil { t.Fatal(err) }
		cmd := exec.Command("bash", "-euo", "pipefail", "-c", gate)
		cmd.Dir = dir
		if err := cmd.Run(); (err == nil) != test.passes {
			t.Fatalf("complete zero evidence passes=%v, result=%v", test.passes, err)
		}
	}
}

func TestCleanPublishHasOnlyAuthorizedPrivateTarget(t *testing.T) {
	data, err := os.ReadFile(".github/workflows/publish.yml")
	if err != nil { t.Fatal(err) }
	text := string(data)
	if strings.Count(text, "image: [cloudrun-fixture]") != 3 {
		t.Fatal("clean run must build, scan and publish only cloudrun-fixture")
	}
	for _, forbidden := range []string{"DOCKERHUB_TOKEN", "ECR_REPOSITORY", "cosign copy \"$staged\""} {
		if strings.Contains(text, forbidden) { t.Fatalf("unauthorized clean destination: %s", forbidden) }
	}
	for _, required := range []string{"receipt-cloudrun-fixture", "scan-cloudrun-fixture", "slsaprovenance1", "--type vuln", "AR_REPOSITORY"} {
		if !strings.Contains(text, required) { t.Fatalf("missing clean evidence contract: %s", required) }
	}
}
