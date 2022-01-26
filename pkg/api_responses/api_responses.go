package api_responses

type GetPlaybookIdByNameResponse struct {
	Results []struct {
		Id              string        `json:"id"`
	} `json:"results"`
}
