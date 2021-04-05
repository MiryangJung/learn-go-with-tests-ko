# WebSockets

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/main/websockets)**

이 챕터에서 우리는 우리의 어플리케이션을 개선하기 위해 어떻게 웹소켓을 사용할 지 배울 것입니다.

## Project recap

우리는 우리의 poker codebase에 두개의 어플리케이션을 가집니다.

-*Command line app*. 사용자에게 게임의 플레이어 수를 입력하라는 메시지를 표시합니다. 그때부터 플레이어에게 시간이 지남에 따라 증가하는 "블라인드 베팅"값이 무엇인지 알려줍니다. 사용자는 언제든지` "{Playername} wins"`를 입력하여 게임을 완료하고 저장에 승자를 기록 할 수 있습니다.
-*Web app*. 사용자가 게임의 승자를 기록하고 리그 테이블을 표시 할 수 있습니다. 명령 줄 앱과 동일한 저장소를 공유합니다.

## Next steps

제품 소유자는 command line app에 만족하지만 해당 기능을 브라우저에 가져오는 것을 더 선호합니다. 그녀는 사용자가 플레이어의 수를 적을 수 있는 텍스트 상자가 있는 웹페이지를 상상합니다. 텍스트 상자 양식을 제출할 때, 페이지는 블라인드 값을 보여주고 적절한 때에 자동으로 블라인드 값을 업데이트 합니다. Command line app과 마찬가지로 사용자는 승자를 선언 할 수 있으며 데이터베이스에 저장됩니다.

겉으로는 아주 간단하게 들리지만 항상 그렇듯이 소프트웨어 작성에 대해 _반복적_ 접근 방식을 강조해야합니다.

