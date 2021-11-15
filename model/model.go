package model

type TodoModel struct {
	Title string `bson:"title"`
}

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Error   interface{} `json:"error,omitempty"`
}

type UpdateResult struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

type DeleteResult struct {
	ID string `json:"id"`
}
