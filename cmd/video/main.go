package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	"git.campmon.com/golang/corekit/proc"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/platforms/dji/tello"
)

var drone = tello.NewDriver("8888")

const offset = 32767.0

func main() {

	mplayer := exec.Command("mplayer", "-fps", "60", "-")

	mplayerIn, err := mplayer.StdinPipe()
	if nil != err {
		panic(err)
	}

	if err := mplayer.Start(); err != nil {
		panic(err)
	}

	drone.On(tello.ConnectedEvent, func(data interface{}) {
		drone.SetVideoEncoderRate(tello.VideoBitRateAuto)
		drone.StartVideo()
		// it need to send `StartVideo` to the drone every 100ms
		gobot.Every(100*time.Millisecond, func() {
			if err := drone.StartVideo(); nil != err {
				fmt.Printf("fail to start video on drone:%s\n", err)
			}
		})
	})

	drone.On(tello.VideoFrameEvent, func(data interface{}) {
		pkt := data.([]byte)
		if len(pkt) > 0 {
			if _, err := mplayerIn.Write(pkt); err != nil {
				fmt.Printf("err:%s\n", err)
			}
		}
	})

	drone.On(tello.FlightDataEvent, func(data interface{}) {
		fd := data.(*tello.FlightData)
		printFlightData(fd)
	})

	if err := drone.Start(); nil != err {
		fmt.Printf("drone err:%s \n", err)
	}
	donechan := make(chan struct{})
	go droneCtrl(donechan)
	// robot := gobot.NewRobot("bot", []gobot.Device{drone})
	// robot.Start()
	if err := mplayer.Wait(); nil != err {
		panic(err)
	}

	proc.WaitForTermination()
	close(donechan)
	// stop drone
	drone.Halt()

}
func droneCtrl(done chan struct{}) {
	fmt.Println("wait for input")
	r := bufio.NewReader(os.Stdin)

	for {
		select {
		case <-done:
			return
		default:
			item, err := r.ReadString('\n')
			if nil != err {
				fmt.Println("err:", err)
			}
			switch string(bytes.TrimSuffix([]byte(item), []byte{'\n'})) {
			case "t":
				if err := drone.TakeOff(); nil != err {
					fmt.Printf("fail to take off ,err : %s", err)
				}
			case "w":
				if err := drone.Up(5); nil != err {
					fmt.Printf("drone fail to go up,err:%s", err)
				}
			case "s":
				if err := drone.Down(5); nil != err {
					fmt.Printf("drone fail to go down,err:%s", err)
				}
			case "a":
				if err := drone.Left(5); nil != err {
					fmt.Printf("drone fail to go Left,err:%s", err)
				}
			case "d":
				if err := drone.Right(5); nil != err {
					fmt.Printf("drone fail to go right,err:%s", err)
				}
			case "l":
				if err := drone.Land(); nil != err {
					fmt.Printf("drone fail to land,err:%s", err)
				}
			}

		}
	}

}

func printFlightData(d *tello.FlightData) {
	if d.BatteryLow {
		fmt.Printf(" -- Battery low: %d%% --\n", d.BatteryPercentage)
	}

	// 	displayData := `
	// Height:         %d
	// Ground Speed:   %d
	// Light Strength: %d
	// `
	//fmt.Printf(displayData, d.Height, d.GroundSpeed, d.LightStrength)
}
