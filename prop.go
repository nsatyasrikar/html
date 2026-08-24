package html

func Prop[T any](value T) *T { return &value }

type Script string
