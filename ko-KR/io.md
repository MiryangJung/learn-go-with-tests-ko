# 입출력과 정렬

**[이 챕터의 모든 코드는 여기에서 확인할 수 있다.](https://github.com/quii/learn-go-with-tests/tree/main/io)**

[이전 챕터에서](json.md) 우리는 어플리케이션에 새로운 엔드포인트 `league`를 추가하는 과정을 계속 반복해왔다. 그 과정에서 우리는 JSON을 다루는 법, 타입을 임베딩하는 법 그리고 라우팅하는 법을 배웠다.

프로덕트 오너는 서버가 재시작 될 때 소프트웨어가 점수들을 잃을까 약간 불안해한다. 왜냐하면 우리는 스토어를 인메모리로 구현했기 때문이다. 게다가 그녀는 우리가 `league` 엔드포인트가 이긴 횟수를 기준으로 정렬한 선수들을 반환해야 하는 것을 해석하지 못하는 것이 만족스럽지 않다!

## 지금까지의 코드

```go
// server.go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
)

// PlayerStore는 선수들에 대한 점수 정보를 저장한다
type PlayerStore interface {
    GetPlayerScore(name string) int
    RecordWin(name string)
    GetLeague() []Player
}

// Player는 이긴 횟수와 함께 이름을 저장한다
type Player struct {
    Name string
    Wins int
}

// PlayerServer는 사용자 정보를 위한 HTTP 인터페이스이다
type PlayerServer struct {
    store PlayerStore
    http.Handler
}

const jsonContentType = "application/json"

// NewPlayerServer는 라우팅이 정의된 PlayerServer를 생성한다
func NewPlayerServer(store PlayerStore) *PlayerServer {
    p := new(PlayerServer)

    p.store = store

    router := http.NewServeMux()
    router.Handle("/league", http.HandlerFunc(p.leagueHandler))
    router.Handle("/players/", http.HandlerFunc(p.playersHandler))

    p.Handler = router

    return p
}

func (p *PlayerServer) leagueHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("content-type", jsonContentType)
    json.NewEncoder(w).Encode(p.store.GetLeague())
}

func (p *PlayerServer) playersHandler(w http.ResponseWriter, r *http.Request) {
    player := r.URL.Path[len("/players/"):]

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

```go
// in_memory_player_store.go
package main

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
    return &InMemoryPlayerStore{map[string]int{}}
}

type InMemoryPlayerStore struct {
    store map[string]int
}

func (i *InMemoryPlayerStore) GetLeague() []Player {
    var league []Player
    for name, wins := range i.store {
        league = append(league, Player{name, wins})
    }
    return league
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
    i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
    return i.store[name]
}
```

```go
// main.go
package main

import (
    "log"
    "net/http"
)

func main() {
    server := NewPlayerServer(NewInMemoryPlayerStore())

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

코드에 해당하는 테스트들은 챕터 상단의 링크에서 확인할 수 있다.

## 데이터를 저장한다.

사용할 수 있는 데이터베이스는 많지만 우리는 매우 간단한 접근을 할 것 이다. 우리는 이 어플리케이션의 데이터를 JSON 파일의 형태로 저장할 것이다.

이 접근 방식은 데이터를 매우 이동이 쉬운 형태로 유지하고 상대적으로 구현하기 쉽게 만든다.

이는 확장에 특별히 좋은 형태는 아니지만 프로토타입으로 지금으로서는 충분하다. 우리는 `PlayerStore`를 추상화했기 때문에 만약 우리의 환경이 변하고 더 이상 적절하지 않다면 간단히 다른 무언가로 간단히 변경할 수 있다.

지금 당장 우리는 `InMemoryPlayerStore`를 유지할 것이기 때문에 통합 테스트들은 우리가 새로운 스토어를 개발하는 동안에도 계속 통과할 것이다. 우리는 새로운 구현이 통합 테스트를 충분히 통과할 것이라는 확신을 가지게 됬을 때 이를 교체하고 `InMemoryPlayerStore`를 삭제할 것이다.

## 테스트부터 작성하기

이제 당신은 데이터를 읽고(`io.Reader`), 쓰는(`io.Writer`) 표준 라이브러리들의 인터페이스와 실제 파일들을 사용하지 않고 이런 기능들을 테스트하기 위해 표준 라이브러리를 사용하는 법에 익숙해져야 한다.

이 작업을 완료하기 위해 우리는 `PlayStore`를 구현해야 한다. 그리고 스토어가 우리가 구현해야 하는 메서드를 호출할 수 있도록 하는 테스트를 작성해야 한다. `GetLeague`부터 시작해보자.

```go
//file_system_store_test.go
func TestFileSystemStore(t *testing.T) {

    t.Run("league from a reader", func(t *testing.T) {
        database := strings.NewReader(`[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)

        store := FileSystemPlayerStore{database}

        got := store.GetLeague()

        want := []Player{
            {"Cleo", 10},
            {"Chris", 33},
        }

        assertLeague(t, got, want)
    })
}
```

우리는 `FileSystemPlayerStore`가 데이터를 읽을 수 있도록 하는 `Reader`를 반환하는 `strings.NewReader`를 사용중이다. `main`에 파일을 추가할 것이고, 이 파일 또한 `Reader`이다.

## 테스트 실행해보기

```
# github.com/quii/learn-go-with-tests/io/v1
./file_system_store_test.go:15:12: undefined: FileSystemPlayerStore
```

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

새로운 파일에 `FileSystemPlayerStore`를 정의한다.

```go
//file_system_store.go
type FileSystemPlayerStore struct {}
```

다시 시도한다.

```
# github.com/quii/learn-go-with-tests/io/v1
./file_system_store_test.go:15:28: too many values in struct initializer
./file_system_store_test.go:17:15: store.GetLeague undefined (type FileSystemPlayerStore has no field or method GetLeague)
```

우리가 `Reader`를 넘겨줬지만 입력을 기대하지 않고, `GetLeague`가 아직 정의되어 있지 않기 때문에 컴파일 에러가 발생한다.

```go
//file_system_store.go
type FileSystemPlayerStore struct {
    database io.Reader
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
    return nil
}
```

한번 더 시도한다...

```
=== RUN   TestFileSystemStore//league_from_a_reader
    --- FAIL: TestFileSystemStore//league_from_a_reader (0.00s)
        file_system_store_test.go:24: got [] want [{Cleo 10} {Chris 33}]
