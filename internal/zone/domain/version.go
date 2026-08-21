package domain

import (
	"fmt"
	"sort"
	"time"
)

type Version struct {
	ID, ZoneID, Author, Reason string
	Serial                     uint32
	Status                     Status
	Records                    int
	Digest                     string
	CreatedAt                  time.Time
	PublishedAt                *time.Time
}

func (v Version) Validate() error {
	if v.ID == "" || v.ZoneID == "" {
		return fmt.Errorf("version identity required")
	}
	if v.Serial == 0 {
		return fmt.Errorf("version serial required")
	}
	if v.Records < 0 {
		return fmt.Errorf("invalid record count")
	}
	return nil
}
func SortVersions(v []Version) []Version {
	out := append([]Version{}, v...)
	sort.Slice(out, func(i, j int) bool { return out[i].Serial > out[j].Serial })
	return out
}
func (v Version) Published() bool { return v.Status == Published && v.PublishedAt != nil }
func (v Version) Summary() string {
	return fmt.Sprintf("%s serial=%d records=%d", v.ID, v.Serial, v.Records)
}
