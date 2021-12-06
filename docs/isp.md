# SOLID на практике в Golang: Принцип разделения интерфейсов

Мы продолжим наш обзор принципов SOLID, рассмотрев тот, который оказывает 
наиболее значительное влияние на разработку кода — принцип разделения интерфейсов.

![intro](images/isp/intro.jpeg)
*Фото [Mike Winkler](https://unsplash.com/@ahmeyer) из [Unsplash](https://unsplash.com/)*

Каждый раз когда кто-то начинает заниматься программированием, подход при обучении
один и тот же: первые несколько недель или даже месяцев всё сводится к алгоритмам
и перестройке мозга, чтобы он привык к такому образу мышления.

В какой-то момент происходит введение в объектно-ориентированное программирование. 
Если учителя слишком долго тянут с переходом на эту тему, то бывает очень сложно 
перестроиться после функционального программирования.

Но в какой-то момент мы свыкаемся с мыслью, что мы должны использовать объекты.
И мы начинаем использовать их там, где нужно, и, конечно же, там, где они не нужны.
Затем мы узнаем об абстракциях и о том, как сделать так, чтобы наш код можно было
повторно использовать.

Тогда мы можем начать неправильно использовать абстракции, добавляя их везде и
для всего. Чтобы сделать наш код повторно используемым, мы обобщаем его, 
замедляя наше дальнейшее развитие.

Рано или поздно мы начинаем понимать, что должен существовать предел для этих
обобщений. К счастью вопрос о его нахождении уже решён — используйте принцип 
разделения интерфейсов, который обозначает букву I в слове SOLID.

> Другие статьи из цикла SOLID:
> 1. [SOLID на практике в Golang: Принцип единой ответственности](https://levelup.gitconnected.com/practical-solid-in-golang-single-responsibility-principle-20afb8643483)
> 2. [SOLID на практике в Golang: Принцип открытости/закрытости](https://levelup.gitconnected.com/practical-solid-in-golang-open-closed-principle-1dd361565452)
> 3. [SOLID на практике в Golang: Принцип подстановки Барбары Лисков](https://levelup.gitconnected.com/practical-solid-in-golang-liskov-substitution-principle-e0d2eb9dd39)

> Некоторые статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Агрегат](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)
> 4. [DDD на практике в Golang: Репозиторий](https://levelup.gitconnected.com/practical-ddd-in-golang-repository-d308c9d79ba7)
> 5. ...

## Когда мы не соблюдаем принцип разделения интерфейсов

> Делайте интерфейсы небольшими, чтобы пользователи не зависели от ненужных им вещей.

[Дядя Боб](https://twitter.com/unclebobmartin) придумал этот принцип, и более подробную информацию о нем вы можете 
найти в его [блоге](https://blog.cleancoder.com/uncle-bob/2020/10/18/Solid-Relevance.html). Этот принцип четко определяет свои требования, вероятно, 
лучше всего по сравнению с другими принципами SOLID.

Из этого простого утверждения о том, что интерфейсы должны быть как можно меньше,
не следует вывод, что нужно плодить интерфейсы с одним методом. Необходимо 
учитывать контекст, объединяя характерные для него функции, которым он должен 
удовлетворять.

Давайте рассмотрим приведенный ниже код:

```go
type Money struct {
    // какие-то поля
}

type Product struct {
    // какие-то поля
}

type Wallet struct {
    // какие-то поля
}

func (w *Wallet) Deduct(money Money) error {
    return nil
}

type DiscountPolicy struct {
    // какие-то поля
}

func (d *DiscountPolicy) IsApplicableFor(customer *PremiumCustomer, product Product) bool {
    return true
}

type User interface {
    AddToShoppingCart(product Product)
    IsLoggedIn() bool
    Pay(money Money) error
    HasPremium() bool
    HasDiscountFor(product Product) bool
    //
    // некоторые дополнительные методы
    //
}

type ShoppingCart struct {

}

func (s ShoppingCart) Add(product Product) {

}

type Guest struct {
    cart ShoppingCart
    //
    // некоторые дополнительные поля
    //
}

func (g *Guest) AddToShoppingCart(product Product) {
    g.cart.Add(product)
}

func (g *Guest) IsLoggedIn() bool {
    return false
}

func (g *Guest) Pay(Money) error {
    return errors.New("user is not logged in")
}

func (g *Guest) HasPremium() bool {
    return false
}

func (g *Guest) HasDiscountFor(product Product) bool {
    return false
}

type NormalCustomer struct {
    cart ShoppingCart
    wallet Wallet
    //
    // некоторые дополнительные поля
    //
}

func (c *NormalCustomer) AddToShoppingCart(product Product) {
    c.cart.Add(product)
}

func (c *NormalCustomer) IsLoggedIn() bool {
    return true
}

func (c *NormalCustomer) Pay(money Money) error {
    return c.wallet.Deduct(money)
}

func (c *NormalCustomer) HasPremium() bool {
    return false
}

func (c *NormalCustomer) HasDiscountFor(product Product) bool {
    return false
}

type PremiumCustomer struct {
    cart ShoppingCart
    wallet Wallet
    policies []DiscountPolicy
    //
    // некоторые дополнительные поля
    //
}

func (c *PremiumCustomer) AddToShoppingCart(product Product) {
    c.cart.Add(product)
}

func (c *PremiumCustomer) IsLoggedIn() bool {
    return true
}

func (c *PremiumCustomer) Pay(money Money) error {
    return c.wallet.Deduct(money)
}

func (c *PremiumCustomer) HasPremium() bool {
    return true
}

func (c *PremiumCustomer) HasDiscountFor(product Product) bool {
    for _, p := range c.policies {
        if p.IsApplicableFor(c, product) {
            return true
        }
    }
    
    return false
}

type UserService struct {
    //
    // какие-то поля
    //
}

func (u *UserService) Checkout(ctx context.Context, user User, product Product) error {
    if !user.IsLoggedIn() {
        return errors.New("user is not logged in")
    }
    
    var money Money
    //
    // какие-то вычисления
    //
    if user.HasDiscountFor(product) {
        //
        // применить скидку
        //
    }
    return user.Pay(money)
}
```

Предположим, мы создаем приложение для онлайн торговли. Решить эту задачу можно 
различными способами. Например, одним из способов — создать интерфейс `User`, как
мы сделали в примере выше. Этот интерфейс содержит множество функций, которые 
могут понадобиться пользователю.

Пользователь (`User`) на нашей платформе может добавить товар (`Product`) в 
корзину (`ShoppingCart`). Он может купить его. Он может получить скидку на 
конкретный товар (`Product`). Единственная проблема заключается в том, что все
это может делать только определенный (премиум) пользователь (`User`).

Фактическими реализациями этого интерфейса являются три структуры. Первая — 
структура `Guest`. Это пользователь (`User`), который не осуществил вход в 
нашу систему, но, по крайней мере, может добавить товар (`Product`) в корзину
(`ShoppingCart`).

Вторая реализация — это обычный покупатель (`NormalCustomer`). Ему доступны все 
функции гостя (`Guest`) плюс возможность покупать товары (`Product`). Третья 
реализация — это премиум покупатель (`PremiumCustomer`), который может 
использовать все функции нашей системы.

Теперь взгляните на эти три структуры. Только премиум покупателю (`PremiumCustomer`)
нужны все три метода (`Pay`, `HasPremium`, `HasDiscountFor`). Возможно, мы могли 
бы определить их все для обычного покупателя (`NormalCustomer`), но точно нам 
вряд ли понадобится больше двух для гостя (`Guest`).

Методы `HasPremium` и `HasDiscountFor` не имеют смысла для гостя (`Guest`). Если эта 
структура описывает пользователя (`User`), который не вошёл в систему, зачем нам 
вообще реализовывать методы для скидок?

Здесь возможно нам стоило использовать функцию `panic` с ошибкой 
"method is not implemented" - это было бы более уместно в данном фрагменте кода. 
Но обычно мы даже не должны вызывать метод `HasPremium` для гостя (`Guest`).

Всё это было сделано для того, чтобы в одном и том же месте, одним и тем же
кодом (`UserService`) можно было обрабатывать все типы пользователей (`User`).
Но из-за этого нам нужно реализовать кучу неиспользуемых методов.

Итак, для повторного использования кода, нам пришлось:

1. Создать несколько структур с неиспользуемыми методами.
2. Как-то пометить методы структуры, которые нельзя использовать.
3. Значительно увеличить размер кода для unit тестирования.
4. Использовать не присущий объектам полиморфизм.
5. ...

Итак, давайте отрефакторим этот безобразно написанный фрагмент кода.

## Как соблюсти принцип разделения интерфейсов

> Создавая интерфейсы, объединяйте минимальную группу характерных для него функции.

Не нужно делать каких-то великих открытий. Всё что нужно — определить минимальный 
интерфейс, обеспечивающий полный набор характерных для него функций. Взгляните 
на код, представленный ниже:

```go
type Money struct {
    // какие-то поля
}

type Product struct {
    // какие-то поля
}

type Wallet struct {
    // какие-то поля
}

func (w *Wallet) Deduct(money Money) error {
    return nil
}

type DiscountPolicy struct {
    // какие-то поля
}

func (d *DiscountPolicy) IsApplicableFor(customer *PremiumCustomer, product Product) bool {
    return true
}

type ShoppingCart struct {
}

func (s ShoppingCart) Add(product Product) {

}

type User interface {
    AddToShoppingCart(product Product)
    //
    // некоторые дополнительные методы
    //
}

type LoggedInUser interface {
    User
    Pay(money Money) error
    //
    // некоторые дополнительные методы
    //
}

type PremiumUser interface {
    LoggedInUser
    HasDiscountFor(product Product) bool
    //
    // некоторые дополнительные методы
    //
}

type Guest struct {
    cart ShoppingCart
    //
    // некоторые дополнительные поля
    //
}

func (g *Guest) AddToShoppingCart(product Product) {
    g.cart.Add(product)
}

type NormalCustomer struct {
    cart   ShoppingCart
    wallet Wallet
    //
    // некоторые дополнительные поля
    //
}

func (c *NormalCustomer) AddToShoppingCart(product Product) {
    c.cart.Add(product)
}

func (c *NormalCustomer) Pay(money Money) error {
    return c.wallet.Deduct(money)
}

type PremiumCustomer struct {
    cart     ShoppingCart
    wallet   Wallet
    policies []DiscountPolicy
    //
    // некоторые дополнительные поля
    //
}

func (c *PremiumCustomer) AddToShoppingCart(product Product) {
    c.cart.Add(product)
}

func (c *PremiumCustomer) Pay(money Money) error {
    return c.wallet.Deduct(money)
}

func (c *PremiumCustomer) HasDiscountFor(product Product) bool {
    for _, p := range c.policies {
        if p.IsApplicableFor(c, product) {
            return true
        }
    }
    
    return false
}

type UserService struct {
    //
    // какие-то поля
    //
}

func (u *UserService) Checkout(ctx context.Context, user User, product Product) error {
    loggedIn, ok := user.(LoggedInUser)
    if !ok {
        return errors.New("user is not logged in")
    }
    
    var money Money
    //
    // какие-то вычисления
    //
    if premium, ok := loggedIn.(PremiumUser); ok && premium.HasDiscountFor(product) {
        //
        // применить скидку
        //
    }
    
    return loggedIn.Pay(money)
}
```

Теперь вместо одного, у нас есть три интерфейса. В `PremiumUser` встраивается 
`LoggedInUser`, в который встраивается `User`. Кроме того, каждый из них добавляет
один новый метод.

`User` теперь описывает только покупателей, которые ещё не аутентифицировались 
на нашей платформе. Мы знаем, что такой тип покупателей может использовать 
функции `ShoppingCart`.

Новый интерфейс `LoggedInUser` описывает всех наших аутентифицированных покупателей,
а интерфейс `PremiumUser` - всех аутентифицированных покупателей с платной премиальной 
учетной записью.

Обратите внимание: мы действительно добавили ещё два интерфейса, но удалили два 
метода: `IsLoggedIn` и `HasPremium`. Эти методы перестали являться частью 
нашего интерфейса. Но как организовать работу без них?

Как видите в `UserService` вместо того, чтобы использовать методы, 
возвращающие булевый результат, мы просто уточняем подтип интерфейса `User`. 
Если `User` реализует `LoggedInUser`, мы знаем, что говорим об аутентифицированном
покупателе.

Кроме того, если `User` реализует `PremiumUser`, то это покупатель с 
премиум-аккаунтом. Таким образом, уточняя тип, мы уже проверяем некоторые 
бизнес-правила.

Избавившись от этих двух методов все предыдущие структуры стали более легковесными.
Теперь у каждой из них не по пять методов, многие из которых вообще не используются,
а только действительно нужные им методы.

## Другие примеры

Хотя всегда хорошо создавать небольшие и гибкие интерфейсы, мы должны вводить 
их с учетом назначения. Не имеет особого смысла добавлять небольшие простые 
интерфейсы, если реализовывать их все вместе в одной структуре.

Посмотрите на пример, показанный ниже:

```go
// слишком сильное разбиение
type UserWithFirstName interface {
    FirstName() string
}

type UserWithLastName interface {
    LastName() string
}

type UserWithFullName interface {
    FullName() string
}

// оптимальное разбиение
type UserWithName interface {
    FirstName() string
    LastName() string
    FullName() string
}
```

Это тот случай, когда мы слишком сильно разбиваем интерфейс. Да, мы можем задать 
интерфейс для каждого метода, определяя их теперь как интерфейсы-роли. Такие
интерфейсы с одним методом иногда хороши, но не в этом случае.

Очевидно, что если покупатель зарегистрирован на нашей платформе, ему нужно 
будет указать свое имя и фамилию для выставления счетов. Итак, нашему 
покупателю (`User`) потребуются методы `FirstName` и `LastName`, а вместе с ними,
естественно, `FullName`.

В этом случае разделение этих трех методов на три интерфейса не имеет смысла,
поскольку эти три метода всегда используются вместе. Итак, это плохой пример 
интерфейса с одним методом.

Что может быть хорошим примером?

```go
package io

type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type Closer interface {
    Close() error
}

type Seeker interface {
    Seek(offset int64, whence int) (int64, error)
}

type WriteCloser interface {
    Writer
    Closer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}

// и так далее
```

Прекрасным примером в Go является пакет [ввода-вывода](https://pkg.go.dev/io). В нем находится код и 
интерфейсы для обработки операций ввода-вывода и, вероятно, все разработчики Go 
использовали этот пакет хотя бы один раз.

Он предоставляет интерфейсы `Reader`, `Writer`, `Closer`, `Seeker`. В каждом из них 
определен только один метод: чтение (`Read`), запись (`Write`), закрытие (`Close`) и
поиск (`Seek`). Мы используем их для чтения, записи, поиска фрагмента байтов от 
и до в определенном источнике и закрытия этого источника.

Чтобы обеспечить большую гибкость для таких источников, все функции размещены в 
интерфейсах. Затем они создают вместе более сложные интерфейсы, такие как 
`WriteCloser`, `ReadWriteCloser` и так далее.

## Заключение

Принцип разделения интерфейсов — четвертый принцип SOLID, обозначающий букву I 
в слове SOLID. Он приучает делать наши интерфейсы как можно меньше.

Когда мы хотим отделить объекты одного типа от другого, мы можем использовать 
различные интерфейсы. Мы не должны делать наши интерфейсы слишком маленькими и 
они должны предоставлять полный набор характерных для них функций.

> Другие статьи из цикла SOLID:
> 1. [SOLID на практике в Golang: Принцип единой ответственности](https://levelup.gitconnected.com/practical-solid-in-golang-single-responsibility-principle-20afb8643483)
> 2. [SOLID на практике в Golang: Принцип открытости/закрытости](https://levelup.gitconnected.com/practical-solid-in-golang-open-closed-principle-1dd361565452)
> 3. [SOLID на практике в Golang: Принцип подстановки Барбары Лисков](https://levelup.gitconnected.com/practical-solid-in-golang-liskov-substitution-principle-e0d2eb9dd39)

> Некоторые статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Агрегат](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)
> 4. [DDD на практике в Golang: Репозиторий](https://levelup.gitconnected.com/practical-ddd-in-golang-repository-d308c9d79ba7)
> 5. ...