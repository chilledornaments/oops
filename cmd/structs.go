package main

type newSecret struct {
	Secret string `json:"secret"`
}

type createTemplateData struct {
	CreateEndpoint string
}

type successTemplateData struct {
	URL string
}

type successJSON struct {
	URL string `json:"url"`
}