```

## 테스트를 통과하는 최소한의 코드 작성하기

우리는 전에 리더로부터 JSON을 읽어왔다.

```go
//file_system_store.go
func (f *FileSystemPlayerStore) GetLeague() []Player {
    var league []Player
    json.NewDecoder(f.database).Decode(&league)
    return league
}
```

테스트는 통과할 것이다.

## 리팩터링 하기

우리는 이전에 이것을 _했었다_! 서버를 위한 우리의 테스트 코드는 응답으로부터 JSON을 디코딩 해야한다.

함수에 DRY(Do not Repeat Yourself)를 적용해보자.

`league.go`라는 새로운 파일을 생성해 안에 넣는다.

```go
//league.go
func NewLeague(rdr io.Reader) ([]Player, error) {
    var league []Player
    err := json.NewDecoder(rdr).Decode(&league)
    if err != nil {
        err = fmt.Errorf("problem parsing league, %v", err)
    }

    return league, err
}
```

구현과 `server_test.go` 안에 있는 `getLeagueFromResponse` 테스트 헬퍼에서 이를 호출한다.

```go
//file_system_store.go
func (f *FileSystemPlayerStore) GetLeague() []Player {
    league, _ := NewLeague(f.database)
    return league
}
```

우리는 아직 파싱 에러를 처리할 방법을 가지고 있지는 않지만 계속 진행해보자.

### 문제 찾기

우리의 구현에는 한가지 흠이 있다. 무엇보다도 우리 스스로 `io.Reader`가 어떻게 정의되어 있는지 다시 생각해보자.

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

파일에서 보이듯, 당신은 끝까지 바이트 단위로 읽어나가는 것을 생각해낼 수 있을 것이다. 만약 `Read`를 두 번 시도한다면 어떻게 될까?

아래 내용을 현재 테스트의 끝에 추가하자.

```go
//file_system_store_test.go

// 다시 읽는다.
got = store.GetLeague()
assertLeague(t, got, want)
```

테스트가 통과하기를 원하지만, 만약 당신이 테스트를 실행했다면 통과하지 못했을 것이다.

문제는 `Reader`가 끝에 다다랐을 때 더 이상 읽을 것이 없다는 것이다. 우리는 처음으로 돌아가라고 말해줄 방법이 필요하다.

[ReadSeeker](https://golang.org/pkg/io/#ReadSeeker)는 이 문제를 해결하도록 도와줄 수 있는 표준 라이브러리에 있는 다른 인터페이스이다.

```go
type ReadSeeker interface {
    Reader
    Seeker
}
```

임베딩을 기억하는가? 이것은 `Reader`와 [`Seeker`](https://golang.org/pkg/io/#Seeker)로 구성된 인터페이스이다.

```go
type Seeker interface {
    Seek(offset int64, whence int) (int64, error)
}
```

좋아 보인다. `FileSystemPlayerStore`을 이 인터페이스로 바꿀 수 있을까?

```go
//file_system_store.go
type FileSystemPlayerStore struct {
    database io.ReadSeeker
}

