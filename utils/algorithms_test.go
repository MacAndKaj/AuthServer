package utils

import (
	"testing"
)

var findTestingSets = []struct {
	elem  int
	set   []int
	index int
}{
	{
		elem:  5745,
		set:   []int{55, 512, 5745, 324, 122, 5654},
		index: 2,
	},
	{
		elem:  5745,
		set:   []int{5745, 324, 122, 5654},
		index: 0,
	},
	{
		elem:  5745,
		set:   []int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 55, 512, 5745},
		index: 13,
	},
	{
		elem:  0,
		set:   []int{},
		index: -1,
	},
	{
		elem:  5745,
		set:   []int{5745},
		index: 0,
	},
}

func TestFind(t *testing.T) {
	for _, set := range findTestingSets {
		i := Find(set.set, set.elem)
		if i != set.index {
			t.Errorf("Incorrect index: %d. Expected: %d", i, set.index)
		}
	}
}
