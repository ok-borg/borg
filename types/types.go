package types

import (
	"time"
)

type Problem struct {
	Id            string     `json:"Id"`
	Title         string     `json:"Title,omitempty"`
	Solutions     []Solution `json:"Solutions,omitempty"`
	ImportMeta    ImportMeta `json:"ImportMeta,omitempty"`
	CreatedBy     string     `json:"CreatedBy,omitempty"`
	Created       time.Time  `json:"Created,omitempty"`
	LastUpdatedBy string     `json:"LastUpdatedBy,omitempty"`
	LastUpdated   time.Time  `json:"Updated,omitempty"`
}

type ImportMeta struct {
	Source int    `json:"Source"` // enum, 0 stackoverflow
	Id     string `"json:"Id"`
}

type Solution struct {
	Body  []string `json:"Body,omitempty"`
	Score int      `json:"Score,omitempty"`
}

type Solutions []Solution

func (a Solutions) Len() int           { return len(a) }
func (a Solutions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Solutions) Less(i, j int) bool { return a[i].Score > a[j].Score }
