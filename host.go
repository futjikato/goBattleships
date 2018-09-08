package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/futjikato/goBattleships/battleships"
	"github.com/futjikato/goBattleships/messages"
)

type screenState string

const (
	lobbyScreen screenState = "lobby"
	setupScreen screenState = "setup"
	gameScreen  screenState = "game"
)

type hostGameState struct {
	active       bool
	secondPlayer string
	screen       screenState
	board        *battleships.Board
	ready        bool
	clientReady  bool
	turn         bool
}

func Host(hostName string, port int) {
	outMsgChan := make(chan *messages.Msg)
	inMsgChan := make(chan *messages.Msg)
	go startServer(port, outMsgChan, inMsgChan)

	inputChan := make(chan *messages.Msg)
	go input(inputChan)

	gs := &hostGameState{
		active: true,
		screen: lobbyScreen,
		board:  battleships.NewBoard(),
		turn:   true,
	}

	for gs.active {
		select {
		case val := <-inMsgChan:
			switch val.MsgType {
			case "join":
				if gs.screen == lobbyScreen {
					gs.secondPlayer = val.Payload
				} else {
					fmt.Printf("join from %s in wrong game state\n", val.Payload)
				}
				outMsgChan <- &messages.Msg{
					MsgType: "host",
					Payload: hostName,
				}
			case "ready":
				gs.clientReady = true
				fmt.Println(gs.secondPlayer + " is ready")
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
			case "start":
				if gs.screen == lobbyScreen {
					gs.screen = setupScreen
					outMsgChan <- &messages.Msg{
						MsgType: "started",
						Payload: "",
					}
					fmt.Print("Starting game. Setup ships.\nType \"place <Label> <X> <Y> <Orientation>\" \"place CA 2 3 o\"\nShip Labels are:\n")
					fmt.Println("CA - Carrier ( Length: 5 )")
					fmt.Println("BA - Battleship ( Length: 4 )")
					fmt.Println("CR - Cruiser ( Length: 3 )")
					fmt.Println("SU - Submarine ( Length: 3 )")
					fmt.Println("DE - Destroyer ( Length: 2 )")
					fmt.Println("You have one of each.")
				} else {
					fmt.Println("Cant start from non lobby state")
				}
			case "place":
				if gs.screen == setupScreen {
					handleShipPlacement(gs.board, in)
					if gs.board.Ready() {
						gs.ready = true
						outMsgChan <- &messages.Msg{
							MsgType: "ready",
							Payload: "",
						}
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

		if gs.clientReady && gs.ready && gs.screen == setupScreen {
			fmt.Println("All players ready. Starting game")
			gs.screen = gameScreen
		}

		fmt.Println("----------------")
		switch gs.screen {
		case lobbyScreen:
			fmt.Printf("Lobby:\n#1 %s\n#2 %s\n", hostName, gs.secondPlayer)
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

func startServer(port int, outMsgChan, inMsgChan chan *messages.Msg) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("Client connection error: %s\n", err)
		}
		go handleConnection(conn, outMsgChan, inMsgChan)
	}
}

func handleConnection(c net.Conn, outMsgChan <-chan *messages.Msg, inMsgChan chan<- *messages.Msg) {
	go func(out chan<- *messages.Msg) {
		remain := make([]byte, 0)
		for {
			messages, newRemain := messages.ReadFrom(c, remain)
			remain = newRemain
			for _, msg := range messages {
				out <- msg
			}
		}
	}(inMsgChan)

	for {
		msg := <-outMsgChan
		msg.WriteTo(c)
	}
}

func input(inputChannel chan *messages.Msg) {
	for {
		buf := make([]byte, 512)
		i, err := os.Stdin.Read(buf)
		if i > 0 {
			input := strings.Trim(string(buf[:i]), "\n\r\t")
			parts := strings.SplitN(input, " ", 2)
			if len(parts) == 1 {
				inputChannel <- &messages.Msg{
					MsgType: parts[0],
				}
			} else if len(parts) == 2 {
				inputChannel <- &messages.Msg{
					MsgType: parts[0],
					Payload: parts[1],
				}
			}
		}
		if err != nil {
			fmt.Println(err)
		}
	}
}
