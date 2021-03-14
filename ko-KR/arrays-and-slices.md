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

 우리와 같은 경우, 아무도 우리의 함수를 사용하지 않기 때문에 두개의 함수를 유지하는것 보다 1개만 가지도록 합시다.
 
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



It is important to question the value of your tests. It should not be a goal to have as many tests as possible, but rather to have as much confidence as possible in your code base. Having too many tests can turn in to a real problem and it just adds more overhead in maintenance. Every test has a cost.

In our case, you can see that having two tests for this function is redundant. If it works for a slice of one size it's very likely it'll work for a slice of any size (within reason).

Go's built-in testing toolkit features a coverage tool, which can help identify areas of your code you have not covered. I do want to stress that having 100% coverage should not be your goal, it's just a tool to give you an idea of your coverage. If you have been strict with TDD, it's quite likely you'll have close to 100% coverage anyway.

