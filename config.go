package rest

// VersionIn -- where is the version In the HTTP
type VersionIn string

const (
	// VersionInPath like:
	// GET /basePath/v1 HTTP/1.1
	// Host: 127.0.0.1
	VersionInPath VersionIn = "path"
	// VersionInHeader like:
	// GET /basePath HTTP/1.1
	// Host: 127.0.0.1
	// version: v1
	VersionInHeader VersionIn = "header"
)

// Config config for server
// Fields:
//   Schemes -- such as SCHEMES_HTTP, SCHEMES_HTTPS, SCHEMES_HTTP_AND_HTTPS
//   Host -- such as "127.0.0.1:8000"
//   DomainName -- domain name
//   BasePath -- such as "/apis"
//   Version -- such as "v1.0.0"
//   VersionIn -- where is the version In the HTTP (path, or header).
//   Title -- the title of server
//   Description -- the description of server
//   CrossOrigin -- Cross-Origin Resource Sharing
//   OpenDocsSwaggerJSONUrl -- open the url /docs/swagger.json
//   DocsSwaggerJSONUrl -- the url to get openAPI(swagger) document, default value is /docs/swagger.json
type Config struct {
	Schemes                []Scheme
	Host                   string
	DomainName             string
	BasePath               string
	Version                string
	VersionIn              VersionIn
	Title                  string
	Description            string
	CrossOrigin            bool
	OpenDocsSwaggerJSONUrl bool
	DocsSwaggerJSONUrl     string
}
