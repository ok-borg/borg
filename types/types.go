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
	Source int    `json:"Source,omitempty"` // enum, 0 stackoverflow
	Id     string `json:"Id,omitempty"`
}

type Solution struct {
	Body  []string `json:"Body,omitempty"`  // this was a mistake to make it a string - after db correction and refactoring should get rid of it
	Score int      `json:"Score,omitempty"` // this has values in the DB but they are not being used for anything. should nuke it. editing an entry nukes it anyway
}

type Solutions []Solution

func (a Solutions) Len() int           { return len(a) }
func (a Solutions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Solutions) Less(i, j int) bool { return a[i].Score > a[j].Score }
