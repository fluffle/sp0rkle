package db

// implements some basic storable types for IRC

type StorableNick struct {
	Nick, Ident, Host string
}

type StorableChan struct {
	Chan string
}
