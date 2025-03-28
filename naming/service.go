package naming

// 提供一个对inter.ServiceRegistration的默认实现
type DefaultServiceRegistration struct {
	Id        string
	Name      string
	Address   string
	Port      int
	Protocol  string
	Namespace string
	Tags      []string
	Meta      map[string]string
}

func (d DefaultServiceRegistration) ServiceID() string {
	return d.Id
}

func (d DefaultServiceRegistration) ServiceName() string {
	return d.Name
}

func (d DefaultServiceRegistration) GetMeta() map[string]string {
	return d.Meta
}

func (d DefaultServiceRegistration) GetTags() []string {
	return d.Tags
}

func (d DefaultServiceRegistration) PublicPort() int {
	return d.Port
}

func (d DefaultServiceRegistration) PublicAddress() string {
	return d.Address
}

func (d DefaultServiceRegistration) GetProtocol() string {
	return d.Protocol
}
