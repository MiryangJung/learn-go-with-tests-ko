# 동시성

**[이 챕터에서 사용되는 모든 코드는 여기서 찾을 수 있다.](https://github.com/quii/learn-go-with-tests/tree/main/concurrency)**

생각해 보자: 동료가 URL 목록의 상태를 확인하는 기능인 `CheckWebsites` 함수를 작성했다.

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
  results := make(map[string]bool)

  for _, url := range urls {
    results[url] = wc(url)
  }

  return results
}
```

이 코드는 각 URL을 확인하여 map으로 구성된 boolean 값 - 올바른 응답에는 `true`, 잘못된 응답에는 `false` 을 반환한다.

당신은 또한 `WebsiteChecker`를 통과해야 한다. 해당 함수는 단일의 URL을 필요로 하고 boolean 값을 반환한다. 이 기능은 모든 웹 사이트를 확인하는 데 사용된다.

[의존성 주입(DI)][DI]을 통해 실제 HTTP 호출 없이 기능을 테스트할 수 있어 안정적이고 빠르게 사용할 수 있다.

이것이 그들이 썼던 테스트이다:

```go
package concurrency

import (
  "reflect"
  "testing"
)

func mockWebsiteChecker(url string) bool {
  if url == "waat://furhurterwe.geds" {
    return false
  }
  return true
}

func TestCheckWebsites(t *testing.T) {
  websites := []string{
    "http://google.com",
    "http://blog.gypsydave5.com",
    "waat://furhurterwe.geds",
  }

  want := map[string]bool{
    "http://google.com":          true,
    "http://blog.gypsydave5.com": true,
    "waat://furhurterwe.geds":    false,
  }

  got := CheckWebsites(mockWebsiteChecker, websites)

  if !reflect.DeepEqual(want, got) {
    t.Fatalf("Wanted %v, got %v", want, got)
  }
}
```

해당 함수는 생산 중이고 수백 개의 웹사이트들을 확인하는 데 사용되고 있다. 하지만 이 작업이 느리다고 당신의 동료의 불만이 쌓이기 시작한다. 그래서 그들은 이 기능의 속도를 높여달라고 요청했다.

## 테스트를 작성해 보자

변화에 대한 효과를 보기 위해 기준(benchmark)을 사용하여 `CheckWebsites`의 속도를 테스트해보겠다.

```go
package concurrency

import (
  "testing"
  "time"
)

func slowStubWebsiteChecker(_ string) bool {
  time.Sleep(20 * time.Millisecond)
  return true
}

