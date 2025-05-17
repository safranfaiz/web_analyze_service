package response

type SuccessResponse struct {
	HtmlVersion         string    `json:"htmlVersion"`
	Title               string    `json:"title"`
	ServiceTime         int64     `json:"serviceTime"`
	WebPageExtractTime  int64     `json:"webPageExtractTime"`
	Headings            []Heading `json:"headings"`
	Urls                []Url     `json:"urls"`
	HasLogin            bool      `json:"hasLogin"`
	ExecutedUrl         string    `json:"executedUrl"`
	BasePath            string    `json:"basePath"`
	AppExecuteTotalTime int64     `json:"appExecuteTotalTime"`
}

type Heading struct {
	Tag  string `json:"tag"`
	Text string `json:"text"`
}

type Url struct {
	Url              string `json:"url"`
	Accessible       bool   `json:"accessible"`
	Type             string `json:"type"`
	Status           int    `json:"status"`
	UrlExecutionTime int64  `json:"urlExecutionTime"`
}
