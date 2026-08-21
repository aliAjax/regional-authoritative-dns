package domain

type Permission struct {
	Subject, Zone, Action string
	Allow                 bool
}

func (p Permission) Matches(zone, action string) bool {
	return (p.Zone == "*" || p.Zone == zone) && (p.Action == "*" || p.Action == action)
}
