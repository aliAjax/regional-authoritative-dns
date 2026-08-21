# Bug Reproduction

Resolver and publication paths ignore canceled contexts. Canceled requests can read or write stale cache and publication state, and events missing identity or action are considered ready.

Cancel cache get/put and publication record/save requests, then submit an incomplete event. Canceled work remains visible or the incomplete event is recorded, and the resolver/publication tests fail.

