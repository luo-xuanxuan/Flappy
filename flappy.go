package main

import (
	"log"
	"math"
	"net"
	"time"

	"github.com/go-gl/glfw/v3.3/glfw"
)

type Pipe struct {
	Position float64
	Opening  int64
}

var (
	gameboard [128][64]bool

	flappyVirtualPosition = 32.0
	flappyVisualPosition  = 32
	flappyVelocity        = 0.0
	flappyJump            = 1.0
	gravity               = 1.5

	pipeSpeed    = 1.0
	pipeGap      = 20
	pipeDistance = 20
)

var serverAddr = "192.168.1.134:4210"

var isJumpPressed = false

var lastFrameTime int64 = 0

func processFlappy(input []glfw.Action) {

	var now int64 = time.Now().UnixMilli()

	if lastFrameTime == 0 {
		lastFrameTime = now
	}

	var delta float64 = float64(now-lastFrameTime) / 1000.0
	lastFrameTime = now

	if (input[0] == 1) && (!isJumpPressed) {
		flappyVelocity -= flappyJump
		isJumpPressed = true
	}
	if input[0] == 0 {
		isJumpPressed = false
	}
	// Apply gravity
	flappyVelocity += gravity * delta
	flappyVirtualPosition += flappyVelocity

	// Collision with the ground
	if int(flappyVirtualPosition) >= len(gameboard[0])-2 {
		flappyVirtualPosition = float64(len(gameboard[0]) - 2)
		flappyVelocity = 0
	}

	// Collision with the ceiling
	if flappyVirtualPosition < 0 {
		flappyVirtualPosition = 0
		flappyVelocity = 0
	}

	flappyVisualPosition = int(math.Round(flappyVirtualPosition))

	updateGameBoard()

	addr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		log.Fatal(err)
	}

	// Create a UDP connection
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write(gameBoardToByteArray())
	if err != nil {
		log.Fatal(err)
	}

	//clear console and write debug values?
	/*clearConsole()
	fmt.Println("Virtual : " + strconv.FormatFloat(flappyVirtualPosition, 'f', -1, 64))
	fmt.Println("Visual  : " + strconv.Itoa(flappyVisualPosition))
	fmt.Println("Velocity: " + strconv.FormatFloat(flappyVelocity, 'f', -1, 64))
	fmt.Println("Delta   : " + strconv.FormatFloat(delta, 'f', -1, 64))*/

	time.Sleep(time.Millisecond)

}

func updateGameBoard() {
	// Clear the board
	for i := range gameboard {
		for j := range gameboard[i] {
			gameboard[i][j] = false
		}
	}

	// Set the new position of Flappy
	if flappyVisualPosition >= 0 && flappyVisualPosition < len(gameboard[0]) {
		gameboard[5][flappyVisualPosition] = true
		gameboard[5][flappyVisualPosition+1] = true
		gameboard[4][flappyVisualPosition] = true
		gameboard[4][flappyVisualPosition+1] = true
	}
}

func gameBoardToByteArray() []byte {
	bytes := make([]byte, 8192/8) // 8192 bits are needed, 8192 bits / 8 bits per byte = 1024 bytes
	for j := 0; j < 64; j++ {
		for i := 0; i < 128; i++ {
			if gameboard[i][j] {
				byteIndex := (j*128 + i) / 8
				bitIndex := (j*128 + i) % 8
				bytes[byteIndex] |= (1 << bitIndex)
			}
		}
	}
	return bytes
}

/*
func clearConsole() {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("cmd", "/c", "cls")
	default:
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	cmd.Run()
}*/
