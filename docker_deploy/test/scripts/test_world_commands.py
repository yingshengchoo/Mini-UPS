import socket
import world_ups_1_pb2  # Make sure this is generated from your .proto file
import world_amazon_1_pb2 as amazon_pb2
import time
from google.protobuf.internal.decoder import _DecodeVarint32
from google.protobuf.internal.encoder import _EncodeVarint
from io import BytesIO

# World Command Interaction Overview
# Note: UPS and AMAZON commands are inside A/UCommands and World Responses are in A/UResponses
#
# # # # # # # # # # # # # # # # # # # # # # # # # # # # #
#      UPS     |     AMAZON    |          World         #
#-------------------------------------------------------#
#   UConnect   |    AConnect   |                        # Truck and warehouses initialized here (需要在)
#              |               | UConnected, AConnected #
#              | APurchaseMore |                        # ¡¡¡Must have stock to pack!!! (要有貨才能包裝)
#              |               |      APurchaseMore     # if there is stock already, no need to purchase (如果有貨可以不用購買)
#              |     APack     |                        #
#              |               |         APacked        #
#   UGoPickUp  |               |                        # ¿¿¿UGoPickup should only be done when packages are packed??? (貨車需要再包裹包奘好才能到倉庫)
#              |               |        UFinished       # <-- Status = "ARRIVE WAREHOUSE"
#              |     ALoad     |                        #
#              |               |         ALoaded        #
#  UGoDeliver  |               |                        #
#              |               |        UFinished       # <-- Status = "IDLE"
# UDisconnect  |  ADisconnect  |                        #
#              |               |       finished x2      # <-- finished = True
# # # # # # # # # # # # # # # # # # # # # # # # # # # # #
#
# Aside from the commands above : A/UQuery (UPS returns UTruck, Amazon returns APackage)
# Aside from the responses above: A/UErr, APackage, UTruck
# HOST = 'vcm-46755.vm.duke.edu'
HOST = 'vcm 47478.vm.duke.edu'
UPS_PORT = 12345
AMAZON_PORT = 23456
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
    sock = socket.create_connection((HOST, AMAZON_PORT))
    print("Sending AConnect...")
    aconnect = amazon_pb2.AConnect()
    aconnect.isAmazon = True

    aconnect.worldid = 1
    wh = aconnect.initwh.add()
    wh.id = 1
    wh.x = 10
    wh.y = 20

    send_msg(sock, aconnect)
    response = recv_msg(sock, amazon_pb2.AConnected)
    print("AConnected:", response)
    return sock


def connect_ups(worldid):
    sock = socket.create_connection((HOST, UPS_PORT))
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
    #make sure that it receives packages so that Truck can go pick it up
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
    while True:
        response = recv_msg(sock, world_ups_1_pb2.UResponses)
        if response is None:
            print("UPS socket closed unexpectedly.")
            return False

        print("UPS received UResponses:", response)
        for completed in response.completions:
            if completed.status == "ARRIVE WAREHOUSE":
                print("Package packed and ready:", completed)

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

def simulate_amazon_disconnect(sock):
    cmd = amazon_pb2.ACommands()
    cmd.disconnect = True
    send_msg(sock, cmd)
    while True:
        response = recv_msg(sock, amazon_pb2.AResponses)
        if response is None:
            print("Amazon socket closed unexpectedly.")
            return False

        print("Amazon received AResponses:", response)

        if response.finished == True:
            print("Amazon disconnected gracefully.")
            return True

        time.sleep(0.5)

def simulate_ups_disconnect(sock):
    cmd = world_ups_1_pb2.UCommands()
    cmd.disconnect = True
    send_msg(sock, cmd)
    while True:
        response = recv_msg(sock, world_ups_1_pb2.UResponses)
        if response is None:
            print("UPS socket closed unexpectedly.")
            return False

        print("UPS received UResponses:", response)

        if response.finished == True:
            print("UPS disconnected gracefully.")
            return True

        time.sleep(0.5)



def main():
    print("Full simulation beginning...")
    amazon_sock = connect_amazon()
    # ups_sock = connect_ups(worldid)

    # simulate_amazon_flow(amazon_sock)
    # simulate_ups_pickup(ups_sock)
    # simulate_amazon_load_and_put(amazon_sock)
    # simulate_ups_deliver(ups_sock)
    # simulate_amazon_disconnect(amazon_sock)
    # simulate_ups_disconnect(ups_sock)
    # amazon_sock.close()
    # ups_sock.close()
    while True :
        pass

    print("Simulation complete")

if __name__ == "__main__":
    main()