func (f *FileSystemPlayerStore) GetLeague() []Player {
    f.database.Seek(0, 0)
    league, _ := NewLeague(f.database)
    return league
}
```

테스트를 실행해보자. 테스트가 통과했다! 운이 좋게도 우리가 테스트에 사용한 `string.NewReader`도 `ReadSeeker`를 구현하고 있어서 더 이상 변경할 필요가 없다.

다음으로 우리는 `GetPlayerScore`를 구현할 것이다.

## 테스트부터 작성하기

```go
//file_system_store_test.go
t.Run("get player score", func(t *testing.T) {
    database := strings.NewReader(`[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)

    store := FileSystemPlayerStore{database}

    got := store.GetPlayerScore("Chris")

    want := 33

    if got != want {
        t.Errorf("got %d want %d", got, want)
    }
})
```

## 테스트 실행해보기

```
./file_system_store_test.go:38:15: store.GetPlayerScore undefined (type FileSystemPlayerStore has no field or method GetPlayerScore)
```

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

우리는 테스트가 컴파일 될 수 있도록 새로운 타입에 메서드를 추가해야 한다.

```go
//file_system_store.go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {
    return 0
}
```

이제 컴파일은 성공하고 테스트는 실패한다.

```
=== RUN   TestFileSystemStore/get_player_score
    --- FAIL: TestFileSystemStore//get_player_score (0.00s)
        file_system_store_test.go:43: got 0 want 33
```

## 테스트를 통과하는 최소한의 코드 작성하기

우리는 league를 순회하며 선수를 찾고 그들의 점수를 반환할 수 있다.

```go
//file_system_store.go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

    var wins int

    for _, player := range f.GetLeague() {
        if player.Name == name {
            wins = player.Wins
            break
        }
    }

    return wins
}
```

## 리팩터링 하기

당신은 테스트를 보조하기 위한 많은 리팩토링 방법들을 봤을 것이기 때문에 당신이 해낼 수 있도록 남겨둘 것이다.

```go
//file_system_store_test.go
t.Run("get player score", func(t *testing.T) {
    database := strings.NewReader(`[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)

    store := FileSystemPlayerStore{database}

    got := store.GetPlayerScore("Chris")
    want := 33
    assertScoreEquals(t, got, want)
})
```

마지막으로 우리는 `RecorddWin`으로 점수들을 기록해야 한다.

## 테스트부터 작성하기

우리의 접근 방법은 쓰기에 상당히 근시안적이다. 우리는 파일 안에 있는 JSON의 한 "줄"만을 간단히 업데이트할 수 없다. 때문에 우리는 모든 쓰기마다 우리 데이터 베이스 전체의 새로운 표현을 저장해야만 합니다.

어떻게 쓰기를 할 수 있을까? 우리는 보통 `Writer`를 사용했지만, 우리는 이미 우리만의 `ReadSeeker`를 가지고 있다. 잠재적으로 우리는 2개의 의존성을 가질 수 있지만, 표준 라이브러리는 이미 파일에 필요한 모든 작업을 수행할 수 있는 `ReadWriteSeeker` 인터페이스를 가지고 있다.

우리의 타입을 바꿔보자.

```go
//file_system_store.go
type FileSystemPlayerStore struct {
    database io.ReadWriteSeeker
}
```

컴파일이 되는지 확인하자.

```go
./file_system_store_test.go:15:34: cannot use database (type *strings.Reader) as type io.ReadWriteSeeker in field value:
    *strings.Reader does not implement io.ReadWriteSeeker (missing Write method)
./file_system_store_test.go:36:34: cannot use database (type *strings.Reader) as type io.ReadWriteSeeker in field value:
    *strings.Reader does not implement io.ReadWriteSeeker (missing Write method)
```

`strings.Reader`가 `ReadWriteSeeker`를 구현하지 못한다는 것은 그리 놀라운 일이 아니다. 그렇다면 무엇을 해야할까?

우리는 두가지를 선택할 수 있다.

- 각각의 테스트를 위한 임시 파일을 생성한다. `*os.File`은 `ReadWriteSeeker`를 구현한다. Create a temporary file for each test. `*os.File` implements `ReadWriteSeeker`. 이 방법의 장점은 더 통합 테스트에 가까워진다는 것이다. 우리는 실제 파일 시스템에서 읽고 쓰고 있기 때문에 더 높은 수준의 신뢰성을 얻을 수 있습니다. 단점은 단위 테스트가 더 빠르고 일반적으로 간단하기 때문에 더 선호된다는 것이다. 또한 임시 파일들을 만들어내고 테스트 이후에 파일들이 지워졌는지 확인하기 위한 일들을 더 해야만한다.
- 우리는 써드파티 라이브러리를 사용한다. [Mattetti](https://github.com/mattetti)는 우리에게 필요한 인터페이스가 구현되어 있으면서 파일 시스템을 건드리지 않는 [filebuffer](https://github.com/mattetti/filebuffer)를 작성했다.

이 중 특별히 틀린 답이 있다 생각하지는 않지만, 써드파티 라이브러리를 사용하는 것을 선택한다면 의존성 관리에 대해서 설명을 해야만 한다! 그러니 우리는 파일을 사용할 것이다.

테스트를 추가하기 전에 `os.File`을 `strings.Reader`로 바꿔서 테스트가 컴파일 될 수 있도록 해야한다.

데이터가 포함된 임시 파일을 생성하는 헬퍼 함수를 만들어보자.

```go
//file_system_store_test.go
func createTempFile(t testing.TB, initialData string) (io.ReadWriteSeeker, func()) {
    t.Helper()

    tmpfile, err := ioutil.TempFile("", "db")

    if err != nil {
        t.Fatalf("could not create temp file %v", err)
    }

    tmpfile.Write([]byte(initialData))

    removeFile := func() {
    	tmpfile.Close()
        os.Remove(tmpfile.Name())
    }

    return tmpfile, removeFile
}
```

[TempFile](https://golang.org/pkg/io/ioutil/#TempDir)은 우리가 사용할 수 있는 임시 파일을 생성한다. 우리가 넘긴 `"db"` 값은 만들어질 임의 파일 이름에 붙는 접두사입니다. 이렇게 함으로써 다른 파일들과 우연히 충돌하는 것을 방지합니다.

당신은 `ReadWriteSeeker`(파일) 뿐만 아니라 함수 또한 반환하고 있다는 것을 알고 있어야 합니다. 테스트가 끝나면 파일이 삭제되어야 한다는 것을 확실히 해야합니다. 에러가 발생하기 쉽고 리더에 무관심하기 때문에 테스트에 파일들의 세부사항들을 유출하고 싶지 않다. `removeFile` 함수를 반환함으로써, 헬퍼 안의 세부사항들을 관리할 수 있고 모든 호출자는 `defer cleanDatabase()`만 실행하기만 하면 된다.

```go
//file_system_store_test.go
func TestFileSystemStore(t *testing.T) {

    t.Run("league from a reader", func(t *testing.T) {
        database, cleanDatabase := createTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
        defer cleanDatabase()

        store := FileSystemPlayerStore{database}

        got := store.GetLeague()

        want := []Player{
            {"Cleo", 10},
            {"Chris", 33},
        }

        assertLeague(t, got, want)

        // read again
        got = store.GetLeague()
        assertLeague(t, got, want)
    })

    t.Run("get player score", func(t *testing.T) {
        database, cleanDatabase := createTempFile(t, `[
            {"Name": "Cleo", "Wins": 10},
            {"Name": "Chris", "Wins": 33}]`)
        defer cleanDatabase()

        store := FileSystemPlayerStore{database}

        got := store.GetPlayerScore("Chris")
        want := 33
        assertScoreEquals(t, got, want)
    })
}
```

테스트를 실행하면 통과될 것이다! 상당히 많은 변경사항들이 있지만 드디어 인터페이스의 정의가 끝났다는 느낌이 들고, 이제부터 새로운 테스트를 추가하기가 매우 쉬워졌다.

기존 선수들의 승리를 기록하는 첫번째 반복을 시작해보자.

```go
//file_system_store_test.go
t.Run("store wins for existing players", func(t *testing.T) {
    database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
    defer cleanDatabase()

    store := FileSystemPlayerStore{database}

    store.RecordWin("Chris")

    got := store.GetPlayerScore("Chris")
    want := 34
    assertScoreEquals(t, got, want)
})
```

## 테스트 실행해보기

`./file_system_store_test.go:67:8: store.RecordWin undefined (type FileSystemPlayerStore has no field or method RecordWin)`

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

새로운 메서드를 추가한다.

```go
//file_system_store.go
func (f *FileSystemPlayerStore) RecordWin(name string) {

}
```

```
=== RUN   TestFileSystemStore/store_wins_for_existing_players
    --- FAIL: TestFileSystemStore/store_wins_for_existing_players (0.00s)
        file_system_store_test.go:71: got 33 want 34
```

구현이 비어 있어서 오래된 점수가 반환된다.

## 테스트를 통과하는 최소한의 코드 작성하기

```go
//file_system_store.go
func (f *FileSystemPlayerStore) RecordWin(name string) {
    league := f.GetLeague()

    for i, player := range league {
        if player.Name == name {
            league[i].Wins++
        }
    }

    f.database.Seek(0,0)
    json.NewEncoder(f.database).Encode(league)
}
```

당신은 내가 왜 `player.Wins++`가 아닌 `league[i].Wins++` 실행하는지 스스로 되묻고 있을 수도 있다.

당신이 슬라이스에 `범위`를 지정하면 루프의 현재 인덱스(우리의 경우 `i`)와 이 인덱스의 요소의 _복사본_ 을 반환받는다. 복사본 `Wins` 값의 변경은 우리가 반복중인 `league` 슬라이스에 아무런 영향을 주지 않는다. 때문에, `league[i]`를 이용해 실제 값에 대한 참조를 얻어오고 그 값을 변경해야 한다.

테스트를 실행하면 통과될 것이다!

## 리팩터링 하기

`GetPlayerScore`와 `RecordWin`에서 플레이어를 이름으로 찾기 위해 `[]Player`를 반복시킨다.

우리는 이 공통 코드를 `FileSystemStore` 내부에서 리팩터링할 수 있겠지만, 나에게는 이 코드가 새로운 타입으로 만들 수 있는 유용한 코드가 될 것이라는 느낌이 든다. 지금까지의 "League" 작업은 `[]Player`와 함께 였지만, 새로운 타입인 `League`를 생성할 수 있다. 이렇게 하면 다른 개발자들이 이해하기도 쉬울 것이고 우리가 사용할 수 있도록 새로운 메서드를 붙일 수도 있을 것이다.

`league.go` 안에 다음 코드를 추가한다.

```go
//league.go
type League []Player

func (l League) Find(name string) *Player {
    for i, p := range l {
        if p.Name==name {
            return &l[i]
        }
    }
    return nil
}
```

이젠 `League`를 가진 누구나 쉽게 선수를 주어진 선수를 찾을 수 있습니다.

`PlayerStore` 인터페이스가 `[]Player`가 아닌 `League`를 반환하도록 변경하자. 다시 테스트를 실행하면 인터페이스를 변경했기 때문에 컴파일 에러를 얻을 것이다. 하지만 매우 고치기 쉽다; 그냥 반환 타입을 `[]Player`에서 `League`로 변경해라.

이렇게 되면 `file_system_store` 안에 있는 메서드를 간단하게 할 수 있다.

```go
//file_system_store.go
func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

    player := f.GetLeague().Find(name)

    if player != nil {
        return player.Wins
    }

    return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
    league := f.GetLeague()
    player := league.Find(name)

    if player != nil {
        player.Wins++
    }

    f.database.Seek(0, 0)
    json.NewEncoder(f.database).Encode(league)
}
```

이렇게 바꾸니 훨씬 더 좋아 보이고, `League`와 관련된 다른 유용한 기능들을 리팩터링 하는 방법에 대해서 알 수 있었다.

이제 우리는 새로운 선수들의 승리를 기록하는 시나리오를 처리해야 한다.

## 테스트부터 작성하기

```go
//file_system_store.go
t.Run("store wins for new players", func(t *testing.T) {
    database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
    defer cleanDatabase()

    store := FileSystemPlayerStore{database}

    store.RecordWin("Pepper")

    got := store.GetPlayerScore("Pepper")
    want := 1
    assertScoreEquals(t, got, want)
})
```

## 테스트 실행해보기

```
=== RUN   TestFileSystemStore/store_wins_for_new_players#01
    --- FAIL: TestFileSystemStore/store_wins_for_new_players#01 (0.00s)
        file_system_store_test.go:86: got 0 want 1
