package main

func Abs[T int | float64](num T) T {
	return 0 - num
}
