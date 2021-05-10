# HTTP Server

**[이 챕터의 모든 코드는 여기에서 확인할 수 있다.](https://github.com/quii/learn-go-with-tests/tree/main/http-server)**

사용자가 자신이 얼마나 많은 게임을 이겼는지 확인할 수 있는 웹서버를 만들어야 한다. 

-   `GET /players/{name}` 는 전체 승점을 리턴해야 한다. 
-   `POST /players/{name}` 를 하나 보낼때 마다 전체 승점이 하나씩 증가해야 한다. 

TDD 방식을 따라서 동작하는 소프트웨어를 가능한 빨리 만든 다음, 목표로 하는 구현을 완료할 때까지 반복적으로 작은 개선을 해나갈 것이다. 이렇게 구현을 하면

- 어떤 순간에도, 문제 발생 범위가 작게 유지된다
- rabbit holes(*실제 구현 이외에 너무 시간을 빼앗김)에 빠지지 않는다
- 문제가 발생하여 이전으로 되돌아가도 낭비한 작업이 적게 된다

## 적색, 녹색, 리팩터 

이 책 내내, 테스트를 먼저 작성하고 실패하게 하고 (적색), 동작하는 _최소한_ 의 코드를 짠 다음 (녹색), 리팩터링하는 TDD 프로세스를 강조해왔다. 

최소한의 코드만을 짠다는 원칙은 TDD가 가져다주는 안전성의 측면에서 중요하다. "적색"에서 가능한 빨리 벗어나려 해야 한다. 

켄트 벡은 다음과 같이 말했다.

> 무슨 짓을 해서라도 테스트를 빨리 통과하라

테스트 통과를 위해 저지른 잘못들은 리팩터링으로 고쳐나가면 되며, 리팩터링은 테스트를 통해 안전하게 진행할 수 있다. 

### 이렇게 하지 않는다면? 

적색 상태에서 수정을 더욱 많이 할 수록 테스트로 커버되지 않는 더 많은 문제가 추가되기 쉽다. 

말하고자 하는 것은, 테스트를 통과하는 유용한 코드를 조금씩, 반복적으로 써나가라는 것이다. 이렇게 하면 몇 시간씩 래빗홀에 빠지지 않게 된다. 

### 닭과 달걀

어떻게 하나씩 구현할 수 있을까? 승점이 하나도 저장되어 있지 않으면 `GET` player 할 수 없고, `GET` 엔드포인트가 없으면, `POST`가 동작하는지 알기 어렵다. 

이럴 때에 _mocking_ 이 필요하다. 

- `GET` 은 player의 점수를 얻기위해 `PlayerStore` _같은 것_ 이 필요하며 인터페이스여야 한다. 실제 저장 코드를 만들 필요없이 간단한 스텁(stub)을 생성하여 테스트 할 수 있기 때문이다. 
- `POST` 가 `PlayerStore`를 호출할 때에 제대로 저장하는지 _훔쳐볼_ 수 있어야 한다. 저장과 검색의 구현은 커플링되지 않도록 할 것이다. 
- 작동하는 소프트웨어를 빨리 만들기 위해, 매우 간단한 인-메모리 구현을 한 다음, 원하는 저장 메커니즘에 기반한 구현을 할 것이다. 

## Write the test first

테스트를 짜고 하드코딩된 값을 리턴하게 구현하여 테스트를 통과하자. 여기서부터 시작이다. 켄트 벡은 이를 "꾸며대기(Faking it)" 라고 불렀다. 테스트를 통과시키고 난 다음에는 테스트 코드를 추가하여 하드코딩 만으로는 통과하지 못하게 만든다. 

이렇게 작은 단계를 수행하는 것이, 어플리케이션의 로직에 대한 큰 걱정 없이 전체 프로젝트의 구조가 정확하게 작동하게 만드는 중요한 시작이 된다. 

Go에서 웹서버를 생성하려면 [ListenAndServe](https://golang.org/pkg/net/http/#ListenAndServe)를 호출하면 된다. 

```go
func ListenAndServe(addr string, handler Handler) error
```

특정 포트를 리스닝하는 웹서버가 시작되며, 모든 request에 대해 고루틴이 생성되며 그 위에서 [`Handler`](https://golang.org/pkg/net/http/#Handler)가 실행된다. 

```go
type Handler interface {
	ServeHTTP(ResponseWriter, *Request)
}
```

하나의 타입이 `ServeHTTP` 메서드를 구현하면 Handler 인터페이스를 구현한 것이다. `ServeHTTP` 메서드는 두 개의 인자를 가지는데 첫번째는 _response를 쓰는 곳_ 이고, 두번째는 서버가 받은 HTTP request이다. 

`server_test.go` 라는 파일을 만들고 `PlayerServer`라는 함수를 테스트하는 코드를 작성하자. `PlayerServer` 함수는 두 개의 인자를 가진다. request는 player의 승점 `"20"`을 받아야 한다. 

```go
func TestGETPlayers(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		got := response.Body.String()
		want := "20"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
```

서버를 테스트하려면 `Request`를 서버로 보내고 handler가 `ResponseWriter`에 무엇을 쓰는지 알아야 한다. 

- `http.NewRequest`로 request를 만들었다. 첫번째 인자는 request의 method(ex. "GET", "POST", etc.)이고, 두번째는 request의 경로인데 이 인자가 `nil`이라는 것은 request의 body에 아무것도 없기 때문이다. 
- `net/http/httptest`패키지에는 `ResponseRecorder`라는 스파이가 이미 있으며, response에 무엇을 썼는지 분석할 수 있는 유용한 메서드가 많다. 


## Try to run the test

`./server_test.go:13:2: undefined: PlayerServer`

## Write the minimal amount of code for the test to run and check the failing test output

컴파일러의 에러 출력만 보아도 문제를 해결할 수 있다. 

`server.go` 파일을 생성하고 `PlayerServer`를 정의하자. 

```go
func PlayerServer() {}
```

다시 시도해보자.

```
./server_test.go:13:14: too many arguments in call to PlayerServer
    have (*httptest.ResponseRecorder, *http.Request)
    want ()
```

함수에 인자를 추가하자.

```go
import "net/http"

func PlayerServer(w http.ResponseWriter, r *http.Request) {

}
```

컴파일에 성공하고, 테스트가 실패할 것이다. 

```
=== RUN   TestGETPlayers/returns_Pepper's_score
    --- FAIL: TestGETPlayers/returns_Pepper's_score (0.00s)
        server_test.go:20: got '', want '20'
```

## Write enough code to make it pass

DI 장에서 `Greet` 함수를 가진 HTTP 서버를 짰던 기억이 날 것이다. net/http 패키지의 `ResponseWriter`는 `Writer`도 구현되어 있다. 따라서 `fmt.Fprintf`를 이용해 문자열을 HTTP response로 보낼 수 있다. 

```go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}
```

이제 테스트를 통과할 것이다.

## 비계(scaffolding)를 완성하라

이제는 실제 애플리케이션으로 연결해야 한다. 이것이 중요한 이유는 (번역: 이 부분 이해가 잘 안됨)

- _실제 작동하는 소프트웨어_ 를 가지게 될 것이고, 이를 위한 테스트를 짜지는 않을 것이다. 작동하는 코드를 보는건 좋다. 
- 리팩터링을 하는 건, 프로그램의 구조를 바꾸는 것과 같다. 변경사항는 점진적인 개발의 하나로서 애플리케이션에 반영될 것이다. 

`main.go` 파일을 만들고 코드를 작성하자. 

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	handler := http.HandlerFunc(PlayerServer)
	if err := http.ListenAndServe(":5000", handler); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

현 시점에서 애플리케이션은 하나의 파일로 구현되어 있다. 하지만 더 큰 프로젝트에서는 여러 파일로 나누고 싶어질 것이다. 

애플리케이션을 실행하려면, `go build` 명령으로 디렉토리 안의 모든 `.go` 파일들로 프로그램을 빌드한다음 `./myprogram` 을 실행하면 된다. 


### `http.HandlerFunc`

앞서서 서버를 만들려면 `Handler` 인터페이스가 필요하다고 했었다. _일반적으로_ `struct`를 만들고 ServeHTTP 메서드를 구현하여 인터페이스를 구현한다. 하지만 struct의 용도는 데이터를 담는 것인데 _현재는_ 아무 state가 없기에 struct를 만들기 머뭇거려진다. 

[HandlerFunc](https://golang.org/pkg/net/http/#HandlerFunc)를 사용해서 이 문제를 비껴갈 수 있다. 

> HandlerFunc 타입은 평범한 함수들을 HTTP 핸들러로 쓸 수 있게 해주는 어댑터이다. 만약에 f가 적합한 시그니처를 가진 함수라면, HandlerFunc(f)는 f를 호출하는 핸들러이다. (역주: 여기서 f가 타입 HandlerFunc로 타입 컨버젼이 되었다.)

```go
type HandlerFunc func(ResponseWriter, *Request)
```

문서를 보면 `HandlerFunc`는 이미 `ServeHTTP` 메서드가 구현되어 있다. 
`PlayerServer`를 `HandlerFunc`로 타입 컨버젼하면 `Handler`를 구현한 셈이 된다. 

### `http.ListenAndServe(":5000"...)`

`ListenAndServe`는 `Handler`가 리스닝할 포트를 지정한다. 이미 리스닝중인 포트라면 `error`를 리턴한다. 에러는 `if`문을 이용해서 에러를 잡고 로깅을 할 수 있다. 

_또 다른_ 테스트를 작성해서 하드 코딩된 값보다 나은 구현을 해보자. 

## Write the test first

다른 player의 승점을 확인하는 테스트를 작성할 것이다. 이 테스트는 하드코딩한 코드로는 통과하지 못한다. 

```go
t.Run("returns Floyd's score", func(t *testing.T) {
	request, _ := http.NewRequest(http.MethodGet, "/players/Floyd", nil)
	response := httptest.NewRecorder()

	PlayerServer(response, request)

	got := response.Body.String()
	want := "10"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
})
```

이런 생각이 들 수 있다. 

> 어떤 player의 승점이 얼마인지를 제어할 수 있는 저장소 개념이 필요하다. 테스트를 할 때에 임의의 값이 들어가는 건 어색하다


논리적으로 가능한 최소의 단계를 밟으려 하고 있다는 걸 잊지 말자. 일단은 상수값을 리턴하는 것을 개선하는 것에만 집중하자. 


## Try to run the test

```
=== RUN   TestGETPlayers/returns_Pepper's_score
    --- PASS: TestGETPlayers/returns_Pepper's_score (0.00s)
=== RUN   TestGETPlayers/returns_Floyd's_score
    --- FAIL: TestGETPlayers/returns_Floyd's_score (0.00s)
        server_test.go:34: got '20', want '10'
```

## Write enough code to make it pass

```go
//server.go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	if player == "Pepper" {
		fmt.Fprint(w, "20")
		return
	}

	if player == "Floyd" {
		fmt.Fprint(w, "10")
		return
	}
}
```

request의 URL을 보고 어떤 값을 리턴할지를 결정하게 하였다. player 승점의 저장 및 연동을 고려하다 보면 다음 단계는 _라우팅_ 이 될 것같다.

변경된 만큼 승점을 저장하게 하려 했다면 이보다 훨씬 많은 수정을 해야 했을 것이다. **하지만 이런 구현이 우리의 최종 목표를 향한 훨씬 작은, 테스트를 기반으로한 단계이다."

지금은 라우팅 라이브러리를 이용하고 싶은 유혹을 참아내고, 테스트를 통과하는 최소의 단계만을 생각하자. 

`r.URL.Path` 는 request의 경로를 리턴하며, 우리는 [`strings.TrimPrefix`](https://golang.org/pkg/strings/#TrimPrefix) 를 이용해 `/players/` 를 잘라내어 요청한 player만을 얻을 수 있다. 단단한 코드라 볼 수는 없지만 당장은 동작한다. 

## Refactor

승점을 가져오는 부분을 별도 함수로 추출하여 `PlayerServer`를 단순화 시켜보자. 

```go
//server.go
func PlayerServer(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	fmt.Fprint(w, GetPlayerScore(player))
}

func GetPlayerScore(name string) string {
	if name == "Pepper" {
		return "20"
	}

	if name == "Floyd" {
		return "10"
	}

	return ""
}
```

테스트에서 helper를 이용해 반복을 줄일 수 있다. (DRY up - Don't repeat yourself)

```go
//server_test.go
func TestGETPlayers(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		assertResponseBody(t, response.Body.String(), "10")
	})
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
```

물론 제대로 된 구현은 아니다. 서버가 점수를 알고 있다는건 아무래도 이상하다. 

리팩터링을 진행하다 보면 무엇을 개선해야 할지 보인다. 

승점 계산 부분을 main에서 `GetPlayerScore` 함수로 옮겼지만, 함수보다는 인터페이스를 이용하는게 적절해 보인다. 

리팩터링하여 옮긴 함수를 인터페이스로 만들어 보자. 

```go
type PlayerStore interface {
	GetPlayerScore(name string) int
}
```

`PlayerServer` 가 `PlayerStore`를 쓸 수 있으려면 레퍼런스가 필요하다. `PlayerServer`가 구조체가 되게 아키텍처를 바꿀 시점이다. 

```go
type PlayerServer struct {
	store PlayerStore
}
```

새로운 구조체에 메서드를 추가해서, `Handler` 인터페이스를 구현하고, 핸들러 코드를 넣어주자. 

```go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

유일한 차이점은 로컬 함수(삭제 예정)가 아닌 `store.GetPlayserScore` 메서드를 호출한다는 것이다. 

서버의 전체코드를 보자. 

```go
//server.go
type PlayerStore interface {
	GetPlayerScore(name string) int
}

type PlayerServer struct {
	store PlayerStore
}

func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

### 이슈 수정

제법 수정을 하고 보니 컴파일이 안될 것이다. 우선 컴파일이 되도록 수정하자. 

`./main.go:9:58: type PlayerServer is not an expression`

테스트에서 `PlayerServer`를 생성하는 것이 아니라 메서드인 `ServerHTTP`를 호출하여야 한다. 

```go
//server_test.go
func TestGETPlayers(t *testing.T) {
	server := &PlayerServer{}

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "10")
	})
}
```

_아직까지도_ 저장소를 만들지 않았다는 것에 주목하자. 어떻게든 컴파일 성공부터 하는 것이다. 

컴파일이 되도록 한 다음에, 테스트를 통과하게 하는 거다. 이 순서대로 코딩하는 습관이 몸에 베어야 한다. 

컴파일도 되지 않았는데 (stub 저장소 같은) 기능을 추가하는 것은 _훨씬_ 복잡한 컴파일 문제를 만들 수 있다. 

이제 `main.go` 는 컴파일 되지 않을것이다. 

```go
func main() {
	server := &PlayerServer{}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

마침내, 컴파일에 성공하지만, 이번에는 테스트를 실패한다. 

```
=== RUN   TestGETPlayers/returns_the_Pepper's_score
panic: runtime error: invalid memory address or nil pointer dereference [recovered]
    panic: runtime error: invalid memory address or nil pointer dereference
```

아직 테스트에 `PlayerStore`를 전달하지 않았기 때문이다. stub를 만들 차례다. 

```go
//server_test.go
type StubPlayerStore struct {
	scores map[string]int
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}
```

테스트를 위해 `map`으로 빠르고 쉬운 stub key/value 저장소를 만들 수 있다. 저장소를 만들고 `PlayerServer`로 전달하자. 

```go
//server_test.go
func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := &PlayerServer{&store}

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertResponseBody(t, response.Body.String(), "10")
	})
}
```

테스트를 통과했고 코드도 보기 좋아졌다. 저장소 덕분에 코드의 _의도_ 가 선명해졌다. _이런 데이터가 `PlayerStore`_ 에 있으니 `PlayerServer` 를 이용해 원하는 response를 받을 수 있다고 말해주는 것이다. 

### 애플리케이션 실행

모든 테스트를 통과했으니 리팩터링을 완료하기 위해 애플리케이션의 동작을 확인해보자. 프로그램은 시작하겠지만 `http://localhost:5000/players/Pepper` 로 request를 하면 끔찍한 response를 받을 것이다. 

`PlayerStore`를 전달하지 않았기 때문이다. 

아직 의미있는 데이터를 저장하지 않고 있기에 `PlayerStore` 구현은 조금 곤란하다. 우선은 하드코딩을 해두자. 

```go
//main.go
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return 123
}

func main() {
	server := &PlayerServer{&InMemoryPlayerStore{}}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

`go build` 를 실행하고 `http://localhost:5000/players/Pepper` URL로 request 하면 `"123"`이 회신된다. 멋지진 않지만 현재로선 이게 최선이다. 

다음에 할 만한 것들은 다음과 같다. 

- player가 존재하지 않을 경우의 처리
- `POST /players/{name}` 에 대한 처리
- 메인 애플리케이션이 시작했지만 실제 동작하지 않아서 불편함. 문제점을 확인하려면 매번 테스트를 실행하여야 한다. 

`POST` 처리를 하고 싶지만, 존재하지 않는 player 처리를 먼저하는게 지금까지 구현한 것과 연관도 있어서 적절하게 느껴진다. 나머지는 이후에 구현한다. 

## Write the test first

존재하지 않는 player 처리 테스트를 추가한다.

```go
//server_test.go
t.Run("returns 404 on missing players", func(t *testing.T) {
	request := newGetScoreRequest("Apollo")
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	got := response.Code
	want := http.StatusNotFound

	if got != want {
		t.Errorf("got status %d want %d", got, want)
	}
})
```

## Try to run the test

```
=== RUN   TestGETPlayers/returns_404_on_missing_players
    --- FAIL: TestGETPlayers/returns_404_on_missing_players (0.00s)
        server_test.go:56: got status 200 want 404
```

## Write enough code to make it pass

```go
//server.go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	w.WriteHeader(http.StatusNotFound)

	fmt.Fprint(w, p.store.GetPlayerScore(player))
}
```

때로는 TDD 신봉자들이 "테스트를 통과할 최소한의 코드만 짜라"고 하는게 지나치게 현학적으로 느껴져서 눈이 동그래진다. 

이 코드가 매우 좋은 예이다. 정말 최소한의 코드만 짰다.(그리고 올바른 구현도 아니다). **모든 response**를 `StatusNotFound` 로 보내버리는 것이다. 그럼에도 불구하고 모든 테스트를 통과한다!

**이렇게 최소한의 코드를 짜서 테스트를 통과하면, 테스트들 간의 차이를 선명하게 볼 수 있다.** 여기서는 player가 _존재_ 한다면, `StatusOK`를 받아야 한다는 것을 단정하지 않은 것이다. 

다른 두 테스트가 status를 체크하도록 수정하자. 

아래가 새로이 수정한 테스트이다. 

```go
//server_test.go
func TestGETPlayers(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{
			"Pepper": 20,
			"Floyd":  10,
		},
	}
	server := &PlayerServer{&store}

	t.Run("returns Pepper's score", func(t *testing.T) {
		request := newGetScoreRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "20")
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request := newGetScoreRequest("Floyd")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
		assertResponseBody(t, response.Body.String(), "10")
	})

	t.Run("returns 404 on missing players", func(t *testing.T) {
		request := newGetScoreRequest("Apollo")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusNotFound)
	})
}

func assertStatus(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Errorf("did not get correct status, got %d, want %d", got, want)
	}
}

func newGetScoreRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/players/%s", name), nil)
	return req
}

func assertResponseBody(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
```

모든 테스트의 상태를 체크하는데 도움이 될 `assertStatus` 함수를 만들었다. 

첫 두 개의 테스트는 200이 아닌 404를 받아서 실패한다. `PlayerServer`를 수정해서 승점이 0이면 찾지 못했다고 `StatusNotFound`를 회신하게 하자. 


```go
//server.go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}
```

### 승점 저장

저장소에서 승점을 가져올 수 있게 되었으니 이제 새로운 승점을 저장할 수 있게 만들어보자.

## Write the test first

```go
//server_test.go
func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
	}
	server := &PlayerServer{&store}

	t.Run("it returns accepted on POST", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodPost, "/players/Pepper", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)
	})
}
```

특정 라우트로 POST를 보낼 경우에, 올바른 status code를 받는지부터 확인하자. `GET /players/{name}` 와는 다른 종류의 request를 받아서 처리하는 기능을 구현해야 한다. 이게 동작하면 핸들러에서 승점과 연동하는 분을 확인할 것이다. 

## Try to run the test

```
=== RUN   TestStoreWins/it_returns_accepted_on_POST
    --- FAIL: TestStoreWins/it_returns_accepted_on_POST (0.00s)
        server_test.go:70: did not get correct status, got 404, want 202
```

## Write enough code to make it pass

테스트부터 만드는 것은 신중하게 문제를 만드는 것이다. `if` 문으로 request의 method를 구분하여 해결해보자. 

```go
//server.go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		w.WriteHeader(http.StatusAccepted)
		return
	}

	player := strings.TrimPrefix(r.URL.Path, "/players/")

	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}
```

## Refactor

핸들러가 지저분하게 구현되어 있다. 코드를 나누어 알아보기 편하게 함수들로 만들자. 

```go
//server.go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodPost:
		p.processWin(w)
	case http.MethodGet:
		p.showScore(w, r)
	}

}

func (p *PlayerServer) showScore(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter) {
	w.WriteHeader(http.StatusAccepted)
}
```

`ServeHTTP`의 라우팅이 좀더 잘 이해된다. 다음 반복에는 `processWin` 함수 내부의 저장부분을 구현한다.

그 다음엔 서버가 `POST /players/{name}`를 받으면 `PlayerStore`가 승점을 저장하라는 요청을 듣는지 체크할 것이다. 

## Write the test first

`RecordWin` 메서드를 `StubPlayerStore`에 추가한 다음 호출해보자. 

```go
//server_test.go
type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}
```

이번에는 호출 횟수를 확인하는 테스트를 가장 처음에 추가해보자.
Now extend our test to check the number of invocations for a start

```go
//server_test.go
func TestStoreWins(t *testing.T) {
	store := StubPlayerStore{
		map[string]int{},
	}
	server := &PlayerServer{&store}

	t.Run("it records wins when POST", func(t *testing.T) {
		request := newPostWinRequest("Pepper")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusAccepted)

		if len(store.winCalls) != 1 {
			t.Errorf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
		}
	})
}

func newPostWinRequest(name string) *http.Request {
	req, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("/players/%s", name), nil)
	return req
}
```

## Try to run the test

```
./server_test.go:26:20: too few values in struct initializer
./server_test.go:65:20: too few values in struct initializer
```

## Write the minimal amount of code for the test to run and check the failing test output

`StubPlayerStore`에 필드를 추가했으니, 생성할때의 코드도 수정한다. 

```go
//server_test.go
store := StubPlayerStore{
	map[string]int{},
	nil,
}
```

```
--- FAIL: TestStoreWins (0.00s)
    --- FAIL: TestStoreWins/it_records_wins_when_POST (0.00s)
        server_test.go:80: got 0 calls to RecordWin want 1
```

## Write enough code to make it pass

정확한 값이 아니라 호출 횟수만을 단정하기에 첫 반복은 간단했다.

만약 `RecordWin`을 호출할 수 있게 되면, 인터페이스를 변경해서 `PlayerStore`의 개념에 대해 `PlayerServer`를 수정할 필요가 있다. 

```go
//server.go
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
}
```

이렇게 하면 `main`은 컴파일되지 않는다.

```
./main.go:17:46: cannot use InMemoryPlayerStore literal (type *InMemoryPlayerStore) as type PlayerStore in field value:
    *InMemoryPlayerStore does not implement PlayerStore (missing RecordWin method)
```

컴파일러는 무엇이 문제인지 알려준다. `InMemoryPlayerStore`에 `RecordWin` 메서드를 추가해주자.

```go
//main.go
type InMemoryPlayerStore struct{}

func (i *InMemoryPlayerStore) RecordWin(name string) {}
```

테스트해보면 컴파일은 성공하지만 테스트는 실패한다. 

이제 `PlayerStore`가 `RecordWin` 메서드를 가지고 있으니 `PlayerServer`에서 호출할 수 있다. 

```go
//server.go
func (p *PlayerServer) processWin(w http.ResponseWriter) {
	p.store.RecordWin("Bob")
	w.WriteHeader(http.StatusAccepted)
}
```

테스트를 실행하면 통과할 것이다. 하지만 `RecordWin`에 넣으려는 이름이 `"Bob"`은 아니었으니 테스트를 좀더 다듬어보자. 

## Write the test first

```go
//server_test.go
t.Run("it records wins on POST", func(t *testing.T) {
	player := "Pepper"

	request := newPostWinRequest(player)
	response := httptest.NewRecorder()

	server.ServeHTTP(response, request)

	assertStatus(t, response.Code, http.StatusAccepted)

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != player {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], player)
	}
})
```

`winCalls` 슬라이스에는 하나의 원소가 있어야 하고, 그 원소가 `player`와 같아야 테스트를 통과한다. 

## Try to run the test

```
=== RUN   TestStoreWins/it_records_wins_on_POST
    --- FAIL: TestStoreWins/it_records_wins_on_POST (0.00s)
        server_test.go:86: did not store correct winner got 'Bob' want 'Pepper'
```

## Write enough code to make it pass

```go
//server.go
func (p *PlayerServer) processWin(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
```

`processWin` 메서드의 코드를 수정해서 `http.Request`를 받아 URL에서 player 이름을 추출하게 하였다. 이제 `store` 의 `RecordWin` 메서드를 player 이름으로 호출하고 테스트를 통과할 것이다. 

## Refactor

DRY(Don't Repeat Yourself). 반복되는 코드를 줄여보자. player 이름을 추출하는 코드를 `ServeHTTP`로 옮겼다. 

```go
//server.go
func (p *PlayerServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	player := strings.TrimPrefix(r.URL.Path, "/players/")

	switch r.Method {
	case http.MethodPost:
		p.processWin(w, player)
	case http.MethodGet:
		p.showScore(w, player)
	}
}

func (p *PlayerServer) showScore(w http.ResponseWriter, player string) {
	score := p.store.GetPlayerScore(player)

	if score == 0 {
		w.WriteHeader(http.StatusNotFound)
	}

	fmt.Fprint(w, score)
}

func (p *PlayerServer) processWin(w http.ResponseWriter, player string) {
	p.store.RecordWin(player)
	w.WriteHeader(http.StatusAccepted)
}
```

테스트는 통과했지만 아직은 소프트웨어가 작동하지 않는다. `PlayStore`를 제대로 구현하지 않았기 때문이다. 하지만 핸들러에 집중했기에 어떤 인터페이스가 필요한지 명확히 할 수 있었다. 실행 시작부에서부터 디자인했다면 쉽지 않았을 것이다. 

`InMemoryPlayerStore` 부터 _테스트를 짤 수도 있었다._ 하지만 `InMemoryPlayerStore`는, 데이터베이스와 같이, 제대로 승점을 저장하도록 변경할 때까지 임시로 사용하는 것이다. 


이제 `PlayerServer`와 `InMemoryPlayerStore` 사이의 _integration test_ 를 짜서, 기능을 끝낼 것이다. 이 테스트를 통해, `InMemoryPlayStore`를 바로 테스트 하는 것과 달리, 애플리케이션이 제대로 동작한다는 확신을 얻을 수 있다. 그 뿐 아니라, `PlayStore`를 데이터베이스로 구현하게 될 때에 같은 integration test로 테스트 할 수 있다. 


### 통합 테스트

통합 테스트는 시스템의 큰 범위를 테스트하기에 유용하지만 염두에 둘 것이 있다. 

-   작성하기 어렵다.
-   실패하면 원인을 알기 어렵기에 수정도 어럽다. (통합 테스트의 컴포넌트 사이의 버그인 경우가 많다)
-   테스트 수행이 느린 경우가 있다. (데이터베이스와 같은 "진짜" 컴포넌트를 사용하기 때문이다)

이런 이유로 _테스트 피라미드_ 를 알아볼 것을 추천한다. 

## Write the test first

간결하게, 리팩터링이 끝난 최종 통합테스트를 보여주겠다. 

```go
//server_integration_test.go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
	store := InMemoryPlayerStore{}
	server := PlayerServer{&store}
	player := "Pepper"

	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))
	server.ServeHTTP(httptest.NewRecorder(), newPostWinRequest(player))

	response := httptest.NewRecorder()
	server.ServeHTTP(response, newGetScoreRequest(player))
	assertStatus(t, response.Code, http.StatusOK)

	assertResponseBody(t, response.Body.String(), "3")
}
```

-   통합하려는 두 개의 컴포넌트를 생성한다. `InMemoryPlayerStore` 와 `PlayerServer`.
-   세 개의 request를 보내어 `player`의 세 개의 승점을 저장하게 한다. 통합이 잘 되었는지 여부와는 무관하기에 Status code는 일단 무시하자. 
-   `response` 변수에 response를 저장해서 `player`의 승점을 확인한다.

## Try to run the test

```
--- FAIL: TestRecordingWinsAndRetrievingThem (0.00s)
    server_integration_test.go:24: response body is wrong, got '123' want '3'
```

## Write enough code to make it pass

조금의 자유를 누려보자. 테스트 없이는 부담스러울 정도의 코드를 짜본다. 

_이렇게 할 수도 있다!_ 제대로 동작하는지 확인하는 테스트들이 있긴 하지만, 우리가 작업해온 `InMemoryPlayerStore` 와는 상관이 없다. 

이러다가 구현이 꼬여버렸다면 이전의 마지막으로 되돌리면 된다. 그러고 다시 `InMemoryPlayerStore` 주위의 구체적인 유닛 테스트를 좀 더 짜보면서 해법을 찾아내자. 

```go
//in_memory_player_store.go
func NewInMemoryPlayerStore() *InMemoryPlayerStore {
	return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct {
	store map[string]int
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
	i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
	return i.store[name]
}
```

-   데이터를 저장해야 하기에 `map[string]int` 를 `InMemoryPlayerStore` 구조체에 추가했다.
-   편의를 위해 저장소를 초기화하는 `NewInMemoryPlayerStore` 를 추가하고, 통합 테스트가 이를 사용하게 수정한다. 
  
```go
//server_integration_test.go
store := NewInMemoryPlayerStore()
server := PlayerServer{store}
```

-   나머지 코드는 `map`을 감싼 것이다. 

통합 테스트를 통과했다. 이제 `main`이 `NewInMemoryPlayStore()`를 사용하게 바꿔주기만 하면 된다. 


```go
//main.go
package main

import (
	"log"
	"net/http"
)

func main() {
	server := &PlayerServer{NewInMemoryPlayerStore()}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

빌드하고 실행한 다음, `curl`로 테스트 해보자. 

-   `curl -X POST http://localhost:5000/players/Pepper` 를 여러 번 실행한다. player 이름을 바꿔가며 해도 좋다.
-   `curl http://localhost:5000/players/Pepper` 로 승점을 확인하자 

지금까지 잘 해왔다. REST같은 서비스를 만들었다. 좀 더 개선해서 프로그램 실행 이후에도 데이터를 저장할 수 있게 할 수 있겠다. 

-   저장소를 고른다. (Bolt? Mongo? Postgres? File system?)
-   `PostgresPlayerStore`를 만들고 `PlayerStore`를 구현한다.
-   기능을 TDD로 개발하여 작동을 확인한다. 
-   통합 테스트에 추가하여 테스트를 통과하는지 확인한다.
-   마지막으로 `main`에 통합해준다. 

## Refactor

거의 끝나간다. 아래와 같은 동시성 문제가 나지 않도록 대비하자

```
fatal error: concurrent map read and map write
```

뮤텍스를 추가해서 동시성에 안전하게 만들자. 특히 `RecordWin` 함수의 counter 를 챙기자. sync 장에서 mutexes 에 대해 더 알아보자. 

## Wrapping up

### `http.Handler`

-   웹서버를 만들기 위해 다음의 인터페이스를 구현한다.
-   `http.HandlerFunc`를 이용해 일반적인 함수를 `http.Handler`로 쓴다. 
-   `httptest.NewRecorder`를 `ResponseWriter` 로써 전달하여 핸들러의 response를 훔쳐본다. 
-   `http.NewRequest`를 이용하여 시스템으로 들어올 request를 생성한다. 

### 인터페이스, 목업, 그리고 DI(Dependency Inversion)

-   조금씩 반복해가며 시스템을 만든다. 
-   실제 저장소없이, 저장소가 필요한 핸들러를 만든다. 
-   
-   Allows you to develop a handler that needs a storage without needing actual storage
-   TDD를 통해 원하는 인터페이스를 만들어낸다. 

### 문제를 만들고, 리팩터링하기 (그리고 소스 관리 시스템에 commit 한다)

-   컴파일 실패와 테스트 실패를 적색 경보라 생각하고, 최대한 빠르게 빠져나온다.
-   적색 경보에서 빠져나올 수 있는 최소한의 코드를 짠다. _그리고 나서_ 리팩터링하고 코드를 다듬는다. 
-   컴파일이 안되고, 테스트가 실패하는 동안에, 지나치게 많은 변경을 하면 문제가 복잡해질 위험이 커진다. 
-   이러한 문제 해결법을 고수하면, 작은 테스트를 짤 수 밖에 없고, 그 결과 작은 수정만을 하게 된다. 그 결과로 복잡한 시스템에서도 꾸준히 안정적으로 개발을 할 수 있다. 
