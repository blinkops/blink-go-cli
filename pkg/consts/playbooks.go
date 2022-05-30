package consts

const (
	GetAutomationPackByDisplayNameQueryFormat = `{"limit":1,"offset":0,"filter":{"display_name":{"$eq":"%s"}},"ci":true,"select":["id"]}`
)
