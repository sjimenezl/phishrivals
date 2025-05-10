package certstream

import (
	"context"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sjimenezl/phishrivals/internal/models"
)

const (
	pingPeriod time.Duration = 15 * time.Second
)

func RunLocalCertstream(ctx context.Context, skipHeartbeats bool) (chan *models.DomainInfo, chan error) {
	outputStream := make(chan *models.DomainInfo)
	errStream := make(chan error)

	go func() {
		defer close(outputStream)
		defer close(errStream)

		for {
			conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/", nil)

			if err != nil {
				errStream <- errors.Wrap(err, "Error connecting to certstream! Sleeping a few seconds and reconnecting... ")
				time.Sleep(5 * time.Second)
				continue
			}

			done := make(chan struct{})
			go keepAlive(conn, done)

			for {
				var msg models.CertstreamMessage
				conn.SetReadDeadline(time.Now().Add(15 * time.Second))
				if err := conn.ReadJSON(&msg); err != nil {
					errStream <- errors.Wrap(err, "error parsing certstream JSON")
					conn.Close()
					close(done)
					break
				}

				if msg.MessageType == "heartbeat" && skipHeartbeats {
					continue
				}

				// convert the timestamps
				nb := time.Unix(int64(msg.Data.LeafCert.NotBefore), 0)
				// na := time.Unix(int64(msg.Data.LeafCert.NotAfter), 0)

				for _, domain := range msg.Data.LeafCert.AllDomains {
					outputStream <- &models.DomainInfo{
						Domain:  domain,
						Created: &nb,
						Source:  "certstream",
					}
				}
			}
		}
	}()

	return outputStream, errStream
}

func keepAlive(conn *websocket.Conn, done <-chan struct{}) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			conn.WriteMessage(websocket.PingMessage, nil)
		case <-done:
			return
		}
	}
}
