# 구조체, 메서드 & 인터페이스 (Structs, methods & interfaces)

**[이 챕터에서 사용되는 모든 코드는 여기서 찾을 수 있다.](https://github.com/quii/learn-go-with-tests/tree/main/structs)**

높이와 너비가 주어진 사각형의 둘레를 계산하는 기하학 코드가 필요하다고 가정한다. `Perimeter(width float64, height float64)` 함수를 작성할 수 있다. 여기서 `float64`는 `123.45` 같은 부동 소수점 수에 대한 것이다.

지금쯤이면 TDD 주기가 꽤 익숙할 것이다.

## 테스트부터 작성하기

```go
func TestPerimeter(t *testing.T) {
    got := Perimeter(10.0, 10.0)
    want := 40.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}
```

새로운 형식의 문자열을 보면 `f`는 `float64`에 대한 것이고, `.2`는 소수점 2자리 출력을 의미한다.

## 테스트 실행해보기

`./shapes_test.go:6:9: undefined: Perimeter`

## 컴파일이 되는 최소한의 코드를 작성하고, 테스트 실패 출력을 확인하기

```go
func Perimeter(width float64, height float64) float64 {
    return 0
}
```

결과 `shapes_test.go:10: got 0.00 want 40.00` 가 출력된다.

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func Perimeter(width float64, height float64) float64 {
    return 2 * (width + height)
}
```

지금까지는 정말 쉽다. 이제 직사각형의 면적을 반환하는 `Area(width, height float64)` 함수를 만들어 본다.

TDD 주기에 따라 직접 작성한다.

이와 같은 코드로 테스트가 끝나야 한다.

```go
func TestPerimeter(t *testing.T) {
    got := Perimeter(10.0, 10.0)
    want := 40.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}

func TestArea(t *testing.T) {
    got := Area(12.0, 6.0)
    want := 72.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}
```

```go
func Perimeter(width float64, height float64) float64 {
    return 2 * (width + height)
}

func Area(width float64, height float64) float64 {
    return width * height
}
```

## 리팩터링 하기

코드는 제대로 작동하지만, 직사각형에 대한 명시적인 내용이 없다. 부주의한 개발자가 삼각형의 너비와 높이를 이러한 함수에 사용 할 수 있는데, 함수가 잘못된 값을 반환할 것이다.

`RectangleArea`와 같이 기능을 좀 더 구체적으로 지정할 수 있다. 더 나은 해결책은 이 개념을 캡슐화하는 `Rectangle`이라고 불리는 자기만의 _type_ 을 정의하는 것이다.

**struct**를 사용해서 간단한 유형을 만들 수 있다. [A struct](https://golang.org/ref/spec#Struct_types)는 데이터를 저장할 수 있는 명명된 필드의 집합이다.

이와 같이 struct를 선언한다.

```go
type Rectangle struct {
    Width float64
    Height float64
}
```

이제 일반 `float64` 대신 `Rectangle`을 사용하도록 코드를 리팩터링한다.

```go
func TestPerimeter(t *testing.T) {
    rectangle := Rectangle{10.0, 10.0}
    got := Perimeter(rectangle)
    want := 40.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}

func TestArea(t *testing.T) {
    rectangle := Rectangle{12.0, 6.0}
    got := Area(rectangle)
    want := 72.0

    if got != want {
        t.Errorf("got %.2f want %.2f", got, want)
    }
}
```

코드를 고치기 전에 테스트를 실행하고, 다음과 같은 유용한 에러를 얻어야 한다.

```text
./shapes_test.go:7:18: not enough arguments in call to Perimeter
    have (Rectangle)
    want (float64, float64)
```

`myStruct.field` 구문을 사용하여 구조체 필드에 접근할 수 있다.

두 가지 함수를 변경하여 테스트 코드를 수정한다.

```go
func Perimeter(rectangle Rectangle) float64 {
    return 2 * (rectangle.Width + rectangle.Height)
}

func Area(rectangle Rectangle) float64 {
    return rectangle.Width * rectangle.Height
}
```

함수에 `Rectangle`을 전달하는 것이 의도에 더 명확하게 전달하지만 앞으로 배워가는 구조체를 사용하는 것이 더 많은 이점이 있다는것에 동의해주길 바란다.

다음 필요조건은 원에 대한 `Area`함수를 작성하는 것이다. 

## 테스트부터 작성하기

```go
func TestArea(t *testing.T) {

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        got := Area(rectangle)
        want := 72.0

        if got != want {
            t.Errorf("got %g want %g", got, want)
        }
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        got := Area(circle)
        want := 314.1592653589793

        if got != want {
            t.Errorf("got %g want %g", got, want)
        }
    })

}
```

보다시피, 'f'가 'g'로 대체되었다. 'f'를 사용하면 정확한 십진수를 알 수가 없고, 'g'를 사용하면 오류메시지에서 완전한 십진수를 얻을 수 있다 \([fmt options](https://golang.org/pkg/fmt/)\).

## 테스트 실행해보기

`./shapes_test.go:28:13: undefined: Circle`

## 컴파일이 되는 최소한의 코드를 작성하고, 테스트 실패 출력을 확인하기

`Circle` 타입을 정의할 필요가 있다.

```go
type Circle struct {
    Radius float64
}
```

이제 테스트를 다시 실행한다.

`./shapes_test.go:29:14: cannot use circle (type Circle) as type Rectangle in argument to Area`

일부 프로그래밍 언어를 사용하면 다음과 같은 코드를 사용할 수 있다.

```go
func Area(circle Circle) float64 { ... }
func Area(rectangle Rectangle) float64 { ... }
```

하지만 Go에서는 사용할 수 없다.

`./shapes.go:20:32: Area redeclared in this block`

두 가지 방법이 있다.

* 동일한 이름의 함수를 다른 _packages_ 로 선언할 수 있다. 그래서 새로운 패키지로 `Area(Circle)`을 만들 수 있지만, 여기서는 너무 과한 느낌이다.
* 대신 새로 정의된 유형을 [_methods_](https://golang.org/ref/spec#Method_declarations) 정의할 수 있다. 

### 메서드란 무엇인가?

지금까지 _functions_ 을 쓰고 몇가지 방법을 사용해 왔다. `t.Errorf`를 부를 때 `t` \(`testing.T`\)의 인스턴스에서 `Errorf`메서드라고 한다.

메서드는 수신기가 있는 함수이다. 메서드 선언은 메서드 이름인 식별자를 메서드에 바인딩하고 메서드를 수신기의 기본 유형과 연결한다.

메서드는 함수와 매우 유사하지만 특정 유형의 인스턴스에서 메서드를 호출해야 호출된다.

`Area(rectangle)`과 같이 원하는 곳 어디서나 기능을 호출할 수 있는 경우에는 "things"에 대해서만 메서드를 호출할 수 있다.

예를 들어보면 도움이 되므로 먼저 테스트를 변경하여 메서드를 호출한 다음 코드를 수정한다.

```go
func TestArea(t *testing.T) {

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        got := rectangle.Area()
        want := 72.0

        if got != want {
            t.Errorf("got %g want %g", got, want)
        }
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        got := circle.Area()
        want := 314.1592653589793

        if got != want {
            t.Errorf("got %g want %g", got, want)
        }
    })

}
```

만약 테스트를 실행한다면

```text
./shapes_test.go:19:19: rectangle.Area undefined (type Rectangle has no field or method Area)
./shapes_test.go:29:16: circle.Area undefined (type Circle has no field or method Area)
```

> type Circle has no field or method Area

여기서 컴파일러가 얼마나 훌륭한지 다시 한번 강조하고 싶다. 오류 메시지를 천천히 읽는 것은 매우 중요하기 때문에 장기적으로 도움이 될 것이다.

## 컴파일이 되는 최소한의 코드를 작성하고, 테스트 실패 출력을 확인하기

유형에 몇가지 메서드를 추가하겠다.

```go
type Rectangle struct {
    Width  float64
    Height float64
}

