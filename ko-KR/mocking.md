# Mocking

**[You can find all the code for this chapter here](https://github.com/quii/learn-go-with-tests/tree/main/mocking)**


당신은 3부터 카운트 다운하는 프로그램을 제작할 것을 요청받았다. 이 프로그램은 새 라인에 각각의 수를 쓰면서 1초씩 멈춘다. 0에 도달하면 "Go!"가 인쇄되고 종료된다.

```
3
2
1
Go!
```

이 문제를 해결하려면 `Countdown`이라는 함수를 작성하고`main` 프로그램에 넣을 것이다.
```go
package main

func main() {
    Countdown()
}
```

이것은 매우 사소한 프로그램이지만 완전히 테스트하려면 항상 _iterative_, _test-driven_ 접근 방식을 취해야 한다.
Iterative란 무엇을 의미하는가? 우리는 _유용한 소프트웨어_를 만들기 위해 할 수 있는 가장 작은 단계를 정해야 한다.


우리는 약간의 해킹 후에 이론적으로 작동 할 코드에 오랜 시간을 보내고 싶지 않다. 그것이 종종 개발자들이 토끼 구멍에 빠지는 방식이기 때문이다. **요구 사항을 가능한 한 작게 분할하여 _동작하는 소프트웨어_를 가질 수있는 것은 중요한 기술이다.**

작업을 나누고 반복하는 방법은 다음과 같다.

- Print 3
- Print 3, 2, 1 and Go!
- Wait a second between each line

## Write the test first

소프트웨어는 stdout으로 출력해야하며 DI 섹션에서 이를 테스트하기 위해 DI를 사용하는 방법을 확인했다.

```go
func TestCountdown(t *testing.T) {
    buffer := &bytes.Buffer{}

    Countdown(buffer)

    got := buffer.String()
    want := "3"

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

`buffer`와 같은 것이 익숙하지 않다면 [이전 섹션] (dependency-injection.md)을 다시 읽어라.

우리는 'Countdown'함수가 데이터를 어딘가에 기록하기를 원한다는 것을 알고 있으며 'io.Writer'는 이를 Go의 인터페이스로 나타내는 사실상의 방법이다.

-`main`에서`os.Stdout`를 전송하여 사용자가 terminal에 인쇄 된 카운트 다운을 볼 수 있다.
-테스트에서 우리는 테스트에서 생성되는 데이터를 포착할 수 있도록 `bytes.Buffer`로 보낼 것이다.


## Try and run the test

`./countdown_test.go:11:2: undefined: Countdown`

## Write the minimal amount of code for the test to run and check the failing test output

Define `Countdown`

```go
func Countdown() {}
```

Try again

```go
./countdown_test.go:11:11: too many arguments in call to Countdown
    have (*bytes.Buffer)
    want ()
```

컴파일러가 함수 시그니처가 무엇인지 알려주므로 업데이트해라.

```go
func Countdown(out *bytes.Buffer) {}
```

`countdown_test.go:17: got '' want '3'`

Perfect!

## Write enough code to make it pass

```go
func Countdown(out *bytes.Buffer) {
    fmt.Fprint(out, "3")
}
```

우리는 `io.Writer`를 가지고 있는 `fmt.Fprint`을 사용한다. (`*bytes.Buffer`와 유사) 그리고 `string`을 그곳으로 보다. 테스트는 통과되어야 다.

## Refactor

우리는 `*bytes.Buffer` 대신 범용 인터페이스를 사용하는 것이 더 낫다는 것을 알고있다.

```go
func Countdown(out io.Writer) {
    fmt.Fprint(out, "3")
}
```

테스트를 다시 실행하면 통과해야 한다

문제를 완료하기 위해 이제 함수를 `main`에 연결하자. 우리는 발전하고 있음을 다시 확신 할 수 있도록 하는 동작하는 소프트웨어를 가지고 있다.

```go
package main

import (
    "fmt"
    "io"
    "os"
)

func Countdown(out io.Writer) {
    fmt.Fprint(out, "3")
}

func main() {
    Countdown(os.Stdout)
}
```

프로그램을 실행하고 당신의 작업에 놀라워하자.

그렇다, 이것은 사소한 것처럼 보이지만 이 접근 방식은 모든 프로젝트에 권장하는 것이다. **기능의 작은 부분을 취하고, 그것을 테스트를 통해 전체적으로도 동작하도록 하자. **

다음으로 2,1을 인쇄 한 다음 "Go!"를 인쇄 할 수 있다.

## Write the test first

전체 배관이 제대로 작동하도록 투자함으로써 우리는 솔루션을 안전하고 쉽게 반복할 수 있다. 모든 로직이 테스트 되므로 프로그램이 작동하는지 확인하기 위해 더 이상 프로그램을 중지하고 다시 실행할 필요가 없다.

```go
func TestCountdown(t *testing.T) {
    buffer := &bytes.Buffer{}

    Countdown(buffer)

    got := buffer.String()
    want := `3
2
1
Go!`

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}
```

The backtick syntax is another way of creating a `string` but lets you put things like newlines which is perfect for our test.

## Try and run the test

```
countdown_test.go:21: got '3' want '3
        2
        1
        Go!'
```
## Write enough code to make it pass

```go
func Countdown(out io.Writer) {
    for i := 3; i > 0; i-- {
        fmt.Fprintln(out, i)
    }
    fmt.Fprint(out, "Go!")
}
```

Use a `for` loop counting backwards with `i--` and use `fmt.Fprintln` to print to `out` with our number followed by a newline character. Finally use `fmt.Fprint` to send "Go!" aftward.

## Refactor

There's not much to refactor other than refactoring some magic values into named constants.

```go
const finalWord = "Go!"
const countdownStart = 3

func Countdown(out io.Writer) {
    for i := countdownStart; i > 0; i-- {
        fmt.Fprintln(out, i)
    }
    fmt.Fprint(out, finalWord)
}
```

If you run the program now, you should get the desired output but we don't have it as a dramatic countdown with the 1 second pauses.

Go lets you achieve this with `time.Sleep`. Try adding it in to our code.

```go
func Countdown(out io.Writer) {
    for i := countdownStart; i > 0; i-- {
        time.Sleep(1 * time.Second)
        fmt.Fprintln(out, i)
    }

    time.Sleep(1 * time.Second)
    fmt.Fprint(out, finalWord)
}
```

If you run the program it works as we want it to.

## Mocking

The tests still pass and the software works as intended but we have some problems:
- Our tests take 4 seconds to run.
    - Every forward thinking post about software development emphasises the importance of quick feedback loops.
    - **Slow tests ruin developer productivity**.
    - Imagine if the requirements get more sophisticated warranting more tests. Are we happy with 4s added to the test run for every new test of `Countdown`?
- We have not tested an important property of our function.

We have a dependency on `Sleep`ing which we need to extract so we can then control it in our tests.

If we can _mock_ `time.Sleep` we can use _dependency injection_ to use it instead of a "real" `time.Sleep` and then we can **spy on the calls** to make assertions on them.

## Write the test first

Let's define our dependency as an interface. This lets us then use a _real_ Sleeper in `main` and a _spy sleeper_ in our tests. By using an interface our `Countdown` function is oblivious to this and adds some flexibility for the caller.

```go
type Sleeper interface {
    Sleep()
}
```

I made a design decision that our `Countdown` function would not be responsible for how long the sleep is. This simplifies our code a little for now at least and means a user of our function can configure that sleepiness however they like.

Now we need to make a _mock_ of it for our tests to use.

```go
type SpySleeper struct {
    Calls int
}

func (s *SpySleeper) Sleep() {
    s.Calls++
}
```

_Spies_ are a kind of _mock_ which can record how a dependency is used. They can record the arguments sent in, how many times it has been called, etc. In our case, we're keeping track of how many times `Sleep()` is called so we can check it in our test.

Update the tests to inject a dependency on our Spy and assert that the sleep has been called 4 times.

```go
func TestCountdown(t *testing.T) {
    buffer := &bytes.Buffer{}
    spySleeper := &SpySleeper{}

    Countdown(buffer, spySleeper)

    got := buffer.String()
    want := `3
2
1
Go!`

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }

    if spySleeper.Calls != 4 {
        t.Errorf("not enough calls to sleeper, want 4 got %d", spySleeper.Calls)
    }
}
```

## Try and run the test

```
too many arguments in call to Countdown
    have (*bytes.Buffer, *SpySleeper)
    want (io.Writer)
