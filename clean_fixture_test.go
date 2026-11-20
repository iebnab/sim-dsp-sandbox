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

func TestRawSecretEvidenceMustMatchScannedArtifact(t *testing.T) {
	data, err := os.ReadFile(".github/workflows/publish.yml")
	if err != nil { t.Fatal(err) }
	text := string(data)
	start := strings.Index(text, "# Raw secret evidence validation.")
	end := strings.Index(text, "# End raw secret evidence reduction.")
	if start < 0 || end <= start { t.Fatal("raw secret evidence is not validated before reduction") }
	gate := text[start:end]
	valid := `{"SchemaVersion":2,"Trivy":{"Version":"0.72.0"},"ArtifactName":"repo@sha256:abc","ArtifactType":"container_image","ArtifactID":"image-id","Metadata":{"ImageID":"sha256:abc","RepoDigests":["repo@sha256:abc"]}}`
	for _, raw := range []string{`{}`, `null`, `{"Results":null}`, strings.Replace(valid, `"container_image"`, `null`, 1), strings.Replace(valid, `"0.72.0"`, `"0.71.0"`, 1), strings.ReplaceAll(valid, `repo@sha256:abc`, `wrong@sha256:abc`), strings.TrimSuffix(valid, "}")+`,"Results":{}}`, strings.TrimSuffix(valid, "}")+`,"Results":null}`, strings.TrimSuffix(valid, "}")+`,"Results":[{"Target":"/file","Class":"secret","Secrets":null}]}`} {
		dir := t.TempDir()
		if err := os.WriteFile(dir+"/secrets-raw.json", []byte(raw), 0600); err != nil { t.Fatal(err) }
		cmd := exec.Command("bash", "-euo", "pipefail", "-c", gate)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "ref=repo@sha256:abc")
		if err := cmd.Run(); err == nil { t.Fatalf("invalid raw report accepted: %s", raw) }
	}
	for _, raw := range []string{valid, strings.TrimSuffix(valid, "}")+`,"Results":[]}`, strings.TrimSuffix(valid, "}")+`,"Results":[{"Target":"/file","Class":"secret"}]}`} {
		dir := t.TempDir()
		if err := os.WriteFile(dir+"/secrets-raw.json", []byte(raw), 0600); err != nil { t.Fatal(err) }
		cmd := exec.Command("bash", "-euo", "pipefail", "-c", gate)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "ref=repo@sha256:abc")
		if output, err := cmd.CombinedOutput(); err != nil { t.Fatalf("legitimate zero report refused: %v %s", err, output) }
		if reduced, err := os.ReadFile(dir+"/secrets.json"); err != nil || strings.TrimSpace(string(reduced)) != "[]" { t.Fatal("zero report did not reduce to empty results") }
	}
}
