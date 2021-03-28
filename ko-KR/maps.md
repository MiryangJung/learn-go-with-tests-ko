# 맵

**[이 챕터의 모든 코드는 여기에서 확인할 수 있다](https://github.com/quii/learn-go-with-tests/tree/main/maps)**

[배열과 슬라이스](arrays-and-slices.md)에서 값을 순서대로 저장하는 방법을 다뤘다. 이번에는 항목을 `key`에 따라 저장하고 이렇게 저장한 `key`를 빠르게 찾는 방법을 살펴볼 것이다.

맵은 사전과 비슷한 방식으로 항목을 저장할 수 있어서, `key`는 단어이고 `value`는 정의라는 식으로 생각할 수 있다. 그러므로 우리만의 사전을 만드는 것이 맵을 배우는 가장 좋은 방법이지 않을까?

우선 몇 개의 단어와 이들의 정의가 있는 사전이 있다고 가정해 보자. 단어를 검색하면, 사전은 그 단어의 정의를 반환해야한다.

## 테스트부터 작성하기

`dictionary_test.go` 는

```go
package main

import "testing"

func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    want := "this is just a test"

    if got != want {
        t.Errorf("got %q want %q given, %q", got, want, "test")
    }
}
```

맵을 선언하는 것은 배열을 선언하는 것과 비슷하지만 다음과 같은 점에 있어서 다르다. 맵을 선언하려면 `map`이라는 키워드로 시작하고 두 개의 타입이 있어야한다. 첫번째 타입은 키 타입으로 `[]` 안에 쓰여 있다. 두번째는 값 타입으로, `[]` 바로 다음에 온다.

키 타입은 특별하다. 키 타입에는 오직 비교 가능한 타입만이 올 수 있는데 왜냐하면 두개의 키가 동일한지 판별할 수 없다면 올바른 값을 가져왔는지 확신할 수 있는 방법이 없기 때문이다. [언어 명세]에 비교 가능한 타입이 자세하게 설명되어 있다.

반면에 값 타입으로 무엇이든 원하는 값이 가능하다. 심지어 또다른 맵도 가능하다.

테스트의 나머지는 친숙할 것이다.

## 테스트 실행해보기

`go test`를 실행하면 컴파일러는 `./dictionary_test.go:8:9: undefined: Search`와 함께 실패할 것이다.

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

In `dictionary.go`

```go
package main

func Search(dictionary map[string]string, word string) string {
    return ""
}
```

이번에는 테스트가 *명확한 에러 메시지*와 함께 실패할 것이다.

`dictionary_test.go:12: got '' want 'this is just a test' given, 'test'`.

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func Search(dictionary map[string]string, word string) string {
    return dictionary[word]
}
```

맵에서 값을 가져오는 것은 `map[key]` 배열에서 값을 가져오는 것과 동일하다.

## 리팩터링 하기

```go
func TestSearch(t *testing.T) {
    dictionary := map[string]string{"test": "this is just a test"}

    got := Search(dictionary, "test")
    want := "this is just a test"

    assertStrings(t, got, want)
}

func assertStrings(t testing.TB, got, want string) {
    t.Helper()

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

`assertStrings` 헬퍼를 만듦으로써 구현이 보다 일반적이게 되도록 만들었다.

### 커스텀 타입 사용하기

map에 대한 새로운 타입을 만들고 `Search` 함수를 만듦으로써 위에서 작성한 사전의 사용성을 개선한다.

In `dictionary_test.go`:

```go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    got := dictionary.Search("test")
    want := "this is just a test"

    assertStrings(t, got, want)
}
```

여기서 `Dictionary` 타입을 도입했는데 아직 선언하지 않았다. 그리고 `Dictionary` 인스턴스의 `Search` 함수를 호출하였다.

`assertStrings`를 변경할 필요는 없다.

`dictionary.go` 에서:

```go
type Dictionary map[string]string

func (d Dictionary) Search(word string) string {
    return d[word]
}
```

여기서 `Dictionary` 타입을 생성했는데, `map`을 감싸는 얇은 래퍼로 동작합니다. 새로 정의한 커스텀 타입과 함께, `Search` 함수를 생성할 수 있다.

## 테스트부터 작성하기

기본 검색은 구현하기 매우 쉬웠다. 그러나 만약 딕셔너리에 없는 단어를 검색한다면 어떻게 될까?

실제로 아무 것도 가져올 수 없다. 이래도 프로그램이 계속 동작하게하기 때문에 괜찮지만 더 나은 방법이 있다. `Search` 함수는 단어가 사전에 존재하지 않는다고 알려줄 수 있다. 이 방법으로 사용자가 단어가 존재하지 않는건지 아니면 단지 정의가 없는건지 궁금해하지 않게 된다 (이런 사전은 유용하지 않아 보인다. 그렇지만 다른 사례에서 키가 될 시나리오이다).

```go
func TestSearch(t *testing.T) {
    dictionary := Dictionary{"test": "this is just a test"}

    t.Run("known word", func(t *testing.T) {
        got, _ := dictionary.Search("test")
        want := "this is just a test"

        assertStrings(t, got, want)
    })

    t.Run("unknown word", func(t *testing.T) {
        _, err := dictionary.Search("unknown")
        want := "could not find the word you were looking for"

        if err == nil {
            t.Fatal("expected to get an error.")
        }

        assertStrings(t, err.Error(), want)
    })
}
```

Go에서 이러한 시나리오를 다루는 방법은 두번째 인자인 `Error` 타입을 활용하는 것이다.

`Error`는 `.Error()` 메소드를 통해 문자열로 변환될 수 있다. 이 문자열은 assertion에 넘겨주는 대상이다. 또한 `if` 조건문을 통해 `assertStrings`를 보호함으로써 `error`가 `nil`일 때 `.Error()`를 호출하지 않게끔 보장한다.

## 테스트 실행해보기

위 코드는 컴파일 되지 않다:

```
./dictionary_test.go:18:10: assignment mismatch: 2 variables but 1 values
```

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

```go
func (d Dictionary) Search(word string) (string, error) {
    return d[word], nil
}
```

이번에 작성한 테스트는 보다 명확한 에러 메시지와 함께 실패할 것이다.

`dictionary_test.go:22: expected to get an error.`

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", errors.New("could not find the word you were looking for")
    }

    return definition, nil
}
```

테스트를 통과하기 위해서, 맵 탐색의 흥미로운 특성을 사용했다. 맵은 두 개의 값을 반환한다. 두번째 값은 boolean으로 키를 찾는데 성공했는지를 가리킨다.

이 성질을 이용해서 단어가 존재하지 않는 것과 단어에 정의가 없는 것을 구분할 수 있다.

## 리팩터링 하기

```go
var ErrNotFound = errors.New("could not find the word you were looking for")

