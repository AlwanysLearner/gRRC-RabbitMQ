package Json

type MyMessage struct {
	ProductId int64 `json:"product_id,omitempty""`
	Number    int64 `json:"number,omitempty"`
}
