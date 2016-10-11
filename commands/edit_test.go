package commands

import (
	"github.com/ok-borg/borg/types"
	"reflect"
	"testing"
)

func TestTextToProblem(t *testing.T) {
	ps := []types.Problem{
		{
			Title: "title title",
			Solutions: []types.Solution{
				{
					Body: []string{
						"solution 1",
					},
				},
			},
		},
		{
			Title: "asdsadasd",
			Solutions: []types.Solution{
				{
					Body: []string{
						"a multiline snippet\ngood stuff",
					},
				},
				{
					Body: []string{
						"a multiline snippet again \n why not",
					},
				},
			},
		},
	}
	for _, v := range ps {
		text := problemToText(&v)
		p, err := textToProblem(text)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(v, p) {
			t.Fatal(v, p)
		}
	}
}