```

## Write the minimal amount of code for the test to run and check the failing test output

We need to update `Countdown` to accept our `Sleeper`

```go
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := countdownStart; i > 0; i-- {
        time.Sleep(1 * time.Second)
        fmt.Fprintln(out, i)
    }

    time.Sleep(1 * time.Second)
    fmt.Fprint(out, finalWord)
}
```

If you try again, your `main` will no longer compile for the same reason

```
./main.go:26:11: not enough arguments in call to Countdown
    have (*os.File)
    want (io.Writer, Sleeper)
```

Let's create a _real_ sleeper which implements the interface we need

```go
type DefaultSleeper struct {}

func (d *DefaultSleeper) Sleep() {
    time.Sleep(1 * time.Second)
}
```

We can then use it in our real application like so

```go
func main() {
    sleeper := &DefaultSleeper{}
    Countdown(os.Stdout, sleeper)
}
```

## Write enough code to make it pass

The test is now compiling but not passing because we're still calling the `time.Sleep` rather than the injected in dependency. Let's fix that.

```go
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := countdownStart; i > 0; i-- {
        sleeper.Sleep()
        fmt.Fprintln(out, i)
    }

    sleeper.Sleep()
    fmt.Fprint(out, finalWord)
}
```

The test should pass and no longer take 4 seconds.

### Still some problems

There's still another important property we haven't tested.

`Countdown` should sleep before each print, e.g:

- `Sleep`
- `Print N`
- `Sleep`
- `Print N-1`
- `Sleep`
- `Print Go!`
- etc

Our latest change only asserts that it has slept 4 times, but those sleeps could occur out of sequence.

When writing tests if you're not confident that your tests are giving you sufficient confidence, just break it! (make sure you have committed your changes to source control first though). Change the code to the following

```go
func Countdown(out io.Writer, sleeper Sleeper) {
    for i := countdownStart; i > 0; i-- {
        sleeper.Sleep()
    }

    for i := countdownStart; i > 0; i-- {
        fmt.Fprintln(out, i)
    }

    sleeper.Sleep()
    fmt.Fprint(out, finalWord)
}
```

If you run your tests they should still be passing even though the implementation is wrong.

Let's use spying again with a new test to check the order of operations is correct.

We have two different dependencies and we want to record all of their operations into one list. So we'll create _one spy for them both_.

```go
type CountdownOperationsSpy struct {
    Calls []string
}

