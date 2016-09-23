package types

type Problem struct {
	Id         string     `json:"Id"`
	Title      string     `json:"Title"`
	Solutions  []Solution `json:"Solutions"`
	ImportMeta ImportMeta `json:"ImportMeta"`
}

type ImportMeta struct {
	Source int    `json:"Source"` // enum, 0 stackoverflow
	Id     string `"json:"Id"`
}

type Solution struct {
	Body  []string `json:"Body"`
	Score int      `json:"Score"`
}

type Solutions []Solution

func (a Solutions) Len() int           { return len(a) }
func (a Solutions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Solutions) Less(i, j int) bool { return a[i].Score > a[j].Score }
