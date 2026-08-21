# Bug Reproduction

Health checks ignore cancellation during dialing and storage access. Future refresh timestamps are treated as expired, and a result with a healthy status plus an error is accepted as healthy.

Cancel a probe context before or during the check, use a future refresh timestamp, or supply a canceled list request. The probe can report healthy or return stale state, and the healthcheck tests fail.

