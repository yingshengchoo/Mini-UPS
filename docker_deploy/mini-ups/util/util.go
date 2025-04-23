package util

import (
	"encoding/binary"
	"fmt"
	"mini-ups/config"
	"net"
	"sync/atomic"

	"google.golang.org/protobuf/proto"
)

var UPSConn net.Conn

var HOST = config.GetEnvOrDefault("WORLD_HOST", "vcm-47478.vm.duke.edu:12345") // Changed based on HOST

var seqnum int64

func GenerateSeqNum() int64 {
	return atomic.AddInt64(&seqnum, 1)
}

// reserved for future use
func SendMsg(conn net.Conn, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return err
	}

	// Varint encode the length
	var lenBuf [binary.MaxVarintLen32]byte
	n := binary.PutUvarint(lenBuf[:], uint64(len(data)))

	// Send length followed by data
	if _, err := conn.Write(lenBuf[:n]); err != nil {
		return err
	}
	if _, err := conn.Write(data); err != nil {
		return err
	}

	fmt.Printf("Sent %d bytes header + %d bytes data\n", n, len(data))
	return nil
}

func RecvMsg(conn net.Conn, msg proto.Message) error {
	// Read varint length prefix
	var size uint64
	var err error
	var sizeBuf [1]byte
	var shift uint

	for {
		_, err = conn.Read(sizeBuf[:])
		if err != nil {
			return err
		}
		b := sizeBuf[0]
		size |= uint64(b&0x7F) << shift
		if b < 0x80 {
			break
		}
		shift += 7
		if shift >= 64 {
			return fmt.Errorf("recvMsg: varint size too long")
		}
	}

	// Read full message of decoded size
	data := make([]byte, size)
	total := 0
	for total < int(size) {
		n, err := conn.Read(data[total:])
		if err != nil {
			return err
		}
		total += n
	}

	return proto.Unmarshal(data, msg)
}
