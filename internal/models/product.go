package models

import (
	"fmt"
	"time"
)

type group struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type category struct {
	Id      int    `json:"id"`
	Caption string `json:"caption"`
}

type multilang struct {
	Ru string `json:"ru"`
	Uk string `json:"uk"`
}

type image struct {
	Id           int64  `json:"id"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Url          string `json:"url"`
}

type Product struct {
	Id                   int           `json:"id"`
	ExternalId           interface{}   `json:"external_id"`
	Name                 string        `json:"name"`
	Sku                  string        `json:"sku"`
	Keywords             string        `json:"keywords"`
	Presence             string        `json:"presence"`
	Price                float64       `json:"price"`
	MinimumOrderQuantity interface{}   `json:"minimum_order_quantity"`
	Discount             interface{}   `json:"discount"`
	Prices               []interface{} `json:"prices"`
	Currency             string        `json:"currency"`
	Description          string        `json:"description"`
	Group                group         `json:"group"`
	Category             category      `json:"category"`
	MainImage            string        `json:"main_image"`
	Images               []image       `json:"images"`
	SellingType          string        `json:"selling_type"`
	Status               string        `json:"status"`
	QuantityInStock      interface{}   `json:"quantity_in_stock"`
	MeasureUnit          string        `json:"measure_unit"`
	IsVariation          bool          `json:"is_variation"`
	VariationBaseId      interface{}   `json:"variation_base_id"`
	VariationGroupId     interface{}   `json:"variation_group_id"`
	DateModified         time.Time     `json:"date_modified"`
	Regions              []interface{} `json:"regions"`
	NameMultilang        multilang     `json:"name_multilang"`
}

type Products struct {
	Products []Product `json:"products"`
	GroupId  int       `json:"group_id"`
}

func (p Product) String() string {
	return fmt.Sprintf("ID: %d\tName: %s", p.Id, p.Name)
}
