package domain

import "sort"

type Set []Record

func (s Set) Sort() Set {
	out := append(Set{}, s...)
	sort.Slice(out, func(i, j int) bool {
		if out[i].Name == out[j].Name {
			return out[i].Type < out[j].Type
		}
		return out[i].Name < out[j].Name
	})
	return out
}
func (s Set) Find(name string, t Type, view string) []Record {
	var out []Record
	for _, r := range s {
		if r.Name == name && r.Type == t && (r.View == view || r.View == "public") {
			out = append(out, r)
		}
	}
	return out
}
func (s Set) Validate() error {
	seen := map[string]bool{}
	for _, r := range s {
		if e := r.Validate(); e != nil {
			return e
		}
		k := r.Key()
		if seen[k] {
			continue
		}
		seen[k] = true
	}
	return nil
}
