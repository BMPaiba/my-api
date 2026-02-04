package tcp

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	addr   string
	logger *zap.SugaredLogger
}

func NewServer(addr string, logger *zap.SugaredLogger) *Server {
	return &Server{
		addr:   addr,
		logger: logger,
	}
}

func (s *Server) Start() error {
	// Forzar IPv4 resolviendo la dirección
	tcpAddr, err := net.ResolveTCPAddr("tcp4", s.addr)
	if err != nil {
		return fmt.Errorf("error resolviendo dirección TCP: %w", err)
	}

	listener, err := net.ListenTCP("tcp4", tcpAddr)
	if err != nil {
		return fmt.Errorf("error creando listener TCP: %w", err)
	}
	defer listener.Close()

	s.logger.Infof("Servidor TCP escuchando en %s para dispositivo Aanviz EP300Pro", s.addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			s.logger.Errorf("Error aceptando conexión: %v", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	s.logger.Infof("🔌 Nueva conexión desde dispositivo Aanviz: %s", clientAddr)

	for {
		// Primero lee el header (10 bytes)
		header := make([]byte, 10)
		_, err := io.ReadFull(conn, header)
		if err != nil {
			if err != io.EOF && err != io.ErrUnexpectedEOF {
				s.logger.Errorf("❌ Error leyendo header de %s: %v", clientAddr, err)
			} else {
				s.logger.Infof("🔌 Dispositivo %s cerró la conexión", clientAddr)
			}
			return
		}

		// Verifica start marker
		if header[0] != 0xA5 {
			// Si recibimos bytes inválidos (especialmente 0x00), probablemente
			// la conexión está terminando o hay basura. Cerramos la conexión.
			s.logger.Debugf("Start marker inválido de %s: 0x%02X, cerrando conexión", clientAddr, header[0])
			return
		}

		// Lee la longitud del payload (bytes 8-9)
		payloadLen := binary.LittleEndian.Uint16(header[8:10])

		var fullMessage []byte

		if payloadLen == 0 {
			// Mensaje de solo 10 bytes (sin payload), es completo
			fullMessage = header
		} else {
			// Lee el payload + checksum (2 bytes)
			remaining := make([]byte, int(payloadLen)+2)
			_, err = io.ReadFull(conn, remaining)
			if err != nil {
				s.logger.Errorf("❌ Error leyendo payload de %s: %v", clientAddr, err)
				return
			}

			// Combina header + remaining para el mensaje completo
			fullMessage = append(header, remaining...)
		}

		// Procesar el mensaje completo
		s.processAanvizData(fullMessage, conn, clientAddr)
	}
}

func (s *Server) processAanvizData(data []byte, conn net.Conn, clientAddr string) {
	// Parsear el mensaje
	msg, err := ParseAanvizMessage(data)
	if err != nil {
		s.logger.Warnf("❌ Error parseando mensaje de %s: %v (hex: %s)", clientAddr, err, hex.EncodeToString(data))
		return
	}

	// Log del mensaje recibido
	s.logger.Infof("📨 [%s] %s | Payload: %d bytes", clientAddr, msg.Command.String(), msg.Length)

	// Procesar según el tipo de mensaje
	switch msg.Command {
	case MsgHeartbeat:
		s.logger.Infof("💓 Heartbeat del dispositivo")

	case MsgVerifyRecord:
		s.handleVerifyRecord(msg.Payload, clientAddr)

	default:
		s.logger.Infof("📦 Mensaje tipo 0x%02X | Payload (hex): %s", byte(msg.Command), hex.EncodeToString(msg.Payload))
	}

	// Enviar ACK
	ack := BuildACK(msg.Command)
	if _, err := conn.Write(ack); err != nil {
		s.logger.Errorf("❌ Error enviando ACK: %v", err)
	} else {
		s.logger.Debugf("✓ ACK enviado para comando 0x%02X", byte(msg.Command))
	}
}

func (s *Server) handleVerifyRecord(payload []byte, clientAddr string) {
	record, err := ParseVerifyRecord(payload)
	if err != nil {
		s.logger.Warnf("⚠️  Error parseando registro de verificación: %v", err)
		s.logger.Infof("   📊 Payload (hex): %s", hex.EncodeToString(payload))
		return
	}

	// Determinar modo de verificación según BackupID
	var mode string
	backupLow := record.BackupID & 0x0F         // Bits 0-3
	backupHigh := (record.BackupID >> 4) & 0x0F // Bits 4-7

	if backupHigh >= 1 && backupHigh <= 10 {
		mode = fmt.Sprintf("Huella #%d", backupHigh)
	} else if (record.BackupID & 0x08) != 0 { // bit 3
		mode = "Tarjeta"
	} else if (record.BackupID & 0x04) != 0 { // bit 2
		mode = "Contraseña"
	} else if backupLow > 0 {
		mode = fmt.Sprintf("Otro (0x%02X)", record.BackupID)
	} else {
		mode = "Desconocido"
	}

	// Determinar resultado y tipo de acceso desde RecordType
	doorOpened := (record.RecordType & 0x80) != 0 // bit 7
	attendanceStatus := record.RecordType & 0x0F  // bits 0-3

	var statusStr string
	switch attendanceStatus {
	case 0:
		statusStr = "Check-in"
	case 1:
		statusStr = "Check-out"
	case 2:
		statusStr = "Break-out"
	case 3:
		statusStr = "Break-in"
	case 4:
		statusStr = "Overtime-in"
	case 5:
		statusStr = "Overtime-out"
	default:
		statusStr = fmt.Sprintf("Estado %d", attendanceStatus)
	}

	var emoji string
	if doorOpened {
		emoji = "✅"
	} else {
		emoji = "❌"
	}

	// Convertir timestamp (segundos desde 2000-01-02) a fecha legible
	baseTime := time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC)
	recordTime := baseTime.Add(time.Duration(record.Timestamp) * time.Second)

	s.logger.Infof("%s Usuario: %s | Modo: %s | Estado: %s | Puerta: %v | Hora: %s",
		emoji, record.EmployeeID, mode, statusStr, doorOpened, recordTime.Format("2006-01-02 15:04:05"))

	// TODO: Aquí puedes guardar el registro en la base de datos
}
