package configuration

import (
	"encoding/json"
	"os"
)

type Config struct {
	Host struct {
		Port string `json:"port"`
	} `json:"host"`

	ApiInternal struct {
		UrlInserted string `json:"url_inserted"`
		UrlDeleted  string `json:"url_deleted"`
		UrlUpdated  string `json:"url_updated"`
	} `json:"api_internal"`
}

func LoadConfig(file string, cfg interface{}) error {
	r, err := os.Open(file)
	if err != nil {
		return err
	}
	defer r.Close()

	return json.NewDecoder(r).Decode(&cfg)
}
