package core

// Definição da struct Transaction
type Transaction struct {
	ID        string  `json:"id"`
	Amount    float64 `json:"amount"`
	Sender    string  `json:"sender"`
	Receiver  string  `json:"receiver"`
	Timestamp string  `json:"timestamp"`
}

type Block struct {
	Index        int           `json:"index"`
	Timestamp    string        `json:"timestamp"`
	Transactions []Transaction `json:"transactions"`
	Proof        int           `json:"proof"`
	PreviousHash string        `json:"previous_hash"`
}
