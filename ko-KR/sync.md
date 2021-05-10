# Sync

**[이 챕터의 모든 코드는 여기에서 확인할 수 있다.](https://github.com/quii/learn-go-with-tests/tree/main/sync)**

안전하게 동시적으로 (concurrently) 사용할 수 있는 카운터를 작성한다

먼저 동시적으로 안전하지 않은 카운터를 만들어 싱글-스레드 환경에서 정상적으로 작동하는지 확인하겠다.

이후에 해당 카운터의 불안정함을 여러 고루틴들 (goroutines) 을 통하여 테스트 하여 수정할 것이다.

## Write the test first

카운터를 증가시킨 다음 그 값을 알리는 메서드를 제공하는 API를 작성한다.

```go
func TestCounter(t *testing.T) {
	t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
		counter := Counter{}
		counter.Inc()
		counter.Inc()
		counter.Inc()

		if counter.Value() != 3 {
			t.Errorf("got %d, want %d", counter.Value(), 3)
		}
	})
}
```

## Try to run the test

```
./sync_test.go:9:14: undefined: Counter
```

## Write the minimal amount of code for the test to run and check the failing test output

`Counter`를 정의 (define) 한다.

```go
type Counter struct {

}
```

다시 한번 시도하면 다음과 같이 실패할 것이다.

```
./sync_test.go:14:10: counter.Inc undefined (type Counter has no field or method Inc)
./sync_test.go:18:13: counter.Value undefined (type Counter has no field or method Value)
```

따라서 테스트가 성공할 수 있도록 다음과 같은 메서드들을 정의한다.

```go
func (c *Counter) Inc() {

}

func (c *Counter) Value() int {
	return 0
}
```

테스트를 다시 실행하면 다음과 같이 실패할 것이다.

```
=== RUN   TestCounter
=== RUN   TestCounter/incrementing_the_counter_3_times_leaves_it_at_3
--- FAIL: TestCounter (0.00s)
    --- FAIL: TestCounter/incrementing_the_counter_3_times_leaves_it_at_3 (0.00s)
    	sync_test.go:27: got 0, want 3
```

## Write enough code to make it pass

이는 우리와 같은 Go 전문가들에게는 사소한 일이다. 카운터의 상태를 데이터 타입에 저장한 후 모든 `Inc` 호출에서 증가시키면 된다.

```go
type Counter struct {
	value int
}

func (c *Counter) Inc() {
	c.value++
}

func (c *Counter) Value() int {
	return c.value
}
```

## Refactor

리팩토링 할 내용이 많지는 않지만 `Counter`에 대해 더 많은 테스트를 작성할 예정이므로 `assertCount` 라는 작은 assertion (확인) 함수를 작성하여 테스트가 좀 더 명확해질 수 있도록 작성한다.

```go
t.Run("incrementing the counter 3 times leaves it at 3", func(t *testing.T) {
    counter := Counter{}
    counter.Inc()
    counter.Inc()
    counter.Inc()

    assertCounter(t, counter, 3)
})

func assertCounter(t testing.TB, got Counter, want int)  {
	t.Helper()
	if got.Value() != want {
		t.Errorf("got %d, want %d", got.Value(), want)
	}
}
```

## Next steps

여기까지는 그리 어렵지 않았으나 이제 동시적 환경에서 안전하게 사용될 수 있도록 조건들을 추가해야 한다. 이를 위해 실패할 테스트를 작성한다.

## Write the test first

```go
t.Run("it runs safely concurrently", func(t *testing.T) {
    wantedCount := 1000
    counter := Counter{}

    var wg sync.WaitGroup
    wg.Add(wantedCount)

    for i := 0; i < wantedCount; i++ {
        go func(w *sync.WaitGroup) {
            counter.Inc()
            w.Done()
        }(&wg)
    }
    wg.Wait()

    assertCounter(t, counter, wantedCount)
})
```

위의 테스트는 `wantedCount` 값 만큼 반복문을 실행하며, `counter.Inc()`라는 고루틴 (goroutine) 을 실행 할 것이다.

동시적 프로세스들을 동기화 하는데 간편한 방법인 [`sync.WaitGroup`](https://golang.org/pkg/sync/#WaitGroup)을 이용한다.

> WaitGroup은 관련 고루틴들이 완료되기를 기다립니다. 메인 고루틴은 기다려야할 고루틴의 숫자를 설정하기 위하여 Add를 호출합니다. 이후 각각의 고루틴들은 실행후 완료되었음을 알리기 위해 Done을 호출합니다. 이와 동시에 모든 고루틴이 완료 될때까지 Wait을 사용하여 차단 (block) 할 수 있습니다.

`wg.Wait()`를 사용하여 모든 고루틴들이 `Counter`에 대해 `Inc`를 시도하였음을 assertion 하기 전에 분명히 할 수 있다.

## Try to run the test

```
=== RUN   TestCounter/it_runs_safely_in_a_concurrent_envionment
--- FAIL: TestCounter (0.00s)
    --- FAIL: TestCounter/it_runs_safely_in_a_concurrent_envionment (0.00s)
    	sync_test.go:26: got 939, want 1000
FAIL
```

테스트는 _아마도_ 다른 결과값을 가지고 실패할 것이다. 성공의 여부와 별개로 해당 테스트는 여러 고루틴이 동시에 카운터값을 수정하려 시도할 때 의도한 대로 작동하지 않는다는 것을 보여준다.

## Write enough code to make it pass

간단한 해결책은 `Counter`에 락 (lock), 즉 [`Mutex`](https://golang.org/pkg/sync/#Mutex) (상호 배제) 를 추가하는 것이다.

> Mutex는 상호 배제 잠금으로 0 은 잠금 해제 된 상태를 의미합니다.

```go
type Counter struct {
	mu sync.Mutex
	value int
}

func (c *Counter) Inc() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value++
}
```

이는  `Inc`를 호출하는 첫 번째 고루틴이 `Counter`에 대한 락을 획득함을 의미한다. 다른 고루틴들은 접근 권한을 얻기 위해서 해당 락이 'Unlock' 될 때까지 기다려야 한다.

이제 각각의 고루틴이 해당 값을 바꾸기 위해 차례를 지키며 기다리기에 테스트를 다시 실행하면 성공적으로 통과할 것이다.

## I've seen other examples where the `sync.Mutex` is embedded into the struct.

다음과 같은 예시를 찾아볼 수도 있다.

```go
type Counter struct {
	sync.Mutex
	value int
}
```

이와 같은 예시는 코드를 조금 더 우아하게 만들 수 있다고 주장될 수 있다.

```go
func (c *Counter) Inc() {
	c.Lock()
	defer c.Unlock()
	c.value++
}
```

이와 같은 코드는 프로그래밍이 굉장히 주관적인 원칙에 의해 진행되므로 _좋아_ 보일 수 있으나, 사실 **나쁘고 잘못된** 코드이다.

때때로 사람들은 임베딩 (embedding) 유형은 곧 해당 유형의 메서드가 _공개 인터페이스의 일부_ 가 됨을 잊고 있으며 동시에 그렇게 되지 않기를 원한다. 우리는 Public APIs 에 대해 매우 주의해야 하는데 이는 해당 APIs 공개하는 순간 다른 코드들과 결합되어 (couple) 사용되기 때문이다. 우리는 항상 불필요한 결합을 피하기를 원한다.

`Lock` 과 `Unlock` 의 노출은 이와 관련된 메서드들을 호출하기 시작하는 순간 해당 소프트웨어에 최선의 경우 혼란을, 최악의 경우 잠재적으로 굉장히 해로울 수 있다.

![이와 같은 API를 사용하는 유저가 락의 상태를 잘못된 상태로 사용할 수 있는지를 보여주는 예시](https://i.imgur.com/SWYNpwm.png)

_굉장히 좋지 않은 생각이다_

## Copying mutexes

테스트는 모두 통과했지만, 코드는 아직 위험성을 내포하고 있다.

`go vet`을 호출하면 다음과 같은 오류를 확인할 수 있다.

```
sync/v2/sync_test.go:16: call of assertCounter copies lock value: v1.Counter contains sync.Mutex
sync/v2/sync_test.go:39: assertCounter passes lock by value: v1.Counter contains sync.Mutex
```

[`sync.Mutex`](https://golang.org/pkg/sync/#Mutex) 를 확인하여 이유를 확인할 수 있다.

> Mutex (상호배제)는 처음 사용된 후 복사되어서는 안 됩니다.

`Counter`를 값으로 `assertCounter` 함수에 전달하기에 해당 mutex의 복사본을 생성하게 된다.

이 문제를 해결하기 위해서 `Counter`를 값이아닌 포인터로 전달 할 수 있도록 `assertCounter`의 서명 (signature) 을 다음과 같이 변경한다.

```go
func assertCounter(t *testing.T, got *Counter, want int)
```

`* Counter`가 아닌 `Counter`를 전달하려고 하기 때문에 테스트가 더 이상 컴파일되지 않는다. 이를 해결하기 위해 API 유저에게 유형을 직접 초기화하지 않는 것이 더 낫다는 것을 보여주는 생성자를 만드는 것을 (개인적으로) 선호한다.

```go
func NewCounter() *Counter {
	return &Counter{}
}
```

`Counter`를 초기화 할때 해당 함수를 테스트에서 사용한다.

## Wrapping up

[동기화 패키지](https://golang.org/pkg/sync/) 에서 몇 가지를 다루었다.

-`Mutex`를 사용하여 데이터에 잠금을 추가 할 수 있다.
-`Waitgroup`은 고루틴이 작업을 완료하기를 기다리는 수다.


### When to use locks over channels and goroutines?

[우리는 이전에 첫 번째 동시성 챕터에서 고루틴을 다뤘다](concurrency.md) 이를 통하여 안전한 동시성 코드를 작성할 수 있게 해주는데 굳이 lock을 사용할 필요가 있을까?
[go wiki에는 이 주제에 대한 전용 페이지가 있다; Mutex 또는 채널](https://github.com/golang/go/wiki/MutexOrChannel)

> Go 초보자의 일반적인 실수는 단지 가능하거나 재미있다는 이유로 채널과 고루틴을 과도하게 사용하는 것이다. 해당 문제에 가장 적합하다는 판단이 들 때 sync.Mutex를 사용하는 것을 두려워하지 말아야 한다. Go는 문제를 가장 잘 해결하는 도구를 사용하고 한 가지 스타일의 코드를 강요하지 않도록 하는 데에 있어서 실용적이기 때문이다.

이를 다르게 말하면:

- **데이터의 소유권을 전달할 때 채널 사용**
- **상태 관리를 위할 때 Mutexes 사용**

### go vet

`go vet`은 당신의 코드에 존재하는 미묘한 버그들에 대해 주의를 줄 수 있으니 유저들이 당신의 코드로 고통받지 않도록 빌드 스크립트에 사용하는 것을 잊지 말아야 한다.

### Don't use embedding because it's convenient

- 임베딩이 공용 API에 미치는 영향에 대해 생각해보라.
- _정말_ 해당 메서드들을 노출시켜 유저들이 해당 코드들을 다른 코드와 결합할 수 있도록 하고 싶은가?
- Mutex는 굉장히 예측할 수 없으며 이상한 형태로 잠재적 재앙이 될 수 있다. Mutex의 잠금을 해제하는 악의적인 코드를 상상해보아야 한다. 이는 추적 하기 어려운 굉장히 이상한 버그를 일으킬 것이다.
