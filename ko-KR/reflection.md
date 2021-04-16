# Reflection

**[이 챕터의 모든 코드는 여기에서 확인할 수 있다.](https://github.com/MiryangJung/learn-go-with-tests-ko/tree/main/select)**

[From Twitter](https://twitter.com/peterbourgon/status/1011403901419937792?s=09)

> golang 챌린지: 구조체 `x`를 받고 내부에서 찾을 수 있는 문자열 필드를 위한 `fn`을 호출하는 함수 `walk(x interface{}, fn func(string))`을 작성하라. 난이도: 재귀적

우리는 위의 문제를 해결하기 위해 _reflaction_을 사용할 것이다.

> 컴퓨팅에서 Reflection이란 자신의 구조, 특히 타입을 통해서 검토할 수있는 프로그램의 능력을 말한다. 이것은 메타프로그램의 한 형태이다. 이것은 또한 혼란의 큰 원인이 되기도 한다.

From [The Go Blog: Reflection](https://blog.golang.org/laws-of-reflection)

## `interface` 란 무엇인가?

우리는 Go에서 함수가 `string`, `int` 그리고 우리의 자료형인 `BackAccount`와 같이 알려진 자료형으로 동작한다는 측면에서 type-safety의 편리성을 느껴왔다.

이것은 우리가 쉽게 문서화 할 수 있다는 것과 만약 함수에 잘못된 자료형을 전달하는 경우, 컴파일러가 이를 알아 낼 것임을 의미한다.

컴파일 시 자료형에 대해 알지 못하는 함수를 작성하려는 상황을 접할 수 있다.

Go는 이 문제를 해결하기 위해 _모든_ 자료형이라 생각할 수 있는 `interface{}`라는 자료형을 제공한다.

따라서, `walk(x interface{}, fn func(string))`은 `x`로 어떠한 값도 받을 수 있다.

### 그렇다면 모든 것에 `interface`를 사용하고 정말 유연한  함수를 갖는 건 어떨까?

- `interface`를 사용하는 함수의 사용자는 type-safery를 잃게 된다. 만약 `string` 형인 `Foo.bar`를 함수에 전달하도록 의도했지만 `int`형의 `Foo.baz`가 전달됐다면? 컴파일러는 그 실수를 알려줄 수 없을 것이다. 또한 함수에 _무엇_이 잘되어야 하는지도 알 수 없다. 예시로 함수가 `UserService`를 수용한다는 걸 아는 것은 매우 도움이 된다.
- 함수의 작성자로서, 전달될 _어떠한 것_에 대해 검사할 수 있어야 하며 그 자료형은 무엇인지, 그것으로 무엇을 할 수 있는지를 알아야 한다. 이것을 위해 _reflection_을 이용한다. 이는 상당히 익숙치 않고 읽기 어려울 수 있으며, 일반적으로 성능이 저하된다(런타임에 검사를 해야함).

간략히 말해, 정말 필요할 때에 refection을 사용한다.

만약 다형성을 가진 함수(polymorphic functions)를 원한다면, 인터페이스(`interface`가 아님)를 중심으로 설계할 수 있는지 고려한다. 그러면 사용자는 그 함수가 동작하는데 필요한 방법을 구현 할 때 여러 자료형을 통해 함수를 사용할 수 있다. 

우리의 함수는 다른 많은 것들과 함게 동작해야 할 것이다. 항상 그랬듯, 우리는 우리가 지원하고자 하는 새로운 것에 대한 테스트를 작성하고 끝날 때까지 리팩토링하는 반복적인 접근법을 취할 것이다.

## 테스트부터 작성하기 

우리는 내부에 string 필드(`x`)를 갖는 구조체와 함께 함수를 호출할 것이다. 그러면 전달된 함수(`fn`)에서 그것이 호출되는 지 확인할 수 있다.

```go
func TestWalk(t *testing.T) {

    expected := "Chris"
    var got []string

    x := struct {
        Name string
    }{expected}

    walk(x, func(input string) {
        got = append(got, input)
    })

    if len(got) != 1 {
        t.Errorf("wrong number of function calls, got %d want %d", len(got), 1)
    }
}
```

- 우리는 `walk`를 통해 `fn`에 들어오는 문자열을 담는 문자열 슬라이스를 저장하고자 한다. 이전 장에서는 함수/메소드의 호출부에 전용 자료형을 만들었지만, 이 경우에는 단지 `got`에 접근하는 익명 함수 `fn`을 전달한다.
- 우리는 가장 단순한 방법을 위해 string 자료형인 `Name`을 갖는 익명 `구조체`를 사용한다.
- 마지막으로, `x`와 함께 `walk`를 호출하고 `got`의 길이를 확인한다. 그리고 우리가 아주 기본적인 일을 하게 될 때, assertions에 대해 조금 더 자세히 알아본다.

## 테스트 실행해보기

```
./reflection_test.go:21:2: undefined: walk
```

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

`walk`에 대한 정의가 필요하다.

```go
func walk(x interface{}, fn func(input string)) {

}
```

테스트를 다시 수행한다.

```
=== RUN   TestWalk
--- FAIL: TestWalk (0.00s)
    reflection_test.go:19: wrong number of function calls, got 0 want 1
FAIL
```

## 테스트를 통과하는 최소한의 코드 작성하기

테스트 통과를 위해 아무 문자열을 통해 호출 할 수 있다.

```go
func walk(x interface{}, fn func(input string)) {
    fn("I still can't believe South Korea beat Germany 2-0 to put them last in their group")
}
```

이제 테스트는 통과할 것이다. 이제 필요한 다음 일은 `fn`이 어떤 것과 호출될 것인지 조금 더 정확하게 선언하는 것이다.

## 테스트부터 작성해보기

`fn`에 전달된 문자열이 올바른지 확인하기 위해 다음의 코드를 추가한다.

```go
if got[0] != expected {
    t.Errorf("got %q, want %q", got[0], expected)
}
```

## 테스트 실행해보기

```
=== RUN   TestWalk
--- FAIL: TestWalk (0.00s)
    reflection_test.go:23: got 'I still can't believe South Korea beat Germany 2-0 to put them last in their group', want 'Chris'
FAIL
```

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)
    field := val.Field(0)
    fn(field.String())
}
```

이 코드는 매우 위험하고 단순하지만, 우리가 "빨간색"(테스트 실패)에 있을 때 우리의 목표는 가능한 최소한의 코드를 작성하는 것임을 기억한다. 그런 다음 우리의 우려를 해결하기 위해 더 많은 테스트를 작성한다.

우리는 `x`와 그 속성을 알아보기 위해 reflection을 이용한다.

[reflect 패키지](https://godoc.org/reflect)는 주어진 변수의 `값`을 전달하는 `ValueOf`함수를 갖는다. 이것은 우리에게 값을 알아볼 방법을 제공하고 우리가 그 다음 줄에 사용한 것처럼 그 값의 필드까지도 포함한다.

그런 다음 전달된 값에 대한 매우 낙관적인 가정을 한다.

- 첫번째 필드를 찾아보고, panic을 일으킬 필드는 없을지도 모른다.
- 그런 다음, 문자열을 기본값으로 전달하는 `String()`을 호출하고 만약 해당 필드가 문자열이 아닌 다른 값이라면 문제가 될 것임을 안다.

## 리팩터링 하기

우리의 코드가 단순한 케이스에서는 통과하지만 많은 단점을 가지고 있다는 것을 안다.

우리는 여러 다른 값을 전달하는 테스트를 작성하고 `fn`과 함께 호출되는 문자열 집합을 확인할 것이다.

우리는 새로운 시나리오를 더 쉽게 테스트하기 위해 테스트를 표 기반 테스트로 리팩토링해야 한다.

```go
func TestWalk(t *testing.T) {

    cases := []struct{
        Name string
        Input interface{}
        ExpectedCalls []string
    } {
        {
            "Struct with one string field",
            struct {
                Name string
            }{ "Chris"},
            []string{"Chris"},
        },
    }

    for _, test := range cases {
        t.Run(test.Name, func(t *testing.T) {
            var got []string
            walk(test.Input, func(input string) {
                got = append(got, input)
            })

            if !reflect.DeepEqual(got, test.ExpectedCalls) {
                t.Errorf("got %v, want %v", got, test.ExpectedCalls)
            }
        })
    }
}
```

이제 우리는 하나 이상의 문자열 필드를 가질 때 어떤 일이 일어나는 지에 대한 시나리오를 쉽게 추가할 수 있다.

## 테스트부터 작성하기

다음의 시나리오를 `cases`에 추가한다.

```go
{
    "Struct with two string fields",
    struct {
        Name string
        City string
    }{"Chris", "London"},
    []string{"Chris", "London"},
}
```

## 테스트 실행해보기

```
=== RUN   TestWalk/Struct_with_two_string_fields
    --- FAIL: TestWalk/Struct_with_two_string_fields (0.00s)
        reflection_test.go:40: got [Chris], want [Chris London]
```

## 테스트를 통과하는 최소한의 코드 작성하기 

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i:=0; i<val.NumField(); i++ {
        field := val.Field(i)
        fn(field.String())
    }
}
```

`val`은 값 내부 필드의 수를 반환하는 `NumField`메소드를 갖는다. 이것을 통해 필드를 순회하여 테스트에 통과하는 `fn`을 호출할 수 있다.

## 리팩터링 하기

코드를 개선할 수 있는 분명한 요인이 있는 것 같지는 않으니 계속 진행한다.

`walk`의 다음 단점은 모든 필드를 `string`으로 간주하는 것이다. 다음 시나리오에 대한 테스트를 작성해본다.

## 테스트부터 작성하기

다음의 케이스를 추가한다.

```go
{
    "Struct with non string field",
    struct {
        Name string
        Age  int
    }{"Chris", 33},
    []string{"Chris"},
},
```

## 테스트 실행해보기

```
=== RUN   TestWalk/Struct_with_non_string_field
--- FAIL: TestWalk/Struct_with_non_string_field (0.00s)
    reflection_test.go:46: got [Chris <int Value>], want [Chris]
```

## 테스트를 통과하는 최소한의 코드 작성하기 

이제 필드의 자료형이 `string`인지 확인할 필요가 있다.

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        if field.Kind() == reflect.String {
            fn(field.String())
        }
    }
}
```

[Kind](https://godoc.org/reflect#Kind)를 통해 확인할 수 있다.

## 리팩터링 하기

지금까지는 코드가 충분히 적당한 것으로 보인다.

다음 시나리오는 만약 `struct`가 "flat" 하지 않은 경우이다. 다른 말로, 만약 `struct`가 nested 필드를 갖는다면 어떻게 되는지이다.

## 테스트부터 작성하기

우리는 자료형을 임시방편으로 선언하기 위해 익명 구조체 구문을 사용해왔고 다음과 같이 계속할 수 있다.

```go
{
    "Nested fields",
    struct {
        Name string
        Profile struct {
            Age  int
            City string
        }
    }{"Chris", struct {
        Age  int
        City string
    }{33, "London"}},
    []string{"Chris", "London"},
},
```

하지만 우리가 내부 익명 구조체 구문을 사용할 때 약간의 혼란이 있을 수 있다. [여기에 해당 구문을 더 훌륭하게 만들 제안이 있다.](https://github.com/golang/go/issues/12854)

이제 이 시나리오를 위한 알려진 타입을 만들고 테스트에서 참조하는 방식으로 변경해보자. 테스트를 위한 우리의 코드가 테스트 외부에 있다는 점에서 약간의 기만이 있지만, 독자들은 초기화를 통해 `구조체`의 구조를 추론할 수 있어야 한다.

다음의 타입 선언을 테스트 파일 어딘가에 추가해보자.

```go
type Person struct {
    Name    string
    Profile Profile
}

