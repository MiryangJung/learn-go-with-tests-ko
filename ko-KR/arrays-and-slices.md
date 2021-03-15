# 배열과 슬라이스

**[이 챕터의 모든 코드는 여기에서 확인할 수 있습니다.]((https://github.com/MiryangJung/learn-go-with-tests-ko/tree/maㅑㅜ/arrays)**

배열은 같은 타입의 여러 요소들을 한 변수에 특정한 순서대로 저장할 수 있게 해줍니다.   

우리가 배열을 다룰때에, 배열들에 대해 반복하는것은 매우 흔한 일입니다. 그래서 우리가[새롭게 배운 `for`](iteration.md)을 사용하여 `Sum`함수를 만들것입니다. `Sum`함수는 숫자 배열들을 가지고 총합을 return할것입니다.   

우리 `테스트 기반 설계` 스킬을 활용해봅시다.   


## 테스트를 먼저 작성하세요   


작업을 할 새로운 폴더를 생성합시다. `sum_test.go`라는 파일을 생성한 후 아래 코드를 집어넣습니다.
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

배열은 변수로 선언될때 우리가 정한 고정된 크기를 가집니다. 우리는 두가지 방법으로 배열을 선언할 수 있습니다. 
* \[N\]type{value1, value2, ..., valueN} e.g. `numbers := [5]int{1, 2, 3, 4, 5}`
* \[...\]type{value1, value2, ..., valueN} e.g. `numbers := [...]int{1, 2, 3, 4, 5}`

[문자열 포멧팅에 대해서 더 읽어봅시다.](https://golang.org/pkg/fmt/)

## 테스트를 실행해 보세요

 `go test`로 실행을 시키면 다음과 같이 컴파일 에러가 날것입니다. `./sum_test.go:10:15: undefined: Sum`

 ## 테스트를 실행할 최소한의 코드를 작성하고 실패한 테스트 출력을 확인하세요.

 `sum.go`를 다음과 같이 한다면,

```go
package main

func Sum(numbers [5]int) int {
	return 0
}
```

당신은  테스트의 조금더 명확한 에러 메세지를 확인할 수 있습니다.

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
배열의 특정한 인덱스값을 가져오려고 한다면, `array[index]`라는 syntax를  사용하면 됩니다. 지금과 같은 경우에는, 우 `for`을 이용하여 5번을 반복하며 배열의 각 항목들을 sum에 더할것 입니다.

## 리펙토링 
우리 `range`를 사용하여 코드를 간추려봅시다.

```go
func Sum(numbers [5]int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
```
`range`는 배열에 반복적인 처리를 할 수 있게 합니다. `range`는 호출될때 마다 인덱스와 값, 총 두개의 값을 리턴합니다. 우리는 여백 식별자인 `_`를 사용하여 인덱스의 값을 무시할것 입니다.

### 배열과 그 타입

배열의 흥미로운 특성은, 배열의 크기는 배열의 타입으로 인코딩 된다는 점이다. 따라서 `[5]int`를 요구하는 함수에 `[4]int`가 입력될 경우 컴파일 되지 않을것입니다. 이 이유는 바로 타입이 다르기 때문입니다, 이는 `int`를 요구하는  함수에 `string`을 입력한것과 같습니다.

당신은 이제 배열이 고정된 크기를 가진다는 점때문에 복잡하다고 생각하겠지만, 대부분의 경우 당신은 배열을 사용하지 않을것입니다.

Go언어에서의 슬라이스는 collection의 크기를 같이 인코딩 하지않아 아무 크기나 가질 수 있습니다.

다음요건은 달라지는 크기의 collection의 합을 구하는것입니다.

## 테스트를 먼저 작성하세요

우리는 이제 collection의 크기가 자유로운 `slice`타입을 이용할 것 입니다. 슬라이스의 syntax는  배열과 매우 유사합니다, 그저 선언할때 크기를 생략하기만 하면 됩니다.

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

위 코드는 컴파일 되지 않습니다.
`./sum_test.go:22:13: cannot use numbers (type []int) as type [5]int in argument to Sum`

 ## 테스트를 실행할 최소한의 코드를 작성하고 실패한 테스트 출력을 확인하세요.
 
 여기서의 문제는 둘 중 하나이다.
 - 기존에 존제하던 API의 인수인 `Sum`을 배열대신 슬라이스로 바꾸어 주는 것이다. 만약 우리가 이방법을 사용한다면 우리는 다른 사람의 하루를 망쳤을수도 있습니다. 이는 우리의 다른 테스트들을 컴파일 되지 않도록 하기 때문입니다.
 - 새로운 함수를 만듭니다.

 지금과 같은 경우, 아무도 우리의 함수를 사용하지 않을거기 때문에 두개의 함수를 유지하는것 보다 1개만 가지도록 합시다.
 
```go
func Sum(numbers []int) int {
	sum := 0
	for _, number := range numbers {
		sum += number
	}
	return sum
}
```
당신이 컴파일 하려고 한다면, 이는 아직도 컴파일 되지 않을것 입니다. 당신은 배열 이 아닌 슬라이스로 pass하려고 한다면 첫번째 테스트를 수정해야 합니다.

## 통과할 수 있도록 충분한 코드를 작성하세요
알고보니 우리는 컴파일러의 문제만 해결했으면 됩니다.테스트는 통과했습니다.

## 리펙토링 
우리는 이미 `Sum`을 리펙토링을 하였으며, 우리는 그저 배열을 슬라이스로 바꾸기만 했습니다. 그러므로 이 부분에서는 많이 할것이 없습니다. 우리는 리펙토링 단계에서 테스트 코드를 소홀히 하면 안되고 이 단계에서 우리가 해야 할 것이 있다는것을 기억해야합니다.

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

당신의 테스트값에 대하여 의문을 갖는것은 중요합니다. 테스트 

It is important to question the value of your tests. It should not be a goal to have as many tests as possible, but rather to have as much confidence as possible in your code base. Having too many tests can turn in to a real problem and it just adds more overhead in maintenance. Every test has a cost.

지금과 같은 경우, 한 함수에 대하여 2개의 테스트를 가지고 있는것은 매우 불필요합니다. 만약 함수가 한가지 크기의 슬라이스로 통과를 한다면, 이는 다른 크기의 슬라이스도 높은 확률로 통과 할수 있다는 뜻입니다(합당한 범위 내에서).

Go언어의 내장된 테스팅 도구는 커버리지 도구가 있습니다, 당신의 다루지 않는 부분의 코드를  확인할 수 있도록 돕습니다. 저는  100% 커버리지가 당신의 목표가  아니라는것을 강조하고 싶습니다. 이는 그저 당신에게 어느정도가 커버리지가 되는지 알려주기 위함입니다. 당신이 TDD에 관하여 엄격하였다면, 높은 확률로 당신의 커버리지는 100%로 끝날거기때입니다.

아래를 실행해 봅시다.

`go test -cover`

당신은 다음을 확인할 수 있습니다.

```bash
PASS
coverage: 100.0% of statements
```

이제 테스트 하나를 지우고 커버리지를 다시 확인합시다.

우리는 이제 잘 테스트된 함수를 가지고 있습니다, 이제 다음도전을 하기 전에 우리의 멋진 작업을 커밋하도록 합시다.

우리는 `SumAll`이라는 여러 슬라이스를 받고 그들의 합을 슬라이스들로 리턴할 새로운 함수가 필요합니다.

예시로,

`SumAll([]int{1,2}, []int{0,9})` 는 `[]int{3, 9}`을 리턴할 것입니다.

또는

`SumAll([]int{1,1,1})` 는 `[]int{3}`를 리턴할 것입니다.

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

우리는 SumAll을 테스트가 요구하는대로 정의를 해야합니다.

Go 언어에서는 변동가능한 수의 인수들을 받을 수 있는 [_variadic functions_](https://gobyexample.com/variadic-functions)을 사용할 수 있게 해줍니다.

```go
func SumAll(numbersToSum ...[]int) (sums []int) {
	return
}
```

컴파일을 시도해 보세요, 하지만 우리의 테스트는 아직 컴파일 하지 않습니다!

`./sum_test.go:26:9: invalid operation: got != want (slice can only be compared to nil)`

Go언어에서는 slice를 다룰때에 등호를 사용할 수 없습니다. 당신은 `got`슬라이스와 `want`슬라이스를 반복하여 그들의 값을 비교하는 함수를 만들 수 있습니다. 하지만 이 방법은 편리하지 않습니다. 그래서 우리는 `아무` 타입이나 비교할 수 있는 [`reflect.DeepEqual`][deepEqual] 를 사용할 수 있습니다.

```go
func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := []int{3, 9}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```

\(파일 상단에 `import reflect`가 있어야 `DeepEqual`를 사용할 수 있습니다.\)

우리가 집고 넘어가야하는 부분은 `reflect.DeepEqual`는 "type safe"하지 않다는 것이다. 당신이 멍청한 짓을 했어도 코드는 그냥 컴파일 될 것 이다.이것을 확인하고 싶다면 테스트를 잠시동안 다음과 같이 바꿔도보록 합시다.
```go
func TestSumAll(t *testing.T) {

	got := SumAll([]int{1, 2}, []int{0, 9})
	want := "bob"

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
```
우리가 위에서 하고자 하는것은 `슬라이스`와 `문자열`을 비교하려고 하고 있다는 것이다. 이것은 터무니 없는 일이지만 테스트는 컴파일 되었다. 따라서 `reflect.DeepEqual`를 사용하여 슬라이스를 \(그리고 다른 문자열들을\) 비교하는것은 매우 편리하지만 사용할때에 주의를 요한다.

테스트를 다시 원상복구 시키고 실행한다면, 다음과 같은 출력값을 보게 될것이다.

`sum_test.go:30: got [] want [3 9]`

## 통과할 수 있도록 충분한 코드를 작성하세요



What we need to do is iterate over the varargs, calculate the sum using our
`Sum` function from before and then add it to the slice we will return

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

Lots of new things to learn!

There's a new way to create a slice. `make` allows you to create a slice with
a starting capacity of the `len` of the `numbersToSum` we need to work through.

You can index slices like arrays with `mySlice[N]` to get the value out or
assign it a new value with `=`

The tests should now pass

## Refactor

As mentioned, slices have a capacity. If you have a slice with a capacity of
2 and try to do `mySlice[10] = 1` you will get a _runtime_ error.

However, you can use the `append` function which takes a slice and a new value,
returning a new slice with all the items in it.

```go
func SumAll(numbersToSum ...[]int) []int {
	var sums []int
	for _, numbers := range numbersToSum {
		sums = append(sums, Sum(numbers))
	}

	return sums
}
```

In this implementation, we are worrying less about capacity. We start with an
empty slice `sums` and append to it the result of `Sum` as we work through the varargs.

Our next requirement is to change `SumAll` to `SumAllTails`, where it now
calculates the totals of the "tails" of each slice. The tail of a collection is
all the items apart from the first one \(the "head"\)

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

Rename the function to `SumAllTails` and re-run the test

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

Slices can be sliced! The syntax is `slice[low:high]` If you omit the value on
one of the sides of the `:` it captures everything to the side of it. In our
case, we are saying "take from 1 to the end" with `numbers[1:]`. You might want to
invest some time in writing other tests around slices and experimenting with the
slice operator so you can be familiar with it.

## 리팩토링

Not a lot to refactor this time.

What do you think would happen if you passed in an empty slice into our
function? What is the "tail" of an empty slice? What happens when you tell Go to
capture all elements from `myEmptySlice[1:]`?

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

오, 이런! 우리는 테스트가 _컴파일_ 되었는지 유의해야 합니다. 위 에러는 런타임 에러 입니다. 

Compile time errors are our friend because they help us write software that
works, runtime errors are our enemies because they affect our users.

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

우리의 코드는 assertion 부분에서 코드가 반복되는 일이 또 발생하였습니다. 우리 이 부분을 함수로 추출해봅시다.

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

A handy side-effect of this is this adds a little type-safety to our code. If
a silly developer adds a new test with `checkSums(t, got, "dave")` the compiler
will stop them in their tracks.

```bash
$ go test
./sum_test.go:52:21: cannot use "dave" (type string) as type []int in argument to checkSums
```

## 마무리 하며

우리가 배운것들은 다음과 같습니다.

* 배열
* 슬라이스
  * 슬라이스를 만드는 여러가지 방법
  * 스라이스들이 어떻게 _고정된_ 크기를 가지고 있지만 `append`를 사용하여 기존의 슬라이스를 새로운 슬라이스로 만드는방법
  * 슬라이스를 자르는 방법
* `len`을 이용하여 배열과 슬라이스의 크기를 받기
* 검사 방법 도구
* `reflect.DeepEqual`이 왜 유용하며, 왜 코드의 type-safety를 저해 할  수 있는지

우리는 정수의 배열과 슬라이스를 사용하였지만, 다른 타입들에도 사용될 수 있습니다 배열과 슬라이스 까지도. 그래서 당신은 당신은 변수를 다음과 같이 정의하면 됩니다.
`[][]string` 당신이 필요하다면.

[슬라이스에 관한 Go언어 블로그를 읽어보시길 바랍니다.][blog-slice] 슬라이스에 대해 더 깊게 알고 싶으시다면. 더 많은 테스트를 쓰고,당신이 읽으면서 배운것들을 입증해봅시다.

테스트를 작성하는 것 보다 Go언어를 쉽게 테스트 하는 방법은 Go playground를 이용하는 것입니다. 당신은 많은 것들을 시도해 볼 수 있으며, 당신이 질문을 해야 할 경우 당신의 코드를 쉽게 공유할 수 있습니다. [당신이 슬라이스를 실험 해볼 수 있도록 go playground를 만들어 두었습니다.](https://play.golang.org/p/ICCWcRGIO68)

[여기 예시가 있습니다,](https://play.golang.org/p/bTrRmYfNYCp) 에레이를 슬라이스화하고 이 슬라이스를 변경하는것이 본래의 배열에 어떠한 영향을 주는지; 하지만 `복제`된 슬라이스는 본래의 배열에 영향을 주지 못합니다.

[또다른 예시가 있습니다,](https://play.golang.org/p/Poth8JS28sc) of why it's a good idea
to make a copy of a slice after slicing a very large slice.

[for]: ../iteration.md#
[blog-slice]: https://blog.golang.org/go-slices-usage-and-internals
[deepEqual]: https://golang.org/pkg/reflect/#DeepEqual
[slice]: https://golang.org/doc/effective_go.html#slices