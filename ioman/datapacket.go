package ioman

//DataPacket ..
type DataPacket struct {
	Valid   bool
	FlowErr error
	Flow    float32
}
