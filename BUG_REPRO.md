# Bug Reproduction

Zone validation and publication allow frozen or archived state transitions, miss a CNAME conflict, and accept a published version without its timestamp.

Attempt the version and publish operations from frozen or archived zones, include a conflicting CNAME, and omit the publication timestamp. The invalid transitions are accepted and the zone tests fail.

