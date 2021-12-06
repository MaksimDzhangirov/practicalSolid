# SOLID на практике в Golang: принцип инверсии зависимостей

Мы продолжаем обзор принципов SOLID, рассматривая тот, который оказывает 
наиболее значительное влияние на unit тестирование в Go - принцип инверсии 
зависимостей.

![intro](images/dip/intro.jpeg)

*Фото [Jonny Gios](https://unsplash.com/@supergios) из [Unsplash](https://unsplash.com/)*

Изучение нового языка программирования часто не является чем-то сложным. Я
часто слышу: «Первый язык программирования вы выучите за год. Второй — за месяц. 
Третий за неделю, а потом каждый следующий — за день».

Это конечно преувеличение, но в некоторых случаях оно не так уж далеко от 
истины. Например, перейти на язык, относительно похожий на предыдущие (Java и 
C#), можно довольно просто.

Но иногда переключиться сложно, даже когда мы переходим от одного 
объектно-ориентированного языка на другой. На скорость перехода влияют многие 
особенности языка, такие как сильная или слабая типизация, есть ли у языка 
интерфейсы, абстрактные классы или классы вообще.

Иногда трудности возникают сразу после перехода и мы применяем практики, 
используемые в новом языке. Но с некоторыми проблемами мы сталкиваемся позже, 
например, во время unit тестирования. После этого мы понимаем почему принцип 
инверсии зависимостей важен, особенно в Go.

> Другие статьи из цикла SOLID:
> 1. [SOLID на практике в Golang: Принцип единой ответственности](https://levelup.gitconnected.com/practical-solid-in-golang-single-responsibility-principle-20afb8643483)
> 2. [SOLID на практике в Golang: Принцип открытости/закрытости](https://levelup.gitconnected.com/practical-solid-in-golang-open-closed-principle-1dd361565452)
> 3. [SOLID на практике в Golang: Принцип подстановки Барбары Лисков](https://levelup.gitconnected.com/practical-solid-in-golang-liskov-substitution-principle-e0d2eb9dd39)
> 4. [SOLID на практике в Golang: Принцип разделения интерфейсов](https://levelup.gitconnected.com/practical-solid-in-golang-interface-segregation-principle-f272c2a9a270)

> Некоторые статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Агрегат](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)
> 4. [DDD на практике в Golang: Репозиторий](https://levelup.gitconnected.com/practical-ddd-in-golang-repository-d308c9d79ba7)
> 5. ...
>
> Прим. пер. Их перевод доступен по [адресу](https://github.com/MaksimDzhangirov/practicalDDD).

## Когда мы не соблюдаем принцип инверсии зависимостей

> Модули верхних уровней не должны зависеть от модулей нижних уровней. Оба 
> типа модулей должны зависеть от абстракций. Абстракции не должны зависеть от
> деталей. Детали должны зависеть от абстракций.

Выше приведено определение принципа инверсии зависимостей (DIP). 
[Дядя Боб](https://twitter.com/unclebobmartin) представил его в своей 
[статье](https://web.archive.org/web/20110714224327/http://www.objectmentor.com/resources/articles/dip.pdf). 
Более подробно он описывает этот принцип в своём 
[блоге](https://blog.cleancoder.com/uncle-bob/2016/01/04/ALittleArchitecture.html).

Итак, как правильно понимать это определение, особенно в контексте Go? 
Во-первых, мы должны принять *Абстракцию* как 
[концепцию](https://eng.libretexts.org/Courses/Delta_College/C_-_Data_Structures/06%3A_Abstraction_Encapsulation/1.01%3A_Difference_between_Abstraction_and_Encapsulation) 
ООП. Мы используем такую концепцию, чтобы выявить основные поведения и скрыть 
детали их реализации.

Во-вторых, дадим определение модулям верхних и нижних уровней. Модули верхних 
уровней в контексте Go - это программные компоненты, используемые на высших 
уровнях приложения, например код, используемый на уровне представления.

Это также может быть код, близкий к высшим уровням, например, код бизнес-логики
или какие-то компоненты, реализующую особую логику работы. Важно понимать, 
что этот слой влияет на бизнес-логику нашего приложения.

С другой стороны, программные компоненты низкого уровня — это в основном 
небольшие фрагменты кода, поддерживающие более высокоуровневые. Они скрывают 
технические детали различных инфраструктурных интеграций.

Например, это может быть структура, в которой хранится логика для получения 
данных из базы данных, отправки SQS сообщения, получения значения из Redis или
отправки HTTP-запроса к внешнему API.

Итак, что же происходит, когда мы нарушаем принцип инверсии зависимостей и 
наш высокоуровневый компонент зависит от низкоуровневого? Разберем следующий
пример:

```go
// инфраструктурный уровень

type UserRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
    return &UserRepository{
        db: db,
    }
}

func (r *UserRepository) GetByID(id uint) (*domain.User, error) {
    user := domain.User{}
    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err
    }
    
    return &user, nil
}

// уровень предметной области

type User struct {
    ID uint `gorm:"primaryKey;column:id"`
    // какие-то поля
}

// уровень прикладных операций

type EmailService struct {
    repository *infrastructure.UserRepository
    // какой-то отправитель электронных писем
}

func NewEmailService(repository *infrastructure.UserRepository) *EmailService {
    return &EmailService{
        repository: repository,
    }
}

func (s *EmailService) SendRegistrationEmail(userID uint) error {
    user, err := s.repository.GetByID(userID)
    if err != nil {
        return err
    }
    
    fmt.Println(user)
    // отправляем электронное письмо
    return nil
}
```

В вышеприведенном фрагменте мы определили высокоуровневый компонент 
`EmailService`. Эта структура относится к уровню прикладных операций и отвечает 
за отправку электронной почты новым зарегистрированным клиентам.

Идея заключается в том, что у нас есть метод `SendRegistrationEmail`, который 
ожидает ID пользователя (`User`). Под капотом, он извлекает пользователя 
(`User`) из `UserRepository`, а затем (вероятно) передаёт его в какой-то 
сервис `EmailSender` для отправки электронной почты.

Часть касающаяся `EmailSender` сейчас нас не интересует. Давайте вместо этого 
сконцентрируемся на `UserRepository`. Эта структура представляет собой 
репозиторий, который взаимодействует с базой данных, поэтому она принадлежит к 
инфраструктурному уровню.

Итак, похоже, что наш высокоуровневый компонент `EmailService` зависит от 
низкоуровневого компонента `UserRepository`. Действительно, если мы не 
определим подключение к базе данных, мы не можем инициализировать нашу 
структуру для отправки писем.

Такой анти-шаблон сразу же влияет на unit тестирование в Go. Предположим, мы 
хотим протестировать `EmailService`, как показано во фрагменте кода ниже:

```go
import (
    "fmt"
    "testing"
    
    "github.com/DATA-DOG/go-sqlmock"
    "github.com/MaksimDzhangirov/practicalSOLID/dip/badCode1/infrastructure"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

func TestEmailService_SendRegistrationEmail(t *testing.T) {
    db, mock, err := sqlmock.New()
    assert.NoError(t, err)
    
    dialector := mysql.New(mysql.Config{
        DSN:        "dummy",
        DriverName: "mysql",
        Conn:       db,
    })
    finalDB, err := gorm.Open(dialector, &gorm.Config{})
    
    repository := infrastructure.NewUserRepository(finalDB)
    service := NewEmailService(repository)
    fmt.Println(service, mock)
    //
    // большой фрагмент кода для имитации SQL-запросов
    //
    // а затем сам тест
}
```

В отличие от некоторых языков, таких как PHP, мы не можем просто имитировать в 
Go всё, что захотим. Имитация в Go основана на использовании интерфейсов, для 
которых мы можем определить фиктивную реализацию, но не нельзя сделать то же 
самое для структур.

Итак, мы не можем имитировать UserRepository, поскольку это структура. В таком 
случае нам нужно сымитировать работу нижнего уровня, в данном случае, объекта 
подключения [Gorm](https://gorm.io/index.html), что можно сделать с помощью 
пакета [SQLMock](https://github.com/DATA-DOG/go-sqlmock).

Но даже такой вариант не является надежным и эффективным способом тестирования.
Нам нужно имитировать слишком много SQL-запросов и слишком много знать о схеме
базы данных. Любое изменение внутри базы данных требует модификации unit тестов.

Unit тестирование на самом деле даже не является самой большой проблемой. Что 
если мы захотим использовать другое хранилище для данных, например, 
[Cassandra](https://cassandra.apache.org/_/index.html)? Поскольку в нашем 
хранилище содержатся данные клиентов возможно в будущем мы планируем сделать 
его распределенным?

Если такое произойдет и мы будем использовать эту реализацию `UserRepository`, 
то потребуется отрефакторить много кода.

Теперь мы видим, к чему может привести зависимость высокоуровневого компонента 
от низкоуровневого. А что насчет абстракций, основанных на деталях? Давайте 
посмотрим на код ниже:

```go
// уровень предметной области

type User struct {
    ID uint `gorm:"primaryKey;column:id"`
    // какие-то поля
}

type UserRepository interface {
	GetByID(id uint) (*User, error)
}
```

Чтобы исправить первую проблему с высокоуровневыми и низкоуровневыми компонентами, 
мы должны начать с определения несколько интерфейсов. В этом случае мы можем 
определить `UserRepository` как интерфейс на уровне предметной области.

Таким образом, это позволяет отделить `EmailService` от базы данных, но еще не 
полностью. Посмотрите на структуру `User`. В ней по-прежнему есть дескрипторы 
полей, определяющие связь с базой данных.

И даже если такая структура находится внутри уровня предметной области, она 
всё равно содержит технические детали об инфраструктуре. Наш новый интерфейс 
`UserRepository` (абстракция) зависит от структуры `User`, связанной с базой 
данных (детали), и мы по-прежнему нарушаем принцип инверсии зависимостей (DIP).

Изменение схемы базы данных неизбежно изменит наш интерфейс. Этот интерфейс 
по-прежнему может использовать ту же структуру `User`, но она будет содержать 
изменения произошедшие на низком уровне.

В конце концов этот рефакторинг ничего не изменил. Мы всё ещё движемся в 
неправильном направлении и это приводит к следующим последствиям:

1. Мы не можем правильно протестировать нашу бизнес-логику или логику приложения.
2. Любое изменение движка базы данных или структуры таблицы влияет на более высокие уровни.
3. Мы не можем легко переключиться на другой тип хранилища.
4. Наша модель сильно связана с хранилищем данных.
5. ...

Итак, давайте еще раз проведем рефакторинг этого фрагмента кода.

## Как соблюсти принцип инверсии зависимостей

> Модули верхних уровней не должны зависеть от модулей нижних уровней. 
> **Оба типа модулей должны зависеть от абстракций.** Абстракции не должны 
> зависеть от деталей. **Детали должны зависеть от абстракций.**

Давайте вернемся к исходному определению инверсии зависимостей и рассмотрим 
выделенные жирным шрифтом предложения. Они уже дают некоторые направления для 
рефакторинга.

Мы должны определить некоторую абстракцию (интерфейс), от которой будут зависеть 
оба наших компонента, `EmailService` и `UserRepository`. Кроме того, такая 
абстракция не должна полагаться на какие-либо технические детали (например, 
объект Gorm).

Рассмотрим приведенный ниже код:

```go
// инфраструктурный уровень

type UserGorm struct {
    // какие-то поля
}

func (g UserGorm) ToUser() *domain.User {
    return &domain.User{
        // какие-то поля
    }
}

type UserDatabaseRepository struct {
    db *gorm.DB
}

var _ domain.UserRepository = &UserDatabaseRepository{}

/*
type UserRedisRepository struct {

}

type UserCassandraRepository struct {

}
*/

func NewUserDatabaseRepository(db *gorm.DB) domain.UserRepository {
    return &UserDatabaseRepository{
        db: db,
    }
}

func (r *UserDatabaseRepository) GetByID(id uint) (*domain.User, error) {
    user := UserGorm{}
    err := r.db.Where("id = ?", id).First(&user).Error
    if err != nil {
        return nil, err
    }
    
    return user.ToUser(), nil
}

// уровень предметной области

type User struct {
    // какие-то поля
}

type UserRepository interface {
    GetByID(id uint) (*User, error)
}

// уровень прикладных операций

type EmailService struct {
    repository domain.UserRepository
    // какой-то отправитель электронных писем
}

func NewEmailService(repository domain.UserRepository) *EmailService {
    return &EmailService{
        repository: repository,
    }
}

func (s *EmailService) SendRegistrationEmail(userID uint) error {
    user, err := s.repository.GetByID(userID)
    if err != nil {
        return err
    }

    fmt.Println(user)
    // отправляем электронное письмо
    return nil
}
```

В нем мы видим интерфейс `UserRepository` как компонент, который зависит от 
структуры `User`, и оба они находятся на уровне предметной области.

В структуре `User` теперь нет дескрипторов, определяющих связь со схемой базы 
данных. Для этого мы используем структуру `UserGorm`. Она находится на 
инфраструктурном уровне. В ней есть метод `ToUser`, который преобразует её в 
структуру `User`.

В этом случае мы можем использовать `UserGorm` внутри `UserDatabaseRepository`, 
которая является фактической реализацией `UserRepository`.

На уровнях предметной области и прикладных операций мы зависим только от интерфейса 
`UserRepository` и [сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3) 
`User`, которые описаны на уровне предметной области.

Внутри инфраструктурного уровня мы можем определить столько реализаций для 
`UserRepository` сколько захотим. Например, `UserFileRepository` или 
`UserCassandraRepository`.

Высокоуровневый компонент (`EmailService`) зависит от абстракции — он содержит 
поле с типом `UserRepository`. Но как низкоуровневый компонент зависит от 
абстракции?

В Go структуры [неявно](https://tour.golang.org/methods/10) реализуют интерфейсы. 
Это означает, что нам не нужно добавлять код, в котором `UserDatabaseRepository` 
явно реализует `UserRepository`, но мы можем добавить проверку с 
[пустым идентификатором](https://yourbasic.org/golang/underscore/).

Благодаря такому подходу нам будет проще контролировать наши зависимости. Наши 
структуры зависят от интерфейсов, и всякий раз, когда мы хотим изменить нашу 
общую зависимость, мы можем определить различные реализации и внедрить их.

Такой способ часто используется во многих фреймворках и с его помощью 
реализуется шаблон внедрения зависимостей. На Go реализовано множество DI 
библиотек, например, от [Facebook](https://github.com/facebookarchive/inject),
[Wire](https://github.com/google/wire) или 
[Dingo](https://github.com/i-love-flamingo/dingo).

Как теперь обстоит дело с unit тестированием? Рассмотрим пример.

```go
type GetByIDFunc func (id uint) (*domain.User, error)

func (f GetByIDFunc) GetByID(id uint) (*domain.User, error) {
    return f(id)
}

func TestEmailService_SendRegistrationEmail(t *testing.T) {
    service := NewEmailService(GetByIDFunc(func(id uint) (*domain.User, error) {
        return nil, errors.New("error")
    }))
    fmt.Println(service)
    //
    // и после этого просто вызываем сервис
}
```

Проведя такой рефакторинг мы можем задать `GetByIDFunc` как новый тип, определяющий 
функцию из `UserRepository`, которую мы хотим имитировать. В Go так часто 
делают: определяют тип-функцию и связывают с ним метод, реализующий интерфейс.

Теперь наши тесты намного элегантнее и эффективнее. Мы можем внедрить 
различные реализации `UserRepository` для любого варианта использования и 
контролировать результат теста.

## Другие примеры

Мы можем столкнуться с нарушением DIP в других компонентах, а не только в 
структурах. Например, это может быть происходить в обычных независимых 
функциях:

```go
type User struct {
    // какие-то поля
}

type UserJSON struct {
    // какие-то поля
}

func (j UserJSON) ToUser() *User {
    return &User{
        // какие-то поля
    }
}

func GetUser(id uint) (*User, error) {
    filename := fmt.Sprintf("user_%d.json", id)
    data, err := ioutil.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    var user UserJSON
    err = json.Unmarshal(data, &user)
    if err != nil {
        return nil, err
    }
    
    return user.ToUser(), nil
}
```

Итак, мы хотим считать данные пользователя (`User`). Для этого мы можем 
использовать файлы и формат JSON. Метод `GetUser` читает из файла и преобразует 
содержимое файла в фактический объект `User`.

Сам метод зависит от наличия файлов, и если мы хотим правильно его протестировать, 
нам нужно будет использовать такие файлы. Таким образом, для такого метода 
неудобно писать тесты, например, проверяющие правила валидации данных, 
переданных для заполнения структуры, если мы добавим их позже в метод `GetUser`.

Опять же, наш код зависит от слишком большого количества деталей, и было бы 
неплохо использовать некоторые абстракции:

```go
type User struct {
    // какие-то поля
}

type UserJSON struct {
    // какие-то поля
}

func (j UserJSON) ToUser() *User {
    return &User{
        // какие-то поля
    }
}

func GetUserFile(id uint) (io.Reader, error) {
    filename := fmt.Sprintf("user_%d.json", id)
    file, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    
    return file, nil
}

func GetUserHTTP(id uint) (io.Reader, error) {
    uri := fmt.Sprintf("http://some-api.com/users/%d", id)
    resp, err := http.Get(uri)
    if err != nil {
        return nil, err
    }
    
    return resp.Body, nil
}

func GetDummyUser(userJSON UserJSON) (io.Reader, error) {
    data, err := json.Marshal(userJSON)
    if err != nil {
        return nil, err
    }
    
    return bytes.NewReader(data), nil
}

func GetUser(reader io.Reader) (*User, error) {
    data, err := ioutil.ReadAll(reader)
    if err != nil {
        return nil, err
    }
    
    var user UserJSON
    err = json.Unmarshal(data, &user)
    if err != nil {
        return nil, err
    }
    
    return user.ToUser(), nil
}
```

В этой новой реализации метод `GetUser` полагается на экземпляр, реализующий 
интерфейс `Reader`. Это интерфейс из основного пакета Go, [IO](https://pkg.go.dev/io).

Здесь мы можем определить множество различных способов, обеспечивающих 
реализацию интерфейса `Reader`, например `GetUserFile`, `GetUserHTTP`, `GetDummyUser` 
(которые мы можем использовать для тестирования метода `GetUser`).

Этот подход мы можем использовать во многих различных ситуациях. Всякий раз, 
когда у нас возникают сложности при создании соответствующего unit теста или 
возникают циклические зависимости в Go, мы должны попытаться сделать код менее 
связным, предоставив интерфейс и столько реализаций, сколько захотим.

## Заключение

Принцип инверсии зависимостей — это последний принцип SOLID, и он обозначает 
букву D в слове SOLID. Он утверждает, что компоненты высокого уровня не должны 
зависеть от компонентов низкого уровня.

Вместо этого все наши компоненты должны зависеть от абстракций или, проще говоря,
интерфейсов. Такие абстракции позволяют нам более гибко использовать наш код и 
соответствующим образом его тестировать.

> Другие статьи из цикла SOLID:
> 1. [SOLID на практике в Golang: Принцип единой ответственности](https://levelup.gitconnected.com/practical-solid-in-golang-single-responsibility-principle-20afb8643483)
> 2. [SOLID на практике в Golang: Принцип открытости/закрытости](https://levelup.gitconnected.com/practical-solid-in-golang-open-closed-principle-1dd361565452)
> 3. [SOLID на практике в Golang: Принцип подстановки Барбары Лисков](https://levelup.gitconnected.com/practical-solid-in-golang-liskov-substitution-principle-e0d2eb9dd39)
> 4. [SOLID на практике в Golang: Принцип разделения интерфейсов](https://levelup.gitconnected.com/practical-solid-in-golang-interface-segregation-principle-f272c2a9a270)

> Некоторые статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Агрегат](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)
> 4. [DDD на практике в Golang: Репозиторий](https://levelup.gitconnected.com/practical-ddd-in-golang-repository-d308c9d79ba7)
> 5. ...
>
> Прим. пер. Их перевод доступен по [адресу](https://github.com/MaksimDzhangirov/practicalDDD).