type Profile struct {
    Age  int
    City string
}
```

이제 우리는 이전보다 훨씬 깔끔하게 케이스를 추가할 수 있다.

```go
{
    "Nested fields",
    Person{
        "Chris",
        Profile{33, "London"},
    },
    []string{"Chris", "London"},
},
```

## 테스트 실행해보기

```
=== RUN   TestWalk/Nested_fields
    --- FAIL: TestWalk/Nested_fields (0.00s)
        reflection_test.go:54: got [Chris], want [Chris London]
```

문제는 타입의 계층 구조 중 첫번째 수준의 필드에서만 반복이된다는 것이다.

## 테스트를 통과하는 최소한의 코드 작성하기 

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        if field.Kind() == reflect.String {
            fn(field.String())
        }

        if field.Kind() == reflect.Struct {
            walk(field.Interface(), fn)
        }
    }
}
```

해결 방법은 꽤 단순하다. `Kind`를 통해 다시 한번 검사하고 만약 그것이 `구조체` 라면 우리는 단지 내부에서 `walk`를 다시 호출하면 된다.

## 리팩터링

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}
```

동일한 값에 대한 한 번 이상의 비교를 해야할 때, _일반적으로_ `switch`구문으로 변경하는 것이 가독성과 확장성을 높일 수 있다.

만약 통과되는 구조체의 값이 포인터라면 어떻게 할까?

## 테스트부터 작성해보기

아래 케이스를 추가한다.

```go
{
    "Pointers to things",
    &Person{
        "Chris",
        Profile{33, "London"},
    },
    []string{"Chris", "London"},
},
```

## 테스트 실행해보기

```
=== RUN   TestWalk/Pointers_to_things
panic: reflect: call of reflect.Value.NumField on ptr Value [recovered]
    panic: reflect: call of reflect.Value.NumField on ptr Value
