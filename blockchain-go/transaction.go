package main

import "fmt"

type Transaction struct {
	Sender    string  `json:"sender"`
	Recipient string  `json:"recipient"`
	Amount    float64 `json:"amount"`
}

func (t *Transaction) String() string {
	return fmt.Sprintf("{sender: %s, recipient: %s, amount: %.2f}", t.Sender, t.Recipient, t.Amount)
}