우선 우리는 HTML을 제공해야합니다. 지금까지 모든 HTTP 엔드 포인트는 일반 텍스트 또는 JSON을 반환했습니다. 우리는 우리가 알고있는 것과 동일한 기술을 사용할 수 있지만 (궁극적으로는 strings이기 때문에) 더 깨끗한 솔루션을 위해 [html/template] (https://golang.org/pkg/html/template/) 패키지를 사용할 수도 있습니다.

또한 브라우저를 새로 고치지 않고도 'The blind is now *y*'라는 메시지를 사용자에게 비동기식으로 보낼 수 있어야합니다. 이를 위해 [WebSockets] (https://en.wikipedia.org/wiki/WebSocket)를 사용할 수 있습니다.

> WebSocket은 단일 TCP 연결을 통해 전이중 통신 채널을 제공하는 컴퓨터 통신 프로토콜입니다.

우리가 여러 기술을 취하고 있다는 점을 감안할 때 가능한 한 적은 양의 유용한 작업을 먼저 수행 한 다음 반복하는 것이 훨씬 더 중요합니다.

따라서 사용자가 우승자를 기록 할 수있는 양식이있는 웹 페이지를 가장 먼저 만들것입니다. Plain form을 사용하는 대신 WebSocket을 사용하여 해당 데이터를 서버로 전송하여 기록합니다.

그 후에 우리는 약간의 인프라 코드가 셋업될 시점까지 블라인드 경고로 처리 할 것입니다.

### 자바스크립트에 대한 테스트는 어떡합니까?

이를 수행하기 위해 작성된 JavaScript가 있지만 테스트 작성은 하지 않을 것입니다.

물론 가능하지만 간결함을 위해 이에 대한 설명은 포함하지 않겠습니다.

죄송합니다. O'Reilly를 로비하여 나에게 돈을 주어 "테스트로 JavaScript 배우기"를 만들게 하십시오.

## 테스트를 가장 먼저 만들어라.

가장 먼저해야 할 일은 사용자가 '/game'을 눌렀을 때 일부 HTML을 제공하는 것입니다

다음으로 웹 서버 관련 코드를 알려드립니다.

```go
type PlayerServer struct {
	store PlayerStore
	http.Handler
}

const jsonContentType = "application/json"

func NewPlayerServer(store PlayerStore) *PlayerServer {
	p := new(PlayerServer)

	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))

	p.Handler = router

	return p
}
```

지금 우리가 할 수 있는 가장 쉬운 일은 우리가 `GET /game` 할 때 `200`을 얻었는지 확인하는 것입니다.

```go
func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := NewPlayerServer(&StubPlayerStore{})

		request, _ := http.NewRequest(http.MethodGet, "/game", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response.Code, http.StatusOK)
	})
}
```

## 테스트를 실행해보세요.
```
--- FAIL: TestGame (0.00s)
=== RUN   TestGame/GET_/game_returns_200
    --- FAIL: TestGame/GET_/game_returns_200 (0.00s)
    	server_test.go:109: did not get correct status, got 404, want 200
```

## 테스트가 통과하도록 충분한 코드를 작성하세요.

우리 서버에는 router setup이 있으므로 비교적 쉽게 고칠 수 있습니다.

Router에 추가하기 위해,

```go
router.Handle("/game", http.HandlerFunc(p.game))
```

그런 다음 `game` 메서드를 작성합니다.

```go
func (p *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
```

## Refactor

The server code is already fine due to us slotting in more code into the existing well-factored code very easily.

We can tidy up the test a little by adding a test helper function `newGameRequest` to make the request to `/game`. Try writing this yourself.

```go
func TestGame(t *testing.T) {
	t.Run("GET /game returns 200", func(t *testing.T) {
		server := NewPlayerServer(&StubPlayerStore{})

		request :=  newGameRequest()
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assertStatus(t, response, http.StatusOK)
	})
}
```

그리고 당신은 내가 `assertStatus`의 `response.Code`를 `response`로 바꿨다는 것을 인지할 것입니다. 나는 그것이 더 읽기 쉬운 것 같기에 바꿨습니다.

이제 엔드 포인트가 HTML을 반환하도록 해야합니다. 아래와 같습니다.

```html

<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Let's play poker</title>
</head>
<body>
<section id="game">
    <div id="declare-winner">
        <label for="winner">Winner</label>
        <input type="text" id="winner"/>
        <button id="winner-button">Declare winner</button>
    </div>
</section>
</body>
<script type="application/javascript">

    const submitWinnerButton = document.getElementById('winner-button')
    const winnerInput = document.getElementById('winner')

    if (window['WebSocket']) {
        const conn = new WebSocket('ws://' + document.location.host + '/ws')

        submitWinnerButton.onclick = event => {
            conn.send(winnerInput.value)
        }
    }
</script>
</html>
```

우리는 매우 간단한 웹페이지를 가지고 있습니다.

사용자가 우승자를 입력하기위한 텍스트 입력
-우승자를 선언하기 위해 클릭 할 수있는 버튼.
-버튼이 눌러졌을 때, 서버와의 WebSocket 연결 수립을 위한 JavaScript 코드

`WebSocket`은 대부분의 최신 브라우저에 내장되어 있으므로 라이브러리를 가져 오는 것에 대해 걱정할 필요가 없습니다. 웹 페이지는 이전 브라우저에서는 작동하지 않지만 이 시나리오에서는 괜찮습니다.

### 올바른 마크 업을 반환하는지 어떻게 테스트합니까?

몇 가지 방법이 있습니다. 책 전체에서 강조해왔듯이 당신이 작성하는 테스트가 비용을 정당화할 만큼의 충분한 가치를 갖는 것이 중요합니다.

1. Selenium과 같은 것을 사용하여 브라우저 기반 테스트를 작성하십시오. 이러한 테스트는 어떤 종류의 실제 웹 브라우저를 시작하고 그와 상호 작용하는 사용자를 시뮬레이션하기 때문에 모든 접근 방식 중 가장 "현실적"입니다. 이러한 이 테스트는 시스템 작동에 대한 확신을 줄 수 있지만 단위 테스트보다 작성하기가 더 어렵고 실행 속도가 훨씬 느립니다. 우리 제품의 목적 상 이것은 과합니다.
2. 정확한 문자열 일치를 수행합니다. 이것은 괜찮을 _수_ 있지만 이러한 종류의 테스트는 결국 매우 취약합니다. 누군가가 마크 업을 변경하는 순간 _실제로 고장난_ 것이 없는데도 테스트가 실패하게됩니다.
3. 올바른 템플릿을 호출하는지 확인합니다. 우리는 HTML을 제공하기 위해 표준 lib의 템플릿 라이브러리를 사용할 것이며 (곧 논의 될 것입니다) _thing_에 삽입하여 HTML을 생성하고 우리가 제대로하고 있는지 확인하기 위해 호출을 감시 할 수 있습니다. 이것은 우리 코드의 디자인에 영향을 미칠 것이지만 실제로 많은 것을 테스트하지는 않습니다. 올바른 템플릿 파일로 부르는 것 이외에는 그렇지 않습니다. 프로젝트에 템플릿이 하나만 있으면 여기서 실패 할 가능성은 낮아 보입니다.
So in the book "Learn Go with Tests" for the first time, we're not going to write a test.

파일에 `game.html`이라는 마크 업을 넣습니다.

다음으로 방금 작성한 endpoint 다음과 같이 변경하십시오.

```go
func (p *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("game.html")

	if err != nil {
		http.Error(w, fmt.Sprintf("problem loading template %s", err.Error()), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
```

[`html/template`] (https://golang.org/pkg/html/template/)은 HTML을 만들기위한 Go 패키지입니다. 우리의 경우`template.ParseFiles`를 호출하여 html 파일의 경로를 제공합니다. 오류가 없다고 가정하면 템플릿을 '실행'하여 'io.Writer'에 기록 할 수 있습니다. 우리의 경우 인터넷에 '쓰기'를 원하므로 'http.ResponseWriter'를 제공합니다.

테스트를 작성하지 않았으므로 웹 서버를 수동으로 테스트하여 원하는대로 작동하는지 확인하는 것이 좋습니다. `cmd/webserver`로 이동하여`main.go` 파일을 실행합니다. `http://localhost:5000/game`을 방문하십시오.

템플릿을 찾을 수 없다는 오류를 대면해야 합니다. 경로를 당신의 폴더에 상대적인 방식으로 변경하거나 `cmd/webserver` 디렉토리에`game.html`의 복사본을 가질 수 있습니다. 나는 프로젝트 루트 내부에 파일에 대한 심볼릭 링크 (`ln -s ../../game.html game.html`)를 생성하기로 선택했습니다. 변경하면 서버를 실행할 때 반영됩니다.

이렇게 변경하고 다시 실행하면 UI가 표시됩니다.

이제 우리는 서버에 대한 WebSocket 연결을 통해 문자열을 얻을 때 게임의 승자로 선언하는지 테스트해야합니다.
## 테스트를 먼저 작성하세요.

처음으로 우리는 WebSocket으로 작업 할 수 있도록 외부 라이브러리를 사용할 것입니다.

`go get github.com/gorilla/websocket`를 실행하세요.

이렇게하면 우수한 [Gorilla WebSocket] (https://github.com/gorilla/websocket) 라이브러리의 코드를 가져옵니다. 이제 새로운 요구 사항에 대한 테스트를 업데이트 할 수 있습니다.
```go
t.Run("when we get a message over a websocket it is a winner of a game", func(t *testing.T) {
    store := &StubPlayerStore{}
    winner := "Ruth"
    server := httptest.NewServer(NewPlayerServer(store))
    defer server.Close()

    wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"

    ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
    if err != nil {
        t.Fatalf("could not open a ws connection on %s %v", wsURL, err)
    }
    defer ws.Close()

    if err := ws.WriteMessage(websocket.TextMessage, []byte(winner)); err != nil {
        t.Fatalf("could not send message over ws connection %v", err)
    }

    AssertPlayerWin(t, store, winner)
})
```

`websocket` 라이브러리를 제대로 import했는지 확인하세요. 내 IDE가 자동으로 import 했으므로 당신의 것도 할 것입니다.

브라우저에서 어떤 일이 발생하는지 테스트하려면 자체 WebSocket 연결을 열고 여기에 작성해야합니다.

이전 테스트는 서버에서 메소드를 호출했지만 이제 서버에 지속적인 연결을 해야합니다. 이를 위해 우리는`http.Handler`를 가져 와서 연결을 수신하는`httptest.NewServer`를 사용합니다.

'websocket.DefaultDialer.Dial'을 사용하여 서버에 메시지를 보내어서 'winner'와 메시지를 보내려고합니다.

마지막으로 플레이어 store에서 승자가 기록되었는지 확인하기 위해 assert할 것입니다.

## 테스트를 실행하려고 시도합니다.
```
=== RUN   TestGame/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game
    --- FAIL: TestGame/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game (0.00s)
        server_test.go:124: could not open a ws connection on ws://127.0.0.1:55838/ws websocket: bad handshake
```

'/ws'에서 WebSocket 연결을 허용하도록 서버를 변경하지 않았으므로 아직 handshacking을 하지 않습니다.

## 통과 할 수 있도록 충분한 코드 작성합니다.

라우터에 다른 목록 추가

```go
router.Handle("/ws", http.HandlerFunc(p.webSocket))
```

그런 다음 새로운 `webSocket` 핸들러를 추가합니다.

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.Upgrade(w, r, nil)
}
```

WebSocket 연결을 수락하기 위해 요청을 'Upgrade'합니다. 이제 테스트를 다시 실행하면 다음 오류로 이동해야 할 것입니다.

```
=== RUN   TestGame/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game
    --- FAIL: TestGame/when_we_get_a_message_over_a_websocket_it_is_a_winner_of_a_game (0.00s)
        server_test.go:132: got 0 calls to RecordWin want 1
```

이제 연결이 열렸으므로 메시지를 듣고 승자로 기록하는 것을 원합니다.

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	conn, _ := upgrader.Upgrade(w, r, nil)
	_, winnerMsg, _ := conn.ReadMessage()
	p.store.RecordWin(string(winnerMsg))
}
```

(예, 우리는 지금 많은 오류를 무시하고 있습니다!)

`conn.ReadMessage ()`는 연결에서 메시지를 기다릴 때 block합니다. 메시지를 받게 되면, 그것을 `RecordWin`에 사용합니다. 이것은 마침내 WebSocket 연결을 닫습니다.

테스트를 시도하고 실행하면 여전히 실패합니다.

문제는 타이밍입니다. WebSocket 연결이 메시지를 읽고 승리를 기록하는 사이에 지연이 있으며 우리의 테스트는 그것이 일어나기 전에 완료됩니다. 최종 assertion 앞에 'time.Sleep'을 짧게 입력하여 이를 테스트 할 수 있습니다.

지금은 그렇게 갑시다. 그러나 테스트에 임의의 sleep을 취하는 것은 ** 매우 나쁜 습관 **이라는 것을 알아야합니다.

```go
time.Sleep(10 * time.Millisecond)
AssertPlayerWin(t, store, winner)
```

## Refactor

우리는 이 테스트가 서버 코드와 테스트 코드 모두에서 작동하도록하기 위해 많은 죄를 지었지만 이것이 우리가 작업하는 가장 쉬운 방법임을 기억하십시오.

우리는 테스트로 뒷받침되는 끔찍한 _작동하는_ 소프트웨어를 가지고 있습니다. 그래서 이제 우리는 그것을 멋지게 만들 수 있고 우연히 어떤 것도 깨뜨리지 않을 것이라는 것을 압니다.

서버 코드부터 시작하겠습니다.

모든 WebSocket 연결 요청에 대해 다시 선언 할 필요가 없기 때문에 `upgrader`를 패키지 내부의 private 값으로 이동할 수 있습니다.

```go
var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := wsUpgrader.Upgrade(w, r, nil)
	_, winnerMsg, _ := conn.ReadMessage()
	p.store.RecordWin(string(winnerMsg))
}
```

`template.ParseFiles ( "game.html")`에 대한 호출은 모든`GET /game`에서 실행됩니다. 즉, 템플릿을 다시 파싱 할 필요가 없더라도 모든 요청에 대해 파일 시스템으로 이동합니다.

`PlayerServer`의 관련 변경 사항은 다음과 같습니다.

```go
type PlayerServer struct {
	store PlayerStore
	http.Handler
	template *template.Template
}

const htmlTemplatePath = "game.html"

func NewPlayerServer(store PlayerStore) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	p.template = tmpl
	p.store = store

	router := http.NewServeMux()
	router.Handle("/league", http.HandlerFunc(p.leagueHandler))
	router.Handle("/players/", http.HandlerFunc(p.playersHandler))
	router.Handle("/game", http.HandlerFunc(p.game))
	router.Handle("/ws", http.HandlerFunc(p.webSocket))

	p.Handler = router

	return p, nil
}

func (p *PlayerServer) game(w http.ResponseWriter, r *http.Request) {
	p.template.Execute(w, nil)
}
```

`NewPlayerServer`의 signature 변경함으로써 이제 컴파일 문제가 발생합니다. 직접 시도하고 수정하거나 어려움이 있다면 소스 코드를 참조하십시오.

테스트 코드를 위해`mustMakePlayerServer (t * testing.T, store PlayerStore) * PlayerServer`라는 helper를 만들어 테스트에서 오류 노이즈를 숨길 수있었습니다.

```go
func mustMakePlayerServer(t *testing.T, store PlayerStore) *PlayerServer {
	server, err := NewPlayerServer(store)
	if err != nil {
		t.Fatal("problem creating player server", err)
	}
	return server
}
```

마찬가지로 WebSocket 연결을 만들 때 불쾌한 오류 노이즈를 숨길 수 있도록 또 다른 helper `mustDialWS`를 만들었습니다.

```go
func mustDialWS(t *testing.T, url string) *websocket.Conn {
	ws, _, err := websocket.DefaultDialer.Dial(url, nil)

	if err != nil {
		t.Fatalf("could not open a ws connection on %s %v", url, err)
	}

	return ws
}
```

마지막으로 테스트 코드에서 메시지 전송을 정리하는 helper를 만들 수 있습니다.

```go
func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
	t.Helper()
	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		t.Fatalf("could not send message over ws connection %v", err)
	}
}
```

이제 테스트가 통과되었습니다. 서버를 실행하고`/game`에서 승자를 선언하십시오. `/league`에 기록 된 것을 볼 수 있습니다. 당첨자를 얻을 때마다 _연결을 닫음_을 기억하세요. 연결을 다시 열려면 페이지를 새로 고침해야합니다.

사용자가 게임의 승자를 기록 할 수 있도록 간단한 웹 양식을 만들었습니다. 사용자가 많은 플레이어와 게임을 시작할 수 있도록 반복하고 서버는 시간이 지남에 따라 블라인드 값이 무엇인지 알려주는 메시지를 클라이언트에 푸시합니다.

우선`game.html`을 업데이트하여 새로운 요구 사항에 맞게 클라이언트 측 코드를 업데이트하십시오.

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Lets play poker</title>
</head>
<body>
<section id="game">
    <div id="game-start">
        <label for="player-count">Number of players</label>
        <input type="number" id="player-count"/>
        <button id="start-game">Start</button>
    </div>

    <div id="declare-winner">
        <label for="winner">Winner</label>
        <input type="text" id="winner"/>
        <button id="winner-button">Declare winner</button>
    </div>

    <div id="blind-value"/>
</section>

<section id="game-end">
    <h1>Another great game of poker everyone!</h1>
    <p><a href="/league">Go check the league table</a></p>
</section>

</body>
<script type="application/javascript">
    const startGame = document.getElementById('game-start')

    const declareWinner = document.getElementById('declare-winner')
    const submitWinnerButton = document.getElementById('winner-button')
    const winnerInput = document.getElementById('winner')

    const blindContainer = document.getElementById('blind-value')

    const gameContainer = document.getElementById('game')
    const gameEndContainer = document.getElementById('game-end')

    declareWinner.hidden = true
    gameEndContainer.hidden = true

    document.getElementById('start-game').addEventListener('click', event => {
        startGame.hidden = true
        declareWinner.hidden = false

        const numberOfPlayers = document.getElementById('player-count').value

        if (window['WebSocket']) {
            const conn = new WebSocket('ws://' + document.location.host + '/ws')

            submitWinnerButton.onclick = event => {
                conn.send(winnerInput.value)
                gameEndContainer.hidden = false
                gameContainer.hidden = true
            }

            conn.onclose = evt => {
                blindContainer.innerText = 'Connection closed'
            }

            conn.onmessage = evt => {
                blindContainer.innerText = evt.data
            }

            conn.onopen = function () {
                conn.send(numberOfPlayers)
            }
        }
    })
</script>
</html>
```

주요 변경 사항은 플레이어 수를 입력하는 섹션과 블라인드 값을 표시하는 섹션을 가져 오는 것입니다. 게임의 단계에 따라 사용자 인터페이스를 표시하거나 숨기는 약간의 logic이 있습니다.

'conn.onmessage'를 통해 수신하는 모든 메시지는 블라인드 경고라고 가정하므로 이에 따라 'blindContainer.innerText'를 설정합니다.

앞 장에서는 `Game`에 대한 아이디어를 소개하여 CLI 코드가 `Game`을 호출할 수 있고 이것이 블라인드 알림 스케쥴링을 포함한 모든 사항을 처리할 수 있도록 했습니다.

```go
type Game interface {
	Start(numberOfPlayers int)
	Finish(winner string)
}
```

사용자가 CLI에서 플레이어 수를 묻는 메시지를 받으면 게임을 `Start`하여 블라인드 경고를 시작할 것이고, 사용자가 승자를 선언하면 `Finish`가 될 것입니다. 이것은 우리가 현재 가지고있는 것과 동일한 요구 사항이며 입력을 얻는 다른 방법입니다. 가능하다면 이 개념을 재사용해야합니다.

`Game`의 "실제"구현은 `TexasHoldem`입니다.

```go
type TexasHoldem struct {
	alerter BlindAlerter
	store   PlayerStore
}
```

`BlindAlerter`에서 `TexasHoldem`을 보냄으로써, _언제나_ 블라인드 경고가 전송되도록 스케쥴할 수 있습니다.

```go
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int)
}
```

그리고 상기시켜 드리자면, 이것은 우리가 CLI에서 사용했던 `BlindAlerter` 구현입니다.

```go
func StdOutAlerter(duration time.Duration, amount int) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(os.Stdout, "Blind is now %d\n", amount)
	})
}
```

이것은 _항상 `os.Stdout`에 경고를 보내기_를 원하기 때문에 CLI에서는 작동하지만 웹 서버에서는 작동하지 않습니다. 모든 요청에 대해 새로운 `http.ResponseWriter`를 얻은 다음`*websocket.Conn`으로 업그레이드합니다. 따라서 종속성을 구성 할 때 우리의 경고가 어디로 이동해야할 지 알 수 없습니다.

따라서 그것을 경고를 위한 목적지를 가지게 하여 웹 서버에서도 재사용할 수 있도록 BlindAlerter.ScheduleAlertAt`을 변경해야합니다.

