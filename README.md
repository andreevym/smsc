# Документация для библиотеки smsc

smsc - это библиотека, предоставляющая возможность отправки SMS-сообщений через HTTP API. Она предоставляет удобные методы для создания и отправки сообщений, а также обработки ответов от сервера.

## Установка

1. Склонируйте репозиторий:

```bash
git clone https://github.com/andreevym/smsc.git
```

2. Перейдите в каталог проекта:
```shell
cd smsc
```

3. Установите зависимости:

```shell
go mod tidy
```

4. Соберите приложение:

```shell
go build
```

## Использование

### Отправка SMS

Пример использования:

```go
package main

import (
	"fmt"

	"github.com/andreevym/smsc"
)

const loginSmscRu = "TODO"          // Логин клиента.
const passwordSmscRu = "TODO"       // Пароль клиента.
const smsSender = "TODO"            // Имя отправителя. (название отправителя)
const smsRecipient = "+71234567890" // Номера телефонов в международном формате.
const smsText = "smsText"           // Текст сообщения.

func main() {
	// Создание клиента (указываем один логин, пароль, название отправителя в смс для получателя)
	client := smsc.NewClient(loginSmscRu, passwordSmscRu, smsSender)

	// Создание сообщения
	message := smsc.NewMessage()
	message.Phones = smsRecipient
	message.Message = smsText

	// Отправка сообщения
	response, err := client.SendMessage(message)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to send message: %w", err))
		return
	}

	// Вывод ответа сервера
	fmt.Println("Ответ сервера:", response)
}
```

Пример запроса в сервис smsc.ru:

```json
{
  "login": "your_login",
  "psw": "your_password",
  "phones": "1234567890",
  "mes": "Hello, World!"
}
```

Пример ответа от сервиса smsc.ru:

```json
{
  "id": 1,
  "cnt": 1,
  "error_code": 0
}
```

### Параметры сообщения

Для отправки SMS поддерживаются различные параметры, такие как:

- `login`: Логин клиента.
- `psw`: Пароль клиента.
- `phones`: Номера телефонов получателей.
- `mes`: Текст сообщения.
- `list`: Список номеров телефонов и соответствующих им сообщений, 
разделенных двоеточием или точкой с запятой и представленный в виде: `phones1:mes1`, `phones2:mes2`

Обязательными параметрами являются `login`, `psw`, `phones` и `mes` либо `login`, `psw` и `list`.

Полный список параметров и их описание [смотрите в документации](https://smsc.ru/api/http/#menu).

### Валидация

Перед отправкой сообщения происходит валидация всех параметров. В случае ошибок валидации возвращается соответствующая ошибка.

## Дополнительные ресурсы

- [Документация API: Полное описание всех параметров и методов API.](https://smsc.ru/api/http/#menu)
- [Примеры использования: Набор примеров кода для быстрого старта.](./example/main.go)
- [Исходный код: Исходный код библиотеки на GitHub.](https://github.com/andreevym/smsc)

## Лицензия
Этот проект лицензирован под [MIT License](LICENSE).