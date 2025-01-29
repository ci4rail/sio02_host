#!/usr/bin/env python3

import socket
import time
# ensure that the io4edge_api package is in the PYTHONPATH
import io4edge_api.tracelet.python.v1.tracelet_pb2 as tracelet_pb2
import struct

# map of client address to Client object
CLIENTS = {}

class Client:
    def __init__(self, address: str):
        self.address = address
        self.last_seq = None
        self.last_msg_ts = None
        
    def process_message(self, message: bytes):
        seq = struct.unpack('<L', message[0:4])[0]
        payload = message[4:]
        
        print(f'Client {self.address} received seq={seq}')
        
        if self.last_seq is None or (seq-self.last_seq >= 1 and seq-self.last_seq <= 100):
            self.last_seq = seq
            self.last_msg_ts = time.time()
            self.process_payload(payload)
        else:
            print(f'Client {self.address} ignore dup message: {seq})')    
        
        # ack message
        ack = struct.pack('<L', seq)
        server_socket.sendto(ack, self.address)
        
    def process_payload(self, data: bytes):
        m = tracelet_pb2.TraceletToServer()
        m.ParseFromString(data)
        print(f'message from {m.tracelet_id}, ts={m.delivery_ts.ToDatetime()}')
        t = m.WhichOneof('type')
        if t == 'location':
            loc = m.location
            print(
                f'   FUSED: valid {loc.fused.valid} {loc.fused.latitude:.2f} {loc.fused.longitude:.2f} eph {loc.fused.eph}\n'
                f'     UWB: valid {loc.uwb.valid} {loc.uwb.x:.2f} {loc.uwb.y:.2f} site:{loc.uwb.site_id} eph {loc.uwb.eph}\n'
                f'    GNSS: valid {loc.gnss.valid} {loc.gnss.latitude:.6f} {loc.gnss.longitude:.6f} eph {loc.gnss.eph:.2f}')
        print(f'   metrics: {m.metrics}')    
        print()
            

if __name__ == '__main__':
    server_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    server_address = ('', 11002)
    server_socket.bind(server_address)
    server_socket.settimeout(5)  

    while True:
        print('Waiting to receive message...')
        try:
            data, address = server_socket.recvfrom(4096)
            #print(f'Received {len(data)} bytes from {address}')
            
            if address in CLIENTS:
                client = CLIENTS[address]
            else:
                print(f'New client {address}')
                client = Client(address)
                CLIENTS[address] = client

            client.process_message(data)
        except socket.timeout:
            print('Timeout waiting for message')
            pass