BlindAlerter.go를 열고 'to io.Writer' 매개 변수를 추가합니다.

```go
type BlindAlerter interface {
	ScheduleAlertAt(duration time.Duration, amount int, to io.Writer)
}

type BlindAlerterFunc func(duration time.Duration, amount int, to io.Writer)

func (a BlindAlerterFunc) ScheduleAlertAt(duration time.Duration, amount int, to io.Writer) {
	a(duration, amount, to)
}
```

`StdoutAlerter`의 아이디어는 우리의 새 모델에 맞지 않으므로 이름을 `Alerter`로 변경합시다.

```go
func Alerter(duration time.Duration, amount int, to io.Writer) {
	time.AfterFunc(duration, func() {
		fmt.Fprintf(to, "Blind is now %d\n", amount)
	})
}
```

컴파일을 시도하면 목적지없이`ScheduleAlertAt`을 호출하기 때문에`TexasHoldem`에서 실패하여 다시 컴파일을 하기 위해서는 _목적지를 `os.Stdout`로 하드 코딩_합니다.

테스트를 실행하면`SpyBlindAlerter`가 더 이상`BlindAlerter`를 구현하지 않기 때문에 실패합니다.`ScheduleAlertAt`의 signature를 업데이트하여이 문제를 해결하고 테스트를 실행하면 여전히 녹색이어야합니다.

