# Bug Reproduction

The TCP DNS handler treats a short read as a complete frame. The rollback handler swallows malformed JSON, JSON errors are built by unsafe string concatenation, and unknown transfer modes are accepted.

Split one TCP DNS frame across reads, send malformed rollback JSON, include quotes in an error message, or submit an unknown transfer mode. The response is truncated, silently defaults, or is invalid JSON, and the transport tests fail.

