// Copyright 2014 Verizon Laboratories. All rights reserved.
package metadata

import "time"

type MessageMetadata struct {
	EnqueuedAt time.Time
	TcpTxDuration time.Duration
	ProdRef interface {}
}

const(
	INVALID_PARTITION int32 = -1
)

func (mm *MessageMetadata) Producer() interface{}{
	return mm.ProdRef
}

func (mm *MessageMetadata) GetMsgDelay() (time.Duration, time.Duration){
	return time.Since(mm.EnqueuedAt), mm.TcpTxDuration
}

func (mm *MessageMetadata) SetBatchNwDelay(tcpTxDuration time.Duration) (){
	mm.TcpTxDuration = tcpTxDuration
}
