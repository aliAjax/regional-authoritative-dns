package adapter

import (
	"github.com/example/regional-authoritative-dns/internal/record/domain"
	"sort"
)

type Change struct {
	Kind          string
	Before, After domain.Record
}

func Diff(a, b domain.Set) []Change {
	am := map[string]domain.Record{}
	bm := map[string]domain.Record{}
	for _, r := range a {
		am[r.Key()] = r
	}
	for _, r := range b {
		bm[r.Key()] = r
	}
	var out []Change
	for k, r := range am {
		if _, ok := bm[k]; !ok {
			out = append(out, Change{Kind: "remove", Before: r})
		}
	}
	for k, r := range bm {
		if old, ok := am[k]; !ok {
			out = append(out, Change{Kind: "add", After: r})
		} else if old.TTL != r.TTL || old.Healthy != r.Healthy {
			out = append(out, Change{Kind: "update", Before: old, After: r})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Kind < out[j].Kind })
	return out
}
func Added(c []Change) int {
	n := 0
	for _, x := range c {
		if x.Kind == "add" {
			n++
		}
	}
	return n
}
func Removed(c []Change) int {
	n := 0
	for _, x := range c {
		if x.Kind == "remove" {
			n++
		}
	}
	return n
}
