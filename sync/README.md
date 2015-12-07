介绍
====

这个包用来在开发调试期，帮助排查程序中的死锁情况。

用法
====

通常我们项目中引用到原生`sync`包的代码会像这样：

```go
package myapp

import "sync"

var MyLock sync.Mutex

func MyFunc() {
	MyLock.Lock()
	defer MyLock.Unlock()

	// .......
}
```

只需要将原来引用`sync`的代码改为引用`github.com/funny/sync`包，不需要修改别的代码：


```go
package myapp

import "github.com/funny/sync"

var MyLock sync.Mutex

func MyFunc() {
	MyLock.Lock()
	defer MyLock.Unlock()

	// .......
}
```

这时候死锁诊断还没有被启用，因为做了条件编译，所以锁的开销跟原生`sync`包是一样的。

当需要编译一个带死锁诊断的版本的时候，在`go build --tags`列表中加入`deadlock`标签。

例如这样：

```
go build -tags deadlock myproject
```

同样这个标签也用于单元测试，否则默认的单元测试会死锁：

```
go test -tags deadlock -v
```


原理
====

在开启死锁检查的时候，系统会维护一份全局的锁等待列表，其次每个锁都会有当前使用者的信息。

当一个goroutine要等待一个锁的时候，系统会到全局的等待列表里面查找当前这个锁的使用者，是否间接或直接的正在等待当前请求锁的这个goroutine。

死锁不一定只发生在两个goroutine之间，极端情况也可能是一个链条状的依赖关系，又或者可能出现自身重复加锁的死锁情况。

当出现死锁的时候，系统将提取死锁链上的所有goroutine的堆栈跟踪信息，方便排查故障原因。

因为需要维护一份全局的锁等待列表，所以这里会出现额外并且集中的一个全局锁开销，会导致明显的程序的并发性能下降。

全局锁的问题还会再继续研究和加以改进，但是目前这个包是不能用于生产环境的，只能用在开发和调试期作为死锁诊断的辅助工具。
