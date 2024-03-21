package smsc

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"time"
)

// Message представляет структуру сообщения.
// Обязательными параметрами являются login, psw, phones и mes либо login, psw и list.
type Message struct {
	// Обязательные параметры
	Login    string `json:"login"`            // Логин клиента.
	Password string `json:"psw"`              // Пароль клиента.
	Phones   string `json:"phones,omitempty"` // Номера телефонов в международном формате.

	// Обязательные либо list
	Message string `json:"mes,omitempty"`  // Текст сообщения.
	List    string `json:"list,omitempty"` // Список сообщений, если не указан текст сообщения.

	// Дополнительные параметры
	ID       string  `json:"id,omitempty"`       // Идентификатор сообщения.
	Sender   string  `json:"sender,omitempty"`   // Имя отправителя.
	Translit int     `json:"translit,omitempty"` // Признак транслитерации.
	TinyURL  int     `json:"tinyurl,omitempty"`  // Признак автоматического сокращения ссылок.
	Time     string  `json:"time,omitempty"`     // Время отправки сообщения.
	TZ       int     `json:"tz,omitempty"`       // Часовой пояс.
	Period   float64 `json:"period,omitempty"`   // Промежуток времени для рассылки.
	Freq     int     `json:"freq,omitempty"`     // Интервал рассылки.
	Flash    int     `json:"flash,omitempty"`    // Признак Flash сообщения.
	Bin      int     `json:"bin,omitempty"`      // Признак бинарного сообщения.
	Push     int     `json:"push,omitempty"`     // Признак WAP-push сообщения.
	HLR      int     `json:"hlr,omitempty"`      // Признак HLR-запроса.
	Ping     int     `json:"ping,omitempty"`     // Признак Ping-SMS.
	MMS      int     `json:"mms,omitempty"`      // Признак MMS-сообщения.
	Mail     int     `json:"mail,omitempty"`     // Признак e-mail сообщения.
	Soc      int     `json:"soc,omitempty"`      // Признак soc-сообщения.
	Viber    int     `json:"viber,omitempty"`    // Признак Viber-сообщения.
	WhatsApp int     `json:"whatsapp,omitempty"` // Признак WhatsApp-сообщения.
	Bot      string  `json:"bot,omitempty"`      // Имя бота Telegram.
	SMSReq   int     `json:"smsreq,omitempty"`   // SMS-запрос.
	FileURL  string  `json:"fileurl,omitempty"`  // URL файла.
	Call     int     `json:"call,omitempty"`     // Признак голосового сообщения.
	Voice    string  `json:"voice,omitempty"`    // Голос.
	Param    string  `json:"param,omitempty"`    // Параметры голосового сообщения.
	Subject  string  `json:"subj,omitempty"`     // Тема MMS или e-mail сообщения.
	Charset  string  `json:"charset,omitempty"`  // Кодировка сообщения.
	Cost     int     `json:"cost,omitempty"`     // Стоимость рассылки.
	Format   int     `json:"fmt,omitempty"`      // Формат ответа сервера.
	Valid    string  `json:"valid,omitempty"`    // Срок "жизни" сообщения.
	MaxSMS   int     `json:"maxsms,omitempty"`   // Максимальное количество SMS.
	ImgCode  string  `json:"imgcode,omitempty"`  // Код "captcha".
	UserIP   string  `json:"userip,omitempty"`   // IP-адрес.
	Err      int     `json:"err,omitempty"`      // Признак добавления списка ошибочных номеров в ответ.
	Op       int     `json:"op,omitempty"`       // Признак добавления информации по каждому номеру в ответ.
	PP       string  `json:"pp,omitempty"`       // Привязка клиента к партнеру.
}

