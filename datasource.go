package graphx

import "encoding/json"

// DataSource
type DataSource struct {
	// unique name for this datasource
	Name string `json:"name"`
	// the datasource type such as "prometheus"
	Type string `json:"type"`
	// the connection string for the above type
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
