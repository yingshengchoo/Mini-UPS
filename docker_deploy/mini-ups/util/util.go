package util

import (
	"encoding/binary"
	"fmt"
	"log"
	"mini-ups/config"
	"net"
	"sync/atomic"

	"google.golang.org/protobuf/proto"
)

//This file contains helper functions and global variables used throughout our UPS.

var WORLD_HOST = config.GetEnvOrDefault("WORLD_HOST", "vcm-47478.vm.duke.edu:12345") // Changed based on HOST
// var UPS_HOST = config.GetEnvOrDefault("UPS_HOST", "vcm-47478.vm.duke.edu:8080")      // Changed based on HOST
var UPS_HOST = config.AppConfig.Ups.Host // Changed based on HOST

var UPSConn net.Conn

var seqnum int64

var globalPackageID int64

// Generates a unique packageID
func GeneratePackageID() int64 {
	return atomic.AddInt64(&globalPackageID, 1)
}

// Geneartea a unique Sequence number.
func GenerateSeqNum() int64 {
	return atomic.AddInt64(&seqnum, 1)
}

// Sends message to the World simulation
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

	log.Printf("Sent %d bytes header + %d bytes data\n", n, len(data))
	return nil
}

// Recieves message from World
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

// simple set datastructure with only the Add function because that is all we need.
type Int64Set struct {
	data map[int64]struct{}
}

// This is a set datastructure to hold INT64.
func NewInt64Set() *Int64Set {
	return &Int64Set{
		data: make(map[int64]struct{}),
	}
}

// Add adds a number to our set.
func (s *Int64Set) Add(val int64) bool {
	if _, exists := s.data[val]; exists {
		return false
	}
	s.data[val] = struct{}{}
	return true
}
