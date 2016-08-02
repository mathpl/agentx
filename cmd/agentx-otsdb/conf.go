package main

type Conf struct {
	Tags map[string]string
	Freq int

	MIBS map[string]MIB
}

type MIB struct {
	BaseOid string
	Metrics []MIBMetric // single key metrics
	//Trees   []MIBTree   // tagged array metrics
}

type MIBMetric struct {
	Metric      string
	Oid         string
	Description string
	Tags        string // static tags to populate for this metric. "direction=in"
}

//type MIBTag struct {
//	Key string
//	Oid string // If present will load from this oid. Use "idx" to populate with index of row instead of another oid.
//}
//
//type MIBTree struct {
//	BaseOid string
//	Tags    []MIBTag
//	Metrics []MIBMetric
//}
