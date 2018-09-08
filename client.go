package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/futjikato/goBattleships/battleships"

	"github.com/futjikato/goBattleships/messages"
)

type clientGameState struct {
	active     bool
	joined     bool
	hostPlayer string
	screen     screenState
	board      *battleships.Board
	hostReady  bool
	ready      bool
	turn       bool
}

func Join(playerName string, address string) {
	outMsgChan := make(chan *messages.Msg)
	inMsgChan := make(chan *messages.Msg)
	go connectServer(address, outMsgChan, inMsgChan)
	outMsgChan <- &messages.Msg{
		MsgType: "join",
		Payload: playerName,
	}

	inputChan := make(chan *messages.Msg)
	go input(inputChan)

	gs := &clientGameState{
		active: true,
		screen: lobbyScreen,
		board:  battleships.NewBoard(),
	}

	for gs.active {
		select {
		case val := <-inMsgChan:
			switch val.MsgType {
			case "host":
				gs.hostPlayer = val.Payload
				gs.joined = true
				fmt.Println("Joined " + gs.hostPlayer)
			case "started":
				gs.screen = setupScreen
				fmt.Print("Host started game. Setup ships.\nType \"place <Label> <X> <Y> <Orientation>\" \"place CA 2 3 o\"\nShip Labels are:\n")
				fmt.Println("CA - Carrier ( Length: 5 )")
				fmt.Println("BA - Battleship ( Length: 4 )")
				fmt.Println("CR - Cruiser ( Length: 3 )")
				fmt.Println("SU - Submarine ( Length: 3 )")
				fmt.Println("DE - Destroyer ( Length: 2 )")
				fmt.Println("You have one of each.")
			case "ready":
				gs.hostReady = true
				fmt.Println(gs.hostPlayer + " is ready")
			case "hitresult":
				handleHitResult(val, gs.board, outMsgChan)
			case "hit":
				handleHit(val, gs.board, outMsgChan)
				if gs.board.Lost() {
					fmt.Println("LOST :(")
					gs.active = false
				} else {
					gs.turn = true
				}
			case "youwin":
				fmt.Println("WON!!!")
				gs.active = false
			default:
				fmt.Printf("Debug unknown message %s: %s\n", val.MsgType, val.Payload)
			}

		case in := <-inputChan:
			switch in.MsgType {
			case "place":
				handleShipPlacement(gs.board, in)
				if gs.board.Ready() {
					gs.ready = true
					outMsgChan <- &messages.Msg{
						MsgType: "ready",
						Payload: "",
					}
				}
			case "hit":
				if gs.screen == gameScreen && gs.turn {
					gs.turn = false
					parts := strings.Split(in.Payload, " ")
					x, _ := strconv.Atoi(parts[0])
					y, _ := strconv.Atoi(parts[1])
					outMsgChan <- &messages.Msg{
						MsgType: "hit",
						Payload: fmt.Sprintf("%d %d", x-1, y-1),
					}
				} else {
					fmt.Println("Not in game mode or not your turn")
				}
			default:
				fmt.Printf("Debug unknown input %s\n", in)
			}
		}

		if gs.hostReady && gs.ready && gs.screen == setupScreen {
			fmt.Println("All players ready. Starting game")
			gs.screen = gameScreen
		}

		fmt.Println("----------------")
		switch gs.screen {
		case lobbyScreen:
			if gs.joined {
				fmt.Printf("Lobby:\n#1 %s\n#2 %s\n", gs.hostPlayer, playerName)
			} else {
				fmt.Printf("Not in lobby")
			}
		case setupScreen:
			fmt.Println(gs.board.DrawPlayerSide())
		case gameScreen:
			if gs.turn {
				fmt.Println("Your turn")
			} else {
				fmt.Println("Enemy turn")
			}
		}
		fmt.Println("----------------")
	}
}

func connectServer(address string, outMsgChan, inMsgChan chan *messages.Msg) {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		panic(err)
	}
	go func(out chan<- *messages.Msg) {
		remain := make([]byte, 0)
		for {
			messages, newRemain := messages.ReadFrom(conn, remain)
			remain = newRemain
			for _, msg := range messages {
				out <- msg
			}
		}
	}(inMsgChan)

	for {
		msg := <-outMsgChan
		msg.WriteTo(conn)
	}
}
