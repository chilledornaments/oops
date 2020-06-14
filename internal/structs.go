package internal

// Secret is our DynamoDB item
type Secret struct {
	Secret     string
	Expiration int64
	OopsID     string
}
