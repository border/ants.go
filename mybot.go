package main

import (
	"os"
    "log"
    "sort"
)

type antDistT struct {
    dist float64
    antLoc Location
    foodLoc Location
}

type antDistListT []antDistT


type MyBot struct {
    debug bool
    orders map[Location]Location
    targets map[Location]Location
    state *State
}

func (a antDistListT) Swap(i, j int) {
    a[i], a[j] = a[j], a[i]
}

func (a antDistListT) Len() int {
    return len(a)
}

func (a antDistListT) Less(i, j int) bool {
    return a[i].dist < a[j].dist
}

//NewBot creates a new instance of your bot
func NewBot(s *State) Bot {
	mb := &MyBot{
        debug: true,
        orders: make(map[Location]Location),
        targets: make(map[Location]Location),
	}
	return mb
}
//Reset clears the Ant for the next turn
func (mb *MyBot) Reset() {
    mb.orders = make(map[Location]Location)
    mb.targets = make(map[Location]Location)
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

// Move Location From source Location to Dest location
func (mb *MyBot) doMoveLocation(loc, dest Location) bool {
    directions := mb.state.Map.FromLocToNewLoc(loc, dest)
    log.Printf("FromLoc: %v, DestLoc: %v directions: %v", loc, dest,  directions)
    for _, direction := range directions {
        if mb.debug {
            row, col := mb.state.Map.FromLocation(loc)
            log.Printf("doMoveLocation Loc(%d, %d) direction: %v", row, col,  direction)
        }
        if mb.doMoveDirection(loc, direction) {
            mb.targets[dest] = loc
            log.Printf("------------------\n eating Fooding mb.targets: %v\n------------------", mb.targets)
            return true
        }
    }
    return false
}


//DoTurn is where you should do your bot's actual work.
func (mb *MyBot) DoTurn(s *State) os.Error {
	//dirs := []Direction{North, East, South, West}
    mb.state = s
    var antDistList antDistListT
    log.Printf("---------------------- DoTurn --------------------")
    log.Printf("-------------------orders: %v-----------------", mb.orders)
    log.Printf("-------------------targets: %v-----------------", mb.targets)

    for foodLoc, is := range s.Map.Food {
        for antLoc, ant := range s.Map.Ants {
            if ant != MY_ANT {
                continue
            }
            if mb.debug {
                log.Printf("My Ant: %d, antLoc: %v, foodLoc: %v, is: %v\n", ant, antLoc, foodLoc, is)
            }
            dist := mb.state.Map.Distance(antLoc, foodLoc)
            var antDist antDistT
            antDist.dist, antDist.antLoc, antDist.foodLoc = dist, antLoc, foodLoc
            antDistList = append(antDistList, antDist)

        }
    }
    sort.Sort(antDistList)
    if mb.debug {
        log.Printf("antDistList after sort %v\n", antDistList)
    }
    for _, antDist := range antDistList {
        if _, is := mb.targets[antDist.foodLoc]; !is {
            var isHas bool
            for _, antLoc := range mb.targets {
                if (antLoc == antDist.antLoc) {
                    isHas = true
                    break
                }
            }
            if (isHas) {
                break
            }
            if mb.debug {
                row, col := mb.state.Map.FromLocation(antDist.antLoc)
                row2, col2 := mb.state.Map.FromLocation(antDist.foodLoc)
                log.Printf("doMoveLocation antLoc(%d, %d), foodLoc(%d, %d)", row, col, row2, col2)
            }
            mb.doMoveLocation(antDist.antLoc, antDist.foodLoc)
        }
    }
	//returning an error will halt the whole program!
	return nil
}
