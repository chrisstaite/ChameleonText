package chameleon

import (
    "time"
	"google.golang.org/appengine/datastore"
)

type TeamType int

const (
    Scout TeamType = iota
    Explorer
    Catcher
)

// Has no parent
type Game struct {
    name string
    number string
    start time.Time
    finish time.Time
}

// Has a Game parent
type Team struct {
    id string
    name string
    type TeamType
}

// Has a Team parent
type TeamMember struct {
    number string
}

// Has a Team parent
type Location struct {
    northing int
    easting int
    time time.Time
}

// Has a Team parent
type Catch struct {
    code string
    caught *datastore.Key
    time time.Time
}

// Has a Team parent
type CaptureCode struct {
    code string
    valid bool
}

// Has a Team parent
type Checkpoint struct {
    code string
    time time.Time
}

// Has a Team parent
type Bonus struct {
    value int
    time time.Time
}
