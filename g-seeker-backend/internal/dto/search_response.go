package dto

type SearchItem struct {
	Name        string `json:"name"`
	Owner       string `json:"owner"`
	URL         string `json:"url"`
	Stars       int    `json:"stars"`
	Description string `json:"description"`
	Reason      string `json:"reason"`
}

type SearchResponse struct {
	Items []SearchItem `json:"items"`
}
