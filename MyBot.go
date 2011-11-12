package main

import (
	"os"
	"rand"
    "log"
)

type MyBot struct {
    debug bool
}

//NewBot creates a new instance of your bot
func NewBot(s *State) Bot {
	mb := &MyBot{
        debug: true,
	}
	return mb
}

//DoTurn is where you should do your bot's actual work.
func (mb *MyBot) DoTurn(s *State) os.Error {
    orders := make(map[Location]Location)
	dirs := []Direction{North, East, South, West}

	for loc, ant := range s.Map.Ants {
		if ant != MY_ANT {
			continue
		}
        if mb.debug {
            log.Printf("My Ant: %d, %v\n", ant, loc)
        }

		//try each direction in a random order
		p := rand.Perm(4)
		for _, i := range p {
			d := dirs[i]
            loc2 := s.Map.Move(loc, d)
            if s.Map.SafeDestination(loc2) {
                if _, is := orders[loc2]; !is {
                    s.IssueOrderLoc(loc, d)
                    orders[loc2] = loc
                }
            }
		}
	}
	//returning an error will halt the whole program!
	return nil
}
