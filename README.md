# sim-dsp-sandbox

Synthetic test fixture for a fictional ad-tech demand-side platform. No real code,
data or secrets.

It is the weakly protected variant used to test a read-only posture product:

- the default branch is protected, but a pull request needs no approving review,
  admins are exempt, and force pushes are allowed;
- the workflow requests a write-all token and uses an unpinned third-party action,
  and runs only on manual dispatch.
