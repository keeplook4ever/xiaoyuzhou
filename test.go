package main

import (
	"fmt"
	"xiaoyuzhou/pkg/logging"
	"xiaoyuzhou/pkg/qiniu"
)

func main11() {
	filepath := "D:\\Users\\liu.ning\\Pictures\\5b38e720e0624f6ca6d6fe51334a0721.jpeg"

	//err := qiniu.UploadLocalImg("D:\\Users\\liu.ning\\Pictures\\5b38e720e0624f6ca6d6fe51334a0721.jpeg")
	err := qiniu.UploadLocalImg(filepath)

	if err != nil {
		logging.Fatal(err)
		return
	}
	fmt.Println("put success")

}
