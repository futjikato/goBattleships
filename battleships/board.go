package battleships

import (
	"fmt"
	"strings"
)

type Board struct {
	maxShips  map[string]int
	ships     []*Ship
	hits      []*Hit
	enemyHits []*Hit
}

var (
	ShipConfigDefault = map[string]int{Carrier: 1, Battleship: 1, Cruiser: 1, Submarine: 1, Destroyer: 1}
)

func NewBoard() *Board {
	return &Board{
		maxShips:  ShipConfigDefault,
		ships:     make([]*Ship, 0),
		hits:      make([]*Hit, 0),
		enemyHits: make([]*Hit, 0),
	}
}

func (b *Board) DrawPlayerSide() string {
	builder := &strings.Builder{}
	shipMatrix := b.getShipMatrix()

	for _, hit := range b.enemyHits {
		i := hit.YCoordinate*10 + hit.XCoordinate
		if _, hasShip := shipMatrix[i]; hasShip == true {
			shipMatrix[i] = " XX "
		} else {
			shipMatrix[i] = " OO "
		}
	}

	builder.WriteString("Your ships\n")
	builder.WriteString("   01  02  03  04  05  06  07  08  09  10")
	for i := 0; i < 100; i++ {
		x := i % 10
		y := i / 10
		if x == 0 {
			builder.WriteString(fmt.Sprintf("\n%2d", y+1))
		}

		if _, hasEntry := shipMatrix[i]; hasEntry == true {
			builder.WriteString(shipMatrix[i])
		} else {
			builder.WriteString(" __ ")
		}

	}

	return builder.String()
}

func (b *Board) DrawEnemySide() string {
	builder := &strings.Builder{}

	hitMatrix := make(map[int]string)
	for _, hit := range b.hits {
		i := hit.YCoordinate*10 + hit.XCoordinate
		if hit.Hit {
			hitMatrix[i] = " XX "
		} else {
			hitMatrix[i] = " OO "
		}
	}

	builder.WriteString("Your hits\n")
	builder.WriteString("   01  02  03  04  05  06  07  08  09  10")
	for i := 0; i < 100; i++ {
		x := i % 10
		y := i / 10
		if x == 0 {
			builder.WriteString(fmt.Sprintf("\n%2d", y+1))
		}

		if o, isHit := hitMatrix[i]; isHit == true {
			builder.WriteString(o)
		} else {
			builder.WriteString(" __ ")
		}

	}

	return builder.String()
}

func (b *Board) AddShip(ship *Ship) error {
	allowed := b.maxShips[ship.Label]
	for _, s := range b.ships {
		if s.Label == ship.Label {
			allowed--
		}
	}

	if allowed < 0 {
		return fmt.Errorf("Max number of ships of this kind reached")
	}

	b.ships = append(b.ships, ship)
	return nil
}

func (b *Board) AddHit(x int, y int, hit bool) {
	b.hits = append(b.hits, &Hit{XCoordinate: x, YCoordinate: y, Hit: hit})
}

func (b *Board) AddEnemyHit(x int, y int) *Hit {
	shipMatrix := b.getShipMatrix()

	hit := false
	i := y*10 + x
	if _, hasShip := shipMatrix[i]; hasShip == true {
		hit = true
	}
	hitObj := &Hit{XCoordinate: x, YCoordinate: y, Hit: hit}
	b.enemyHits = append(b.enemyHits, hitObj)

	return hitObj
}

func (b *Board) Ready() bool {
	maxCopy := make(map[string]int)
	for k, v := range b.maxShips {
		maxCopy[k] = v
	}
	for _, s := range b.ships {
		maxCopy[s.Label]--
	}
	for _, v := range maxCopy {
		if v > 0 {
			return false
		}
	}

	return true
}

func (b *Board) Lost() bool {
	hitsToLose := 0
	for _, ship := range b.ships {
		hitsToLose += ship.Length * ship.Width
	}

	connectedHits := 0
	for _, hit := range b.enemyHits {
		if hit.Hit {
			connectedHits++
		}
	}

	return connectedHits >= hitsToLose
}

func (b *Board) getShipMatrix() map[int]string {
	shipMatrix := map[int]string{}
	for _, ship := range b.ships {
		startCoord := ship.YCoordinate*10 + ship.XCoordinate
		for i := 0; i < ship.Length; i++ {
			if ship.Orientation == Horizontal {
				shipMatrix[startCoord+i] = fmt.Sprintf(" %s ", ship.Label)
			} else {
				shipMatrix[startCoord+i*10] = fmt.Sprintf(" %s ", ship.Label)
			}
		}
	}

	return shipMatrix
}
