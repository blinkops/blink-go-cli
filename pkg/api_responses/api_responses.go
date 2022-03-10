package api_responses

type GetPlaybookIdByNameResponse struct {
	Results []struct {
		Id string `json:"id"`
	} `json:"results"`
}

type CreatePlaybookResponse struct {
	Id string `json:"id"`
}
