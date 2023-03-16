package tarot_service

import (
	"xiaoyuzhou/models"
)

func SetPrice(data map[string]interface{}) error {
	return models.SetPrice(data)
}

func UpdatePrice(data map[string]interface{}) error {
	return models.UpdatePrice(data)
}

func GetPriceTotal() ([]models.Price, error) {
	return models.GetPriceTotal()
}

func GetPaymentPrice(scene string, location string) float32 {
	// enums:"ta_one_high,ta_one_low,ta_three_high,ta_three_low"
	// location: jp,zh,en,tc
	return models.GetPaymentPrice(scene, location)
}
