package model

type Books struct {
	Total int64   `json:"total" desc:"total of zoos"`
	Start int64   `json:"start"`
	Count int64   `json:"count"`
	Books []*Book `json:"books" desc:"books"`
}
