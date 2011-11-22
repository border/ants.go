package main

import (
	"os"
    "log"
    "sort"
//    "rand"
)

var maxint float64 = 9999

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

func (mb *MyBot)ClosestFood(loc Location, filter map[Location]Location) (retLoc Location) {
    minDist := maxint
    for foodLoc, _ := range mb.state.Map.Food {
        if (len(filter) == 0) || (!isInLocMap(filter, foodLoc)) {
            dist := mb.state.Map.Distance(loc, foodLoc)
            if dist < minDist {
                minDist = dist
                retLoc = foodLoc
            }
        }
    }
    return
}

func (mb *MyBot)ClosestEnemyAnt(loc Location, filter map[Location]Location) (retLoc Location) {
    minDist := maxint
    for hillAntLoc, ant := range mb.state.Map.Ants {
        if ant == MY_ANT {
            continue
        }
        if (len(filter) == 0) || (!isInLocMap(filter, hillAntLoc)) {
            dist := mb.state.Map.Distance(loc, hillAntLoc)
            if dist < minDist {
                minDist = dist
                retLoc = hillAntLoc
            }
        }
    }
    return
}

func (mb *MyBot)ClosestEnemyHill(loc Location, filter map[Location]Location) (retLoc Location) {
    minDist := maxint
    for hillLoc, ant := range mb.state.Map.Hills {
        if ant == MY_HILL {
            continue
        }
        if (len(filter) == 0) || (!isInLocMap(filter, hillLoc)) {
            dist := mb.state.Map.Distance(loc, hillLoc)
            if dist < minDist {
                minDist = dist
                retLoc = hillLoc
            }
        }
    }
    return
}

func (mb *MyBot)ClosestUnseen(loc Location, filter map[Location]Location) (retLoc Location) {
    minDist := maxint
    var unseenLoc Location
    for row :=0; row < mb.state.Map.Rows; row++ {
        for col := 0; col < mb.state.Map.Cols; col++ {
            unseenLoc = mb.state.Map.FromRowCol(row, col)
            if (len(filter) == 0) || (!isInLocMap(filter, unseenLoc)) {
                dist := mb.state.Map.Distance(loc, unseenLoc)
                if dist < minDist {
                    minDist = dist
                    retLoc = unseenLoc
                }
            }
        }
    }
    return
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

func isInLocMap(srcLoc map[Location]Location, tLoc Location) bool {
     for _, loc := range srcLoc {
        if (tLoc == loc) {
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
    var hillDistList antDistListT
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

    sort.Sort(antDistList)

    if mb.debug {
        log.Printf("antDistList after sort %v\n", antDistList)
    }
    for _, antDist := range antDistList {
        if _, is := mb.targets[antDist.foodLoc]; !is {
            if !isInLocMap(mb.targets, antDist.antLoc) {
                mb.doMoveLocation(antDist.antLoc, antDist.foodLoc)
            }
        }
    }

    // attack Hills
    for hillLoc, hill := range s.Map.Hills {
        if hill != MY_HILL {
            continue
        }
        for antLoc, ant := range s.Map.Ants {
            if ant != MY_ANT {
                continue
            }

            if !isInLocMap(mb.orders, antLoc) {
                continue
            }
            dist := mb.state.Map.Distance(antLoc, hillLoc)
            var hillDist antDistT
            hillDist.dist, hillDist.foodLoc, hillDist.foodLoc = dist, antLoc, hillLoc
            hillDistList = append(hillDistList, hillDist)
        }
    }
    sort.Sort(hillDistList)
    for _, hillDist := range hillDistList {
        mb.doMoveLocation(hillDist.antLoc, hillDist.foodLoc)
    }

    // The map exploration unseen areas
    for antLoc, _ := range s.Map.Ants {
        if isInLocMap(mb.orders, antLoc) {
            continue
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
                break
            }
        }
    }

    // unblock own hill
    for hillLoc, hill := range s.Map.Hills {
        if hill != MY_HILL {
            continue
        }
        for antLoc, _ := range s.Map.Ants {
            if hillLoc == antLoc  {
                if !isInLocMap(mb.orders, hillLoc) {
                    for _, d := range dirs {
                        if mb.doMoveDirection(hillLoc, d) {
                            break
                        }
                    }
                }
            }
        }
    }
	//returning an error will halt the whole program!
	return nil
}
