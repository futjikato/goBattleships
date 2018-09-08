package battleships

type ShipOrientation string

const (
	Horizontal ShipOrientation = "h"
	Vertical   ShipOrientation = "v"

	Carrier    string = "CA"
	Battleship string = "BA"
	Cruiser    string = "CR"
	Submarine  string = "SU"
	Destroyer  string = "DE"
)

type Ship struct {
	XCoordinate int
	YCoordinate int
	Width       int
	Length      int
	Orientation ShipOrientation
	Label       string
}

func NewCarrier(x int, y int, orientation ShipOrientation) *Ship {
	return &Ship{
		Label:       Carrier,
		Length:      5,
		Width:       1,
		XCoordinate: x - 1,
		YCoordinate: y - 1,
		Orientation: orientation,
	}
}

func NewBattleship(x int, y int, orientation ShipOrientation) *Ship {
	return &Ship{
		Label:       Battleship,
		Length:      4,
		Width:       1,
		XCoordinate: x - 1,
		YCoordinate: y - 1,
		Orientation: orientation,
	}
}

func NewCruiser(x int, y int, orientation ShipOrientation) *Ship {
	return &Ship{
		Label:       Cruiser,
		Length:      3,
		Width:       1,
		XCoordinate: x - 1,
		YCoordinate: y - 1,
		Orientation: orientation,
	}
}

func NewSubmarine(x int, y int, orientation ShipOrientation) *Ship {
	return &Ship{
		Label:       Submarine,
		Length:      3,
		Width:       1,
		XCoordinate: x - 1,
		YCoordinate: y - 1,
		Orientation: orientation,
	}
}

func NewDestroyer(x int, y int, orientation ShipOrientation) *Ship {
	return &Ship{
		Label:       Destroyer,
		Length:      2,
		Width:       1,
		XCoordinate: x - 1,
		YCoordinate: y - 1,
		Orientation: orientation,
	}
}
