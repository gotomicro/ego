// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package emetric

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/gotomicro/ego/core/elog"
)

type tcpConnectionState int

const (
	// TCP_ESTABLISHED
	tcpEstablished tcpConnectionState = iota + 1
	// TCP_SYN_SENT
	tcpSynSent
	// TCP_SYN_RECV
	tcpSynRecv
	// TCP_FIN_WAIT1
	tcpFinWait1
	// TCP_FIN_WAIT2
	tcpFinWait2
	// TCP_TIME_WAIT
	tcpTimeWait
	// TCP_CLOSE
	tcpClose
	// TCP_CLOSE_WAIT
	tcpCloseWait
	// TCP_LAST_ACK
	tcpLastAck
	// TCP_LISTEN
	tcpListen
	// TCP_CLOSING
	tcpClosing
	// TCP_RX_BUFFER
	//tcpRxQueuedBytes
	// TCP_TX_BUFFER
	//tcpTxQueuedBytes
)

type TcpStatCollector struct {
}

// NewTCPStatCollector returns a new Collector exposing network stats.
func NewTCPStatCollector() (*TcpStatCollector, error) {
	return &TcpStatCollector{}, nil
}

func (c *TcpStatCollector) Update() error {
	go func() {
		statsFile := path.Join("/proc", strconv.Itoa(os.Getpid()), "net", "tcp")
		tcpStats, err := getTCPStats(statsFile)
		if err != nil {
			elog.EgoLogger.Error(fmt.Errorf("couldn't get tcpstats: %w", err).Error())
			return
		}
		for {
			for index, value := range tcpStats {
				for addr, number := range value {
					ClientStatsGauge.WithLabelValues(
						"conn_states",
						addr,
						index.String(),
					).Set(number)
				}
				time.Sleep(5 * time.Second)
			}
		}
	}()

	// if enabled ipv6 system
	//tcp6File := procFilePath("net/tcp6")
	//if _, hasIPv6 := os.Stat(tcp6File); hasIPv6 == nil {
	//	tcp6Stats, err := getTCPStats(tcp6File)
	//	if err != nil {
	//		return fmt.Errorf("couldn't get tcp6stats: %w", err)
	//	}
	//
	//	for st, value := range tcp6Stats {
	//		tcpStats[st] += value
	//	}
	//}

	//for st, value := range tcpStats {
	//	ch <- c.desc.mustNewConstMetric(value, st.String())
	//}
	return nil
}

func getTCPStats(statsFile string) (map[tcpConnectionState]map[string]float64, error) {
	file, err := os.Open(statsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return parseTCPStats(file)
}

/*
46: 010310AC:9C4C 030310AC:1770 01
   |      |      |      |      |   |--> connection state（套接字状态）
   |      |      |      |      |------> remote TCP port number（远端端口，主机字节序）
   |      |      |      |-------------> remote IPv4 address（远端IP，网络字节序）
   |      |      |--------------------> local TCP port number（本地端口，主机字节序）
   |      |---------------------------> local IPv4 address（本地IP，网络字节序）
   |----------------------------------> number of entry
*/

func parseTCPStats(r io.Reader) (map[tcpConnectionState]map[string]float64, error) {
	//tcpStats := map[tcpConnectionState]float64{}
	tcpStatsMap := make(map[tcpConnectionState]map[string]float64, 0)
	contents, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(contents), "\n")[1:] {
		parts := strings.Fields(line)
		if len(parts) == 0 {
			continue
		}
		if len(parts) < 5 {
			return nil, fmt.Errorf("invalid TCP stats line: %q", line)
		}

		//qu := strings.Split(parts[4], ":")
		//if len(qu) < 2 {
		//	return nil, fmt.Errorf("cannot parse tx_queues and rx_queues: %q", line)
		//}
		//
		//tx, err := strconv.ParseUint(qu[0], 16, 64)
		//if err != nil {
		//	return nil, err
		//}
		//tcpStats[tcpConnectionState(tcpTxQueuedBytes)] += float64(tx)
		//
		//rx, err := strconv.ParseUint(qu[1], 16, 64)
		//if err != nil {
		//	return nil, err
		//}
		//tcpStats[tcpConnectionState(tcpRxQueuedBytes)] += float64(rx)

		ipv4, _ := parseIpV4(parts[2])
		st, err := strconv.ParseInt(parts[3], 16, 8)
		if err != nil {
			return nil, err
		}

		info, flag := tcpStatsMap[tcpConnectionState(st)]
		if !flag {
			info = make(map[string]float64, 0)
			info[ipv4] = 1
		} else {
			info[ipv4]++
		}
		tcpStatsMap[tcpConnectionState(st)] = info

		//tcpStats[tcpConnectionState(st)]++

	}

	return tcpStatsMap, nil
}

func (st tcpConnectionState) String() string {
	switch st {
	case tcpEstablished:
		return "established"
	case tcpSynSent:
		return "syn_sent"
	case tcpSynRecv:
		return "syn_recv"
	case tcpFinWait1:
		return "fin_wait1"
	case tcpFinWait2:
		return "fin_wait2"
	case tcpTimeWait:
		return "time_wait"
	case tcpClose:
		return "close"
	case tcpCloseWait:
		return "close_wait"
	case tcpLastAck:
		return "last_ack"
	case tcpListen:
		return "listen"
	case tcpClosing:
		return "closing"
	//case tcpRxQueuedBytes:
	//	return "rx_queued_bytes"
	//case tcpTxQueuedBytes:
	//	return "tx_queued_bytes"
	default:
		return "unknown"
	}
}

// 只解析IPV4
// 34190A0A:3D2D
func parseIpV4(s string) (string, error) {
	if len(s) == 13 {
		return "", fmt.Errorf("not ipv4")
	}
	hexIP := s[:len(s)-5]
	hexPort := s[len(s)-4:]
	bytesIP, err := hex.DecodeString(hexIP)
	if err != nil {
		return "", nil
	}
	uint32IP := binary.LittleEndian.Uint32(bytesIP) //转换为主机字节序
	IP := make(net.IP, 4)
	binary.BigEndian.PutUint32(IP, uint32IP)
	port, err := strconv.ParseUint(hexPort, 16, 16)
	return fmt.Sprintf("%s:%d", IP.String(), port), err
}

//func parsePort(portStr string) (int64, error) {
//	return strconv.ParseInt(portStr, 16, 16)
//}