```

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func walk(x interface{}, fn func(input string)) {
    val := reflect.ValueOf(x)

    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}
```

포인터인 `값`에서는 `NumField`를 사용할 수 없다. 우리는 그전에 드러나지 않은 값을 추출할 필요가 있고 그것은 `Elem()`을 통해 할 수 있다.

## 리팩터링

이제 주어진 함수로 주어진 `interface{}`로부터 `refect.Value`를 추출하는 기능을 encapsulate 해보자.

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}

func getValue(x interface{}) reflect.Value {
    val := reflect.ValueOf(x)

    if val.Kind() == reflect.Ptr {
        val = val.Elem()
    }

    return val
}
```

실제로 _더 많은_ 코드를 추가했지만, 이러한 추상화 수준이 옳다고 생각한다.

- `x`의 `reflect.Value`를 얻고 검사할 수 있지만, 나는 그것이 어떻게 되는지 신경쓰지 않아도 된다.
- 필드를 반복하며 그 타입에 따라 필요한 무엇이든지 한다.

다음으로는, 슬라이스를 보완해야한다.

## 테스트부터 작성하기

```go
{
    "Slices",
    []Profile {
        {33, "London"},
        {34, "Reykjavík"},
    },
    []string{"London", "Reykjavík"},
},
```

## 테스트 실행해보기

```
=== RUN   TestWalk/Slices
panic: reflect: call of reflect.Value.NumField on slice Value [recovered]
    panic: reflect: call of reflect.Value.NumField on slice Value
