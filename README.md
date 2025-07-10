# 🏥 Sistema de Gestión Hospitalaria

Sistema de gestión de citas médicas desarrollado con Go y Fiber (backend) y Angular (frontend), conectado a Supabase PostgreSQL.

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
- **Autenticación MFA**: Sistema de autenticación de dos factores con TOTP
- **Gestión de Consultorios**: CRUD completo de consultorios médicos
- **Sistema de Consultas**: Programación y gestión de citas médicas
- **Expedientes Médicos**: Historial clínico de pacientes
- **Recetas Médicas**: Gestión de prescripciones
- **Horarios**: Control de disponibilidad médica
- **Frontend Angular**: Interfaz de usuario moderna y responsiva
- **API REST**: Endpoints bien estructurados
- **CORS**: Habilitado para desarrollo frontend
- **Logging**: Sistema de logs integrado

## 📁 Estructura del Proyecto
Backend-Base-de-datos-main/
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
│   ├── auth.go            # Handlers de autenticación y MFA
│   ├── mfa.go             # Handlers específicos de MFA
│   ├── usuarios.go        # Handlers de usuarios
│   ├── consultorios.go    # Handlers de consultorios
│   ├── consultas.go       # Handlers de consultas
│   ├── expedientes.go     # Handlers de expedientes
│   ├── recetas.go         # Handlers de recetas
│   └── horarios.go        # Handlers de horarios
├── middleware/            # Middlewares
│   ├── auth.go            # Middleware de autenticación
│   ├── ratelimit.go       # Middleware de rate limiting
│   └── response_validator.go # Validador de respuestas
├── utils/                 # Utilidades
│   ├── jwt.go             # Utilidades JWT
│   ├── mfa.go             # Utilidades MFA
│   ├── password.go        # Utilidades de contraseñas
│   └── response_codes.go  # Códigos de respuesta
├── routes/
│   └── routes.go          # Configuración de rutas
├── schemas/
│   └── response_schemas.go # Esquemas de respuesta
└── hospital-frontend/     # Frontend Angular
├── src/
│   ├── app/
│   │   ├── pages/
│   │   │   └── auth/
│   │   │       ├── login/
│   │   │       └── register/
│   │   └── services/
│   │       └── auth.service.ts
│   └── styles.css
├── angular.json
├── package.json
└── tsconfig.json


## 🛠️ Tecnologías

### Backend
- **Go** 1.21+
- **Fiber** v2 (Framework Web)
- **PostgreSQL** (Supabase)
- **JWT** para autenticación
- **TOTP** para MFA
- **bcrypt** para hash de contraseñas

### Frontend
- **Angular** 17+
- **PrimeNG** (Componentes UI)
- **TypeScript**
- **RxJS**
- **Angular Reactive Forms**

## 📋 Requisitos

- **Go** 1.21 o superior
- **Node.js** 18+ y npm
- **Angular CLI** 17+
- **PostgreSQL** (o cuenta de Supabase)
- **Git**

## 🛠️ Instalación

### 1. Clonar el repositorio
```bash
git clone <url-del-repositorio>
cd Backend-Base-de-datos-main



## 📋 Configuración 
2. Configurar el Backend
Instalar dependencias de Go
bash
Run
go mod download
Configurar variables de entorno
Crea un archivo .env en la raíz del proyecto:

env

DB_HOST=tu-host-supabaseDB_PORT=5432DB_USER=tu-usuarioDB_PASSWORD=tu-contraseñaDB_NAME=tu-base-de-datosJWT_SECRET=tu-clave-secreta-jwtPORT=3000
3. Configurar el Frontend
Navegar al directorio del frontend
bash
Run
cd hospital-frontend
Instalar dependencias de Node.js
bash
Run
npm install
Instalar Angular CLI (si no está instalado)
bash
Run
npm install -g @angular/cli
🚀 Uso
Ejecutar el Backend
Desde la raíz del proyecto:

bash
Run
go run main.go
El servidor backend estará disponible en: http://localhost:3000

Ejecutar el Frontend
Desde el directorio hospital-frontend:

bash
Run
cd hospital-frontendng serve
O desde la raíz del proyecto:

bash
Run
cd hospital-frontend && ng serve
El frontend estará disponible en: http://localhost:4200

Desarrollo Completo
Para ejecutar ambos servicios simultáneamente:

Terminal 1 (Backend):

bash
Run
go run main.go
Terminal 2 (Frontend):

bash
Run
cd hospital-frontendng serve
🔐 Funcionalidades de Autenticación
Registro de Usuario
Completa el formulario de registro
Al registrarte exitosamente, aparecerá un diálogo con:
Código QR para configurar MFA
Clave secreta manual
Escanea el QR con una app como Google Authenticator
Continúa al login
Inicio de Sesión
Ingresa tu email y contraseña
Si tienes MFA habilitado, ingresa el código de 6 dígitos
Accede al dashboard
📡 API Endpoints
Autenticación
POST /api/v1/auth/register - Registro de usuario con MFA
POST /api/v1/auth/login - Inicio de sesión con soporte MFA
POST /api/v1/auth/enable-mfa - Habilitar MFA
POST /api/v1/auth/verify-mfa - Verificar código MFA
Usuarios
GET /api/v1/usuarios - Obtener todos los usuarios
GET /api/v1/usuarios/:id - Obtener usuario por ID
POST /api/v1/usuarios - Crear nuevo usuario
PUT /api/v1/usuarios/:id - Actualizar usuario
DELETE /api/v1/usuarios/:id - Eliminar usuario
Consultorios
GET /api/v1/consultorios - Obtener todos los consultorios
GET /api/v1/consultorios/:id - Obtener consultorio por ID
POST /api/v1/consultorios - Crear nuevo consultorio
PUT /api/v1/consultorios/:id - Actualizar consultorio
DELETE /api/v1/consultorios/:id - Eliminar consultorio
Consultas
GET /api/v1/consultas - Obtener todas las consultas
POST /api/v1/consultas - Crear nueva consulta
Expedientes
GET /api/v1/expedientes - Obtener todos los expedientes
GET /api/v1/expedientes/:id - Obtener expediente por ID
POST /api/v1/expedientes - Crear nuevo expediente
PUT /api/v1/expedientes/:id - Actualizar expediente
DELETE /api/v1/expedientes/:id - Eliminar expediente