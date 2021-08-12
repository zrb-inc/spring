package runtime

type BeanFactory interface {
	GetBean(name string) (interface{}, error)
}

type ListableBeanFactory interface {
	BeanFactory
	GetBeans() ([]interface{}, error)
}
