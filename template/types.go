package template

import "github.com/qlcchain/qlc-go-sdk/pkg/types"

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
	accountSeed   string
}

type PtmParam struct {
	Name          string
	Image         string
	ContainerName string
	Command       string
	Port1         string
	Port2         string
	Ipv4Address   string
	Volumes       string
	cEndpoint     string
}

type QlcNode struct {
	HTTPEndpoint  string
	WSEndpoint    string
	ContainerName string
	Account       *types.Account
}

type PtmNode struct {
	EndPoint      string
	ContainerName string
}
