package types

import "github.com/ffenix113/zigbee_home/zcl/cluster"

type Endpoints []*Endpoint

type Endpoint struct {
	Num      int
	Clusters cluster.Clusters
}

func (e *Endpoints) LastEndpoint() *Endpoint {
	if len(*e) == 0 {
		return e.AddEndpoint()
	}

	return (*e)[len(*e)-1]
}

func (e *Endpoints) AddEndpoint() *Endpoint {
	// This will create first endpoint with number 1
	*e = append(*e, &Endpoint{Num: len(*e) + 1})
	return (*e)[len(*e)-1]
}

func (e *Endpoint) HasCluster(clusterID cluster.ID) bool {
	for _, cluster := range e.Clusters {
		if cluster.ID() == clusterID {
			return true
		}
	}

	return false
}

func (e *Endpoint) AddCluster(cluster cluster.Cluster) {
	e.Clusters = append(e.Clusters, cluster)
}
