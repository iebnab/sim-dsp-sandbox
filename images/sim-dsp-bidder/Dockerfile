# Synthetic image fixture, never deployed as an application.
# The deliberately old immutable base exercises vulnerability fix age.
FROM alpine:3.18.0@sha256:02bb6f428431fbc2809c5d1b41eab5a68350194fb508869a33cb1af4444c9b11
# This sentinel is deliberately fake, never issued by AWS and has no secret pair.
RUN printf '%s\n' '# SYNTHETIC ONLY - not a usable credential' \
    'AWS_ACCESS_KEY_ID=AKIAF11SENTINEL000001' > /synthetic-credential.txt
USER 65532:65532
ENTRYPOINT ["/bin/echo", "synthetic bidder stand-in"]
