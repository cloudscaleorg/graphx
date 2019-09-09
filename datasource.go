package graphx

import "encoding/json"

// DataSource
type DataSource struct {
	Name       string `json:"name"`
	ConnString string `json:"connection_string"`
}

func (d *DataSource) ToJSON() ([]byte, error) {
	b, err := json.Marshal(d)
	return b, err
}

func (d *DataSource) FromJSON(b []byte) error {
	err := json.Unmarshal(b, d)
	return err
}
