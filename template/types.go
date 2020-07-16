package template

type NodeParam struct {
	Name          string
	Image         string
	ContainerName string
	Command       string
	PortP2P       string
	PortHttp      string
	PortWs        string
	Ipv4Address   string
	Volumes       string
}

type PTMParam struct {
	Name          string
	Image         string
	ContainerName string
	Command       string
	Port1         string
	Port2         string
	Ipv4Address   string
	Volumes       string
}

type QlcNode struct {
}

type PtmNode struct {
	Url string
}
