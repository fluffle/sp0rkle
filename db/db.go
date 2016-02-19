package db

type Database interface {
	Init(db string) error
	Close()
	C(name string) Collection
}

type Collection interface {
}
