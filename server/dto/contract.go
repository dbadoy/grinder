package dto

type Contract struct {
	TxHash     string
	Candidates []string
}

func (Contract) Index() string {
	return "contracts"
}
