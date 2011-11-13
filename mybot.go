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
    // to track what moves we have issued.
    orders map[Location]Location
    // foodLoc:antLoc
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
    newLoc := mb.state.Map.Move(loc, direction)
    if mb.state.Map.SafeDestination(newLoc) {
        if _, is := mb.orders[newLoc]; !is {
            if mb.debug {
                row, col := mb.state.Map.FromLocation(loc)
                row2, col2 := mb.state.Map.FromLocation(newLoc)
                log.Printf("(1)doMoveDirection antLoc(%d, %d), newLoc(%d, %d)", row, col, row2, col2)
            }
            mb.state.IssueOrderLoc(loc, direction)
            mb.orders[newLoc] = loc
            return true
        }
    }
    return false
}

// Move Location From source Location to Dest location
func (mb *MyBot) doMoveLocation(loc, dest Location) bool {
    directions := mb.state.Map.FromLocToNewLoc(loc, dest)
    for _, direction := range directions {
        if mb.doMoveDirection(loc, direction) {
            if mb.debug {
                row, col := mb.state.Map.FromLocation(loc)
                row2, col2 := mb.state.Map.FromLocation(dest)
                log.Printf("(2)doMoveLocation antLoc(%d, %d), foodLoc(%d, %d) direction: %v", row, col, row2, col2,  direction)
            }
            mb.targets[dest] = loc
            return true
        }
    }
    return false
}


//DoTurn is where you should do your bot's actual work.
func (mb *MyBot) DoTurn(s *State) os.Error {
	dirs := []Direction{North, East, South, West}
    mb.state = s
    var antDistList antDistListT
    log.Printf("---------------------- DoTurn --------------------")
    log.Printf("-------------------targets: %v-----------------", mb.targets)

    // prevent stepping on own hil
    /*
    for hillLoc, _ := range s.Map.Hills {
        mb.orders[hillLoc] = s.Map.FromRowCol(0, 0)
    }
    */
    log.Printf("-------------------orders: %v-----------------", mb.orders)

    for foodLoc, _ := range s.Map.Food {
        for antLoc, ant := range s.Map.Ants {
            if ant != MY_ANT {
                continue
            }
            if mb.debug {
                //log.Printf("My Ant: %d, antLoc: %v, foodLoc: %v, is: %v\n", ant, antLoc, foodLoc, is)
            }
            dist := mb.state.Map.Distance(antLoc, foodLoc)
            var antDist antDistT
            antDist.dist, antDist.antLoc, antDist.foodLoc = dist, antLoc, foodLoc
            antDistList = append(antDistList, antDist)

        }
    }

    // unblock own hill
    for hillLoc, _ := range s.Map.Hills {
        for antLoc, _ := range s.Map.Ants {
            if hillLoc == antLoc {
                var isHas bool
                 for _, loc := range mb.orders {
                    if (hillLoc == loc) {
                        isHas = true
                        break
                    }
                 }
                if (isHas) {
                    break
                }
	            for _, d := range dirs {
                    if mb.doMoveDirection(hillLoc, d) {
                        break
                    }
                }
            }
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
            mb.doMoveLocation(antDist.antLoc, antDist.foodLoc)
        }
    }

    // explore unseen areas
    for antLoc, _ := range s.Map.Ants {
         var isHas bool
         for _, loc := range mb.orders {
            if (antLoc == loc) {
                isHas = true
                break
            }
         }
        if (isHas) {
            break
        }

        var unseenList antDistListT
        for unseen, item := range s.Map.itemGrid {
            if (item != UNKNOWN) {
                continue
            }
            unseenLoc := Location(unseen)

            dist := mb.state.Map.Distance(antLoc, unseenLoc)
            var unseenDist antDistT
            unseenDist.dist, unseenDist.foodLoc = dist, unseenLoc
            unseenList = append(unseenList, unseenDist)
        }

        sort.Sort(unseenList)

        for _, unseenDist := range unseenList {
            if mb.doMoveLocation(antLoc, unseenDist.foodLoc) {
                s.Map.AddWater(antLoc)
                break
            }
        }
    }


	//returning an error will halt the whole program!
	return nil
}
