package models

import (
	"fmt"
	"time"
)

type group struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type category struct {
	ID      int    `json:"id"`
	Caption string `json:"caption"`
}

type image struct {
	ID           int64  `json:"id"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Url          string `json:"url"`
}

type Product struct {
	ID           int       `json:"id"`
	ExternalId   string    `json:"external_id"`
	Name         string    `json:"name"`
	Sku          string    `json:"sku"`
	Keywords     string    `json:"keywords"`
	Presence     string    `json:"presence"`
	Price        float64   `json:"price"`
	Currency     string    `json:"currency"`
	Description  string    `json:"description"`
	Group        group     `json:"group"`
	Category     category  `json:"category"`
	MainImage    string    `json:"main_image"`
	Images       []image   `json:"images"`
	SellingType  string    `json:"selling_type"`
	Status       string    `json:"status"`
	DateModified time.Time `json:"date_modified"`
}

type Products struct {
	Products []Product `json:"products"`
	GroupId  int       `json:"group_id"`
}

type ProductsArray []Product

func (a ProductsArray) Len() int      { return len(a) }
func (a ProductsArray) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ProductsArray) Less(i, j int) bool {
	if a[i].Group.Name == a[j].Group.Name {
		if a[i].Category.Caption == a[j].Category.Caption {
			if a[i].Sku == a[j].Sku {
				if a[i].Price == a[j].Price {
					return a[i].Name < a[j].Name
				}

				return a[i].Price < a[j].Price
			}

			return a[i].Sku < a[j].Sku
		}

		return a[i].Category.Caption < a[j].Category.Caption
	}

	return a[i].Group.Name < a[j].Group.Name
}

func (p Product) String() string {
	return fmt.Sprintf("ID: %d\tНазва: %s", p.ID, p.Name)
}

func (p Product) GetTranslatedStatus() string {
	switch p.Status {
	case "on_display":
		return "опубліковано"
	case "draft":
		return "чернетка"
	case "deleted":
		return "видалено"
	case "not_on_display":
		return "приховано"
	case "editing_required":
		return "потребує редагування"
	case "approval_pending":
		return "очікує модерації"
	case "deleted_by_moderator":
		return "видалено модератором"
	default:
		return p.Status
	}
}
