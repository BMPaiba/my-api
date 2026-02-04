package main

import (
	"fmt"
	"net"
	"time"
)

// Servidor de debug para ver datos crudos del dispositivo Anviz EP300 Pro
func startAnvizDebugServer() {
	// Escuchamos en el puerto 8888 que configuraste en el EP300 Pro
	listener, err := net.Listen("tcp", "0.0.0.0:8888")
	if err != nil {
		fmt.Printf("❌ Error al abrir puerto 8888: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Println("🚀 Servidor Anviz DEBUG a la escucha en el puerto 8888...")
	fmt.Println("⏳ Esperando que el ícono de la nube en el dispositivo se active...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("❌ Error de conexión: %v\n", err)
			continue
		}

		// Manejamos la conexión en una goroutine
		go func(c net.Conn) {
			defer c.Close()
			
			// IP del dispositivo que se conectó
			remoteAddr := c.RemoteAddr().String()
			fmt.Printf("\n✅ ¡DISPOSITIVO CONECTADO! Desde: %s\n", remoteAddr)
			fmt.Printf("⏰ Hora: %s\n", time.Now().Format("2006-01-02 15:04:05"))

			// Buffer para recibir los datos
			buffer := make([]byte, 1024)
			
			// Establecemos un timeout para no quedar bloqueados si no envía nada
			c.SetReadDeadline(time.Now().Add(30 * time.Second))

			for {
				n, err := c.Read(buffer)
				if err != nil {
					fmt.Printf("ℹ️ Conexión finalizada con %s: %v\n", remoteAddr, err)
					return
				}

				// ESTO ES LO QUE BUSCAMOS: Ver los bytes en la terminal
				fmt.Printf("\n📥 [%s] DATA RECIBIDA (%d bytes):\n", time.Now().Format("15:04:05"), n)
				fmt.Printf("HEX: %x\n", buffer[:n])
				fmt.Printf("STR: %q\n", string(buffer[:n]))
				fmt.Println("--------------------------------------------------")
				
				// Resetear el timeout después de recibir datos
				c.SetReadDeadline(time.Now().Add(30 * time.Second))
			}
		}(conn)
	}
}

func main() {
	fmt.Println("================================================")
	fmt.Println("   SERVIDOR DEBUG - ANVIZ EP300 Pro")
	fmt.Println("================================================")
	startAnvizDebugServer()
}
