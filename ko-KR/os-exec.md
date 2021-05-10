# OS Exec

**[이 챕터의 모든 코드는 여기에서 확인할 수 있다.](https://github.com/quii/learn-go-with-tests/tree/main/q-and-a/os-exec)**

[keith6014](https://www.reddit.com/user/keith6014) 라는 유저가 [reddit](https://www.reddit.com/r/golang/comments/aaz8ji/testdata_and_function_setup_help/) 에 다음과 같이 물었다.

> 현재 XML 데이터를 os/exec.Command()를 사용하여 생성하고 있습니다. 해당 콜은 GetData()라는 함수에서 이후에 실행됩니다.

> GetData()를 테스트 하기 위해 직접 테스트 데이터를 만들었습니다.

> 제 _test.go 는 GetData()를 호출하는 TestGetData 를 가지고 있지만, 해당 콜은 os.exec을 사용합니다. 이를 대신해 저는 제 테스트 데이터를 사용하고 싶습니다.

> 이를 위한 가장 좋은 방법은 무엇일까요? GetData()를 호출 시에 "테스트" 플래그 모드를 설정하여 파일을 읽게 해야 할까요? 예: GetData(mode string)

몇 가지 사항들을 이야기해보자.

- 무엇인가 테스트하기 힘든 상황은 대게 관심사의 세분화가 (separation of concenrs) 올바르지 않기 때문이다.
- "테스트 모드"를 코드에 추가하지 말고 대신에 [Dependency Injection](/dependency-injection.md) 를 사용하여 종속성 (dependencies)을 모델링하고 우려하는 부분을 세분화 할 수 있게 해야 한다.

아마 다음의 코드와 같은 형식을 이루지 않을까 하고 짐작해보았다.

```go
type Payload struct {
	Message string `xml:"message"`
}

func GetData() string {
	cmd := exec.Command("cat", "msg.xml")

	out, _ := cmd.StdoutPipe()
	var payload Payload
	decoder := xml.NewDecoder(out)

	// these 3 can return errors but I'm ignoring for brevity
	cmd.Start()
	decoder.Decode(&payload)
	cmd.Wait()

	return strings.ToUpper(payload.Message)
}
```

- `exec.Command` 을 사용하여 외부 커맨드를 해당 프로세스에 사용할 수 있다.
- `io.ReadCloser` 값을 리턴하는 `cmd.StdoutPipe` 의 결과값을 저장한다 (이후에 굉장히 중요해진다)
- 나머지 코드는 [excellent documentation](https://golang.org/pkg/os/exec/#example_Cmd_StdoutPipe)와 거의 유사하다.
    - stdout의 결과값을 `io.ReadCloser` 에 저장한후 `Start` 명령어를 사용한다. 그리고 `Wait` 을 호출하여 모든 데이터가 읽혀지기를 기다린다. 이 두 호출 사이에 데이터를 `Payload` 구조체로 디코딩 (decoding) 한다.

다음은 `msg.xml`에 포함 된 내용이다.

```xml
<payload>
    <message>Happy New Year!</message>
</payload>
```

실제로 작동하는지를 보여주기 위해 다음과 같은 간단한 테스트를 작성해보았다.

```go
func TestGetData(t *testing.T) {
	got := GetData()
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

## Testable code

테스트 가능한 코드는 분리되었으며 단 한 가지의 목적을 가진다. 나에게는 해당 코드에 두 가지 주요 관심사 (concerns)가 존재한다고 보는데 이는 다음과 같다.

1. raw XML 데이터 검색 (retrieving)
2. XML 데이터의 디코딩과 이를 우리의 비즈니스 로직에 적용하는 것 (해당 코드의 경우 `message` 의 `strings.ToUpper`).

첫 번째 부분의 경우 단지 표준 lib에서 예제를 복사하는 것에 불과하다.

The second part is where we have our business logic and by looking at the code we can see where the "seam" in our logic starts; it's where we get our `io.ReadCloser`. We can use this existing abstraction to separate concerns and make our code testable.
두 번째 부분은 비즈니스 로직이 있는 곳이며 코드를 살펴보면 로직의 "이음새 (seam)"가 시작되는 위치를 알 수 있는데 이에 따라 `io.ReadCloser`를 얻을 수 있게 된다. 이 기존 추상화 (abstraction)를 사용하여 관심 사항을 분리하고 코드의 테스트를 가능하게 만들 수 있다.

**GetData의 문제점은 비즈니스 로직이 XML을 얻는 수단과 결합되어 있다는 것이다. 디자인을 더 좋게 만들기 위해 우리는 그것들을 분리해야 한다.**

`TestGetData` 는 두 가지 관심사 간의 통합 테스트 역할을 할 수 있음으로 계속 작동하는지 확인하기 위해 유지한다.

새로 분리된 코드는 다음과 같다.

```go
type Payload struct {
	Message string `xml:"message"`
}

func GetData(data io.Reader) string {
	var payload Payload
	xml.NewDecoder(data).Decode(&payload)
	return strings.ToUpper(payload.Message)
}

func getXMLFromCommand() io.Reader {
	cmd := exec.Command("cat", "msg.xml")
	out, _ := cmd.StdoutPipe()

	cmd.Start()
	data, _ := ioutil.ReadAll(out)
	cmd.Wait()

	return bytes.NewReader(data)
}

func TestGetDataIntegration(t *testing.T) {
	got := GetData(getXMLFromCommand())
	want := "HAPPY NEW YEAR!"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
```

이제 `GetData` 는 테스트 가능하게 만든 `io.Reader` 에서 입력을 가져오므로 더 이상 데이터를 검색하는 방법에 대해 걱정하지 않아도 된다; 사람들은 `Getdata` 함수를 `io.Reader`(굉장히 흔한)를 반환하는 여러 함수에 재사용 할 수 있다. 예를 들어 command line 대신 URL에서 XML 가져오기 시작할 수 있다.

```go
func TestGetData(t *testing.T) {
	input := strings.NewReader(`
<payload>
    <message>Cats are the best animal</message>
</payload>`)

	got := GetData(input)
	want := "CATS ARE THE BEST ANIMAL"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

```

다음은 `GetData`에 대한 단위 테스트의 예다.

Go 테스트 내에서 우려 사항을 분리하고 기존 추상화를 사용하면 중요한 비즈니스 로직을 쉽게 사용할 수 있다.
