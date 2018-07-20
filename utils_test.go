package ehttp

import "testing"

func TestGinPathToSwaggerPath(t *testing.T) {
	//  the gin router , /src/*filepath not supported
	ginRoutes := [...]string{
		"/",
		"/cmd/:tool/:sub",
		"/cmd/:tool/",
		"/search/",
		"/search/:query",
		"/user_:name",
		"/user_:name/about",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/info/:user/public",
		"/info/:user/project/:project",
	}
	swaggerRoutes := [...]string{
		"/",
		"/cmd/{tool}/{sub}",
		"/cmd/{tool}/",
		"/search/",
		"/search/{query}",
		"/user_{name}",
		"/user_{name}/about",
		"/doc/",
		"/doc/go_faq.html",
		"/doc/go1.html",
		"/info/{user}/public",
		"/info/{user}/project/{project}",
	}

	for i := 0; i < len(ginRoutes); i++ {
		path, err := ginPathToSwaggerPath(ginRoutes[i])
		if err != nil {
			testError(t, err)
		} else {
			if path != swaggerRoutes[i] {
				testError(t, path+" != "+swaggerRoutes[i])
			}
		}
		path, err = swaggerPathToGinPath(swaggerRoutes[i])
		if err != nil {
			testError(t, err)
		} else {
			if path != ginRoutes[i] {
				testError(t, path+" != "+ginRoutes[i])
			}
		}
	}

	ginRoutesNotSupport := [...]string{
		"/src/*filepath",
		"/usr/:id/*key",
	}
	for i := 0; i < len(ginRoutesNotSupport); i++ {
		_, err := ginPathToSwaggerPath(ginRoutesNotSupport[i])
		if err == nil {
			testError(t, "the err should not be nil")
		}
	}
}

func TestGetSwaggerTagFormPath(t *testing.T) {
	tests := []struct {
		Path string
		Want string
	}{
		{"/", ""},
		{"/cmd/:tool/:sub", "cmd"},
		{"/cmd/:tool/", "cmd"},
		{"/search/", "search"},
		{"/search/:query", "search"},
		{"/user_:name", "user_"},
		{"/user_:name/about", "user_"},
		{"/doc/", "doc"},
		{"/doc/go_faq.html", "doc"},
		{"/doc/go1.html", "doc"},
		{"/info/{user}/public", "info"},
		{"/info/{user}/project/{project}", "info"},
	}

	for _, test := range tests {
		tag := getSwaggerTagFormPath(test.Path)
		if tag != test.Want {
			testError(t, tag, " != ", test.Want)
		}
	}
}
