package typedef

type Status int

const Divisor = 4 // divide the total width of the viewport by 4

const ( //indices to determine which list is focused
	Todo Status = iota
	InProgress
	Done
)

const (
	ModelEnum Status = iota
	FormEnum
)