'TexasHoldem'이 블라인드 알림을 보낼 위치를 알게 하는 것은 말이되지 않습니다. 이제 게임을 시작할 때 _어디_로 경고 가야할 지 선언하기 위해 `Game`을 업데이트하겠습니다.

```go
type Game interface {
	Start(numberOfPlayers int, alertsDestination io.Writer)
	Finish(winner string)
}
```

컴파일러가 수정해야 할 사항을 알려줍니다. 이 변화는 그렇게 나쁘지 않습니다.

-`TexasHoldem`을 업데이트하여`Game`을 올바르게 구현합니다.
-`CLI`에서 게임을 시작할 때 `out` 속성 (`cli.game.Start (numberOfPlayers, cli.out)`)을 전달합니다.
-`TexasHoldem`의 테스트에서`game.Start (5, ioutil.Discard)`를 사용하여 컴파일 문제를 수정하고 경고 출력을 버리도록 구성합니다.

모든 것이 올바르게 되었다면 모든 것이 녹색이어야합니다! 이제`Server` 내에서`Game`을 사용해 볼 수 있습니다.

## 테스트를 먼저 작성하세요.

`CLI`와`Server`의 요구 사항은 동일합니다! 단순히 전달 메커니즘이 다릅니다.

영감을 얻기위해 'CLI'테스트를 살펴 보겠습니다.

