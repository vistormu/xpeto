package log

type Sink interface {
	Write(frame uint64, records []Record)
	Flush() error
	Sync() error
}
