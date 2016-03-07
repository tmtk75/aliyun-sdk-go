package api

type Http struct {
	Method string `json:"method"`
	Uri    string `json:"uri"`
}

type Member struct {
	Required bool   `json:"required"`
	Type     string `json:"type"`     // e.g) string, integer
	Location string `json:"location"` // e.g) header
	Name     string `json:"name"`     // e.g) x-oss-acl
	Format   string
}

type Ope struct {
	Name  string `json:"name"`
	Http  Http   `json:"http"`
	Input struct {
		Type    string            `json:"structure"`
		Members map[string]Member `json:"members"`
	} `json:"input"`
}

type APIDef struct {
	Format              string         `json:"format"`
	APIVersion          string         `json:"apiVersion"`
	ChecksumFormat      string         `json:"checksumFormat"`
	EndpointPrefix      string         `json:"endpointPrefix"`
	ServiceAbbreviation string         `json:"serviceAbbreviation"`
	ServiceFullName     string         `json:"serviceFullName"`
	SignatureVersion    string         `json:"signatureVersion"`
	TimestampFormat     string         `json:"timstampFormat"`
	Operations          map[string]Ope `json:"operations"`
	// This following field is not in aliyun service definition
	// Provided by logic
	Imports map[string]bool
}