```

## 테스트를 실행할 최소한의 코드를 작성하고 테스트 실패 결과를 확인하기

이것은 이전의 포인터 시나리오와 비슷하다. 우리는 `reflect.Value`에서 `NumField`를 호출 하려하지만, 그것은 구조체가 아니기 때문에 값이 없다.

## 테스트를 통과하는 최소한의 코드 작성하기 

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    if val.Kind() == reflect.Slice {
        for i:=0; i< val.Len(); i++ {
            walk(val.Index(i).Interface(), fn)
        }
        return
    }

    for i := 0; i < val.NumField(); i++ {
        field := val.Field(i)

        switch field.Kind() {
        case reflect.String:
            fn(field.String())
        case reflect.Struct:
            walk(field.Interface(), fn)
        }
    }
}
```

## 리팩터링

이 코드는 작동하지만 조금 지저분하다. 그래도 작동하는 코드가 있으니 편안하게 우리가 좋아하는 방식으로 손볼 수 있다.

조금 추상적으로 생각한다면, 우리는 두 경우 모두에서 `walk`를 호출 하고 싶을 것이다.

- 구조체 내부의 각각의 필드
- 슬라이스 내부의 각각의 `무언가`

현재 우리의 코드는 그렇게 동작하지만, 제대로 reflect하고 있지는 않다. 그래서 그것이 슬라이스(코드의 남은 실행을 멈출 수 있는 `return`이 있는)인지를 처음에 검사하고 슬라이스가 아니라면 구조체라고 가정한다.

이제 다시 코드를 다시 수정해서 타입을 먼저 확인하고 작업을 진행해본다.

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    switch val.Kind() {
    case reflect.Struct:
        for i:=0; i<val.NumField(); i++ {
            walk(val.Field(i).Interface(), fn)
        }
    case reflect.Slice:
        for i:=0; i<val.Len(); i++ {
            walk(val.Index(i).Interface(), fn)
        }
    case reflect.String:
        fn(val.String())
    }
}
```

훨씬 좋아보인다. 만약 구조체 혹은 슬라이스라면 우리는 각각 `walk`를 호출하며 그 값을 순회한다. 그렇지 않고 만약 `relect.String`이라면 `fn`을 호출하면 된다.

아직도 더 개선할 부분이 있어보인다. 필드/값을 순회하는 연산을 반복적으로 하고 `walk`함수를 호출하는데, 이것은 개념적으로 모두 같은 부분이다.

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    numberOfValues := 0
    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice:
        numberOfValues = val.Len()
        getField = val.Index
    }

    for i:=0; i< numberOfValues; i++ {
        walk(getField(i).Interface(), fn)
    }
}
```

