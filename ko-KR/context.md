# 컨텍스트

**[이 챕터의 모든 코드를 이 곳에서 확인할 수 있습니다.](https://github.com/quii/learn-go-with-tests/tree/main/context)**

소프트웨어는 종종 오랜 기간 실행되고 리소스를 많이 점유하는 프로세스(혹은 고루틴)를 실행합니다. 만약 프로세스를 시작한 작업이 어떠한 이유로 인하여 취소되거나 실패할 경우 프로그램에서 실행된 프로세스들을 일관된 방법으로 멈춰줄 필요가 있습니다.

이를 제대로 처리하지 않는다면 곧 여러분은 여러분이 자랑스러워 하는 그 멋진 고 어플리케이션의 성능 문제를 디버그하는 데에 어려움을 겪게 될 것입니다.

이 챕터에서는 `context` 패키지의 도움을 받아 오래 실행되는 프로세스를 관리해 볼 것입니다.

응답으로 보낼 어떠한 데이터를 가져오는 데에 잠재적으로 오랜 시간을 소비할 수도 있는 프로세스들을 실행하는 전형적인 웹 서버와 함께 시작해봅시다.

데이터를 가져오는 도중에 사용자가 요청을 취소하는 경우를 가정해보고 이 때 프로세스가 중단될 수 있도록 해봅시다.

행복한 경로(happy path) 코드를 준비했습니다. 아래의 서버 코드를 확인해봅시다.

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, store.Fetch())
	}
}
```

`Server` 함수는 `Store` 인수를 받은 뒤 `http.HandlerFunc`를 반환합니다. Store는 다음과 같이 정의되어 있습니다:

```go
type Store interface {
	Fetch() string
}
```

반환된 함수는 `store`의 `Fetch` 메서드를 통해 데이터를 얻은 뒤 응답으로 해당 데이터를 출력합니다.

아래는 이 테스트에 쓰인 `Store`의 구현부분입니다.

```go
type StubStore struct {
	response string
}

func (s *StubStore) Fetch() string {
	return s.response
}

func TestServer(t *testing.T) {
	data := "hello, world"
	svr := Server(&StubStore{data})

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()

	svr.ServeHTTP(response, request)

	if response.Body.String() != data {
		t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
	}
}
```

행복한 경로 코드가 준비되었으니 이제는 약간 더 현실성이 있는 경우를 가정해 볼 차례입니다. 사용자가 요청을 취소하기 전까지 `Store`가 `Fetch`를 끝내지 못하는 경우를 생각해봅시다.

## 우선 테스트를 작성해봅시다

핸들러에서 `Store`를 중단시킬 방법이 필요하므로 인터페이스를 알맞게 바꿔줍시다.

```go
type Store interface {
	Fetch() string
	Cancel()
}
```

스파이를 수정하여 `data`를 가져오는 데에 시간이 걸리도록 하고, 취소 여부를 알 수 있도록 해봅시다. 또한 해당 객체가 어떻게 호출되는지 지켜볼 것이므로 이름을 `SpyStore`로 바꿔줍니다. 그리고 `Store` 인터페이스를 구현하도록 `Cancel` 메서드를 추가해줍니다.

```go
type SpyStore struct {
	response string
	cancelled bool
}

func (s *SpyStore) Fetch() string {
	time.Sleep(100 * time.Millisecond)
	return s.response
}

func (s *SpyStore) Cancel() {
	s.cancelled = true
}
```

100 밀리초가 지나기 전에 요청을 취소하는 테스트를 추가해보고, store가 취소되는지 확인해봅시다.

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
      data := "hello, world"
      store := &SpyStore{response: data}
      svr := Server(store)

      request := httptest.NewRequest(http.MethodGet, "/", nil)

      cancellingCtx, cancel := context.WithCancel(request.Context())
      time.AfterFunc(5 * time.Millisecond, cancel)
      request = request.WithContext(cancellingCtx)

      response := httptest.NewRecorder()

      svr.ServeHTTP(response, request)

      if !store.cancelled {
          t.Errorf("store was not told to cancel")
      }
  })
```

