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
		fmt.Println(fmt.Errorf("failed to send message: '%s', %w", smsc.ParseErrorCode(response.ErrorCode), err))
		return
	}

	// Вывод ответа сервера
	fmt.Println("Ответ сервера:", response)
}
