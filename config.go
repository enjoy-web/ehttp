package rest

// Config config for server
// Fields:
//   Schemes -- such as SCHEMES_HTTP, SCHEMES_HTTPS, SCHEMES_HTTP_AND_HTTPS
//   Host -- such as "127.0.0.1:8000"
//   DomainName -- domain name
//   BasePath -- such as "/apis"
//   Version -- such as "v1.0.0"
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
	Title                  string
	Description            string
	CrossOrigin            bool
	OpenDocsSwaggerJSONUrl bool
	DocsSwaggerJSONUrl     string
}
