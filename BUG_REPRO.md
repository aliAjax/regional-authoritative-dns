# Bug Reproduction

Record normalization overwrites an explicitly unhealthy record, and view lookup can mix public and private records. Normalized duplicates and lowercase types are also accepted along the same path.

Trigger the record upsert/parse/lookup flow with an unhealthy record, equivalent names in different case or trailing-dot forms, and both a public and an exact-view record. The result reports the unhealthy address as healthy or returns the wrong view, and the dedicated record tests fail.

