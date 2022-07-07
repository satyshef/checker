package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/satyshef/checker/internal/config"
	"github.com/satyshef/tdbot"
	"github.com/satyshef/tdbot/mimicry"
	"github.com/satyshef/tdlib"
)

func completeTask(b *tdbot.Bot) *tdlib.Error {
	//TODO : Mimicry tempolary disable
	/*
		human := mimicry.NewHuman(b)
		replyToMessage(human)
		sendFriendMessage(human)
	*/

	//Рассылка
	for _, m := range conf.Mailings {
		if !m.Enable {
			continue
		}
		var msg string
		if _, err := os.Stat(m.Message); err == nil {
			msg = loadRandomString(m.Message)
		} else {
			msg = m.Message
		}

		if msg == "" {
			break
		}
		fmt.Printf("Send to %s : %s\n\n", m.Chat, msg)
		_, err := bot.SendMessageToChat(m.Chat, msg, m.Leave)
		if err != nil {
			//return err
			fmt.Printf("Send to chat error : %s\n", err)
			break
		}
	}

	// Жалобы
	for _, r := range conf.Reports {
		if !r.Enable {
			continue
		}
		sendReport(r)
	}

	return nil
}

/*
//Получить входящие сообщения, если есть от друзей то ответить(РАБОТАЕТ ТОЛЬКО С ТЕМИ СООБЩЕНИЯ ЧТО ПРИШЛИ КОГДА БОТ ЗАПУЩЕН)
func receiveMessage(h *mimicry.Human, friends []int32) {
	var messageDummy tdlib.UpdateNewMessage

	rawUpdates := h.Bot.Client.GetRawUpdatesChannel(1)

	for update := range rawUpdates {

		// Show all updates
		//updateNewMessage
		fmt.Printf("%#v\n", update.Data)
		//Если новое сообщение от друга то отвечаем на него
		if update.Data["@type"] == "updateNewMessage" {
			json.Unmarshal(update.Raw, &messageDummy)
			for _, uid := range friends {
				if uid == messageDummy.Message.Sender.(*tdlib.MessageSenderUser).UserID {
					//Отправить
					msg := "Ho ho ho"
					h.Bot.SendMessageByUID(uid, msg, 0)
					break
				}
			}
			//uid := messageDummy.Message.Sender.(*tdlib.MessageSenderUser).UserID
			//fmt.Printf("%#v\n", messageDummy.Message)
			//fmt.Println(messageDummy.Message.Content.GetMessageContentEnum())

		}
	}
}
*/

func sendReport(r config.Report) *tdlib.Error {

	chat, err := bot.GetChat(r.Chat, false)
	if err != nil {
		return err
	}
	if chat.LastMessage == nil {
		return tdlib.NewError(tdbot.ErrorCodeSystem, "LASTMESSAGE_NOT_LOAD", "")
	}
	var msg string
	if _, err := os.Stat(r.Message); err == nil {
		msg = loadRandomString(r.Message)
	} else {
		msg = r.Message
	}
	if msg == "" {
		return tdlib.NewError(tdbot.ErrorCodeSystem, "EMPTY_MESSAGE", "")
	}
	reason := tdlib.NewChatReportReasonCustom(msg)
	//reason := tdlib.NewChatReportReasonSpam()
	_, e := bot.Client.ReportChat(chat.ID, reason, []int64{chat.LastMessage.ID})
	if e != nil {
		return e.(*tdlib.Error)
	}
	fmt.Println("Report OK")
	return nil
}

//Отправить сообщение одному из друзей
func sendFriendMessage(h *mimicry.Human) {
	// Выбрать случайно друга
	uid, err := h.GetRandomFriend()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Friend ID ", uid)

	//Получить текст сообщения
	msg := generateMessageForFriend()
	//Отправить
	h.Bot.SendMessageByUID(uid, msg, 0)

}

func generateMessageForFriend() string {

	var result string
	countWord := mimicry.RandInt(5, 10)
	for i := 0; i < countWord; i++ {
		countLetter := mimicry.RandInt(3, 9)
		result += " " + mimicry.RandString(countLetter)
	}
	return strings.Trim(result, " ")

}

func replyToMessage(h *mimicry.Human) {

	chatList, _ := h.Bot.GetChatList(100)

	fmt.Printf("M %#v\n\n", chatList)
	for _, c := range chatList {
		if c.UnreadCount != 0 && c.Type.GetChatTypeEnum() == tdlib.ChatTypePrivateType {
			h.Bot.Client.ViewMessages(c.ID, 0, []int64{c.LastMessage.ID}, true)
			msg := generateMessageForFriend()
			h.Bot.SendMessageByUID(c.ID, msg, 0)
			//fmt.Printf("\n%#v\n\n", c)
		}

	}

}

func loadRandomString(fileName string) string {

	lines, err := readFileToSlice(fileName)
	if err != nil {
		log.Fatal(err)
	}
	return shuffleArray(lines)[0]
}

func readFileToSlice(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := strings.Trim(scanner.Text(), " \n\t")
		lines = append(lines, text)
	}

	if scanner.Err() != nil {
		return nil, err
	}

	if len(lines) == 0 {
		return nil, fmt.Errorf("Empty file %s", fileName)
	}

	return lines, nil
}

func shuffleArray(src []string) []string {
	final := make([]string, len(src))
	rand.Seed(time.Now().UTC().UnixNano())
	perm := rand.Perm(len(src))

	for i, v := range perm {
		final[v] = src[i]
	}
	return final
}
