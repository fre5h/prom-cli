package models

import "fmt"

type ProductUpdate struct {
	Id    string `json:"id"`
	Sku   string `json:"sku"`
	Price string `json:"price"`
}

func (p ProductUpdate) String() string {
	return fmt.Sprintf("ID: %s\tКод/Артикул: %s\tЦіна: %s\n", p.Id, p.Sku, p.Price)
}
