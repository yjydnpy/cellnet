package gate

import (
	"github.com/davyxu/cellnet"
	"github.com/davyxu/cellnet/proto/coredef"
	"github.com/davyxu/cellnet/socket"
)

var BackendAcceptor cellnet.Peer

// 开启后台服务器的侦听通道
func StartBackendAcceptor(pipe cellnet.EventPipe, address string) {

	BackendAcceptor = socket.NewAcceptor(pipe)

	// 默认开启并发
	BackendAcceptor.EnableConcurrenceMode(true)

	// 关闭客户端连接
	socket.RegisterSessionMessage(BackendAcceptor, "coredef.CloseClientACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*coredef.CloseClientACK)

		if msg.ClientID == 0 {

			// 关闭所有客户端
			ClientAcceptor.IterateSession(func(ses cellnet.Session) bool {

				if DebugMode {
					log.Debugf("backend->client, close clientid %d", msg.ClientID)
				}
				ses.Close()

				return true
			})

		} else {

			// 关闭指定客户端
			clientSes := ClientAcceptor.GetSession(msg.ClientID)

			// 找到连接并关闭
			if clientSes != nil {

				if DebugMode {
					log.Debugf("backend->client, close clientid %d", msg.ClientID)
				}

				clientSes.Close()
			} else if DebugMode {
				log.Debugf("backend->client, client not found, close failed, clientid %d", msg.ClientID)
			}

		}

	})

	// 广播
	socket.RegisterSessionMessage(BackendAcceptor, "coredef.DownstreamACK", func(content interface{}, ses cellnet.Session) {
		msg := content.(*coredef.DownstreamACK)

		pkt := &cellnet.Packet{
			MsgID: msg.MsgID,
			Data:  msg.Data,
		}

		if len(msg.ClientID) == 0 {

			// 广播给所有客户端
			ClientAcceptor.IterateSession(func(ses cellnet.Session) bool {

				if DebugMode {
					log.Debugf("backend->client, msgid: %d clientid %d", msg.MsgID, msg.ClientID)
				}

				ses.RawSend(pkt)

				return true
			})

		} else {

			// 指定客户端发送
			for _, clientid := range msg.ClientID {
				clientSes := ClientAcceptor.GetSession(clientid)

				if clientSes != nil {

					if DebugMode {
						log.Debugf("backend->client, msgid: %d clientid %d", msg.MsgID, msg.ClientID)
					}

					clientSes.RawSend(pkt)

				} else if DebugMode {

					log.Debugf("backend->client, client not found, msgid: %d clientid %d", msg.MsgID, msg.ClientID)
				}
			}
		}

	})

	BackendAcceptor.Start(address)
}
