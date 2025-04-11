import socket
import world_ups_1_pb2  # Make sure this is generated from your .proto file
import time
from google.protobuf.internal.decoder import _DecodeVarint32
from google.protobuf.internal.encoder import _EncodeVarint
from io import BytesIO

HOST = 'vcm-46946.vm.duke.edu'
PORT = 12345

def send_msg(sock, msg):
    data = msg.SerializeToString()
    out = BytesIO()
    _EncodeVarint(out.write, len(data), None)
    out.write(data)
    sock.sendall(out.getvalue())

def recv_msg(sock, message_type):
    var_int_buff = b""
    while True:
        byte = sock.recv(1)
        if not byte:
            return None
        var_int_buff += byte
        msg_len, new_pos = _DecodeVarint32(var_int_buff, 0)
        if new_pos != 0:
            break

    message_data = b""
    while len(message_data) < msg_len:
        chunk = sock.recv(msg_len - len(message_data))
        if not chunk:
            return None
        message_data += chunk

    msg = message_type()
    msg.ParseFromString(message_data)
    return msg

def main():
    with socket.create_connection((HOST, PORT)) as sock:
        # Step 1: Send UConnect
        connect_msg = world_ups_1_pb2.UConnect()
        connect_msg.isAmazon = False
        truck = connect_msg.trucks.add()
        truck.id = 1
        truck.x = 0
        truck.y = 0

        print("Sending UConnect...")
        send_msg(sock, connect_msg)

        # Step 2: Receive UConnected
        response = recv_msg(sock, world_ups_1_pb2.UConnected)
        if response:
            print("Received UConnected:", response)
            worldid = response.worldid
        else:
            print("No response or invalid format.")
            return

        create_truck = world_ups_1_pb2.

        # Step 3: Send UCommands with UGoPickup
        cmd = world_ups_1_pb2.UCommands()
        pickup = cmd.pickups.add()
        pickup.truckid = 1
        pickup.whid = 100
        pickup.seqnum = 1

        cmd.simspeed = 200

        print("Sending UGoPickup...")
        send_msg(sock, cmd)

        # Step 4: Receive UResponses
        resp = recv_msg(sock,  world_ups_1_pb2.UResponses)
        if resp:
            print("Received UResponses:")
            print(resp)
        else:
            print("No UResponses received.")

if __name__ == "__main__":
    main()