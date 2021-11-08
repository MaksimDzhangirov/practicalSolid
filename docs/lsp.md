# SOLID на практике в Golang: Принцип подстановки Барбары Лисков

Мы продолжаем наш обзор SOLID принципов и рассмотрим сегодня тот, который имеет
наиболее сложное определение — принцип подстановки Барбары Лисков.

![intro](images/lsp/intro.jpeg)
*Фото [Markus Spiske](https://unsplash.com/@markusspiske) из [Unsplash](https://unsplash.com/)*

Мне сложно что-то понять просто прочитав это. Очень часто при чтении я понимаю, 
что потерял нить повествования за последние несколько минут. Я могу прочитать 
целую главу так и не поняв о чём она, дойдя до её конца.

Иногда опускаются руки, когда я пытаюсь сосредоточиться на теме повествования, но 
вскоре понимаю, что нужно начать читать заново. Тогда я пытаюсь найти другие 
способы изучить материал.

Впервые у меня возникла такая проблема с чтением при изучении SOLID принципов,
когда я столкнулся с принципом подстановки Барбары Лисков. Определение было (и 
остаётся) слишком сложным с моей точки зрения, по крайней мере в своём 
первоначальном виде.

Как можно догадаться, LSP обозначает букву L в слове SOLID. На самом деле принцип
не так сложно понять (хотя хорошо бы было иметь менее формализованное, 
математическое определение).

> Другие статьи из цикла SOLID:
> 1. [SOLID на практике в Golang: Принцип единой ответственности](https://levelup.gitconnected.com/practical-solid-in-golang-single-responsibility-principle-20afb8643483)
> 2. [SOLID на практике в Golang: Принцип открытости/закрытости](https://levelup.gitconnected.com/practical-solid-in-golang-open-closed-principle-1dd361565452)

> Некоторые статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Агрегат](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)
> 4. [DDD на практике в Golang: Репозиторий](https://levelup.gitconnected.com/practical-ddd-in-golang-repository-d308c9d79ba7)
> 5. ...

## Когда мы не соблюдаем принцип подстановки Барбары Лисков

Впервые мы услышали об этом принципе в 1988 году от Барбары Лисков. Позднее 
[Дядя Боб](https://twitter.com/unclebobmartin) высказал свое [мнение](https://web.archive.org/web/20151128004108/http://www.objectmentor.com/resources/articles/lsp.pdf) 
по этой теме в своей статье и позже использовал 
его как один из принципов SOLID. Посмотрим о чём он гласит:

> Пусть Ф(x) является свойством верным относительно объектов x некоторого типа T.
> Тогда Ф(y) также должно быть верным для объектов y типа S, где S является 
> подтипом типа T.

Не самое удачное определение, не так ли?

Нет, серьезно, что это за определение? При написании этой статьи я все еще не 
смог понять смысл этого определения, несмотря на то, что прекрасно знаю LSP. 
Давайте дадим другое определение:

> Если S является подтипом T, тогда объекты типа T в программе могут быть 
> заменены объектами типа S без каких-либо изменений свойств, которые должна 
> иметь эта программа

Теперь стало немного понятнее о чём речь. Если `ObjectA` - это экземпляр 
`ClassA` и `ObjectB` - экземпляр `ClassB`, и `ClassB` - подтип `ClassA` - если мы
используем `ObjectB` вместо `ObjectA` где-нибудь в коде, то работа приложения не 
должна нарушиться.

Мы говорим здесь о классах и наследовании, двух парадигмах, которых нет в Go. 
Тем не менее, этот принцип можно применить, используя **интерфейсы** и 
**полиморфизм**.

```go
type User struct {
    ID uuid.UUID
    //
    // какие-то поля
    //
}

type UserRepository interface {
    Update(ctx context.Context, user User) error
}

type DBUserRepository struct {
    db *gorm.DB
}

func (r *DBUserRepository) Update(ctx context.Context, user User) error {
    return r.db.WithContext(ctx).Delete(user).Error
}
```

Здесь показан пример кода. И, честно говоря, хуже и глупее, я найти не смог. 
Например, вместо обновления пользователя в базе данных, как следует из названия 
метода `Update`, он удаляет его.

Но вот в чём дело. Мы видим интерфейс `UserRepository`. После интерфейса определена
структура `DBUserRepository`. Хотя эта структура реализует исходный интерфейс,
она не выполняет то, что требует от неё интерфейс.

Она делает совсем не то, чего от неё ожидает интерфейс. В этом суть LSP в Go: 
структура должна соответствовать ожиданиям интерфейса.

Теперь давайте посмотрим на менее очевидные примеры:

```go
type UserRepository interface {
    Create(ctx context.Context, user badCode1.User) (*badCode1.User, error)
    Update(ctx context.Context, user badCode1.User) error
}

type DBUserRepository struct {
    db *gorm.DB
}

func (r *DBUserRepository) Create(ctx context.Context, user badCode1.User) (*badCode1.User, error) {
    err := r.db.WithContext(ctx).Create(&user).Error
    return &user, err
}

func (r *DBUserRepository) Update(ctx context.Context, user badCode1.User) error {
    return r.db.WithContext(ctx).Save(&user).Error
}

type MemoryUserRepository struct {
    users map[uuid.UUID]badCode1.User
}

func (r *MemoryUserRepository) Create(_ context.Context, user badCode1.User) (*badCode1.User, error) {
    if r.users == nil {
        r.users = map[uuid.UUID]badCode1.User{}
    }
    user.ID = uuid.New()
    r.users[user.ID] = user
    
    return &user, nil
}

func (r *MemoryUserRepository) Update(_ context.Context, user badCode1.User) error {
    if r.users == nil {
        r.users = map[uuid.UUID]badCode1.User{}
    }
    r.users[user.ID] = user
    
    return nil
}
```

Здесь у нас определён новый интерфейс `UserRepository` и две его реализации:
`DBUserRepository` и `MemoryUserRepository`. Как видно, для `MemoryUserRepository`
не нужен аргумент `context.Context`, но он задан, чтобы соответствовать 
интерфейсу.

И в этом заключается проблема. Мы адаптировали `MemoryUserRepository`, чтобы
он соответствовал интерфейсу, добавив не присущий ему аргумент. Благодаря этому,
мы можем переключаться между источниками данных в нашем приложении, где один
из источников не является постоянным хранилищем.

Цель шаблона [Репозиторий](https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/infrastructure-persistence-layer-design)
- предоставить интерфейс к соответствующему постоянному хранилищу данных,
например к базе данных. Он не должен играть роль системы кеширования, как здесь,
когда мы храним `Users` в памяти.
  
Иногда это приводит не только к семантическим последствиям, но и влияет на код.
Такие случаи бросаются в глаза при реализации и от них сложнее всего избавится,
поскольку это требует серьёзного рефакторинга кода.

Чтобы продемонстрировать такой случай, рассмотрим [известный пример](http://stg-tud.github.io/sedc/Lecture/ws13-14/3.3-LSP.html#mode=document) с 
геометрическими фигурами. Интересно, что он противоречит геометрическому факту.

```go
type ConvexQuadrilateral interface {
    GetArea() int
}

type Rectangle interface {
    ConvexQuadrilateral
    SetA(a int)
    SetB(b int)
}

type Oblong struct {
    Rectangle
    a int
    b int
}

func (o *Oblong) SetA(a int) {
    o.a = a
}

func (o *Oblong) SetB(b int) {
    o.b = b
}

type Square struct {
    Rectangle
    a int
}

func (o *Square) SetA(a int) {
    o.a = a
}

func (o Square) GetArea() int {
    return o.a * o.a
}

func (o *Square) SetB(b int) {
    //
    // должно ли o.a быть равно b?
    // или оставаться неопределенным?
    //
}
```

В вышеприведённом примере мы видим реализацию геометрических фигур в Go. В геометрии 
мы можем сравнивать выпуклые четырехугольники, используя [подтипы](https://en.wikipedia.org/wiki/Quadrilateral#Convex_quadrilaterals): 
прямоугольник, квадрат.

Если мы перенесём эту логику для реализации вычисления площади в Go код, то
получим фрагмент похожий на показанный выше. Вверху определен интерфейс
`ConvexQuadrilateral`.

В этом интерфейсе задан только один метод, `GetArea`. Как подтип интерфейса
`ConvexQuadrilateral` мы можем определить интерфейс `Rectangle`. У этого 
подтипа используются две стороны для вычисления площади, поэтому нужны методы
`SetA` и `SetB`.

Далее следуют реализации. Первая - `Oblong`, у которой одна из сторон длиннее.
В геометрии — это любой четырёхугольник с прямыми углами, который не является 
квадратом. Реализовать логику работы для этой структуры несложно.

Второй подтип `Rectangle` - `Square`. В геометрии квадрат — это подтип четырёхугольника
с прямыми углами, но если следовать этому принципу при разработке программного обеспечения,
это приведёт только к проблемам при реализации.

У квадрата все стороны равны. Таким образом, метод `SetB` не нужен. Первоначально 
выбрав такие подтипы, в ходе реализации мы поняли, что в нашем коде есть ненужные
методы. Аналогичная проблема возникает, если пойти другим путём:

```go
type ConvexQuadrilateral interface {
    GetArea() int
}

type EquilateralRectangle interface {
    ConvexQuadrilateral
    SetA(a int)
}

type Oblong struct {
    EquilateralRectangle
    a int
    b int
}

func (o *Oblong) SetA(a int) {
    o.a = a
}

func (o *Oblong) SetB(b int) {
    // где определён этот метод?
    o.b = b
}

func (o Oblong) GetArea() int {
    return o.a * o.b
}

type Square struct {
    EquilateralRectangle
    a int
}

func (o *Square) SetA(a int) {
    o.a = a
}

func (o Square) GetArea() int {
    return o.a * o.a
}
```

В этом примере вместо `Rectangle` мы ввели интерфейс `EquateralRectangle`. Геометрически
это прямоугольник с равными сторонами.

В этом случае в нашем интерфейсе определен только метод `SetA`, таким образом,
нам не нужно будет реализовывать лишние методы. Тем не менее, это
нарушает LSP, поскольку мы ввели дополнительный метод `SetB` для `Oblong`, без
которого мы не можем вычислить площадь, даже если наш интерфейс говорит 
обратное.

Итак, постепенно мы начинаем улавливать идею принципа подстановки Барбары Лисков
в Go. Подытожим, что может пойти не так, если мы нарушаем его:

1) создание реализаций, с не присущими для него аргументами.
2) наличие ненужного кода.
3) нарушение ожидаемой последовательности выполнения кода.
4) нарушение требуемой логики работы.
5) возникновение интерфейса, который невозможно реализовать.
6) ...

