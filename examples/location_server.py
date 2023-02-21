#!/usr/bin/env python3

import socketserver
import time
import threading
# ensure that the io4edge_api package is in the PYTHONPATH
import io4edge_api.tracelet.python.v1.tracelet_pb2 as tracelet_pb2
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
        self.command_thread_exit = False
        self.command_thread = threading.Thread(target=self.command_requester)
        self.command_thread.start()
        while True:
            m = self.read_fstream()
            loc = m.location
            #print(loc)
            print(
                f'{m.tracelet_id} {m.receive_ts.ToDatetime()}\n'
                f'   UWB: valid {loc.uwb.valid} {loc.uwb.x:.2f} {loc.uwb.y:.2f} ite:{loc.uwb.site_id}\n'
                f'  GNSS  valid {loc.gnss.valid} {loc.gnss.latitude:.6f} {loc.gnss.longitude:.6f} {loc.gnss.eph:.2f}\n')

        print('exit handler %s\n' % threading.current_thread().name)

    def server_close(self):
        print('server close')
        self.command_thread_exit = True
        self.command_thread.join()
        super().server_close()    

    def command_requester(self):
        while not self.command_thread_exit:
            time.sleep(3)
            statusReq = tracelet_pb2.ServerToTracelet.StatusRequest()
            req = tracelet_pb2.ServerToTracelet(id=1)
            req.status.CopyFrom(statusReq)
            data = req.SerializeToString()
            self.send(data)
        print("exit command requester")

    def send(self, data):
        hdr = struct.pack("<HL", 0xEDFE, len(data))
        self.request.sendall(hdr + data)

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
            loc = tracelet_pb2.TraceletToServer()
            loc.ParseFromString(proto_data)
            return loc
        else:
            raise RuntimeError('bad magic')


class ThreadedTCPServer(socketserver.ThreadingMixIn, socketserver.TCPServer):
    pass


if __name__ == '__main__':
    HOST, PORT = '0.0.0.0', 11002

    # Create the server, binding to localhost on specified port
    server = ThreadedTCPServer((HOST, PORT), MyTCPHandler)



    # Activate the server; this will keep running until you
    # interrupt the program with Ctrl-C
    server.serve_forever()
