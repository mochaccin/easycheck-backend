# Easycheck

Si alguna wea falla taggenme en el ds

## Prerequisitos

- Go
- Configurar el .env

## Como instalar go

- Descargar la version que corresponda a su sistema operativo desde aca [Go Downloads](https://go.dev/dl/)
- Lo instalan como cualquier otro programa

Listo ya tienen go instalado.

Ahora clonan el repo y tienen 2 opciones, buildearlo manual de la siguiente forma

Descargar dependencias

```go
go mod tidy
```
Correr el proyecto

```go
go run cmd/api/main.go
```

La otra opcion es utilizar live reloading, cosa que si quieren cambiar algo no tienen que reiniciar como tal.

Para eso hay que usar go/air

Instalan chocolatey que nos permite instalar herramientas de linux en windows, en este caso el comando make para usar el makefile que usa go air.

```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://community.chocolatey.org/install.ps1'))
```

Despues instalan go air para el hot reloading

```powershell
go install github.com/air-verse/air@latest
```

Finalmente para ejecutar el proyecto abren una terminal en el vscode y escriben

```powershell
air
```
