# Bug Reproduction

Resolver policy selection sorts caller-owned state in place, the limiter stops at the wrong boundary, oversized labels are accepted, and truncation shares packet storage with the input.

Run the concurrent policy and limiter scenarios, parse an oversized label, and mutate a truncated packet after creation. The race detector or resolver tests report the incorrect state or ownership.