func BenchmarkCheckWebsites(b *testing.B) {
  urls := make([]string, 100)
  for i := 0; i < len(urls); i++ {
    urls[i] = "a url"
  }

  for i := 0; i < b.N; i++ {
    CheckWebsites(slowStubWebsiteChecker, urls)
  }
}
```

벤치마크는 100개의 url들을 사용한 `CheckWebsites`와 가짜의 구현체를 사용한 `WebsiteChecker`를 테스트한다.
`slowStubWebsiteChecker`는 의도적으로 느리게 만들었다. 해당 코드는 `time.Sleep`를 사용하여 20ms를 기다렸다가 true 값을 반환한다.

벤치마크를 사용하려면 `go test -bench=.`를 실행하자. (혹은 만약 윈도우 PowerShell을 사용한다면 `go test -bench="."`이다.)

```sh
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v0
BenchmarkCheckWebsites-4               1        2249228637 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v0        2.268s
```

`CheckWebsites`가 2249228637 나노 초로 기준(benchmark)이 되었다 - 2와 1/4초이다.

이것을 좀 더 빨리 만들어 보자.

### 통과할 만큼 충분한 코드를 작성하자

드디어 우리는 이제 동시성에 대해 얘기할 수 있다. 이는 '한 번의 진행에 1개보다 더 많은 일을 하는 것'을 뜻한다. 그리고 이 일은 우리가 매일 자연스럽게 하고 있다.

예를 들어, 나는 아침에 차를 한 잔 만들었다. 나는 주전자를 올려놓고 물이 끓는 동안 냉장고에서 우유를 가지고 왔고, 찬장에서 차를 꺼내고 내가 좋아하는 머그잔을 찾았으며, 컵에 티백을 넣고 물이 다 끓었으면 컵에 물을 따랐다.

내가 _하지 않았던_ 것은 주전자를 올려놓고 거기에 서서 주전자의 물이 끓을 때까지 멍하니 바라보다가, 물이 다 끓으면 모든 일을 하는 것이다.

만약 당신이 첫 번째 방법으로 차를 만드는 방법이 왜 더 빠른지 이해한다면, `CheckWebsites`을 어떻게 더 빠르게 만들지 이해를 할 수 있을 것이다. 다음 웹 사이트에 요청을 보내기 전에 웹 사이트가 응답하기를 기다리는 대신에, 우리가 컴퓨터에게 대기하는 시간 동안 다음 요청을 하도록 만들어 보겠다.

보통 Go에서는 `doSomething()`이라는 함수를 호출했을 때 반환이 될 때까지 기다려야 한다(반환 값이 없다고 하더라도 함수가 끝날 때까지 기다린다). 우리는 이러한 연산을 *동기(blocking)* - 이것은 우리가 끝날 때까지 기다리도록 만든다라고 말한다. 동기적으로 실행되지 않는 연산은 *goroutine*이라고 하는 별도의 프로세스에서 실행된다. Go 코드를 상단부터 아래로 읽어 내려가는 동작을 생각하면, 각 함수를 만날 때마다 코드의 '내부'로 들어가 무슨 기능을 하는지 읽게 된다. 별도의 프로세스가 시작되면 원래 읽던 사람과는 다르게 다른 읽는 사람이 함수 내부를 읽어 내려가는 것과 같다.

Go에게 새로운 goroutine으로 실행하라고 말하기 위해서는 키워드 `go`를 함수 앞에 붙이는 방법: `go doSomething()`으로 함수 호출을 `go` statement로 바꿀 수 있다.

```go
package concurrency

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
  results := make(map[string]bool)

  for _, url := range urls {
    go func() {
      results[url] = wc(url)
    }()
  }

  return results
}
```

goroutine을 시작하는 유일한 방법은 `go`를 함수 호출 앞에 붙이는 것이기 때문에, goroutine을 시작하기 위해 종종 *익명 함수*를 사용하기도 한다. 익명 함수는 정규 함수 선언과 동일하게 보이지만 이름이 없다(당연하다). 위에 적힌 코드의 `for` 반복문 내의 몸체 부분에서 볼 수 있다.

익명 함수에는 유용하게 사용할 수 있는 여러 가지 기능들이 있는데, 이 중 2가지는 위에 사용을 했다. 첫 번째로, 선언된 것과 동시에 실행될 수 있다 - 그래서 익명 함수의 끝에 `()`이 붙어있는 것이다. 두 번째로는 정의된 곳에서의 lexical scope에 대한 접근을 유지한다는 것이다 - 익명 함수를 선언할 때 사용할 수 있는 모든 변수들을 함수 본문에서도 사용할 수 있다.

위에 있는 익명 함수의 몸체 부분은 이전 반복 문의 몸체 부분과 동일하다. 유일한 차이점은 각 반복이 새로운 goroutine으로 시작이 되고, 현재의 프로세스(`WebsiteChecker` 함수)와 동시적으로 실행되어 각 결과를 결과 map에 추가한다는 것이다.

하지만 우리가 `go test`로 실행을 하면:

```sh
--- FAIL: TestCheckWebsites (0.00s)
        CheckWebsites_test.go:31: Wanted map[http://google.com:true http://blog.gypsydave5.com:true waat://furhurterwe.geds:false], got map[]
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v1        0.010s

```

### 잠시 다른 얘기를 하자면...

당신은 이 결과를 얻지 않았을 수 있다. 잠시 후에 얘기할 내용에서도 에러 메시지를 받을 수 있다. 그 메시지를 받더라도 걱정하지 말고 위의 결과를 얻을 때까지 계속해서 _시도_ 해 보라. 혹은 성공 한 척해 보라. 너에게 달렸다. 동시성에 오신 것을 환영한다: 올바르게 처리하지 않으면 무슨 일이 일어날지 예측하기 힘들다. 걱정하지 말라 - 그래서 우리는 동시성을 예측할 수 있게 테스트를 작성하는 것이다.

### ... 다시 돌아와서

`CheckWebsites`가 빈 map 값을 반환하는 것을 볼 수 있다. 무엇이 잘못되었을까?

`for` 반복문이 시작되고 난 후 goroutine들 중 하나도 결괏값을 `results` map에 추가할 시간이 없었다; `WebsiteChecker` 함수가 goroutine들에게는 너무 빨라서 비어있는 map이 반환되는 것이다.

이 점을 고치기 위해 우리는 goroutine들이 일들을 마칠 때까지 기다렸다가 반환하기만 하면 된다. 2초면 되지 않을까?

```go
package concurrency

import "time"

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
  results := make(map[string]bool)

  for _, url := range urls {
    go func() {
      results[url] = wc(url)
    }()
  }

  time.Sleep(2 * time.Second)

  return results
}
```

이제 테스트를 실행해서 값을 얻을 수 있다.(혹은 얻지 못할 수 있다 - 다음과 같이):

```sh
--- FAIL: TestCheckWebsites (0.00s)
        CheckWebsites_test.go:31: Wanted map[http://google.com:true http://blog.gypsydave5.com:true waat://furhurterwe.geds:false], got map[waat://furhurterwe.geds:false]
FAIL
exit status 1
FAIL    github.com/gypsydave5/learn-go-with-tests/concurrency/v1        0.010s
```

이것은 맞지 않다 - 왜 하나의 결과만 얻었을까? 시간을 늘려 시도를 해봐야 할 것 같다 - 원하는 만큼 해 보자. 작동이 되지 않을 것이다. 해당 문제는 변수 `url`이 모든 `for` 반복 때마다 재사용 된다는 것이다 - `urls`에서 매번 새로운 값을 가져간다. 하지만 우리의 각 goroutine들은 각 `url` 변수에 대한 참조를 가지고 있다 - 그들은 그들만의 독립된 복사본을 갖고 있지 않다. 그래서 그들은 모두 `url`이 반복이 끝날 때 갖는 값을 쓰고 있다 - 마지막 url 말이다. 이것이 우리가 결과로 마지막 url만 받은 이유이다.

이것을 고치기 위해:

```go
package concurrency

import (
  "time"
)

type WebsiteChecker func(string) bool

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
  results := make(map[string]bool)

  for _, url := range urls {
    go func(u string) {
      results[u] = wc(u)
    }(url)
  }

  time.Sleep(2 * time.Second)

  return results
}
```

각 익명 함수에 url 매개 변수인 - `u` - 를 부여한 다음 `url`을 인수로 하여 익명 함수를 호출하고, `u`의 값을 goroutine을 실행하는 루프의 반복에 대한 `url` 값으로 고정되도록 한다. `u`는 `url`의 값을 복사한 것이므로 변경되지 않는다.

당신이 운이 좋다면 얻을 것은:

```sh
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v1        2.012s
```

하지만 만약 운이 좋지 않다면 (벤치 마크에서 실행하면 더 많은 시도가 이뤄질 가능성이 높다)

```sh
fatal error: concurrent map writes

