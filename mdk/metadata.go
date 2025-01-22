package mdk

type ModMetadata struct {
	ID          string
	DisplayName string
	Version     string
	Author      string
	License     string
	Description string
	Extra       map[string]interface{}
}
