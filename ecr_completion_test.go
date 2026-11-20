package fixture

import (
	"os"
	"strings"
	"testing"
)

func TestEcrDispatchCopiesOnlyOriginalVerifiedCorpus(t *testing.T) {
	data, err := os.ReadFile(".github/workflows/publish.yml")
	if err != nil { t.Fatal(err) }
	text := string(data)
	if !strings.Contains(text, "workflow_dispatch:") || !strings.Contains(text, "  ecr_completion:") {
		t.Fatal("ECR-only dispatch does not exist")
	}
	job := strings.SplitN(text, "  ecr_completion:", 2)[1]
	for _, required := range []string{"44780460309", "b594998574a73ecfca38933411e83b5bb0daeabac226bf13bf13e8eb92a7745c", "86af7ecfb5be38d58104e78b82f06cee6839cb61b47839d00affbae58a98f5c9", "342843238294.dkr.ecr.eu-central-1.amazonaws.com/sim-supply-chain", "--type vuln", "--type slsaprovenance1", "source-receipts.json", "cosign copy"} {
		if !strings.Contains(job, required) { t.Fatalf("missing original ECR evidence binding: %s", required) }
	}
	for _, forbidden := range []string{"docker build", "trivy image", "cosign sign", "/api/ingest/ci/sbom", "/api/ingest/ci/verification", "docker.io", "workflow_dispatch.inputs"} {
		if strings.Contains(job, forbidden) { t.Fatalf("ECR-only completion widens execution: %s", forbidden) }
	}
}

func TestEcrDispatchCannotTriggerFullProducer(t *testing.T) {
	data, err := os.ReadFile(".github/workflows/publish.yml")
	if err != nil { t.Fatal(err) }
	text := string(data)
	for _, job := range []string{"build", "scan", "publish"} {
		start := strings.Index(text, "  "+job+":")
		if start < 0 { t.Fatalf("missing existing %s job", job) }
		section := strings.SplitN(text[start:], "\n\n  ", 2)[0]
		if !strings.Contains(section, "if: github.event_name == 'push'") { t.Fatalf("%s is not push-only", job) }
	}
}