Итак, снова нужно провести рефакторинг.

## Как соблюсти принцип подстановки Барбары Лисков

> Мы можем реализовать подтипы в Go с помощью интерфейсов, только если 
> они будут соответствовать ожиданиями интерфейса и его методам.

Я не будут приводить здесь код для первого примера, поскольку понятно, что
метод `Update` должен обновлять информацию о пользователе, а не удалять его.

Поэтому давайте перейдём к исправлению второго примера, где даны различные 
реализации интерфейса `UserRepository`:

```go
type UserRepository interface {
    Create(ctx context.Context, user badCode1.User) (*badCode1.User, error)
    Update(ctx context.Context, user badCode1.User) error
}

type MySQLUserRepository struct {
    db *gorm.DB
}

type CassandraUserRepository struct {
    session *gocql.Session
}

type UserCache interface {
    Create(user badCode1.User)
    Update(user badCode1.User)
}

type MemoryUserCache struct {
    users map[uuid.UUID]badCode1.User
}
```

В этом примере мы разделили интерфейс на два с чётким назначением и сигнатурами
различных методов. Теперь у нас есть интерфейс `UserRepository` и `UserCache`.

Цель `UserRepository` теперь определенно заключается в том, чтобы хранить 
пользовательские данные в каком-либо постоянном хранилище. Для этого мы 
подготовили конкретные реализации, например, `MySQLUserRepository` и 
`CassandraUserRepository`.

