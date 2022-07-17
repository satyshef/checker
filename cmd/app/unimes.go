// обработка UniversalMessage
package main

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/satyshef/mslib/unimes"
	"github.com/satyshef/tdbot"
	"github.com/satyshef/tdlib"
)

func generateUnimes(b *tdbot.Bot, msg *tdlib.Message) (*unimes.UniversalMessage, error) {
	var text string
	if msg.Content == nil {
		bot.Logger.Errorf("%#v\n\n", msg)
		return nil, fmt.Errorf("%s", "Message content empty")
	}
	//Определяем тип контента
	switch msg.Content.GetMessageContentEnum() {
	case tdlib.MessagePhotoType:
		text = msg.Content.(*tdlib.MessagePhoto).Caption.Text
	case tdlib.MessageVideoType:
		text = msg.Content.(*tdlib.MessageVideo).Caption.Text
	case tdlib.MessageTextType:
		text = msg.Content.(*tdlib.MessageText).Text.Text
	default:
		fmt.Printf("MESSAGE TYPE: %#v\n\n", msg.Content.GetMessageContentEnum())
		fmt.Printf("UNKNOW MESSAGE: \n %#v\n\n", msg.Content)
		return nil, fmt.Errorf("%s", "UNKNOW MESSAGE")
	}
	//Формируем unimes сообщение
	/*
		m, err := bot.Client.GetMessage(msg.Message.ChatID, msg.Message.ID)
		if err != nil {
			return nil, err
		}
	*/
	//fmt.Printf("SENDER: \n %#v\n\n", msg.Message.Sender)

	sender, err := generateSender(bot, msg.Sender)
	if err != nil {
		return nil, err
	}
	var locale unimes.Destination
	//Если совпадают ID чата и отправителя значит Locale совподает с Sender
	if msg.ChatID == sender.ID {
		locale = sender
	} else {
		locale, err = generateLocale(bot, msg)
		if err != nil {
			return nil, err
		}
	}
	recipient := generateDestinationFromBot(b)
	uniMessage := &unimes.UniversalMessage{
		ID:        msg.ID,
		Sender:    sender,
		Locale:    locale,
		Recipient: recipient,
		Date:      msg.Date,
		Content:   unimes.Content{Type: string(tdlib.MessageTextType), Data: text},
	}
	return uniMessage, nil
}

//генерируем адрес отправителя сообщения
func generateSender(bot *tdbot.Bot, sender tdlib.MessageSender) (unimes.Destination, error) {
	if sender == nil {
		return unimes.Destination{}, fmt.Errorf("%s", "Nil sender value")
	}
	switch sender.GetMessageSenderEnum() {
	case tdlib.MessageSenderChatType:
		return generateDestinationFromChat(bot, sender.(*tdlib.MessageSenderChat).ChatID)
	case tdlib.MessageSenderUserType:
		return generateDestinationFromUser(bot, sender.(*tdlib.MessageSenderUser).UserID)
	}
	return unimes.Destination{}, fmt.Errorf("%s", "Unknown sender type")
}

//генерируем адрес чата где было получено сообщение
func generateLocale(bot *tdbot.Bot, message *tdlib.Message) (unimes.Destination, error) {
	return generateDestinationFromChat(bot, message.ChatID)
}

func generateDestinationFromUser(bot *tdbot.Bot, uid int64) (unimes.Destination, error) {
	var err error
	result := unimes.Destination{}
	user, err := bot.Client.GetUser(uid)
	if err != nil {
		bot.Logger.Errorf("Get User Error : %s[%d]", err, uid)
		return result, err
	}
	result.ID = user.ID
	result.Type = string(user.Type.GetUserTypeEnum())
	result.Service = "telegram"
	result.FirstName = user.FirstName
	result.Lastname = user.LastName
	result.Username = user.Username
	return result, nil
}

//генерируем адрес чата где было получено сообщение
func generateDestinationFromChat(bot *tdbot.Bot, cid int64) (unimes.Destination, error) {
	result := unimes.Destination{}
	chat, err := bot.GetChatFullInfo(cid)
	if err != nil {
		return result, err
	}
	if chat == nil {
		return result, fmt.Errorf("%s", "NIL chat")
	}
	result.ID = chat.ID
	result.Type = string(chat.Type)
	result.Service = "telegram"
	result.FirstName = chat.Name
	result.Lastname = chat.BIO
	result.Username = chat.Address
	return result, nil
}

func generateDestinationFromBot(bot *tdbot.Bot) unimes.Destination {
	result := unimes.Destination{}
	result.ID = bot.Profile.User.ID
	result.Type = string(bot.Profile.User.Type)
	result.Service = "telegram"
	result.FirstName = bot.Profile.User.FirstName
	result.Lastname = bot.Profile.User.LastName
	result.Username = bot.Profile.User.Addr
	return result
}

func send(url string, data []byte) {

	//values := map[string]string{"text": "From body text", "occupation": "gardener"}
	_, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Send ERROR :", err)
	}

	// TODO: сдклать обработку ответов
	/*
		if resp.StatusCode == http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			fmt.Printf("%s\n\n", bodyString)
		}
	*/
	/*
		var res map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&res)
		fmt.Println("Response", res["json"])
	*/

}
