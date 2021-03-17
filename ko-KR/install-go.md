# Install Go, set up environment for productivity

Go에 대한 공식 설치 가이드를 사용할 수 있습니다. [here](https://golang.org/doc/install).

이 가이드에서는 사용자가 패키지 관리자를 사용하고 있다고 가정합니다. 예를들어 [Homebrew](https://brew.sh), [Chocolatey](https://chocolatey.org), [Apt](https://help.ubuntu.com/community/AptGet/Howto)와 [yum](https://access.redhat.com/solutions/9934)이 있습니다.

시연을 위하여 OSX에서 Homebrew를 이용한 설치 과정을 보여드리겠습니다.

## Installation

설치 프로세스는 매우 간단합니다. 일단, Homebrew를 설치하려면 이 명령을 실행해야 합니다. Homebrew는 Xcode에 종속되어 있으므로 Xcode를 먼저 설치해야 합니다.

```sh
xcode-select --install
```

그리고 다음과 같이 homebrew를 설치할 수 있습니다.:

```sh
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
```

이 시점에 Go도 다음과 같이 설치할 수 있습니다.:

```sh
brew install go
```

*패키지 관리자가 권장하는 모든 지침을 따라야 합니다. **Note**: OS 별로 권장사항은 다를 수 있습니다*.

다음과 같이 설치된 것을 확인할 수 있습니다.

```sh
$ go version
go version go1.14 darwin/amd64
```

## Go Environment

### $GOPATH
Go는 독선적입니다.

컨벤션에 따라서, 모든 Go 코드는 하나의 workspace(폴더) 안에 존재해야합니다. 이 workspace는 컴퓨터 어디에나 있을 수 있습니다. workspace를 지정하지 않으면 Go는 `$HOME/go`를 기본 작업 공간으로 간주합니다. workspace는 [GOPATH](https://golang.org/cmd/go/#hdr-GOPATH_environment_variable) 환경변수를 통해 식별되거나 변경될 수 있습니다.

나중에 스크립트, 셸 등에서 사용할 수 있도록 환경 변수를 설정해야 합니다.
당신의 `.bash_profile`에 아래의 expoort를 추가하세요.:

```sh
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

*Note* 이러한 환경변수를 적용하기 위하여는 새 shell을 열어야합니다.

Go는 workspace에 특정 디렉토리 구조가 포함되어 있다고 가정합니다.

Go는 go파일을 다음 세 개의 디렉토리에 배치합니다. 모든 소스 코드는 src에, 패키지 객체는 pkg에, 컴파일된 프로그램은 bin에 위치합니다. 다음과 같이 디렉터리를 생성할 수 있습니다.:

```sh
mkdir -p $GOPATH/src $GOPATH/pkg $GOPATH/bin
```

At this point you can `go get` and the `src/package/bin` will be installed correctly in the appropriate `$GOPATH/xxx` directory.

### Go Modules
Go 1.11에서는 다른 워크 플로우를 가능하게 하는 [Modules](https://github.com/golang/go/wiki/Modules)(모듈)가 소개되었습니다. 이 새로운 접근 방식은 Go 1.16 이후 기본 빌드 모드이므로 `GOPATH`를 사용하는 것은 권장되지 않습니다.

모듈은 종속성 관리, 버전 선택 및 재현 가능한 빌드와 관련된 문제를 해결하는 것을 목표로 합니다. 그리고 `GOPATH` 밖에서도 Go 코드를 실행할 수 있게 합니다.

Modules을 사용하는 것은 꽤나 직관적입니다. 프로젝트의 루트로 'GOPATH' 외부에 있는 원하는 디렉터리를 선택하고 'go mod init' 커맨드로 새 모듈을 만듭니다.

새 모듈을 만들면 모듈 위치, Go 버전, 필요한 디펜던시(빌드에 필요한 다른 모듈)가 포함된 `go.mod` 파일이 만들어집니다.

만약 `<modulepath>`가 지정되지 않는다면 `go mod init`는 디렉터리 구조에서 모듈 경로를 추측하려고 시도하지만 커맨드에 직접 `<modulepath>`를 넘겨주어 `<modulepath>`를 지정할 수도 있습니다.

```sh
mkdir my-project
cd my-project
go mod init <modulepath>
```

`go.mod`파일은 다음과 같을 수 있습니다.:

```
module cmd

go 1.16

```

기본으로 제공되는 help에는 사용 가능한 모든 `go mod` 커멘드에 대한 개요가 나와 있습니다.

```sh
go help mod
go help mod init
```

## Go Editor

어떤 편집기를 선호하는지는 개인에 따라서 매우 다르며, 이미 Go를 지원하는 환경설정이 있을 수 있습니다. 그렇지 않다면 Visual Studio Code와 같은 편집기를 고려할 수 있습니다. Visual Studio Code는 좋은 Go 지원 기능을 가지고 있습니다.

아래의 커맨드를 사용해서 설치할 수 있습니다.:

```sh
brew cask install visual-studio-code
```

VS Code가 정상적으로 설치되었는지 확인하려면 아래의 명령을 터미널에서 실행하여 확인할 수 있습니다.

```sh
code .
```

VS Code는 설정된 소프트웨어가 거의 없는 상태로 제공됩니다. extension을 설치하여 새로운 소프트웨어를 설정할 수 있습니다. Go 지원을 추가하려면 extension을 설치해야 합니다. VS Code에는 [Luke Hoban's Package](https://github.com/golang/vscode-go)와 같이 매우 뛰어고 다양한 기능이 있습니다. 이런 extension은 다음과 같이 설치할 수 있습니다.

```sh
code --install-extension golang.go
```

VS Code로 Go 코드를 처음 열었을 때, Analysis tool이 누락되었다고 나타납니다. Analysis tool은 설치 버튼을 클릭하여 설치할 수 있습니다. [여기](https://github.com/golang/vscode-go/blob/master/docs/tools.md)에서 VS Code에서 설치하고 사용할 수 있는 도구 목록을 확인할 수 있습니다.

## Go Debugger

Go를 디버깅하기 위하여 Delve를 고려할 수 있습니다. (그리고 VS Code와도 연동되어있습니다.) Delve는 아래와 같이 설치할 수 있습니다.:

```sh
go get -u github.com/go-delve/delve/cmd/dlv
```

VS Code에서 Go 디버거를 사용하기 위한 추가적인 설정을 살펴보고싶다면 [VS Code debugging documentation](https://github.com/golang/vscode-go/blob/master/docs/debugging.md)를 참고하세요.

## Go Linting

[GolangCI-Lint](https://golangci-lint.run)를 사용하여 향상된 기본 linter를 구성할 수 있습니다.

linter는 다음와 같이 설치할 수 있습니다.:

```sh
brew install golangci/tap/golangci-lint
```

## Refactoring and your tooling

이 책은 리팩터링의 중요성에 중점을 두고 있습니다.

여러분의 도구들은 더 크고 자신감 있는 리팩터링을 도와줄 수 있습니다.

간단한 키 조합으로 다음을 수행할 수 있을 정도로 편집기에 익숙해져야 합니다.

- **Extract/Inline variable**. 임의의 값에 이름을 붙힐 수 있다면 코드를 빠르게 단순하게 만들 수 있습니다.
- **Extract method/function**. 코드의 한 부분을 메서드나 함수로 추출할 수 있어야합니다.
- **Rename**. 파일에서 자신감있게 심볼의 이름을 변경할 수 있어야합니다.
- **go fmt**. Go에는 `go fmt`라고 불리는 공식 포매터가 있습니다. 에디터는 파일을 저장할 때 마다 `go fmt`를 실행시켜야합니다.
- **Run tests**. 위의 내용 중 하나를 수행한 후 바로 테스트를 다시 실행하여 리팩터링이 어떤 것도 고장나지 않도록 해야 한다는 것은 말할 필요도 없습니다.

또한 코드 작업을 편하게 하기 위해서는 다음을 수행할 수 있어야 합니다.:

- **View function signature** - Go에서 함수를 호출하는 방법을 확신할 수 없습니다. IDE는 문서, 매개변수 및 반환되는 내용에 따라 함수를 설명해야 합니다.
- **View function definition** - 함수가 어떤 것을 하는지 아직 확실하지 않으면 소스 코드로 이동하여 직접 확인할 수 있습니다.
- **Find usages of a symbol** - 호출되는 함수의 컨텍스트를 볼 수 있다면 리팩터링에서 결정을 내리는 경우에 도움이 될 수 있습니다.

도구를 마스터하면 코드에 집중하는데 도움을 주며 컨텍스트 전환을 줄일 수 있습니다.

## Wrapping up

지금 우리는 Go가 설치되어있어야 하며, 편집기와 몇가지 기본적인 도구가 준비되어 있습니다. Go는 third party 제품에 매우 커다란 에코시스템을 갖고있습니다. 우리는 몇가지 유용한 컴포넌트를 확인했으며, 보다 완벽한 리스트를 원한다면 [awesome-go](https://awesome-go.com)를 확인하세요.