func (s *CountdownOperationsSpy) Sleep() {
    s.Calls = append(s.Calls, sleep)
}

func (s *CountdownOperationsSpy) Write(p []byte) (n int, err error) {
    s.Calls = append(s.Calls, write)
    return
}

const write = "write"
const sleep = "sleep"
```

Our `CountdownOperationsSpy` implements both `io.Writer` and `Sleeper`, recording every call into one slice. In this test we're only concerned about the order of operations, so just recording them as list of named operations is sufficient.

We can now add a sub-test into our test suite which verifies our sleeps and prints operate in the order we hope

```go
t.Run("sleep before every print", func(t *testing.T) {
    spySleepPrinter := &CountdownOperationsSpy{}
    Countdown(spySleepPrinter, spySleepPrinter)

    want := []string{
        sleep,
        write,
        sleep,
        write,
        sleep,
        write,
        sleep,
        write,
    }

    if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
        t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
    }
})
```

This test should now fail. Revert `Countdown` back to how it was to fix the test.

We now have two tests spying on the `Sleeper` so we can now refactor our test so one is testing what is being printed and the other one is ensuring we're sleeping in between the prints. Finally we can delete our first spy as it's not used anymore.

```go
func TestCountdown(t *testing.T) {

    t.Run("prints 3 to Go!", func(t *testing.T) {
        buffer := &bytes.Buffer{}
        Countdown(buffer, &CountdownOperationsSpy{})

        got := buffer.String()
        want := `3
2
1
Go!`

        if got != want {
            t.Errorf("got %q want %q", got, want)
        }
    })

    t.Run("sleep before every print", func(t *testing.T) {
        spySleepPrinter := &CountdownOperationsSpy{}
        Countdown(spySleepPrinter, spySleepPrinter)

        want := []string{
            sleep,
            write,
            sleep,
            write,
            sleep,
            write,
            sleep,
            write,
        }

        if !reflect.DeepEqual(want, spySleepPrinter.Calls) {
            t.Errorf("wanted calls %v got %v", want, spySleepPrinter.Calls)
        }
    })
}
```

We now have our function and its 2 important properties properly tested.

## Extending Sleeper to be configurable

A nice feature would be for the `Sleeper` to be configurable. This means that we can adjust the sleep time in our main program.

### Write the test first

Let's first create a new type for `ConfigurableSleeper` that accepts what we need for configuration and testing.

```go
type ConfigurableSleeper struct {
    duration time.Duration
    sleep    func(time.Duration)
}
```

We are using `duration` to configure the time slept and `sleep` as a way to pass in a sleep function. The signature of `sleep` is the same as for `time.Sleep` allowing us to use `time.Sleep` in our real implementation and the following spy in our tests:

```go
type SpyTime struct {
    durationSlept time.Duration
}