```

## 테스트를 통과하는 최소한의 코드 작성하기

`Find`가 선수를 찾을 수 없을 때 `nil`을 반환하는 시나리오를 처리하면 된다.

```go
//file_system_store.go
func (f *FileSystemPlayerStore) RecordWin(name string) {
    league := f.GetLeague()
    player := league.Find(name)

    if player != nil {
        player.Wins++
    } else {
        league = append(league, Player{name, 1})
    }

    f.database.Seek(0, 0)
    json.NewEncoder(f.database).Encode(league)
}
```

행복한 경로(예외나 에러 조건이 없는 경우)는 괜찮아보이기 때문에 이제 우리의 새 `Store`를 통합 테스트에서 사용해볼 수 있다. 이를 통해 소프트웨어가 잘 동작한다는 확신을 얻을 수 있고, 중복인 `InMemoryPlayerStore`를 삭제할 수 있다.

`TestRecordingWinsAndRetrievingThem` 안의 오래된 스토어를 바꾼다.

```go
//server_integration_test.go
database, cleanDatabase := createTempFile(t, "")
defer cleanDatabase()
store := &FileSystemPlayerStore{database}
```

테스트를 실행하면 통과할 것이고, 이제 `InMemoryPlayerStore`. `main.go`는 컴파일 문제가 생겼을 것이다. 컴파일 문제가 생겼다는 것은 이제 "실제" 코드에서 새로운 스토어를 사용해야 한다는 것을 말해준다.

```go
//main.go
package main

