# Configuración de Red para Aanviz EP300Pro

## Configuración del Servidor

El servidor TCP ahora escucha en **0.0.0.0:8888**, lo que significa:
- ✅ Acepta conexiones desde cualquier interfaz de red
- ✅ Dispositivos en la misma red local pueden conectarse
- ✅ No solo escucha en localhost

## Configuración del Firewall de Windows

Para permitir conexiones entrantes en el puerto 8888, ejecuta en PowerShell como Administrador:

```powershell
New-NetFirewallRule -DisplayName "Aanviz TCP Server" -Direction Inbound -LocalPort 8888 -Protocol TCP -Action Allow
```

O manualmente:
1. Abre "Windows Defender Firewall con seguridad avanzada"
2. Clic en "Reglas de entrada" → "Nueva regla"
3. Tipo: Puerto
4. Protocolo: TCP, Puerto: 8888
5. Acción: Permitir la conexión
6. Aplicar a: Dominio, Privado, Público (según tu red)
7. Nombre: "Aanviz TCP Server"

## Configuración del Dispositivo Aanviz EP300Pro

Configura en el dispositivo:
- **IP del servidor**: `192.168.0.25` (tu IP local)
- **Puerto**: `8888`
- **Protocolo**: TCP

## Verificar Configuración

### 1. Verificar que el servidor esté escuchando:
```bash
# En WSL
netstat -tlnp | grep 8888
```

### 2. Verificar desde otro dispositivo en la red:
```bash
# Desde otro equipo en la red
telnet 192.168.0.25 8888
# o
nc -zv 192.168.0.25 8888
```

### 3. Verificar firewall de Windows:
```powershell
# En PowerShell
Get-NetFirewallRule -DisplayName "*Aanviz*"
```

## Variables de Entorno (.env)

Agrega o verifica en tu archivo `.env`:
```env
TCP_PORT=0.0.0.0:8888
```

## Troubleshooting

### Si el dispositivo no puede conectarse:

1. **Firewall bloqueando**: Revisa el firewall de Windows
2. **IP incorrecta**: Verifica con `ipconfig.exe` que la IP sea la correcta
3. **Puerto ocupado**: Verifica que no haya otro servicio en el puerto 8888
4. **Red diferente**: Asegúrate de que el dispositivo Aanviz esté en la misma red (192.168.0.x)

### Verificar que WSL pueda recibir desde la red local:

```bash
# En WSL
sudo apt-get install net-tools
netstat -tlnp | grep 8888
```

## Información de la Red

IP del servidor: `192.168.0.25`  
Puerto TCP: `8888`  
Escucha en: `0.0.0.0:8888` (todas las interfaces)
