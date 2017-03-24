package libsql

func newTransaction(tx sqlTx) Transaction {
	return &transactionImpl{
		Queryer:  newQueryerMixin(tx),
		Preparer: newPreparerMixin(tx),
	}
}

type transactionImpl struct {
	Queryer
	Preparer
}

var _ Transaction = (*transactionImpl)(nil)