func (r Rectangle) Area() float64  {
    return 0
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64  {
    return 0
}
```

메서드를 선언하는 구문은 함수와 거의 동일하며, 이는 메서드가 함수와 정말 유사하기 때문이다. 유일한 차이점은 메서드 수신기 `func (receiverName ReceiverType) MethodName(args)`의 구문이다.

메서드가 해당 유형의 변수에 호출되면 `receiverName` 변수를 통해 해당 데이터에 대한 참조를 얻는다. 다른 많은 프로그래밍 언어에서 이것은 암시적으로 수행되며 `receiverName`을 통해 수신기에 접근한다.

수신자 변수를 유형의 첫 번째 문자로 지정하는 것이 Go의 관례이다.

```go
r Rectangle
```

테스트를 다시 실행하려고 하면 컴파일 오류가 발생한 부분을 제공해야 한다.

## 테스트를 통과하는 최소한의 코드 작성하기

이제 새로운 방법으로 수정하여 직사각형 테스트를 통과하겠다.

```go
func (r Rectangle) Area() float64  {
    return r.Width * r.Height
}
```

테스트를 다시 실행하는 경우 직사각형 테스트는 통과 하지만 원 테스트는 여전히 실패한다.

원의 `Area` 함수가 통과하도록 만들기 위해 `math` 패키지에서 `Pi` 상수를 빌려온다. \(import 잊지 않아야 한다.\).

```go
func (c Circle) Area() float64  {
    return math.Pi * c.Radius * c.Radius
}
```

## 리팩터링 하기

테스트에는 몇가지 중복된 것이 있다.

원하는 것은 도형 모음을 가져와서 도형에 대한 `Area()` 메서드를 호출한 다음 결과를 확인하는 것이다.

`Rectangle`과 `Circle` 테스트를 통과할 수 있지만 도형이 아닌 것을 전달하려고 하면 컴파일하지 못하는 일종의 `checkArea` 기능으로 쓸 수 있기를 바란다.

Go를 사용하면 **interfaces**를 사용하여 이 의도를 코드화할 수 있다.

[Interfaces](https://golang.org/ref/spec#Interface_types) Go와 같이 정적 형식의 언어에서는 매우 강력한 개념으로, 다양한 유형과 함께 사용할 수 있는 함수를 만들고, 여전히 유형 안전성을 유지하면서 고도로 세분화된 코드를 만들 수 있기 때문이다.

테스트 내용을 리팩터링해서 소개하겠다.

```go
func TestArea(t *testing.T) {

    checkArea := func(t testing.TB, shape Shape, want float64) {
        t.Helper()
        got := shape.Area()
        if got != want {
            t.Errorf("got %g want %g", got, want)
        }
    }

    t.Run("rectangles", func(t *testing.T) {
        rectangle := Rectangle{12, 6}
        checkArea(t, rectangle, 72.0)
    })

    t.Run("circles", func(t *testing.T) {
        circle := Circle{10}
        checkArea(t, circle, 314.1592653589793)
    })

}
```

다른 연습문제처럼 도우미 기능을 만들고 있는데 이번에는 도형을 반환하라고 한다. 만약 도형이 아닌 것으로 부르려고 한다면, 컴파일하지 못할 것이다.

어떤 것이 도형이 되는가? 인터페이스 선언을 사용하는 `Shape`가 무엇인지 Go에 알려주기만 하면 된다.

```go
type Shape interface {
    Area() float64
}
```

`Rectangle`과 `Circle`을 만들었던 것처럼 새로운 `type`을 만들고 있지만 이번에는 `struct`가 아닌 `interface` 이다.

이것을 코드에 추가하면 테스트는 통과된다.

### 무엇을 기다리나요?

Go의 `interface`는 대부분의 다른 프로그래밍 언어의 인터페이스와는 상당히 다르다. 보통은 `My type Foo implements interface Bar` 라는 코드를 작성해야 한다.

하지만 우리의 경우에는

* `Rectangle`은 `Area` 메서드를 호출하여 `float64`를 반환하므로 `Shape` 인터페이스를 만족시킨다.
* `Circle`은 `Area` 메서드를 호출하여 `float64`를 반환하므로 `Shape` 인터페이스를 만족시킨다.
* `string`에 해당하는 메서드가 없으므로 인터페이스를 만족하지 않는다.
* 기타 등등

Go에서 **인터페이스 자료형은 암시적** 이다. 전달하는 유형이 인터페이스가 요청하는 유형과 일치하면 컴파일 된다.

### 디커플링(Decoupling)

도우미가 그 도형이 `Rectangle`인지 `Circle`인지 `Triangle`의 도형인지에 대해 어떻게 신경 쓸 필요가 없는지 주목한다. 인터페이스를 선언함으로써 도우미는 구체적인 유형으로부터 _분리_ 되고 단지 그 일을 하는데 필요한 방법을 갖게 된다.

**필요한 것만** 선언하는 인터페이스를 사용하는 접근 방식은 소프트웨어 설계에서 매우 중요하며, 이후 섹션에서 자세히 설명하겠다.

## 추가 리팩터링 하기

이제 구조체에 대해 어느 정도 이해했으므로 "테이블 기반 테스트"를 소개하겠다.

[테이블 기반 테스트](https://github.com/golang/go/wiki/TableDrivenTests) 동일한 방법으로 테스트할 수 있는 테스트 사례 목록을 작성하는 경우 유용하다.

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        shape Shape
        want  float64
    }{
        {Rectangle{12, 6}, 72.0},
        {Circle{10}, 314.1592653589793},
    }

    for _, tt := range areaTests {
        got := tt.shape.Area()
        if got != tt.want {
            t.Errorf("got %g want %g", got, tt.want)
        }
    }

}
```

여기서 유일하게 새로운 구문은 "익명 구조" areaTests를 만드는 것이다. `shape`와 `want`라는 두 개의 필드가 있는 `[]struct`를 사용하여 구조체를 선언한다. 그런 다음 슬라이스를 케이스로 채운다.

다른 슬라이스처럼 테스트를 실행하기 위해 구조체 필드를 사용하여 반복한다.

개발자가 새로운 도형을 도입하고 `Area`를 구현한 후 테스트 케이스에 추가하는 것은 쉽다. 게다가, 만약 `Area`에 있는 버그가 발견된다면, 버그를 고치기 전에 새로운 테스트 케이스를 추가하여 테스트하는 것이 매우 쉽다.

테이블 기반 테스트 도구 모음에서 좋은 항목이 될 수 있지만 테스트에서 extra noise가 필요한지 확인해야 한다.

인터페이스의 다양한 구현을 테스트하려는 경우 또는 함수에 전달된 데이터에 테스트가 필요한 여러 가지 요구 사항이 있는 경우 적합한 항목이다.

이 모든 것을 다른 도형, 즉 삼각형을 추가해서 테스트 한다.

## 테스트부터 작성하기

새로운 도형을 위한 새로운 테스트를 추가하는 것은 매우 쉽다. 코드 `{Triangle{12, 6}, 36.0},`를 추가하기만 하면 된다.

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        shape Shape
        want  float64
    }{
        {Rectangle{12, 6}, 72.0},
        {Circle{10}, 314.1592653589793},
        {Triangle{12, 6}, 36.0},
    }

    for _, tt := range areaTests {
        got := tt.shape.Area()
        if got != tt.want {
            t.Errorf("got %g want %g", got, tt.want)
        }
    }

}
```

## 테스트 실행해보기

테스트를 계속 실행하여 컴파일러가 솔루션을 안내하도록 한다.

## 컴파일이 되는 최소한의 코드를 작성하고, 테스트 실패 출력을 확인하기

`./shapes_test.go:25:4: undefined: Triangle`

아직 Triangle을 정의하지 않았다.

```go
type Triangle struct {
    Base   float64
    Height float64
}
```

다시 테스트를 시작한다.

```text
./shapes_test.go:25:8: cannot use Triangle literal (type Triangle) as type Shape in field value:
    Triangle does not implement Shape (missing Area method)
