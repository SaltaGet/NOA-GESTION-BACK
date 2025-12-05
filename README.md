# NOA GESTION - Backend

Este es el repositorio del backend para la aplicación NOA Gestión, desarrollado en Go con el framework Fiber.

## Estructura del Proyecto

El proyecto sigue una estructura organizada para separar responsabilidades, facilitar el mantenimiento y la escalabilidad.

```
NOA-GESTION-BACK/
├── cmd/
│   └── api/
│       ├── main.go               # Punto de entrada principal de la API.
│       ├── controllers/          # Controladores HTTP que manejan la lógica de las rutas.
│       ├── docs/                 # Archivos generados por Swag para la documentación.
│       ├── initial/              # Paquetes para inicializar configuraciones (ej. email).
│       ├── jobs/                 # Tareas programadas o de un solo uso (ej. migraciones).
│       ├── logging/              # Configuración y helpers para el logging.
│       ├── middleware/           # Middlewares para las rutas de Fiber.
│       └── routes/               # Definición y configuración de las rutas de la API.
│
├── internal/
│   ├── cache/                  # Lógica de interacción con la caché (Redis).
│   ├── database/               # Gestión de la conexión a la base de datos y migraciones.
│   ├── dependencies/           # Contenedor de inyección de dependencias.
│   ├── ports/                  # Interfaces que definen los contratos para los servicios (arquitectura hexagonal).
│   ├── schemas/                # Estructuras de datos (DTOs) para requests, responses y validaciones.
│   ├── test/                   # Pruebas unitarias y de integración.
│   └── validators/             # Funciones de validación reutilizables.
│
├── .env                        # Archivo de variables de entorno (local, no versionado).
├── go.mod                      # Archivo de definición del módulo de Go.
├── go.sum                      # Checksums de las dependencias.
├── grafana.json                # Configuración del dashboard de Grafana para monitoreo.
└── README.md                   # Este archivo.
```

### Descripción de Directorios

#### `cmd/api/`
Contiene el código ejecutable de la aplicación.
- **`main.go`**: Inicializa la aplicación Fiber, configura la base de datos, la caché, el logger, los middlewares y las rutas, y arranca el servidor. También maneja el apagado gracefully.
- **`controllers/`**: Cada archivo corresponde a una entidad del negocio (ej. `auth.go`, `product.go`). Reciben las peticiones HTTP, las procesan (parseo, validación básica) y llaman a los servicios correspondientes.
- **`docs/`**: Contiene los archivos `swagger.json`, `swagger.yaml` y `docs.go` generados por `swag` para la documentación de la API.
- **`middleware/`**: Contiene los middlewares de Fiber que se aplican a las rutas, como el de inyección de dependencias (`InjectionDepends`), logging (`LoggingMiddleware`), limitación de peticiones (`RateLimitMiddleware`), etc.
- **`routes/`**: El archivo `routes.go` define todas las rutas de la API, las agrupa y les asigna sus controladores y middlewares específicos.

#### `internal/`
Contiene la lógica de negocio principal de la aplicación. Al estar en `internal/`, Go previene que otros proyectos importen estos paquetes.
- **`cache/`**: Abstrae la lógica de conexión y operaciones con Redis.
- **`database/`**: Gestiona la conexión con la base de datos principal y las bases de datos de los tenants, incluyendo un sistema de caché de conexiones y un "janitor" para limpiar conexiones inactivas.
- **`dependencies/`**: Define el struct `Application` que actúa como un contenedor para las dependencias (servicios, base de datos), facilitando la inyección de dependencias.
- **`ports/`**: Define las interfaces (puertos en terminología de arquitectura hexagonal) que los servicios deben implementar. Esto desacopla los controladores de la implementación concreta de los servicios.
- **`schemas/`**: Contiene todas las estructuras de datos usadas para la transferencia de información: modelos para las peticiones (`body`), respuestas (`response`), y objetos de transferencia de datos (DTOs) entre capas. También incluye los métodos de validación para estas estructuras.
- **`test/`**: Contiene los tests de la aplicación. El ejemplo `auth_controller_test.go` muestra cómo se realizan tests unitarios de los controladores mockeando los servicios.
- **`validators/`**: Proporciona funciones de ayuda para validaciones comunes y reutilizables, como la validación de un ID.

### Archivos Raíz

- **`.env`**: Almacena variables de entorno para la configuración local (puerto, credenciales de BD, etc.).
- **`go.mod` / `go.sum`**: Gestionan las dependencias del proyecto.
- **`grafana.json`**: Archivo de configuración para importar un dashboard predefinido en Grafana, útil para el monitoreo de la aplicación a través de Prometheus y Loki.