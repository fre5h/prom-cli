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

type image struct {
	Id           int64  `json:"id"`
	ThumbnailUrl string `json:"thumbnail_url"`
	Url          string `json:"url"`
}

type Product struct {
	Id           int       `json:"id"`
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

func (p Product) String() string {
	return fmt.Sprintf("ID: %d\tНазва: %s", p.Id, p.Name)
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