```

`Area()` 메서드에 Triangle이 없기 때문에 도형으로 사용할 수 없음을 알려주고 있으므로 빈 메서드을 추가하여 테스트한다.

```go
func (t Triangle) Area() float64 {
    return 0
}
```

마지막으로 코드가 컴파일되고 오류가 발생한다.

`shapes_test.go:31: got 0.00 want 36.00`

## 테스트를 통과하는 최소한의 코드 작성하기

```go
func (t Triangle) Area() float64 {
    return (t.Base * t.Height) * 0.5
}
```

그리고 테스트를 통과했다!

## 리팩터링 하기

다시 말하지만, 구현은 괜찮지만 테스트 코드는 약간의 개선을 할 수 있다.

값을 입력할때

```go
{Rectangle{12, 6}, 72.0},
{Circle{10}, 314.1592653589793},
{Triangle{12, 6}, 36.0},
```

모든 숫자가 무엇을 나타내는지는 바로 분명하지 않으며 쉽게 이해될 수 있는 테스트를 목표로 해야 한다.

지금까지 `MyStruct{val1, val2}` 구조의 인스턴스를 생성하는 구문만 표시되었지만 선택적으로 필드 이름을 지정할 수 있다.

어떻게 구성되는지 본다.

```go
        {shape: Rectangle{Width: 12, Height: 6}, want: 72.0},
        {shape: Circle{Radius: 10}, want: 314.1592653589793},
        {shape: Triangle{Base: 12, Height: 6}, want: 36.0},
