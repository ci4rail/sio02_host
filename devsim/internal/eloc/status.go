/*
Copyright © 2022 Ci4Rail GmbH
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package eloc

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/ci4rail/io4edge-client-go/client"
	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge-client-go/transport/socket"
	"github.com/ci4rail/sio01_host/devsim/internal/eloc/pb"
)

func (e *Eloc) statusServer(port int) error {

	listener, err := socket.NewSocketListener(fmt.Sprintf(":%d", port))

	if err != nil {
		log.Printf("Failed to create listener: %s", err)
		return err
	}

	for {
		conn, err := socket.WaitForSocketConnect(listener)
		if err != nil {
			log.Fatalf("Failed to wait for connection: %s", err)
		}
		log.Printf("statusServer: new connection from %s!\n", conn.RemoteAddr())

		go func(conn *net.TCPConn) {
			ms := transport.NewFramedStreamFromTransport(conn)
			ch := client.NewChannel(ms)

			e.serveConnection(ch)
		}(conn)
	}
}

func (e *Eloc) serveConnection(ch *client.Channel) {
	defer ch.Close()

	for {
		req := &pb.StatusRequest{}
		err := ch.ReadMessage(req, 0)
		if err != nil {
			if err == io.EOF {
				return
			}
			log.Printf("serveConnection failed to read: %s", err)
			return
		}

		res := &pb.StatusResponse{
			Id:                  req.Id,
			PowerUpCount:        0,
			HasPosition:         e.havePosition,
			HasServerConnection: e.haveServerConnection,
			HasTime:             true,
			ElocModuleStatusOk:  true,
		}

		err = ch.WriteMessage(res)
		if err != nil {
			log.Printf("serveConnection failed to write: %s", err)
			return
		}
	}
}
