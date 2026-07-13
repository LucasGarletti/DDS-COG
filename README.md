Trabajo Final Integrador – Desarrollo de Software 2026

Integrantes

Calle
Garletti
Ohanian

Carrera: Licenciatura en Bioinformática
1. Descripción del proyecto

DDS-COG es una aplicación web para la gestión, administración y compra de entradas para eventos.

El sistema permite a los usuarios:

Registrarse e iniciar sesión.
Consultar eventos disponibles.
Comprar entradas.
Visualizar entradas adquiridas.
Cancelar entradas.
Transferir entradas a otros usuarios.
Crear un itinerario personal para festivales (Bonus Track).

Además, cuenta con funcionalidades administrativas para la gestión completa de eventos y festivales.

2. Funcionalidades principales
Usuarios
Registro de cuenta.
Inicio de sesión mediante JWT.
Consulta de eventos.
Compra de entradas.
Cancelación de entradas.
Transferencia de entradas.
Visualización de entradas adquiridas.
Festivales (Bonus Track)

Para los eventos marcados como festivales:

El administrador puede crear una grilla oficial.
Se pueden cargar artistas, escenarios y horarios.
Los usuarios que posean una entrada activa pueden crear su itinerario personal.
Se detectan conflictos horarios entre actividades.
Se pueden agregar actividades personalizadas.
Administrador
Crear eventos.
Modificar eventos.
Cancelar eventos.
Gestionar la grilla oficial de festivales.
Consultar reportes y estadísticas.
3. Tecnologías utilizadas
Backend
Go
Gin
GORM
JWT
MySQL
Frontend
React
Vite
Testing
Go Testing
httptest
Contenedores
Docker
Docker Compose
Control de versiones
Git
GitHub
4. Arquitectura del sistema
Frontend (React + Vite)
           │
           ▼
Backend (Go + Gin)
           │
           ▼
MySQL

El frontend consume una API REST desarrollada en Go.

La persistencia de datos se realiza mediante MySQL utilizando GORM como ORM.

5. Modelo de datos

Principales entidades:

Users: Representa los usuarios del sistema.

Events: Representa eventos y festivales.

Tickets: Entradas adquiridas por los usuarios.

Festival Schedules: Grilla oficial de festivales administrada por el rol administrador.

User Itineraries: Itinerarios personales de los usuarios para festivales.

Itinerary Items: Elementos individuales que forman parte de un itinerario.

6. Arquitectura del backend

El backend fue desarrollado siguiendo una arquitectura en capas.

config/

Configuración de la aplicación y conexión a la base de datos.

controllers/

Recepción y respuesta de solicitudes HTTP.

dao/

Acceso y persistencia de datos mediante GORM.

domain/

Definición de entidades del sistema.

middlewares/

Autenticación JWT y autorización por roles.

routes/

Definición de endpoints.

services/

Implementación de la lógica de negocio.

utils/

Funciones auxiliares.

main.go

Punto de entrada de la aplicación.

7. Arquitectura del frontend
assets/

Recursos estáticos.

components/

Componentes reutilizables.

pages/

Pantallas principales de la aplicación.

router/

Configuración de rutas.

services/

Comunicación con la API REST.

utils/

Funciones auxiliares.

8. Testing

Cobertura actual: 47.4%

Documentación de testing:

https://docs.google.com/document/d/1XUsHm1oQKEg67baPkc7hhEwIPgR31MwD2DTuUWmXFmE/edit?usp=sharing

Para ejecutar los tests del backend:

go test ./...
9. Diagramas y documentación

Documentación de arquitectura y base de datos:

https://docs.google.com/document/d/1CN9uuMGZDT8DMT86mzMU75OcnNTKK6DSJbilROskOfQ/edit?usp=sharing

10. Requisitos

Para ejecutar el proyecto mediante Docker:

Docker Desktop instalado.
Docker Compose habilitado.

Verificar instalación:

docker --version
docker compose version

11. Ejecución con Docker

El proyecto puede ejecutarse completamente mediante Docker Compose.

Desde la raíz del proyecto: docker compose up --build

Durante la primera ejecución Docker:

Descargará las imágenes necesarias.
Construirá las imágenes del frontend y backend.
Creará la base de datos MySQL.
Ejecutará las migraciones automáticas mediante GORM.
12. Servicios disponibles

Una vez iniciado el sistema:

Frontend
http://localhost:5173
Backend
http://localhost:8080
Base de datos MySQL
Host: localhost
Puerto: 3306
Base: dds_cog
13. Datos iniciales de prueba

Al levantar el sistema con Docker, el backend ejecuta un seed inicial porque
docker-compose.yml define RUN_SEED=true para el servicio backend.

Credenciales de demostracion:

Administrador:
Email: admin@tickgo.com
Password: Admin123!

Cliente:
Email: client@tickgo.com
Password: Client123!

Las contrasenas se guardan hasheadas con la misma funcion usada por el registro
normal. El seed es idempotente: busca usuarios por email, eventos por titulo y
artistas por event_id + artist + start_time, por lo que puede ejecutarse varias
veces sin duplicar registros.

Eventos cargados:

Festival Cosquín Rock
Los Pumas en el Estadio Mario Alberto Kempes
Las Pastillas del Abuelo

Grilla inicial de Cosquín Rock:

Airbag
Dillom
Guasones

14. Comandos útiles de Docker

Levantar el sistema:

docker compose up --build

Detener contenedores:

docker compose down

Detener contenedores y eliminar volumen de MySQL:

docker compose down -v

Ver logs:

docker compose logs -f

Ver estado de contenedores:

docker compose ps

Reconstruir imágenes:

docker compose build --no-cache
docker compose up
15. Persistencia de datos

El sistema utiliza un volumen Docker para almacenar la base de datos.

docker compose down

Conserva los datos.

docker compose down -v

Elimina completamente la base de datos creada dentro de Docker.

La base MySQL utilizada por Docker es independiente de cualquier instalación local de MySQL.

16. Ejecución sin Docker
Backend
cd backend
go run main.go
Frontend
cd frontend
npm install
npm run dev
17. API y autenticación

La autenticación se realiza mediante JWT.

Roles disponibles:

client
admin

Las rutas administrativas requieren autenticación y autorización mediante middleware de roles.

18. Estado actual del proyecto

Funcionalidades implementadas:

Autenticación JWT.
Registro e inicio de sesión.
Gestión de eventos.
Compra de entradas.
Cancelación de entradas.
Transferencia de entradas.
Reportes administrativos.
Gestión de festivales.
Bonus Track de itinerarios personales.
Persistencia en MySQL.
Dockerización completa mediante Docker Compose.
19. Estructura general del proyecto
DDS-COG
│
├── backend
│
├── frontend
│
├── docker-compose.yml
│
├── README.md
│
└── .gitignore
