package models

import "fmt"

type ProductUpdate struct {
	ID    int     `json:"id"`
	Price float64 `json:"price"`
}

func (p ProductUpdate) String() string {
	return fmt.Sprintf("ID: %d\tЦіна: %.2f", p.ID, p.Price)
}
