package utils

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
)

func GetJson(url string, result interface{}) {
	get := fiber.Get(url)
	_, body, _ := get.String()

	err := json.Unmarshal([]byte(body), &result)
	if err != nil {
		panic(err.Error())
	}

}