import (
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

    store := &FileSystemPlayerStore{db}
    server := NewPlayerServer(store)

    if err := http.ListenAndServe(":5000", server); err != nil {
        log.Fatalf("could not listen on port 5000 %v", err)
    }
}
```

- 데이터베이스를 위한 파일을 생성한다.
- `os.OpenFile`의 두번째 인자값으로 파일을 열 수 있는 권한을 정의할 수 있다. `O_RDWR`는 읽고 쓸 수 있다는 것을 의미하고, _그리고_ `os.O_CREATE`는 존재하지 않는 파일을 생성할 수 있다는 것을 의미한다.
- 세번째 인자값은 파일에 대한 권한을 정의하는 것을 의미하고, 우리의 경우 모든 사용자가 파일에 읽고 쓸 수 있다는 것을 의미한다. [(superuser.com에서 더 자세한 설명을 확인할 수 있다.)](https://superuser.com/questions/295591/what-is-the-meaning-of-chmod-666).

이제 프로그램을 실행하는 것은 재시작을 하더라도 파일 내의 데이터를 유지한다. 만세!

## 추가로 리팩터링 하기. 그리고 성능 고려하기

매번 누군가 `GetLeague()` 또는 `GetPlayerScore()`를 호출하면 전체 파일을 읽어와 JSON으로 파싱한다. 하지만 `FileSystemStore`가 league 전체의 상태를 책임을 가지기 때문에 이렇게 할 필요가 없다; 프로그램이 시작할 때와 데이터가 바뀌어서 파일을 업데이트해야할 때만 파일을 읽어야만 한다.

We can create a constructor which can do some of this initialisation for us and store the league as a value in our `FileSystemStore` to be used on the reads instead.

```go
//file_system_store.go
type FileSystemPlayerStore struct {
    database io.ReadWriteSeeker
    league League
}

func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
    database.Seek(0, 0)
    league, _ := NewLeague(database)
    return &FileSystemPlayerStore{
        database:database,
        league:league,
    }
}
```

이 방법으로 디스크를 한 번만 읽으면 된다. 이제 디스크에서 리그를 가져오는 이전의 호출들을 모두 교체할 수 있게 되었고 그냥 `f.league`를 대신에 사용하면 된다.

```go
//file_system_store.go
func (f *FileSystemPlayerStore) GetLeague() League {
    return f.league
}

func (f *FileSystemPlayerStore) GetPlayerScore(name string) int {

    player := f.league.Find(name)

    if player != nil {
        return player.Wins
    }

    return 0
}

func (f *FileSystemPlayerStore) RecordWin(name string) {
    player := f.league.Find(name)

    if player != nil {
        player.Wins++
    } else {
        f.league = append(f.league, Player{name, 1})
    }

    f.database.Seek(0, 0)
    json.NewEncoder(f.database).Encode(f.league)
}
```

만약 테스트를 실행한다면 이제 `FileSystemPlayerStore`를 초기화하는 것에 대한 불평할 것이다. 그러므로 새로운 생성자를 호출하는 것으로 바꾸기만 하면 된다.

### 또 다른 문제

파일을 다루는 방법에 몇가지 더 단순한 것들이 있는데, 이는 매우 끔찍한 버그를 만들어 _낼 수_ 있다.

`RecordWin` 할 때, 파일의 처음으로 `Seek`하기 위해 돌아가고 새로운 데이터를 쓴다-하지만 새로운 데이터가 이 전에 있었던 것들보다 더 작다면 어떻게 될까?

현재의 경우, 이건 불가능하다. 점수를 수정하거나 삭제하지 않기 때문에 데이터가 커질수 밖에 없다. 그러나 코드를 이렇게 두는 것은 책임감이 없어 보인다; 삭제 시나리오가 생겨나지 않을 것이라는 것을 생각하지 않을 수 없다.

그렇다면 어떻게 테스트 해야할까? 첫 번째로 우리는 코드를 리팩터링해서 _우리가 작성한 코드로부터 데이터 쓰기_ 에 대한 관심사를 분리해야한다. 이렇게 하면 우리가 원하는대로 동작하는지 확인하기 위해 나눠서 테스트할 수 있다.

"쓰기는 처음부터" 기능을 캡슐화하기 위해 새로운 타입을 만들 것이다. 나는 이것을 `Tape`이라고 부를 것이다. 다음과 같은 새로운 파일을 생성한다:

```go
//tape.go
package main

import "io"

type tape struct {
    file io.ReadWriteSeeker
}

