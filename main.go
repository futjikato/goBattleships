package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/futjikato/goBattleships/battleships"
	"github.com/futjikato/goBattleships/messages"

	"github.com/antham/strumt"
)

const (
	hostType    string = "host"
	clientType  string = "client"
	botType     string = "bot"
	botHostType string = "hostbot"
)

type config struct {
	hostOrClient  string
	hostPort      int
	clientAddress string
	playerName    string
}

func main() {
	p := strumt.NewPromptsFromReaderAndWriter(bufio.NewReader(os.Stdin), os.Stdout)
	c := &config{}

	p.AddLinePrompter(NewChoice("hostorjoin", "[H]ost or [J]oin or [L]et bot host or [B]ot", map[string]string{
		"H": "hostport",
		"J": "joinaddress",
		"L": "hostport",
		"B": "joinaddress",
	}, map[string]string{"H": hostType, "J": clientType, "B": botType, "L": botHostType}, &c.hostOrClient))
	p.AddLinePrompter(NewRangePromter("hostport", "Enter port to listen to", "name", 9000, 9999, &c.hostPort))
	p.AddLinePrompter(NewTextPromter("joinaddress", "Enter host address (<ip>:<port>)", "name", &c.clientAddress))
	p.AddLinePrompter(NewTextPromter("name", "Enter your name", "", &c.playerName))
	p.SetFirst("hostorjoin")
	p.Run()

	if c.hostOrClient == hostType {
		Host(c.playerName, c.hostPort)
	} else if c.hostOrClient == clientType {
		Join(c.playerName, c.clientAddress)
	} else if c.hostOrClient == botType {
		Bot(c.playerName, c.clientAddress)
	} else {
		HostBot(c.playerName, c.hostPort)
	}
}

func handleHit(msg *messages.Msg, b *battleships.Board, outMsgChan chan *messages.Msg) {
	parts := strings.Split(msg.Payload, " ")
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	res := b.AddEnemyHit(x, y)
	hitStr := "0"
	if res.Hit {
		hitStr = "1"
	}
	fmt.Println(b.DrawPlayerSide())
	outMsgChan <- &messages.Msg{
		MsgType: "hitresult",
		Payload: fmt.Sprintf("%d %d %s", x, y, hitStr),
	}
	if b.Lost() {
		outMsgChan <- &messages.Msg{
			MsgType: "youwin",
			Payload: "",
		}
	}
}

func handleHitResult(msg *messages.Msg, b *battleships.Board, outMsgChan chan *messages.Msg) {
	parts := strings.Split(msg.Payload, " ")
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])
	hit := parts[2] == "1"
	b.AddHit(x, y, hit)
	fmt.Println(b.DrawEnemySide())
}
