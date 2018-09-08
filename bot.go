package main

import (
	"fmt"

	"github.com/futjikato/goBattleships/battleships"
	"github.com/futjikato/goBattleships/messages"
)

type botGameState struct {
	active    bool
	joined    bool
	otherName string
	screen    screenState
	board     *battleships.Board
	turn      bool
	hitIter   int
}

func Bot(name, address string) {
	outMsgChan := make(chan *messages.Msg)
	inMsgChan := make(chan *messages.Msg)
	go connectServer(address, outMsgChan, inMsgChan)
	outMsgChan <- &messages.Msg{
		MsgType: "join",
		Payload: name,
	}

	gs := &botGameState{
		active:  true,
		screen:  lobbyScreen,
		board:   battleships.NewBoard(),
		hitIter: 0,
	}

	for gs.active {
		select {
		case val := <-inMsgChan:
			switch val.MsgType {
			case "host":
				gs.otherName = val.Payload
				gs.joined = true
				fmt.Println("Joined " + gs.otherName)
			case "started":
				gs.screen = setupScreen
				cx := 1
				for shipType, amount := range battleships.ShipConfigDefault {
					for i := 0; i < amount; i++ {
						switch shipType {
						case battleships.Carrier:
							gs.board.AddShip(battleships.NewCarrier(cx, 1, battleships.Vertical))
							cx++
						case battleships.Battleship:
							gs.board.AddShip(battleships.NewBattleship(cx, 1, battleships.Vertical))
							cx++
						case battleships.Cruiser:
							gs.board.AddShip(battleships.NewCruiser(cx, 1, battleships.Vertical))
							cx++
						case battleships.Submarine:
							gs.board.AddShip(battleships.NewSubmarine(cx, 1, battleships.Vertical))
							cx++
						case battleships.Destroyer:
							gs.board.AddShip(battleships.NewDestroyer(cx, 1, battleships.Vertical))
							cx++
						}
					}
				}
				outMsgChan <- &messages.Msg{
					MsgType: "ready",
					Payload: "",
				}
				fmt.Println(gs.board.DrawPlayerSide())
			case "ready":
				fmt.Println(gs.otherName + " is ready")
			case "hitresult":
				handleHitResult(val, gs.board, outMsgChan)
			case "hit":
				handleHit(val, gs.board, outMsgChan)
				if gs.board.Lost() {
					fmt.Println("LOST :(")
					gs.active = false
				} else {
					hitx := gs.hitIter % 10
					hity := gs.hitIter / 10
					outMsgChan <- &messages.Msg{
						MsgType: "hit",
						Payload: fmt.Sprintf("%d %d", hitx, hity),
					}
					gs.hitIter++
				}
			case "youwin":
				fmt.Println("WON!!!")
				gs.active = false
			default:
				fmt.Printf("Debug unknown message %s: %s\n", val.MsgType, val.Payload)
			}
		}
	}
}
