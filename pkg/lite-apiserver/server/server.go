package server

type LiteServer struct {
	stopCh <-chan struct{}
}
