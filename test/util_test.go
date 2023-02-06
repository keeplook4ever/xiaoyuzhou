package test

import (
	"fmt"
	"testing"
	"time"
	"xiaoyuzhou/pkg/util"
)

func TestSwap(t *testing.T) {
	starNum := util.RandFromRange(100, 800)
	time.Sleep(time.Microsecond)
	readNum := util.RandFromRange(100, 800)

	fmt.Printf("Origin: readNum:%d, starNum:%d", readNum, starNum)
	if starNum > readNum {
		util.SwapTwoInt(&starNum, &readNum)
	}
	fmt.Printf("Swaped: readNum:%d, starNum:%d", readNum, starNum)
}
