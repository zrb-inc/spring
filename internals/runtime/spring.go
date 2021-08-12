package runtime

type Spring struct {
}

type ApplicationContext interface {
	ListableBeanFactory
}

type CompileApplicationContext struct {
	GeneralApplicationContext
}

func (c *CompileApplicationContext) PushDefinition(name string, d BeanDefinition) {
	if _, ok := c.Definitions[name]; !ok {
		c.Definitions[name] = d
	}
}

var CompileApp = CompileApplicationContext{
	GeneralApplicationContext: GeneralApplicationContext{
		Definitions: map[string]BeanDefinition{},
	},
}

type GeneralApplicationContext struct {
	Definitions map[string]BeanDefinition
}

func (g *GeneralApplicationContext) Config() {
	g.Definitions = CompileApp.Definitions
}

func (g *GeneralApplicationContext) GetBean(name string) (interface{}, error) {
	var b BeanDefinition

	found := false

	for k, v := range g.Definitions {
		if k == name || v.GetId() == name {
			b = v
			found = true
		}
	}

	if !found {
		return nil, BeanNotFound
	}
	fn := b.GetConstruct()

	instance, err := fn()
	if err != nil {
		return nil, FailedConstruct
	}
	return instance, nil
}

func (g *GeneralApplicationContext) GetBeans() ([]interface{}, error) {

	instances := []interface{}{}

	for _, v := range g.Definitions {
		fn := v.GetConstruct()
		instance, err := fn()
		if err != nil {
			return nil, FailedConstruct
		}
		instances = append(instances, instance)
	}
	return instances, nil
}
