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
	"strconv"
)

func iso8583Handler(request *[]byte) *[]byte {
	var isoMsg = *request

	fmt.Printf("Iso Request Message: %s\n", hex.EncodeToString(isoMsg))

	fields := parseRequest(isoMsg)

	// Prepare a response
	mType, _ := strconv.Atoi(fields[0])
	mType += 10
	fields[0] = fmt.Sprintf("%04d", mType)

	fields[37] = "123456789012"
	fields[38] = "123456"
	fields[39] = "00"

	fmt.Println(fields)

	isoPacked := packResponse(fields)

	rspPacked := make([]byte, 0)

	var msgLen = len(isoPacked) + 19
	rspPacked = append(rspPacked, byte(msgLen/256))
	rspPacked = append(rspPacked, byte(msgLen%256))
	rspPacked = append(rspPacked, 0x01)
	rspPacked = append(rspPacked, 0x00)
	rspPacked = append(rspPacked, []byte("123456789012")...)

	tmp, _ := hex.DecodeString(fields[0])
	rspPacked = append(rspPacked, tmp...)

	tmp, _ = hex.DecodeString(fields[3])
	rspPacked = append(rspPacked, tmp...)
	rspPacked = append(rspPacked, isoPacked...)

	return &rspPacked
}

func parseRequest(isoMsg []byte) map[int]string {
	fileds := make(map[int]string)

	idx := 0
	fileds[0] = string(hex.EncodeToString(isoMsg[0:2]))
	idx += 2

	var bitmap = isoMsg[2:10]
	idx += 8

	for i := 1; i < 65; i++ {
		ei := (i - 1) / 8
		ii := (i - 1) % 8

		if bitmap[ei]&(0x80>>uint(ii)) != 0 {
			switch i {
			case 2, 32, 35:
				l, _ := strconv.Atoi(hex.EncodeToString(isoMsg[idx : idx+1]))
				idx++

				if l%2 != 0 {
					l = l/2 + 1
				} else {
					l = l / 2
				}

				fileds[i] = hex.EncodeToString(isoMsg[idx : idx+l])
				idx += l
				break
			case 3, 11, 12:
				fileds[i] = hex.EncodeToString(isoMsg[idx : idx+3])
				idx += 3
				break
			case 4:
				fileds[i] = hex.EncodeToString(isoMsg[idx : idx+6])
				idx += 6
				break
			case 13, 14, 22, 49:
				fileds[i] = hex.EncodeToString(isoMsg[idx : idx+2])
				idx += 2
				break
			case 25:
				fileds[i] = hex.EncodeToString(isoMsg[idx : idx+1])
				idx += 1
				break
			case 37:
				fileds[i] = string(isoMsg[idx : idx+12])
				idx += 12
				break
			case 38:
				fileds[i] = string(isoMsg[idx : idx+6])
				idx += 6
				break
			case 39:
				fileds[i] = string(isoMsg[idx : idx+2])
				idx += 2
				break
			case 41:
				fileds[i] = string(isoMsg[idx : idx+8])
				idx += 8
				break
			case 42:
				fileds[i] = string(isoMsg[idx : idx+15])
				idx += 15
				break
			case 43:
				fileds[i] = string(isoMsg[idx : idx+40])
				idx += 40
				break
			case 52:
				fileds[i] = hex.EncodeToString(isoMsg[idx : idx+8])
				idx += 8
				break
			case 57:
				fileds[i] = string(isoMsg[idx : idx+16])
				idx += 16
				break
			case 48, 55, 62, 63:
				l, _ := strconv.Atoi(hex.EncodeToString(isoMsg[idx : idx+2]))
				idx += 2

				fileds[i] = hex.EncodeToString(isoMsg[idx : idx+l])
				idx += l
				break
			}
		}
	}

	fmt.Println(fileds)
	return fileds
}

func packResponse(fields map[int]string) []byte {
	isoPacked := make([]byte, 0)
	idx := 0

	tmp, _ := hex.DecodeString(fields[0])
	isoPacked = append(isoPacked, tmp...)
	idx += 2
	isoPacked = append(isoPacked, make([]byte, 8)...)
	idx += 8

	for i := 1; i < 65; i++ {
		value, isKeyPresent := fields[i]
		if isKeyPresent {
			ei := (i - 1) / 8
			ii := (i - 1) % 8

			isoPacked[ei+2] |= 0x80 >> uint(ii)

			switch i {
			case 2, 32, 35:
				l := len(value)
				tmp, _ = hex.DecodeString(fmt.Sprintf("%02d", l))
				isoPacked = append(isoPacked, tmp...)
				idx++

				tmp, _ = hex.DecodeString(value)
				isoPacked = append(isoPacked, tmp...)
				idx += l / 2
				break
			case 37, 38, 39, 41, 42, 43:
				tmp = []byte(value)
				isoPacked = append(isoPacked, tmp...)
				idx += len(tmp)
				break
			case 55, 62, 63:
				l := len(value)
				tmp, _ = hex.DecodeString(fmt.Sprintf("%04d", l))
				isoPacked = append(isoPacked, tmp...)
				idx += 2

				tmp, _ = hex.DecodeString(value)
				isoPacked = append(isoPacked, tmp...)
				idx += l / 2
				break
			default:
				tmp, _ = hex.DecodeString(value)
				isoPacked = append(isoPacked, tmp...)
				idx += len(tmp)
				break
			}
		}
	}

	fmt.Printf("Iso Response Message: %s\n", hex.EncodeToString(isoPacked))

	isoPacked, _ = DesEncryption(isoPacked)

	return isoPacked
}