func (d Dictionary) Search(word string) (string, error) {
    definition, ok := d[word]
    if !ok {
        return "", ErrNotFound
    }

    return definition, nil
}
```

`Search` 함수의 매직 에러를 별개의 변수로 뽑아냄으로써 이 에러를 제거할 수 있다. 이것은 더 나은 테스트를 만들 수 있도록 한다.

```go
t.Run("unknown word", func(t *testing.T) {
    _, got := dictionary.Search("unknown")

    assertError(t, got, ErrNotFound)
})
}

func assertError(t testing.TB, got, want error) {
    t.Helper()

    if got != want {
        t.Errorf("got error %q want %q", got, want)
    }
}
```

새 헬퍼를 만든 덕에 테스트가 더 간결해질 수 있었다. `ErrNotFound` 변수를 사용함으로써 에러 문자열이 나중에 바뀌더라도 테스트가 실패하지 않게했다.

## 테스트부터 작성하기

사전을 검색하는 훌륭한 방법을 갖추었다. 그러나 우리의 사전에 새 단어를 추가하는 방법이 없다.

```go
func TestAdd(t *testing.T) {
    dictionary := Dictionary{}
    dictionary.Add("test", "this is just a test")

    want := "this is just a test"
    got, err := dictionary.Search("test")
    if err != nil {
        t.Fatal("should find added word:", err)
    }

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

이 테스트는 `Search` 함수를 활용하여 사전 검사를 보다 쉽게했다.

## 테스트 실행해보기

In `dictionary.go`

```go
func (d Dictionary) Add(word, definition string) {
}
```

테스트는 이제 실패할 것이다

```
dictionary_test.go:31: should find added word: could not find the word you were looking for
```

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func (d Dictionary) Add(word, definition string) {
    d[word] = definition
}
```

맵에 추가하는 것은 배열과 유사하다. 키를 명시하고 값을 같게하면 된다.

### 포인터, 복사, 그 외

맵의 흥미로운 특성은 그것의 주소를 전달(예컨데 `&myMap`)하지 않고서도 수정할 수 있다는 것이다.

"레퍼런스 타입"처럼 느껴질 수 있는데, [Dave Cheney가 설명한 바에 따르면](https://dave.cheney.net/2017/04/30/if-a-map-isnt-a-reference-variable-what-is-it) 그렇지 않다.

맵에 함수/메소드를 전달하게되면 실제로 맵을 복사하지만 단지 포인터 부분만 해당한다. 데이터를 갖고 있는 하부 자료 구조는 복사하지 않는다.

맵에 관해 유의할 점은 `nil` 값이 가능하다는 점이다. 읽기 작업을 수행할 때 `nil` 맵은 빈 맵과 동일하게 동작하지만 `nil` 맵에 쓰기 작업을 시도한다면 이는 런타임 패닉을 일으키는 원인이 된다. 맵에 관해서는 [여기]서 더 알아볼 수 있다.

따라서, 절대로 빈 맵을 초기화 해서는 안 된다:

```go
var m map[string]string
```

대신에, 위에서 해본 것처럼 빈 맵을 초기화 할 수 있는데 아니면 맵을 새로 생성하는 `make` 키워드를 사용할 수 있다:

```go
var dictionary = map[string]string{}

// OR

var dictionary = make(map[string]string)
```

두 방법은 빈 `hash map`을 생성하고 `dictionary`가 이를 가리키게 한다. 이것은 절대로 런타임 패닉을 발생하지 않도록 보장하는 방법이다.

## 리팩터링 하기

구현에 리팩터링할 게 그리 많지 않지만 테스트는 보다 간결하게 만들 수 있다.

```go
func TestAdd(t *testing.T) {
    dictionary := Dictionary{}
    word := "test"
    definition := "this is just a test"

    dictionary.Add(word, definition)

    assertDefinition(t, dictionary, word, definition)
}

func assertDefinition(t testing.TB, dictionary Dictionary, word, definition string) {
    t.Helper()

    got, err := dictionary.Search(word)
    if err != nil {
        t.Fatal("should find added word:", err)
    }

    if definition != got {
        t.Errorf("got %q want %q", got, definition)
    }
}
```

단어와 정의를 위한 변수를 만들었고, 정의를 검사하는 로직을 별도의 헬퍼 함수로 뺴내었다.

`Add` 함수는 괜찮아보인다. 단지, 추가하고자 하는 값이 이미 존재하는 경우를 고려하지 않았다.

맵은 값이 이미 존재할 경우에 에러를 발생하지 않는다. 대신에, 프로그램은 계속 돌아가며 새로 입력한 값으로 덮어씌워진다. 실제로 편리한 점이긴 하지만 함수 이름이 덜 정밀해지는 지점이기도 하다. `Add` 함수는 이미 존재하는 값을 수정해서는 안 된다. 사전에 새 단어만 추가해야한다.

## 테스트부터 작성하기

```go
func TestAdd(t *testing.T) {
    t.Run("new word", func(t *testing.T) {
        dictionary := Dictionary{}
        word := "test"
        definition := "this is just a test"

        err := dictionary.Add(word, definition)

        assertError(t, err, nil)
        assertDefinition(t, dictionary, word, definition)
    })

    t.Run("existing word", func(t *testing.T) {
        word := "test"
        definition := "this is just a test"
        dictionary := Dictionary{word: definition}
        err := dictionary.Add(word, "new test")

        assertError(t, err, ErrWordExists)
        assertDefinition(t, dictionary, word, definition)
    })
}
...
func assertError(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}
```

이 테스트를 위해 `Add` 함수가 에러를 반환하도록 수정하였는데 이는 새 에러 값인 `ErrWordExists`를 검증한다. 이전 테스트를 수정하여 `nil` 에러를 검사하도록 했고 `assertError` 함수도 마찬가지다.

## 테스트 실행해보기

`Add` 함수에서 값을 반환하도록 하지 않게 만들었기 때문에 컴파일러는 실패할 것이다.

```
./dictionary_test.go:30:13: dictionary.Add(word, definition) used as value
./dictionary_test.go:41:13: dictionary.Add(word, "new test") used as value
```

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

`dictionary.go` 파일에서

```go
var (
    ErrNotFound   = errors.New("could not find the word you were looking for")
    ErrWordExists = errors.New("cannot add word because it already exists")
)

func (d Dictionary) Add(word, definition string) error {
    d[word] = definition
    return nil
}
```

새로운 에러를 두개 추가했다. 여전히 값을 변경하고 있으며, `nil` 에러를 반환한다.

```
dictionary_test.go:43: got error '%!q(<nil>)' want 'cannot add word because it already exists'
dictionary_test.go:44: got 'new test' want 'this is just a test'
```

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func (d Dictionary) Add(word, definition string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        d[word] = definition
    case nil:
        return ErrWordExists
    default:
        return err
    }

    return nil
}
```

이제 `switch` 구문을 사용해서 에러를 매칭해보겠다. 이런 식으로 `switch` 구문을 활용하면 추가적인 안전망을 제공하는데, `Search` 함수가 `ErrNotFound`가 아닌 에러를 반환하는 경우가 이에 해당한다.

## 리팩터링 하기

리팩터링할 게 그리 많지 않다. 그러나 에러 활용도가 커져감에 따라 약간의 수정을 해보겠다.

```go
const (
    ErrNotFound   = DictionaryErr("could not find the word you were looking for")
    ErrWordExists = DictionaryErr("cannot add word because it already exists")
)

type DictionaryErr string

func (e DictionaryErr) Error() string {
    return string(e)
}
```

에러를 상수로 만들었다. 이는 `error` 인터페이스를 구현하는 우리만의 `DictionaryErr` 타입을 만드는데 필요하다. [Dave Cheney의 훌륭한 글]에서 자세한 내용을 읽을 수 있다.

다음으로, 단어의 정의를 `Update`하는 함수를 만들어보자.

## 테스트부터 작성하기

```go
func TestUpdate(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{word: definition}
    newDefinition := "new definition"

    dictionary.Update(word, newDefinition)

    assertDefinition(t, dictionary, word, newDefinition)
}
```

다음에 구현할 내용에서 `Update` 함수는 `Add` 함수와 매우 밀접하게 관련이 있다.

## 테스트 실행해보기

```
./dictionary_test.go:53:2: dictionary.Update undefined (type Dictionary has no field or method Update)
```

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

우리는 이와같은 에러를 어떻게 처리해야할 지 이미 알고 있다. 함수를 정의해야한다.

```go
func (d Dictionary) Update(word, definition string) {}
```

이 자리에서, 단어의 정의를 변경해야할 필요가 눈에 보이게 됐다.

```
dictionary_test.go:55: got 'this is just a test' want 'new definition'
```

## 테스트를 통과하는 최소한의 코드 작성하기

`Add` 함수에서의 문제를 고치면서 이런 경우에 무엇을 해야할지 본 적이 있다. 그러모르 `Add` 함수와 엄청 유사한 무언가를 구현해보자.

```go
func (d Dictionary) Update(word, definition string) {
    d[word] = definition
}
```

간단한 변경인지라 리팩터링해야할 것이 없다. 그러나 `Add` 함수와 같은 문제가 있다. 만약 새로운 단어를 전달한다면, `Update`는 사전에 이를 추가한다.

## 테스트부터 작성하기

```go
t.Run("existing word", func(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    newDefinition := "new definition"
    dictionary := Dictionary{word: definition}

    err := dictionary.Update(word, newDefinition)

    assertError(t, err, nil)
    assertDefinition(t, dictionary, word, newDefinition)
})

t.Run("new word", func(t *testing.T) {
    word := "test"
    definition := "this is just a test"
    dictionary := Dictionary{}

    err := dictionary.Update(word, definition)

    assertError(t, err, ErrWordDoesNotExist)
})
```

단어가 존재하지 않는 경우에 관한 또다른 에러 타입을 추가했다. 또한 `Update` 함수를 수정하여 `error` 값을 반환하게 하였다.

## 테스트 실행해보기

```
./dictionary_test.go:53:16: dictionary.Update(word, "new test") used as value
./dictionary_test.go:64:16: dictionary.Update(word, definition) used as value
./dictionary_test.go:66:23: undefined: ErrWordDoesNotExist
```

이번에는 세개의 에러가 나왔는데, 우리는 어떻게 처리해야할 지 알고있다.

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

```go
const (
    ErrNotFound         = DictionaryErr("could not find the word you were looking for")
    ErrWordExists       = DictionaryErr("cannot add word because it already exists")
    ErrWordDoesNotExist = DictionaryErr("cannot update word because it does not exist")
)

func (d Dictionary) Update(word, definition string) error {
    d[word] = definition
    return nil
}
```

우리만의 에러타입을 추가했으며 `nil` 에러를 리턴하게 했다.

이러한 변화들로 이제 매우 분명한 에러 메시지를 받았다:

```
dictionary_test.go:66: got error '%!q(<nil>)' want 'cannot update word because it does not exist'
```

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func (d Dictionary) Update(word, definition string) error {
    _, err := d.Search(word)

    switch err {
    case ErrNotFound:
        return ErrWordDoesNotExist
    case nil:
        d[word] = definition
    default:
        return err
    }

    return nil
}
```

이 함수는 `Add` 함수와 거의 동일해보인다. `dictionery`를 업데이트할 때와 에러를 리턴할 때를 뺴면 말이다.

### 업데이트 함수를 위한 새로운 에러 타입을 선언하는 것에 관한 note

`ErrNotFound`를 재사용하고 새로운 에러 타입을 추가하지 않을 수도 있다. 그러나, 업데이트에 실패했을 때 정확한 에러를 받는 것이 종종 더 나을 때가 있다.

구체적인 에러는 무엇이 잘못됐는지 더 많은 정보를 준다. 웹앱에서의 예를 보자:

> You can redirect the user when `ErrNotFound` is encountered, but display an error message when `ErrWordDoesNotExist` is encountered.

다음으로, 사전에서 단어를 `Delete`하는 함수를 만들어보자.

## 테스트부터 작성하기

```go
func TestDelete(t *testing.T) {
    word := "test"
    dictionary := Dictionary{word: "test definition"}

    dictionary.Delete(word)

    _, err := dictionary.Search(word)
    if err != ErrNotFound {
        t.Errorf("Expected %q to be deleted", word)
    }
}
```

테스트는 단어와 함께 `Dictionary`를 생성하고선 단어가 지워졌는지 확인한다.

## 테스트 실행해보기

`go test`를 실행하면 다음 메시지가 나온다:

```
./dictionary_test.go:74:6: dictionary.Delete undefined (type Dictionary has no field or method Delete)
```

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

```go
func (d Dictionary) Delete(word string) {

}
```

이 함수를 추가한 이후에, 테스트는 단어를 삭제하지 않았다고 알려줄 것이다.

```
dictionary_test.go:78: Expected 'test' to be deleted
```

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func (d Dictionary) Delete(word string) {
    delete(d, word)
}
```

Go는 맵에 사용가능한 `delete`라는 내장 함수가 있다. 이 함수는 두개의 인자를 받는다. 첫번째 인자는 맵이고 두번째 인자는 삭제할 키다.

`delete` 함수는 아무 것도 반환하지 않으며, 우리의 `Delete` 함수도 같은 방식에 기초할 것이다. 존재하지 않는 값을 삭제하는 일은 아무런 영향이 없다. `Update`와 `Add` 함수와는 달리 API를 에러를 포함해서 복잡하게 할 필요가 없다.

## 정리

이번 섹션에서 많은 걸 다루었다. 우리의 사전을 위해 CRUD(Create, Read, Update 그리고 Delete) API 전부를 만들었다. 이 과정을 통해 다음과 같은 방법을 배웠다:

* 맵 만들기
* 맵에서 항목 검색하기
* 맵에 항목 추가하기
* 맵에 항목 갱신하기
* 맵에서 항목 제거하기
* 에러에 관해 더 배웠다
	* 상수인 에러를 만드는 방법
	* 에러 래퍼(error wrappers) 작성하기

[언어 명세]: https://golang.org/ref/spec#Comparison_operators
[여기]: https://blog.golang.org/go-maps-in-action
[Dave Cheney의 훌륭한 글]: https://dave.cheney.net/2016/04/07/constant-errors
