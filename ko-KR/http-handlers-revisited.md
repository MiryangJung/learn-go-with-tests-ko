# HTTP Handlers Revisited

**[해당 챕터의 모든 코드는 여기에서 확인할 수 있다](https://github.com/quii/learn-go-with-tests/tree/main/q-and-a/http-handlers-revisited)**

이 책은 이미 [HTTP 핸들러 테스트하기](http-server.md)에 대한 챕터를 가지고 있지만, 여기에서는 그것들을 디자인하는 것에 대해 더 광범위한 논의를 할 것이기에 테스트가 간단하다.

현실적인 예시를 살펴보고 단일 책임 원칙 및 관심사의 분리와 같은 원칙을 적용하여 설계 방식을 어떻게 개선할 수 있는지를 살펴보겠다. 이러한 원칙들은 [인터페이스](structs-methods-and-interfaces.md)와 [의존성 주입 (dependency injection)](dependency-injection.md) 을 사용하여 실현할 수 있다. 이와 더불어 우리는 핸들러 테스트가 사실 꽤 사소하다는 점을 보여줄 것이다.

![Go 커뮤니티에 흔히 올라오는 질문을 그림화 한 것](amazing-art.png)

Go 커뮤니티에서 HTTP 핸들러 테스트에 대한 질문은 되풀이 되는 것으로 보이는데, 이는 사람들이 이를 설계하는 방법을 잘못 이해하고 있다는 점을 의미한다고 개인적으로 생각한다.

따라서 사람들이 겪는 어려움은 대체로 실제 테스트 작성이 아닌 코드 디자인에 대한 것이다. 이 책에서 자주 강조하지만:

> 만약 당신이 테스트로 인해 고통을 받는다면 이를 인지하고 당신의 코드 디자인에 대해 생각해 보자.

## 예시

[Santosh Kumar가 보낸 트윗](https://twitter.com/sntshk/status/1255559003339284481)

> mongodb 종속성이있는 http 핸들러는 어떻게 테스트합니까?

다음의 코드를 살펴 보자

```go
func Registration(w http.ResponseWriter, r *http.Request) {
	var res model.ResponseResult
	var user model.User

	w.Header().Set("Content-Type", "application/json")

	jsonDecoder := json.NewDecoder(r.Body)
	jsonDecoder.DisallowUnknownFields()
	defer r.Body.Close()

	// check if there is proper json body or error
	if err := jsonDecoder.Decode(&user); err != nil {
		res.Error = err.Error()
		// return 400 status codes
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// Connect to mongodb
	client, _ := mongo.NewClient(options.Client().ApplyURI("mongodb://127.0.0.1:27017"))
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err := client.Connect(ctx)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)
	// Check if username already exists in users datastore, if so, 400
	// else insert user right away
	collection := client.Database("test").Collection("users")
	filter := bson.D{{"username", user.Username}}
	var foundUser model.User
	err = collection.FindOne(context.TODO(), filter).Decode(&foundUser)
	if foundUser.Username == user.Username {
		res.Error = UserExists
		// return 400 status codes
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		res.Error = err.Error()
		// return 400 status codes
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}
	user.Password = string(pass)

	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		res.Error = err.Error()
		// return 400 status codes
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(res)
		return
	}

	// return 200
	w.WriteHeader(http.StatusOK)
	res.Result = fmt.Sprintf("%s: %s", UserCreated, insertResult.InsertedID)
	json.NewEncoder(w).Encode(res)
	return
}
```

먼저 위의 함수가 혼자 수행해야 하는 모든 일을 나열해 보자.

1. HTTP 응답을 작성하고 헤더, 상태 코드 등을 전송
2. 요청 본문을` User`로 디코딩
3. 데이터베이스에 대한 연결 (및 관련 모든 세부 정보)
4. 데이터베이스 질의 및 결과에 따라 일부 비즈니스 로직 적용
5. 암호 생성
6. 레코드 삽입

너무 많다.

## HTTP 핸들러란 무엇이며 무엇을 하는가?

Go의 특징들을 잠깐 내려놓아 보자. 내가 어떤 언어로 작업했는지와 상관없이  [관심사의 분리](https://en.wikipedia.org/wiki/Separation_of_concerns) 와 [단일 책임 원칙](https://en.wikipedia.org/wiki/Single-responsibility_principle) 에 대해 생각하는 것은 항상 도움이 되었다.

이는 당신이 해결하려는 문제에 따라 적용하기 쉽지 않을 수 있다. 정확히 _책임_ 이란 무엇인가?

당신이 얼마나 추상적으로 생각하는지에 따라 책임의 정의가 불분명 할 수도 있으며 때때로 처음 예상한 정의가 틀릴 수도 있다.

다행스럽게도 HTTP 핸들러와 관련해서 본인은 그것들이 어떻게 작동하는지에 대해 꽤 잘 이해하고 있다고 느끼고 있는데, 프로젝트 종류와 관계없이:

1. HTTP 요청을 수락, parse (구문 분석) 하고 유효성을 검사한다.
2. 1번에서 얻은 데이터로 `ImportantBusinessLogic` 을 수행하기 위해 `ServiceThing` 을 호출한다.
3. `ServiceThing` 이 반환하는 내용에 따라 적절한 `HTTP` 응답을 보낸다.

모든 HTTP 핸들러가 _대략_ 이와 같은 형태를 가져야 한다는것은 아니지만 개인적으로 100번 중 99번은 이에 해당했다.

관심사들을 분리함으로써:

 - 테스트 핸들러는 가벼워지며 소수의 관심사에 집중한다.
 - 중요하게도 `ImportantBusinessLogic` 테스트가 더 이상 `HTTP`와 관련이 없으며 비즈니스 로직을 깔끔하게 테스트 할 수 있다는 것이다.
 - `ImportantBusinessLogic`을 수정하지 않고도 다른 문맥에서 사용할 수 있다.
 - `ImportantBusinessLogic` 이 수정되더라도 인터페이스가 동일하게 유지되는 한 관련 핸들러를 변경할 필요가 없다.

## Go의 핸들러

[`http.HandlerFunc`](https://golang.org/pkg/net/http/#HandlerFunc)

> HandlerFunc 유형은 일반 함수를 HTTP 핸들러로 사용할 수 있는 어댑터 역할을 수행한다.

`type HandlerFunc func(ResponseWriter, *Request)`

독자여, 숨 한번 쉬고 위의 코드를 보자. 무엇을 알아차렸나?

**이것은 몇 가지의 인수를 취하는 함수이다**

해당 코드에는 프레임워크 마법도, 주석도, 마법 콩도, 아무것도 없다.

이것은 단지 함수일 뿐이며,  _우리는 함수를 테스트하는 방법을 알고 있다_.

위의 주석이 꽤 정확함을 알 수 있다.:

- 검사, 파싱 및 유효성 검사를 할 수 있는 데이터 묶음인 [`http.Request`](https://golang.org/pkg/net/http/#Request) 을 사용한다.
- > [HTTP 응답을 생성하기 위해 HTTP 핸들러가 사용하는 `http.ResponseWriter` 인터페이스](https://golang.org/pkg/net/http/#ResponseWriter)

### 왕초보 테스트 예시

```go
func Teapot(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusTeapot)
}

func TestTeapotHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	res := httptest.NewRecorder()

	Teapot(res, req)

	if res.Code != http.StatusTeapot {
		t.Errorf("got status %d but wanted %d", res.Code, http.StatusTeapot)
	}
}
```

우리의 함수를 테스트하기 위해 _호출_ 해보자.

테스트를 위해 `http.ResponseRecorder` 를 `http.ResponseWriter` 인수로 전달하고 함수는 이를 사용하여 `HTTP` 응답을 작성한다. 기록자 (recorder)는 전송된 내용을 기록 (또는 _spy_ on) 한 다음 assertions 을 작성한다.

## 핸들러에서 `ServiceThing` 호출하기

TDD 튜토리얼에 대한 일반적인 불만은 항상 "너무 단순" 하고 "충분히 현실적" 이지 못하다는 것이다. 이에 대한 내 대답은 다음과 같다 :

> 당신이 언급한 예제처럼 당신의 모든 코드가 읽고 테스트하기가 간단하다면 좋지 않을까?

이는 우리가 직면한 가장 큰 과제 중 하나로서 항상 노력해야 한다. 이와 같은 코드를 디자인하는 것이 _가능_ (반드시 쉽다고는 할 수 없지만) 하므로 우리는 좋은 소프트웨어 엔지니어링 원칙을 연습하고 적용하여 해당 코드 디자인이 읽고 테스트하기 쉬울 수 있도록 해야 한다.

이전의 핸들러가 수행하는 작업의 요점을 되풀이하자면:

1. HTTP 응답을 작성하고 헤더, 상태 코드 등의 전송
2. 요청 본문을 `User`로 디코딩
3. 데이터베이스에 연결 (및 관련 모든 세부 정보)
4. 데이터베이스 질의 및 결과에 따라 일부 비즈니스 로직 적용
5. 암호 생성
6. 레코드 삽입

개인적으로 원하는 좀 더 이상적인 관심사의 분리는 다음과 같다:

1. 요청 본문을 `User`로 디코딩.
2. `UserService.Register(user)` 의 호출 (`ServiceThing` 에 해당한다).
3. 해당 호출에서 오류가 발생한다면 (주어진 예시에서는 항상 `400 BadRequest` 을 전송하는데 이는 옳지 않다고 생각한다) _지금은_  `500 Internal Server Error` 에 대한 포괄적인 catch 핸들러를 가지도록 하겠다. 모든 오류에 대해 `500` 을 반환하는 것은 끔찍한 API 된다는 점을 분명히 하고 싶다. 이후 우리는 아마도  [error types](error-types.md) 에서 조금 더 정교한 에러 핸들러를 작성할 수 있을 것이다.
4. 해당 호출에서 오류가 발생하지 않았다면 응답 본문으로 ID와 함께 `201 Created` 을 전송한다 (위의 3번과 같이 간단히 임시적으로 말이다.).

간결함을 위해 이곳에서는 일반적인 TDD 절차들을 다루지 않겠다. 원한다면 다른 챕터에서 예제를 찾아보자.

### 새로운 디자인

```go
type UserService interface {
	Register(user User) (insertedID string, err error)
}

type UserServer struct {
	service UserService
}

func NewUserServer(service UserService) *UserServer {
	return &UserServer{service: service}
}

func (u *UserServer) RegisterUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	// request parsing and validation
	var newUser User
	err := json.NewDecoder(r.Body).Decode(&newUser)

	if err != nil {
		http.Error(w, fmt.Sprintf("could not decode user payload: %v", err), http.StatusBadRequest)
		return
	}

	// call a service thing to take care of the hard work
	insertedID, err := u.service.Register(newUser)

	// depending on what we get back, respond accordingly
	if err != nil {
		//todo: handle different kinds of errors differently
		http.Error(w, fmt.Sprintf("problem registering new user: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, insertedID)
}
```

우리의 `RegisterUser` 메서드는`http.HandlerFunc`의 모양과 일치하므로 이제 다음으로 넘어가 보자. 인터페이스로 캡처되는 `UserService` 에 대한 종속성을 가지는 새로운 유형 `UserServer` 에 대한 메서드를 첨부해두었다.

인터페이스는 우리의 `HTTP` 문제가 어느 특정한 구현체에서 분리될 수 있게 하는 환상적인 방법이다.  해당 의존성에 대해 메서드를 간단히 호출 할 수 있으며 우리는 사용자가 _어떻게_ 등록되는지 신경 쓸 필요가 없다.

만약 해당 방식에 대해 더 자세히 알아보고 싶다면 [Dependency Injection](dependency-injection.md) 챕터와  ["Build an application" 부문의 HTTP Server 섹션](http-server.md) 에서 확인 할 수 있다.

유저 등록과 관련된 특정 구현 세부 사항을 분리했으므로 이제 핸들러 코드 작성이 간단하며 앞에서 설명한 책임들을 준수한다.

### 테스트!

우리의 테스트는 이제 간단해졌다.

```go
type MockUserService struct {
	RegisterFunc    func(user User) (string, error)
	UsersRegistered []User
}

func (m *MockUserService) Register(user User) (insertedID string, err error) {
	m.UsersRegistered = append(m.UsersRegistered, user)
	return m.RegisterFunc(user)
}

func TestRegisterUser(t *testing.T) {
	t.Run("can register valid users", func(t *testing.T) {
		user := User{Name: "CJ"}
		expectedInsertedID := "whatever"

		service := &MockUserService{
			RegisterFunc: func(user User) (string, error) {
				return expectedInsertedID, nil
			},
		}
		server := NewUserServer(service)

		req := httptest.NewRequest(http.MethodGet, "/", userToJSON(user))
		res := httptest.NewRecorder()

		server.RegisterUser(res, req)

		assertStatus(t, res.Code, http.StatusCreated)

		if res.Body.String() != expectedInsertedID {
			t.Errorf("expected body of %q but got %q", res.Body.String(), expectedInsertedID)
		}

		if len(service.UsersRegistered) != 1 {
			t.Fatalf("expected 1 user added but got %d", len(service.UsersRegistered))
		}

		if !reflect.DeepEqual(service.UsersRegistered[0], user) {
			t.Errorf("the user registered %+v was not what was expected %+v", service.UsersRegistered[0], user)
		}
	})

	t.Run("returns 400 bad request if body is not valid user JSON", func(t *testing.T) {
		server := NewUserServer(nil)

		req := httptest.NewRequest(http.MethodGet, "/", strings.NewReader("trouble will find me"))
		res := httptest.NewRecorder()

		server.RegisterUser(res, req)

		assertStatus(t, res.Code, http.StatusBadRequest)
	})

	t.Run("returns a 500 internal server error if the service fails", func(t *testing.T) {
		user := User{Name: "CJ"}

		service := &MockUserService{
			RegisterFunc: func(user User) (string, error) {
				return "", errors.New("couldn't add new user")
			},
		}
		server := NewUserServer(service)

		req := httptest.NewRequest(http.MethodGet, "/", userToJSON(user))
		res := httptest.NewRecorder()

		server.RegisterUser(res, req)

		assertStatus(t, res.Code, http.StatusInternalServerError)
	})
}
```

이제 우리의 핸들러는 특정 스토리지 구현에 연결되지 않았다. 이는 `MockUserService`를 작성하여 간단하고 빠른 단위 테스트를 만들어 특정 책임들을 수행하는 데 도움이 된다.

### 데이터베이스 코드는 어떤가? 당신은 요령을 피우고 있다!

이는 모두 매우 의도적이다. 우리는 비즈니스 로직, 데이터베이스, 연결 등과 관련된 HTTP 핸들러를 원하지 않는다.

이렇게 함으로써 핸들러를 지저분한 사항으로부터 독립 했으며 _더불어_ 더 관련 없는 HTTP 사항들과 결합하지 않기 때문에 지속성 레이어와 비즈니스 로직을 더 쉽게 테스트 할 수 있도록 하였다.

이제 우리가 해야 할 일은 우리가 사용하려는 데이터베이스를 사용하여 `UserService` 를 구현하는 것이다.

```go
type MongoUserService struct {
}

func NewMongoUserService() *MongoUserService {
	//todo: pass in DB URL as argument to this function
	//todo: connect to db, create a connection pool
	return &MongoUserService{}
}

func (m MongoUserService) Register(user User) (insertedID string, err error) {
	// use m.mongoConnection to perform queries
	panic("implement me")
}
```

우리는 이것을 개별적으로 테스트 할 수 있고 만약  `main` 이 만족스럽다면 우리가 작동하는 애플리케이션을 위해 두 유닛을 함께 snap 할 수 있다.

```go
func main() {
	mongoService := NewMongoUserService()
	server := NewUserServer(mongoService)
	http.ListenAndServe(":8000", http.HandlerFunc(server.RegisterUser))
}
```

### 적은 노력으로 더 견고하고 확장 가능한 디자인 만들기

이러한 원칙은 단기적으로 우리의 삶을 더 쉽게 만들 뿐만 아니라 향후 시스템을 더 쉽게 확장 할 수 있도록 한다.

이 시스템을 추가로 반복할 때 사용자에게 등록 확인 이메일을 보내길 원하는 것은 더는 놀라운 일이 아니다.

오래된 디자인에서는 핸들러와 관련 테스트를 변경해야 했다. 이는 때때로 코드의 일부를 유지 및 보수할 수 없게 되는 방식으로, 이미 그런 식으로 설계되었기 때문에 점점 더 많은 기능이 도입된다. 이로 인해 "HTTP 처리기"가 모든 것을 처리하게 되는 것이다.

인터페이스를 사용하여 관심사를 분리하면 핸들러를 _전혀_ 수정할 필요가 없는데 이는 비즈니스 로직이 등록 절차와 아무런 관련이 없기 때문이다.

## 마치며

Go의 HTTP 핸들러를 테스트하는 것은 어렵지 않다고 하지만 좋은 소프트웨어를 설계하는 것은 어려울 수 있다!

사람들은 HTTP 핸들러가 특별하다고 생각하는 실수를 범하고 이를 작성할 때 좋은 소프트웨어 엔지니어링 관행들을 버림으로써 테스트를 어렵게 만든다.

다시 말한다.  ** Go의 http 핸들러는 함수일 뿐이다 ** 분명한 책임과 관심 사항을 잘 분리하여 다른 함수처럼 작성한다면 테스트하는 데 문제가 없으며 더 건강한 코드 베이스를 가질 수 있다.