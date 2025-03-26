package core

type Blockchain struct {
	Chain               []Block       `json:"chain"`
	CurrentTransactions []Transaction `json:"current_transactions"`
}

func NewBlockchain() *Blockchain {
	blockchain := &Blockchain{
		Chain:               []Block{},
		CurrentTransactions: []Transaction{},
	}
	// Cria o bloco gênesis
	blockchain.NewBlock(100, "1")
	return blockchain
}

func (bc *Blockchain) NewBlock(proof int, previousHash string) Block {
	block := Block{
		Index:        len(bc.Chain) + 1,
		Timestamp:    "data/hora atual",
		Transactions: bc.CurrentTransactions,
		Proof:        proof,
		PreviousHash: previousHash,
	}
	bc.CurrentTransactions = []Transaction{}
	bc.Chain = append(bc.Chain, block)
	return block
}

// Outras funções de validação e consenso
