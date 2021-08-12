package runtime

// BeanDefinition defines a bean's name ,id, dependency, and construct method
type BeanDefinition interface {
	SetBeanPackageName(string)
	GetBeanPackageName() string

	SetId(id string)
	GetId() string

	SetDependsOn(dependsOn ...string)
	GetDependsOn() []string

	IsPrimary() bool

	SetConstruct(func() (interface{}, error))
	GetConstruct() func() (interface{}, error)
}

type GeneralBeanDefinition struct {
	PackageName string
	Id          string
	DependsOn   []string
	Primary     bool
	Constrction func() (interface{}, error)
}

func (g *GeneralBeanDefinition) SetBeanPackageName(name string) {
	g.PackageName = name
}

func (g *GeneralBeanDefinition) GetBeanPackageName() string {
	return g.PackageName
}

func (g *GeneralBeanDefinition) SetId(id string) {
	g.Id = id
}

func (g *GeneralBeanDefinition) GetId() string {
	return g.Id
}

func (g *GeneralBeanDefinition) SetDependsOn(dependsOn ...string) {
	g.DependsOn = dependsOn
}

func (g *GeneralBeanDefinition) GetDependsOn() []string {
	return g.DependsOn
}

func (g *GeneralBeanDefinition) IsPrimary() bool {
	return g.Primary
}

func (g *GeneralBeanDefinition) SetConstruct(fn func() (interface{}, error)) {
	g.Constrction = fn
}

func (g *GeneralBeanDefinition) GetConstruct() func() (interface{}, error) {
	return g.Constrction
}
