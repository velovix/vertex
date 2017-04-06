package main

var lastUID uidGenerator = 0

type uidGenerator int

func (uidG *uidGenerator) uid() int {
	if *uidG == 0 {
		lastUID++
		(*uidG) = lastUID
	}

	return int(*uidG)
}
