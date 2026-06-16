package courier

import (
	"time"
)

type Courier struct {
	ID        int       
	Name      string    
	Phone     string    
	Status    CourierStatus    
	CreatedAt time.Time 
	UpdatedAt time.Time 
}

type CourierStatus string

const (
	StatusFree 	 = "free"
	StatusBusy   = "busy"
	StatusPaused = "paused"
)