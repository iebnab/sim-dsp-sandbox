# Synthetic image fixture, never deployed as an application.
# The deliberately old immutable base exercises vulnerability fix age.
FROM alpine:3.18.0@sha256:02bb6f428431fbc2809c5d1b41eab5a68350194fb508869a33cb1af4444c9b11
# This sentinel is deliberately fake and not an issued service-account credential.
RUN printf '%s\n' '{"type": "service_account", "project_id": "synthetic-f11-sentinel", "private_key_id": "NOT-A-REAL-KEY"}' \
    > /synthetic-credential.json
USER 65532:65532
ENTRYPOINT ["/bin/echo", "synthetic bidder stand-in"]
