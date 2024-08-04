package ipc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
)

func (ipc *SocketConn) CloseSocket() {
	if ipc.socket != nil {
		_ = ipc.socket.Close()
		ipc.socket = nil
	}
}

// Read the socket response
func (ipc *SocketConn) Read() string {
	buf := make([]byte, 512)
	payloadLength, err := ipc.socket.Read(buf)

	if err != nil {
		fmt.Println("Nothing to read")
	}

	buffer := new(bytes.Buffer)
	for i := 8; i < payloadLength; i++ {
		buffer.WriteByte(buf[i])
	}

	return buffer.String()
}

// Send opcode and payload to the unix socket
func (ipc *SocketConn) Send(opcode int, payload string) (string, error) {
	log.Printf("ipc.Send(): Sending payload: \n%v", payload)

	buf := new(bytes.Buffer)

	err := binary.Write(buf, binary.LittleEndian, int32(opcode))
	if err != nil {
		fmt.Println(err)
	}

	err = binary.Write(buf, binary.LittleEndian, int32(len(payload)))
	if err != nil {
		fmt.Println(err)
	}

	buf.Write([]byte(payload))

	if ipc == nil {
		fmt.Println("ipc.Send(): Tried to send message to unreferenced socket conn, skipping")
		return "", nil
	}

	_, err = ipc.socket.Write(buf.Bytes())
	if err != nil {
		fmt.Println("ipc.Send(): Error writing message", err)
		return "", err
	}

	return ipc.Read(), nil
}
