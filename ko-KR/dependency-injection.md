# 의존성 주입

**[이 챕터의 모든 코드는 여기에서 확인할 수 있다.](https://github.com/quii/learn-go-with-tests/tree/main/di)**

의존성 주입을 위해서는 인터페이스에 대한 이해가 필요하므로 이전에 structs 섹션을 읽었다고 가정한다.

프로그래밍 커뮤니티에는 의존성 주입과 관련해서  _많은_  오해가 있다. 이 가이드는 당신에게 어떻게 아래의 항목들이 가능한지 알려줄 것이다.

* 프레임워크가 필요하지 않다.
* 디자인을 지나치게 복잡하게 하지 않는다.
* 테스트를 용이하게 한다.
* 그것이 훌륭한 범용 함수를 작성하게 할 것이다.

Hello-world 장에서 했던 것처럼 누군가를 맞이하는 함수를 작성하고 싶지만 이번에는 _실제 print_를 테스트할 것이다.

요약하자면 다음은 그 함수의 모습이다.

```go
func Greet(name string) {
	fmt.Printf("Hello, %s", name)
}
```

그러나 이것을 어떻게 테스트할 수 있는가? `fmt.Printf`를 호출하면 stdout으로 인쇄된다. 이는 테스트 프레임 워크를 사용하여 캡처하기가 매우 어렵다.

우리가 해야할 일은 print 하는 것의 의존성을 **주입** \(인자를 넘기는 것을 그냥 fancy 하게 표현한 것이다\) 할 수 있도록 한다.

**우리의 함수는 **_**어디에서**_** 또는 **_**어떻게**_** print가 발생하는 지를 신경 쓸 필요가 없다. 그래서 우리는 구체적인 type보다는 **_**interface**_** type을 허용해야 한다.**

그렇게 한다면, 우리가 제어하는 어떤 것으로 출력하도록 구현을 변경하여 테스트할 수 있다. "실생활"에서는 stdout에 쓰는 것을 주입한다.

`fmt.Printf`의 소스 코드를 보면 연결하는 방법을 알 수 있다.

```go
// 쓰여진 바이트 수와 발생한 error를 반환합니다.
func Printf(format string, a ...interface{}) (n int, err error) {
	return Fprintf(os.Stdout, format, a...)
}
```

흥미롭다! 내부적으로 `Printf`는 `os.Stdout`을 전달하는`Fprintf`를 단지 호출한다.

`os.Stdout`은 정확히 무엇인가? `Fprintf`는 첫 번째 argument로 무엇을 전달받을 것으로 예상하는가?

```go
func Fprintf(w io.Writer, format string, a ...interface{}) (n int, err error) {
	p := newPrinter()
	p.doPrintf(format, a)
	n, err = w.Write(p.buf)
	p.free()
	return
}
```

`io.Writer`이다.

```go
type Writer interface {
	Write(p []byte) (n int, err error)
}
```

당신이 Go 코드를 더 많이 작성하면 이 interface를 많이 보게 될 것이다. 왜냐하면 "data를 어딘가에 넣는 것"을 잘 표현하는 좋은 general purpose interface이기 때문이다.

그래서 우리는 어딘가에 인사말을 보내기 위해 궁극적으로 `Writer`를 사용하고 있다는 것을 안다. 기존 추상화를 사용하여 코드를 테스트할 수 있고 더 재사용 가능하게 만들어 보자.

## 먼저, 테스트하기

```go
func TestGreet(t *testing.T) {
	buffer := bytes.Buffer{}
	Greet(&buffer, "Chris")

	got := buffer.String()
	want := "Hello, Chris"

	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

`bytes` 패키지의 `buffer` type은 `Writer` 인터페이스를 구현한다.

그래서 우리는 테스트에서 buffer를 `Writer`로서 사용할 것이다. 그리고 `Greet` 호출한 후에, 무엇이 buffer에 쓰여졌는 지 확인할 수 있다.

## 테스트 실행해보기

테스트는 컴파일되지 않는다.

```text
./di_test.go:10:7: too many arguments in call to Greet
    have (*bytes.Buffer, string)
    want (string)
```

## 테스트 실행을 위해 최소한의 코드를 작성하고 실패한 테스트 출력을 확인하자.

_컴파일러의 말을 듣고_ 문제를 해결하자.

```go
func Greet(writer *bytes.Buffer, name string) {
	fmt.Printf("Hello, %s", name)
}
```

`Hello, Chris di_test.go:16: got '' want 'Hello, Chris'`

테스트에 실패했다. 이름이 인쇄되었지만 stdout으로 간다.

## 테스트를 통과하는 최소한의 코드 작성하기

테스트에서 writer를 사용하여 버퍼에 인사말을 보낸다. `fmt.Fprintf`는`fmt.Printf`와 비슷하지만, 문자열을 보낼 곳 `Writer`를 가진다. 반면, `fmt.Printf`는 기본적으로 stdout을 사용한다.

```go
func Greet(writer *bytes.Buffer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}
```

테스트는 이제 통과된다.

## 리팩터링 하기

이전에 컴파일러는`bytes.Buffer`에 대한 포인터를 전달하라고 했다. 이것은 기술적으로 정확하지만, 그다지 유용하지는 않다.

이를 증명하기 위해 `Greet` 함수를 표준 출력으로 인쇄하려는 Go 애플리케이션에 연결해보자.

```go
func main() {
	Greet(os.Stdout, "Elodie")
}
```

`./di.go:14:7: cannot use os.Stdout (type *os.File) as type *bytes.Buffer in argument to Greet`

앞서 논의했듯이`fmt.Fprintf`는 `io.Writer`로 그것을 구현하는 `os.Stdout`와 `bytes.Buffer`를 전달할 수 있도록 한다.

좀 더 범용적인 인터페이스를 사용하도록 코드를 변경하면 이제 테스트와 애플리케이션 모두에서 사용할 수 있다.

```go
package main

import (
    "fmt"
    "os"
    "io"
)

func Greet(writer io.Writer, name string) {
    fmt.Fprintf(writer, "Hello, %s", name)
}

func main() {
	Greet(os.Stdout, "Elodie")
}
```

## io.Writer 에 대해 더 알아보기

`io.Writer`를 사용하여 데이터를 쓸 수 있는 다른 곳이 있는가? 우리의 `Greet` 함수는 얼마나 general purpose 인가?

### The internet

다음을 실행해라.

```go
package main

import (
	"fmt"
	"io"
	"net/http"
)

func Greet(writer io.Writer, name string) {
	fmt.Fprintf(writer, "Hello, %s", name)
}

func MyGreeterHandler(w http.ResponseWriter, r *http.Request) {
	Greet(w, "world")
}

func main() {
	http.ListenAndServe(":5000", http.HandlerFunc(MyGreeterHandler))
}
```

프로그램을 실행시키고 [http://localhost:5000](http://localhost:5000) 주소로 가라. Greet 함수가 사용되는 것을 볼 수 있다.

HTTP 서버는 이후 장에서 다룰 것이므로 세부 사항에 대해 너무 걱정하지 마라.

HTTP 핸들러를 작성할 때 요청을 위해 만들어진 `http.ResponseWriter`와`http.Request`를 받습니다. 서버를 구현할 때 writer를 사용하여 응답을 _write_한다.

당신은 `http.ResponseWriter` 또한 `io.Writer`를 구현한다고 추측할 것이고 그것이 우리가 handler 내에서 `Greet` 함수를 재사용할 수 있는 이유이다.

## 정리

첫 번째 코드는 제어 할 수 없는 곳에 데이터를 기록했기 때문에 테스트하기가 쉽지 않았다.

_테스트에 의해 동기 부여받았다._ 코드를 리팩토링하여 _어느 곳에_ 데이터가 쓰여질 지를 **종속성 주입**을 통해 제어할 수 있게 되었다.

* **코드를 테스트해라** 함수를 _쉽게_ 테스트할 수 없다면, 이는 일반적으로 함수 _또는_ 전역 상태에 연결된 종속성 때문이다. 예를 들어 서비스 계층에서 사용되는 글로벌 데이터베이스 연결 풀이 있는 경우 테스트하기가 어려울 수 있으며 실행 속도가 느려진다. DI는 당신이 데이터베이스 의존성에 \(인터페이스를 통해\) 주입하도록 동기를 부여한다. 그런 다음 테스트에서 제어할 수 있는 무언가로 mock 할 수 있다.
* **관심사를 분리해라**, _데이터가 가는 목적지_를 _데이터를 어떻게 만들 것인가_와 분리한다. 메소드/함수가 너무 많은 책임이 있다고 느낀 적이 있다면 \(데이터 생성 _그리고_ db에 쓰기? HTTP 요청 처리 _그리고_ 도메인 레벨 로직 수행? \) DI가 아마도 필요한 도구일 것이다.
* **다른 컨텍스트에서 코드를 재사용 할 수 있도록 허용하라** 우리 코드를 사용할 수 있는 첫 번째 "새로운" 컨텍스트는 테스트 내부이다. 그러나 향후 누군가가 당신의 함수로 새로운 것을 시도하고 싶다면 그들 자신의 의존성을 주입 할 수 있다.

### Mocking은 어떠한가? 당신이 DI를 위해 필요하며 또한 그것은 악마라고 들었다.

Mocking은 나중에 자세히 다룰 것이다 \(그리고 악마가 아니다\). Mocking을 사용하면 당신이 주입하는 실제 내용을 가짜 버전으로 대체해서 테스트에서 당신이 제어하고 조사할 수 있도록 한다. 우리의 경우에는 표준 라이브러리에 사용 가능하도록 어떤 것이 준비되어 있다.

### Go 표준 라이브러리는 정말 훌륭하다. 그것을 공부해라.

`io.Writer` 인터페이스에 익숙해지면 테스트에서 `bytes.Buffer`를 `Writer`로 사용할 수 있으며 명령 줄 앱 또는 웹 서버에서 표준 라이브러리의 다른 `Writer`를 사용하여 함수를 사용할 수 있다.

표준 라이브러리와 친해질수록 이러한 범용 인터페이스를 더 많이 볼 수 있다. 그러면 당신의 코드에서 그것을 재사용하여 당신의 소프트웨어를 많은 맥락에서 재사용 가능하게 만들 수 있다.

이 예제는 [Go 프로그래밍 언어] (https://www.amazon.co.uk/Programming-Language-Addison-Wesley-Professional-Computing/dp/0134190440)의 장에서 크게 영향을 받았다. 가서 사라!
