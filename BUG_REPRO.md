# Bug Reproduction

DNSSEC digest canonicalization is order- and case-sensitive, inactive keys are accepted, empty record sets verify, and key activity windows do not enforce creation and expiry boundaries.

Verify equivalent records in different order and name case, then use inactive, empty, not-yet-created, or expired keys. The digest or validation result is unsafe and the DNSSEC tests fail.

