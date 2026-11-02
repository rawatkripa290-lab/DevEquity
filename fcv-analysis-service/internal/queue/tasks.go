package queue

import "encoding/json"

const TypeIngestRepository = "repository:ingest"

type IngestPayload struct {
	AnalysisID string `json:"analysis_id"`
	Owner      string `json:"owner"`
	Repo       string `json:"repo"`
}

func EncodeIngest(p IngestPayload) ([]byte, error) { return json.Marshal(p) }
func DecodeIngest(data []byte) (IngestPayload, error) {
	var p IngestPayload
	err := json.Unmarshal(data, &p)
	return p, err
}
