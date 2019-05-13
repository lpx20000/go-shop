package models

type Dictionary struct {
	Id    uint   `json:"id"`
	Type  string `json:"type"`
	Name  string `json:"name"`
	Value string `json:"value"`
}
