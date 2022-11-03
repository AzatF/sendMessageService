package main

import (
	"bytes"
	"flag"
	"fmt"
	gomail "gopkg.in/mail.v2"
	"html/template"
	"log"
	"mail/internal/base"
	"mail/internal/config"
	"mail/pkg/logging"
	"net/smtp"
	"os"
	"path"
	"strconv"
	"time"
)

var cfgPath string

func init() {
	flag.StringVar(&cfgPath, "config", path.Join("etc", ".env"), "config file path")
}

func main() {

	cfg := config.GetConfig(cfgPath)
	logger := logging.GetLogger("trace")
	db, err := base.NewDataBase(cfg, logger)
	if err != nil {
		logger.Error(err)
	}

	f, err := os.ReadFile(path.Join(cfg.DataPath, "newUsers.data"))
	if err != nil {
		logger.Error(err)
	}

	err = db.AddNewSubscribers(f)
	if err != nil {
		logger.Error(err)
	}

	users, err := db.GetSubscribers()
	if err != nil {
		logger.Error(err)
	}

	for _, v := range users {
		fmt.Printf("Имя: %s\tАдрес почты: %s\tМесяц рождения: %s\n", v.FirstName, v.Email, v.BirthDay.Month())
	}

	//------------------------- Birth Day Mail  ---------------------------------------------------

	//user := base.Subscribers{
	//	Id:         1,
	//	FirstName:  "Василий",
	//	SecondName: "Иванов",
	//	FatherName: "Павлович",
	//	Email:      "sardelka@mail.ru",
	//	BirthDay:   time.Now().UTC(),
	//	Sex:        "male",
	//}
	//subject := "Поздравление!"
	//first := "Дорогой друг " + user.FirstName
	//second := "В этот знаменательный день... " + user.BirthDay.Format(config.StructDateFormat)
	//third := "Псс, подарки нннада?"
	//fourth := "Пошли..."
	//sendBirthDayMail(cfg, user, first, second, third, fourth, cfg.RecipientEmail, subject)

	//------------------------  Info Mail ------------------------------------------------------------

	user := base.Subscribers{
		Id:         1,
		FirstName:  "Василий",
		SecondName: "Иванов",
		FatherName: "Павлович",
		Email:      "sardelka@mail.ru",
		BirthDay:   time.Now().UTC(),
		Sex:        "male",
	}
	subject := "Информация!"
	first := "В своём стремлении повысить качество жизни, они забывают, " +
		"что разбавленное изрядной долей эмпатии, рациональное мышление предопределяет высокую " +
		"востребованность поставленных обществом задач. Принимая во внимание показатели успешности, " +
		"современная методология разработки требует определения и уточнения системы массового участия. " +
		"Сложно сказать, почему представители современных социальных резервов будут преданы социально-демократической анафеме."
	second := "Приятно, граждане, наблюдать, как диаграммы связей, инициированные исключительно синтетически, призваны к ответу! " +
		"С другой стороны, синтетическое тестирование способствует подготовке и реализации первоочередных требований. " +
		"Лишь активно развивающиеся страны третьего мира и по сей день остаются уделом либералов, которые жаждут быть призваны к ответу."
	third := "Не следует, однако, забывать, что внедрение современных методик способствует подготовке и реализации форм воздействия. " +
		"Также как сплочённость команды профессионалов представляет собой интересный эксперимент проверки кластеризации усилий"
	fourth := "С наилучшими пожеланиями!"
	infoMail(cfg, user, first, second, third, fourth, cfg.RecipientEmail, subject)

}

