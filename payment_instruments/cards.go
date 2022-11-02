package payment_instruments

import "context"

type Card struct {
}

func (c Card) Tokenize(ctx context.Context) (string, error) {
	return "", nil
}

func (c Card) Detokenize(ctx context.Context, token string) (string, error) {
	return "", nil
}