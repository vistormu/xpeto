package render

type OrderKey uint64

func NewOrderKey(layer uint16, order uint16, tie uint32) OrderKey {
	return OrderKey((uint64(layer) << 48) | (uint64(order) << 32) | uint64(tie))
}
