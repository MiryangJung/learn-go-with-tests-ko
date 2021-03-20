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

## 테스트를 실행해보기

`./racer_test.go:14:9: undefined: Racer`

## 테스트를 실행할 최소한의 코드를 작성하고 실패한 테스트 출력을 확인하기

```go
func Racer(a, b string) (winner string) {
	return
}
```

`racer_test.go:25: got '', want 'http://www.quii.co.uk'`

## 통과할 만큼 충분한 코드를 작성하기

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

## 리팩토링

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

우리는 가짜 서버를 `makeDelayedServer`라는 함수로 리팩토링 하여 테스트에서 흥미롭지 않은 코드를 옮기고 반복을 줄였다.

### `defer`

함수 호출 앞에 `defer`를 붙이면 해당 함수를 _포함하는 함수의 끝에서_ 호출합니다.

때로는 파일을 닫거나 서버를 닫는 것과 같은 리소스를 정리하여 포트가 계속 수신하지 않도록 해야 한다.

함수가 끝날 때 실행되기를 원하지만, 나중에 코드를 읽는 사람을 위해 서버를 생성한 위치 근처에 명령어를 보관한다.

리팩토링은 개선된 것이며 지금까지 다루었던 Go 기능을 고려할 때 합리적인 해결책이지만 해결책을 더 간단하게 만들 수 있다.

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

If you recall from the concurrency chapter, you can wait for values to be sent to a channel with `myVar := <-ch`. This is a _blocking_ call, as you're waiting for a value.

What `select` lets you do is wait on _multiple_ channels. The first one to send a value "wins" and the code underneath the `case` is executed.

We use `ping` in our `select` to set up two channels for each of our `URL`s. Whichever one writes to its channel first will have its code executed in the `select`, which results in its `URL` being returned (and being the winner).

After these changes, the intent behind our code is very clear and the implementation is actually simpler.

### 시간초과

Our final requirement was to return an error if `Racer` takes longer than 10 seconds.

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

We've made our test servers take longer than 10s to return to exercise this scenario and we are expecting `Racer` to return two values now, the winning URL (which we ignore in this test with `_`) and an `error`.

## 테스트를 실행해보기

`./racer_test.go:37:10: assignment mismatch: 2 variables but 1 values`

## 테스트를 실행할 최소한의 코드를 작성하고 실패한 테스트 출력을 확인하기

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

Change the signature of `Racer` to return the winner and an `error`. Return `nil` for our happy cases.

The compiler will complain about your _first test_ only looking for one value so change that line to `got, _ := Racer(slowURL, fastURL)`, knowing that we should check we _don't_ get an error in our happy scenario.

If you run it now after 11 seconds it will fail.

```
--- FAIL: TestRacer (12.00s)
    --- FAIL: TestRacer/returns_an_error_if_a_server_doesn't_respond_within_10s (12.00s)
        racer_test.go:40: expected an error but didn't get one
```

## 통과할 만큼 충분한 코드를 작성하기

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

`time.After` is a very handy function when using `select`. Although it didn't happen in our case you can potentially write code that blocks forever if the channels you're listening on never return a value. `time.After` returns a `chan` (like `ping`) and will send a signal down it after the amount of time you define.

For us this is perfect; if `a` or `b` manage to return they win, but if we get to 10 seconds then our `time.After` will send a signal and we'll return an `error`.

### 느린 테스트

The problem we have is that this test takes 10 seconds to run. For such a simple bit of logic, this doesn't feel great.

What we can do is make the timeout configurable. So in our test, we can have a very short timeout and then when the code is used in the real world it can be set to 10 seconds.

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

Our tests now won't compile because we're not supplying a timeout.

Before rushing in to add this default value to both our tests let's _listen to them_.

- Do we care about the timeout in the "happy" test?
- The requirements were explicit about the timeout.

Given this knowledge, let's do a little refactoring to be sympathetic to both our tests and the users of our code.

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

Our users and our first test can use `Racer` (which uses `ConfigurableRacer` under the hood) and our sad path test can use `ConfigurableRacer`.

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

I added one final check on the first test to verify we don't get an `error`.

## 마무리

### `select`

- Helps you wait on multiple channels.
- Sometimes you'll want to include `time.After` in one of your `cases` to prevent your system blocking forever.

### `httptest`

- A convenient way of creating test servers so you can have reliable and controllable tests.
- Using the same interfaces as the "real" `net/http` servers which is consistent and less for you to learn.
