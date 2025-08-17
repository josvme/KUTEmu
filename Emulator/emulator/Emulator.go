package emulator

import (
	"fmt"
	"log"
	"maps"
	"os"
	"riscv/instructions"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type Emulator struct {
	cpu      *instructions.Cpu
	window   *sdl.Window
	renderer *sdl.Renderer
	texture  *sdl.Texture
	running  bool
}

const VIRT_DRAM = 0x80000000
const VIRT_OPENSBI_START = 0x80200000
const VIRT_VIRTIO = 0x10001000
const DTB = 0x87e00000
const SCREEN_WIDTH = 320
const SCREEN_HEIGHT = 200

func NewEmulator() *Emulator {
	uart := instructions.NewUART()
	clint := &instructions.Clint{}
	plic := &instructions.Plic{
		Priority:  1,
		Pending:   make([]uint32, 0),
		Enable:    0,
		Threshold: 1,
		Claim:     0,
		Complete:  0,
	}
	screen := make(map[uint32]uint32)
	disp := &instructions.Display{
		Screen: screen,
		Mutex:  sync.Mutex{},
	}
	memory := &instructions.Memory{Map: make(map[uint32]byte), Uart: uart, Plic: plic, Clint: clint, Display: disp}
	csr := &instructions.CSR{
		Registers: make([]uint32, 4096),
	}
	cpu := &instructions.Cpu{
		PC:          0x80000000,
		Registers:   [32]uint32{},
		Memory:      memory,
		CSR:         csr,
		CurrentMode: 3,
	}
	memory.SetCpu(cpu)
	cpu.Memory = memory
	return &Emulator{
		cpu:      cpu,
		window:   nil,
		renderer: nil,
		texture:  nil,
		running:  true,
	}
}

func (e *Emulator) UpdateTime() {
	for {
		time.Sleep(1 * time.Millisecond)
		e.cpu.Memory.Clint.Mtime = uint64(time.Now().UnixMilli())
	}
}

func (e *Emulator) Run() {
	//f, _ := os.Create("e.log")
	//defer f.Close()

	//body, _ := os.ReadFile("./../C/OS/fw_dynamic.bin")
	path := os.Getenv("OBJ_PATH")
	//path = "/home/josv/Projects/RiscV/Tests/os.img"
	path = "/home/josv/Projects/RiscV/Tests/doom-riscv.bin"
	body, _ := os.ReadFile(path)
	_ = e.cpu.Memory.LoadBytes(body, VIRT_DRAM)

	memory := e.cpu.Memory
	cpu := e.cpu

	//body, _ = os.ReadFile("two.dtb")
	//_ = memory.LoadBytes(body, DTB)

	var b [4]byte
	go func() {
		e.initialize()
		for {
			e.drawScreen()
		}
	}()

	// Update time goroutines
	go e.UpdateTime()

	for {
		b[0] = memory.ReadByte(cpu.PC)
		b[1] = memory.ReadByte(cpu.PC + 1)
		b[2] = memory.ReadByte(cpu.PC + 2)
		b[3] = memory.ReadByte(cpu.PC + 3)

		// fail on empty instructions
		if b[0] == 0 && b[1] == 0 && b[2] == 0 && b[3] == 0 {
			fmt.Printf("\nFail: Empty instruction. PC: %x\n", cpu.PC)
			cpu.PC += 4
			break
		}
		inst := instructions.DecodeBytes(b)

		//fmt.Println(fmt.Sprintf("Executing Bytes: %x, Opcode: %s, Operation: #%v PC: %x ", instructions.TransformLittleToBig(b), inst.Operation(), inst, cpu.PC))
		//mstatus := instructions.ToMStatusReg(cpu.CSR.GetValue(instructions.MSTATUS, cpu.CurrentMode, &cpu))
		//fmt.Println(fmt.Sprintf("before mstatus: %x", mstatus))

		// f.WriteString(fmt.Sprintf("0x%x\n", cpu.PC))
		_ = cpu.ExecInst(inst)

		//mstatus = instructions.ToMStatusReg(cpu.CSR.GetValue(instructions.MSTATUS, cpu.CurrentMode, &cpu))
		//fmt.Println(fmt.Sprintf("after mstatus: %x", mstatus))

		if e.cpu.Memory.Uart.DataExistsToRead() && cpu.CSR.Registers[instructions.MIE] > 0 {
			// trigger a plic interrupt
			e.cpu.Memory.Plic.TriggerInterrupt(10, cpu)
		}
		_ = cpu.HandleInterrupts(inst.Operation())

		//mstatus = instructions.ToMStatusReg(cpu.CSR.GetValue(instructions.MSTATUS, cpu.CurrentMode, &cpu))
		//fmt.Println(fmt.Sprintf("post interrupt mstatus: %x", mstatus))
		// Handle interrupts / exceptions
	}
}

func (e *Emulator) initialize() {
	// Initialize SDL
	if err := sdl.Init(sdl.INIT_VIDEO); err != nil {
		log.Fatalf("Failed to initialize SDL: %s\n", err)
	}

	// Create a window
	window, err := sdl.CreateWindow("RiscV32", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 320, 200, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Failed to create window: %s\n", err)
	}
	e.window = window

	// Create a renderer
	renderer, err := sdl.CreateRenderer(e.window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Failed to create renderer: %s\n", err)
	}
	e.renderer = renderer

	// Create a texture to manipulate pixels
	texture, err := e.renderer.CreateTexture(sdl.PIXELFORMAT_ARGB8888, sdl.TEXTUREACCESS_STREAMING, SCREEN_WIDTH, SCREEN_HEIGHT)
	if err != nil {
		log.Fatalf("Failed to create texture: %s\n", err)
	}
	e.texture = texture
	e.running = true
}

func (e *Emulator) drawScreen() {
	// Lock the texture to directly modify its pixels
	pixels, pitch, err := e.texture.Lock(nil)
	if err != nil {
		log.Fatalf("Failed to lock texture: %s\n", err)
	}
	e.cpu.Memory.Display.Mutex.Lock()
	Screen := maps.Clone(e.cpu.Memory.Display.Screen)
	e.cpu.Memory.Display.Mutex.Unlock()
	// Set random colors for each pixel
	for y := uint32(0); y < SCREEN_HEIGHT; y++ {
		for x := uint32(0); x < SCREEN_WIDTH; x++ {
			offset := y*uint32(pitch) + x*4 // Calculate the offset in the pixel buffer
			// Generate random ARGB values
			a := byte((Screen[y*SCREEN_WIDTH+x] & 0xFF000000) >> 24)
			r := byte((Screen[y*SCREEN_WIDTH+x] & 0x00FF0000) >> 16)
			g := byte((Screen[y*SCREEN_WIDTH+x] & 0x0000FF00) >> 8)
			b := byte(Screen[y*SCREEN_WIDTH+x] & 0x000000FF)
			pixels[offset] = b
			pixels[offset+1] = g
			pixels[offset+2] = r
			pixels[offset+3] = a
		}
	}

	// Unlock the texture
	e.texture.Unlock()

	// Clear the renderer and copy the texture to it
	e.renderer.Clear()
	e.renderer.Copy(e.texture, nil, nil)
	e.renderer.Present()
	// 60 FPS
	sdl.Delay(30)
	//e.cpu.Memory.Display.Screen[0] = 1 << 16
}
