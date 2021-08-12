package runtime

import "errors"

var (
	BeanNotFound    = errors.New("BeanNotFound err")
	FailedConstruct = errors.New("FailedConstruct err")
)
