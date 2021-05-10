# Command line and project structure

**[이 챕터의 모든 코드는 여기에서 확인할 수 있다.](https://github.com/quii/learn-go-with-tests/tree/main/command-line)**

Product Owner는 command line에서 작동하는 두번째 애플리케이션을 도입해서 _피봇_ 하고자 한다. (스타트업 기업들이 초기의 사업의 목표나 서비스 운영방식등을 바꿔 다른 사업으로 이전하는 것)

우선 당장은 사용자가 `Ruth wins`라고 입력하면 그 선수의 승리가 기록되기만 하면 된다. 궁극적으로는 이를 사용해 사용자들이 포커를 칠 수 있도록 돕는 프로그램을 만드는 것이다.

Product Owner는 리그가 새로운 애플리케이션에 기록된 승수에 따라 업데이트할 수록 두 애플리케이션이 서로 공유하고 있는 데이터베이스가 있기를 원한다.

## 코드 다시 한 번 보기

HTTP 서버를 실행하는 `main.go`라는 파일이 있는 애플리케이션이 있다. 이번 챕터에서 다룰 내용을 위해서 만들었던 HTTP 서버 자체에는 관심이 없지만 이 서버가 추상화한 `PlayerStore`에는 관심이 있다.

```go
type PlayerStore interface {
	GetPlayerScore(name string) int
	RecordWin(name string)
	GetLeague() League
}
```

이전 챕터에서 우리는 `PlayerStore` 인터페이스로 구현한 `FileSystemPlayerStore`를 만들었다. 이 중 일부를 새로 만들 애플리케이션에서 재사용할 수 있어야 한다.

## 먼저 프로젝트를 리팩터링 해보자

이제 우리 프로젝트는 기존의 웹 서버와 command line 애플리케이션, 총 두 개의 binary가 필요하다.

우리는 새로운 일에 몰두하기 전에 두 개의 binary가 있을 수 있도록 프로젝트를 먼저 구성해야 한다.

지금까지 모든 코드는 한 폴더에 있었고 주소는 아래와 같았다.

`$GOPATH/src/github.com/your-name/my-app`

Go를 사용해서 애플리케이션을 만들기 위해서는, `package main` 안에 `main` 함수가 있어야 한다. 지금까지 우리의 모든 도메인 코드는 `package main`안에 있었고 `func main`은 모든 것을 참조할 수 있었다.

지금까지는 괜찮았고 패키지 구조로만 본다면 너무 지나치지 않은 좋은 습관이다. 표준 라이브러리를 천천히 훑어보면 많은 폴더와 구조에서는 거의 보이지 않을 것이다.

감사하게도, _우리가 필요할 때_ 구조를 추가하는 것은 매우 간단하다.

현재 프로젝트 안에는 `cmd` 폴더와 그 안에 `webserver`폴더를 같이 만들어라 (예시: `mkdir -p cmd/webserver`).

그리고 그 안에 `main.go` 파일을 이동시켜라.

`tree`가 설치되어 있다면 돌렸을 때 아래와 같은 폴더 구조를 가져야 한다.

```
.
├── file_system_store.go
├── file_system_store_test.go
├── cmd
│   └── webserver
│       └── main.go
├── league.go
├── server.go
├── server_integration_test.go
├── server_test.go
├── tape.go
└── tape_test.go
```

이제 애플리케이션과 라이브러리 코드가 효과적으로 분리되어 있지만, 패키지 이름을 몇 개 변경해야 한다. 우리가 Go 애플리케이션을 빌드할 때는 그 패키지는 _무조건_ `main`이어야만 하는 것을 기억하자.

다른 모든 코드가 `poker` 패키지 안에 있도록 바꾸자.

마지막으로 이 패키지를 `main.go`에서 불러와서 그 패키지를 사용해서 웹 서버를 만들 수 있다. 그 때 라이브러리 코드는 `poker.FunctionName`과 같이 사용할 수 있다

패키지가 저장되어있는 주소들은 당신의 컴퓨터에서는 다를 수 있겠지만 이와 비슷해야만 한다.

```go
//cmd/webserver/main.go
package main

import (
	"github.com/quii/learn-go-with-tests/command-line/v1"
	"log"
	"net/http"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

전체 경로를 다 적어야 하는 것은 약간 불편해 보일 수 있지만 이 방법으로 공개적으로 사용가능한 라이브러리를 불러올 수 있다.

도메인 코드를 별도의 패키지로 분리하고 깃허브와 같은 공개적인 리포지토리에 커밋함으로써 Go 개발자들은 우리가 작성한 기능들이 있는 패키지를 불러오는 자신들의 코드를 작성할 수 있다. 너가 처음 이 코드를 돌리면 그 패키지들이 없다고 에러 메세지가 뜨겠지만 당신은 `go get`을 실행시키기만 하면 된다.

추가로 사용자들은 [godoc.org에 있는 문서를](https://godoc.org/github.com/quii/learn-go-with-tests/command-line/v1) 볼 수도 있다.

### 최종 확인

-   프로젝트 최상위 폴더로 가서 `go test`를 실행시켜서 아직도 모든 테스트가 통과하는지 확인하자.
-   `cmd/webserver`로 가서 `go run main.go`를 실행시키자.
    -   `http://localhost:5000/league`로 들어가서 아직도 작동하는 지 확인하자.

### 코드 구조 훑어보기

테스트를 작성하기 전에 우리 프로젝트가 빌드할 새로운 애플리케이션을 추가하자. `cmd` 폴더 안에 `cli`라는 새로운 폴더를 만들고 그 안에 아래와 같이 `main.go` 파일을 추가하자.

```go
//cmd/cli/main.go
package main

import "fmt"

func main() {
	fmt.Println("Let's play poker")
}
```

우리가 가장 먼저 할 요구사항은 사용자가 `{PlayerName} wins`를 입력할 때 승리를 기록하는 것이다.

## 테스트부터 작성하기

우선 포커를 `Play`할 수 있게 해주는 `CLI`라는 것을 만들 필요가 있다는 것을 안다. 그것은 사용자가 입력한 값을 읽고 `PlayerStore`에 승리를 기록해야 한다.

너무 들어가기 이전에 우리가 원하는 `PlayerStore`와 통합이 되는지를 체크하는 테스트를 작성해보자.

`CLI_test.go` 안을 보면 (`cmd`폴더 안이 아닌 프로젝트 루트 폴더)

```go
//CLI_test.go
package poker

import "testing"

func TestCLI(t *testing.T) {
	playerStore := &StubPlayerStore{}
	cli := &CLI{playerStore}
	cli.PlayPoker()

	if len(playerStore.winCalls) != 1 {
		t.Fatal("expected a win call but didn't get any")
	}
}
```

-   우리는 다른 테스트에서 `StubPlayerStore`을 사용할 수 있다.
-   우리는 아직 존재하지 않는 `CLI` 타입에 의존 변수를 제공한다.
-   아직 작성되지 않은 `PlayPoker` 메소드를 호출한다.
-   승리가 기록되었는지 체크한다.

## 테스트 실행해보기

```
# github.com/quii/learn-go-with-tests/command-line/v2
./cli_test.go:25:10: undefined: CLI
```

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

이 쯤 되면 당신은 의존 변수를 위해 각각의 필드가 있는 `CLI` 구조체와 메소드를 추가할만큼 충분히 편안해야 한다.

당신의 코드는 아래와 같이 작성되어야 한다.

```go
//CLI.go
package poker

type CLI struct {
	playerStore PlayerStore
}

func (cli *CLI) PlayPoker() {}
```

우리는 테스트를 실패하는지 체크할 만큼만 테스트를 작성하면 된다는 것을 기억하자.

```
--- FAIL: TestCLI (0.00s)
    cli_test.go:30: expected a win call but didn't get any
FAIL
```

## 테스트를 통과하는 최소한의 코드를 작성하기

```go
//CLI.go
func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Cleo")
}
```

위의 코드는 테스트를 통과해야만 한다.

다음으로 우리는 `Stdin` (유저로부터 들어오는 입력값)을 읽는 것을 시뮬레이션하기 위해서 우리는 특정 선수들의 승리를 기록할 수 있어야 한다.

이것을 하기 위해서 우리 테스트 코드를 좀 더 확장시켜보자.

## 테스트부터 작성하기

```go
//CLI_test.go
func TestCLI(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &StubPlayerStore{}

	cli := &CLI{playerStore, in}
	cli.PlayPoker()

	if len(playerStore.winCalls) < 1 {
		t.Fatal("expected a win call but didn't get any")
	}

	got := playerStore.winCalls[0]
	want := "Chris"

	if got != want {
		t.Errorf("didn't record correct winner, got %q, want %q", got, want)
	}
}
```

`os.stdin`은 `main`에서 사용자의 입력값을 가져오기 위해서 사용할 패키지다. 속에 들여다보면 `*File`인데 이는 우리가 알고 있듯이 입력값을 받아오는 편리한 방법으로 `io.Reader`를 구현한다는 뜻이다.

우리는 사용자가 입력하기로 예상한 값을 채우기 위해 `strings.NewReader`를 사용해서 테스트 안에 `io.Reader`를 만든다.

## 테스트 실행해보기

`./CLI_test.go:12:32: too many values in struct initializer`

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

`CLI`안에 새로운 의존 변수들을 넣어야 한다.

```go
//CLI.go
type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}
```

## 테스트를 통과하는 최소한의 코드 작성하기

```
--- FAIL: TestCLI (0.00s)
    CLI_test.go:23: didn't record the correct winner, got 'Cleo', want 'Chris'
FAIL
```

가장 쉬운 일부터 하는 것을 기억하자.

```go
func (cli *CLI) PlayPoker() {
	cli.playerStore.RecordWin("Chris")
}
```

테스트는 통과했다. 다음으로 진짜 코드를 작성하기 위해서 다른 테스트들을 넣기 전에 먼저 리팩터링부터 해보자.

## 리팩터링 하기

`server_test.go`에서 우리는 승리가 기록되는지를 체크하는 코드를 이미 작성했다. 이 코드를 helper로 만들어서 반복을 줄이자.

```go
//server_test.go
func assertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}
```

`server_test.go` 와 `CLI_test.go`에 있는 코드를 대체하면 테스트는 아래와 같아야 한다.

```go
//CLI_test.go
func TestCLI(t *testing.T) {
	in := strings.NewReader("Chris wins\n")
	playerStore := &StubPlayerStore{}

	cli := &CLI{playerStore, in}
	cli.PlayPoker()

	assertPlayerWin(t, playerStore, "Chris")
}
```

이제 실제로 입력값을 읽게 하는 _다른_ 테스트를 다른 사용자의 입력값과 함께 작성해보자.

## 테스트부터 작성하기

```go
//CLI_test.go
func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &StubPlayerStore{}

		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &StubPlayerStore{}

		cli := &CLI{playerStore, in}
		cli.PlayPoker()

		assertPlayerWin(t, playerStore, "Cleo")
	})

}
```

## 테스트 실행해보기

```
=== RUN   TestCLI
--- FAIL: TestCLI (0.00s)
=== RUN   TestCLI/record_chris_win_from_user_input
    --- PASS: TestCLI/record_chris_win_from_user_input (0.00s)
=== RUN   TestCLI/record_cleo_win_from_user_input
    --- FAIL: TestCLI/record_cleo_win_from_user_input (0.00s)
        CLI_test.go:27: did not store correct winner got 'Chris' want 'Cleo'
FAIL
```

## 테스트를 통과하는 최소한의 코드 작성하기

우리는 `io.Reader`로부터 입력값을 읽기 위해서 [`bufio.Scanner`](https://golang.org/pkg/bufio/) 사용할 것이다.

> Package bufio는 버퍼가 있는 I/O를 구현한다. 이것은 인터페이스를 구현하는 다른 객체(Reader 혹은 Writer)를 만들면서 io.Reader나 io.Writer 객체를 감싸고 버퍼와 텍스트 I/O를 위한 몇몇 도움을 제공한다.

코드를 아래와 같이 업데이트 하자.

```go
//CLI.go
type CLI struct {
	playerStore PlayerStore
	in          io.Reader
}

func (cli *CLI) PlayPoker() {
	reader := bufio.NewScanner(cli.in)
	reader.Scan()
	cli.playerStore.RecordWin(extractWinner(reader.Text()))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}
```

이제 테스트들은 통과할 것이다.

-   `Scanner.Scan()` 새로운 줄이 나올 때까지 읽을 것이다.
-   그 때 우리는 Scanner가 읽은 `string`을 리턴하기 위해 `Scanner.Text()`를 사용할 수 있다.

이제 통과하는 테스트들이 있으니깐 이를 `main`에 작성해야 한다. 우리는 항상 가능한 빨리 완전히 작동가능한 소프트웨어를 만드는 것을 갈망해야하는 것을 기억하자.

`main.go`파일 안에 아래와 같이 입력하고 실행시키자. (의존성을 해결하기 위해 당신 컴퓨터에 맞춰서 주소를 변경해야 할 수도 있다.)

```go
package main

import (
	"fmt"
	"github.com/quii/learn-go-with-tests/command-line/v3"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")

	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	game := poker.CLI{store, os.Stdin}
	game.PlayPoker()
}
```

아래와 같은 에러메시지가 뜰 것이다.

```
command-line/v3/cmd/cli/main.go:32:25: implicit assignment of unexported field 'playerStore' in poker.CLI literal
command-line/v3/cmd/cli/main.go:32:34: implicit assignment of unexported field 'in' in poker.CLI literal
```

우리가 `playerStore` 필드들과 `CLI`안에 `in`을 값을 대입하려고 했기 때문에 생겨난 일이다. 그것들은 내보내지지 않은 (unexported), private 필드들이다. 테스트 코드들은 `CLI` (`poker`)와 같은 패키지에 있었기 때문에 이렇게 _할 수 있었다_. 하지만 `main`은 `main` 패키지에 있기 때문에 접근 권한이 없다.

이는 너가 쓴 코드들을 _통합_ 하는 것이 얼마나 중요한지를 알려준다. 우리는 `CLI`의 의존 변수들을 private으로 올바르게 만들었다 (왜냐하면 `CLI`를 사용하는 사용자들에게 이 변수들을 보여주고 싶지 않기 때문이다). 하지만 사용자들이 이것을 만들 방법을 아직 제공하지 않았다.

이러한 문제를 더 일찍 알 방법이 없었을까?

### `package mypackage_test`

지금까지 있었던 다른 예시들에서는 우리는 테스트 파일을 만들 때 우리가 테스트를 하고 있는 같은 패키지에서 테스트를 선언했다.

이것은 충분히 괜찮고 우리가 패키지 내부의 무언가를 테스트하고 싶어하는 이상한 경우에 내보내지 않는 (unexported) 타입에 접근할 수 있다는 것을 의미합니다.

그러나 일반적으로 내부의 무언가를 테스트하지 않겠다고 해왔는데, Go는 이를 강제로 하는 데 도움이 될 수 있을까요? 만약 `main`처럼 우리가 접근이 가능한 내보내진 (exported) 타입들에만 테스트가 가능하다면 어떨까요?

만약 당신이 여러 개의 패키지가 있는 프로젝트를 작성한다면, 테스트 패키지 이름 뒤에 `_test`를 붙일 것을 강력하게 추천한다. 이렇게 한다면, 당신은 당신의 패키지에 public 타입들에만 접근이 가능해 질 것이다. 이는 public API들만 테스트한다는 규칙을 강제하는데 큰 도움이 된다. 당신이 아직도 패키지 내부에 있는 것을 테스트하고 싶다면, 테스트하고 싶은 패키지에 따로 테스트를 만들 수 있다.

테스트 기반 개발(TDD)의 단점은 코드를 테스트할 수 없다면 코드를 사용하는 사람들이 그것을 가지고 와서 쓰기 어려울 수 있다는 것이다. `package foo_test`를 사용함으로써 마치 당신이 당신의 패키지를 사용하는 사람들처럼 불러와서 당신의 코드를 테스트하게끔 강요하게 해 도움을 줄 것이다.

`main`을 고치기 전에 `CLI_test.go`의 테스트들의 패키지를 `poker_test`로 변경하자.

만약 너가 괜찮은 IDE를 쓰고 있다면 코드의 빨간 줄들이 갑자기 많이 보일 것이다. 당신이 컴파일러를 실행시키면 아래와 같은 에러들을 발견할 것이다.

```
./CLI_test.go:12:19: undefined: StubPlayerStore
./CLI_test.go:17:3: undefined: assertPlayerWin
./CLI_test.go:22:19: undefined: StubPlayerStore
./CLI_test.go:27:3: undefined: assertPlayerWin
```

우리는 이제 패키지 디자인에 대해 더 많은 질문이 생기게 되었다. 우리 소프트웨어를 테스트 하기 위해서 우리는 `CLI_test`에서 더이상 사용할 수 없는 내보내지지 않은 stub (테스트를 위해 작성한 임시 코드)과 helper 함수들을 가지게 되었다. 왜냐하면 이 helper들은 `poker`패키지안에 `_test.go`파일에 정의되어 있기 때문이다.

#### 우리는 과연 stub과 helper들을 'public'으로 만들어야 할까?

이것은 각자의 주관에 따라 다르다. 어떤 사람은 테스트를 용이하게 하기 위해 패키지의 API를 오염시키고 싶지 않다고 주장할 수 있다.

Mitchell Hashimoto의 ["Go 언어 고급 테스팅"](https://speakerdeck.com/mitchellh/advanced-testing-with-go?slide=53) 프리젠테이션을 보면, Hashicorp에서 어떤 방식으로 이렇게 하는지를 설명한다. 그래서 패키지의 사용자들은 stub들을 다시 만들지 않고도 테스트 코드들을 작성할 수 있다. 우리의 코드를 예로 든다면, `poker` 패키지를 사용하는 개발자들이 그들의 코드에서 작동하기를 희망하는 `PlayerStore` stub을 만들지 않게 하는 것을 의미한다.

일화적으로 나는 다른 공유 피키지에서 이 기법을 사용했고 그것이 다른 개발자들이 내가 만든 패키지들을 사용할 때 시간을 절약할 수 있다는 점에서 매우 유용하다는 것이 증명되었다.

그러니깐 `testing.go`라는 파일을 만들고 그 안에 stub과 helper들을 만들어보자.

```go
//testing.go
package poker

import "testing"

type StubPlayerStore struct {
	scores   map[string]int
	winCalls []string
	league   []Player
}

func (s *StubPlayerStore) GetPlayerScore(name string) int {
	score := s.scores[name]
	return score
}

func (s *StubPlayerStore) RecordWin(name string) {
	s.winCalls = append(s.winCalls, name)
}

func (s *StubPlayerStore) GetLeague() League {
	return s.league
}

func AssertPlayerWin(t testing.TB, store *StubPlayerStore, winner string) {
	t.Helper()

	if len(store.winCalls) != 1 {
		t.Fatalf("got %d calls to RecordWin want %d", len(store.winCalls), 1)
	}

	if store.winCalls[0] != winner {
		t.Errorf("did not store correct winner got %q want %q", store.winCalls[0], winner)
	}
}

// 해야 할 일: 다른 helper 함수들을 직접 만들어보자.
```

이 패키지를 불러오는 사람들한테 보이게 하려면 helper 함수들을 public으로 만들어야 한다 (함수의 첫 글자를 대문자로 하면 내보낼 수 있다는 것을 기억하자).

마치 다른 패키지에서 코드를 사용하는 것처럼 `CLI` 테스트에서 그 코드를 불러와야 한다.

```go
//CLI_test.go
func TestCLI(t *testing.T) {

	t.Run("record chris win from user input", func(t *testing.T) {
		in := strings.NewReader("Chris wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := &poker.CLI{playerStore, in}
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Chris")
	})

	t.Run("record cleo win from user input", func(t *testing.T) {
		in := strings.NewReader("Cleo wins\n")
		playerStore := &poker.StubPlayerStore{}

		cli := &poker.CLI{playerStore, in}
		cli.PlayPoker()

		poker.AssertPlayerWin(t, playerStore, "Cleo")
	})

}
```

우리가 `main`에서 가지고 있던 문제들이 똑같이 있다는 걸 볼 수 있다.

```
./CLI_test.go:15:26: implicit assignment of unexported field 'playerStore' in poker.CLI literal
./CLI_test.go:15:39: implicit assignment of unexported field 'in' in poker.CLI literal
./CLI_test.go:25:26: implicit assignment of unexported field 'playerStore' in poker.CLI literal
./CLI_test.go:25:39: implicit assignment of unexported field 'in' in poker.CLI literal
```

이 문제를 가장 쉽게 우회하는 방법은 다른 타입들 처럼 생성자를 만드는 것이다. `CLI`도 바꿔야 하는데 그렇게 하면 reader대신에 `bufio.Scanner`를 생성할 때 자동으로 감싸지면서 가질 수 있게 된다.

```go
//CLI.go
type CLI struct {
	playerStore PlayerStore
	in          *bufio.Scanner
}

func NewCLI(store PlayerStore, in io.Reader) *CLI {
	return &CLI{
		playerStore: store,
		in:          bufio.NewScanner(in),
	}
}
```

이렇게 함으로써 사용자의 입력값을 받는 우리의 코드를 단순화하고 리팩터링 할 수 있다.

```go
//CLI.go
func (cli *CLI) PlayPoker() {
	userInput := cli.readLine()
	cli.playerStore.RecordWin(extractWinner(userInput))
}

func extractWinner(userInput string) string {
	return strings.Replace(userInput, " wins", "", 1)
}

func (cli *CLI) readLine() string {
	cli.in.Scan()
	return cli.in.Text()
}
```

생성자를 쓰도록 테스트를 변경하고 통과하는 테스트로 돌아와야 한다.

마지막으로 우리의 새로운 `main.go` 파일로 가서 우리가 방금 만든 생성자를 사용해야 한다.

```go
//cmd/cli/main.go
game := poker.NewCLI(store, os.Stdin)
```

이제 실행시켜보자. "Bob wins"라고 입력해라.

### 리팩터링 하기

파일을 열고 사용자의 입력값에서 `file_system_store`를 만드는 각각의 애플리케이션에 반복되는 코드들이 좀 있다. 이것은 우리의 패키지 디자인의 약간의 문제가 있다고 느껴지게 하기 때문에 주소를 불러들여 파일을 열고 `PlayerStore`를 리턴하는 식으로 중복된 코드들을 하나의 함수로 만들어야 한다.

```go
//file_system_store.go
func FileSystemPlayerStoreFromFile(path string) (*FileSystemPlayerStore, func(), error) {
	db, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return nil, nil, fmt.Errorf("problem opening %s %v", path, err)
	}

	closeFunc := func() {
		db.Close()
	}

	store, err := NewFileSystemPlayerStore(db)

	if err != nil {
		return nil, nil, fmt.Errorf("problem creating file system player store, %v ", err)
	}

	return store, closeFunc, nil
}
```

우리의 모든 애플리케이션들을 store를 만들어내는 이 함수를 사용하도록 리팩터링하자.

#### CLI 애플리케이션 코드

```go
//cmd/cli/main.go
package main

import (
	"fmt"
	"github.com/quii/learn-go-with-tests/command-line/v3"
	"log"
	"os"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	fmt.Println("Let's play poker")
	fmt.Println("Type {Name} wins to record a win")
	poker.NewCLI(store, os.Stdin).PlayPoker()
}
```

#### 웹 서버 애플리케이션 코드

```go
//cmd/webserver/main.go
package main

import (
	"github.com/quii/learn-go-with-tests/command-line/v3"
	"log"
	"net/http"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSystemPlayerStoreFromFile(dbFileName)

	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := poker.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

다른 사용자 인터페이스를 가졌음에도 불구하고 코드의 구성은 거의 비슷한 이 대칭성을 느껴야 한다. 우리의 디자인이 지금까지 괜찮다는 걸 느끼게 한다.
또한 `FileSystemPlayerStoreFromFile`이 파일을 닫는 함수를 리턴한다는 것을 알아야 한다. 그렇기 때문에 우리는 store를 다 쓰고 나서 열었던 파일을 닫을 수 있다.

## 마무리

### 패키지 구조

우리는 이 챕터에서 지금까지 작성해왔던 도메인 코드들을 재사용해서 두 개의 애플리케이션을 만들려고 했다. 이를 위해서 패키지 구조를 다시 변경해서 각각의 `main`을 위한 별도의 폴더를 가지게 되었다.

이를 통해 우리는 내보내지 않은 (unexported) 변수들로 인해 코드를 통합하는 문제에 부딪혔으며, 이는 작은 단위로 일을 하고 자주 코드를 통합해야하는 것이 얼마나 중요한지를 더욱 입증한다.

우리는 `mypackage_test`가 어떻게 당신의 코드와 함께 사용할 다른 패키지들과 같은 경험을 제공하는 테스트 환경을 만드는데 도움을 주는지 배웠다. 그리고 이것이 코드가 작동하는지(혹은 작동하지 않는지)와 코드를 통합할 때 생기는 지를 빠르고 쉽게 찾을 수 있는지도 배웠다.

### 사용자의 입력값 읽기

우리는 `io.Reader`를 이용해서 `os.Stdin`에서 입력값을 읽는 것이 얼마나 쉬운지 보여줬다. 그리고 사용자의 입력값을 각각의 줄로 나눠서 쉽게 읽기 위해서 `bufio.Scanner`를 사용했다.

### 간단한 추상화는 코드 재사용을 간단하게 한다.

`PlayerStore`를 새로운 애플리케이션에서 사용하는데 거의 노력이 들지 않았다 (패키지를 한 번 바꿨을 뿐이다). stub 코드를 public으로 하기로 결정했기 때문에 결국 테스트 또한 매우 쉬웠다.