```

[예제별 테스트 기반 개발](https://g.co/kgs/yCzDLF) Kent Beck은 몇 가지 테스트를 반영하고 다음과 같이 주장한다.

> 테스트는 더 명확하게 말해준다, as if it were an assertion of truth, **작업의 연속이 아니다**

\(강조\)

이제 테스트는 \(최소한의 예제 목록\) shapes와 areas에 대한 올바른 결과(assertion of truth)를 만든다.

## 테스트 출력이 유용한지 확인하기

아까 `Triangle`을 실행하다가 실패한 테스트에서 `shapes_test.go:31: got 0.00 want 36.00`를 출력했다.

`Triangle`과 관련 있다는 것은 `Triangle`을 가지고 작업하고 있었기 때문에 알고 있다. 하지만 테이블에 있는 20개의 케이스 중 하나에서 버그가 시스템에서 발생한다면? 어떤 경우에 실패했는지 개발자가 어떻게 알 수 있을까? 이것은 개발자에게 좋은 경험이 아니다. 실제로 실패한 사례를 찾기 위해 수동으로 검토해야 한다.

오류 메시지를 `%#v got %g want %g`로 변경할 수 있습니다. `%#v` 형식 문자열은 필드 값이 있는 구조를 출력하여 개발자가 테스트 중인 속성을 한 눈에 볼 수 있도록 한다. 