goroutine 8 [running]:
runtime.throw(0x12c5895, 0x15)
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/panic.go:605 +0x95 fp=0xc420037700 sp=0xc4200376e0 pc=0x102d395
runtime.mapassign_faststr(0x1271d80, 0xc42007acf0, 0x12c6634, 0x17, 0x0)
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:783 +0x4f5 fp=0xc420037780 sp=0xc420037700 pc=0x100eb65
github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1(0xc42007acf0, 0x12d3938, 0x12c6634, 0x17)
        /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x71 fp=0xc4200377c0 sp=0xc420037780 pc=0x12308f1
runtime.goexit()
        /usr/local/Cellar/go/1.9.3/libexec/src/runtime/asm_amd64.s:2337 +0x1 fp=0xc4200377c8 sp=0xc4200377c0 pc=0x105cf01
created by github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker
        /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xa1

        ... many more scary lines of text ...
```

이 말이 길고 무섭지만, 숨을 쉬면서 스택 추적(stacktrace)을 읽기만 하면 된다: `fatal error: concurrent map writes`. 가끔, 테스트를 할 때, 2개의 goroutine들이 같은 시간에 결과 map에 쓸 때가 있다. Go에 있는 Map은 한 번에 여러 개를 쓰는 것을 싫어하기 때문에 `fetal error`가 난 것이다.

이것을 _경쟁 상태_ 라고 하는데, 이는 소프트웨어 출력이 제어할 수 없는 이벤트의 타이밍과 시퀀스에 종속될 때 발생하는 버그이다. 각 goroutine이 결과 map에 쓰는 시간을 정확하게 제어할 수 없기 때문에, 두 개의 goroutine들이 동시에 결과 map을 쓰는 것에 취약하다.

Go에 내장되어 있는 [_race detector_][godoc_race_detector]는 경쟁 상태를 알려주는 데 도움을 준다. 이 기능을 사용하려면, 테스트를 `race` 옵션과 함께 실행하면 된다: `go test -race`.

당신은 이렇게 생긴 결과물을 받을 것이다:

```sh
==================
WARNING: DATA RACE
Write at 0x00c420084d20 by goroutine 8:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:774 +0x0
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x82

