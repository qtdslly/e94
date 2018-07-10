package apimodel

type Video struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	ContentType uint8     `json:"content_type"`
	Score       string    `json:"score"`
	ThumbX      string    `json:"thumb_x"`
	ThumbY      string    `json:"thumb_y"`
	PublishDate string    `json:"publish_date"`
	Year        string    `json:"year"`
	Language    string    `json:"language"`
	Country     string    `json:"country"`
	Directors   string    `json:"directors"`
	Actors      string    `json:"actors"`
	Tags        string    `json:"tags"`
	PageUrl     string    `json:"page_url"`
	Provider    uint32    `json:"provider"`
}

func VideoToHtml(src Video) string {
	result := "<div>"
	result = result + "<img src='" + src.ThumbY + "' alt='" + src.Title + "' />"
	result = result + "<h3>" + src.Title + "</h3>"
	result = result + "<span>" + src.Description + "</span>"
	result = result + "</div>"

	return result
}