func (s *SpyTime) Sleep(duration time.Duration) {
    s.durationSlept = duration
}
```

With our spy in place, we can create a new test for the configurable sleeper.

```go
func TestConfigurableSleeper(t *testing.T) {
    sleepTime := 5 * time.Second

    spyTime := &SpyTime{}
    sleeper := ConfigurableSleeper{sleepTime, spyTime.Sleep}
    sleeper.Sleep()

    if spyTime.durationSlept != sleepTime {
        t.Errorf("should have slept for %v but slept for %v", sleepTime, spyTime.durationSlept)
    }
}
```

There should be nothing new in this test and it is setup very similar to the previous mock tests.

### Try and run the test
```
sleeper.Sleep undefined (type ConfigurableSleeper has no field or method Sleep, but does have sleep)

```

You should see a very clear error message indicating that we do not have a `Sleep` method created on our `ConfigurableSleeper`.

### Write the minimal amount of code for the test to run and check failing test output
```go
func (c *ConfigurableSleeper) Sleep() {
}
```

With our new `Sleep` function implemented we have a failing test.

```
countdown_test.go:56: should have slept for 5s but slept for 0s
```

### Write enough code to make it pass

All we need to do now is implement the `Sleep` function for `ConfigurableSleeper`.

```go
func (c *ConfigurableSleeper) Sleep() {
    c.sleep(c.duration)
}
```

With this change all of the tests should be passing again and you might wonder why all the hassle as the main program didn't change at all. Hopefully it becomes clear after the following section.

### Cleanup and refactor

The last thing we need to do is to actually use our `ConfigurableSleeper` in the main function.

```go
func main() {
    sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
    Countdown(os.Stdout, sleeper)
}
```

If we run the tests and the program manually, we can see that all the behavior remains the same.

Since we are using the `ConfigurableSleeper`, it is now safe to delete the `DefaultSleeper` implementation. Wrapping up our program and having a more [generic](https://stackoverflow.com/questions/19291776/whats-the-difference-between-abstraction-and-generalization) Sleeper with arbitrary long countdowns.

## But isn't mocking evil?

You may have heard mocking is evil. Just like anything in software development it can be used for evil, just like [DRY](https://en.wikipedia.org/wiki/Don%27t_repeat_yourself).

People normally get in to a bad state when they don't _listen to their tests_ and are _not respecting the refactoring stage_.

If your mocking code is becoming complicated or you are having to mock out lots of things to test something, you should _listen_ to that bad feeling and think about your code. Usually it is a sign of

- The thing you are testing is having to do too many things (because it has too many dependencies to mock)
    - Break the module apart so it does less
- Its dependencies are too fine-grained
    - Think about how you can consolidate some of these dependencies into one meaningful module
- Your test is too concerned with implementation details
    - Favour testing expected behaviour rather than the implementation

Normally a lot of mocking points to _bad abstraction_ in your code.

**What people see here is a weakness in TDD but it is actually a strength**, more often than not poor test code is a result of bad design or put more nicely, well-designed code is easy to test.

### But mocks and tests are still making my life hard!

Ever run into this situation?

- You want to do some refactoring
- To do this you end up changing lots of tests
- You question TDD and make a post on Medium titled "Mocking considered harmful"

This is usually a sign of you testing too much _implementation detail_. Try to make it so your tests are testing _useful behaviour_ unless the implementation is really important to how the system runs.

It is sometimes hard to know _what level_ to test exactly but here are some thought processes and rules I try to follow:

- **The definition of refactoring is that the code changes but the behaviour stays the same**. If you have decided to do some refactoring in theory you should be able to make the commit without any test changes. So when writing a test ask yourself
    - Am I testing the behaviour I want, or the implementation details?
    - If I were to refactor this code, would I have to make lots of changes to the tests?
- Although Go lets you test private functions, I would avoid it as private functions are implementation detail to support public behaviour. Test the public behaviour. Sandi Metz describes private functions as being "less stable" and you don't want to couple your tests to them.
- I feel like if a test is working with **more than 3 mocks then it is a red flag** - time for a rethink on the design
- Use spies with caution. Spies let you see the insides of the algorithm you are writing which can be very useful but that means a tighter coupling between your test code and the implementation. **Be sure you actually care about these details if you're going to spy on them**

#### Can't I just use a mocking framework?

Mocking requires no magic and is relatively simple; using a framework can make mocking seem more complicated than it is. We don't use automocking in this chapter so that we get:

- a better understanding of how to mock
- practise implementing interfaces

In collaborative projects there is value in auto-generating mocks. In a team, a mock generation tool codifies consistency around the test doubles. This will avoid inconsistently written test doubles which can translate to inconsistently written tests.

You should only use a mock generator that generates test doubles against an interface. Any tool that overly dictates how tests are written, or that use lots of 'magic', can get in the sea.

## Wrapping up

### More on TDD approach

- When faced with less trivial examples, break the problem down into "thin vertical slices". Try to get to a point where you have _working software backed by tests_ as soon as you can, to avoid getting in rabbit holes and taking a "big bang" approach.
- Once you have some working software it should be easier to _iterate with small steps_ until you arrive at the software you need.

> "When to use iterative development? You should use iterative development only on projects that you want to succeed."

Martin Fowler.

### Mocking

- **Without mocking important areas of your code will be untested**. In our case we would not be able to test that our code paused between each print but there are countless other examples. Calling a service that _can_ fail? Wanting to test your system in a particular state? It is very hard to test these scenarios without mocking.
- Without mocks you may have to set up databases and other third parties things just to test simple business rules. You're likely to have slow tests, resulting in **slow feedback loops**.
- By having to spin up a database or a webservice to test something you're likely to have **fragile tests** due to the unreliability of such services.

Once a developer learns about mocking it becomes very easy to over-test every single facet of a system in terms of the _way it works_ rather than _what it does_. Always be mindful about **the value of your tests** and what impact they would have in future refactoring.

In this post about mocking we have only covered **Spies** which are a kind of mock. The "proper" term for mocks though are "test doubles"

[> Test Double is a generic term for any case where you replace a production object for testing purposes.](https://martinfowler.com/bliki/TestDouble.html)

Under test doubles, there are various types like stubs, spies and indeed mocks! Check out [Martin Fowler's post](https://martinfowler.com/bliki/TestDouble.html) for more detail.
