package ehttp

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type corsInfos struct {
	options *corsInfo
	get     *corsInfo
	put     *corsInfo
	post    *corsInfo
	patch   *corsInfo
	delete  *corsInfo
}

func (c *corsInfos) OPTIONS() *corsInfo {
	if c.options == nil {
		c.options = &corsInfo{}
	}
	return c.options
}

func (c *corsInfos) GET() *corsInfo {
	if c.get == nil {
		c.get = &corsInfo{}
	}
	return c.get
}

func (c *corsInfos) PUT() *corsInfo {
	if c.put == nil {
		c.put = &corsInfo{}
	}
	return c.put
}

func (c *corsInfos) POST() *corsInfo {
	if c.post == nil {
		c.post = &corsInfo{}
	}
	return c.post
}

func (c *corsInfos) PATCH() *corsInfo {
	if c.patch == nil {
		c.patch = &corsInfo{}
	}
	return c.patch
}
func (c *corsInfos) DELETE() *corsInfo {
	if c.delete == nil {
		c.delete = &corsInfo{}
	}
	return c.delete
}

type corsInfo struct {
	Methods     map[string]bool
	Headers     map[string]bool
	Origins     map[string]bool
	Credentials bool
}

func (c *corsInfo) addMethod(method string) {
	if c.Methods == nil {
		c.Methods = make(map[string]bool, 0)
	}
	c.Methods[method] = true
}

func (c *corsInfo) addHeader(header string) {
	if c.Headers == nil {
		c.Headers = make(map[string]bool, 0)
	}
	c.Headers[header] = true
}

func (c *corsInfo) addOrigin(origin string) {
	if c.Origins == nil {
		c.Origins = make(map[string]bool, 0)
	}
	c.Origins[origin] = true
}

func (c corsInfo) toAccessControlAllow() *accessControlAllow {
	access := &accessControlAllow{Credentials: c.Credentials}
	access.setMethods(c.Methods)
	access.setHeaders(c.Headers)
	access.setOrigin(c.Origins)
	return access
}

type accessControlAllows struct {
	OPTIONS *accessControlAllow
	GET     *accessControlAllow
	PUT     *accessControlAllow
	POST    *accessControlAllow
	DELETE  *accessControlAllow
}

type accessControlAllow struct {
	Methods     string
	Headers     string
	Origins     map[string]bool
	Credentials bool
}

// cors Cross-Origin Resource Sharing
// (CORS) is a mechanism that uses additional HTTP headers to tell a browser
// to let a web application running at one origin (domain) have permission to access selected resources
// from a server at a different origin.
func (a *accessControlAllow) cors(c *gin.Context) error {
	c.Writer.Header().Set("Access-Control-Allow-Methods", a.Methods)
	c.Writer.Header().Set("Access-Control-Allow-Headers", a.Headers)

	if _, ok := a.Origins["*"]; ok {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	} else {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			_, ok := a.Origins[origin]
			if ok {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			} else {
				return errors.New("Origin " + origin + " is not allow")
			}
		}
	}

	if a.Credentials {
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	return nil
}

// sliceToStringWithDeduplication convert to a deduplicated string
func (a *accessControlAllow) toDeduplicatedString(m map[string]bool) string {
	var newStr string
	i := 0
	for str := range m {
		if i != 0 {
			newStr += ","
		} else {
			i++
		}
		newStr += str
	}
	return newStr
}

func (a *accessControlAllow) setMethods(methods map[string]bool) {
	a.Methods = a.toDeduplicatedString(methods)
}

func (a *accessControlAllow) setHeaders(headers map[string]bool) {
	a.Headers = a.toDeduplicatedString(headers)
}

func (a *accessControlAllow) setOrigin(origins map[string]bool) {
	a.Origins = origins
}
