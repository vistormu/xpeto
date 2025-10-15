package asset

type EventAssetAdded struct {
	Handle Handle
}

type EventAssetModified struct {
	Handle Handle
}

type EventAssetRemoved struct {
	Handle Handle
}
