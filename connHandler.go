/*
 *   Copyright (c) 2024 thearistotlemethod@gmail.com
 *   All rights reserved.

 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at

 *   http://www.apache.org/licenses/LICENSE-2.0

 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"sync/atomic"
	"time"
)

type connHandler struct {
	conn net.Conn
}

func (hdl connHandler) handler() {
	// Close the connection when we're done
	defer func() {
		atomic.AddInt32(&activeConnectionCount, -1)
		hdl.conn.Close()
	}()

	fmt.Printf("New connection from %s\n", hdl.conn.RemoteAddr().String())

	atomic.AddInt32(&activeConnectionCount, 1)

	for {
		buf := readRequest(hdl.conn)
		if buf == nil {
			return
		}

		// Print the incoming data
		fmt.Printf("Received: %s\n", hex.EncodeToString(buf))

		request, err := DesDecryption(buf[21:])
		if err != nil {
			fmt.Println("Decryption failed:", err)
			return
		}

		reponse := *iso8583Handler(&request)

		fmt.Printf("Send: %s\n", hex.EncodeToString(reponse))
		_, err = hdl.conn.Write(reponse)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func readRequest(conn net.Conn) []byte {
	// Read incoming data
	buf := make([]byte, 1024*10)
	idx := 0
	length := 0
	for {
		e := conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		if e != nil {
			fmt.Println("SetReadDeadline failed:", e)
			return nil
		}

		l, e := conn.Read(buf[idx:])
		if e != nil {
			fmt.Println(e)
			return nil
		}

		if idx == 0 {
			if l >= 2 {
				length = int(buf[0])*256 + int(buf[1])
				length += 2
			} else {
				fmt.Println("Invalid data length")
				return nil
			}
			idx += l
		} else {
			idx += l
		}

		if idx >= length {
			break
		}
	}

	return buf[0:length]
}
