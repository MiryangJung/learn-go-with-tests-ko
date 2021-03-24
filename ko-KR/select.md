# 선택

**[이 챕터의 모든 코드는 여기에서 확인할 수 있다.](https://github.com/MiryangJung/learn-go-with-tests-ko/tree/main/select)**

HTTP GET으로 두 개의 URL을 가지고 먼저 반환된 URL을 반환하여 "경쟁"하는 `WebSiteRacer`라는 함수를 만들라는 요청을 받았다. 10초 이내에 반환되는 항목이 없으면 오류를 반환해야 한다.

이를 위해 우리는 아래 목록을 사용해야 한다.

- `net/http`을 사용해 HTTP 호출을 한다.
- `net/http/httptest`를 사용해 테스트를 한다.
- 고루틴을 사용한다.
- 프로세스를 동기화하기 위해 `select` 한다.

## 테스트부터 작성하기

움직이기 위해 단순한 것부터 시작한다.

```go
func TestRacer(t *testing.T) {
	slowURL := "http://www.facebook.com"
	fastURL := "http://www.quii.co.uk"

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

우리는 이것이 완벽하지 않고 문제가 있다는 것을 알고 있지만 그것은 우리를 움직일 것이다. 처음부터 완벽하게 만드는 것에 너무 얽매이지 않는 것이 중요하다.

## 테스트 실행해보기

`./racer_test.go:14:9: undefined: Racer`

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

```go
func Racer(a, b string) (winner string) {
	return
}
```

`racer_test.go:25: got '', want 'http://www.quii.co.uk'`

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func Racer(a, b string) (winner string) {
	startA := time.Now()
	http.Get(a)
	aDuration := time.Since(startA)

	startB := time.Now()
	http.Get(b)
	bDuration := time.Since(startB)

	if aDuration < bDuration {
		return a
	}

	return b
}
```

각 URL에 대해 아래의 작업을 한다.

1. `URL`을 가져 오기 전에 `time.Now()`를 사용하여 기록한다.
2. 그런 다음 [`http.Get`](https://golang.org/pkg/net/http/#Client.Get)을 사용하여 `URL`의 내용을 가져온다. 이 함수는 [`http.Response`](https://golang.org/pkg/net/http/#Response)와 `error`를 반환하지만 지금까지는이 값에 관심이 없다.
3. `time.Since`는 시작 시간을 받으며 차이인 `time.Duration`을 반환한다.

일단 이 작업을 완료하면 가장 빠른 시간을 확인하기 위해 단순하게 걸린 시간을 비교한다.

### 문제점

이렇게 하면 테스트가 통과될 수도 있고 통과되지 못할 수도 있다. 문제는 우리가 우리의 논리를 시험하기 위해 실제 웹사이트에 손을 뻗는다는 것이다.

HTTP를 사용하는 테스트 코드는 매우 일반적이기 때문에 Go는 표준 라이브러리에 테스트하는 데 도움이 되는 도구를 가지고 있다.

mocking과 의존성 주입 챕터에서는 아래 같은 이유로 외부 서비스에 의존하지 않는 것이 얼마나 이상적일 수 있는지 살펴봤다.

- 느리다.
- 믿을 수 없다.
- 엣지 케이스를 테스트할 수 없다.

표준 라이브러리에서는 모의 HTTP 서버를 쉽게 만들 수 있는 [`net/http/httptest`](https://golang.org/pkg/net/http/httptest/)라고 불리는 패키지가 있다.

테스트를 mock을 사용하여 제어할 수 있는 안정적인 서버를 확보하도록 변경해보자.

```go
func TestRacer(t *testing.T) {

	slowServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(20 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))

	fastServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	slowURL := slowServer.URL
	fastURL := fastServer.URL

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}

	slowServer.Close()
	fastServer.Close()
}
```

문법이 조금 복잡해 보일 수 있지만 천천히 해보자.

`httptest.NewServer`는 *익명 함수*를 이용해 보내는 `http.HandlerFunc`를 받는다.

`http.HandlerFunc`의 타입은 `type HandlerFunc func(ResponseWriter, *Request)`와 같다.

실제로는 `ResponseWriter`와 `Request`를 받는 함수가 필요하다는 것이다. 이것은 HTTP 서버에는 그렇게 놀라운 일이 아니다.

여기에는 추가적인 마법이 없다는 것이 밝혀졌다. **이것은 Go에서 _실제_ HTTP 서버를 작성하는 방법이기도 하다**. 유일한 차이점은 `httptest.NewServer`로 감싸는 것이다. 이것은 요청을 대기할 열린 포트를 찾고 테스트가 끝나면 닫을 수 있기 때문에 테스트와 함께 사용하기가 더 쉽다.

두 서버 내에서 느린 서버는 다른 서버보다 느리게 만들라는 요청을 받으면 짧은 `time.Sleep`을 만든다. 그런 다음 두 서버 모두 `w.WriteHeader (http.StatusOK)`와 함께 `OK` 응답을 호출자에게 반환한다.

테스트를 다시 실행하면 확실히 통과 할 것이며 더 빨라질 것이다. 의도적으로 테스트를 실패시키기 위해 sleep을 사용해라.

## 리팩터링 하기

프로덕션 코드와 테스트 코드 모두에 약간의 중복이 있다.

```go
func Racer(a, b string) (winner string) {
	aDuration := measureResponseTime(a)
	bDuration := measureResponseTime(b)

	if aDuration < bDuration {
		return a
	}

	return b
}

func measureResponseTime(url string) time.Duration {
	start := time.Now()
	http.Get(url)
	return time.Since(start)
}
```

이런 건조(DRY-ing up)는 `Racer` 코드를 훨씬 쉽게 읽을 수 있게 만든다.

```go
func TestRacer(t *testing.T) {

	slowServer := makeDelayedServer(20 * time.Millisecond)
	fastServer := makeDelayedServer(0 * time.Millisecond)

	defer slowServer.Close()
	defer fastServer.Close()

	slowURL := slowServer.URL
	fastURL := fastServer.URL

	want := fastURL
	got := Racer(slowURL, fastURL)

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func makeDelayedServer(delay time.Duration) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(delay)
		w.WriteHeader(http.StatusOK)
	}))
}
```

우리는 가짜 서버를 `makeDelayedServer`라는 함수로 리팩터링 하여 테스트에서 흥미롭지 않은 코드를 옮기고 반복을 줄였다.

### `defer`

함수 호출 앞에 `defer`를 붙이면 해당 함수를 _포함하는 함수의 끝에서_ 호출합니다.

때로는 파일을 닫거나 서버를 닫는 것과 같은 리소스를 정리하여 포트가 계속 수신하지 않도록 해야 한다.

함수가 끝날 때 실행되기를 원하지만, 나중에 코드를 읽는 사람을 위해 서버를 생성한 위치 근처에 명령어를 보관한다.

리팩터링은 개선된 것이며 지금까지 다루었던 Go 기능을 고려할 때 합리적인 해결책이지만 해결책을 더 간단하게 만들 수 있다.

### 동기화 프로세스

- Go가 동시성이 뛰어나지만 웹 사이트의 속도를 차례로 테스트하는 이유는 무엇일까? 두 가지를 동시에 확인할 수 있어야한다.
- 우리는 요청의 *정확한 응답 시간*에 대해 신경 쓰지 않고 어떤 것이 먼저 돌아 오는지 알고 싶다.

이를 위해 동기화 프로세스를 정말 쉽고 명확하게 하는 데 도움이 되는 select라는 새로운 구조를 도입할 것이다.

```go
func Racer(a, b string) (winner string) {
	select {
	case <-ping(a):
		return a
	case <-ping(b):
		return b
	}
}

func ping(url string) chan struct{} {
	ch := make(chan struct{})
	go func() {
		http.Get(url)
		close(ch)
	}()
	return ch
}
```

#### `ping`

`chan struct {}`를 생성하고 그것을 반환하는 `ping` 함수를 정의했다.

이 경우 채널에 어떤 타입이 전송되는지는 _신경_ 쓰지 않고 _단지 완료되었음을 알리고_ 채널을 닫으면 완벽하게 작동한다!

왜 `bool`과 같은 다른 타입이 아닌 `struct {}`일까? `chan struct {}`는 메모리 관점에서 사용할 수 있는 가장 작은 데이터 타입이므로 `bool`에 비해 할당이 없다. 닫은 후에 채널에 아무것도 보내지 않는데 왜 할당해야 할까요?

동일한 함수 내에서 한번 `http.Get(url)`을 완료하면 해당 채널로 신호를 보내는 고루틴을 시작한다.

##### 항상 `make` 함수로 채널 만들기

`var ch chan struct {}`를 선언하는 것보다 채널을 만들 때 `make`를 어떻게 사용해야 하는지 주의해야 한다. `var`를 사용할 때 변수는 타입의 "비어있는" 값으로 초기화된다. 따라서 `string`의 경우 `""`, `int`는 0 등으로 초기화된다.

채널의 경우 비어있는 값은 `nil`이고 `<-`을 사용해 보내려고 하면 `nil` 채널로 보낼 수 없기 때문에 영원히 차단된다.

[Go Playground에서 이 동작을 확인할 수 있다.](https://play.golang.org/p/IIbeAox5jKA)

#### `select`

동시성 챕터에서 생각해 보면, `myVar := <- ch`를 사용해 값이 채널로 전송 될 때까지 기다릴 수 있다. 값을 기다리고 있으므로 _차단_ 호출이다.

`select`을 사용하면 _여러_ 채널에서 대기 할 수 있습니다. 처음으로 값을 보내는 항목이 "승리"하고 `case` 아래의 코드가 실행된다.

`select`에서 `ping`을 사용하여 각 `URL`에 대해 두 개의 채널을 설정한다. 어느 쪽이 먼저 채널에 작성하든 `select`에서 코드가 실행되어 `URL`이 반환되고 승자가 된다.

이러한 변경을 한 후 코드의 의도는 매우 명확하고 구현이 실제로 더 간단하다.

### 시간초과

마지막 요구 사항은 `Racer`가 10 초 이상 걸리면 오류를 반환하는 것이다.

## 테스트부터 작성하기

```go
t.Run("returns an error if a server doesn't respond within 10s", func(t *testing.T) {
	serverA := makeDelayedServer(11 * time.Second)
	serverB := makeDelayedServer(12 * time.Second)

	defer serverA.Close()
	defer serverB.Close()

	_, err := Racer(serverA.URL, serverB.URL)

	if err == nil {
		t.Error("expected an error but didn't get one")
	}
})
```

우리는 테스트 서버가 이 시나리오를 실행하는 데 10 초 이상 걸리도록 만들었으며 이제 `Racer`가 두 개의 값, 즉 승리 URL (이 테스트에서는 `_`로 무시함)과 `error`를 반환할 것으로 예상한다.

## 테스트 실행해보기

`./racer_test.go:37:10: assignment mismatch: 2 variables but 1 values`

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

```go
func Racer(a, b string) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	}
}
```

`Racer`의 시그니처를 변경하여 승자와 `error`를 반환한다. 행복한 케이스에 대해서는 `nil`을 반환한다.

컴파일러는 하나의 값만 찾는 *첫 번째 테스트*에 대해 불평할 것이므로 해당 줄을 `got, _ := Racer(slowURL, fastURL)`로 변경하여 행복한 시나리오에서 오류가 발생하지 _않는지_ 확인해야 한다.

테스트를 실행하면 11초 뒤에 실패할 것이다.

```
--- FAIL: TestRacer (12.00s)
    --- FAIL: TestRacer/returns_an_error_if_a_server_doesn't_respond_within_10s (12.00s)
        racer_test.go:40: expected an error but didn't get one
```

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func Racer(a, b string) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(10 * time.Second):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}
```

`time.After`는 `select`를 사용할 때 매우 편리한 기능이다. 우리의 경우에는 발생하지 않았지만 수신중인 채널이 값을 반환하지 않으면 영원히 차단되는 코드를 잠재적으로 작성할 수 있다. `time.After`는 `chan` (`ping`과 같은)을 반환하고 정의한 시간 후에 신호를 보낸다.

우리에게 이것은 완벽하다. `a` 또는 `b`가 반환하면 승리하지만 10 초가 되면 `time.After`가 신호를 보내고 오류를 반환하게 된다.

### 느린 테스트

문제는 이 테스트를 실행하는 데 10 초가 걸린다는 것이다. 그런 간단한 논리로는 기분이 좋지 않다.

우리가 할 수 있는 일은 시간제한을 구성 가능하게 만드는 것이다. 따라서 테스트에서 매우 짧은 시간제한을 가질 수 있으며 코드가 실제 세계에서 사용될 때 10 초로 설정할 수 있다.

```go
func Racer(a, b string, timeout time.Duration) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}
```

타임 아웃을 제공하지 않기 때문에 이제 테스트가 컴파일되지 않는다.

이 기본값을 두 테스트에 모두 추가하기 전에 _기다려_ 보겠습니다.

- "행복한"테스트의 시간 초과에 대해 신경 쓸지?
- 제한 시간에 대한 요구 사항이 명시되어 있다.

이 지식을 감안할 때 테스트와 코드 사용자 모두에게 공감할 수 있도록 약간의 리팩터링을 해보겠다.

```go
var tenSecondTimeout = 10 * time.Second

func Racer(a, b string) (winner string, error error) {
	return ConfigurableRacer(a, b, tenSecondTimeout)
}

func ConfigurableRacer(a, b string, timeout time.Duration) (winner string, error error) {
	select {
	case <-ping(a):
		return a, nil
	case <-ping(b):
		return b, nil
	case <-time.After(timeout):
		return "", fmt.Errorf("timed out waiting for %s and %s", a, b)
	}
}
```

사용자와 첫 번째 테스트에서는 `Racer` (내부에서 `ConfigurableRacer` 사용)를 사용할 수 있고 슬픈 경로 테스트에서는 `ConfigurableRacer`를 사용할 수 있다.

```go
func TestRacer(t *testing.T) {

	t.Run("compares speeds of servers, returning the url of the fastest one", func(t *testing.T) {
		slowServer := makeDelayedServer(20 * time.Millisecond)
		fastServer := makeDelayedServer(0 * time.Millisecond)

		defer slowServer.Close()
		defer fastServer.Close()

		slowURL := slowServer.URL
		fastURL := fastServer.URL

		want := fastURL
		got, err := Racer(slowURL, fastURL)

		if err != nil {
			t.Fatalf("did not expect an error but got one %v", err)
		}

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns an error if a server doesn't respond within the specified time", func(t *testing.T) {
		server := makeDelayedServer(25 * time.Millisecond)

		defer server.Close()

		_, err := ConfigurableRacer(server.URL, server.URL, 20*time.Millisecond)

		if err == nil {
			t.Error("expected an error but didn't get one")
		}
	})
}
```

`error`가 없는지 확인하기 위해 첫 번째 테스트에서 최종 확인을 추가했다.

## 정리

### `select`

- 여러 채널에서 대기할 수 있다.
- 때로는 시스템이 영원히 차단되는 것을 방지하기 위해 `cases` 중 하나에 `time.After`를 포함하고 싶을 것이다.

### `httptest`

- 신뢰할 수 있고 제어 가능한 테스트를 수행할 수 있도록 테스트 서버를 만드는 편리한 방법이다.
- 일관되고 배우기 어려운 "실제" `net/http` 서버와 동일한 인터페이스를 사용한다.
