# Changelog

Todos los cambios notables de este proyecto serán documentados en este archivo.



### Agregado
- README.md estructurado con documentación completa
- CHANGELOG.md para seguimiento de cambios
- Documentación de API endpoints
- Documentación de modelos de datos
- Estructura del proyecto documentada

## [1.0.0] - 2024-01-XX

### Agregado
- Sistema base de gestión hospitalaria
- API REST con Fiber framework
- Conexión a base de datos PostgreSQL/Supabase
- Modelos de datos para el sistema hospitalario:
  - Usuarios (pacientes, médicos, enfermeras, admin)
  - Consultorios
  - Consultas médicas
  - Expedientes médicos
  - Recetas médicas
  - Horarios
- Handlers CRUD para todas las entidades
- Sistema de rutas organizadas
- Middleware de CORS habilitado
- Sistema de logging integrado
- Configuración mediante variables de entorno
- Compilación a ejecutable

### Características Técnicas
- Framework: Go Fiber v2.52.0
- Base de datos: PostgreSQL con driver lib/pq
- Gestión de configuración: godotenv
- Arquitectura: REST API
- Patrón: MVC (Model-View-Controller)

### Endpoints Implementados
- `/api/v1/usuarios/*` - Gestión completa de usuarios
- `/api/v1/consultorios/*` - Gestión completa de consultorios
- `/api/v1/consultas/*` - Gestión de consultas médicas
- `/api/v1/expedientes/*` - Gestión completa de expedientes
- `/api/v1/recetas/*` - Gestión de recetas médicas
- `/api/v1/horarios/*` - Gestión de horarios médicos

### Seguridad
- Variables de entorno para credenciales sensibles
- CORS configurado para desarrollo
- Validación de datos en endpoints

---

## Formato de Versionado

Este proyecto utiliza [Semantic Versioning](https://semver.org/):
- **MAJOR**: Cambios incompatibles en la API
- **MINOR**: Funcionalidad agregada de manera compatible
- **PATCH**: Correcciones de bugs compatibles

## Tipos de Cambios

- **Agregado**: para nuevas características
- **Cambiado**: para cambios en funcionalidad existente
- **Obsoleto**: para características que serán removidas
- **Removido**: para características removidas
- **Corregido**: para corrección de bugs
- **Seguridad**: para vulnerabilidades

## Convenciones de Commits

Para mantener un historial limpio, se recomienda usar:
- `feat:` para nuevas características
- `fix:` para corrección de bugs
- `docs:` para cambios en documentación
- `style:` para cambios de formato
- `refactor:` para refactorización de código
- `test:` para agregar o modificar tests
- `chore:` para tareas de mantenimiento