만약 `값`이 `reflect.String`이라면 평소처럼 그냥 `fn`을 호출한다.

그렇지 않다면, `switch`를 통해 타입에 의존한 두 가지 것을 추출한다.

- 몇 개의 필드가 있는지
- 어떻게 `값`(`필드` 또는 `인덱스`)을 추출할 것인지

이것을 정의하게되면 우리는 `numberOfValues`만큼 순회하며 `getField`함수의 결과와 함께 `walk`를 호출할 수 있다.

이제 배열을 처리하는 일은 간단하다.

## 테스트부터 작성하기 

아래 케이스를 추가한다.

```go
{
    "Arrays",
    [2]Profile {
        {33, "London"},
        {34, "Reykjavík"},
    },
    []string{"London", "Reykjavík"},
},
```

## 테스트 실행해보기

```
=== RUN   TestWalk/Arrays
    --- FAIL: TestWalk/Arrays (0.00s)
        reflection_test.go:78: got [], want [London Reykjavík]
```

## 테스트를 통과하는 최소한의 코드 작성하기 

배열은 슬라이스와 동일하게 처리될 수 있으므로 그냥 콤마와 함께 케이스를 추가한다.

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    numberOfValues := 0
    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice, reflect.Array:
        numberOfValues = val.Len()
        getField = val.Index
    }

    for i:=0; i< numberOfValues; i++ {
        walk(getField(i).Interface(), fn)
    }
}
```

다음 우리가 다룰 타입은 `map`이다.

## 테스트부터 작성하기

```go
{
    "Maps",
    map[string]string{
        "Foo": "Bar",
        "Baz": "Boz",
    },
    []string{"Bar", "Boz"},
},
```

## 테스트 실행해보기

```
=== RUN   TestWalk/Maps
    --- FAIL: TestWalk/Maps (0.00s)
        reflection_test.go:86: got [], want [Bar Boz]
```

## 테스트를 통과하는 최소한의 코드 작성하기 

조금만 추상적으로 생각해보면 `map`은 `struct`와 굉장히 유사하다는 것을 알 수 있다. 단지 컴파일 과정에서 키들을 알 수 없다는 것 뿐이다.

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    numberOfValues := 0
    var getField func(int) reflect.Value

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        numberOfValues = val.NumField()
        getField = val.Field
    case reflect.Slice, reflect.Array:
        numberOfValues = val.Len()
        getField = val.Index
    case reflect.Map:
        for _, key := range val.MapKeys() {
            walk(val.MapIndex(key).Interface(), fn)
        }
    }

    for i:=0; i< numberOfValues; i++ {
        walk(getField(i).Interface(), fn)
    }
}
```

하지만 설계상 인덱스별로 맵의 값을 가지고 올 수 없다. 오직 _키_ 를 통해 가능하므로 애석하게도 추상화을 깨뜨린 것이다.

## 리팩터링 하기

어떤가? 괜찮은 추상화였다고 생각했었는데 이제 우리의 코드는 약간 찌질하게 느껴진다.

_괜찮다!_  리팩터링은 여정이고 때로는 실수를 하기도 한다. TDD의 중요한 점은 이런 것들을 시험해 볼 수 잇는 자유를 준다는 것이다.

테스트를 통해 뒷받침되는 작은 단계를 밟는 것은 결코 돌이킬 수 없는 상황이 아니다. 리팩터링 전의 상태로 되돌려보자.

```go
func walk(x interface{}, fn func(input string)) {
    val := getValue(x)

    walkValue := func(value reflect.Value) {
        walk(value.Interface(), fn)
    }

    switch val.Kind() {
    case reflect.String:
        fn(val.String())
    case reflect.Struct:
        for i := 0; i< val.NumField(); i++ {
            walkValue(val.Field(i))
        }
    case reflect.Slice, reflect.Array:
        for i:= 0; i<val.Len(); i++ {
            walkValue(val.Index(i))
        }
    case reflect.Map:
        for _, key := range val.MapKeys() {
            walkValue(val.MapIndex(key))
        }
    }
}
```

`val`에서 `reflect.Value`를 추출하기 위해 `switch`에서 `walk`를 호출하는(DRY) `walkValue`를 도입했다.

### 마지막 문제

Go에서 맵은 순서를 보장하지 않는다는 걸 기억해라. 따라서 우리는 특정한 순서 내에서 `fn`이 호출되는 것으로 단언했기 때문에 테스트는 때때로 실패할 것이다.

