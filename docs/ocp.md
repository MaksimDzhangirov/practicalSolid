# SOLID на практике в Golang: Принцип открытости/закрытости

Мы продолжим наш обзор SOLID принципов, рассмотрев принцип, позволяющий вносить
изменения в приложения, при этом ничего не сломав — принцип открытости/закрытости.

![intro](images/ocp/intro.jpeg)

*Фото [Tekton](https://unsplash.com/@tekton_tools) из [Unsplash](https://unsplash.com/)*

Множество различных подходов и принципов позволяют улучшать наш код в 
долгосрочной перспективе. Некоторые из них хорошо известны в сообществе 
разработчиков программного обеспечения, а некоторые почему-то остаются 
незамеченными.

На мой взгляд, так обстоит дело с Принципом открытости/закрытости (OCP), который 
обозначает букву O в слове SOLID. Из моего опыта, только люди, желающие изучить
SOLID, действительно понимают, что означает этот принцип.

Если Вы использовали шаблон проектирования ["Стратегия"](https://refactoring.guru/design-patterns/strategy),
то применяли этот принцип на практике, даже не осознавая этого. Тем не менее, 
шаблон "Стратегия" - это всего лишь одна из областей применения OCP.

В этой статье мы попытаемся понять зачем нужно использовать этот принцип. Как
обычно, все примеры будут на Go.

> Другие статьи из цикла SOLID:
> 1. [SOLID на практике в Golang: Принцип единой ответственности](https://levelup.gitconnected.com/practical-solid-in-golang-single-responsibility-principle-20afb8643483)
> Некоторые статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Агрегат](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)
> 4. [DDD на практике в Golang: Репозиторий](https://levelup.gitconnected.com/practical-ddd-in-golang-repository-d308c9d79ba7)
> 5. ...
> 
> Прим. пер. Их перевод доступен по [адресу](https://github.com/MaksimDzhangirov/practicalDDD).

## Когда мы не соблюдаем принцип открытости/закрытости

> У Вас должна быть возможность расширять поведение системы, не изменяя её.

Вышеизложенное требование OCP, [дядя Боб](https://twitter.com/unclebobmartin) 
привёл в своём [блоге](http://blog.cleancoder.com/uncle-bob/2014/05/12/TheOpenClosedPrinciple.html). 
Мне нравится такой способ определения принципа открытости/закрытости, поскольку
он демонстрирует всё его красоту.

На первый взгляд это требование кажется абсурдным. Действительно, как можно
что-то расширить, не модифицируя? Я имею в виду, как что-то поменять, не меняя?

Взгляните на код, показанный ниже, чтобы понять, как некоторые структуры могут
не соблюдать этот принцип и возможные последствия.

```go
package badCode

import (
    "net/http"
    
    "github.com/ahmetb/go-linq"
    "github.com/gin-gonic/gin"
)

type PermissionChecker struct {
    //
    // какие-то поля
    //
}

func (c *PermissionChecker) HasPermission(ctx *gin.Context, name string) bool {
    var permissions []string
    switch ctx.GetString("authType") {
    case "jwt":
        permissions = c.extractPermissionsFromJwt(ctx.Request.Header)
    case "basic":
        permissions = c.getPermissionsForBasicAuth(ctx.Request.Header)
    case "applicationKey":
        permissions = c.getPermissionsForApplicationKey(ctx.Query("applicationKey"))
    }
    
    var result []string
    linq.From(permissions).
        Where(func(permission interface{}) bool {
            return permission.(string) == name
        }).ToSlice(&result)
    
    return len(result) > 0
}

func (c *PermissionChecker) getPermissionsForApplicationKey(key string) []string {
    var result []string
    //
    // получаем права доступа из key
    //
    return result
}

func (c *PermissionChecker) getPermissionsForBasicAuth(h http.Header) []string {
    var result []string
    //
    // получаем права доступа из заголовка
    //
    return result
}

func (c *PermissionChecker) extractPermissionsFromJwt(h http.Header) []string {
    var result []string
    //
    // извлекаем права доступа из JWT
    //
    return result
}
```

В вышеприведенном примере показана одна структура `PermissionChecker`. Она 
проверяет есть ли необходимые права доступа к какому-либо ресурсу, зависящие 
от контекста (`Context`) веб-приложения, поддерживаемого пакетом 
[Gin](https://github.com/gin-gonic/gin).

Здесь у нас есть основной метод HasPermission, который проверяет есть ли в
контексте (`Context`) определённые поля, связанные с правами доступа.

Процедура извлечения прав доступа из контекста (`Context`) может меняться в зависимости
от того, авторизуется ли пользователь с помощью JWT токена, базовой авторизации
или ключа приложения. Внутри структуры представлены различные способы извлечения
среза с правами доступа.

Если мы соблюдаем принцип единой ответственности, `PermissionChecker` отвечает
за определение того, находятся ли права доступа внутри контекста (`Context`), и
никак не связан с процессом авторизации.

Возможно процесс авторизации определен где-то еще, в какой-то другой структуре, 
может быть, даже в другом модуле. Таким образом, если мы хотим расширить 
процесс авторизации где-то, нам также необходимо адаптировать логику здесь.

Предположим, мы хотим расширить логику авторизации и добавить новые 
возможности, например, сохранение пользовательских данных в сессии или 
использовать [Дайджест-аутентификацию](https://httpwg.org/specs/rfc7616.html).
В этом случае нам также необходимо внести изменения в `PermissionChecker`.

Такая реализация порождает целый ряд проблем:

1. `PermissionChecker` содержит логику, которая уже присутствует где-либо ещё.
2. Любое изменение логики авторизации, которая может находится в другом модуле,
   требует адаптации в `PermissionChecker`.
3. Чтобы добавить новый способ извлечения прав доступа, нам всегда необходимо
   модифицировать `PermissionChecker`.
4. Логика внутри `PermissionChecker` будет неизбежно расти с каждым новым способом 
   авторизации.
5. Unit тест для `PermissionChecker` будет содержать слишком много технических деталей
   о различных способах извлечения прав доступа.
6. ...

Итак, опять у нас появился код, который нужно отрефакторить.

## Как соблюсти принцип открытости/закрытости

> Принцип открытости/закрытости гласит, что программные структуры должны быть 
> **открыты** для расширения, но **закрыты** для модификации.

Вышеприведенное утверждение подсказывает нам как должен выглядеть наш код, чтобы
он соблюдал OCP. Такой код должен позволять вносить расширения извне.

В объектно-ориентированном программировании мы поддерживаем такие расширения, 
используя разные реализации для одного и того же интерфейса. Другими словами, 
с помощью [полиморфизма](https://en.wikipedia.org/wiki/Polymorphism_(computer_science)).

```go
type PermissionProvider interface {
    Type() string
    GetPermissions(ctx *gin.Context) []string
}

type PermissionChecker struct {
    providers []PermissionProvider
    //
    // какие-то поля
    //
}

func (c *PermissionChecker) HasPermission(ctx *gin.Context, name string) bool {
    var permissions []string
    for _, provider := range c.providers {
        if ctx.GetString("authType") != provider.Type() {
            continue
        }
    
        permissions = provider.GetPermissions(ctx)
        break
    }
    
    var result []string
    linq.From(permissions).
            Where(func(permission interface{}) bool {
                return permission.(string) == name
            }).ToSlice(&result)
    
    return len(result) > 0
}
```

В вышеприведенном примере показан один из способов соблюдения OCP. Адаптер
`PermissionChecker` не содержит технических подробностей об извлечении прав 
доступа из контекста (`Context`).

Вместо этого, мы вводим новый интерфейс `PermissionProvider`. Здесь будет 
храниться логика для извлечения прав доступа различными способами.

Например, можно реализовать `JwtPermissionProvider`, `ApiKeyPermissionProvider` 
или `AuthBasicPermissionProvider`. Теперь модуль, отвечающий за авторизацию, 
может также содержать логику, отвечающую за извлечение прав доступа.

Это означает, что мы можем хранить логику, связанную с авторизацией пользователей
в одном месте, не размазывая её по всему коду.

С другой стороны, наша основная задача — расширить `PermissionChecker` без 
необходимости его модификации — теперь выполнима. Мы можем инициализировать
`PermissionChecker` с любым количеством `PermissionProvider`ов.

Допустим, нам нужно добавить возможность получения прав доступа из ключа сессии.
В этом случае нам нужно реализовать новый SessionPermissionProvider, который 
будет извлекать cookie из контекста (`Context`) и использовать его для 
получения прав доступа из `SessionStore`.

У нас появилась возможность расширять `PermissionChecker` всякий раз, когда 
это необходимо, не изменяя его внутреннюю логику. Теперь нам должно быть понятно,
что означает быть открытым для расширения и закрытым для модификации.

## Другие примеры

Предыдущую проблему можно решить немного иначе. Давайте посмотрим на следующий
фрагмент кода:

```go
type PermissionProvider interface {
    Type() string
    GetPermissions(ctx *gin.Context) []string
}

type PermissionChecker struct {
    //
    // какие-то поля
    //
}

func (c *PermissionChecker) HasPermission(ctx *gin.Context, provider PermissionProvider, name string) bool {
    permissions := provider.GetPermissions(ctx)
    
    var result []string
    linq.From(permissions).
            Where(func(permission interface{}) bool {
                return permission.(string) == name
            }).ToSlice(&result)
    
    return len(result) > 0
}
```

Здесь мы удалили срез из `PermissionProviders` в `PermissionChecker`. Вместо
этого мы передаём необходимый `provider` в качестве аргумента метода 
`HasPermission`.

Мне больше нравится первый способ, но этот тоже можно использовать, в 
зависимости от нашего приложения.

Мы можем применять принцип открытости/закрытости к методам, а не только к 
структурам. Примером может служить код ниже:

```go
type City struct {
    name      string
    latitude  float64
    longitude float64
    country   string
}

func GetCities(sourceType string, source string) ([]City, error) {
    var data []byte
    var err error

    if sourceType == "file" {
        data, err = ioutil.ReadFile(source)
        if err != nil {
            return nil, err
        }
    } else if sourceType == "link" {
        resp, err := http.Get(source)
        if err != nil {
            return nil, err
        }

        data, err = ioutil.ReadAll(resp.Body)
        if err != nil {
            return nil, err
        }
        defer resp.Body.Close()
    }

    var cities []City
    err = yaml.Unmarshal(data, &cities)
    if err != nil {
        return nil, err
    }

    return cities, nil
}
```

Функция `GetCities` считывает список городов из некоторого источника. Этим 
источником может быть файл или какой-либо ресурс в Интернете. Кроме того, мы
можем захотеть в будущем считывать данные из памяти, из Redis или из любого
другого источника.

Так или иначе, лучше сделать процесс чтения необработанных данных немного 
более абстрактным. С учетом сказанного, мы можем передавать способ чтения в 
качестве аргумента метода.

```go
type DataReader func(source string) ([]byte, error)

func ReadFromFile(fileName string) ([]byte, error) {
    data, err := ioutil.ReadFile(fileName)
    if err != nil {
        return nil, err
    }
   
    return data, nil
}

func ReadFromLink(link string) ([]byte, error) {
    resp, err := http.Get(link)
    if err != nil {
        return nil, err
    }
   
    data, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
   
    return data, nil
}

func GetCities(reader DataReader, source string) ([]City, error) {
    data, err := reader(source)
    if err != nil {
        return nil, err
    }
   
    var cities []City
    err = yaml.Unmarshal(data, &cities)
    if err != nil {
        return nil, err
    }
    
    return cities, nil
}
```

Как видно из вышеприведенного решения, в Go мы можем определить новый тип, в
который представляет собой функцию. Здесь мы описали новый тип `DataReader`, 
использующийся для чтения необработанных данных из некоторого источника.

Новые методы `ReadFromFile` и `ReadFromLink` являются фактическими реализациями 
типа `DataReader`. Метод `GetCities` ожидает фактическую реализацию 
`DataReader` в качестве аргумента, который затем вызывается внутри тела функции,
чтобы получить необработанные данные.

Как видите, основная задача OCP - обеспечить гибкость нашему коду и 
пользователям нашего кода. Наши библиотеки представляют действительную 
ценность, если кто-то может расширить их без копирования репозитория, 
осуществления pull-запросов или внесения каких-либо изменений в них. 

## Заключение

Принцип открытости/закрытости — это второй принцип SOLID, обозначающийся в нём
буквой O. Он утверждает, что мы всегда должны расширять наши структуры кода, а
не изменять их.

Мы должны использовать полиморфизм для удовлетворения этого требования. Наш код 
должен предоставлять простой интерфейс для добавления такой расширяемости.
