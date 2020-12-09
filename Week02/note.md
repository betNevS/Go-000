# Week02-Go语言实践之error

### Error

Go error 就是一个普通的接口，普通的值。

```go
type error interface {
	Error() string
}	
```

在开发中，经常使用`errors.New()`来返回一个error，特别简单。

```go
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

func New(text string) error {
	return &errorString{text}
}
```

Go语言中有panic机制，但是panic和其他语言的exception是完全不一样的，`Go panic` 意味着挂掉，不能假设让调用者来解决panic，一旦使用panic就意味着代码不能运行了。对于真正意外的情况，那些表示不可恢复的程序错误，例如索引越界、不可恢复的环境问题、栈溢出，我们才使用 panic。对于其他的错误情况，我们应该是期望使用 error 来进行判定。

>特别在http服务中，一定不要在请求中处理业务逻辑时抛出panic，一旦抛出panic，整个http服务就挂了。

### Sentinel Error

预定义的特定错误，叫做`sentinel error`，在基础类库中大量使用，使用sentinel error是最不灵活的错误处理策略，因为调用方比较使用==将结果与预先声明的值进行比较。当想要加入更多上下文到错误中时，则相当于新增了一个error，就会造成返回不是同一个error，破坏相等性的检查。

```go
// ErrShortWrite means that a write accepted fewer bytes than requested
// but failed to return an explicit error.
var ErrShortWrite = errors.New("short write")

// ErrShortBuffer means that a read required a longer buffer than was provided.
var ErrShortBuffer = errors.New("short buffer")

// EOF is the error returned by Read when no more input is available.
// Functions should return EOF only to signal a graceful end of input.
// If the EOF occurs unexpectedly in a structured data stream,
// the appropriate error is either ErrUnexpectedEOF or some other error
// giving more detail.
var EOF = errors.New("EOF")
```

sentinel error具有以下两个问题：

* Sentinel error 成为API 公共部分。比如 io.Reader，像 io.Copy 这类函数需要 reader 的实现者比如返回 io.EOF 来告诉调用者没有更多数据了，但这又不是错误。
* Sentinel error 在两个包之间创建了依赖。

结论：开发业务逻辑时尽可能避免 sentinel error，虽然标准库中使用了它，但并不是一个效仿的模式。

### Error Type

Error type 是实现了error接口的自定义类型，例如

```go
type MyError struct {
	Msg string
	File string
	Line int
}

func (e *MyError) Error() string {
	return fmt.Sprintf("%s:%d:%s", e.File, e.Line, e.Msg)
}
```

因为 MyError 是一个 type，所以调用者可以使用断言来换成这个类型 `err.(type)`，虽然 error type 可以携带更多的上下文信息，但是这种模型会导致和调用者产生强耦合， 从而使的API变得脆弱。

### Opaque Error

不透明错误处理是指因为虽然您知道发生了错误，但您没有能力看到错误的内部。作为调用者，关于操作的结果，您所知道的就是它起作用了，或者没有起作用(成功还是失败)。这就是不透明错误处理的全部功能–只需返回错误而不假设其内容。

```go
x, err := bar.Foo()
if err != nil {
	return err
}
```

但是在一些情况下，这种二分处理错误的方式是不够的，有时是需要知道具体的错误场景然后做其他的操作，这种情况下断言错误实现了特定的行为，而不是断言错误是特定的类型或值。

```go
type temporary interface {
	Temporary() bool
}

func IsTemporary(err error) bool {
	te, ok := err.(temporary)
	return ok && te.Temporary()
}
```

> 这里的关键是，这个逻辑可以在不导入定义错误的包或者实际上不了解 err 的底层类型的情况下实现——我们只对它的行为感兴趣。

### Wrap Error

> You should only handle errors once. Handling an error means inspecting the error value, and making a single decision.

建议使用 `github.com/pkg/errors` 来 wrap error.

```go
func testWrap() error {
	err := getError()
	return errors.Wrap(err, "add error context")
}
```

使用场景：

* 在应用代码中，使用 `github.com/pkg/errors` 中的 errors.New() 返回错误。
* 如果调用其他包内的函数，直接简单返回。
* 如果和基础库进行协作，考虑使用 errors.Wrap 考虑保存堆栈信息。
* 直接返回错误，而不是每个错误的地方导出打日志。
* 在程序顶部使用 `%+v` 记录错误。
* 使用 errors.Cause 获取 root error，再和 sentinel error 判定。

总结：

* Packages that are reusable across many projects only return root error values. 
* If the error is not going to be handled, wrap and return up the call stack.
* Once an error is handled, it is not allowed to be passed up the call stack any longer.

go1.13为 errors 和 fmt 标准库包中引入了新特性，其中最重要的一条是：包含另一个错误的error可以实现返回底层错误的Unwrap方法。如果e1.Unwrap()返回e2，那么我们说e1包装了e2。

```go
旧：if err == ErrNotFound {...}
新：if errors.Is(err, ErrNotFound) {...}

旧：if e, ok := err.(*QueryError);ok {...}
新：if errors.As(err, &e) {...}
```

在Go 1.13中 fmt.Errorf 新增了`%w`，用`%w`包装错误可用于 `errors.Is` 以及 `errors.As`：

```go
err := fmt.Errorf("access denied: %w", ErrPermission)
if errors.Is(err, ErrPermission) ...	
```

### Error & github.com/pkg/errors

在Go1.13之后，error已经挺好用了，但是依旧有一个问题，就是没有带上堆栈信息，出现问题不好定位到具体的文件，因此对于Go语言中error的实践还是建议 `github.com/pkg/errors` ，这个库已经兼容了标准库的error。

```go
package main

import (
	"errors"
	"fmt"

	xerrors "github.com/pkg/errors"
)

var ErrMy = errors.New("my error")

func main() {
	err := test2()
	fmt.Printf("%+v\n", err)
	fmt.Println(errors.Is(err, ErrMy))

}

func test2() error {
	return test1()
}

func test1() error {
	return test0()
}

func test0() error {
	return xerrors.Wrapf(ErrMy, "wrap %s", "dd")
}

```

### [作业](README.md)

### [课堂问题](https://shimo.im/docs/R6gP8qyvWqJrgRCk)



