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
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/ci4rail/io4edge-client-go/client"
	"github.com/ci4rail/io4edge-client-go/transport"
	"github.com/ci4rail/io4edge-client-go/transport/socket"
	pb "github.com/ci4rail/io4edge_api/tracelet/go/tracelet"
	"google.golang.org/protobuf/types/known/timestamppb"
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

// publish location to server periodically
func (e *Tracelet) locationClient(locationServerAddress string) error {
	go func() {
		for {
			log.Printf("try to connect to %v\n", locationServerAddress)
			ch, err := channelFromSocketAddress(locationServerAddress)

			if err == nil {
				defer ch.Close()
				quit := make(chan bool)
				var wg sync.WaitGroup
				wg.Add(1)
				go e.commandHandler(ch, quit, &wg)
				for {
					m := e.makeLocationMessage()
					t2s := e.makeTraceletToServerMessage(0)
					t2s.Type = &pb.TraceletToServer_Location_{Location: m}

					fmt.Printf("locationClient WriteMessage: %v\n", t2s)

					err := ch.WriteMessage(t2s)
					if err != nil {
						log.Printf("locationClient WriteMessage failed, %v\n", err)
						break
					}
					time.Sleep(500 * time.Millisecond)
				}
				select {
				case quit <- true:
				default:
				}
				wg.Wait()
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	return nil
}

func (e *Tracelet) commandHandler(ch *client.Channel, quit chan bool, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-quit:
			log.Print("commandHandler quit\n")
			return
		default:
			m := &pb.ServerToTracelet{}
			err := ch.ReadMessage(m, 0)
			if err != nil {
				log.Printf("commandHandler ReadMessage failed, %v. Exit command handler\n", err)
				return
			}
			t2s := e.makeTraceletToServerMessage(m.Id)

			switch x := m.Type.(type) {
			case *pb.ServerToTracelet_Location:
				{
					fmt.Printf("commandHandler Location: %v\n", m)
					m := e.makeLocationMessage()
					t2s.Type = &pb.TraceletToServer_Location_{Location: m}
				}
			case *pb.ServerToTracelet_Status:
				{
					fmt.Printf("commandHandler Status: %v\n", m)
					m := &pb.TraceletToServer_StatusResponse{
						PowerUpCount:     123,
						HasTime:          true,
						UwbModuleStatus:  0,
						GnssModuleStatus: 0,
						Imu1Status:       0,
						TachoStatus:      777,
					}
					t2s.Type = &pb.TraceletToServer_Status{Status: m}
				}
			default:
				{
					fmt.Printf("commandHandler unknown message type: %v\n", x)
					continue
				}
			}
			err = ch.WriteMessage(t2s)
			if err != nil {
				log.Printf("commandHandler WriteMessage failed, %v\n", err)
			}
		}
	}
}

func (e *Tracelet) makeTraceletToServerMessage(id int32) *pb.TraceletToServer {
	return &pb.TraceletToServer{
		Id:         id,
		TraceletId: e.deviceID,
		ReceiveTs:  timestamppb.Now(),
	}
}

func (e *Tracelet) makeLocationMessage() *pb.TraceletToServer_Location {
	e.locMutex.Lock()
	defer e.locMutex.Unlock()
	return &pb.TraceletToServer_Location{
		Gnss: &pb.TraceletToServer_Location_Gnss{
			Valid:     e.loc.gnss_valid,
			Latitude:  e.loc.gnss_lat,
			Longitude: e.loc.gnss_lon,
			Altitude:  e.loc.gnss_alt,
			Eph:       0.4,
			Epv:       2.5,
		},
		Uwb: &pb.TraceletToServer_Location_Uwb{
			Valid:             e.loc.uwb_valid,
			X:                 e.loc.uwb_x,
			Y:                 e.loc.uwb_y,
			Z:                 e.loc.uwb_z,
			SiteId:            0x1234,
			LocationSignature: 0x12345678ABCD,
			CovXx:             11.1,
			CovXy:             12.2,
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

			time.Sleep(1500 * time.Millisecond)

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
