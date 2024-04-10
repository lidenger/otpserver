package store

type StoreConf interface {
	CloseStore()
	TestStore() error
}
