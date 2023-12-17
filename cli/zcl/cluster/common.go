package cluster

type Cluster interface {
	ID() ID
	CAttrType() string
	CVarName() string

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
