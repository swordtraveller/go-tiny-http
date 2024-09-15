package http

import (
	"io"
	"strings"
)

type Request struct {
	Method     string
	URL        *URL
	Proto      string
	Header     map[string][]string
	Body       ReadCloser
	Form       map[string][]string
	RequestURI string
}

type ReadCloser struct {
	value    []byte
	progress int
}

func (rc ReadCloser) Read(p []byte) (n int, err error) {
	curProgress := 0
	for index, _ := range p {
		if rc.progress == len(rc.value) {
			return curProgress, io.EOF
		} else {
			p[index] = rc.value[rc.progress]
			rc.progress++
			curProgress++
		}
	}
	return curProgress, nil
}

func (r *Request) FormValue(key string) string {
	// 懒计算 / 惰性求值
	// lazy evaluation
	if r.Form == nil {
		r.formValue()
	}
	value, ok := r.Form[key]
	if ok {
		return value[0]
	}
	return ""
}

func (r *Request) formValue() {
	r.Form = getKeyValueFromForm(r.URL.RawQuery)
}

func getKeyValueFromForm(input string) map[string][]string {
	result := make(map[string][]string)
	terms := strings.Split(input, "&")
	for _, term := range terms {
		keyAndValue := strings.Split(term, "=")
		if len(keyAndValue) > 1 {
			tmp := make([]string, 0)
			tmp = append(tmp, keyAndValue[1])
			result[keyAndValue[0]] = tmp
		}
	}
	return result
}
