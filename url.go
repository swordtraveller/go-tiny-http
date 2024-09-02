package http

import "strings"

type Values map[string][]string

type URL struct {
	Path     string
	RawQuery string
}

func (u *URL) Query() Values {
	v, _ := ParseQuery(u.RawQuery)
	return v
}

func ParseQuery(input string) (Values, error) {
	values := make(Values)
	terms := strings.Split(input, "&")
	for _, term := range terms {
		keyAndValue := strings.Split(term, "=")
		if len(keyAndValue) > 1 {
			tmp := make([]string, 0)
			tmp = append(tmp, keyAndValue[1])
			values[keyAndValue[0]] = tmp
		}
	}
	return values, nil
}

func (v Values) Get(key string) string {
	value := v[key]
	if len(value) == 0 {
		return ""
	}
	return value[0]
}
