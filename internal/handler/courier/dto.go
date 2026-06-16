package courier

type Courier struct {
	ID        int    	`json:"id"`
	Name      string 	`json:"name"`
	Phone     string 	`json:"phone"`
	Status    string 	`json:"status"`
}

type CreateRequest struct {
	Name      string	`json:"name"`
	Phone     string 	`json:"phone"`
	Status    string	`json:"status"`
}

type UpdateRequest Courier