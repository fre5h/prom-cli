package models

import "fmt"

type Group struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Description   string `json:"description"`
	Image         string `json:"image"`
	ParentGroupId int    `json:"parent_group_id"`
}

type Groups struct {
	Groups []Group `json:"groups"`
}

func (g Group) String() string {
	return fmt.Sprintf("ID: %d\tНазва: %s", g.Id, g.Name)
}