С другой стороны, мы четко понимаем, что интерфейс `UserCache` нужен для 
временного хранения пользовательских данных в каком-то кеше. В качестве конкретной 
реализации мы можем использовать `MemoryUserCache`.

Теперь перейдём к примеру с геометрическими фигурами. Тут ситуация немного 
сложнее:

```go
type ConvexQuadrilateral interface {
    GetArea() int
}

type EquilateralQuadrilateral interface {
    ConvexQuadrilateral
    SetA(a int)
}

type NonEquilateralQuadrilateral interface {
    ConvexQuadrilateral
    SetA(a int)
    SetB(b int)
}

type NonEquiangularQuadrilateral interface {
    ConvexQuadrilateral
    SetAngle(angle float64)
}

type Oblong struct {
    NonEquilateralQuadrilateral
    a int
    b int
}

type Square struct {
    EquilateralQuadrilateral
    a int
}

type Parallelogram struct {
    NonEquilateralQuadrilateral
    NonEquiangularQuadrilateral
    a     int
    b     int
    angle float64
}

type Rhombus struct {
    EquilateralQuadrilateral
    NonEquiangularQuadrilateral
    a     int
    angle float64
}
```

Чтобы поддерживать подтипы геометрических фигур в Go, мы должны учесть все их 
особенности. Таким образом, у нас не будет нарушений логики или лишних методов.

