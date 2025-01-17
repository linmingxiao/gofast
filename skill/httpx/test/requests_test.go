package test

import (
	"github.com/qinchende/gofast/skill/httpx"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseForm(t *testing.T) {
	var v struct {
		Name    string  `form:"name"`
		Age     int     `form:"age"`
		Percent float64 `form:"percent,NA"`
	}

	r, err := http.NewRequest(http.MethodGet, "http://hello.com/a?name=hello&age=18&percent=3.4", nil)
	assert.Nil(t, err)
	assert.Nil(t, httpx.Parse(r, &v))
	assert.Equal(t, "hello", v.Name)
	assert.Equal(t, 18, v.Age)
	assert.Equal(t, 3.4, v.Percent)
}

func TestParseHeader(t *testing.T) {
	m := httpx.ParseHeader("key=value;")
	assert.EqualValues(t, map[string]string{
		"key": "value",
	}, m)
}

func TestParseFormOutOfRange(t *testing.T) {
	var v struct {
		Age int `form:"age,range=[10:20)"`
	}

	tests := []struct {
		url  string
		pass bool
	}{
		{
			url:  "http://hello.com/a?age=5",
			pass: false,
		},
		{
			url:  "http://hello.com/a?age=10",
			pass: true,
		},
		{
			url:  "http://hello.com/a?age=15",
			pass: true,
		},
		{
			url:  "http://hello.com/a?age=20",
			pass: false,
		},
		{
			url:  "http://hello.com/a?age=28",
			pass: false,
		},
	}

	for _, test := range tests {
		r, err := http.NewRequest(http.MethodGet, test.url, nil)
		assert.Nil(t, err)

		err = httpx.Parse(r, &v)
		if test.pass {
			assert.Nil(t, err)
		} else {
			assert.NotNil(t, err)
		}
	}
}

func TestParseMultipartForm(t *testing.T) {
	var v struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}

	body := strings.Replace(`----------------------------220477612388154780019383
Content-Disposition: form-data; name="name"

kevin
----------------------------220477612388154780019383
Content-Disposition: form-data; name="age"

18
----------------------------220477612388154780019383--`, "\n", "\r\n", -1)

	r := httptest.NewRequest(http.MethodPost, "http://localhost:3333/", strings.NewReader(body))
	r.Header.Set(httpx.ContentType, "multipart/form-data; boundary=--------------------------220477612388154780019383")

	assert.Nil(t, httpx.Parse(r, &v))
	assert.Equal(t, "kevin", v.Name)
	assert.Equal(t, 18, v.Age)
}

func TestParseMultipartFormWrongBoundary(t *testing.T) {
	var v struct {
		Name string `form:"name"`
		Age  int    `form:"age"`
	}

	body := strings.Replace(`----------------------------22047761238815478001938
Content-Disposition: form-data; name="name"

kevin
----------------------------22047761238815478001938
Content-Disposition: form-data; name="age"

18
----------------------------22047761238815478001938--`, "\n", "\r\n", -1)

	r := httptest.NewRequest(http.MethodPost, "http://localhost:3333/", strings.NewReader(body))
	r.Header.Set(httpx.ContentType, "multipart/form-data; boundary=--------------------------220477612388154780019383")

	assert.NotNil(t, httpx.Parse(r, &v))
}

func TestParseJsonBody(t *testing.T) {
	var v struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	body := `{"name":"kevin", "age": 18}`
	r := httptest.NewRequest(http.MethodPost, "http://localhost:3333/", strings.NewReader(body))
	r.Header.Set(httpx.ContentType, httpx.ApplicationJson)

	assert.Nil(t, httpx.Parse(r, &v))
	assert.Equal(t, "kevin", v.Name)
	assert.Equal(t, 18, v.Age)
}

func TestParseRequired(t *testing.T) {
	v := struct {
		Name    string  `form:"name"`
		Percent float64 `form:"percent"`
	}{}

	r, err := http.NewRequest(http.MethodGet, "http://hello.com/a?name=hello", nil)
	assert.Nil(t, err)
	assert.NotNil(t, httpx.Parse(r, &v))
}

func TestParseOptions(t *testing.T) {
	v := struct {
		Position int8 `form:"pos,enum=1|2"`
	}{}

	r, err := http.NewRequest(http.MethodGet, "http://hello.com/a?pos=4", nil)
	assert.Nil(t, err)
	assert.NotNil(t, httpx.Parse(r, &v))
}

func BenchmarkParseRaw(b *testing.B) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/a?name=hello&age=18&percent=3.4", nil)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		v := struct {
			Name    string  `form:"name"`
			Age     int     `form:"age"`
			Percent float64 `form:"percent,NA"`
		}{}

		v.Name = r.FormValue("name")
		v.Age, err = strconv.Atoi(r.FormValue("age"))
		if err != nil {
			b.Fatal(err)
		}
		v.Percent, err = strconv.ParseFloat(r.FormValue("percent"), 64)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseAuto(b *testing.B) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/a?name=hello&age=18&percent=3.4", nil)
	if err != nil {
		b.Fatal(err)
	}

	for i := 0; i < b.N; i++ {
		v := struct {
			Name    string  `form:"name"`
			Age     int     `form:"age"`
			Percent float64 `form:"percent,NA"`
		}{}

		if err = httpx.Parse(r, &v); err != nil {
			b.Fatal(err)
		}
	}
}
