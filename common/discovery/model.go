package discovery

import "encoding/json"

type EndPointInfo struct {
	IP       string                 `json:"ip"`
	Port     string                 `json:"port"`
	MetaData map[string]interface{} `json:"meta"`
}

func UnMarshal(data []byte) (*EndPointInfo, error) {
	e := &EndPointInfo{}
	err := json.Unmarshal(data, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *EndPointInfo) Marshal() string {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}
	return string(data)
}
