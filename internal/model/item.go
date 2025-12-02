package model

type Item struct{}

type Attribute struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	IsSystem bool   `json:"isSystem"`
}
