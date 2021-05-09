package providers

type GenericPlane struct {
	FQDN string
	CUID string
	Cores int
	RAM int
	Disk int
	//NetworkInterfaces []GenericNetworkInterface
}

//type GenericNetworkInterface struct {
//	Name string
//	IP string
//	Subnet string
//	MTU int
//	Gateway string
//}
