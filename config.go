package rest

// Config config for server
// Fields:
//   Schemes -- such as SCHEMES_HTTP, SCHEMES_HTTPS, SCHEMES_HTTP_AND_HTTPS
//   Host -- such as "127.0.0.1:8000"
//   DomainName -- domain name, like github.com
//   BasePath -- such as "/apis"
//   Version -- such as "v1.0.0"
//   Title -- the title of server
//   Description -- the description of server
//   AllowOrigin -- Allow Origin
//   Origins -- ( Cross-Origin Resource Sharing ) Access-Control-Allow-Origin
//   OpenAPIDocumentURL -- open the url /docs/swagger.json
//   APIDocumentURL -- the url to get openAPI(swagger) document, default value is /docs/swagger.json
type Config struct {
	Schemes            []Scheme
	Host               string
	DomainName         string
	BasePath           string
	Version            string
	Title              string
	Description        string
	AllowOrigin        bool
	Origins            []string
	OpenAPIDocumentURL bool
	APIDocumentURL     string
}
