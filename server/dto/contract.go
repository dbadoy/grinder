package dto

type Contract struct {
	TxHash     string
	Candidates []string

	// We need to determine which logic contract
	// the Proxy contract is connected to.
	RelateAddress []string
}

func (Contract) Index() string {
	return "contracts"
}