[고 블로그: 컨텍스트](https://blog.golang.org/context)에 따르면

> 컨텍스트 패키지는 기존 컨텍스트들로 부터 새 컨텍스트들을 파생시키는 함수를 제공하며, 이를 사용할 경우 해당 컨텍스트들은 트리를 형성하게 됩니다: 어떠한 컨텍스트가 취소될 경우, 그것으로 부터 파생된 모든 컨텍스트들 또한 취소됩니다.

이 점을 숙지하여 주어진 요청에 대한 취소가 해당 요청의 호출 스택을 따라 전파되어 모든 컨텍스트들이 취소될 수 있도록 컨텍스트를 파생시키는 것이 중요합니다.

`request`에서 새로운 `cancellingCtx`를 파생시킴과 동시에 `cancel` 함수를 얻게 됩니다. 다음으로 `time.AfterFunc`를 통해 해당 함수가 5 밀리초 후에 호출되도록 설정해봅시다. 마지막으로 `request.WithContext`를 통해 새로 얻은 컨텍스트를 사용해봅시다.

## 테스트를 시도해봅시다

예상대로 테스트는 실패할 것입니다.

```go
--- FAIL: TestServer (0.00s)
    --- FAIL: TestServer/tells_store_to_cancel_work_if_request_is_cancelled (0.00s)
    	context_test.go:62: store was not told to cancel
```

## 테스트가 통과할만큼만 코드를 작성해봅시다

TDD를 숙지하며 테스트를 작성해봅시다. *최소한의* 코드를 추가하여 테스트가 통과하도록 해봅시다.

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		store.Cancel()
		fmt.Fprint(w, store.Fetch())
	}
}
```

위와 같이 수정함으로써 테스트를 통과하기는 하지만 기분이 그리 좋지만은 않습니다! 당연한 얘기이지만 *모든 요청*에 대하여 데이터를 가져오기도 전에 `Store`를 취소하여서는 안됩니다.

좋습니다! TDD를 숙지함으로써 테스트의 결점이 보이기 시작했습니다.

행복한 경로 테스트를 수정하여 `Store`가 취소되지 않았음을 assert 하도록 해봅시다.

```go
t.Run("returns data from store", func(t *testing.T) {
    data := "hello, world"
    store := &SpyStore{response: data}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)
    response := httptest.NewRecorder()

    svr.ServeHTTP(response, request)

    if response.Body.String() != data {
        t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
    }

    if store.cancelled {
        t.Error("it should not have cancelled the store")
    }
})
```

두 테스트를 실행해보면 행복한 경로 테스트는 이제 실패할 것입니다. 조금 더 합리적인 구현이 필요한 시점입니다.

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		data := make(chan string, 1)

		go func() {
			data <- store.Fetch()
		}()

		select {
		case d := <-data:
			fmt.Fprint(w, d)
		case <-ctx.Done():
			store.Cancel()
		}
	}
}
```

여기서 우리는 무엇을 했나요?

`context` 에게는 `Done()`이라는 메서드가 있으며 이는 컨텍스트가 "완료"되거나 "취소"될 경우 신호를 받는 채널을 반환합니다. 우리는 해당 신호를 대기하여 해당 신호가 올 경우 `store.Cancel`을 호출하고 싶지만, 만약 `Store`가 그 전에 `Fetch`를 완료했을 경우에는 그 신호를 무시해주고 싶습니다.

이를 위해 고루틴에서 `Fetch`를 호출한 뒤 새로 만들어줄 채널인 `data`에 결과를 보내줍니다. 그리고 `select` 문을 사용하여 두 비동기 프로세스를 경합시킨 뒤 응답을 출력하거나 `Cancel`을 수행합시다.

## 리팩토링

스파이에 assertion 메서드들을 추가하여 테스트 코드를 리팩토링 해봅시다.

```go
type SpyStore struct {
	response  string
	cancelled bool
	t         *testing.T
}

func (s *SpyStore) assertWasCancelled() {
	s.t.Helper()
	if !s.cancelled {
		s.t.Errorf("store was not told to cancel")
	}
}

func (s *SpyStore) assertWasNotCancelled() {
	s.t.Helper()
	if s.cancelled {
		s.t.Errorf("store was told to cancel")
	}
}
```

