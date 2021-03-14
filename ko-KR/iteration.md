# 반복 

**[이 장의 모든 코드는 여기에서 찾을 수 있습니다](https://github.com/quii/learn-go-with-tests/tree/main/for)**

Go에서 반복적인 작업을 하기 위해서는 `for`가 필요합니다. Go에는 `while`, `do`, `until` 같은 키워드가 없고 오직 `for`만 사용할 수 있습니다. 그건 좋은 일입니다.

문자를 5번 반복하는 함수를 위한 테스트를 작성해 보겠습니다.

여기까지는 새로운 게 없으니, 연습 삼아 작성해 보겠습니다.

## 먼저 테스트를 작성하겠습니다

```go
package iteration

import "testing"

func TestRepeat(t *testing.T) {
	repeated := Repeat("a")
	expected := "aaaaa"

	if repeated != expected {
		t.Errorf("expected %q but got %q", expected, repeated)
	}
}
```

## 테스트를 실행해 보세요

`./repeat_test.go:6:14: undefined: Repeat`

## 테스트를 실행하기 위한 약간의 코드를 작성하고 실패하는 테스트 출력을 확인할 수 있습니다

_수련을 계속 하세요!_ 테스트 실패가 적절하게 된 것인지에 대해서 지금 당장 어떤 것도 새로 알려고 할 필요가 없습니다.

컴파일될 수 있게 해서 작성한 테스트가 잘 작성되었는지 확인 하는 것으로 충분합니다.

```go
package iteration

func Repeat(character string) string {
	return ""
}
```

기본적인 문제에 대해서 테스트를 작성할 정도로 Go를 이미 알고 있는거로 괜찮지 않습니까? 즉, 이제 프로덕션 코드를 하고 싶은 대로 작성할 수 있고 원하는 동작이 잘 동작하는지 알 수 있습니다.

`repeat_test.go:10: expected 'aaaaa' but got ''`

## 테스트를 통과할 수 있는 충분한 코드 작성

`for` 문법은 대체로 C와 비슷한 언어들이 따르는 매우 평범한 형태입니다.


```go
func Repeat(character string) string {
	var repeated string
	for i := 0; i < 5; i++ {
		repeated = repeated + character
	}
	return repeated
}
```

C, Java 또는 JavaScript 같은 다른 언어들과 달리 세 개의 컴포넌트를 둘러싼 중괄호가 없고 중괄호 `{ }`가 항상 필요합니다. 행에서 무슨 일이 발생하고 있는지 궁금할 것입니다.

```go
	var repeated string
```

변수를 초기화하고 선언하기 위해서 `:=`를 사용해 왔습니다. 그러나 `:=`는 [두 단계를 간단하게 줄여 줍니다](https://gobyexample.com/variables). 여기에서는 `string` 변수만 선언하고 있습니다. 그래서 명시적인 버전입니다. `var`를 사용해서 함수를 선언할 수도 있다는 걸 나중에 보게 될 것입니다.

테스트를 실행하면 통과할 것입니다.

for 반복의 추가적인 형태는 [여기](https://gobyexample.com/for)에서 설명하고 있습니다.

## 코드개선

이제 리팩토링을 하고 다른 구조의 `+=` 할당 연산자를 도입할 차례입니다.

```go
const repeatCount = 5

func Repeat(character string) string {
	var repeated string
    for i := 0; i < repeatCount; i++ {
        repeated += character
    }
    return repeated
}
```

`+=`은 _"추가와 할당 연산자"_ 라고 불리고, 오른쪽 피연산자를 왼쪽의 피연산자에 더하고 결과를 왼쪽 피연산자에 할당합니다. 정수 형태의 타입들에서 동작을 합니다.

### 성능측정

Go에서 [benchmarks](https://golang.org/pkg/testing/#hdr-Benchmarks) 작성은 언어의 또 다른 1급 기능이고 테스트 작성과 매우 비슷합니다.

```go
func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a")
	}
}
```

테스트와 매우 비슷한 코드를 볼 수 있습니다.

`testing.B`는 애매하게 이름 지어진 `b.N`에 접근할 수 있게 합니다.

성능측정 코드가 실행되면 `b.N` 번 실행되고 얼마나 오래 걸렸는지 측정합니다.

코드가 실행된 횟수는 문제가 되지 않습니다. 프레임워크는 괜찮은 결과를 구하기 위해 "좋은" 값을 정합니다.

성능측정을 하기 위해서 `go test -bench=.`를 합니다 (혹시 Windows Powershell이라면 `go test -bench="."` 합니다)

```text
goos: darwin
goarch: amd64
pkg: github.com/quii/learn-go-with-tests/for/v4
10000000           136 ns/op
PASS
```

`136 ns/op`의 의미는 함수가 \(내 컴퓨터에서\) 실행하는데 평균 136 나노초가 걸린다는 것입니다. 제법 괜찮은 겁니다! 테스트를 위해서 10000000번을 실행했습니다.

_NOTE_ 벤치마크는 기본적으로 순차적으로 실행됩니다.

## 연습문제

* 호출자에서 문자가 반복되는 횟수를 지정할 수 있도록 테스트를 변경하고 코드도 수정합니다
* 함수를 문서화하기 위하여 `ExampleRepeat`를 작성합니다
* [strings](https://golang.org/pkg/strings) 패키지를 찾아보세요. 쓸모 있다고 생각되는 함수를 찾아보고 여기에서 한 것 같이 테스트를 작성해서 실험해보세요. 표준 라이브러리를 배우는데 시간을 쏟다 보면 시간이 지난 다음 진짜로 보상을 받게 될 것입니다.

## 마무리

* 더 많은 TDD 연습
* `for`에 대한 배움
* 벤치마크 작성 방법을 배움
