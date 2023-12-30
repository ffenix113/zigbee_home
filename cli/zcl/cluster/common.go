package cluster

const (
	Server          Side = 1 << iota
	Client          Side = 1 << iota
	ClientAndServer      = Server | Client
)

type Side uint8

type Provider interface {
	Clusters() Clusters
}

type Cluster interface {
	ID() ID
	CAttrType() string
	CVarName() string
	// ReportAttrCount returns how many attributes are reportable
	ReportAttrCount() int
	// Side tells if the cluster is client and/or server. ZCL 1.3.
	Side() Side
}

type Clusters []Cluster

func (c Clusters) ReportAttrCount() (count int) {
	for _, cluster := range c {
		count += cluster.ReportAttrCount()
	}

	return count
}

func (c Clusters) Servers() (count int) {
	for _, cluster := range c {
		if cluster.Side()&Server == Server {
			count++
		}
	}

	return count
}

func (c Clusters) Clients() (count int) {
	for _, cluster := range c {
		if cluster.Side()&Client == Client {
			count++
		}
	}

	return count
}