이 문제를 해결하기 위해, 맵에 대한 단언(assertion)을 순서를 신경 쓰지 않는 새로운 테스트로 옮길 필요가 있다.

```go
t.Run("with maps", func(t *testing.T) {
    aMap := map[string]string{
        "Foo": "Bar",
        "Baz": "Boz",
    }

    var got []string
    walk(aMap, func(input string) {
        got = append(got, input)
    })

    assertContains(t, got, "Bar")
    assertContains(t, got, "Boz")
})
```

`assertContains`의 정의는 아래와 같다.

```go
func assertContains(t testing.TB, haystack []string, needle string)  {
    t.Helper()
    contains := false
    for _, x := range haystack {
        if x == needle {
            contains = true
        }
    }
    if !contains {
        t.Errorf("expected %+v to contain %q but it didn't", haystack, needle)
    }
}
```

다음 우리가 다룰 타입은 `chan`이다.

## 테스트부터 작성하기 

```go
t.Run("with channels", func(t *testing.T) {
		aChannel := make(chan Profile)

		go func() {
			aChannel <- Profile{33, "Berlin"}
			aChannel <- Profile{34, "Katowice"}
			close(aChannel)
		}()

		var got []string
		want := []string{"Berlin", "Katowice"}

		walk(aChannel, func(input string) {
			got = append(got, input)
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
```

## 테스트 실행해보기

```
--- FAIL: TestWalk (0.00s)
    --- FAIL: TestWalk/with_channels (0.00s)
        reflection_test.go:115: got [], want [Berlin Katowice]
```

## 테스트를 통과하는 최소한의 코드 작성하기 

채널이 Recv()로 종료되기 전까지 채널을 통해 전달된 모든 값을 순회할 수 있다.

```go
func walk(x interface{}, fn func(input string)) {
	val := getValue(x)

	walkValue := func(value reflect.Value) {
		walk(value.Interface(), fn)
	}

	switch val.Kind() {
	case reflect.String:
		fn(val.String())
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			walkValue(val.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			walkValue(val.Index(i))
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			walkValue(val.MapIndex(key))
		}
	case reflect.Chan:
		for v, ok := val.Recv(); ok; v, ok = val.Recv() {
			walk(v.Interface(), fn)
		}
	}
}
```
다음 우리가 다룰 타입은 `func`이다.

## 테스트부터 작성하기 

```go
t.Run("with function", func(t *testing.T) {
		aFunction := func() (Profile, Profile) {
			return Profile{33, "Berlin"}, Profile{34, "Katowice"}
		}

		var got []string
		want := []string{"Berlin", "Katowice"}

		walk(aFunction, func(input string) {
			got = append(got, input)
		})

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
```

## 테스트 실행해보기

```
--- FAIL: TestWalk (0.00s)
    --- FAIL: TestWalk/with_function (0.00s)
        reflection_test.go:132: got [], want [Berlin Katowice]
```

## 테스트를 통과하는 최소한의 코드 작성하기 

인자값을 갖는 함수는 이 시나리오 상 알맞지 않다고 보인다. 하지만 우리는 임의의 리턴값도 허용해야 한다.

```go
func walk(x interface{}, fn func(input string)) {
	val := getValue(x)

	walkValue := func(value reflect.Value) {
		walk(value.Interface(), fn)
	}

	switch val.Kind() {
	case reflect.String:
		fn(val.String())
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			walkValue(val.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			walkValue(val.Index(i))
		}
	case reflect.Map:
		for _, key := range val.MapKeys() {
			walkValue(val.MapIndex(key))
		}
	case reflect.Chan:
		for v, ok := val.Recv(); ok; v, ok = val.Recv() {
			walk(v.Interface(), fn)
		}
	case reflect.Func:
		valFnResult := val.Call(nil)
		for _, res := range valFnResult {
			walk(res.Interface(), fn)
		}
	}
}
```

## 정리

- `reflect`패키지의 몇가지 개념을 설명했다.
- 임의의 데이터 구조를 살펴보기 위해 재귀를 사용했다.
- 나쁜 리팩토링을 경험했지만 이에 대해 크게 당황하지 않았다. 테스트를 반복적으로 하는 것은 그리 큰 일이 아니다.
- 이 글은 reflection에 작은 관점만을 담고 있다. [Go 블로그에 더 세부적인 내용을 담고 있는 좋은 글들이 있다.](https://blog.golang.org/laws-of-reflection)
- 이제 reflection에 대해 알았으니, 이것을 사용하지 않도록 최선을 다한다.