스파이를 생성할 때 `*testing.T`를 잊지 말고 넘겨주도록 합시다.

```go
func TestServer(t *testing.T) {
	data := "hello, world"

	t.Run("returns data from store", func(t *testing.T) {
		store := &SpyStore{response: data, t: t}
		svr := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)
		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		if response.Body.String() != data {
			t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
		}

		store.assertWasNotCancelled()
	})

	t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
		store := &SpyStore{response: data, t: t}
		svr := Server(store)

		request := httptest.NewRequest(http.MethodGet, "/", nil)

		cancellingCtx, cancel := context.WithCancel(request.Context())
		time.AfterFunc(5*time.Millisecond, cancel)
		request = request.WithContext(cancellingCtx)

		response := httptest.NewRecorder()

		svr.ServeHTTP(response, request)

		store.assertWasCancelled()
	})
}
```

위의 접근 방식은 작동하기는 하지만 자연스러운가요?

웹 서버가 `Store`를 취소하는데에 직접 관여하는 것이 적절하다고 생각하시나요? `Store`가 다른 실행 속도가 느린 프로세스들에 의존하는 경우 어떻게 될까요? `Store.Cancel`이 올바르게 파생 컨텍스트들에게 취소를 전파하도록 해야할 필요가 있습니다.

`context`를 사용하는 주요 이유 중의 하나는 일관된 취소를 수행하기 위함입니다.

[공식 고 문서에 의하면](https://golang.org/pkg/context/)

> 서버로 들어오는 요청들에 대해 컨텍스트를 생성하는 것이 좋고 내보내는 함수는 컨텍스트를 인수로 받는 것이 좋습니다. 또한 두 과정 사이의 함수들을 호출할 때 해당 컨텍스트를 반드시 전파하여야 하며, 선택적으로 해당 컨텍스트를 WithCancel, WithDeadline, WithTimeout, 혹은 WithValue를 이용해 파생시킨 컨텍스트를 사용할 수도 있습니다. 컨텍스트가 취소될 때 해당 컨텍스트를 상속한 모든 컨텍스트들 또한 취소됩니다.

다시 [고 블로그: 컨텍스트](https://blog.golang.org/context)를 살펴보면:

> 구글에서는 고 프로그래머들로 하여금 모든 들어오는 요청과 나가는 요청 함수들의 첫번째 인수를 컨텍스트로 하도록 규정합니다. 이는 여러 팀에서 개발된 고 코드들이 서로 잘 작동하도록 합니다. 컨텍스트는 간단한 방법을 통해 시간 초과와 취소를 관리할 수 있도록 하며, 보안 자격 증명과 같은 중요한 값들이 고 프로그램내에서 올바르게 넘겨질 수 있도록 합니다.

(잠시 시간을 내어 모든 함수가 컨텍스트를 넘겨줄 경우 가져올 영향과 그것을 인간공학적인 관점에서 생각해 봅시다.)

약간 불편하게 느껴지시나요? 좋습니다. 불편할지라도 해당 접근 방식을 따라하여 `Store`에 `context`를 넘겨줌으로써 관여하게 해봅시다. 이를 통해 `Store`는 해당 `context`를 그것에 의존하는 것들에 넘겨줄 수 있게 되고, 그 컨텍스트들은 그것들을 취소하는 데에 관여하게 됩니다.

## 우선 테스트를 작성해봅시다

각 구성요소가 관여하는 부분이 바뀌었으므로 테스트 또한 수정해줍니다. 핸들러가 담당하는 부분은 이제 단지 컨텍스트를 `Store`를 통해 전파시키는 것과 `Store`가 취소될 경우 보내지는 오류를 처리하는 것입니다.

`Store` 인터페이스를 수정하여 새로 관여하는 부분을 담당할 수 있도록 합니다.

```go
type Store interface {
	Fetch(ctx context.Context) (string, error)
}
```

우선 핸들러 내부의 코드를 지워줍시다.

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	}
}
```

`SpyStore`를 수정해줍시다.

```go
type SpyStore struct {
	response string
	t        *testing.T
}

