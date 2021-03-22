# 테스트와 함께 Go 배우기

<p align="center">
  <img src="../red-green-blue-gophers-smaller.png" />
</p>

[Art by Denise](https://twitter.com/deniseyu21)

[![Build Status](https://travis-ci.org/quii/learn-go-with-tests.svg?branch=main)](https://travis-ci.org/quii/learn-go-with-tests)
[![Go Report Card](https://goreportcard.com/badge/github.com/quii/learn-go-with-tests)](https://goreportcard.com/report/github.com/quii/learn-go-with-tests)

## 번역

- 원본 : [english](https://quii.gitbook.io/learn-go-with-tests/)

- [한국어](https://miryang.gitbook.io/learn-go-with-tests)
- [中文](https://studygolang.gitbook.io/learn-go-with-tests)
- [Português](https://larien.gitbook.io/aprenda-go-com-testes/)
- [日本語](https://andmorefine.gitbook.io/learn-go-with-tests/)

## 왜

* 테스트를 작성하며 Go언어를 배우기
* **TDD(테스트 주도 개발)의 기반 다지기**. Go는 학습하는 것이 쉽고, 테스트 기능이 내장되어 있어 TDD를 배우기에 좋은 언어이다.
* 튼튼하고 충분히 테스트 된 시스템을 Go언어로 작성하게 될 것이라고 확신한다.
* [비디오를 보거나 유닛 테스트와 TDD가 중요한 이유에 대해 읽어보십시오.](why.md)

## 목차

### Go 기초

1. [Go 설치](install-go.md) - 생산성을 위한 환경 준비
2. [Hello, world](hello-world.md) - 변수 선언, 상수, if/else 조건문, switch, 첫 go 프로그램과 첫 테스트를 작성, 서브테스트 구문과 클로저
3. [정수](integers.md) - 함수 선언 구문의 자세한 내용과 코드 문서를 개선하는 새로운 방법 학습
4. [반복](iteration.md) - `for` 학습과 벤치마킹
5. [배열과 슬라이스](arrays-and-slices.md) - 배열, 슬라이스, `len`, 인자, `range` 학습 및 테스트 커버리지
6. [구조체, 메소드 & 인터페이스](structs-methods-and-interfaces.md) - `struct`, 메소드, `interface` 학습 및 테이블 기반 테스트
7. [포인터 & 에러](pointers-and-errors.md) - 포인터와 에러 학습
8. [맵](maps.md) - 맵 데이터 구조체에 값을 저장하는 방법 학습
9. [의존성 주입](dependency-injection.md) - 의존성 주입, 인터페이스 사용과의 관계 및 io 기본에 대해 학습
10. [Mocking](mocking.md) - 테스트되지 않은 기존 코드를 가져와 mocking과 함께 의존성 주입을 사용하여 테스트
11. [동시성](concurrency.md) - 소프트웨어를 더 빠르게 만들기 위해 동시성 코드를 작성하는 방법 학습
12. [select](select.md) - 비동기식 프로세스를 우아하게 동기화하는 방법 학습
13. [리플렉션](reflection.md) - 리플렉션 학습
13. [sync](sync.md) - `WaitGroup` 과 `Mutex` 를 포함한 sync 패키지의 일부 기능 학습
13. [Context](context.md) - context 패키지를 사용하여 장기 실행 프로세스 관리 및 취소
14. [속성 기반 테스트](roman-numerals.md) - Roman Numerals kata로 TDD를 연습하고, 속성 기반 테스트에 대한 간단한 소개
15. [Maths](math.md) - `math` 패키지를 사용하여 SVG 시계 그리기

### 어플리케이션 만들기

이제 _Go 기초_ 부분을 잘 소화했으며, 언어의 주요 기능과 TDD 작동 방식에 대한 탄탄한 기반이 마련되었다.

이번 섹션에는 어플리케이션 개발이 진행된다.

각 장은 이전 장에서 이어지며, 제품 소유자의 지시에 따라 어플리케이션의 기능을 확장한다.

좋은 코드를 작성하기 위해 새로운 개념들이 소개될 예정이지만, 대부분은 Go 표준 라이브러리로 수행할 수 있다.

이 과정을 끝내면, 테스트와 함께 Go 어플리케이션을 효과적으로 작성하는 방법을 잘 이해할 수 있다.

* [HTTP 서버](http-server.md) - HTTP 요청을 받고, 응답을 보내는 어플리케이션 생성
* [JSON, routing 및 embedding](json.md) - 엔드포인트에 JSON을 반환하고, 라우팅하는 방법
* [IO 및 sorting](io.md) - 디스크에서 데이터를 읽고, 데이터 정렬 다루기
* [Command line & project structure](command-line.md) - 하나의 코드 베이스에서 여러 어플리케이션을 보조하고, 커맨드 라인에서 입력 읽기
* [Time](time.md) - 스케쥴 작업들을 위해 `time` 패키지 사용
* [웹소켓](websockets.md) - 웹소켓을 사용하는 서버를 만들고 테스트하는 방법 학습

### 질문과 답

* [OS exec](os-exec.md) - OS 커맨트 명령을 실행하여 데이터를 가져오고, 로직 테스트가 가능한 비즈니스를 유지하는 방법 예제
* [Error types](error-types.md) - 테스트를 개선하고 코드를 사용하기 쉽게 할 수 있도록 자신만의 에러 타입을 만드는 방법 예제
* [Context-aware Reader](context-aware-reader.md) - Learn how to TDD augmenting `io.Reader` with cancellation. Based on [Context-aware io.Reader for Go](https://pace.dev/blog/2020/02/03/context-aware-ioreader-for-golang-by-mat-ryer)
* [Revisiting HTTP Handlers](http-handlers-revisited.md) - HTTP 핸들러를 테스트하는 것은 많은 개발자의 골칫거리로 보인다. 이 챕터에서는 핸들러를 올바르게 설계하는 사안에 대해서 살펴본다.

### Meta / Discussion

* [Why](why.md) - 비디오를 보거나 유닛 테스트와 TDD가 중요한 이유에 대해 읽어보십시오.
* [Intro to generics](intro-to-generics.md) - 제네릭 인수를 사용해 함수를 작성하고, 자신만의 제네릭 구조체(structure)를 만드는 방법 학습
* [Anti-patterns](anti-patterns.md) - TDD와 유닛 테스트 anti-patterns에 관한 짧은 챕터

## 기여

* _이 프로젝트는 진행 중이다._ 만약 기여하고 싶다면, 연락하십시오.
* [contributing.md](https://github.com/quii/learn-go-with-tests/tree/842f4f24d1f1c20ba3bb23cbc376c7ca6f7ca79a/contributing.md) 를 읽으십시오.
* 아이디어가 있다면 이슈를 등록하십시오.

## 누구를 위해

* Go 학습에 관심이 있는 사람들
* 이미 Go를 알고 있지만, TDD로 테스팅을 학습하고 싶은 사람들

## 필요한 것

* 컴퓨터!
* [Go 설치](https://golang.org/)
* 에디터
* 프로그래밍 경험. `if`, 변수, 함수 등을 이해할 수 있는지
* 터미널에 익숙한지

## 번역

번역 작업은 [여기](https://github.com/MiryangJung/learn-go-with-tests-ko) 에서 이루어지고 있습니다.

* 수정이 필요하거나 번역에 참여하고 싶다면 [여기](https://github.com/MiryangJung/learn-go-with-tests-ko) 에 PR을 보내주세요.

### 번역자

- [MiryangJung](https://github.com/miryangjung)
- [hjjg200](https://github.com/hjjg200)
- [MinsuKim](https://github.com/alstn2468)
- [HyunseokJeong](https://github.com/nicewook)
- [JunpyoChang](https://github.com/changjunpyo)
- [JaehoonKim](https://github.com/jaehoonn)
- [YoungjunChoi](https://github.com/WHYjun)
- [Seongbin Cho](https://github.com/chobobdev)
- [SangjunCha](https://github.com/SangjunCha-dev)
- [hywn](https://github.com/14hhan)


## 피드백

* 이슈를 등록하거나 PR를 보내세요. [여기](https://github.com/quii/learn-go-with-tests) 또는 [tweet me @quii](https://twitter.com/quii)

[MIT license](LICENSE.md)

[Logo is by egonelbre](https://github.com/egonelbre) What a star!

