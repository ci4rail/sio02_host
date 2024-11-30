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
	"math/rand"
	"sync"
	"time"

	"github.com/ci4rail/io4edge-client-go/client"
	pb "github.com/ci4rail/io4edge_api/tracelet/go/tracelet"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type location struct {
	uwbValid  bool
	uwbX      float64
	uwbY      float64
	uwbZ      float64
	gnssValid bool
	gnssLat   float64
	gnssLon   float64
	gnssAlt   float64
	gnssFix   int32
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
				for {
					m := e.makeLocationMessage()
					t2s := e.makeTraceletToServerMessage(0)
					t2s.Type = &pb.TraceletToServer_Location_{Location: m}
					t2s.Metrics = makeMetricsMessage()

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
			time.Sleep(1000 * time.Millisecond)
		}
	}()
	return nil
}

func (e *Tracelet) makeTraceletToServerMessage(id int32) *pb.TraceletToServer {
	return &pb.TraceletToServer{
		Id:         id,
		TraceletId: e.deviceID,
		Ignition:   true,
		DeliveryTs: timestamppb.Now(),
	}
}

func (e *Tracelet) makeLocationMessage() *pb.TraceletToServer_Location {
	e.locMutex.Lock()
	defer e.locMutex.Unlock()
	return &pb.TraceletToServer_Location{
		Gnss: &pb.TraceletToServer_Location_Gnss{
			Valid:     e.loc.gnssValid,
			Latitude:  e.loc.gnssLat,
			Longitude: e.loc.gnssLon,
			Altitude:  e.loc.gnssAlt,
			Eph:       0.4,
			Epv:       2.5,
			FixType:   e.loc.gnssFix,
		},
		Uwb: &pb.TraceletToServer_Location_Uwb{
			Valid:             e.loc.uwbValid,
			X:                 e.loc.uwbX,
			Y:                 e.loc.uwbY,
			Z:                 e.loc.uwbZ,
			SiteId:            0x1234,
			LocationSignature: 0x12345678ABCD,
			Eph:               0.6,
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
				uwbValid:  true,
				uwbX:      5.0,
				uwbY:      6.21,
				uwbZ:      7.5,
				gnssValid: false,
				gnssLat:   49.425111,
				gnssLon:   11.077378,
				gnssAlt:   350.0,
				gnssFix:   0,
			}
			e.locMutex.Lock()
			e.loc = loc
			e.locMutex.Unlock()

			time.Sleep(1000 * time.Millisecond)

			loc = location{
				uwbValid:  false,
				uwbX:      0,
				uwbY:      1100,
				uwbZ:      888,
				gnssValid: true,
				gnssLat:   49.425111,
				gnssLon:   11.077378,
				gnssAlt:   350.0,
				gnssFix:   2,
			}
			e.locMutex.Lock()
			e.loc = loc
			e.locMutex.Unlock()

			time.Sleep(1500 * time.Millisecond)

		}
	}()

}

// generate some random metrics
func makeMetricsMessage() *pb.TraceletMetrics {
	return &pb.TraceletMetrics{
		Health__Type__UwbComm:     1,
		Health__Type__UwbFirmware: 0,
		Health__Type__GnssComm:    1,
		FreeHeapBytes:             int64(rand.Intn(1000) + 20000),

		GnssFixTypeEnum:                 int64(rand.Intn(6)),
		GnssHeading__Info__HeadVehValid: 1,
		GnssHeading__Info__HeadVeh:      float64(rand.Intn(360)),
		GnssHeading__Info__HeadMot:      float64(rand.Intn(360)),
		NtripIsConnected:                0,
		SpeedMetersPerSecond:            float64(rand.Intn(100)),
	}
}

func channelFromSocketAddress(address string) (*client.Channel, error) {
	c, err := client.NewUDPClientFromSocketAddress(address)
	if err != nil {
		return nil, errors.New("can't create UDP client: " + err.Error())
	}

	return c.Ch, nil
}
