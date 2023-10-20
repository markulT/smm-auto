package payments

type CardData struct {
	CardNumber string `json:"cardNumber"`
	ExpMonth int64 `json:"expMonth"`
	ExpYear int64 `json:"expYear"`
	CVC string `json:"cvc"`
}
