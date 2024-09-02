package http

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	HOSTNAME    = "0.0.0.0"
	PROTOCOL    = "tcp"
	VERSION_1_1 = "HTTP/1.1"
)

type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}

var router = map[string]func(ResponseWriter, *Request){}

func HandleFunc(pattern string, handler func(ResponseWriter, *Request)) {
	router[pattern] = handler
}

func ListenAndServe(addr string, srvHandler Handler) error {

	// Create a server
	// 创建服务器
	listener, err := net.Listen(PROTOCOL, HOSTNAME+addr)
	if err != nil {
		return err
	}
	defer listener.Close()

	for {

		// Establish a connection
		// 建链
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		// Create a goroutine to handle the connection
		// 协程处理连接
		go func() {
			defer func() {
			}()
			handleConn(conn, srvHandler)
		}()

	}

	// Blocking
	// 阻塞
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
	return nil

}

func handleConn(conn net.Conn, srvHandler Handler) {

	reader := bufio.NewReader(conn)

	var req Request

	// request line
	// 请求行
	reqLine, err := reader.ReadString('\n')
	if err == io.EOF {
		// Use return, don't use panic-recover
		// The performance of return is better than panic, and here my recover doesn't do any meaningful work
		return
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
	reqLineFields := strings.Split(reqLine, " ")
	if len(reqLineFields) < 3 {
		panic("request line fields are too few")
	}
	req.Method = reqLineFields[0]
	PathAndQuery := strings.Split(reqLineFields[1], "?")
	if len(PathAndQuery) > 0 {
		req.URL.Path = PathAndQuery[0]
		if len(PathAndQuery) > 1 {
			req.URL.RawQuery = PathAndQuery[1]
		}
	} else {
		req.URL.Path = reqLineFields[1]
	}
	req.Proto = reqLineFields[2]

	// request headers
	// 请求头
	bodySize := 0
	req.Header = make(map[string][]string)
	for {
		HeaderLine, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		// "\r\n" marks the end of the request headers section
		// 读取到"\r\n"时结束
		if HeaderLine == "\r\n" {
			break
		}
		HeaderLine = strings.TrimRight(HeaderLine, "\r\n")
		key, value := getKeyValue(HeaderLine)
		req.Header[key] = append(req.Header[key], value)
		// Retrieve the size of the request body from the request header named "Content-Length"
		// The keys of the request header are case-insensitive
		// 从"Content-Length"请求头中获取请求体大小
		// 请求头的键名是不区分大小写的
		if strings.ToLower(key) == "content-length" {
			bodySize, err = strconv.Atoi(strings.TrimRight(value, "\r\n"))
			if err != nil {
				panic(err)
			}
		}
	}

	// request body
	// 请求体
	if bodySize > 0 {
		p := make([]byte, bodySize)
		_, err := reader.Read(p)
		if err != nil {
			panic(err)
		}
		req.Body.value = p
	}

	// response
	// 响应
	var rw responseWriter
	rw.SetStatusCode(200)

	// srvHandler != nil when http.ListenAndServe(":1234", handler) if handler is not nil
	// srvHandler is nil when http.ListenAndServe(":1234", nil)
	// srvHandler 即用户调用http.ListenAndServe(addr, handler)时传入的第二个参数
	if srvHandler != nil {
		srvHandler.ServeHTTP(&rw, &req)
	} else {
		handler, ok := router[req.URL.Path]
		if !ok {
			panic(fmt.Sprintf("Route %s does not exist! 路由%s不存在！", req.URL.Path, req.URL.Path))
		}
		handler(&rw, &req)
	}

	// respLine := "HTTP/1.1 200 OK\r\n"
	respLine := fmt.Sprintf("%s %d %s\r\n", VERSION_1_1, rw.StatusCode, MessageMap[rw.StatusCode])
	respHeaders := ""
	header := rw.Header()
	hasContentLengthHeader := false
	for k, v := range header {
		if k == "Content-Length" {
			hasContentLengthHeader = true
		}
		respHeaders = respHeaders + fmt.Sprintf("%s: %s\r\n", k, v)
	}
	if !hasContentLengthHeader {
		respHeaders = respHeaders + fmt.Sprintf("Content-Length: %d\r\n", rw.ContentLength)
	}
	resp := respLine + respHeaders + "\r\n" + rw.ResponseBody
	conn.Write([]byte(resp))
	conn.Close()
}

func getKeyValue(input string) (string, string) {

	// state
	// 状态
	const (
		ReadingKey = iota
		EatingSpace
		ReadingValue
	)

	state := ReadingKey
	raw := []byte(input)
	var key string
	var value string
	index := -1
Loop:
	for _, b := range raw {
		index++
		switch state {
		case ReadingKey:
			if b == ':' {
				state = EatingSpace
				continue
			} else {
				key = key + string(b)
			}
		case EatingSpace:
			if b == ' ' {
				continue
			} else {
				state = ReadingValue
				break Loop
			}
		case ReadingValue:
			break Loop
		}
	}
	value = input[index:]
	return key, value

}
