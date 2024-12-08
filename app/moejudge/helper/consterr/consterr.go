package consterr

type ConstErr string

func (c ConstErr) Error() string {
	return *(*string)(&c)
}
