package entities

// Typedef to parse the JSON body of the request correctly.
type CWFBody_t struct {
	File string `json:"file"`
	Content string `json:"content"`
}
