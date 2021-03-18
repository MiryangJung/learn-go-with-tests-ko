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
3. `time.Since`는 시작 시간을 기록하고 차이인 `time.Duration`을 반환한다.

일단 이 작업을 완료하면 가장 빠른 시간을 확인하기 위해 단순하게 걸린 시간을 비교한다.

### 문제점

This may or may not make the test pass for you. The problem is we're reaching out to real websites to test our own logic.

Testing code that uses HTTP is so common that Go has tools in the standard library to help you test it.

In the mocking and dependency injection chapters, we covered how ideally we don't want to be relying on external services to test our code because they can be

- Slow
- Flaky
- Can't test edge cases

In the standard library, there is a package called [`net/http/httptest`](https://golang.org/pkg/net/http/httptest/) where you can easily create a mock HTTP server.

Let's change our tests to use mocks so we have reliable servers to test against that we can control.

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

The syntax may look a bit busy but just take your time.

`httptest.NewServer` takes an `http.HandlerFunc` which we are sending in via an _anonymous function_.

`http.HandlerFunc` is a type that looks like this: `type HandlerFunc func(ResponseWriter, *Request)`.

All it's really saying is it needs a function that takes a `ResponseWriter` and a `Request`, which is not too surprising for an HTTP server.

It turns out there's really no extra magic here, **this is also how you would write a _real_ HTTP server in Go**. The only difference is we are wrapping it in an `httptest.NewServer` which makes it easier to use with testing, as it finds an open port to listen on and then you can close it when you're done with your test.

Inside our two servers, we make the slow one have a short `time.Sleep` when we get a request to make it slower than the other one. Both servers then write an `OK` response with `w.WriteHeader(http.StatusOK)` back to the caller.

If you re-run the test it will definitely pass now and should be faster. Play with these sleeps to deliberately break the test.

## 리팩토링

We have some duplication in both our production code and test code.

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

This DRY-ing up makes our `Racer` code a lot easier to read.

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

We've refactored creating our fake servers into a function called `makeDelayedServer` to move some uninteresting code out of the test and reduce repetition.

### `defer`

By prefixing a function call with `defer` it will now call that function _at the end of the containing function_.

Sometimes you will need to cleanup resources, such as closing a file or in our case closing a server so that it does not continue to listen to a port.

You want this to execute at the end of the function, but keep the instruction near where you created the server for the benefit of future readers of the code.

Our refactoring is an improvement and is a reasonable solution given the Go features covered so far, but we can make the solution simpler.

### 동기화 프로세스

- Why are we testing the speeds of the websites one after another when Go is great at concurrency? We should be able to check both at the same time.
- We don't really care about _the exact response times_ of the requests, we just want to know which one comes back first.

To do this, we're going to introduce a new construct called `select` which helps us synchronise processes really easily and clearly.

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

We have defined a function `ping` which creates a `chan struct{}` and returns it.

In our case, we don't _care_ what type is sent to the channel, _we just want to signal we are done_ and closing the channel works perfectly!

Why `struct{}` and not another type like a `bool`? Well, a `chan struct{}` is the smallest data type available from a memory perspective so we
get no allocation versus a `bool`. Since we are closing and not sending anything on the chan, why allocate anything?

Inside the same function, we start a goroutine which will send a signal into that channel once we have completed `http.Get(url)`.

##### 항상 `make` 함수로 채널 만들기

Notice how we have to use `make` when creating a channel; rather than say `var ch chan struct{}`. When you use `var` the variable will be initialised with the "zero" value of the type. So for `string` it is `""`, `int` it is 0, etc.

For channels the zero value is `nil` and if you try and send to it with `<-` it will block forever because you cannot send to `nil` channels

[You can see this in action in The Go Playground](https://play.golang.org/p/IIbeAox5jKA)

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
