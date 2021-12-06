# SOLID на практике в Golang: Принцип единой ответственности

Начнём наш обзор основных принципов разработки программного обеспечения с 
самого известного из них — принципа единой ответственности.

![intro](images/srp/intro.jpeg)

*Фото [Hunter Haley](https://unsplash.com/@hnhmarketing) из [Unsplash](https://unsplash.com/)*

Возможностей для совершения какого-то крупного прорыва в области разработки 
программного обеспечения не так уж много. Чаще всего они возникают из-за 
перестройки логики после неправильного начального обучения или восприятия
недостающего элемента в наших знаниях.

Мне нравится это чувство более глубокого осмысления. Иногда во время 
программирования, иногда при чтении книги или статьи в Интернете, а иногда
буквально сидя в автобусе.

Часто мы слышим внутренний голос при этом. Ах, да, вот как это должно работать.
Внезапно все прошлые ошибки имеют логическую причину. Все будущие требования 
приобретают форму.

У меня произошёл такой прорыв при изучении SOLID принципов. Эти принципы 
впервые были представлены в [документе](https://web.archive.org/web/20150906155800/http://www.objectmentor.com/resources/articles/Principles_and_Patterns.pdf)
[Дяди Боба](https://twitter.com/unclebobmartin). Позднее он сформировал их в 
своей книге ["Чистая архитектура"](https://www.amazon.com/dp/935286512X).

В этой статье я планирую начать обзор всех SOLID принципов, приводя примеры на 
Go. Первый в списке, обозначающий букву *S* в слове *SOLID* - это принцип единой
ответственности.

> Возможно вас также заинтересуют принципы предметно-ориентированного проектирования? 
> Ознакомьтесь с этим циклом статей:
> 1. [DDD на практике в Golang: Объект-значение](https://levelup.gitconnected.com/practical-ddd-in-golang-value-object-4fc97bcad70)
> 2. [DDD на практике в Golang: Сущности](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
> 3. [DDD на практике в Golang: Сервисы предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-service-4418a1650274)
> 4. [DDD на практике в Golang: Событие предметной области](https://levelup.gitconnected.com/practical-ddd-in-golang-domain-event-de02ad492989)
> 5. [DDD на практике в Golang: Модуль](https://levelup.gitconnected.com/practical-ddd-in-golang-module-51edf4c319ec)
> 6. [DDD на практике в Golang: Агрегат](https://levelup.gitconnected.com/practical-ddd-in-golang-aggregate-de13f561e629)
> 7. [DDD на практике в Golang: Фабрика](https://levelup.gitconnected.com/practical-ddd-in-golang-factory-5ba135df6362)
> 8. [DDD на практике в Golang: Репозиторий](https://levelup.gitconnected.com/practical-ddd-in-golang-repository-d308c9d79ba7)
> 9. [DDD на практике в Golang: Спецификация](https://levelup.gitconnected.com/practical-ddd-in-golang-specification-6523d14438e6)
> 
> Прим. пер. Их перевод доступен по [адресу](https://github.com/MaksimDzhangirov/practicalDDD).

## Когда мы не соблюдаем единую ответственность

> Принцип единой ответственности (SRP) гласит, что у каждого программного модуля
> должна быть одна и только одна причина для изменения.

Вышеприведенное предложение [написано](http://blog.cleancoder.com/uncle-bob/2014/05/08/SingleReponsibilityPrinciple.html)
самим дядей Бобом. Его смысл первоначально касался модулей и разделению обязанностей
путём сопоставления их с повседневной работой организации.

Сегодня SRP имеет широкий спектр и затрагивает различные аспекты программного
обеспечения. Мы можем использовать его в классах, функциях, модулях. И, 
естественно, в Go, мы можем использовать его в структуре.

```go
type EmailGorm struct {
    gorm.Model
    From    string
    To      string
    Subject string
    Message string
}

type EmailService struct {
    db           *gorm.DB
    smtpHost     string
    smtpPassword string
    smtpPort     int
}

func NewEmailService(db *gorm.DB, smtpHost string, smtpPassword string, smtpPort int) *EmailService {
    return &EmailService{
        db:           db,
        smtpHost:     smtpHost,
        smtpPassword: smtpPassword,
        smtpPort:     smtpPort,
    }
}

func (s *EmailService) Send(from string, to string, subject string, message string) error {
    email := EmailGorm{
        From: from,
        To:   to,
        Subject: subject,
        Message: message,
    }
    
    err := s.db.Create(&email).Error
    if err != nil {
        log.Println(err)
        return err
    }
    
    auth := smtp.PlainAuth("", from, s.smtpPassword, s.smtpHost)
    
    server := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)
    
    err = smtp.SendMail(server, auth, from, []string{to}, []byte(message))
    if err != nil {
        log.Println(err)
        return err
    }
    
    return nil
}
```

Давайте рассмотрим вышеприведенный фрагмент кода. У нас есть структура 
`EmailService` с одним методом, `Send`. Мы используем этот сервис для отправки
электронных писем. Кажется в коде никаких проблем нет, но мы понимаем, что этот
код ломает все принципы SRP, если копнуть глубже.

`EmailService` отвечает не только за отправку электронных писем, но и за хранение 
сообщения в базе данных **и** отправку его по протоколу SMTP.

Присмотритесь в предыдущее предложение. Слово "и" не зря выделено жирным 
шрифтом. При описании сервиса, выполняющего что-то одно, мы бы вряд ли 
использовали это слово.

> Как только в описании ответственности некоторой структуры кода требуется
> использовать слово "и", это уже нарушает принцип единой ответственности.

В нашем примере мы нарушили SRP сразу на нескольких уровнях. Во-первых, на 
функциональном уровне. Функция `Send` отвечает за хранение сообщения в базе 
данных и отправку электронной почты по протоколу SMTP.

Во-вторых, на уровне структуры `EmailService`. Как мы уже сказали, он 
ответственен за хранение в базе данных и отправку электронных писем. 

К чему приведёт использование такого кода?

1. При изменении структуры таблицы или тип хранилища, нам нужно будет поменять
   код для отправки электронных писем через SMTP.
2. Когда мы захотим интегрировать [Mailgun](https://www.mailgun.com/) или 
   [Mailjet](https://www.mailjet.de/), нужно будет изменить код, хранящий данные
   в БД MySQL.
3. Если мы захотим отправлять электронную почту в приложении другим способом, 
   то такой способ должен иметь логику для сохранения данных в БД.
4. Предположим, мы решили разделить разработку приложения между двумя командами:
   одна будет отвечать за поддержку базы данных, другая — за интеграцию с 
   провайдерами электронной почты. В этом случае они будут работать с одним и 
   тем же кодом.
5. Этот сервис практически невозможно протестировать с помощью unit тестов.
6. ...

Таким образом, давайте проведем рефакторинг этого кода.

## Как соблюсти единую ответственность

Чтобы разделить ответственность в этом случае и создать фрагменты кода, у которых
есть только одна причина для существования, мы должны определить структуру для 
каждого из них.

Практически это означает наличие отдельной структуры для хранения данных в 
некотором хранилище и другой структуры для отправки электронных писем,
используя некоторую интеграцию с провайдерами электронной почты. Более подробно
это описано в коде ниже:

```go
type EmailGorm struct {
    gorm.Model
    From    string
    To      string
    Subject string
    Message string
}

type EmailRepository interface {
    Save(from string, to string, subject string, message string) error
}

type EmailDBRepository struct {
    db *gorm.DB
}

func NewEmailRepository(db *gorm.DB) EmailRepository {
    return &EmailDBRepository{
        db: db,
    }
}

func (r *EmailDBRepository) Save(from string, to string, subject string, message string) error {
    email := EmailGorm{
        From:    from,
        To:      to,
        Subject: subject,
        Message: message,
    }
   
    err := r.db.Create(&email).Error
    if err != nil {
        log.Println(err)
        return err
    }
   
    return nil
}

type EmailSender interface {
    Send(from string, to string, subject string, message string) error
}

type EmailSMTPSender struct {
    smtpHost     string
    smtpPassword string
    smtpPort     int
}

func NewEmailSender(smtpHost string, smtpPassword string, smtpPort int) EmailSender {
    return &EmailSMTPSender{
        smtpHost:     smtpHost,
        smtpPassword: smtpPassword,
        smtpPort:     smtpPort,
    }
}

func (s *EmailSMTPSender) Send(from string, to string, subject string, message string) error {
    auth := smtp.PlainAuth("", from, s.smtpPassword, s.smtpHost)

    server := fmt.Sprintf("%s:%d", s.smtpHost, s.smtpPort)

    err := smtp.SendMail(server, auth, from, []string{to}, []byte(message))
    if err != nil {
        log.Println(err)
        return err
    }

    return nil
}

type EmailService struct {
    repository EmailRepository
    sender     EmailSender
}

func NewEmailService(repository EmailRepository, sender EmailSender) *EmailService {
    return &EmailService{
        repository: repository,
        sender:     sender,
    }
}

func (s *EmailService) Send(from string, to string, subject string, message string) error {
    err := s.repository.Save(from, to, subject, message)
    if err != nil {
        return err
    }

    return s.sender.Send(from, to, subject, message)
}
```

Здесь у нас две новые структуры. Первая `EmailDBRepository` - реализация 
интерфейса `EmailRepository`. Она включает поддержку сохранения данных в 
соответствующей БД.

Вторая структура — это `EmailSMTPSender`, реализующая интерфейс `EmailSender`.
Она отвечает только за отправку электронной почты по протоколу SMTP.

Наконец, новый `EmailService` содержит вышеприведенные интерфейсы и делегирует 
запрос на отправку электронной почты.

Может возникнуть вопрос: почему `EmailService` по-прежнему имеет несколько ответственностей, 
поскольку он по прежнему содержит логику для хранения и отправки электронных 
писем? Похоже, мы только что создали абстракцию, но ответственности остались?

На самом деле это не так. `EmailService` не несет ответственности за хранение и 
отправку электронных писем. Он делегирует их структурам `EmailDBRepository` и
`EmailSMTPSender`. В его ответственность входит делегирование запросов на 
обработку писем соответствующим сервисам.

> В этом заключается разница в случае когда сервис содержит ответственность и
> делегирует её. Если адаптация конкретного кода может удалить всю ответственность,
> то мы говорим о содержании. Если эта ответственность всё равно существует 
> даже после удаления конкретного кода, то мы говорим о делегировании.

Если мы полностью удалим `EmailService`, то у нас по-прежнему будет код, отвечающий
за хранение данных в БД и отправку писем через SMTP. Это означает, что мы
определенно можем сказать, что `EmailService` не берёт на себя эти две 
ответственности.

## Другие примеры

Как видно из примеров выше, SRP применяется ко многим различным аспектам
написания кода, а не только к структурам. Мы видели, что можем его нарушить для
функции, но в том примере, SRP уже было нарушено внутри структуры.

Чтобы лучше понять как SRP принцип применяется к функциям, давайте посмотрим на
пример ниже:

```go
import (
    "github.com/golang-jwt/jwt"
    "net/http"
)

func extractUsername(header http.Header) string {
    raw := header.Get("Authorization")
    parser := &jwt.Parser{}
    token, _, err := parser.ParseUnverified(raw, jwt.MapClaims{})
    if err != nil {
        return ""
    }
    
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return ""
    }
    
    return claims["username"].(string)
}
```

Функция `extractUsername` содержит немного строк. Она извлекает необработанный 
[JWT](https://jwt.io/) токен из HTTP заголовка **и** возвращает значение имени пользователя, если 
оно существует внутри него.

И снова вы могли заметить слово "и", выделенное жирным шрифтом. Неважно как мы 
сформулируем описание, но мы не можем избежать использования слова "и" для 
описания того что делает метод.

Вместо того, чтобы пытаться по-другому сформулировать описание, давайте 
попытаемся по-другому написать сам код метода. Ниже предложен следующий вариант:

```go
import (
    "github.com/golang-jwt/jwt"
    "net/http"
)

func extractRawToken(header http.Header) string {
    return header.Get("Authorization")
}

func extractClaims(raw string) jwt.MapClaims {
    parser := &jwt.Parser{}
    token, _, err := parser.ParseUnverified(raw, jwt.MapClaims{})
    if err != nil {
        return nil
    }
   
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return nil
    }
   
    return claims
}

func extractUsername(header http.Header) string {
    raw := extractRawToken(header)
    claims := extractClaims(raw)
    if claims == nil {
        return ""
    }
	
    return claims["username"].(string)
}
```

Теперь у нас две новые функции. Первая, `extractRawToken`, отвечает за 
извлечение необработанного токена JWT из HTTP заголовка. Если мы изменим ключ
в заголовке, где хранится токен, то нужно внести изменения только в один метод.

Вторая - `extractClaims`. Этот метод отвечает за извлечение `claims` из 
необработанного JWT токена. Наконец, наша старая функция `extractUsername` 
использует конкретное значение из `claims` после делегирования запросов на 
извлечение токена соответствующим методам.

Можно привести другие примеры нарушения SRP. Многие из них регулярно встречаются 
в нашем коде. Иногда некоторые фреймворки навязывают использовать неправильный
подход или мы слишком ленивы, чтобы написать правильную реализацию.

```go
type User struct {
    db        *gorm.DB
    Username  string
    Firstname string
    Lastname  string
    Birthday  time.Time
    //
    // какие-то другие поля
    //
}

func (u User) IsAdult() bool {
    return u.Birthday.AddDate(18, 0, 0).Before(time.Now())
}

func (u User) Save() error {
    return u.db.Exec("INSERT INTO users ...", u.Username, u.Firstname, u.Lastname, u.Birthday).Error
}
```

В приведенном выше примере показана типичная реализация шаблона 
[Active Record](https://en.wikipedia.org/wiki/Active_record_pattern). В нашем
случае мы также добавили бизнес-логику в структуру `User`, а не просто сохранили
данные в БД.

Здесь мы смешали назначение шаблонов Active Record и 
[Сущность](https://levelup.gitconnected.com/practical-ddd-in-golang-entity-40d32bdad2a3)
(`Entity`) из предметно-ориентированного проектирования. Правильно было бы создать
отдельные структуры: одну для сохранения данных в БД, а вторую — для выполнения 
роли `Сущности`. Та же ошибка встречается в примере, показанном ниже:

```go
type Wallet struct {
    gorm.Model
    Amount     int `gorm:"column:amount"`
    CurrencyID int `gorm:"column:currency_id"`
}

func (w *Wallet) Withdraw(amount int) error {
    if amount > w.Amount {
        return errors.New("there is no enough money in wallet")
    }
    
    w.Amount -= amount
    
    return nil
}
```

Здесь у нас опять две ответственности, но теперь вторая (сопоставляющая данным 
таблицу в БД с помощью пакета [Gorm](https://gorm.io/index.html)) выражена не в виде кода, а через 
дескрипторы Go.

Даже сейчас структура `Wallet` нарушает принцип SRP, поскольку играет несколько 
ролей. Если мы меняем схему базы данных, нам нужно модифицировать эту структуру.
Если мы изменим бизнес-правила списания денег, то нам нужно будет скорректировать
этот класс.

```go
type Transaction struct {
	gorm.Model
	Amount     int       `gorm:"column:amount" json:"amount" validate:"required"`
	CurrencyID int       `gorm:"column:currency_id" json:"currency_id" validate:"required"`
	Time       time.Time `gorm:"column:time" json:"time" validate:"required"`
}
```

Вышеприведенный фрагмент кода — ещё один пример нарушения SRP. И, на мой взгляд,
наиболее катастрофичный. Тяжело представить себе структуру меньших размеров с 
ещё большим числом ответственностей.

Взглянув на структуру `Transaction` мы понимаем, что она в ней прописано 
сопоставление с таблицей в БД, она определяет формат JSON ответа в REST API, а
поскольку присутствуют дескрипторы `validate`, она может использоваться для 
проверки JSON тела в API запросе. *Одна структура, выполняющая столько задач*.

> Прим. пер. В оригинале видна отсылка к [Кольцу Всевластия](https://ru.wikipedia.org/wiki/%D0%9A%D0%BE%D0%BB%D1%8C%D1%86%D0%BE_%D0%92%D1%81%D0%B5%D0%B2%D0%BB%D0%B0%D1%81%D1%82%D1%8C%D1%8F)
> из "Властелина Колец".

Во всех этих примерах рано или поздно нужно будет провести рефакторинг. Они 
являются скрытыми проблемами, пока остаются в таком виде, которые скоро начнут 
нарушать нашу логику.

## Заключение

Принцип единой ответственности — первый из SOLID принципов. Он обозначает букву
S в слове SOLID. В нём утверждается, что одна структура кода должна иметь 
только одну причину для существования.

Мы называем эти причины ответственностью. Структура может иметь ответственность
или делегировать её. Всякий раз, когда наша структура содержит несколько 
ответственностей, мы должны осуществить рефакторинг этого фрагмента кода.