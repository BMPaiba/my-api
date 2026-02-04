package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("🚀 Escuchando en puerto 8888...")

	listener, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		go func(c net.Conn) {
			defer c.Close()
			fmt.Printf("\n✅ CONECTADO: %s\n", c.RemoteAddr())

			buffer := make([]byte, 1024)
			for {
				n, err := c.Read(buffer)
				if err != nil {
					fmt.Printf("❌ Desconectado\n")
					return
				}

				// Log de data cruda
				fmt.Printf("\n📦 %d bytes:\n", n)
				fmt.Printf("HEX: %x\n", buffer[:n])
				fmt.Printf("ASCII: %q\n\n", buffer[:n])
			}
		}(conn)
	}
}