```go
t.Run("start game with 3 players and finish game with 'Chris' as winner", func(t *testing.T) {
    game := &GameSpy{}

    out := &bytes.Buffer{}
    in := userSends("3", "Chris wins")

    poker.NewCLI(in, out, game).PlayPoker()

    assertMessagesSentToUser(t, out, poker.PlayerPrompt)
    assertGameStartedWith(t, game, 3)
    assertFinishCalledWith(t, game, "Chris")
})
```

'GameSpy'를 사용하여 유사한 결과를 테스트 할 수있을 것 같습니다.

이전 websocket 테스트를 다음으로 교체하십시오.

```go
t.Run("start a game with 3 players and declare Ruth the winner", func(t *testing.T) {
    game := &poker.GameSpy{}
    winner := "Ruth"
    server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
    ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

    defer server.Close()
    defer ws.Close()

    writeWSMessage(t, ws, "3")
    writeWSMessage(t, ws, winner)

    time.Sleep(10 * time.Millisecond)
    assertGameStartedWith(t, game, 3)
    assertFinishCalledWith(t, game, winner)
})
```

-논의했듯이 스파이`Game`을 만들고`mustMakePlayerServer`에 전달합니다 (이를 지원하도록 helper를 업데이트해야합니다).
-그런 다음 게임에 위해 웹 소켓 메시지를 보냅니다.
-마지막으로 우리가 기대하는대로 게임이 시작되고 끝났다고 assert합니다.

