# 배열과 슬라이스

**[이 챕터의 모든 코드는 여기에서 확인할 수 있다.](https://github.com/MiryangJung/learn-go-with-tests-ko/tree/main/arrays)**

배열은 같은 타입의 여러 요소들을 한 변수에 특정한 순서대로 저장할 수 있게 해준다.

우리가 배열을 다룰 때에, 배열들에 대해 반복하는 것은 매우 흔한 일이다. 그래서 우리가 [새롭게 배운 `for`](iteration.md)을 사용하여 `Sum`함수를 만들 것이다. `Sum`함수는 숫자 배열들을 가지고 총합을 리턴할 것이다.

우리 `테스트 기반 설계` 스킬을 활용할 것이다.


## 테스트를 먼저 작성하세요   


작업을 할 새로운 폴더를 생성한다. `sum_test.go`라는 파일을 생성한 후 아래 코드를 집어넣는다.
```go
package main

import "testing"

func TestSum(t *testing.T) {

	numbers := [5]int{1, 2, 3, 4, 5}

	got := Sum(numbers)
	want := 15

	if got != want {
		t.Errorf("got %d want %d given, %v", got, want, numbers)
	}
}
```

배열은 변수로 선언될 때 우리가 정한 고정된 크기를 가진다. 우리는 두 가지 방법으로 배열을 선언할 수 있다.
* \[N\]type{value1, value2, ..., valueN} e.g. `numbers := [5]int{1, 2, 3, 4, 5}`
* \[...\]type{value1, value2, ..., valueN} e.g. `numbers := [...]int{1, 2, 3, 4, 5}`

[여기에서 문자열 포멧팅에 대해서 더 읽어볼 수 있다.](https://golang.org/pkg/fmt/)

## 테스트를 실행해 보세요

 `go test`로 실행을 시키면 다음과 같이 컴파일 에러가 날것이다. `./sum_test.go:10:15: undefined: Sum`

 ## 테스트를 실행할 최소한의 코드를 작성하고 실패한 테스트 출력을 확인하세요.

 `sum.go`를 다음과 같이 한다면,

```go
package main

func Sum(numbers [5]int) int {
	return 0
}
```
당신은 테스트의 조금 더 명확한 에러 메시지를 확인할 수 있다.
`sum_test.go:13: got 0 want 15 given, [1 2 3 4 5]`

## 통과할 수 있도록 충분한 코드를 작성하세요

```go
func Sum(numbers [5]int) int {
	sum := 0
	for i := 0; i < 5; i++ {
		sum += numbers[i]
	}
	return sum
}
```
배열의 특정한 인덱스 값을 가져오려고 한다면, `array[index]`라는 syntax를 사용하면 된다. 지금과 같은 경우에는, 우 `for`을 이용하여 5번을 반복하며 배열의 각 항목들을 sum에 더할 것이다.
## 리팩토링 
우리는 `range`를 사용하여 코드를 간추릴 것이다.

```go
func Sum(numbers [5]int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
```
`range`는 배열에 반복적인 처리를 할 수 있게 한다. `range`는 호출될 때마다 인덱스와 값, 총 두 개의 값을 리턴한다. 우리는 여백 식별자인 `_`를 사용하여 인덱스의 값을 무시할 것이다.

### 배열과 그 타입

배열의 흥미로운 특성은, 배열의 크기는 배열의 타입으로 인코딩 된다는 점이다. 따라서 `[5] int`를 요구하는 함수에 `[4] int`가 입력될 경우 컴파일 되지 않을 것이다. 이 이유는 바로 타입이 다르기 때문이다, 이는 `int`를 요구하는 함수에 `string`을 입력한 것과 같다.

당신은 이제 배열이 고정된 크기를 가진다는 점 때문에 복잡하다고 생각하겠지만, 대부분의 경우 당신은 배열을 사용하지 않을 것이다.

Go 언어에서의 슬라이스는 컬렉션의 크기를 같이 인코딩 하지 않아 아무 크기나 가질 수 있다.

다음 요건은 달라지는 크기의 컬렉션의 합을 구하는 것이다.

## 테스트를 먼저 작성하세요

우리는 이제 컬렉션의 크기가 자유로운 `slice`타입을 이용할 것이다. 슬라이스의 syntax는 배열과 매우 유사하다, 그저 선언할 때 크기를 생략하기만 하면 된다.

 `myArray := [3]int{1,2,3}` 대신 `mySlice := []int{1,2,3}` 처럼

 ```go
func TestSum(t *testing.T) {

	t.Run("collection of 5 numbers", func(t *testing.T) {
		numbers := [5]int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("collection of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

}
```

## 테스트를 실행해 보세요

위 코드는 컴파일 되지 않는다.
`./sum_test.go:22:13: cannot use numbers (type []int) as type [5]int in argument to Sum`

 ## 테스트를 실행할 최소한의 코드를 작성하고 실패한 테스트 출력을 확인하세요.
 여기서의 문제는 둘 중 하나이다.
- 기존에 존 제하던 API의 인수인 `Sum`을 배열 대신 슬라이스로 바꾸어 주는 것이다. 만약 우리가 이 방법을 사용한다면 우리는 다른 사람의 하루를 망쳤을 수도 있다. 이는 우리의 다른 테스트들을 컴파일 되지 않도록 하기 때문이다.
- 새로운 함수를 만드는 것이다.

지금과 같은 경우, 아무도 우리의 함수를 사용하지 않을 거기 때문에 두 개의 함수를 유지하는 것보다 1개만 가지도록 할 것이다.
```go
func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
```
당신이 컴파일 하려고 한다면, 이는 아직도 컴파일 되지 않을 것이다. 당신은 배열 이 아닌 슬라이스로 pass 하려고 한다면 첫 번째 테스트를 수정해야 한다.

## 통과할 수 있도록 충분한 코드를 작성하세요
알고 보니 우리는 컴파일러의 문제만 해결했으면 됐다. 테스트는 통과했다.

## 리팩토링 
우리는 이미 `Sum`을 리팩토링을 하였으며, 우리는 그저 배열을 슬라이스로 바꾸기만 했다. 그러므로 이 부분에서는 많이 할 것이 없다. 우리는 리팩토링 단계에서 테스트 코드를 소홀히 하면 안 되고 이 단계에서 우리가 해야 할 것이 있다는 것을 기억해야 한다.
```go
func TestSum(t *testing.T) {

	t.Run("collection of 5 numbers", func(t *testing.T) {
		numbers := []int{1, 2, 3, 4, 5}

		got := Sum(numbers)
		want := 15

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

	t.Run("collection of any size", func(t *testing.T) {
		numbers := []int{1, 2, 3}

		got := Sum(numbers)
		want := 6

		if got != want {
			t.Errorf("got %d want %d given, %v", got, want, numbers)
		}
	})

}
```

당신의 테스트 값에 대하여 의문을 갖는 것은 중요하다. 테스트 많이 한다고 좋은 것이 아니다, 당신에 코드에 자신감을 가지는 것이 더 중요하다. 많은 양의 테스트는 문제를 낳으며 유지할 때 그저 오버헤드를 증가시키기만 한다. 모든 테스트는 비용이다.

지금과 같은 경우, 한 함수에 대하여 2개의 테스트를 가지고 있는 것은 매우 불필요하다. 만약 함수가 한 가지 크기의 슬라이스로 통과를 한다면, 이는 다른 크기의 슬라이스도 높은 확률로 통과할 수 있다는 뜻이다(단, 합당한 범위 내에서).

Go 언어의 내장된 테스팅 도구는 커버리지 도구가 있다, 당신의 다루지 않는 부분의 코드를 확인할 수 있도록 돕는다. 저는 100% 커버리지가 당신의 목표가 아니라는 것을 강조하고 싶다. 이는 그저 당신에게 어느 정도가 커버리지가 되는지 알려주기 위함이다. 당신이 TDD에 관하여 엄격하였다면, 높은 확률로 당신의 커버리지는 100%로 끝날 것이기 때문이다.
아래를 실행한다.

`go test -cover`

당신은 다음을 확인할 수 있다.

```bash
PASS
coverage: 100.0% of statements
```

이제 테스트 하나를 지우고 커버리지를 다시 확인한다.

우리는 이제 잘 테스트된 함수를 가지고 있다, 이제 다음 도전을 하기 전에 우리의 멋진 작업을 커밋한다.

우리는 `SumAll`이라는 여러 슬라이스를 받고 그들의 합을 슬라이스들로 리턴할 새로운 함수가 필요하다.

예시로,

`SumAll([]int{1,2}, []int{0,9})` 는 `[]int{3, 9}`을 리턴할 것이다.

또는

`SumAll([]int{1,1,1})` 는 `[]int{3}`를 리턴할 것이다.

## 테스트를 먼저 작성하세요

```go
func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
```

## 테스트를 실행해 보세요

`./sum_test.go:23:9: undefined: SumAll`

## 테스트를 실행할 최소한의 코드를 작성하고 실패한 테스트 출력을 확인하세요.

우리는 SumAll을 테스트가 요구하는 대로 정의를 해야 한다.

Go 언어에서는 변동 가능한 수의 인수들을 받을 수 있는 [_variadic functions_](https://gobyexample.com/variadic-functions)을 사용할 수 있게 해준다.

```go
func SumAll(numbersToSum ...[]int) (sums []int) {
	return
}
```

컴파일을 시도해 보아라, 하지만 우리의 테스트는 아직 컴파일 되지 않는다!

`./sum_test.go:26:9: invalid operation: got != want (slice can only be compared to nil)`

Go 언어에서는 slice를 다룰 때에 등호를 사용할 수 없다. 당신은 `got`슬라이스와 `want`슬라이스를 반복하여 그들의 값을 비교하는 함수를 만들 수 있다. 하지만 이 방법은 편리하지 않다. 그래서 우리는 `any` 타입이나 비교할 수 있는 [`reflect.DeepEqual`][deepEqual]를 사용할 수 있다.

```go
func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```
\(파일 상단에 `import reflect`가 있어야 `DeepEqual`를 사용할 수 있다.\)

우리가 집고 넘어가야 하는 부분은 `reflect.DeepEqual`는 "type safe"하지 않다는 것이다. 당신이 멍청한 짓을 했어도 코드는 그냥 컴파일 될 것이다. 이것을 확인하고 싶다면 테스트를 잠시 동안 다음과 같이 바꾼다.

```go
func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := "bob"

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```
우리가 위에서 하고자 하는 것은 `슬라이스`와 `문자열`을 비교하려고 하고 있다는 것이다. 이것은 터무니없는 일이지만 테스트는 컴파일 되었다. 따라서 `reflect.DeepEqual`를 사용하여 슬라이스를 \(그리고 다른 문자열들을\) 비교하는 것은 매우 편리하지만 사용할 때에 주의를 요한다.

테스트를 다시 원상복구 시키고 실행한다면, 다음과 같은 출력값을 보게 될 것이다.

`sum_test.go:30: got [] want [3 9]`

## 통과할 수 있도록 충분한 코드를 작성하세요

우리는 `varags`들을 반복적으로 처리하며, 우리의 `Sum`한수를 사용해서 합을 계산을 한 후 우리가 리턴할 슬라이스에 추가해야한다.


```go
func SumAll(numbersToSum ...[]int) []int {
	lengthOfNumbers := len(numbersToSum)
	sums := make([]int, lengthOfNumbers)

	for i, numbers := range numbersToSum {
		sums[i] = Sum(numbers)
	}

	return sums
}
```

우리는 새로 배울 것들이 아주 많다!

슬라이스를 만드는 새로운 방법이 있다. `make`함수는 우리가 작업해야 할 `numbersToSum`의 `len`만큼의 시작 크기를 가지는 슬라이스를 만들 수 있도록 해준다. 당신은 배열에서 `mySlice[N]` 하듯이 슬라이스에서도 인덱스를 사용할 수 있으며, 새로운 값을 `=`를 이용해서 지정해 줄 수 있다.

이제 테스트는 통과할 것이다.

## 리팩토링

앞서 말했듯이 슬라이스도 크기가 있다. 당신이 크기가 2인 슬라이스에서 `mySlice[10] = 1`를 시도한다면 이것은 _런타임_ 에러가 날것이다.

하지만 당신은 슬라이스와 새로운 값을 받고 새로운 값이 들어간 슬라이스를 리턴하는 `append`함수를 이용할 수 있다.

```go
func SumAll(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}

	return sums
}
```
위 실행에서는, 우리는 크기를 너무 걱정할 필요가 없다. 우리는 비어있는 `sums`라는 슬라이스에 `Sum`함수의 결괏값을 vararngs들을 처리할 때마다 `append`를 사용할 것이기 때문이다. 컬렉션의 꼬리 란 제일 첫 번째인 \("머리"\)를 제외한 나머지 것들을 이야기한다.

우리의 다음 요구 사항은 `SumAll` 을 모든 슬라이스의 꼬리들의 합을 계산하는 `SumAllTails`로 바꾸는 것이다.

## 테스트를 먼저 작성하세요

```go
func TestSumAllTails(t *testing.T) {
	got := SumAllTails([]int{1, 2}, []int{0, 9})
	want := []int{2, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

## 테스트를 실행해 보세요

`./sum_test.go:26:9: undefined: SumAllTails`

## 테스트를 실행할 최소한의 코드를 작성하고 실패한 테스트 출력을 확인하세요.

함수를  `SumAllTails` 로 이름을 변경하고 테스트를 다시 실행한다.

`sum_test.go:30: got [3 9] want [2 9]`

## 통과할 수 있도록 충분한 코드를 작성하세요

```go
func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		tail := numbers[1:]
		sums = append(sums, Sum(tail))
	}

	return sums
}
```
슬라이스들도`slice[low:high]`를 통해 잘릴 수 있다. 당신이 `:`의 옆에 값을 누락시킨다면, 이것은 그 값부터 끝까지의 값을 가져온다. 우리와 같은 경우 "1부터 끝까지 가져와"를 `numbers[1:]`과 같이 표현할 수 있다. 당신은 시간을 조금 투자해서 다른 테스트 등을 통하여 슬라이스와 슬라이스 오퍼레이터와 익숙해져도 좋다.

## 리팩토링

이번에는 리팩토링할 것이 거의 없다.

만약 비어있는 슬라이스를 우리의 함수에 넣는다면 어떻게 될 것인가? 비어있는 슬라이스의 꼬리는 무엇인가? 당신이 Go 언어에게 `myEmptySlice[1:]`있는 모든 요소들을 달라고 한다면 어떻게 될 것 같은가?

## 테스트를 먼저 작성하세요

```go
func TestSumAllTails(t *testing.T) {

	t.Run("make the sums of some slices", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	})

}
```

## 테스트를 실행해 보세요

```text
panic: runtime error: slice bounds out of range [recovered]
    panic: runtime error: slice bounds out of range
