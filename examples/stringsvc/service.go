package somesvc

import "io"

type SomeService interface {
	SomeSliceGetter() []io.Reader
}