## 테스트를 실행하세요.

다른 테스트에서`mustMakePlayerServer` 와 관련한 많은 컴파일 오류를 가질 것입니. Unexported 변수인 'dummyGame'을 도입하고 컴파일하지 않는 모든 테스트를 통해 사용합니다.

```go
var (
	dummyGame = &GameSpy{}
)
```

The final error is where we are trying to pass in `Game` to `NewPlayerServer` but it doesn't support it yet

```
./server_test.go:21:38: too many arguments in call to "github.com/quii/learn-go-with-tests/WebSockets/v2".NewPlayerServer
	have ("github.com/quii/learn-go-with-tests/WebSockets/v2".PlayerStore, "github.com/quii/learn-go-with-tests/WebSockets/v2".Game)
	want ("github.com/quii/learn-go-with-tests/WebSockets/v2".PlayerStore)
```

## Write the minimal amount of code for the test to run and check the failing test output

Just add it as an argument for now just to get the test running

```go
func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
```

마침내!

```
=== RUN   TestGame/start_a_game_with_3_players_and_declare_Ruth_the_winner
--- FAIL: TestGame (0.01s)
    --- FAIL: TestGame/start_a_game_with_3_players_and_declare_Ruth_the_winner (0.01s)
    	server_test.go:146: wanted Start called with 3 but got 0
    	server_test.go:147: expected finish called with 'Ruth' but got ''
FAIL
```

## Write enough code to make it pass

We need to add `Game` as a field to `PlayerServer` so that it can use it when it gets requests.

```go
type PlayerServer struct {
	store PlayerStore
	http.Handler
	template *template.Template
	game Game
}
```

(We already have a method called `game` so rename that to `playGame`)

Next lets assign it in our constructor

```go
func NewPlayerServer(store PlayerStore, game Game) (*PlayerServer, error) {
	p := new(PlayerServer)

	tmpl, err := template.ParseFiles(htmlTemplatePath)

	if err != nil {
		return nil, fmt.Errorf("problem opening %s %v", htmlTemplatePath, err)
	}

	p.game = game

	// etc
```

Now we can use our `Game` within `webSocket`.

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	conn, _ := wsUpgrader.Upgrade(w, r, nil)

	_, numberOfPlayersMsg, _ := conn.ReadMessage()
	numberOfPlayers, _ := strconv.Atoi(string(numberOfPlayersMsg))
	p.game.Start(numberOfPlayers, ioutil.Discard) //todo: Don't discard the blinds messages!

	_, winner, _ := conn.ReadMessage()
	p.game.Finish(string(winner))
}
```

Hooray! The tests pass.

We are not going to send the blind messages anywhere _just yet_ as we need to have a think about that. When we call `game.Start` we send in `ioutil.Discard` which will just discard any messages written to it.

For now start the web server up. You'll need to update the `main.go` to pass a `Game` to the `PlayerServer`

```go
func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := poker.NewFileSystemPlayerStore(db)

	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	game := poker.NewTexasHoldem(poker.BlindAlerterFunc(poker.Alerter), store)

	server, err := poker.NewPlayerServer(store, game)

	if err != nil {
		log.Fatalf("problem creating player server %v", err)
	}

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
```

Discounting the fact we're not getting blind alerts yet, the app does work! We've managed to re-use `Game` with `PlayerServer` and it has taken care of all the details. Once we figure out how to send our blind alerts through to the web sockets rather than discarding them it _should_ all work.

Before that though, let's tidy up some code.

## Refactor

The way we're using WebSockets is fairly basic and the error handling is fairly naive, so I wanted to encapsulate that in a type just to remove that messiness from the server code. We may wish to revisit it later but for now this'll tidy things up a bit

```go
type playerServerWS struct {
	*websocket.Conn
}