```
오, 이런! 우리는 테스트가 _컴파일_ 되었는지 유의해야 한다. 위 에러는 런타임 에러이다.

컴파일 타임 에러는 소프트웨어를 만드는 데 있어서 도움을 주기 때문에 우리의 친구이다, 하지만 런타임 에러는 우리의 사용자에게 영향을 주기 때문에 우리의 적이다.

## 통과할 수 있도록 충분한 코드를 작성하세요

```go
func SumAllTails(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		if len(numbers) == 0 {
			sums = append(sums, 0)
		} else {
			tail := numbers[1:]
			sums = append(sums, Sum(tail))
		}
	}

	return sums
}
```

## 리팩토링

우리의 코드는 assertion 부분에서 코드가 반복되는 일이 또 발생하였다. 우리 이 부분을 함수로 추출할 것이다.

```go
func TestSumAllTails(t *testing.T) {

	checkSums := func(t testing.TB, got, want []int) {
		t.Helper()
		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v want %v", got, want)
		}
	}

	t.Run("make the sums of tails of", func(t *testing.T) {
		got := SumAllTails([]int{1, 2}, []int{0, 9})
		want := []int{2, 9}
		checkSums(t, got, want)
	})

	t.Run("safely sum empty slices", func(t *testing.T) {
		got := SumAllTails([]int{}, []int{3, 4, 5})
		want := []int{0, 9}
		checkSums(t, got, want)
	})

}
```
이것의 유용한 부작용은 우리의 코드에 type-saftey를 조금 올려준다는 것이다. 만약 어리석은 개발자가 `checkSums(t, got, "dave")`를 사용하여 새로운 테스트를 한다면 컴파일러가 스스로 테스트들을 멈출 것이다.


```bash
$ go test
./sum_test.go:52:21: cannot use "dave" (type string) as type []int in argument to checkSums
```

## 마무리 하며

우리가 배운 것들은 다음과 같다.

* 배열
* 슬라이스
* 슬라이스를 만드는 여러 가지 방법
* 슬라이스들이 어떻게 _고정된_ 크기를 가지고 있지만 `append`를 사용하여 기존의 슬라이스를 새로운 슬라이스로 만드는 방법
* 슬라이스를 자르는 방법
* `len`을 이용하여 배열과 슬라이스의 크기를 받기
* 검사 방법 도구
* `reflect.DeepEqual`이 왜 유용하며, 왜 코드의 type-safety를 저해할 수 있는가

우리는 정수의 배열과 슬라이스를 사용하였지만, 다른 타입들에도 사용될 수 있다 배열과 슬라이스까지도. 그래서 당신은 당신은 변수를 다음과 같이 정의하면 된다.
`[][] string` 당신이 필요하다면.

[슬라이스에 관한 Go 언어 블로그를 읽어보길 바란다.][blog-slice] 슬라이스에 대해 더 깊게 알고 싶다면. 더 많은 테스트를 쓰고, 당신이 읽으면서 배운 것들을 입증해라.

테스트를 작성하는 것보다 Go 언어를 쉽게 테스트하는 방법은 Go playground를 이용하는 것이다. 당신은 많은 것들을 시도해 볼 수 있으며, 당신이 질문을 해야 할 경우 당신의 코드를 쉽게 공유할 수 있다. [당신이 슬라이스를 실험해 볼 수 있도록 go playground를 만들어 두었다.](https://play.golang.org/p/ICCWcRGIO68)

[여기 예시가 있다,](https://play.golang.org/p/bTrRmYfNYCp) 배열을 자르고이 슬라이스를 변경하는 것이 본래의 배열에 어떠한 영향을 주는지 알 수 있다.(하지만 `복제`된 슬라이스는 본래의 배열에 영향을 주지 못한다.)

[또 다른 예시가 있다,](https://play.golang.org/p/Poth8JS28sc) 매우 큰 슬라이스를 자르고 난 후 그들을 복제해두면 좋은 이유에 대해서 알 수 있다.

[for]: ../iteration.md#
[blog-slice]: https://blog.golang.org/go-slices-usage-and-internals
[deepEqual]: https://golang.org/pkg/reflect/#DeepEqual
[slice]: https://golang.org/doc/effective_go.html#slices