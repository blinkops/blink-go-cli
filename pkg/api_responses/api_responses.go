package api_responses

type GetIdByNameResponse struct {
	Results []struct {
		Id string `json:"id"`
	} `json:"results"`
}

type CreateResponseWithId struct {
	Id string `json:"id"`
}
