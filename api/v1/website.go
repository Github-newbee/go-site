package v1

type WebsiteRequest struct {
	Name        string `json:"name"`
	Url         string `json:"url"`
	Icon        string `json:"icon"`
	CategoryID  string `json:"category_id"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}
