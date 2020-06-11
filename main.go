package main

import (
	"./lib"
	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

var (
	TOKEN = os.Getenv("TOKEN")
	db = lib.SetupDB()
)


//データベースの状態を確認
func gin_start() {
	router := gin.Default()
	router.Static("./js", "./js")
	router.LoadHTMLGlob("templates/*.html")
	router.GET("/", func(ctx *gin.Context) {
		data := get_format_data()
		ctx.HTML(200, "index.html", gin.H{"data": data})
	})

	router.Run()
}

func get_format_data() string {
	list := ""
	for key, value := range db.Msgs {
		list += key + ":" + value + "\n"
	}
	return list
}

func main() {
	go gin_start()
	dg, err := discordgo.New("Bot " + TOKEN)
	if err != nil {
		log.Fatal("error:start\n", err)
	}

	dg.AddHandler(messageCreate)
	dg.AddHandler(messageDelete)
	err = dg.Open()
	if err != nil {
		log.Fatal("error:wss\n", err)
	}
	log.Println("Ready!")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}

func stringInMap(s string, e map[string]string) bool {
	for k := range e {
		if k == s {
			return true
		}
	}
	return false
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	} else if m.Content == "登録一覧" {
		list := get_format_data()
		s.ChannelMessageSend(m.ChannelID, list)

	} else if strings.Contains(m.Content, "登録:") {
		if strings.Count(m.Content, ":") == 1 {
			s.ChannelMessageSend(m.ChannelID, "ちゃんと送ってね")
			return
		}
		_name := strings.SplitN(m.Content, ":", 3)
		db.Add_msg(_name[1], _name[2])
		s.ChannelMessageSend(m.ChannelID, "["+_name[1]+"]"+"と言った時に["+_name[2]+"]と返って来るようにしました")
	} else if strings.Contains(m.Content, "削除:") {
		_name := strings.SplitN(m.Content, ":", 2)
		if stringInMap(_name[1], db.Msgs) {
			db.Delete_msg(_name[1])
			s.ChannelMessageSend(m.ChannelID, "["+_name[1]+"]"+"と言った時に何も返ってこないようにしました")
		} else {
			s.ChannelMessageSend(m.ChannelID, "["+_name[1]+"]という言葉は登録さ>れていません")
		}
	} else if strings.Contains(m.Content, "oji:") {
		_name := strings.SplitN(m.Content, ":", 2)
		oji, _ := lib.Ojichat(_name[1])
		s.ChannelMessageSend(m.ChannelID, oji)
	} else if m.Content == "realface" {
		lib.Realface("disgord.jpeg")
		i, _ := os.Open("disgord.jpeg")
		s.ChannelFileSend(m.ChannelID, "disgord.jpeg", i)

	} else if stringInMap(m.Content, db.Msgs) {
		s.ChannelMessageSend(m.ChannelID, db.Msgs[m.Content])
	}
}

func messageDelete(s *discordgo.Session, m *discordgo.MessageDelete) {
	s.ChannelMessageSend(m.ChannelID, "誰かがメッセージを取り消しました！")
}