В этом случае мы вводим три новых интерфейса: `EquilateralQuadrilateral` 
(четырёхугольник с четырьмя равными сторонами), `NonEquilateralQuadrilateral` (четырёхугольник,
у которого равны две пары сторон), и `NonEquiangularQuadrilateral` (четырёхугольник, где 
равны две пары углов).

У каждого из этих интерфейсов будут дополнительные методы, необходимые для получения
данных, требуемых для вычисления площади.

Теперь мы можем определить `Square`, используя только метод `SetA`, `Oblong` - с помощью 
`SetA` и `SetB`, а для `Parallelogram` понадобятся вышеперечисленные методы плюс
`SetAngle`. Здесь мы руководствовались не принципом, что одна геометрическая фигура 
является частным случаем другой, а свойствами, описывающими их.

В обоих случаях мы модифицировали код, чтобы он всегда соответствовал ожиданиями 
конечного пользователя. В нём нет лишних методов и нарушений логики. Такой код
будет работать стабильно.

## Заключение

Принцип подстановки Барбары Лисков учит нас правильно работать с подтипами. Мы 
никогда не должны принудительно применять полиморфизм, даже если это отражает то, 
что происходит на самом деле.

LSP обозначает букву L в слове SOLID. Хотя он связан с наследованием и классами,
которые не поддерживаются в Go, мы всё же можем использовать этот принцип для
полиморфизма и интерфейсов.

> Другие статьи из цикла SOLID:
> 1. [SOLID на практике в Golang: Принцип единой ответственности](https://levelup.gitconnected.com/practical-solid-in-golang-single-responsibility-principle-20afb8643483)
> 2. [SOLID на практике в Golang: Принцип открытости/закрытости](https://levelup.gitconnected.com/practical-solid-in-golang-open-closed-principle-1dd361565452)

> Некоторые статьи из DDD цикла:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Агрегат](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)
> 4. [DDD на практике в Golang: Репозиторий](https://levelup.gitconnected.com/practical-ddd-in-golang-repository-d308c9d79ba7)
> 5. ...