package main

import (
	"reflect"
	"sort"
	"testing"
)

func Test_equal(t *testing.T) {
	scoreList := []int{2, 6, 4, 8}
	sortedScoreList := make([]int, len(scoreList))
	copy(sortedScoreList, scoreList)
	sort.Ints(sortedScoreList)
	t.Log(sortedScoreList)
	if !reflect.DeepEqual(sortedScoreList, scoreList) {
		t.Error("not equal")
	}
	t.Log(scoreList)
	t.Log(sortedScoreList)
}
