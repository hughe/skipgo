package skipgo

type Option func(ho headOptions)

func Height(height int) Option {
	return func(ho headOptions) {
		ho.setHeight(height)
	}
}

func Base(base int) Option {
	return func(ho headOptions) {
		ho.setBase(base)
	}
}
