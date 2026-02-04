package tcp

import (
	"encoding/binary"
	"fmt"
)

// AanvizMessageType representa los tipos de mensajes del protocolo Aanviz
type AanvizMessageType byte

const (
	// Comandos comunes del Aanviz EP300Pro
	MsgHeartbeat      AanvizMessageType = 0x7F
	MsgVerifyRecord   AanvizMessageType = 0xDF
	MsgUserInfo       AanvizMessageType = 0x03
	MsgTimeSync       AanvizMessageType = 0x04
	MsgDeviceInfo     AanvizMessageType = 0x05
	MsgFingerTemplate AanvizMessageType = 0x06
	MsgAccessControl  AanvizMessageType = 0x07
)

// String retorna una descripción del tipo de mensaje
func (t AanvizMessageType) String() string {
	switch t {
	case MsgHeartbeat:
		return "Heartbeat"
	case MsgVerifyRecord:
		return "Registro de Verificación"
	case MsgUserInfo:
		return "Info de Usuario"
	case MsgTimeSync:
		return "Sincronización de Tiempo"
	case MsgDeviceInfo:
		return "Info del Dispositivo"
	case MsgFingerTemplate:
		return "Plantilla de Huella"
	case MsgAccessControl:
		return "Control de Acceso"
	default:
		return fmt.Sprintf("Desconocido (0x%02X)", byte(t))
	}
}

// AanvizMessage representa la estructura básica de un mensaje Aanviz
type AanvizMessage struct {
	StartMarker byte
	Command     AanvizMessageType
	Length      uint16
	Payload     []byte
	Checksum    byte
}

// ParseAanvizMessage parsea los datos recibidos del dispositivo
func ParseAanvizMessage(data []byte) (*AanvizMessage, error) {
	if len(data) < 10 {
		return nil, fmt.Errorf("mensaje demasiado corto: %d bytes", len(data))
	}

	// Formato Aanviz EP300Pro:
	// Byte 0: Start marker (0xA5)
	// Bytes 1-4: Device ID
	// Byte 5: Command
	// Bytes 6-7: Reserved/Flags (0x0000)
	//
	// Si el mensaje tiene exactamente 10 bytes (sin payload):
	//   Bytes 8-9: Checksum
	//
	// Si el mensaje tiene más de 10 bytes (con payload):
	//   Bytes 8-9: Length (little-endian) - longitud del payload
	//   Bytes 10-n: Payload
	//   Últimos 2 bytes: Checksum

	msg := &AanvizMessage{
		StartMarker: data[0],
		Command:     AanvizMessageType(data[5]),
	}

	// Si el mensaje tiene exactamente 10 bytes, no hay payload
	if len(data) == 10 {
		msg.Length = 0
		msg.Payload = nil
		msg.Checksum = data[8]
		return msg, nil
	}

	// Leer longitud del payload (bytes 8-9 en little-endian)
	msg.Length = binary.LittleEndian.Uint16(data[8:10])

	// Calcular posiciones
	payloadStart := 10
	payloadEnd := payloadStart + int(msg.Length)

	// Verificar que tengamos suficientes bytes
	expectedTotal := payloadEnd + 2 // payload + 2 bytes de checksum
	if len(data) < expectedTotal {
		return nil, fmt.Errorf("longitud de mensaje inválida: length indica %d bytes de payload, pero solo hay %d bytes totales (esperado %d)",
			msg.Length, len(data), expectedTotal)
	}

	// Extraer payload
	if msg.Length > 0 {
		msg.Payload = data[payloadStart:payloadEnd]
	}

	// Extraer checksum (últimos 2 bytes)
	msg.Checksum = data[len(data)-2]

	return msg, nil
}

// VerifyChecksum verifica el checksum del mensaje
func (m *AanvizMessage) VerifyChecksum(data []byte) bool {
	var sum byte
	for _, b := range data {
		sum ^= b
	}
	return sum == m.Checksum
}

// BuildACK construye un ACK específico para el comando recibido
func BuildACK(cmd AanvizMessageType) []byte {
	// Formato ACK Aanviz:
	// Byte 0: 0xA5 (start marker)
	// Bytes 1-4: 0x00 0x00 0x00 0x01 (device ID respuesta)
	// Byte 5: Comando recibido
	// Bytes 6-7: 0x00 0x00 (flags)
	// Bytes 8-9: 0x00 0x00 (sin payload)
	return []byte{0xA5, 0x00, 0x00, 0x00, 0x01, byte(cmd), 0x00, 0x00, 0x00, 0x00}
}

// VerifyRecord representa un registro de verificación del dispositivo
type VerifyRecord struct {
	EmployeeID string  // User ID (de 5 bytes)
	Timestamp  uint32  // Segundos desde 2000-01-02
	BackupID   byte    // Método de verificación
	RecordType byte    // Tipo de registro y resultado
	WorkType   [3]byte // Código de trabajo
}

// ParseVerifyRecord parsea un registro de verificación
func ParseVerifyRecord(payload []byte) (*VerifyRecord, error) {
	if len(payload) < 14 {
		return nil, fmt.Errorf("payload de registro de verificación inválido (recibido %d bytes, esperado al menos 14)", len(payload))
	}

	// Formato según manual (14 bytes):
	// Bytes 0-4: Employee ID (5 bytes)
	// Bytes 5-8: Timestamp (little-endian, segundos desde 2000-01-02)
	// Byte 9: Backup ID (método de verificación)
	// Byte 10: Record Type
	// Bytes 11-13: Work Type (3 bytes)

	// Extraer Employee ID (5 bytes) y convertir a string
	// El ID puede ser numérico, así que lo convertimos de bytes
	employeeID := extractEmployeeID(payload[0:5])

	record := &VerifyRecord{
		EmployeeID: employeeID,
		Timestamp:  binary.LittleEndian.Uint32(payload[5:9]),
		BackupID:   payload[9],
		RecordType: payload[10],
	}

	copy(record.WorkType[:], payload[11:14])

	return record, nil
}

// extractEmployeeID extrae el ID del empleado de 5 bytes
func extractEmployeeID(data []byte) string {
	// El ID puede estar almacenado como número o string
	// Intentamos convertirlo a número primero
	id := binary.LittleEndian.Uint32(data[0:4])
	if id > 0 && id < 100000000 { // Rango razonable para un ID numérico
		return fmt.Sprintf("%d", id)
	}

	// Si no, tratarlo como string ASCII
	var result []byte
	for _, b := range data {
		if b > 0 && b < 127 { // ASCII válido
			result = append(result, b)
		}
	}

	if len(result) > 0 {
		return string(result)
	}

	// Último recurso: convertir a hex
	return fmt.Sprintf("%X", data)
}
