package tarot_service

import (
	"reflect"
	"xiaoyuzhou/models"
	"xiaoyuzhou/pkg/logging"
)

func SetPrice(data map[string]interface{}) error {
	logging.Debugf("data: %v, type: %s", data, reflect.TypeOf(data))
	return models.SetPrice(data)

}

func UpdatePrice(data map[string]interface{}) error {
	return models.UpdatePrice(data)
}

func GetPrice() (*models.Price, error) {
	return models.GetPrice()
}

func GetPaymentPrice(cardType, highOrLow string) float32 {
	return models.GetPaymentPrice(cardType, highOrLow)
}
