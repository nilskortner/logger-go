package mpscchannel

type ChannelView[T comparable] struct {
	data chan T
	next *ChannelView[T]
}