func infoMail(cfg *config.Config, user base.Subscribers, first, second, third, fourth, recipient, subject string) {

	var b bytes.Buffer
	t, err := template.ParseFiles("info.html")
	if err != nil {
		log.Println("error parse html: ", err)
		return
	}

	err = t.Execute(&b, struct{ Name, First_text, Second_text, Third_text, Fourth_text string }{
		Name:        user.FirstName,
		First_text:  first,
		Second_text: second,
		Third_text:  third,
		Fourth_text: fourth,
	})
	if err != nil {
		return
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", cfg.SenderEmail)
	msg.SetHeader("To", user.Email)
	msg.SetHeader("Subject", subject)
	msg.SetAddressHeader("To", user.Email, "Получатель")
	msg.SetBody("text/html", b.String())

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.SenderEmail, cfg.SenderPass)

	if err = d.DialAndSend(msg); err != nil {
		log.Println("error send message: ", err)
		return
	}

}

func sendBirthDayMail(cfg *config.Config, user base.Subscribers, first, second, third, fourth, recipient, subject string) {

	var b bytes.Buffer
	t, err := template.ParseFiles("birthday2.html")
	if err != nil {
		log.Println("error parse html: ", err)
		return
	}

	err = t.Execute(&b, struct{ First_text, Second_text, Third_text, Fourth_text string }{
		First_text:  first,
		Second_text: second,
		Third_text:  third,
		Fourth_text: fourth,
	})
	if err != nil {
		return
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", cfg.SenderEmail)
	msg.SetHeader("To", user.Email)
	msg.SetHeader("Subject", subject)
	msg.SetAddressHeader("To", user.Email, "Получатель")
	msg.SetBody("text/html", b.String())

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.SenderEmail, cfg.SenderPass)

	if err = d.DialAndSend(msg); err != nil {
		log.Println("error send message: ", err)
		return
	}

}

func sendGoMail(cfg *config.Config, htmlModel string, recipient string, subject string, body string) {

	var b bytes.Buffer
	t, err := template.ParseFiles(htmlModel)
	if err != nil {
		log.Println("error parse html: ", err)
		return
	}
	//t.Execute(&b, struct{ Name, Text string }{Name: "Здравствуй Азат!", Text: body})
	t.Execute(&b, struct{ First_text, Second_text string }{First_text: "Здравствуй Азат!", Second_text: body})

	msg := gomail.NewMessage()
	msg.SetHeader("From", cfg.SenderEmail)
	msg.SetHeader("To", recipient)
	msg.SetHeader("Subject", subject)
	msg.SetAddressHeader("To", cfg.RecipientEmail, "Получатель")
	msg.SetBody("text/html", b.String())
	//msg.Attach("badcode.jpg")

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.SenderEmail, cfg.SenderPass)

	if err = d.DialAndSend(msg); err != nil {
		log.Println("error send message: ", err)
		return
	}

}

func mailSendHtml(cfg *config.Config, htmlModel string, recipient []string, subject string, body string) {

	var b bytes.Buffer
	t, err := template.ParseFiles(htmlModel)
	if err != nil {
		log.Println("error parse html: ", err)
		return
	}
	err = t.Execute(&b, struct{ Name, Text string }{Name: "Азат!", Text: body})
	if err != nil {
		return
	}

	address := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	auth := smtp.PlainAuth(cfg.RecipientEmail, cfg.SenderEmail, cfg.SenderPass, cfg.Host)
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	//message := "Subject: " + subject + "\r\n" + body
	message := "Subject: " + subject + "\r\n" + headers + "\n\n" + b.String()

	err = smtp.SendMail(address, auth, cfg.RecipientEmail,
		recipient, []byte(message))
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
}

func mailSend(cfg *config.Config, recipient []string, subject string, body string) {

	address := cfg.Host + ":" + strconv.Itoa(cfg.Port)
	auth := smtp.PlainAuth(cfg.RecipientEmail, cfg.SenderEmail, cfg.SenderPass, cfg.Host)
	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	//message := "Subject: " + subject + "\r\n" + body
	message := "Subject: " + subject + "\r\n" + headers + "\n\n" + body

	err := smtp.SendMail(address, auth, cfg.RecipientEmail,
		recipient, []byte(message))
	if err != nil {
		fmt.Println("error: ", err)
		return
	}
}
