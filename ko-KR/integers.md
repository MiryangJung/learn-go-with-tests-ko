# 정수 (Integers)

**[이 챕터의 모든 코드는 여기에서 확인할 수 있습니다.](https://github.com/MiryangJung/learn-go-with-tests-ko/tree/master/integers)**

정수는 예상한 대로 동작합니다. 시도하기 위해 `Add` 함수를 작성합니다. `adder_test.go`라는 테스트 파일을 생성한 후 이 코드를 작성합니다.

**주의:** Go의 소스 파일은 각 폴더에 하나의 `package`만 가질 수 있습니다. 파일이 별도로 구성되어 있는지 확인해야 합니다. [여기에 좋은 예시가 있습니다.](https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project)

## 테스트 먼저 작성하세요.

```go
package integers

import "testing"

func TestAdder(t *testing.T) {
	sum := Add(2, 2)
	expected := 4

	if sum != expected {
		t.Errorf("expected '%d' but got '%d'", expected, sum)
	}
}
```

`%q`가 아닌 `%d`로 형식 문자열을 사용하고 있음을 알 수 있습니다. 왜냐하면 문자열보다는 정수를 출력하기 원하기 때문입니다.

또한 더 이상 main 패키지를 사용하지 않으며 대신 `integers`라는 패키지를 정의했습니다. 이름에서 알 수 있듯이 `Add`와 같은 정수를 처리하기 위한 함수를 그룹화합니다.

## 테스트를 시도해보세요.

`go test` 명령어를 이용해 테스트를 실행합니다.

컴파일 에러를 확인할 수 있습니다.

`./adder_test.go:6:9: undefined: Add`

## 테스트를 실행할 최소한의 코드를 작성하고 실패한 테스트 출력을 확인하세요.

컴파일러를 만족시킬 수 있는 최소한의 코드를 작성합니다. _그것이 전부입니다._ - 우리의 테스트 테스트가 올바른 이유로 실패하는 것을 확인하길 원하는 것을 기억해야 합니다.

```go
package integers

func Add(x, y int) int {
	return 0
}
```

두 개 이상의 같은 타입의 인자를 갖는다면 \(지금 케이스에서는 두 개의 정수\) `(x int, y, int)`보다 `(x, y int)`로 짧게 할 수 있습니다.

이제 테스트를 실행하고 테스트가 올바르게 잘못된 내용을 보고하고 있다는 사실에 만족해야 합니다.

`adder_test.go:10: expected '4' but got '0'`

If you have noticed we learnt about _named return value_ in the [last](hello-world.md#one...last...refactor?) section but aren't using the same here. It should generally be used when the meaning of the result isn't clear from context, in our case it's pretty much clear that `Add` function will add the parameters. You can refer [this](https://github.com/golang/go/wiki/CodeReviewComments#named-result-parameters) wiki for more details.

## Write enough code to make it pass

In the strictest sense of TDD we should now write the _minimal amount of code to make the test pass_. A pedantic programmer may do this

```go
func Add(x, y int) int {
	return 4
}
```

Ah hah! Foiled again, TDD is a sham right?

We could write another test, with some different numbers to force that test to fail but that feels like [a game of cat and mouse](https://en.m.wikipedia.org/wiki/Cat_and_mouse).

Once we're more familiar with Go's syntax I will introduce a technique called _"Property Based Testing"_, which would stop annoying developers and help you find bugs.

For now, let's fix it properly

```go
func Add(x, y int) int {
	return x + y
}
```

If you re-run the tests they should pass.

## Refactor

There's not a lot in the _actual_ code we can really improve on here.

We explored earlier how by naming the return argument it appears in the documentation but also in most developer's text editors.

This is great because it aids the usability of code you are writing. It is preferable that a user can understand the usage of your code by just looking at the type signature and documentation.

You can add documentation to functions with comments, and these will appear in Go Doc just like when you look at the standard library's documentation.

```go
// Add takes two integers and returns the sum of them.
func Add(x, y int) int {
	return x + y
}
```

### Examples

If you really want to go the extra mile you can make [examples](https://blog.golang.org/examples). You will find a lot of examples in the documentation of the standard library.

Often code examples that can be found outside the codebase, such as a readme file often become out of date and incorrect compared to the actual code because they don't get checked.

Go examples are executed just like tests so you can be confident examples reflect what the code actually does.

Examples are compiled \(and optionally executed\) as part of a package's test suite.

As with typical tests, examples are functions that reside in a package's `_test.go` files. Add the following `ExampleAdd` function to the `adder_test.go` file.

```go
func ExampleAdd() {
	sum := Add(1, 5)
	fmt.Println(sum)
	// Output: 6
}
```

(If your editor doesn't automatically import packages for you, the compilation step will fail because you will be missing `import "fmt"` in `adder_test.go`. It is strongly recommended you research how to have these kind of errors fixed for you automatically in whatever editor you are using.)

If your code changes so that the example is no longer valid, your build will fail.

Running the package's test suite, we can see the example function is executed with no further arrangement from us:

```bash
$ go test -v
=== RUN   TestAdder
--- PASS: TestAdder (0.00s)
=== RUN   ExampleAdd
--- PASS: ExampleAdd (0.00s)
```

Please note that the example function will not be executed if you remove the comment `//Output: 6`. Although the function will be compiled, it won't be executed.

By adding this code the example will appear in the documentation inside `godoc`, making your code even more accessible.

To try this out, run `godoc -http=:6060` and navigate to `http://localhost:6060/pkg/`

Inside here you'll see a list of all the packages in your `$GOPATH`, so assuming you wrote this code in somewhere like `$GOPATH/src/github.com/{your_id}` you'll be able to find your example documentation.

If you publish your code with examples to a public URL, you can share the documentation of your code at [pkg.go.dev](https://pkg.go.dev/). For example, [here](https://pkg.go.dev/github.com/quii/learn-go-with-tests/integers/v2) is the finalised API for this chapter. This web interface allows you to search for documentation of standard library packages and third-party packages.

## Wrapping up

What we have covered:

-   More practice of the TDD workflow
-   Integers, addition
-   Writing better documentation so users of our code can understand its usage quickly
-   Examples of how to use our code, which are checked as part of our tests
