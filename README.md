# go-tiny-http

A lightweight http library with concise implementation, which can seamlessly replace native libraries in most cases.  
一个轻量级的http库，实现简明，在大部分情况下可以无缝替换原生库。  

Used to demonstrate the principles of the HTTP protocol. Learn tricks by comparing with the standard library.  
用于演示HTTP协议的原理。通过和标准库对比来学习技巧。  

## Getting Started 起步

### Getting go-tiny-http 获取go-tiny-http

Using the net/http library, we can create a demo like this:  
使用net/http库，我们可以创建这样一个例子：  

```go
package main

import (
	"fmt"
	"io"

	"net/http"
)

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":1234", nil)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World!")
	b, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	body := string(b)
	fmt.Println(body)
}
```

Then we test:  
然后我们进行测试：  

On Windows:  
在Windows上：  

```powershell
Invoke-RestMethod -Uri "127.0.0.1:1234/hello" -Method POST -Body "wow!"
```

Or on Linux:  
或者在Linux上：  

```bash
curl -XPOST -d "wow!" "127.0.0.1:1234/hello"
```

The HTTP response is as follows:  
HTTP响应如下：  
```
Hello World!
```

Now we replace "net/http" with "github.com/swordtraveller/go-tiny-http":  
现在我们把"net/http"替换为"github.com/swordtraveller/go-tiny-http"：  

Install go-tiny-http:  
安装go-tiny-http：  

```powershell
go get github.com/swordtraveller/go-tiny-http
```

Replace like this:  
像下面这样替换：  
```
// "net/http"
"github.com/swordtraveller/go-tiny-http"
```

The new code is as follows:  
新的代码如下：  
```go
package main

import (
	"fmt"
	"io"

	"github.com/swordtraveller/go-tiny-http"
)

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":1234", nil)
}

func helloHandler(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "Hello World!")
	b, err := io.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	body := string(b)
	fmt.Println(body)
}

```

Tested again and the results were consistent with before.  
再次测试，结果与之前一致。  

## 自动化测试

我们可以在 Github 仓库的 `Actions` 板块中点击 `New workflow` ，选择一个模板，这里选用 Github 自动推荐的 "Go" 模板，点击 `Configure` ，即可添加新的自动化任务。

有几个值得注意的部分。  

首先是触发的时机，可以在 push 或者 PR 时触发：  

```yaml
on:
  push:
    branches: [ "master" ]
  pull_request:
    branches: [ "master" ]
```

如果想改为对所有分支生效，可以在 `branches` 字段填写 `"*"` 值。   

然后模板中有一个名为 "build" 的 job ，它有许多 steps （阶段），其中两个重要的是 "Build" 和 "Test" ，用来执行编译和测试步骤，在 `run` 字段里可以输入命令内容。  

```yaml
- name: Build
  run: go build -v ./...

- name: Test
  run: go test -v ./...
```
