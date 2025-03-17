package api

type HttpResponse struct {
	IsError bool        `json:"is_error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Auth struct {
	Identifier string `json:"identifier"`
	SecretKey  string `json:"secret_key"`
}

type Category struct {
	Category string `bson:"category" json:"category"`
	Sum      string `bson:"sum" json:"sum"`
}

type Sum struct {
	Type       string     `bson:"type" json:"type"`
	Sum        string     `bson:"sum" json:"sum"`
	Categories []Category `bson:"categories" json:"categories"`
}
