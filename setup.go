package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/futjikato/goBattleships/messages"

	"github.com/futjikato/goBattleships/battleships"
)

func handleShipPlacement(b *battleships.Board, placeMsg *messages.Msg) {
	placeArgs := strings.Split(placeMsg.Payload, " ")
	if len(placeArgs) != 4 {
		fmt.Println("Invalid number of arguments")
		return
	}
	xc, err := strconv.Atoi(placeArgs[1])
	if err != nil {
		fmt.Printf("Cannot convert X-coordinate: %s\n", err)
		return
	}
	yc, err := strconv.Atoi(placeArgs[2])
	if err != nil {
		fmt.Printf("Cannot convert Y-coordinate: %s\n", err)
		return
	}
	orientation := battleships.ShipOrientation(placeArgs[3])
	switch placeArgs[0] {
	case battleships.Carrier:
		if err = b.AddShip(battleships.NewCarrier(xc, yc, orientation)); err != nil {
			fmt.Printf("Cannot add carrier: %s\n", err)
		}
	case battleships.Battleship:
		if err = b.AddShip(battleships.NewBattleship(xc, yc, orientation)); err != nil {
			fmt.Printf("Cannot add battleship: %s\n", err)
		}
	case battleships.Cruiser:
		if err = b.AddShip(battleships.NewCruiser(xc, yc, orientation)); err != nil {
			fmt.Printf("Cannot add cruiser: %s\n", err)
		}
	case battleships.Submarine:
		if err = b.AddShip(battleships.NewSubmarine(xc, yc, orientation)); err != nil {
			fmt.Printf("Cannot add submarine: %s\n", err)
		}
	case battleships.Destroyer:
		if err = b.AddShip(battleships.NewDestroyer(xc, yc, orientation)); err != nil {
			fmt.Printf("Cannot add destroyer: %s\n", err)
		}
	default:
		fmt.Printf("Unknown ship label %s\n", placeArgs[0])
	}
}
