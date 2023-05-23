package cache_service

import "time"

type XingzuoInput struct {
	Name string
}

func (x *XingzuoInput) GetXingzuoKey() string {
	dateNow := time.Now().Format("0102")
	return x.Name + "_" + dateNow
}
