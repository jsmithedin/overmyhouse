package main

import (
	"bufio"
	"flag"
	"log"
	"net"
	"reflect"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var magicTimestampMLAT = []byte{0xFF, 0x00, 0x4D, 0x4C, 0x41, 0x54}

const (
	aisCharset = "@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_ !\"#$%&'()*+,-./0123456789:;<=>?"
)

var (
	listenAddr = flag.String("bind", "127.0.0.1:8081", "\":port\" or \"ip:port\" to bind the server to")
	baseLat    = flag.Float64("baseLat", 55.910838, "latitude used for distance calculation")
	baseLon    = flag.Float64("baseLon", -3.236900, "longitude for distance calculation")
	mode       = flag.String("mode", "overhead", "overhead or table")
)

func main() {
	log.SetOutput(&lumberjack.Logger{
		Filename:   "overmyhouse.log",
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	})

	log.Println("Starting to watch over my house")

	flag.Parse()

	var knownAircraft KnownAircraft
	var tweetedAircraft TweetedAircraft

	server, _ := net.Listen("tcp", *listenAddr)
	conns := startServer(server)

	ticker := time.NewTicker(500 * time.Millisecond)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				switch *mode {
				case "table":
					printAircraftTable(&knownAircraft)
				default:
					printOverhead(&knownAircraft, &tweetedAircraft)
					tweetedAircraft.PruneTweeted()
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	for {
		go handleConnection(<-conns, &knownAircraft)
	}
}

func startServer(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	go func() {
		for {
			client, _ := listener.Accept()
			if client == nil {
				continue
			}
			ch <- client
		}
	}()
	return ch
}

func handleConnection(conn net.Conn, knownAircraft *KnownAircraft) {
	reader := bufio.NewReader(conn)

	var bufferedMessage []byte

	for {
		currentMessage, _ := reader.ReadBytes(0x1A)

		if len(currentMessage) == 0 {
			break
		}

		if bufferedMessage == nil {
			bufferedMessage = currentMessage
		} else {
			bufferedMessage = append(bufferedMessage, currentMessage...)
		}

		parseBuffer := false
		if currentMessage[0] == 0x31 || currentMessage[0] == 0x32 ||
			currentMessage[0] == 0x33 || currentMessage[0] == 0x34 {
			parseBuffer = true
		}
		if !parseBuffer {
			continue
		}

		message := bufferedMessage
		bufferedMessage = nil

		msgType := message[0]
		var msgLen int

		switch msgType {
		case 0x31:
			continue
		case 0x32:
			continue
		case 0x33:
			msgLen = 22
		case 0x34:
			continue
		default:
			continue
		}

		if len(message) < msgLen {
			continue
		}

		// Not sure if MLAT stuff is necessary
		var timestamp time.Time
		isMlat := reflect.DeepEqual(message[1:7], magicTimestampMLAT)
		if !isMlat {
			timestamp = parseTime(message[1:7])
			_ = timestamp // Why?!
		}

		msgContent := message[8 : len(message)-1]

		parseModeS(msgContent, isMlat, knownAircraft)
	}
}