Previous write at 0x00c420084d20 by goroutine 7:
  runtime.mapassign_faststr()
      /usr/local/Cellar/go/1.9.3/libexec/src/runtime/hashmap_fast.go:774 +0x0
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker.func1()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12 +0x82

Goroutine 8 (running) created at:
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xc4
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.TestWebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker_test.go:27 +0xad
  testing.tRunner()
      /usr/local/Cellar/go/1.9.3/libexec/src/testing/testing.go:746 +0x16c

Goroutine 7 (finished) created at:
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.WebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11 +0xc4
  github.com/gypsydave5/learn-go-with-tests/concurrency/v3.TestWebsiteChecker()
      /Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker_test.go:27 +0xad
  testing.tRunner()
      /usr/local/Cellar/go/1.9.3/libexec/src/testing/testing.go:746 +0x16c
==================
```

세부 내용이 나왔고, 읽기 힘들다 - 하지만 `WARNING: DATA RACE`는 꽤 모호하지 않다. 오류 본문을 읽어보면 2가지의 다른 goroutine들이 map에 쓰려고 하는 것을 볼 수 있다. 

`Write at 0x00c420084d20 by goroutine 8:`

은 아래와 같은 메모리 블록에 쓰고 있다.

`Previous write at 0x00c420084d20 by goroutine 7:`

그 윗줄을 보면 몇 번째 줄의 코드에서 일어난 일인지 볼 수 있다:

`/Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:12`

그리고 goroutine 7번과 8번이 시작되는 코드 라인은:

`/Users/gypsydave5/go/src/github.com/gypsydave5/learn-go-with-tests/concurrency/v3/websiteChecker.go:11`

당신이 알아야 하는 것들은 모두 터미널에 출력 되어있다 - 당신이 해야 할 일은 이것을 읽을 만큼 참을 성이 있는 것이다.

### 채널

우리는 _채널_ 을 사용하여 goroutine들을 조직화함으로써 이 경쟁 상태를 해결할 수 있다. 채널들은 값을 수신하거나 전송할 수 있는 Go 데이터 구조이다. 이 연산들은, 세부 정보와 함께 서로 다른 프로세스 간의 통신을 가능하게 한다.

이 경우 우리는 부모 프로세스와 url을 사용하여 `WebsiteChecker` 함수를 수행하게 하는 각 goroutine들 간의 통신에 대해 생각해 보려 한다.

```go
package concurrency

type WebsiteChecker func(string) bool
type result struct {
  string
  bool
}

