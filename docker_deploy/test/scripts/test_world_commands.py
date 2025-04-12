import socket
import world_ups_1_pb2  # Make sure this is generated from your .proto file
import world_amazon_1_pb2 as amazon_pb2
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


def connect_amazon():
    sock = socket.create_connection(("vcm-46946.vm.duke.edu", 23456))
    print("Sending AConnect...")
    aconnect = amazon_pb2.AConnect()
    aconnect.isAmazon = True
    wh = aconnect.initwh.add()
    wh.id = 1
    wh.x = 10
    wh.y = 20

    send_msg(sock, aconnect)
    response = recv_msg(sock, amazon_pb2.AConnected)
    print("AConnected:", response)
    return sock, response.worldid


def connect_ups(worldid):
    sock = socket.create_connection(("vcm-46946.vm.duke.edu", 12345))
    print("Sending UConnect...")
    uconnect = world_ups_1_pb2.UConnect()
    uconnect.isAmazon = False
    uconnect.worldid = worldid
    truck = uconnect.trucks.add()
    truck.id = 1
    truck.x = 0
    truck.y = 0

    send_msg(sock, uconnect)
    response = recv_msg(sock, world_ups_1_pb2.UConnected)
    print("UConnected:", response)
    return sock


def simulate_amazon_flow(sock, seq_start=1):
    cmd = amazon_pb2.ACommands()

    # 1. PurchaseMore
    buy = cmd.buy.add()
    buy.whnum = 1
    buy.seqnum = seq_start
    prod = buy.things.add()
    prod.id = 1001
    prod.description = "Test Product"
    prod.count = 2

    # 2. Pack
    pack = cmd.topack.add()
    pack.whnum = 1
    pack.shipid = 999
    pack.seqnum = seq_start + 1
    item = pack.things.add()
    item.id = 1001
    item.description = "Test Product"
    item.count = 2

    send_msg(sock, cmd)
    print("Sent APurchaseMore and APack")
    # make sure that it receives packages so that Truck can go pick it up
    while True:
        response = recv_msg(sock, amazon_pb2.AResponses)
        if response is None:
            print("Amazon socket closed unexpectedly.")
            return False

        print("Amazon received AResponses:", response)

        for packed in response.ready:
            if packed.shipid == 999:
                print("Package packed and ready:", packed)
                return True

        time.sleep(0.5)


def simulate_ups_pickup(sock, seq_start=1):
    cmd = world_ups_1_pb2.UCommands()

    # 3. UPS goes to pick up
    pickup = cmd.pickups.add()
    pickup.truckid = 1
    pickup.whid = 1
    pickup.seqnum = seq_start

    send_msg(sock, cmd)
    print("Sent UGoPickup")
    recv_msg(sock, world_ups_1_pb2.UResponses)
    while True:
        response = recv_msg(sock, world_ups_1_pb2.UResponses)
        if response is None:
            print("UPS socket closed unexpectedly.")
            return False

        print("UPS received UResponses:", response)
        for completed in response.completions:
            if completed.status == "ARRIVE WAREHOUSE":
                print("Package packed and ready:", completed)
                return True

        time.sleep(0.5)



def simulate_amazon_load_and_put(sock, seq_start=3):
    cmd = amazon_pb2.ACommands()

    # 4. APutOnTruck
    put = cmd.load.add()
    put.whnum = 1
    put.truckid = 1
    put.shipid = 999
    put.seqnum = seq_start

    send_msg(sock, cmd)
    print("Sent APutOnTruck")
        # make sure that it receives packages so that Truck can go pick it up
    while True:
        response = recv_msg(sock, amazon_pb2.AResponses)
        if response is None:
            print("Amazon socket closed unexpectedly.")
            return False

        print("Amazon received AResponses:", response)

        for loaded in response.loaded:
            if loaded.shipid == 999:
                print("Package packed and ready:", loaded)
                return True

        time.sleep(0.5)


def simulate_ups_deliver(sock, seq_start=2):
    cmd = world_ups_1_pb2.UCommands()

    # 5. UGoDeliver
    deliver = cmd.deliveries.add()
    deliver.truckid = 1
    deliver.seqnum = seq_start
    loc = deliver.packages.add()
    loc.packageid = 999
    loc.x = 15
    loc.y = 25

    send_msg(sock, cmd)
    print("Sent UGoDeliver")
    #
    while True:
        response = recv_msg(sock, world_ups_1_pb2.UResponses)
        if response is None:
            print("UPS socket closed unexpectedly.")
            return False

        print("UPS received UResponses:", response)

        for delivered in response.delivered:
            if delivered.packageid == 999:
                print("Package Delivered:", delivered)
                return True

        time.sleep(0.5)


def main():
    print("Full simulation beginning...")
    amazon_sock, worldid = connect_amazon()
    ups_sock = connect_ups(worldid)

    simulate_amazon_flow(amazon_sock)
    simulate_ups_pickup(ups_sock)
    simulate_amazon_load_and_put(amazon_sock)
    simulate_ups_deliver(ups_sock)

    amazon_sock.close()
    ups_sock.close()
    print("Simulation complete")

if __name__ == "__main__":
    main()