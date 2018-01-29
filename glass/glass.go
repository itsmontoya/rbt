package glass

// New will return a new Glass
func New() (gp *Glass, err error) {
	var g Glass
	gp = &g
	return
}

// Glass is a database which utilizes whiskey as it's sorting algorithm
type Glass struct {
}