// JSON метод для преобразования сообщения в JSON
func (m *Message) JSON() (string, error) {
	// Преобразовываем сообщение в JSON
	jsonData, err := json.Marshal(m)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// Validate проверяет корректность значений полей структуры Message.
func (m *Message) Validate() error {
	err := validateRequiredFields(m)
	if err != nil {
		return fmt.Errorf("message validation: %w", err)
	}

	// Проверка времени
	if m.Time != "" {
		if err := validateTime(m.Time); err != nil {
			return fmt.Errorf("неверное время отправки сообщения: %v", err)
		}
	}

	if m.Cost != 0 && (m.Cost < 0 || m.Cost > 3) {
		return errors.New("неверное значение стоимости")
	}
	if m.Format != 0 && (m.Format < 0 || m.Format > 3) {
		return errors.New("неверное значение формата ответа сервера")
	}
	if m.Valid != "" {
		validRegex := regexp.MustCompile(`^\d{1,2}:\d{2}$`)
		if !validRegex.MatchString(m.Valid) {
			return errors.New("неверный формат срока жизни сообщения, используйте формат чч:мм")
		}
	}
	if m.MaxSMS != 0 && (m.MaxSMS != 0 && m.MaxSMS < 0) {
		return errors.New("максимальное количество SMS должно быть неотрицательным")
	}
	if m.Time != "" {
		if err := validateTime(m.Time); err != nil {
			return err
		}
	}
	if m.TZ != 0 && (m.TZ < -12 || m.TZ > 12) {
		return errors.New("неверное значение часового пояса")
	}
	if m.Period != 0 && (m.Period < 0.1 || m.Period > 720) {
		return errors.New("промежуток времени должен быть в диапазоне от 0.1 до 720 часов")
	}
	if m.Freq != 0 && (m.Freq < 1 || m.Freq > 1440) {
		return errors.New("интервал отправки должен быть в промежутке от 1 до 1440 минут")
	}
	if m.Flash != 0 && (m.Flash < 0 || m.Flash > 1) {
		return errors.New("неверное значение признака Flash сообщения")
	}
	if m.Bin != 0 && (m.Bin < 0 || m.Bin > 2) {
		return errors.New("неверное значение признака бинарного сообщения")
	}
	if m.Push != 0 && (m.Push < 0 || m.Push > 1) {
		return errors.New("неверное значение признака WAP-push сообщения")
	}
	if m.HLR != 0 && (m.HLR < 0 || m.HLR > 1) {
		return errors.New("неверное значение признака HLR-запроса")
	}
	if m.Ping != 0 && (m.Ping < 0 || m.Ping > 1) {
		return errors.New("неверное значение признака Ping-SMS")
	}
	if m.MMS < 0 || m.MMS > 1 {
		return errors.New("неверное значение признака MMS-сообщения")
	}
	if m.Mail < 0 || m.Mail > 1 {
		return errors.New("неверное значение признака e-mail сообщения")
	}
	if m.Soc < 0 || m.Soc > 1 {
		return errors.New("неверное значение признака soc-сообщения")
	}
	if m.Viber < 0 || m.Viber > 1 {
		return errors.New("неверное значение признака viber-сообщения")
	}
	if m.WhatsApp < 0 || m.WhatsApp > 1 {
		return errors.New("неверное значение признака whatsapp-сообщения")
	}
	// Проверка дополнительных параметров
	if m.SMSReq != 0 && (m.SMSReq < 10 || m.SMSReq > 999) {
		return errors.New("неверное значение для параметра smsreq, должно быть в диапазоне от 10 до 999")
	}
	if m.FileURL != "" && len(m.FileURL) < 101 {
		return errors.New("неверный http-адрес файла, минимальный размер файла должен быть 101 байт")
	}
	if m.Call < 0 || m.Call > 1 {
		return errors.New("неверное значение признака голосового сообщения")
	}
	if m.Voice != "" && !isValidVoice(m.Voice) {
		return errors.New("неверное значение голоса")
	}
	if m.Param != "" && !isValidParam(m.Param) {
		return errors.New("неверное значение параметров голосового сообщения")
	}
	if m.Valid != "" {
		if err := validateValid(m.Valid); err != nil {
			return fmt.Errorf("неверное значение срока 'жизни' сообщения: %v", err)
		}
	}
	return nil
}

func validateRequiredFields(m *Message) error {
	// Проверка обязательных полей
	if m.Login == "" || m.Password == "" {
		return fmt.Errorf("логин и пароль являются обязательными параметрами")
	}
	if m.Phones == "" && m.List == "" {
		return fmt.Errorf("номера телефонов (phones) или список сообщений (list) являются обязательными параметрами")
	}
	if m.List == "" && m.Message == "" {
		return fmt.Errorf("текст сообщения (mes) или список сообщений (list) являются обязательными параметрами")
	}
	if m.Message != "" && len(m.Message) > 1000 {
		return errors.New("максимальная длина текста сообщения - 1000 символов")
	}
	return nil
}

func validateTime(timeStr string) error {
	timeFormats := []string{
		"0201061504",
		"02.01.06 15:04",
		"15-17",
		"0ts",
		"+5",
	}

	for _, format := range timeFormats {
		_, err := time.Parse(format, timeStr)
		if err == nil {
			return nil
		}
	}

	return errors.New("неверный формат времени")
}

func isValidVoice(voice string) bool {
	validVoices := map[string]bool{
		"m":  true,
		"m2": true,
		"m3": true,
		"m4": true,
		"w":  true,
		"w2": true,
		"w3": true,
		"w4": true,
	}

	_, ok := validVoices[voice]
	return ok
}

func isValidParam(param string) bool {
	validParamRegex := regexp.MustCompile(`^\d+,\d+,\d+$`)
	return validParamRegex.MatchString(param)
}

func validateValid(valid string) error {
	if valid == "" {
		return nil
	}

	if valid == "0" {
		return nil
	}

	validTime, err := time.Parse("15:04", valid)
	if err != nil {
		return errors.New("неверный формат срока 'жизни' сообщения")
	}

	if validTime.Hour() < 1 || validTime.Hour() > 24 || validTime.Minute() < 0 || validTime.Minute() > 59 {
		return errors.New("неверное значение срока 'жизни' сообщения")
	}

	return nil
}

func NewMessage() *Message {
	return &Message{}
}