테스트 사례의 가독성을 높이기 위해 `want`필드를 `hasArea`와 같이 좀 더 설명적인 것으로 바꿀 수 있다.

테이블 기반 테스트의 마지막 팁은 `t.Run`을 사용하고 테스트 케이스 이름을 지정하는 것이다. 

각각의 케이스를 `t.Run`으로 실행하면 케이스 이름을 출력하기 때문에 실패시 테스트 출력이 명확해진다.

```text
--- FAIL: TestArea (0.00s)
    --- FAIL: TestArea/Rectangle (0.00s)
        shapes_test.go:33: main.Rectangle{Width:12, Height:6} got 72.00 want 72.10
```

또한 `go test -run TestArea/Rectangle`를 사용하여 특정 테스트를 실행할 수 있다. 

다음은 위처럼 출력하는 최종 테스트 코드이다.

```go
func TestArea(t *testing.T) {

    areaTests := []struct {
        name    string
        shape   Shape
        hasArea float64
    }{
        {name: "Rectangle", shape: Rectangle{Width: 12, Height: 6}, hasArea: 72.0},
        {name: "Circle", shape: Circle{Radius: 10}, hasArea: 314.1592653589793},
        {name: "Triangle", shape: Triangle{Base: 12, Height: 6}, hasArea: 36.0},
    }

    for _, tt := range areaTests {
        // using tt.name from the case to use it as the `t.Run` test name
        t.Run(tt.name, func(t *testing.T) {
            got := tt.shape.Area()
            if got != tt.hasArea {
                t.Errorf("%#v got %g want %g", tt.shape, got, tt.hasArea)
            }
        })
    }

}
```

## 정리

기본적인 수학 문제에 대한 해결책을 반복하고 테스트에 의해 파생된 언어의 새로운 특징을 배우는 TDD 연습이었다. 

* 관련 데이터를 함께 묶고 코드의 의도를 명확히 할 수 있는 자신만의 데이터 유형을 만들기 위한 구조체 선언
* 다양한 유형 \([파라미터의 다형성](https://en.wikipedia.org/wiki/Parametric_polymorphism)\)에서 사용할 수 있는 함수를 정의할 수 있도록 인터페이스 선언
* 데이터 유형에 메서드를 추가하고 인터페이스를 구현하는 메서드 추가
* 테이블 기반 테스트를 통해 코드를 보다 명확히하고 확장 및 유지 관리하기 쉬움

구조체, 메서드 & 인터페이스는 중요한 장이다. 왜냐하면 이제 자신의 유형을 정의하기 시작했기 때문이다. Go와 같이 정적으로 입력된 언어에서는 이해하기 쉬운 소프트웨어를 만들고, 조립하고, 테스트하기 위해 자신만의 유형을 설계할 수 있는 능력이 필수적이다.

인터페이스는 시스템의 다른 부분으로부터 복잡성을 숨길 수 있는 훌륭한 도구다. 테스트 도우미 _code_ 는 정확한 도형을 알 필요가 없고 단지 그 영역을 "묻는" 방법만 알 필요가 있다.

Go에 익숙해지면 인터페이스와 표준 라이브러리의 실제 강점을 볼 수 있다. _어디에서나_ 사용되는 표준 라이브러리에 정의된 인터페이스에 대해 배울 수 있으며, 이러한 인터페이스를 자신의 유형에 맞게 구현하면 우수한 많은 기능을 빠르게 재사용할 수 있다.
