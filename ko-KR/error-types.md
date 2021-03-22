# 에러 타입(Error types)

**[이 글에 소개된 모든 코드는 여기서 확인할 수 있다](https://github.com/MiryangJung/learn-go-with-tests-ko/tree/master/q-and-a/error-types)**

**자신만의 에러 타입을 만드는 것은, 코드를 정돈하고, 사용하고, 테스트하기 쉽게 만드는 우아한 방법일 수 있다.**

Gopher Slack에서 Pedro가 다음과 같이 물었다.

> `fmt.Errorf("%s must be foo, got %s", bar, baz)`와 같은 에러를 만들고 있는 경우, string 값을 비교하지 않고 동일한지 테스트(test equality)할 수 있는 방법이 있는가?

이 아이디어를 살펴볼 때 도움이 될 함수를 만들어보자.

```go
// DumbGetter will get the string body of url if it gets a 200
func DumbGetter(url string) (string, error) {
	res, err := http.Get(url)

	if err != nil {
		return "", fmt.Errorf("problem fetching from %s, %v", url, err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("did not get 200 from %s, got %d", url, res.StatusCode)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body) // ignoring err for brevity

	return string(body), nil
}
```

실패할 수 있는 이유가 여러가지인 테스트를 작성하는 건 일반적이지 않으며, 각 시나리오를 올바르게 처리하고 싶다.

Pedro가 말했듯, status 에러 테스트를 그렇게 작성할 수 있다.

```go
t.Run("when you don't get a 200 you get a status error", func(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusTeapot)
	}))
	defer svr.Close()

	_, err := DumbGetter(svr.URL)

	if err == nil {
		t.Fatal("expected an error")
	}

	want := fmt.Sprintf("did not get 200 from %s, got %d", svr.URL, http.StatusTeapot)
	got := err.Error()

	if got != want {
		t.Errorf(`got "%v", want "%v"`, got, want)
	}
})
```

이 테스트는 항상 `StatusTeapot`을 반환하는 서버를 만들고, 해당 URL을 `DumbGetter`의 인자로 사용하여 `200` 이외의 응답을 올바르게 처리할 수 있는 지 확인한다.

## 이 테스트 방법의 문제점

이 책은 *테스트에 귀를 기울이는 것을* 강조하려고 하는데, 이 테스트는 좋은 *느낌은* 아니다:

- 테스트를 위해 프로덕션 코드와 같은 문자열을 만들었다
- 읽고 쓰기 귀찮다
- 정확한 에러 메세지 문자열은 *실제로 관련 있는가* ?

이것은 무엇을 시사하고 있는가? 테스트의 인간공학(ergonomics)은 코드를 사용하려고 하는 다른 코드를 반영한다.

코드 사용자는 반환되는 특정 종류의 에러에 어떻게 반응할까? 그들이 할 수 있는 최선은 에러가 발생하기 극도로 쉽고, 무시무시하게 쓰여진 에러 문자열을 살펴보는 것뿐이다.

## 우리가 해야하는 것

TDD와 함께라면, 다음과 같은 사고방식을 갖게 되는 이점이 있다:

> *나는* 이 코드를 어떻게 사용하고 싶은가?

우리가 `DumbGetter`에 할 수 있는 것은, 타입 시스템을 이용해 사용자로 하여금 발생한 에러의 종류를 이해하는 방법을 제공하는 것이다.

만약 `DumbGetter`가 다음과 같은 타입을 반환한다면 어떨까?

```go
type BadStatusError struct {
	URL    string
	Status int
}
```

마법의 문자열이 아니라, 실제로 사용되는 *데이터가* 있다.

이 필요성을 반영하도록 기존 테스트를 바꿔보자.

```go
t.Run("when you don't get a 200 you get a status error", func(t *testing.T) {

	svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusTeapot)
	}))
	defer svr.Close()

	_, err := DumbGetter(svr.URL)

	if err == nil {
		t.Fatal("expected an error")
	}

	got, isStatusErr := err.(BadStatusError)

	if !isStatusErr {
		t.Fatalf("was not a BadStatusError, got %T", err)
	}

	want := BadStatusError{URL: svr.URL, Status: http.StatusTeapot}

	if got != want {
		t.Errorf("got %v, want %v", got, want)
	}
})
```

`BadStatusError`가 에러 인터페이스를 구현하도록 만들어야 한다.

```go
func (b BadStatusError) Error() string {
    return fmt.Sprintf("did not get 200 from %s, got %d", b.URL, b.Status)
}
```

### 이 테스트는 무엇을 할까?

에러의 정확한 문자열을 체크하는 대신, 에러에 대응하는 [타입 assertion(type assertion)](https://tour.golang.org/methods/15)을 실행해 에러가 `BadStatusError`인지 확인한다. 이는 에러의 *종류를* 분명하게 하고자 하는 염원을 반영하고 있다. assertion이 통과되었다고 가정하면, 에러의 프로퍼티(property)가 올바른지 확인할 수 있다.

테스트를 실행하면, 올바른 종류의 에러가 반환되지 않았다는 것을 알 수 있다.

```
--- FAIL: TestDumbGetter (0.00s)
    --- FAIL: TestDumbGetter/when_you_dont_get_a_200_you_get_a_status_error (0.00s)
    	error-types_test.go:56: was not a BadStatusError, got *errors.errorString
```

우리가 만든 타입을 사용하도록 에러 처리 코드를 수정해, `DumbGetter`을 고치자.

```go
if res.StatusCode != http.StatusOK {
    return "", BadStatusError{URL: url, Status: res.StatusCode}
}
```

이 수정은 몇 가지 *현실적으로 긍정적인 효과를* 불러온다.

- `DumbGetter` 함수가 보다 단순해졌다. 에러 문자열의 복잡함과 관련 없어지고, `BadStatusError`를 만들 뿐이다.
- 이제 테스트는 코드 사용자가 단순한 로깅보다 더 정교한 에러 처리를 원할 경우, 사용자가 *할 수 있는* 작업을 반영(그리고 문서화)한다. 타입 assertion을 하는 것만으로도, 에러의 프로퍼티에 쉽게 접근할 수 있다.
- 그럼에도 여전히 "단순한" `error`이기 때문에, 사용자가 원한다면 다른 `error`처럼 콜 스택에 쌓거나, 로그로 기록할 수 있다.

## 정리

다수의 에러 조건을 테스트하고 있는 스스로를 눈치 챘다면, 에러 메세지를 비교한다는 덫에 빠지지 마라.

이는 불완전하고 읽고 쓰기 어려운 테스트로 이어지고, 발생한 에러의 종류에 따라 다른 작업을 시작할 필요가 있는 경우, 코드 사용자가 겪게 될 어려움을 반영한다.

항상 *당신이* 어떻게 코드를 사용하고 싶은 지를 테스트에 반영하라. 이러한 관점에서, 에러를 캡슐화하기 위해 에러 타입을 작성하는 것을 고려하라. 이를 통해 코드 사용자가 다양한 종류의 에러를 보다 쉽게 처리 할 수 있으며, 에러 처리 코드를 보다 간단하고 읽기 쉽게 쓸 수 있게 된다.

## 보충

Go 1.13에서는, 표준 라이브러리의 에러를 다루는 새로운 방법이 있다. [Go 블로그](https://blog.golang.org/go1.13-errors)

```go
t.Run("when you don't get a 200 you get a status error", func(t *testing.T) {

    svr := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
        res.WriteHeader(http.StatusTeapot)
    }))
    defer svr.Close()

    _, err := DumbGetter(svr.URL)

    if err == nil {
        t.Fatal("expected an error")
    }

    var got BadStatusError
    isBadStatusError := errors.As(err, &got)
    want := BadStatusError{URL: svr.URL, Status: http.StatusTeapot}

    if !isBadStatusError {
        t.Fatalf("was not a BadStatusError, got %T", err)
    }

    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
})
```

이 경우, [`errors.As`](https://golang.org/pkg/errors/#example_As)를 사용해 에러를 커스텀 타입으로 추출한다. 성공을 나타내기 위해 `bool`을 반환하고, 이를 `got`으로 추출한다.

