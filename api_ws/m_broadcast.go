package api_ws

import ()

//------------------------------------------------------------
// Websocket Broadcast
//------------------------------------------------------------

type Broadcast struct {
    Op string
    Data interface{}
}

func (b *Broadcast) IsCloseMessage() bool {
    if b.Op == "" && b.Data == nil {
        return true
    }
    return false
}

// Message that causes broadcaster shut down
var CloseBroadcastMessage = &Broadcast{"", nil}

