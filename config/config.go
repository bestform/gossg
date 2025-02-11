package config

import (
	"encoding/json"
	"os"
)

const (
	CONTENT_TYPE_HTML       = "html"
	CONTENT_TYPE_MARKDOWN   = "markdown"
	CONTENT_TYPE_BLOG_INDEX = "blogindex"
)

type Config struct {
	Title      string `json:"title"`
	HeaderFile string `json:"header_file"`
	FooterFile string `json:"footer_file"`
	Pages      []Page `json:"pages"`
}

type Page struct {
	TargetFilename string `json:"target_filename"`
	Title          string `json:"title"`
	ContentType    string `json:"content_type"`
	ContentPath    string `json:"content_path"`
	Pages          []Page `json:"pages"`
}

func ReadConfig(path string) (Config, error) {
	var config Config
	content, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	json.Unmarshal(content, &config)

	return config, nil
}
