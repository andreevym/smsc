package smsc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Константы для типов ошибок
const (
	ErrorBadRequest          = 1
	ErrorInvalidCredentials  = 2
	ErrorInsufficientFunds   = 3
	ErrorIPBlocked           = 4
	ErrorInvalidDateFormat   = 5
	ErrorMessageForbidden    = 6
	ErrorInvalidPhoneNumber  = 7
	ErrorMessageNotDelivered = 8
	ErrorDuplicateRequest    = 9
)

var errorDescriptions = map[int]string{
	ErrorBadRequest:          "Ошибка в параметрах",
	ErrorInvalidCredentials:  "Неверный логин или пароль",
	ErrorInsufficientFunds:   "Недостаточно средств на счете",
	ErrorIPBlocked:           "IP-адрес временно заблокирован",
	ErrorInvalidDateFormat:   "Неверный формат даты",
	ErrorMessageForbidden:    "Сообщение запрещено",
	ErrorInvalidPhoneNumber:  "Неверный формат номера телефона",
	ErrorMessageNotDelivered: "Сообщение не может быть доставлено",
	ErrorDuplicateRequest:    "Дублирование запроса",
}

// Response структура для ответа сервера
type Response struct {
	ID        int        `json:"id"`         // идентификатор сообщения, переданный Клиентом или назначенный Сервером автоматически.
	Cnt       int        `json:"cnt"`        // количество частей (при отправке SMS-сообщения) либо количество секунд (при голосовом сообщении (звонке)).
	Cost      string     `json:"cost"`       // стоимость SMS-сообщения.
	Balance   string     `json:"balance"`    // новый баланс Клиента.
	Error     string     `json:"error"`      // текст ошибки
	ErrorCode int        `json:"error_code"` // код ошибки в статусе.
	Phones    []struct { // номер телефона.
		Phone  string `json:"phone"`  // номер телефона.
		Mccmnc string `json:"mccmnc"` // числовой код страны абонента плюс числовой код оператора абонента.
		Cost   string `json:"cost"`   // стоимость SMS-сообщения.
		Status string `json:"status"` // код статуса SMS-сообщения.
		Error  string `json:"error"`  // код ошибки в статусе.
	} `json:"phones"`
}

// Константа с адресом для отправки сообщений
const smscURL = "https://smsc.ru/rest/send/"

// Client клиент для отправки SMS через сервис smsc.ru
type Client struct {
	login    string // Логин клиента.
	password string // Пароль клиента.
	sender   string // Имя отправителя. (название отправителя)
}

// NewClient создание нового клиента
func NewClient(login string, password string, sender string) *Client {
	return &Client{
		login:    login,
		password: password,
		sender:   sender,
	}
}

// SendMessage метод для отправки сообщения
func (c *Client) SendMessage(message *Message) (Response, error) {
	if message == nil {
		return Response{}, errors.New("message for send can't be nil")
	}

	if message.Login == "" && message.Password == "" {
		message.Login = c.login
		message.Password = c.password
	}

	if err := message.Validate(); err != nil {
		message.Login = "do not show real login for err and log"
		message.Password = "do not show real password for err and log"
		json, _ := message.JSON()
		return Response{}, fmt.Errorf("failed to validate sending message '%s': %w", json, err)
	}

	// Преобразование сообщения в JSON
	jsonData, err := message.JSON()
	if err != nil {
		return Response{}, err
	}

	// Формирование запроса POST
	req, err := http.NewRequest(http.MethodPost, smscURL, strings.NewReader(jsonData))
	if err != nil {
		return Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Отправка запроса
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	body := resp.Body
	defer body.Close()

	bytes, err := io.ReadAll(body)
	if err != nil {
		return Response{}, err
	}
	if len(bytes) == 0 {
		return Response{}, errors.New("empty response from server")
	}

	// Проверка статуса ответа сервера
	if resp.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("неправильный статус ответа от сервера: %d", resp.StatusCode)
	}

	// Обработка ответа сервера
	var response Response
	if err := json.Unmarshal(bytes, &response); err != nil {
		return Response{}, fmt.Errorf("failed to parse json: '%s', %w", string(bytes), err)
	}

	// Проверка наличия ошибки в ответе
	if len(response.Error) > 0 || response.ErrorCode > 0 {
		return Response{}, errors.New(response.Error)
	}

	// Проверка наличия ошибки в ответе
	for _, phone := range response.Phones {
		if len(phone.Error) != 0 {
			return Response{}, errors.New(phone.Error)
		}
	}

	return response, nil
}

// ParseErrorCode Метод для парсинга кода ошибки
func ParseErrorCode(errorCode int) string {
	return errorDescriptions[errorCode]
}