func (t *tape) Write(p []byte) (n int, err error) {
    t.file.Seek(0, 0)
    return t.file.Write(p)
}
```

지금은 `Seek`를 캡슐화했으니, `Write`를 구현하고 있다는 것에 주의해야한다. `FileSystemStore`가 `Writer`를 대신에 참조로 가질 수 있다는 것을 의미한다.

```go
//file_system_store.go
type FileSystemPlayerStore struct {
    database io.Writer
    league   League
}
```

`Tape`를 사용하도록 생성자를 업데이트한다.

```go
//file_system_store.go
func NewFileSystemPlayerStore(database io.ReadWriteSeeker) *FileSystemPlayerStore {
    database.Seek(0, 0)
    league, _ := NewLeague(database)

    return &FileSystemPlayerStore{
        database: &tape{database},
        league:   league,
    }
}
```

마침내 우리는 `RecordWin` 호출로부터 `Seek`를 제거함으로써 우리가 원했던  놀라운 성과를 얻을 수 있다. 맞다, 그렇게 크게 느껴지지는 않는다, 하지만 이것은 쵯환 우리가 만약 다른 종류의 쓰기를 할 경우 우리가 원하는대로 동작하는 우리만의 `Write`에 의존할 수 있다는 것이다. 추가로 이제부터 잠재적으로 문제가 있는 코드를 각각 테스트할 수 있고 수정할 수 있게 되었다.

파일의 전체 내용을 원래 내용보다 작도록 업데이트하는 테스트를 작성해보자.

## 테스트부터 작성하기

우리의 테스트는 몇가지 내용을 가진 파일을 생성할 것이고, `tape`을 이용해 이 파일에 쓰고, 파일에 무엇이 있는지 전체를 다시 읽어온다. `tape_test.go` 내부를 살펴보자: 

```go
//tape_test.go
func TestTape_Write(t *testing.T) {
    file, clean := createTempFile(t, "12345")
    defer clean()

    tape := &tape{file}

    tape.Write([]byte("abc"))

    file.Seek(0, 0)
    newFileContents, _ := ioutil.ReadAll(file)

    got := string(newFileContents)
    want := "abc"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

## 테스트 실행해보기

```
=== RUN   TestTape_Write
--- FAIL: TestTape_Write (0.00s)
    tape_test.go:23: got 'abc45' want 'abc'
```

생각처럼 되었다! 우리가 원하는대로 데이터를 쓰지만, 나머지 기존 데이터를 남겨둔다.

## 테스트를 통과하는 최소한의 코드 작성하기

`os.File`은 truncate 함수를 가지고 있고 이를 활용하면 효과적으로 파일을 비울 수 있다. 우리는 원하는 것을 얻기 위해서는 이 함수를 호출할 수 있어야 한다.

다음과 같이 `tape`을 변경한다:

```go
//tape.go
type tape struct {
    file *os.File
}

func (t *tape) Write(p []byte) (n int, err error) {
    t.file.Truncate(0)
    t.file.Seek(0, 0)
    return t.file.Write(p)
}
```

컴파일러는 `io.ReadWriteSeeker`가 예상되지만 `*os.File`을 보내고 있는 많은 곳들에서 실패할 것이다. 지금까지는 이런 문제들을 스스로 수정할 수 있었지만, 만약 막힌다면 소스코드를 확인해라.

리팩토링을 했다면 `TestTape_Write` 테스트는 통과될 것이다!

### 다른 작은 리팩터링 하기

`RecordWin`에는 `json.NewEncoder(f.database).Encode(f.league)` 라는 코드 라인이 있다.

하지만 매 쓰기마다 새로운 인코더를 생성할 필요가 없이, 생성자 안에 초기화해서 대신 사용할 수 있다.

`Encoder`의 참조를 타입에 저장하고 생성자 내에서 초기화한다.

```go
//file_system_store.go
type FileSystemPlayerStore struct {
    database *json.Encoder
    league   League
}

func NewFileSystemPlayerStore(file *os.File) *FileSystemPlayerStore {
    file.Seek(0, 0)
    league, _ := NewLeague(file)

    return &FileSystemPlayerStore{
        database: json.NewEncoder(&tape{file}),
        league:   league,
    }
}
```

`RecordWin`에서 이를 사용한다.

```go
func (f *FileSystemPlayerStore) RecordWin(name string) {
	player := f.league.Find(name)

	if player != nil {
		player.Wins++
	} else {
		f.league = append(f.league, Player{name, 1})
	}

	f.database.Encode(f.league)
}
```

## 몇 가지 규칙을 어기진 않았는가? private한 것들을 테스트 했는가? 인터페이스가 없는가?

### private 타입 테스트

_일반적으로_ private한 것들을 테스트 하지 않는 것을 선호한다. 왜냐하면 이것은 테스트와 구현을 너무 강하게 연결할 수 있기 때문이다. 그리고 이는 나중에 리팩터링하는 것을 방해할 수 있다.

그러나 테스트는 우리에게 _확신_ 을 준다는 것을 잊으면 안된다.

우리는 구현에 어떠한 기능의 변경 또는 삭제가 일어났을때 구현이 동작한다는 확신을 가질 수 없다. 특히 초기 접근 방법에 대한 문제점을 인지하지 못하는 여러명이 함께 작업하는 경우, 코드를 이대로 내버려둘 수 없다.

드디어 마지막 테스트이다! 작동 방식을 변경하기로 결정했다면 테스트를 삭제하는 것은 더 이상 재앙이 아닐 것이다. 그러나 나중에 유지보수하는 사람을 위한 최소한의 요구사항을 가지게 된다.

### 인터페이스

우리는 우리의 새로운 `PlayStore`를 유닛 테스트하기 위한 가장 쉬운 방법인 `io.Reader`를 사용해서 코드를 만들기 시작했다. 코드를 개발해나가며 `io.ReaderWriter`에서 `io.ReadWriteSeeker`로 옮겨갔고, `*os.File`와는 별개로 실제로 구현된 것이 표준 라이브러리 안에는 없다는 것을 알게 되었다. 우리는 직접 만들거나 오픈소스를 사용하는 것으로 결정을 내렸을 수도 있었지만 테스트를 위한 임시 파일을 만드는 것이 실용적이라고 느꼈다.

마지막으로 `*os.File`에도 있는 `Truncate`가 필요하다. 이런 요구사항들을 만족시키는 자체 인터페이스를 만들기 위한 옵션이었을 것이다.

```go
type ReadWriteSeekTruncate interface {
    io.ReadWriteSeeker
    Truncate(size int64) error
}
```

그러나 우리는 이를 통해 무엇을 얻었는가? 우리는 _mock을 만들지 않았다_ 는 것과 **file system** 스토어가 `*os.File`외의 다른 타입을 가지는 것이 비현실적이기 때문에 인터페이스의 다형성이 필요하지 않다는 것을 명심해야 한다.

우리가 했던 것처럼 타입을 자르고 변경하고 실험하는 것을 두려워하지 말아야 한다. 정적 타입 언어를 사용하는 가장 큰 이점은 컴파일러가 모든 변경사항들로 당신을 도와줄 수 있다는 점이다.

## 에러 핸들링

정렬을 작업하기 전에 현재의 코드에 만족하며 가지고 있는 기술 부채를 모두 제거했다는 것을 확실히 해야 한다. 동작하는 소프트웨어에 가능한 빨리 도달해야 한다(red state에서 벗어나야 한다)는 것은 중요한 원칙이지만 에러 케이스들을 무시해야한다는 것을 의미하지는 않는다.

`FileSystemStore.go`로 돌아가보면 생성자 안에는 `league, _ := NewLeague(f.database)`가 있다.

`NewLeague`는 제공받는 `io.Reader`로부터 league를 파싱하지 못하는 경우 에러를 반환할 수 있다.

그 때는 이미 실패 테스트들을 했기 때문에 에러를 무시하는 것이 실용적이었다. 만약 이것을 한번에 다뤘다면, 한번에 두 가지를 효율적으로 해결했을 것이다.

생성자가 에러를 반환할 수 있도록 만들어보자.

```go
//file_system_store.go
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {
    file.Seek(0, 0)
    league, err := NewLeague(file)

    if err != nil {
        return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
    }

    return &FileSystemPlayerStore{
        database: json.NewEncoder(&tape{file}),
        league:   league,
    }, nil
}
```

도움이 되는 에러 메시지를 제공하는 것은 매우 중요하다는 것을 명심해야 한다(테스트 처럼). 인터넷 상의 사람들이 농담삼아 말하는 대부분의 Go 코드는 다음과 같다:

```go
if err != nil {
    return err
}
```

**저런 코드는 100% 자연스럽지 않다.** 문맥적 정보(i.e 에러를 만들기 위해 당신이 하는 것)를 에러 메시지에 더하는 것은 소프트웨어를 더 쉽게 운영할 수 있도록 도와준다.

컴파일하면 에러가 나올 것이다.

```
./main.go:18:35: multiple-value NewFileSystemPlayerStore() in single-value context
./file_system_store_test.go:35:36: multiple-value NewFileSystemPlayerStore() in single-value context
./file_system_store_test.go:57:36: multiple-value NewFileSystemPlayerStore() in single-value context
./file_system_store_test.go:70:36: multiple-value NewFileSystemPlayerStore() in single-value context
./file_system_store_test.go:85:36: multiple-value NewFileSystemPlayerStore() in single-value context
./server_integration_test.go:12:35: multiple-value NewFileSystemPlayerStore() in single-value context
```

main에서 우리는 에러를 출력하면서 프로그램을 종료하는 것을 원한다.

```go
//main.go
store, err := NewFileSystemPlayerStore(db)

if err != nil {
    log.Fatalf("problem creating file system player store, %v ", err)
}
```

테스트에서는 에러가 없다는 것을 assert 해야한다. 이를 도와주는 헬퍼를 만들 수 있다.

```go
//file_system_store_test.go
func assertNoError(t testing.TB, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("didn't expect an error but got one, %v", err)
    }
}
```

다른 컴파일 문제들을 이 헬퍼를 이용해 통과시킬 수 있다. 드디어 실패 테스트를 할 수 있다:

```
=== RUN   TestRecordingWinsAndRetrievingThem
--- FAIL: TestRecordingWinsAndRetrievingThem (0.00s)
    server_integration_test.go:14: didn't expect an error but got one, problem loading player store from file /var/folders/nj/r_ccbj5d7flds0sf63yy4vb80000gn/T/db841037437, problem parsing league, EOF
```

파일이 비어있기 때문에 league를 파싱할 수 없다. 이전에는 에러를 모두 무시했기 때문에 에러가 없었다.

이제 유효한 JSON을 넣어 큰 통합 테스트를 수정해보자:

```go
//server_integration_test.go
func TestRecordingWinsAndRetrievingThem(t *testing.T) {
    database, cleanDatabase := createTempFile(t, `[]`)
    //etc...
```

이제 모든 테스트가 통과했기 때문에, 파일이 비어있을 때의 시나리오를 해결해야 한다.

## 테스트부터 작성하기

```go
//file_system_store_test.go
t.Run("works with an empty file", func(t *testing.T) {
    database, cleanDatabase := createTempFile(t, "")
    defer cleanDatabase()

    _, err := NewFileSystemPlayerStore(database)

    assertNoError(t, err)
})
```

## 테스트 실행해보기

```
=== RUN   TestFileSystemStore/works_with_an_empty_file
    --- FAIL: TestFileSystemStore/works_with_an_empty_file (0.00s)
        file_system_store_test.go:108: didn't expect an error but got one, problem loading player store from file /var/folders/nj/r_ccbj5d7flds0sf63yy4vb80000gn/T/db019548018, problem parsing league, EOF
```

## 테스트를 통과하는 최소한의 코드 작성하기

다음과 같이 생성자를 변경한다.

```go
//file_system_store.go
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

    file.Seek(0, 0)

    info, err := file.Stat()

    if err != nil {
        return nil, fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
    }

    if info.Size() == 0 {
        file.Write([]byte("[]"))
        file.Seek(0, 0)
    }

    league, err := NewLeague(file)

    if err != nil {
        return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
    }

    return &FileSystemPlayerStore{
        database: json.NewEncoder(&tape{file}),
        league:   league,
    }, nil
}
```

`file.Stat`은 파일의 상태를 반환하고, 우리는 이를 이용해 파일의 크기를 확인할 수 있다. 만약 파일이 비어있다면, 비어있는 배열을 `Write`하고 시작점으로 `Seek`해서 남아있는 코드를 준비한다.

## 리팩터링 하기

생성자가 이제 조금 지저분해졌기 때문에, 초기화 코드를 함수로 추출해보자:

```go
//file_system_store.go
func initialisePlayerDBFile(file *os.File) error {
    file.Seek(0, 0)

    info, err := file.Stat()

    if err != nil {
        return fmt.Errorf("problem getting file info from file %s, %v", file.Name(), err)
    }

    if info.Size()==0 {
        file.Write([]byte("[]"))
        file.Seek(0, 0)
    }

    return nil
}
```

```go
//file_system_store.go
func NewFileSystemPlayerStore(file *os.File) (*FileSystemPlayerStore, error) {

    err := initialisePlayerDBFile(file)

    if err != nil {
        return nil, fmt.Errorf("problem initialising player db file, %v", err)
    }

    league, err := NewLeague(file)

    if err != nil {
        return nil, fmt.Errorf("problem loading player store from file %s, %v", file.Name(), err)
    }

    return &FileSystemPlayerStore{
        database: json.NewEncoder(&tape{file}),
        league:   league,
    }, nil
}
```

## 정렬하기

프로덕트 오너는 `/league`가 선수들이 그들의 점수에 따라 높은점수에서 낮은 점수로 정렬되어서 반환되기를 원한다.

여기서 중요한 것은 소프트웨어 어디에서 이 동작이 일어나야 할지를 결정하는 것이다. 만약 우리가 "실제" 데이터베이스를 사용한다면 우리는 `ORDER BY`와 같은 것을 사용할 것이고 매우 빠르게 동작할 것이다. 이러한 이유로 `PlayerStore`의 구현들이 책임을 가져야할 것으로 느껴진다.

## 테스트부터 작성하기

`TestFileSystemStore`에 있는 첫 번째 테스트의 assertion을 변경할 수 있다.

```go
//file_system_store_test.go
t.Run("league sorted", func(t *testing.T) {
    database, cleanDatabase := createTempFile(t, `[
        {"Name": "Cleo", "Wins": 10},
        {"Name": "Chris", "Wins": 33}]`)
    defer cleanDatabase()

    store, err := NewFileSystemPlayerStore(database)

    assertNoError(t, err)

    got := store.GetLeague()

    want := []Player{
        {"Chris", 33},
        {"Cleo", 10},
    }

    assertLeague(t, got, want)

    // read again
    got = store.GetLeague()
    assertLeague(t, got, want)
})
```

JSON이 들어오는 순서는 잘못된 순서이고 우리의 `요구사항`은 호출한 사람에게 올바른 순서로 반환되는지를 확인합니다.

## 테스트 실행해보기

```
=== RUN   TestFileSystemStore/league_from_a_reader,_sorted
    --- FAIL: TestFileSystemStore/league_from_a_reader,_sorted (0.00s)
        file_system_store_test.go:46: got [{Cleo 10} {Chris 33}] want [{Chris 33} {Cleo 10}]
        file_system_store_test.go:51: got [{Cleo 10} {Chris 33}] want [{Chris 33} {Cleo 10}]
```

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func (f *FileSystemPlayerStore) GetLeague() League {
    sort.Slice(f.league, func(i, j int) bool {
        return f.league[i].Wins > f.league[j].Wins
    })
    return f.league
}
```

[`sort.Slice`](https://golang.org/pkg/sort/#Slice)

> Slice는 제공된 슬라이스를 제공된 less function으로 정렬한다.

간단하다!

## 정리

### 우리가 다룬 것

- `Seeker` 인터페이스와 `Reader`와 `Writer`의 관계
- 파일을 이용한 작업
- 지저분한 모든 것들을 숨기고 있는 파일들을 테스트하기 위해 사용성 좋은 헬퍼를 작성
- 슬라이스 정렬을 위한 `sort.Slice`
- 어플리케이션의 구조적 변화를 안전하게 하기 위해 컴파일러를 활용

### 위반한 규칙들

- 소프트웨어 엔지니어링의 대부분의 규칙들은 실제 규칙들이 아니라 단지 80%의 시간에만 해당되는 모범 사례이다.
- 내부 함수를 테스트 하지 않는다는 이전의 "규칙들"이 우리에게 그리 도움이 되지 않았던 시나리오를 발견했기 때문에 규칙을 위반했다.
- 규칙을 위반할 때에는 당신이 만들어낸 trade-off에 대해 이해하는 것이 중요하다. 우리의 경우, 단지 한개의 테스트였을 뿐이고 그렇지 않았다면 시나리오를 실행하기 매우 어려웠을 것이기 때문에 괜찮았다.
- 규칙을 위반할 수 있기 위해서는 **규칙을 먼저 이해하야만 합니다**. 기타를 배우는 것으로 비유를 들어보자. 당신이 스스로에 대해 얼마나 창의적이라고 생각하는지는 중요하지 않다. 먼저 기초를 이해하고 연습하는 것이 더 중요하다.

### 우리의 소프트웨어가 도달한 위치

- 선수를 생성하고 그들의 점수를 증가시킬수 있는 HTTP API가 있다.
- 리그에 있는 모든 선수들의 점수를 JSON으로 반환할 수 있다.
- 데이터는 JSON 파일로 유지된다.
