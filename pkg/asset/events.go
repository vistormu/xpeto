package asset

type AssetEvent struct {
	Handle Handle
	Kind   AssetEventKind
}

type AssetEventKind uint8

const (
	Added AssetEventKind = iota
	Modified
	Removed
)