func newPlayerServerWS(w http.ResponseWriter, r *http.Request) *playerServerWS {
	conn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("problem upgrading connection to WebSockets %v\n", err)
	}

	return &playerServerWS{conn}
}

func (w *playerServerWS) WaitForMsg() string {
	_, msg, err := w.ReadMessage()
	if err != nil {
		log.Printf("error reading from websocket %v\n", err)
	}
	return string(msg)
}
```

Now the server code is a bit simplified

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayersMsg)
	p.game.Start(numberOfPlayers, ioutil.Discard) //todo: Don't discard the blinds messages!

	winner := ws.WaitForMsg()
	p.game.Finish(winner)
}
```

Once we figure out how to not discard the blind messages we're done.

### Let's _not_ write a test!

Sometimes when we're not sure how to do something, it's best just to play around and try things out! Make sure your work is committed first because once we've figured out a way we should drive it through a test.

The problematic line of code we have is

```go
p.game.Start(numberOfPlayers, ioutil.Discard) //todo: Don't discard the blinds messages!
```

We need to pass in an `io.Writer` for the game to write the blind alerts to.

Wouldn't it be nice if we could pass in our `playerServerWS` from before? It's our wrapper around our WebSocket so it _feels_ like we should be able to send that to our `Game` to send messages to.

Give it a go:

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayersMsg)
	p.game.Start(numberOfPlayers, ws)
	//etc...
```

The compiler complains

```
./server.go:71:14: cannot use ws (type *playerServerWS) as type io.Writer in argument to p.game.Start:
	*playerServerWS does not implement io.Writer (missing Write method)
```

It seems the obvious thing to do, would be to make it so `playerServerWS` _does_ implement `io.Writer`. To do so we use the underlying `*websocket.Conn` to use `WriteMessage` to send the message down the websocket

```go
func (w *playerServerWS) Write(p []byte) (n int, err error) {
	err = w.WriteMessage(websocket.TextMessage, p)

	if err != nil {
		return 0, err
	}

	return len(p), nil
}
```

This seems too easy! Try and run the application and see if it works.

Beforehand edit `TexasHoldem` so that the blind increment time is shorter so you can see it in action

```go
blindIncrement := time.Duration(5+numberOfPlayers) * time.Second // (rather than a minute)
```

You should see it working! The blind amount increments in the browser as if by magic.

Now let's revert the code and think how to test it. In order to _implement_ it all we did was pass through to `StartGame` was `playerServerWS` rather than `ioutil.Discard` so that might make you think we should perhaps spy on the call to verify it works.

Spying is great and helps us check implementation details but we should always try and favour testing the _real_ behaviour if we can because when you decide to refactor it's often spy tests that start failing because they are usually checking implementation details that you're trying to change.

Our test currently opens a websocket connection to our running server and sends messages to make it do things. Equally we should be able to test the messages our server sends back over the websocket connection.

## Write the test first

We'll edit our existing test.

Currently our `GameSpy` does not send any data to `out` when you call `Start`. We should change it so we can configure it to send a canned message and then we can check that message gets sent to the websocket. This should give us confidence that we have configured things correctly whilst still exercising the real behaviour we want.

```go
type GameSpy struct {
	StartCalled     bool
	StartCalledWith int
	BlindAlert      []byte

	FinishedCalled   bool
	FinishCalledWith string
}
```

Add `BlindAlert` field.

Update `GameSpy` `Start` to send the canned message to `out`.

```go
func (g *GameSpy) Start(numberOfPlayers int, out io.Writer) {
	g.StartCalled = true
	g.StartCalledWith = numberOfPlayers
	out.Write(g.BlindAlert)
}
```

This now means when we exercise `PlayerServer` when it tries to `Start` the game it should end up sending messages through the websocket if things are working right.

Finally we can update the test

```go
t.Run("start a game with 3 players, send some blind alerts down WS and declare Ruth the winner", func(t *testing.T) {
    wantedBlindAlert := "Blind is 100"
    winner := "Ruth"

    game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
    server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
    ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

    defer server.Close()
    defer ws.Close()

    writeWSMessage(t, ws, "3")
    writeWSMessage(t, ws, winner)

    time.Sleep(10 * time.Millisecond)
    assertGameStartedWith(t, game, 3)
    assertFinishCalledWith(t, game, winner)

    _, gotBlindAlert, _ := ws.ReadMessage()

    if string(gotBlindAlert) != wantedBlindAlert {
        t.Errorf("got blind alert %q, want %q", string(gotBlindAlert), wantedBlindAlert)
    }
})
```

- We've added a `wantedBlindAlert` and configured our `GameSpy` to send it to `out` if `Start` is called.
- We hope it gets sent in the websocket connection so we've added a call to `ws.ReadMessage()` to wait for a message to be sent and then check it's the one we expected.

## Try to run the test

You should find the test hangs forever. This is because `ws.ReadMessage()` will block until it gets a message, which it never will.


## Write the minimal amount of code for the test to run and check the failing test output

We should never have tests that hang so let's introduce a way of handling code that we want to timeout.

```go
func within(t testing.TB, d time.Duration, assert func()) {
	t.Helper()

	done := make(chan struct{}, 1)

	go func() {
		assert()
		done <- struct{}{}
	}()

	select {
	case <-time.After(d):
		t.Error("timed out")
	case <-done:
	}
}
```

What `within` does is take a function `assert` as an argument and then runs it in a go routine. If/When the function finishes it will signal it is done via the `done` channel.

While that happens we use a `select` statement which lets us wait for a channel to send a message. From here it is a race between the `assert` function and `time.After` which will send a signal when the duration has occurred.

Finally I made a helper function for our assertion just to make things a bit neater

```go
func assertWebsocketGotMsg(t *testing.T, ws *websocket.Conn, want string) {
	_, msg, _ := ws.ReadMessage()
	if string(msg) != want {
		t.Errorf(`got "%s", want "%s"`, string(msg), want)
	}
}
```

Here's how the test reads now

```go
t.Run("start a game with 3 players, send some blind alerts down WS and declare Ruth the winner", func(t *testing.T) {
    wantedBlindAlert := "Blind is 100"
    winner := "Ruth"

    game := &GameSpy{BlindAlert: []byte(wantedBlindAlert)}
    server := httptest.NewServer(mustMakePlayerServer(t, dummyPlayerStore, game))
    ws := mustDialWS(t, "ws"+strings.TrimPrefix(server.URL, "http")+"/ws")

    defer server.Close()
    defer ws.Close()

    writeWSMessage(t, ws, "3")
    writeWSMessage(t, ws, winner)

    time.Sleep(tenMS)

    assertGameStartedWith(t, game, 3)
    assertFinishCalledWith(t, game, winner)
    within(t, tenMS, func() { assertWebsocketGotMsg(t, ws, wantedBlindAlert) })
})
```

Now if you run the test...

```
=== RUN   TestGame
=== RUN   TestGame/start_a_game_with_3_players,_send_some_blind_alerts_down_WS_and_declare_Ruth_the_winner
--- FAIL: TestGame (0.02s)
    --- FAIL: TestGame/start_a_game_with_3_players,_send_some_blind_alerts_down_WS_and_declare_Ruth_the_winner (0.02s)
    	server_test.go:143: timed out
    	server_test.go:150: got "", want "Blind is 100"
