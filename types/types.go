package types

type Problem struct {
	Id         string     `json:"id"`
	Title      string     `json:"title"`
	Solutions  []Solution `json:"solutions"`
	ImportMeta ImportMeta `json:"importMeta"`
}

type ImportMeta struct {
	Source int    `json:"source"` // enum, 0 stackoverflow
	Id     string `"json:"id"`
}

type Solution struct {
	Body  []string `json:"body"`
	Score int      `json:"score"`
}

type Solutions []Solution

func (a Solutions) Len() int           { return len(a) }
func (a Solutions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Solutions) Less(i, j int) bool { return a[i].Score > a[j].Score }
