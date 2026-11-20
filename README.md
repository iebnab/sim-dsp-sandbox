# sim-dsp-sandbox

Synthetic test fixture for a fictional ad-tech demand-side platform. No real code,
data or secrets.

It is the weakly protected variant used to test a read-only posture product:

- the default branch is protected, but a pull request needs no approving review,
  admins are exempt, and force pushes are allowed;
- the workflow requests a write-all token and uses an unpinned third-party action,
  and runs only on manual dispatch.

The signed supply-chain fixture images use the main-only publish workflow. The
sim-dsp-tunnel-client image is a minimal non-root HTTP and batch-job fixture:
it listens on 0.0.0.0:$PORT (default 8080) and returns constant `ok` for every
request. `--job` exits successfully without opening a listener. Cloud Run Jobs
may instead override the command with `/bin/sh -c 'exit 0'`.

It makes no outbound calls and contains no issued credentials. The deliberately
fake service-account-shaped file is retained only to exercise value-free F11
secret findings. The workflow scans the immutable build digest, waits for the
platform receipt, signs and attests it, and verifies signature plus vulnerability
and SLSA attestations before registry copies. An Artifact Registry copy is not
destination-side verification; fixture admission uses only the exact verified
digest and never changes a default deny policy to allow arbitrary images.

Image build runs the three local Go fixture assertions before compiling the
static binary. The bidder's inert Log4j jar remains a never-executed KEV witness.

The clean Binary Authorization fixture is distinct: `cloudrun-fixture` uses only
the pinned static Go binary on scratch, runs as UID/GID 65532, and contains no
planted credential file or vulnerable base packages. Its publish workflow runs
only that image and refuses nonzero vulnerability or secret evidence before
uploading a complete scan. Signing, both attestations, exact receipt matching
and verification remain mandatory. It copies only to the already authorized
private Artifact Registry repository; no Docker Hub or ECR target is used.

The original planted bidder/tunnel sources remain unchanged. Their complete
workflow and sentinel guards are retained byte-for-byte at
`.github/f11-publish.yml.disabled`, outside GitHub's executable workflow directory.
The scratch fixture's job mode is `--job`; it has no shell to override.

The bounded ECR completion is a manual dispatch of the same publish workflow,
not a rebuild or a new scan/signing run. Its fixed matrix copies only the two
already verified F11 digests from run 44780460309 into the exact private ECR
repository, using fresh release-environment OIDC. It validates and retains the
original receipts, verifies both source and destination signature/vulnerability/
SLSA evidence, and records the new destination execution separately. It uploads
no new scan or same-execution platform verification assertion. Automatic push
triggers are disabled during this completion so deploying it cannot rerun the
full producer. There are no caller-supplied repository, digest or source-run inputs.
