#!/usr/bin/env python3

import socketserver
import threading
import tracelet_location_pb2
import struct


class MyTCPHandler(socketserver.BaseRequestHandler):
    """
    The RequestHandler class for our server.

    It is instantiated once per connection to the server, and must
    override the handle() method to implement communication to the
    client.
    """

    def handle(self):
        print('new handler %s\n' % threading.current_thread().name)
        while True:
            loc = self.read_fstream()
            print(
                f'{loc.tracelet_id} {loc.receive_ts.ToDatetime()} {loc.x:.2f} {loc.y:.2f} {loc.z:.2f} '
                f'cov: {loc.cov_xx:.2f} {loc.cov_xy:.2f} {loc.cov_yy:.2f} site:{loc.site_id} sign: {loc.location_signature}')

        print('exit handler %s\n' % threading.current_thread().name)

    def rcv_all(self, n):
        remaining = n
        buf = bytearray()
        while remaining > 0:
            data = self.request.recv(remaining)
            buf.extend(data)
            remaining -= len(data)
        return buf

    def read_fstream(self):
        hdr = self.rcv_all(6)
        if hdr[0:2] == b'\xfe\xed':
            len = struct.unpack('<L', hdr[2:6])[0]
            # print(f'len={len} {hdr[0:6]}')
            proto_data = self.rcv_all(len)
            loc = tracelet_location_pb2.LocationReport()
            loc.ParseFromString(proto_data)
            return loc
        else:
            raise RuntimeError('bad magic')


class ThreadedTCPServer(socketserver.ThreadingMixIn, socketserver.TCPServer):
    pass


if __name__ == '__main__':
    HOST, PORT = '0.0.0.0', 11002

    # Create the server, binding to localhost on port 9999
    server = ThreadedTCPServer((HOST, PORT), MyTCPHandler)

    # Activate the server; this will keep running until you
    # interrupt the program with Ctrl-C
    server.serve_forever()
