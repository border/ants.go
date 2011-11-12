package main

import (
	"os"
    "log"
)

type MyBot struct {
    debug bool
    orders map[Location]Location
    state *State
}

//NewBot creates a new instance of your bot
func NewBot(s *State) Bot {
	mb := &MyBot{
        debug: true,
        orders: make(map[Location]Location),
	}
	return mb
}
//Reset clears the Ant for the next turn
func (mb *MyBot) Reset() {
    mb.orders = make(map[Location]Location)
}

// track all moves, prevent collisions
func (mb *MyBot) doMoveDirection(loc Location, direction Direction) bool {
    loc2 := mb.state.Map.Move(loc, direction)
    if mb.state.Map.SafeDestination(loc2) {
        if _, is := mb.orders[loc2]; !is {
            mb.state.IssueOrderLoc(loc, direction)
            mb.orders[loc2] = loc
            return true
        }
    }
    return false
}

//DoTurn is where you should do your bot's actual work.
func (mb *MyBot) DoTurn(s *State) os.Error {
	dirs := []Direction{North, East, South, West}
    mb.state = s

	for loc, ant := range s.Map.Ants {
		if ant != MY_ANT {
			continue
		}
        if mb.debug {
            log.Printf("My Ant: %d, %v\n", ant, loc)
        }

		//try each direction in a random order
		for _, d := range dirs {
            if mb.doMoveDirection(loc, d) {
                break
            }
		}
	}
	//returning an error will halt the whole program!
	return nil
}
