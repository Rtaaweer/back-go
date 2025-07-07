# 🏥 Sistema de Gestión Hospitalaria

Sistema de gestión de citas médicas desarrollado con Go y Fiber, conectado a Supabase PostgreSQL.

## 📋 Tabla de Contenidos

- [Características](#-características)
- [Estructura del Proyecto](#-estructura-del-proyecto)
- [Tecnologías](#-tecnologías)
- [Requisitos](#-requisitos)
- [Instalación](#️-instalación)
- [Configuración](#-configuración)
- [Uso](#-uso)
- [API Endpoints](#-api-endpoints)
- [Modelos de Datos](#-modelos-de-datos)
- [Contribución](#-contribución)
- [Licencia](#-licencia)

## 🚀 Características

- **Gestión de Usuarios**: Pacientes, médicos, enfermeras y administradores
- **Gestión de Consultorios**: CRUD completo de consultorios médicos
- **Sistema de Consultas**: Programación y gestión de citas médicas
- **Expedientes Médicos**: Historial clínico de pacientes
- **Recetas Médicas**: Gestión de prescripciones
- **Horarios**: Control de disponibilidad médica
- **API REST**: Endpoints bien estructurados
- **CORS**: Habilitado para desarrollo frontend
- **Logging**: Sistema de logs integrado

## 📁 Estructura del Proyecto
menchaca/
├── .env                    # Variables de entorno
├── .vscode/
│   └── settings.json      # Configuración de VS Code
├── README.md              # Documentación del proyecto
├── CHANGELOG.md           # Registro de cambios
├── go.mod                 # Dependencias de Go
├── go.sum                 # Checksums de dependencias
├── main.go                # Punto de entrada de la aplicación
├── hospital-system.exe    # Ejecutable compilado
├── config/
│   └── database.go        # Configuración de base de datos
├── models/                # Modelos de datos
│   ├── usuario.go         # Modelo de usuario
│   ├── consultorio.go     # Modelo de consultorio
│   ├── consulta.go        # Modelo de consulta
│   ├── expediente.go      # Modelo de expediente
│   ├── receta.go          # Modelo de receta
│   └── horario.go         # Modelo de horario
├── handlers/              # Controladores de la API
│   ├── usuarios.go        # Handlers de usuarios
│   ├── consultorios.go    # Handlers de consultorios
│   ├── consultas.go       # Handlers de consultas
│   ├── expedientes.go     # Handlers de expedientes
│   ├── recetas.go         # Handlers de recetas
│   └── horarios.go        # Handlers de horarios
└── routes/
└── routes.go          # Configuración de rutas


## 🛠️ Tecnologías

- **Backend**: Go 1.21+
- **Framework Web**: Fiber v2
- **Base de Datos**: PostgreSQL (Supabase)
- **ORM**: SQL nativo
- **Configuración**: godotenv
- **CORS**: Fiber CORS middleware
- **Logging**: Fiber Logger middleware

## 📋 Requisitos

- Go 1.21 o superior
- PostgreSQL (o cuenta de Supabase)
- Git

## 🛠️ Instalación

### 1. Clonar el repositorio
```bash
git clone <url-del-repositorio>
cd menchaca

## API Endpoints
### Usuarios
- GET /api/v1/usuarios - Obtener todos los usuarios
- GET /api/v1/usuarios/:id - Obtener usuario por ID
- POST /api/v1/usuarios - Crear nuevo usuario
- PUT /api/v1/usuarios/:id - Actualizar usuario
- DELETE /api/v1/usuarios/:id - Eliminar usuario
### Consultorios
- GET /api/v1/consultorios - Obtener todos los consultorios
- GET /api/v1/consultorios/:id - Obtener consultorio por ID
- POST /api/v1/consultorios - Crear nuevo consultorio
- PUT /api/v1/consultorios/:id - Actualizar consultorio
- DELETE /api/v1/consultorios/:id - Eliminar consultorio
### Consultas
- GET /api/v1/consultas - Obtener todas las consultas
- POST /api/v1/consultas - Crear nueva consulta
### Expedientes
- GET /api/v1/expedientes - Obtener todos los expedientes
- GET /api/v1/expedientes/:id - Obtener expediente por ID
- POST /api/v1/expedientes - Crear nuevo expediente
- PUT /api/v1/expedientes/:id - Actualizar expediente
- DELETE /api/v1/expedientes/:id - Eliminar expediente