func CheckWebsites(wc WebsiteChecker, urls []string) map[string]bool {
  results := make(map[string]bool)
  resultChannel := make(chan result)

  for _, url := range urls {
    go func(u string) {
      resultChannel <- result{u, wc(u)}
    }(url)
  }

  for i := 0; i < len(urls); i++ {
    r := <-resultChannel
    results[r.string] = r.bool
  }

  return results
}
```

`results` map과 더불어 이제는 같은 방법으로 `만든(make)` `resultChannel`이 있다. `chan result`는 채널의 타입이다 - `result` 채널의. 새로운 타입인 `result`는 `WebsiteChecher`의 반환 값과 확인 중인 url을 연결하기 위해 만들어졌다 - 이것은 `string`과 `bool`로 이루어졌다. 두 값 중 어느 것도 이름을 붙일 필요가 없기 때문에, 각각의 값은 구조 내에서 익명으로 되어 있다; 이것은 값의 이름을 무엇으로 붙여야 할지 알기 어려울 때 유용할 수 있다.

이제 url을 사용하여 반복할 때, `map`에 바로 적는 것 대신에 `wc`로 각 요청 때마다 `result` 구조를 `resultChannel`에 _보내는 수식_ 과 함께 보낸다. 이것은 `<-` 연산자를 사용하고, 좌측에 있는 채널과 우측의 값을 사용한다:

```go
// 보내는 수식
resultChannel <- result{u, wc(u)}
```

다음 `for` 반복문은 각 url 들에 대해 1번씩 반복된다. 내부에서는 _받는 수식_ 을 사용하고 있는데, 이 식은 채널에서 수신한 값을 변수에 할당한다. 이것 또한 `<-` 연산자를 사용하지만, 2개의 피연산자들의 위치가 뒤바뀐다: 채널이 우측에 위치하고 우리가 할당할 변수는 좌측에 위치한다:

```go
// 받는 수식
r := <-resultChannel
```

그런 다음 수신한 `result`를 사용하여 map을 갱신한다.

채널로 결과를 보내는 것으로, 우리는 결과 map에 쓰는 각 시간들을 제어할 수 있고, 한 번에 하나씩 이루어지는 것을 확실하게 할 수 있다. 각각이 `wc`를 호출하고, 각각 결과 채널로 보내지만, 이 일은 자체 프로세스 내에서 병렬적으로 수행되어 결괏값이 받는 수식을 사용하여 결과 채널에서 값을 추출할 때 각 결과가 한 번에 하나씩 처리된다.

병행적으로 수행되나 연속적으로 수행되었던 일을 바로잡아 우리는 더 빨리 만들고 싶던 부분의 코드를 병행적으로 만들었다. 그리고 우리는 채널을 사용하여 관련된 여러 프로세스들과 소통했다.

벤치마크를 실행하면:

```sh
pkg: github.com/gypsydave5/learn-go-with-tests/concurrency/v2
BenchmarkCheckWebsites-8             100          23406615 ns/op
PASS
ok      github.com/gypsydave5/learn-go-with-tests/concurrency/v2        2.377s
```
23406615 nanoseconds - 0.023 seconds, about one hundred times as fast as
original function. A great success.

## 정리

해당 활동은 TDD에 있어 평소보다 조금 더 가벼운 주제다. 어떤 면에서 우리는 `CheckWebsites` 함수의 긴 리팩터링에 참여하고 있다; 입력과 출력은 변하지 않고, 더 빨라졌을 뿐이다. 그러나 우리가 작성한 벤치마크와 함께 시행한 테스트는, 소프트웨어가 여전히 작동 중이라는 신뢰를 유지하는 방식으로 `Check Website`를 리팩터링 할 수 있게 해주었고 실제로 더 빨라졌음을 보여주었다.

더 빨리 만들기 위해 우리가 배운 것

- *goroutines*, Go에 있는 동시성의 기본 단위, 같은 시간에 1개보다 많은 웹사이트를 확인할 수 있게 해준다.
- *익명 함수*, 각 웹사이트를 확인하는 동시성 프로세스를 시작하기 위해 사용했다.
- *채널*, 다양한 프로세스들을 정리하고 통신할 수 있도록 도와주고, *경쟁 조건*의 버그를 피할 수 있게 해준다.
- *the race detector*은 동시적인 코드에 대한 문제를 디버깅하는 데 도움을 준다.

### 빨리 만들기

소프트웨어 구축 방법의 한 가지 공식인 애자일 방법은(종종 Kent Beck에게서 잘못 이해된다): 

> [Make it work, make it right, make it fast(만들고, 바르게 하고, 빠르게 동작하도록 만들라)][wrf]

'work'는 테스트들을 통과하게 만드는 것이고, 'right'는 코드를 리팩토링하는 것, 그리고 'fast'는 코드를 최적화하는 것, 예를 들어 빠르게 실행되는 것이다. 우리는 그 코드를 바르게 만들어야 'make it fast(빠르게 동작하도록 만들기)'를 할 수 있다. 우리에게 주어진 것은 이미 작동 중임을 증명한 코드였기 때문에 리팩터링할 필요는 없어서 행운이다. 앞의 2단계를 수행하기 전에 'make it fast'를 시도하면 안 된다. 왜냐하면

> [Premature optimization is the root of all evil(조급한 최적화는 모든 악의 근원이다)][popt]
> -- Donald Knuth

[DI]: dependency-injection.md
[wrf]: http://wiki.c2.com/?MakeItWorkMakeItRightMakeItFast
[godoc_race_detector]: https://blog.golang.org/race-detector
[popt]: http://wiki.c2.com/?PrematureOptimization