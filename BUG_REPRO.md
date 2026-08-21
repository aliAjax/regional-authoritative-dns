# Bug Reproduction

Authoritative resolution tries a parent zone before a more specific child zone. Query defaults, zone ordering, and serial wraparound handling are also inconsistent.

Create parent and child zones with the same queried name, use mixed-case and missing-dot queries, list zones repeatedly, and compare serials across wraparound. The parent NXDOMAIN wins or the boundary behavior is unstable, and the resolver/zone tests fail.

