package cluster

type Provider interface {
	Clusters() Clusters
}

type Cluster interface {
	ID() ID
	CAttrType() string
	CVarName() string
	// Reports tells if this cluster sends data(i.e. state/measurements) to coordinator
	Reports() bool
}

type Clusters []Cluster

func (c Clusters) ReportAttrCount() (count int) {
	for _, cluster := range c {
		if cluster.Reports() {
			count++
		}
	}

	return count
}