func (s *SpyStore) Fetch(ctx context.Context) (string, error) {
	data := make(chan string, 1)

	go func() {
		var result string
		for _, c := range s.response {
			select {
			case <-ctx.Done():
				s.t.Log("spy store got cancelled")
				return
			default:
				time.Sleep(10 * time.Millisecond)
				result += string(c)
			}
		}
		data <- result
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case res := <-data:
		return res, nil
	}
}
```

스파이를 수정하여 `context`를 실제로 사용하는 메서드처럼 바꿔봅시다.

결과값의 문자열을 한글자씩 덧붙이는 느린 모의 프로세스 고루틴을 만듭시다. 고루틴이 완료됨과 동시에 `data` 채널에 결과값을 보내주고 고루틴에서 `ctx.Done`을 대기하여 값이 들어올 경우 작업을 중단하게 해봅시다.

마지막으로 추가적인 `select` 문을 사용하여 해당 고루틴이 완료되거나 취소되는 것을 기다리도록 합니다.

이전의 접근 방식과 비슷하지만 이번에는 고에 내장된 동시성 기능을 사용하여 두개의 비동기 프로세스를 경합시켜 반환할 값을 정하게 됩니다.

`context`를 사용하는 함수들과 메서드들을 만들 경우 아래와 비슷한 접근 방식을 사용하게 되므로 꼭 작동 방식을 이해해보도록 합시다.

이제 테스트들을 수정해줄 차례입니다. 이전에 진행한 취소 테스트를 지워줌으로써 행복한 경로 테스트를 수정해 봅시다.

```go
t.Run("returns data from store", func(t *testing.T) {
    data := "hello, world"
    store := &SpyStore{response: data, t: t}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)
    response := httptest.NewRecorder()

    svr.ServeHTTP(response, request)

    if response.Body.String() != data {
        t.Errorf(`got "%s", want "%s"`, response.Body.String(), data)
    }
})
```

## 테스트를 시도해봅시다

```
=== RUN   TestServer/returns_data_from_store
--- FAIL: TestServer (0.00s)
    --- FAIL: TestServer/returns_data_from_store (0.00s)
    	context_test.go:22: got "", want "hello, world"
```

## 테스트가 통과할만큼만 코드를 작성해봅시다

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, _ := store.Fetch(r.Context())
		fmt.Fprint(w, data)
	}
}
```

행복한 경로는 이제 다시... 행복합니다. 이번에는 다른 테스트를 수정해 줄 차례입니다.

## 우선 테스트를 작성해봅시다

오류가 발생할 경우 응답으로 아무것도 출력되지 않는 것을 테스트해야 합니다. 안타깝게도 `httptest.ResponseRecorder`는 이러한 기능을 제공하지 않으므로 스파이를 추가하여 이를 테스트 해봅시다.

```go
type SpyResponseWriter struct {
	written bool
}

func (s *SpyResponseWriter) Header() http.Header {
	s.written = true
	return nil
}

func (s *SpyResponseWriter) Write([]byte) (int, error) {
	s.written = true
	return 0, errors.New("not implemented")
}

func (s *SpyResponseWriter) WriteHeader(statusCode int) {
	s.written = true
}
```

위의 `SpyResponseWriter`는 `http.ResponseWriter` 인터페이스를 구현하기에 테스트에 사용할 수 있습니다.

```go
t.Run("tells store to cancel work if request is cancelled", func(t *testing.T) {
    store := &SpyStore{response: data, t: t}
    svr := Server(store)

    request := httptest.NewRequest(http.MethodGet, "/", nil)

    cancellingCtx, cancel := context.WithCancel(request.Context())
    time.AfterFunc(5*time.Millisecond, cancel)
    request = request.WithContext(cancellingCtx)

    response := &SpyResponseWriter{}

    svr.ServeHTTP(response, request)

    if response.written {
        t.Error("a response should not have been written")
    }
})
```

## 테스트를 시도해봅시다

