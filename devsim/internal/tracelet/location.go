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

package tracelet

import (
	"errors"
	"log"
	"time"

	"github.com/ci4rail/io4edge-client-go/client"
	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge-client-go/transport/socket"
	pb "github.com/ci4rail/io4edge_api/tracelet/go/tracelet"
)

type location struct {
	uwb_valid  bool
	uwb_x      float64
	uwb_y      float64
	uwb_z      float64
	gnss_valid bool
	gnss_lat   float64
	gnss_lon   float64
	gnss_alt   float64
}

func (e *Tracelet) locationClient(locationServerAddress string) error {
	go func() {
		for {
			e.haveServerConnection = false
			ch, err := channelFromSocketAddress(locationServerAddress)

			if err == nil {
				defer ch.Close()
				e.haveServerConnection = true

				for {
					m := e.makeLocationMessage()

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

func (e *Tracelet) makeLocationMessage() *pb.TraceletToServer_Location {
	e.locMutex.Lock()
	defer e.locMutex.Unlock()
	return &pb.TraceletToServer_Location{
		Gnss: &pb.TraceletToServer_Location_Gnss{
			Valid:    e.loc.gnss_valid,
			Latitude: e.loc.gnss_lat, 
			Longitude: e.loc.gnss_lon, 
			Altitude: e.loc.gnss_alt,
			Eph: 0.4,
			Epv: 2.5,
		},
		Uwb:         &pb.TraceletToServer_Location_Uwb{
			Valid: e.loc.uwb_valid,
			X:     e.loc.uwb_x,
			Y:     e.loc.uwb_y,
			Z:     e.loc.uwb_z,
			SiteId: 0x1234,
			LocationSignature: 0x12345678ABCD,
			CovXx: 11.1,
			CovXy: 12.2,
		},
		Direction:   pb.TraceletToServer_Location_NO_DIRECTION,
		Speed:       9,
		Mileage:     50899,
		Temperature: 34.5,
	}
}

func (e *Tracelet) locationGenerator() {
	go func() {

		for {
			loc := location{
				uwb_valid:  true,
				uwb_x:      5.0,
				uwb_y:      6.21,
				uwb_z:      7.5,
				gnss_valid: false,
				gnss_lat:   49.425111,
				gnss_lon:   11.077378,
				gnss_alt:   350.0,
			}
			e.locMutex.Lock()
			e.loc = loc
			e.locMutex.Unlock()

			time.Sleep(1000 * time.Millisecond)

			loc = location{
				uwb_valid:  false,
				uwb_x:      0,
				uwb_y:      1100,
				uwb_z:      888,
				gnss_valid: true,
				gnss_lat:   49.425111,
				gnss_lon:   11.077378,
				gnss_alt:   350.0,
			}
			e.locMutex.Lock()
			e.loc = loc
			e.locMutex.Unlock()

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
