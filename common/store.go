package common

// Storage types (indicating where files are stored)
type StoreType int

const (
	_ StoreType = iota
	// StoreLocal : Local node
	StoreLocal
	// StoreCeph : Ceph cluster
	StoreCeph
	// StoreOSS : Aliyun OSS
	StoreOSS
	// StoreMix : Mixed (Ceph and OSS)
	StoreMix
	// StoreAll : Store a copy of data in all types of storage
	StoreAll
)