```
=== RUN   TestServer
=== RUN   TestServer/tells_store_to_cancel_work_if_request_is_cancelled
--- FAIL: TestServer (0.01s)
    --- FAIL: TestServer/tells_store_to_cancel_work_if_request_is_cancelled (0.01s)
    	context_test.go:47: a response should not have been written
```

## 테스트가 통과할만큼만 코드를 작성해봅시다

```go
func Server(store Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := store.Fetch(r.Context())

		if err != nil {
			return // todo: log error however you like
		}

		fmt.Fprint(w, data)
	}
}
```

이제 서버 코드가 매우 간단해진 것을 확인할 수 있습니다. 직접적으로 취소하는 데에 관여하지 않고 단순히 `context`를 넘겨주고 이후로 호출된 함수들에게 의존하기 때문입니다.

## 마치며

### 이 챕터에서 다룬 내용

- 클라이언트가 요청을 취소할 때의 HTTP 핸들러를 테스트 하는 방법
- 컨텍스트를 사용하여 취소를 관리하는 법
- `context`를 인수로 받는 함수를 작성하고 고루틴, `select` 문, 채널을 이용하여 해당 컨텍스트를 취소하는 방법
- 구글의 가이드라인에 나와있는 데로 요청에 대한 호출 스택에 유효한 컨텍스트를 전파하여 취소를 관리하는 법
- `http.ResponseWriter`의 스파이를 작성하는 법

### context.Value는 무엇인가요?

필자와 [Michal Štrba](https://faiface.github.io/post/context-should-go-away-go2/)는 이에 대해 비슷한 의견을 가지고 있습니다.

> 내 회사(존재하지는 않지만)에서 ctx.Value를 사용하면 해고당합니다

몇몇 엔지니어들은 `context`를 통해 값들을 전해주는 것이 *편리하다는* 이유로 옹호하곤 합니다.

하지만 편의성은 종종 나쁜 코드를 만들어냅니다.

`context.Values`는 단순히 타입이 지정되지 않은 맵이기 때문에 타입 안전성이 보장되지 않고 실제로는 가지고 있지 않은 값을 처리해줘야 하는 문제점을 가지고 있습니다. 한 모듈에서 다른 모듈로 보낼 경우 맵 키들의 대응 목록을 만들어줘야 하고, 누군가 코드를 수정하기 시작하면 문제가 발생하기 시작합니다.

다시 말해, **함수에 값을 넘겨주려면 `context.Value`를 사용하지 말고 타입이 지정된 인수로 넘겨줘야 합니다**. 이는 정적으로 해당 부분이 검수되게 하고 모든 사람이 문서를 읽을 수 있도록 합니다.

#### 하지만...

트레이스 id와 같이 요청과 관련없는 정보를 이용할 때에는 도움이 될 수 있습니다. 호출 스택의 모든 함수에서 해당 정보를 필요로 할 가능성은 낮은 데다 이를 함수 인수로 포함할 경우 함수의 시그니쳐가 복잡해질 수 있기 때문입니다.

[Jack Lindamood는 **Context.Value는 제어하기 보다는 정보만 제공해야한다고 주장합니다**](https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39)

> context.Value의 내용은 사용자를 위한 것이 아니라 관리자를 위한 것입니다. 기대되거나 문서화된 결과 값에 필요한 입력값이 되어서는 절대 안됩니다.

### 추가 자료

- 필자는 [Michal Štrba의 Go 2에서는 컨텍스트가 없어져야 합니다](https://faiface.github.io/post/context-should-go-away-go2/)를 흥미롭게 읽었습니다. 그가 주장하는 바는 `context`를 모든 곳에서 넘겨줘야하는 것은 탐탁하지 않고 이는 곧 고 언어가 가진 취소를 관리하는 데에 있어서의 부족함을 드러낸다는 것입니다. 그는 또한 이러한 문제점이 라이브러리 레벨이 아닌 언어적인 레벨에서 수정되기를 바랍니다. 이러한 문제점이 해결되기 전까지는 오래 실행되는 프로세스를 관리하는 데에 있어 `context`는 필요한 존재입니다.
- [고 블로그에서 `context`를 사용해야하는 이유와 몇몇 예제를 추가적으로 다루고 있습니다.](https://blog.golang.org/context)
