package types

type Problem struct {
	Id         string
	Title      string
	Solutions  []Solution
	ImportMeta ImportMeta
}

type ImportMeta struct {
	Source int // enum, 0 stackoverflow
	Id     string
}

type Solution struct {
	Body  []string
	Score int
}

type Solutions []Solution

func (a Solutions) Len() int           { return len(a) }
func (a Solutions) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Solutions) Less(i, j int) bool { return a[i].Score > a[j].Score }
