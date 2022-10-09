package main

type status int

const divisor = 4 // divide the total width of the viewport by 4

const ( //indices to determine which list is focused
	todo status = iota
	inProgress
	done
)

const (
	model status = iota
	form
)