```

## Write enough code to make it pass

Finally we can now change our server code so it sends our WebSocket connection to the game when it starts

```go
func (p *PlayerServer) webSocket(w http.ResponseWriter, r *http.Request) {
	ws := newPlayerServerWS(w, r)

	numberOfPlayersMsg := ws.WaitForMsg()
	numberOfPlayers, _ := strconv.Atoi(numberOfPlayersMsg)
	p.game.Start(numberOfPlayers, ws)

	winner := ws.WaitForMsg()
	p.game.Finish(winner)
}
```

## Refactor

The server code was a very small change so there's not a lot to change here but the test code still has a `time.Sleep` call because we have to wait for our server to do its work asynchronously.

We can refactor our helpers `assertGameStartedWith` and `assertFinishCalledWith` so that they can retry their assertions for a short period before failing.

Here's how you can do it for `assertFinishCalledWith` and you can use the same approach for the other helper.

```go
func assertFinishCalledWith(t testing.TB, game *GameSpy, winner string) {
	t.Helper()

	passed := retryUntil(500*time.Millisecond, func() bool {
		return game.FinishCalledWith == winner
	})

	if !passed {
		t.Errorf("expected finish called with %q but got %q", winner, game.FinishCalledWith)
	}
}
```

Here is how `retryUntil` is defined

```go
func retryUntil(d time.Duration, f func() bool) bool {
	deadline := time.Now().Add(d)
	for time.Now().Before(deadline) {
		if f() {
			return true
		}
	}
	return false
}
```

## Wrapping up

Our application is now complete. A game of poker can be started via a web browser and the users are informed of the blind bet value as time goes by via WebSockets. When the game finishes they can record the winner which is persisted using code we wrote a few chapters ago. The players can find out who is the best (or luckiest) poker player using the website's `/league` endpoint.

Through the journey we have made mistakes but with the TDD flow we have never been very far away from working software. We were free to keep iterating and experimenting.

The final chapter will retrospect on the approach, the design we've arrived at and tie up some loose ends.

We covered a few things in this chapter

### WebSockets

- Convenient way of sending messages between clients and servers that does not require the client to keep polling the server. Both the client and server code we have is very simple.
- Trivial to test, but you have to be wary of the asynchronous nature of the tests

### Handling code in tests that can be delayed or never finish

- Create helper functions to retry assertions and add timeouts.
- We can use go routines to ensure the assertions don't block anything and then use channels to let them signal that they have finished, or not.
- The `time` package has some helpful functions which also send signals via channels about events in time so we can set timeouts
