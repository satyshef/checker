// Чекер профилей телеграм
package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/satyshef/checker/internal/config"
	"github.com/satyshef/tdbot"
	"github.com/satyshef/tdbot/profile"
	"github.com/satyshef/tdlib"
)

var (
	configPath  string
	conf        *config.Config
	profileDir  string
	useMimicry  bool
	bot         *tdbot.Bot
	runProcess  bool
	stopProcess chan bool
	countGood   int
	lockAccs    []string
	interval    int
)

func init() {
	flag.StringVar(&configPath, "c", "./data/config.toml", "Файл конфигурации")
	flag.StringVar(&profileDir, "p", "./profiles", "Путь к директории с профилями")
	flag.BoolVar(&useMimicry, "m", false, "Использовать мимикрию")
	flag.IntVar(&interval, "i", 0, "Интервал перебора профилей (сек)")
	flag.Parse()
	conf = config.LoadConfig(configPath)
	stopProcess = make(chan bool)
	profile.AddTail(&profileDir)
}

func main() {

	//если указанный путь и есть профиль
	if profile.IsProfile(profileDir) {
		err := checkProf(profileDir)
		if err == nil {
			fmt.Println("OK")
		} else {
			fmt.Println(err)
		}
		return
	}

	//Получаем список профилей в алфовитном порядке
	profList := profile.GetList(profileDir, profile.SORT_TIME_ASC)
	for n, phone := range profList {
		fmt.Printf("#%d\n", n+1)
		err := checkProf(profileDir + phone)
		if err == nil {
			fmt.Printf("Success\n\n\n")
		} else {
			if err.Error() == "Profile is already in use" {
				lockAccs = append(lockAccs, profileDir+phone)
			}
		}
	}

	fmt.Println("Finish:")
	fmt.Println("Locked:")
	for _, a := range lockAccs {
		fmt.Println(a)
	}

	fmt.Printf("Good - %d\n All - %d\n", countGood, len(profList))
}

func checkProf(profDir string) error {
	bot = nil
	runProcess = true
	//attempts := 0
	prof, err := profile.Get(profDir, 1)
	if err != nil {
		fmt.Println("PROF ERROR : ", err.Error())
		return err
	}
	defer prof.Close()

	// Init bot
	bot = tdbot.New(prof)
	bot.Client.AddEventHandler(eventCatcher)
	/*
		go func() {
			for {
				attempts++
				if attempts > 3 {
					break
				}
				e := bot.Start()
				if e != nil {
					fmt.Printf("START ERROR : %#v\n\n", e)

					switch e.Code {
					case tdlib.ErrorCodeLogout:
						bot.Logger.Errorln("Profile logout")
						bot.ProfileToLogout()
						goto Exit
					//Если таймаут делаем еще одну попытку
					case tdlib.ErrorCodeTimeout:
						bot.Logger.Errorf("TIMEOUT : %#v\n", e)
						continue
					case profile.ErrorCodeDirNotExists:
						goto Exit
					default:
						//fmt.Printf("START ERROR : %#v\n\n", e)
						bot.Logger.Errorf("START : %#v\n", e)
						goto Exit
					}
				} else {
					// Если бот запустился тогда стопаем его
				}
			Exit:
				//runProcess = false
				stopProcess <- true
				//	break

			}

		}()
	*/

	e := bot.Start()
	if e != nil {
		fmt.Printf("START ERROR : %#v\n\n", e)
	} else {
		//fmt.Println("SUPER")
	}
	//<-stopProcess

	/*
		for bot != nil {
			fmt.Println(bot.Status)
			time.Sleep(time.Second * 1)
		}

		for runProcess || !bot.IsRun() {
			time.Sleep(time.Second * 1)
		}
		//<-stopProcess
		fmt.Println("Good")

		if useMimicry {
			// Run task
			completeTask(bot)
		}

		countGood++
		bot.Stop()
		//runProcess = false
	*/
	/*
		if bot.Status == tdbot.StatusReady {
			break
		}
	*/
	//fmt.Printf("STATUS %#v\n", bot.Status)
	//time.Sleep(time.Second * 1)

	return nil
}

// обработчик событий Телеграм клиента
func eventCatcher(tdEvent *tdlib.SystemEvent) *tdlib.Error {
	//bot.Logger.Errorf("New Event %#v\n\n", tdEvent)
	switch tdEvent.Type {
	case tdlib.EventTypeRequest:
		return requestHandler(tdEvent)
	case tdlib.EventTypeResponse:
		return responseHandler(tdEvent)
	case tdlib.EventTypeError:
		return errorHandler(tdEvent.Data.(tdlib.Error))
	}

	return nil
}

// оброботчик запросов к серверу Телеграм
func requestHandler(tdEvent *tdlib.SystemEvent) *tdlib.Error {
	//bot.Logger.Errorf("New Request %#v\n\n", tdEvent)
	return nil
}

//Обработчик ответов на запросы Telegram client (не используется)
func responseHandler(tdEvent *tdlib.SystemEvent) *tdlib.Error {
	var err *tdlib.Error
	//bot.Logger.Infof("Response %#v\n\n", tdEvent)
	/*
		switch response["@type"].(string) {
		case "authorizationStateWaitPhoneNumber",
			"authorizationStateWaitCode":
			bot.Logger.Errorln("Profile logout")
			bot.Stop()
			bot.ProfileToLogout()
			runProcess = false

		}
	*/

	switch tdEvent.Name {
	case "getAuthorizationState":
		if state, ok := tdEvent.Data.(map[string]interface{})["@type"]; ok {
			switch state.(string) {
			case string(tdlib.AuthorizationStateWaitCodeType),
				string(tdlib.AuthorizationStateWaitPhoneNumberType),
				string(tdlib.AuthorizationStateLoggingOutType):
				bot.Logger.Errorln("Profile logout")
				//bot.Stop()
				bot.ProfileToLogout()
				//bot = nil
				//stopProcess <- true
			}
		}
	case tdbot.EventNameBotReady:
		// TODO: переделать задержку (так что бы бот оставался в сети)
		time.Sleep(time.Second * time.Duration(interval))
		if useMimicry {
			// Run task
			err = completeTask(bot)
		}

		if err == nil { // Ok
			countGood++
			bot.Stop()
		}
	}
	return err
}

//Обработчик ошибок Telegram client. Пересмотреть логику функции
func errorHandler(e tdlib.Error) *tdlib.Error {
	switch e.Code {
	case tdlib.ErrorCodeFloodLock:
		return bot.ProfileToSpam()
	case tdlib.ErrorCodeTimeout:
		fmt.Println("TIMEOUT")
		bot.Restart()
	default:
		fmt.Printf("CHECKER ERROR HANDLER :  %#v\n\n", e)
	}
	return nil
}
