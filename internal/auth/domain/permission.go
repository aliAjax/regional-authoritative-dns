package domain

type Permission struct {
	Subject, Zone, Action string
	Allow                 bool
}

func (p Permission) Matches(zone, action string) bool {
	// A deny rule never authorizes: Allow=false means the rule is an explicit
	// override, so it must not be treated as a granting match.
	if !p.Allow {
		return false
	}
	return (p.Zone == "*" || p.Zone == zone) && (p.Action == "*" || p.Action == action)
}
