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
	"errors"
	"log"
	"time"

	"github.com/ci4rail/io4edge-client-go/client"
	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge-client-go/transport/socket"
	"github.com/ci4rail/sio01_host/devsim/internal/eloc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type location struct {
	havePos bool
	x       float64
	y       float64
	z       float64
}

func (e *Eloc) locationClient(locationServerAddress string) error {
	go func() {
		for {
			e.haveServerConnection = false
			ch, err := channelFromSocketAddress(locationServerAddress)

			if err == nil {
				defer ch.Close()
				e.haveServerConnection = true

				for {
					// Wait for location report
					loc := <-e.loc

					m := &pb.LocationReport{
						ReceiveTs:         timestamppb.Now(),
						TraceletId:        e.deviceID,
						X:                 loc.x,
						Y:                 loc.y,
						Z:                 loc.z,
						SiteId:            12345,
						LocationSignature: 0x12345678ABCDEF,
						CovXx:             11.2,
						CovXy:             22.4,
						CovYy:             -33,
					}
					err := ch.WriteMessage(m)
					if err != nil {
						log.Printf("locationClient WriteMessage failed, %v\n", err)
						break
					}
				}
			}

			time.Sleep(500 * time.Millisecond)
		}
	}()
	return nil
}

func (e *Eloc) locationGenerator() {
	go func() {
		loc := &location{x: -100, y: 0, z: 2}
		stepX := 3.2
		stepY := 3.2

		for {
			// simulate "no satlet reception" for x positions >=80
			loc.havePos = loc.x < 80.0
			e.havePosition = loc.havePos

			log.Printf("locationGenerator: havePos=%t x=%.2f y=%.2f\n", loc.havePos, loc.x, loc.x)

			if loc.havePos {

				// send location to client, don't block if client isn't ready
				select {
				case e.loc <- *loc:
				default:
					log.Printf("locationGenerator: client not ready")
				}

			}

			if loc.y >= 100 || loc.y <= -100 {
				stepY = -stepY
			}
			if loc.x >= 100 || loc.y <= -100 {
				stepX = -stepX
			}
			loc.x += stepX
			loc.y += stepY

			time.Sleep(1000 * time.Millisecond)
		}
	}()

}

func channelFromSocketAddress(address string) (*client.Channel, error) {
	t, err := socket.NewSocketConnection(address)
	if err != nil {
		return nil, errors.New("can't create connection: " + err.Error())
	}
	ms := transport.NewFramedStreamFromTransport(t)
	ch := client.NewChannel(ms)

	return ch, nil
}
