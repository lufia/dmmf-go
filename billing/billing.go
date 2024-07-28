package billing

type Price interface {
	Value() float64
}

type Amount float64

func Sum[S ~[]P, P Price](prices S) Amount {
	var sum float64
	for _, v := range prices {
		sum += v.Value()
	}
	return Amount(sum)
}
