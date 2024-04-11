package storeconf

type Status interface {
	GetStoreType() string
	CloseStore()
	TestStore() error
}
