#  포인터 & 에러

**[이 챕터의 모든 코드는 여기에서 확인할 수 있다.](https://github.com/quii/learn-go-with-tests/tree/main/pointers)**

지난번 섹션에서 저희는 한 개념 안에 여러 많은 수의 값을 포함 할 수 있는 구조체에 대해서 배웠다.

어떤 상황에서는 구조체를 상태를 관리하는 데 있어 사용할 수 있고, 다른 방법으로는  유저가 제어하는 데로 상태를 바꿀 수 있도록 하는 메서드를 노출하는 방법이 있다.   

**핀테크는 Go를 좋아한다** 그리고 음.. 비트코인도? 따라서 우리가 얼마나 대단한 은행 시스템을 만들 수 있는지 보일 것이다.

`Wallet` 구조체를 만들고 `Bitcoin`을 입금해 보자.

## 테스트부터 작성하기

```go
func TestWallet(t *testing.T) {

    wallet := Wallet{}

    wallet.Deposit(10)

    got := wallet.Balance()
    want := 10

    if got != want {
        t.Errorf("got %d want %d", got, want)
    }
}
```

[이전 예제](./structs-methods-and-interfaces.md)에서 필드에 접근할 때 직접적으로 필드 이름에 접근했었지만, 우리의 _매우 보안적인 wallet_ 에서는 우리의 내부 상태를 밖으로 노출하기를 원하지 않는다. 우리는 메서드를 통해서 접근을 제어하기를 원한다.

## 테스트 실행해보기

`./wallet_test.go:7:12: undefined: Wallet`

## 컴파일이 되는 최소한의 코드를 작성하고, 테스트 실패 출력을 확인하기 

컴파일러는 `Wallet`이 무엇인지 모르기 때문에 알려줘야 한다.

```go
type Wallet struct { }
```

이제 우리의 wallet을 만들었으니, 테스트를 실행시켜서 확인한다.

```go
./wallet_test.go:9:8: wallet.Deposit undefined (type Wallet has no field or method Deposit)
./wallet_test.go:11:15: wallet.Balance undefined (type Wallet has no field or method Balance)
```

우리는 위의 메서드를 정의해야 한다.

해야 할 일은 테스트가 충분히 동작하도록 하는 것임을 기억해야 한다. 우리의 테스트가 깔끔한 오류 메시지와 함께 정확하게 실패하도록 해야 한다.
```go
func (w Wallet) Deposit(amount int) {

}

func (w Wallet) Balance() int {
    return 0
}
```

만약 이 문법이 익숙하지 않다면 이전 섹션으로 돌아가서 구조체에 대해 다시 읽으시길 바란다.

이제 테스트는 컴파일되어 아래와 같이 동작한다.

`wallet_test.go:15: got 0 want 10`

## 테스트를 통과하는 최소한의 코드 작성하기 

우리의 구조체에 상태를 저장하기 위해 일종의 _balance_ 변수가 필요하다. 

```go
type Wallet struct {
    balance int
}
```

Go에서는 만약 symbol(변수, 타입, 함수 등)이 소문자로 시작한다면 _그것이 정의된 패키지 밖에서는_ private 하다.

우리의 예제에서 우리의 메서드만 이 변수를 조작할 수 있도록 하고 다른 것은 조작하지 못하도록 하길 원한다.

우리는 내부 `balance` 필드에 구조체 안에 있는 "receiver" 변수를 통해서만 접근할 수 있다는 것을 기억 해야 한다.

```go
func (w Wallet) Deposit(amount int) {
    w.balance += amount
}

func (w Wallet) Balance() int {
    return w.balance
}
```

핀테크에서의 경력이 유지될 수 있도록 보안에 주의한 뒤, 테스트를 실행하고 통과하는 테스트를 즐기도록 하자.

`wallet_test.go:15: got 0 want 10`

### ????

혼란스럽게도, 우리의 코드는 제대로 작동하는 것처럼 보일 수 있는데, 우리가 잔액에 새로운 비용을 추가한다면 위 balance 메서드는 현재 잔액의 상태를 반환해야 한다.

Go에서는, **함수나 메서드를 호출하는 경우 인자(arguments) 는** _**복사된다**_.

다음 함수를 호출 할 때 `func (w Wallet) Deposit(amount int)`에 `w`는 메서드를 호출하는 것의 복사본이다.

너무 컴퓨터 과학적으로 깊게 들어가지 않고 설명하면, 당신이 값을 생성하면 - wallet 같이, 그것은 메모리 어딘가에 저장된다. 당신은 `&myVal`와 같은 방식으로 값의 메모리 주소 bit를 찾을 수 있다.

당신의 코드에 프린트문을 추가하여 실험해보자.

```go
func TestWallet(t *testing.T) {

    wallet := Wallet{}

    wallet.Deposit(10)

    got := wallet.Balance()

    fmt.Printf("address of balance in test is %v \n", &wallet.balance)

    want := 10

    if got != want {
        t.Errorf("got %d want %d", got, want)
    }
}
```

```go
func (w Wallet) Deposit(amount int) {
    fmt.Printf("address of balance in Deposit is %v \n", &w.balance)
    w.balance += amount
}
```

이스케이프 문자 `\n`은 , 메모리 주소를 아웃풋으로 출력한 뒤 줄 바꿈을 해준다. `&`라는 심볼의 주소로 어떤 것에 대한 포인터를 얻는다.

이제 테스트를 새로 실행시킨다.

```text
address of balance in Deposit is 0xc420012268
address of balance in test is 0xc420012260
```

두 잔액의 주솟 값이 다른 것을 확인 할 수 있다. 따라서 우리가 만약 코드 안에서 잔액 값을 바꾸어 주는 것은 테스트로부터 받은 복사본에 작업하는 것이다. 결국 테스트에서는 잔액 값이 변화하지 않는다.

우리는 이것을 _포인터_ 로 해결 할 수 있다. [포인터](https://gobyexample.com/pointers)는 특정한 값을 _가리키고_ 따라서 그 값을 변화시킬 수 있다. 따라서 Wallet의 복사본을 갖지 않고, 우리는 wallet을 가리키는 포인터를 얻게 되어 값을 바꿀 수 있다.

```go
func (w *Wallet) Deposit(amount int) {
    w.balance += amount
}

func (w *Wallet) Balance() int {
    return w.balance
}
```

리시버 타입의 차이는 `Wallet`이 아니라 `*Wallet`이라고 쓰고 이것은 "wallet에 대한 포인터"라고 얘기할 수 있다.

새로 테스트를 재실행해 보면 통과할 것이다.

이게 왜 통과했지? 라고 의문을 가지고 아래처럼 우리는 포인터의 역참조(dereference)를 함수에서 사용해야 하는 것 아닌가 생각할 수 있다.

```go
func (w *Wallet) Balance() int {
    return (*w).balance
}
```

하지만 우린 객체를 직접적으로 접근하여 다룬 것처럼 보인다. 사실, 위의 `(*w)`를 사용한 코드는 완벽하게 타당하다. 그러나, Go 언어의 개발자들은 이 표기가 쓰기 귀찮은 것이라고 생각했고 그래서 Go에서는 특별한 역참조에 대한 명시 없이 `w.balance`라고 쓰는 것을 허용했다. 
이 구조체에 대한 포인터는 다음과 같이: _구조체 포인터_ 라고 불리고 [자동 역참조](https://golang.org/ref/spec#Method_values)가 된다.

기술적으로는 `Balance`라는 메서드는 포인터 리시버를 사용할 필요가 없고 balance의 복사본을 사용하여도 문제는 없다. 그러나 관습적으로 당신은 메서드 리시버 타입을 하나의 통일성 있게 가져가야 한다.

## 리팩터링 하기

우리는 비트코인 지갑을 만든다고 했지만 비트코인에 대한 언급은 지금까지 하지 않았다. 우리는 지금까지 `int`를 사용하였는데 그것은 무언가를 세는 데 있어서 좋은 타입이기 때문이다.

`구조체`를 추가로 생성해 사용하는 것은 좀 과하다고 생각할 수 있다. `int`만으로 동작하는 데는 문제가 없지만 그것을 잘 설명해주지 못하는 점이 있다.

기존의 존재하는 타입으로 새로운 타입을 만들어 주자.

문법은 다음과 같다 `type MyName OriginalType`

```go
type Bitcoin int

type Wallet struct {
    balance Bitcoin
}

func (w *Wallet) Deposit(amount Bitcoin) {
    w.balance += amount
}

func (w *Wallet) Balance() Bitcoin {
    return w.balance
}
```

```go
func TestWallet(t *testing.T) {

    wallet := Wallet{}

    wallet.Deposit(Bitcoin(10))

    got := wallet.Balance()

    want := Bitcoin(10)

    if got != want {
        t.Errorf("got %d want %d", got, want)
    }
}
```

`Bitcoin`을 만들기 위해서는 `Bitcoin(999)`와 같이 사용하면 된다.


이렇게 하므로 우리는 새로운 타입을 만들어 그 타입 위에 _메서드들_ 을 정의할 수 있다. 이것은 존재하는 타입에서 당신이 원하는 어떤 특정 도메인에 특화된 기능을 추가하는 경우 유용하다.

비트코인에 [Stringer](https://golang.org/pkg/fmt/#Stringer)를 구현해 보자.

```go
type Stringer interface {
        String() string
}
```
위 인터페이스는 `fmt` 패키지에 정의되어 있고 프린트에서 `%s` 포맷의 스트링을 사용하는 경우 당신의 타입이 어떻게 출력될지 정의한다.

```go
func (b Bitcoin) String() string {
    return fmt.Sprintf("%d BTC", b)
}
```

위에서 보이듯이, 타입 별칭(type alias)에서 새로운 메서드를 생성하는 문법과 구조체에서의 경우가 똑같은 것을 알 수 있다.

다음은 우리의 테스트에서 `String()`을 사용하도록 포맷 스트링을 바꿔준다.

```go
    if got != want {
        t.Errorf("got %s want %s", got, want)
    }
```
이것의 결과를 보기 위해, 일부러 테스트를 실패하도록 하면 아래의 결과를 확인 할 수 있다.

`wallet_test.go:18: got 10 BTC want 20 BTC`

이것은 우리의 테스트가 어떻게 진행되고 있는지 명확하게 보여준다.

다음 요구사항은 `Withdraw` 함수를 위한 것이다.

## 테스트부터 작성하기

`Deposit()`에 대부분 반대로 적용한다.

```go
func TestWallet(t *testing.T) {

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}

        wallet.Deposit(Bitcoin(10))

        got := wallet.Balance()

        want := Bitcoin(10)

        if got != want {
            t.Errorf("got %s want %s", got, want)
        }
    })

    t.Run("Withdraw", func(t *testing.T) {
        wallet := Wallet{balance: Bitcoin(20)}

        wallet.Withdraw(Bitcoin(10))

        got := wallet.Balance()

        want := Bitcoin(10)

        if got != want {
            t.Errorf("got %s want %s", got, want)
        }
    })
}
```

## 테스트 실행해보기

`./wallet_test.go:26:9: wallet.Withdraw undefined (type Wallet has no field or method Withdraw)`

## 컴파일이 되는 최소한의 코드를 작성하고, 테스트 실패 출력을 확인하기 

```go
func (w *Wallet) Withdraw(amount Bitcoin) {

}
```

`wallet_test.go:33: got 20 BTC want 10 BTC`

## 테스트를 통과하는 최소한의 코드 작성하기 

```go
func (w *Wallet) Withdraw(amount Bitcoin) {
    w.balance -= amount
}
```

## 리팩터링 하기

우리의 테스트에 중복이 있기 때문에, 중복을 리팩토링하여 제거한다.

```go
func TestWallet(t *testing.T) {

    assertBalance := func(t testing.TB, wallet Wallet, want Bitcoin) {
        t.Helper()
        got := wallet.Balance()

        if got != want {
            t.Errorf("got %s want %s", got, want)
        }
    }

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}
        wallet.Deposit(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw", func(t *testing.T) {
        wallet := Wallet{balance: Bitcoin(20)}
        wallet.Withdraw(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

}
```

만약 `Withdraw`를 계좌에 남아있는 잔액보다 많이 시도하게 된다면 어떻게 될까?  지금까지는, 우리의 요구사항은 초과 인출 시설에 대해서는 가정하지 않았다.

`Withdraw`를 사용하다 문제가 생긴 경우 우리는 어떻게 알려야 할까?

만약 에러를 알려주길 원한다면 Go에서는 관용적으로 함수에서 리턴 값으로 `err`를 보내주어 호출자가 확인하고 행동 할 수 있도록 해준다. 

테스트에서 이것을 시도해 보자.

## 테스트부터 작성하기

```go
t.Run("Withdraw insufficient funds", func(t *testing.T) {
    startingBalance := Bitcoin(20)
    wallet := Wallet{startingBalance}
    err := wallet.Withdraw(Bitcoin(100))

    assertBalance(t, wallet, startingBalance)

    if err == nil {
        t.Error("wanted an error but didn't get one")
    }
})
```

_만약_ 기존의 잔액보다 더 많이 인출을 시도한다면 잔액은 기존과 같게 유지하고 `Withdraw`에서는 에러를 리턴하도록 해야한다.

그래서 우리는 만약 `nil`이 리턴 된다면 테스트가 실패하도록 하여 에러를 체크할 것이다.

`nil`은 다른 프로그래밍 언어에서의 `null`과 동의어다. 에러는 `nil`이 될 수 있는데, 그 이유는 `Withdraw`의 리턴 타입이 `error`이고 이것은 인터페이스이기 때문이다. 만약 인터페이스를 인자나 리턴 값으로 받는 함수를 보게 되면 이것은 `nil`이 될 수 있다(nillable).  

`null`처럼 만약 `nil` 값에 접근하려 하면 런타임 패닉을 던지게 됩니다. 이것은 매우 안 좋다! 반드시 nil인지 확인해야 한다.

## 테스트 실행해보기

`./wallet_test.go:31:25: wallet.Withdraw(Bitcoin(100)) used as value`

위의 말이 아마 좀 확실하지 않아 보일 수 있지만, 이전의 `Withdraw`의 의도는 단지 호출하는 것이였고, 값을 리턴하지 않았다. 컴파일이 되게 하기 위해서는 이 메서드가 리턴 타입을 가지도록 해주어야 한다.

## 컴파일이 되는 최소한의 코드를 작성하고, 테스트 실패 출력을 확인하기 

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {
    w.balance -= amount
    return nil
}
```

다시, 단지 컴파일러를 만족시키는 적당한 코드를 작성하는 것이 매우 중요하다. `Withdraw` 메서드를 `error`를 리턴하도록 수정하고 지금부터는 _어떤 것_ 을 리턴해야하기 때문에 일단 `nil`을 리턴하도록 해보자.

## 테스트를 통과하는 최소한의 코드 작성하기 

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

    if amount > w.balance {
        return errors.New("oh no")
    }

    w.balance -= amount
    return nil
}
```

코드에서 `errors`를 import 해주는 것을 기억해야 한다.

`errors.New`는 당신이 작성한 메시지와 함께 새로운 `error`를 생성하여 준다.

## 리팩터링 하기

에러 체크를 하는 데 있어 테스트를 좀 더 명확하게 읽을 수 있도록 빠르게 테스트 헬퍼(helper)를 만들어 준다.

```go
assertError := func(t testing.TB, err error) {
    t.Helper()
    if err == nil {
        t.Error("wanted an error but didn't get one")
    }
}
```

그리고 우리의 테스트에서 

```go
t.Run("Withdraw insufficient funds", func(t *testing.T) {
    startingBalance := Bitcoin(20)
    wallet := Wallet{startingBalance}
    err := wallet.Withdraw(Bitcoin(100))

    assertBalance(t, wallet, startingBalance)
    assertError(t, err)
})
```

"oh no"라는 에러를 리턴하는 것은 별로 유용하지 않기에 우리는 계속 에러를 반복할 수 있다는 것을 생각해야 한다. 

에러가 궁극적으로 유저에게 전달된다고 가정하면, 단지 에러가 존재하게 두는 것보다는 테스트에서 어떤 종류의 메시지라도 assert 하도록 개선해야 한다.

## 테스트부터 작성하기

helper에서 `string`을 비교하도록 업데이트한다.

```go
assertError := func(t testing.TB, got error, want string) {
    t.Helper()
    if got == nil {
        t.Fatal("didn't get an error but wanted one")
    }

    if got.Error() != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

그다음 호출자를 업데이트한다.

```go
t.Run("Withdraw insufficient funds", func(t *testing.T) {
    startingBalance := Bitcoin(20)
    wallet := Wallet{startingBalance}
    err := wallet.Withdraw(Bitcoin(100))

    assertBalance(t, wallet, startingBalance)
    assertError(t, err, "cannot withdraw, insufficient funds")
})
```

우리는 `t.Fatal`를 사용하였는데 이것은 불리게 된다면 테스트를 중지한다. 주위에 하나도 없는 것이 아니라면, 반환된 오류에 대해 더는 assertions이 일어나게 하고 싶지 않기 때문이다. 이것이 없다면 테스트는 다음 스텝으로 계속 진행되고 nil 포인터에 의해 패닉이 일어난다.

## 테스트 실행해보기

`wallet_test.go:61: got err 'oh no' want 'cannot withdraw, insufficient funds'`

## 테스트를 통과하는 최소한의 코드 작성하기 

```go
func (w *Wallet) Withdraw(amount Bitcoin) error {

    if amount > w.balance {
        return errors.New("cannot withdraw, insufficient funds")
    }

    w.balance -= amount
    return nil
}
```

## 리팩터링 하기

테스트 코드와 `Withdraw` 코드 모두 에러 메시지를 포함하고 있어 중복이 있다.

누군가가 테스트 메시지 워딩을 바꾸길 원한다면 테스트를 실패하도록 하는 것은 매우 귀찮은 일이 될 것이고 워딩을 바꾸는 것은 테스트에서 너무 디테일한 부분이다. 우리는 정확히 어떤 단어인지 _정말로_ 관심이 없고, 특정한 상황에서 인출을 하는 경우, 일종의 의미 있는 에러 메시지를 반환하여 주면 된다.

Go에서는 에러는 값이기 때문에, 우리는 에러를 변수로 리팩토링 할 수 있어 하나의 값으로 에러를 가지고 갈 수 있다.

```go
var ErrInsufficientFunds = errors.New("cannot withdraw, insufficient funds")

func (w *Wallet) Withdraw(amount Bitcoin) error {

    if amount > w.balance {
        return ErrInsufficientFunds
    }

    w.balance -= amount
    return nil
}
```

`var` 키워드는 패키지에서 전역으로 변수를 선언 할 수 있도록 허용한다. 

이제 우리의 `Withdraw` 함수는 매우 깔끔해졌기 때문에 이것은 그 자체로 매우 긍정적인 변화이다.

다음은 테스트 코드에서 특정한 스트링을 사용하는 대신 이 값을 사용하도록 리팩토링한다.

```go
func TestWallet(t *testing.T) {

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}
        wallet.Deposit(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw with funds", func(t *testing.T) {
        wallet := Wallet{Bitcoin(20)}
        wallet.Withdraw(Bitcoin(10))
        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw insufficient funds", func(t *testing.T) {
        wallet := Wallet{Bitcoin(20)}
        err := wallet.Withdraw(Bitcoin(100))

        assertBalance(t, wallet, Bitcoin(20))
        assertError(t, err, ErrInsufficientFunds)
    })
}

func assertBalance(t testing.TB, wallet Wallet, want Bitcoin) {
    t.Helper()
    got := wallet.Balance()

    if got != want {
        t.Errorf("got %q want %q", got, want)
    }
}

func assertError(t testing.TB, got error, want error) {
    t.Helper()
    if got == nil {
        t.Fatal("didn't get an error but wanted one")
    }

    if got != want {
        t.Errorf("got %q, want %q", got, want)
    }
}
```

그리고 이제 테스트는 매우 따라가기 쉬워졌다.

나는 그리고 헬퍼들(helpers)을 메인 테스트 함수에서 빼서 옮겼고, 따라서 다른 사람이 파일을 열었을 때 헬퍼들(helpers)을 먼저 읽기보다 assertions를 먼저 읽을 수 있도록 만들었다. 

테스트의 다른 특징으로는 테스트를 통하여 우리가 _실제_ 코드의 사용법을 이해하도록 도와주고 그래서 코드를 공감할 수 있도록 만들어준다. 여기서 보듯이 개발자는 간단히 우리의 코드를 호출하고 `ErrInsufficientFunds`와 동일한지 확인하고 적절히 행동하면 된다.

### 확인하지 않은 에러

Go 컴파일러가 많은 것을 도와주지만, 때때로 당신이 놓치고 에러 핸들링하기 쉽지 않은 것들이 있다. 

우리가 테스트하지 않은 하나의 시나리오가 있다. 그것을 찾기 위해, 터미널에서 다음과 같은 것을 쳐서 Go에서 이용 가능한 linter 중 하나인 `errcheck`를 설치하자 

`go get -u github.com/kisielk/errcheck`

그 뒤, 당신 코드의 디렉터리 안에서 `errcheck .`을 실행하자.

당신은 다음과 같은 것을 받을 것이다.

`wallet_test.go:17:18: wallet.Withdraw(Bitcoin(10))`

위에서 우리에게 말하고자 하는 것은 우리는 코드의 그 줄에서 반환되는 에러를 확인하지 않았다는 것이다. 이것은 내 컴퓨터에서 코드의 그 라인은 우리의 일반적인 인출 시나리오이고 우리는 `Withdraw`가 성공적인지, 즉 에러가 반환되지 _않았는지_ 확인하지 않았다는 것을 의미한다.

이것이 계좌를 위한 마지막 테스트 코드다.

```go
func TestWallet(t *testing.T) {

    t.Run("Deposit", func(t *testing.T) {
        wallet := Wallet{}
        wallet.Deposit(Bitcoin(10))

        assertBalance(t, wallet, Bitcoin(10))
    })

    t.Run("Withdraw with funds", func(t *testing.T) {
        wallet := Wallet{Bitcoin(20)}
        err := wallet.Withdraw(Bitcoin(10))

        assertBalance(t, wallet, Bitcoin(10))
        assertNoError(t, err)
    })

    t.Run("Withdraw insufficient funds", func(t *testing.T) {
        wallet := Wallet{Bitcoin(20)}
        err := wallet.Withdraw(Bitcoin(100))

        assertBalance(t, wallet, Bitcoin(20))
        assertError(t, err, ErrInsufficientFunds)
    })
}

func assertBalance(t testing.TB, wallet Wallet, want Bitcoin) {
    t.Helper()
    got := wallet.Balance()

    if got != want {
        t.Errorf("got %s want %s", got, want)
    }
}

func assertNoError(t testing.TB, got error) {
    t.Helper()
    if got != nil {
        t.Fatal("got an error but didn't want one")
    }
}

func assertError(t testing.TB, got error, want error) {
    t.Helper()
    if got == nil {
        t.Fatal("didn't get an error but wanted one")
    }

    if got != want {
        t.Errorf("got %s, want %s", got, want)
    }
}
```

## 정리

### 포인터

* Go는 함수/메서드에서 값을 넘겨줄 때 값을 복사하기 때문에 만약 함수에서 그 상태를 바꾸기를 원한다면 바꾸길 원하는 것의 포인터를 받아야 한다. 
* Go에서 값을 복사한다는 사실은 꽤 자주 유용하지만 때때로 당신의 시스템에서 어떤 것의 복사본을 만들기 원하지 않는다면 그 경우, 참조(reference)를 넘겨주어야 한다. 예를 들면, 매우 큰 데이터나 데이터베이스의 커넥션풀 같이 아마 당신이 하나의 인스턴스만 가지려 하는 것들일 수 있다.

### nil

* 포인터는 nil일 수 있다.
* 만약 함수가 어떤 것의 포인터를 반환하였다면 당신은 반드시 그것이 nil인지 아닌지 확인하거나 런타임 예외를 일으켜야 한다. 컴파일러는 이것에서 당신을 도와주지 않는다.
* 당신이 표현하려 하는 값이 없을 수도 있을 때 유용하다.

### 에러

* 에러는 함수/메서드를 호출할 때 실패를 알려주는 방법이다.
* 우리의 테스트 과정을 본다면, 에러에 스트링을 사용하여 체크하는 방법은 매우 유별난 테스트(flaky test)가 된다고 결론을 내렸다. 따라서 우리는 그 대신 의미 있는 값으로 리팩토링하여 테스트를 더 쉽게 할 수 있고 이것을 사용하면 API의 사용자도 더 쉬워질 것이다. 
* 이것은 에러 처리의 끝이 아니며, 당신은 좀 더 복잡한 것을 할 수 있고 이것은 단지 시작이다. 이후 섹션에서 더 많은 전략을 다룰 것이다.
* [에러를 체크하지 마라, 에러를 우아하게 다뤄라](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)

### 기존 타입으로부터 새로운 타입 생성

* 값에 특정 도메인에 의미를 추가하는 데 유용하다.
* 인터페이스를 구현할 수 있도록 한다. 

포인터와 에러는 Go를 작성하는 데 있어서 매우 큰 부분이며 당신은 이것들에 익숙해져야 한다. 만약 당신이 실수하더라도 고맙게도 컴파일러가 _보통_ 문제가 생긴 부분을 도와주기 때문에 시간을 들여 그 에러를 읽어보도록 하자.
