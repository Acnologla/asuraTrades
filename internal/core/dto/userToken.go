package dto

type GenerateUserTokenDTO struct {
	AuthorID string
	OtherID  string
	TradeID  string
}

func NewGenerateUserTokenDTO(authorID, otherID, tradeID string) *GenerateUserTokenDTO {
	return &GenerateUserTokenDTO{
		AuthorID: authorID,
		OtherID:  otherID,
		TradeID:  tradeID,
	}
}
