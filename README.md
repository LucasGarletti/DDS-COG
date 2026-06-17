# DDS-COG
Repositorio del trabajo final de la materia Desarollo de Software, 2026 
Integrantes: Calle, Garletti, Ohanian
Carrera: Lic. Bioinformática


Sistema de gestión y compra de entradas para eventos
1. Descripción
DDS-COG es una aplicación web para la gestión y compra de entradas de eventos.
Permite:
- Registro e inicio de sesión de usuarios.
- Consulta de eventos disponibles.
- Compra de entradas.
- Consulta de entradas adquiridas.
- Cancelación de entradas.
- Transferencia de entradas entre usuarios.
El backend fue desarrollado en Go utilizando Gin y GORM, mientras que el frontend se desarrolla con React y Vite.

2. Tecnologías utilizadas
Backend
- Go
- Gin
- GORM
- JWT
- MySQL
Frontend
- React
- Vite
Testing
- Go Testing
- httptest
Control de versiones
- Git
- GitHub

3. Arquitectura y Base de Datos
https://docs.google.com/document/d/1CN9uuMGZDT8DMT86mzMU75OcnNTKK6DSJbilROskOfQ/edit?usp=sharing 

4. Testing
Cobertura actual: 47.4%
https://docs.google.com/document/d/1XUsHm1oQKEg67baPkc7hhEwIPgR31MwD2DTuUWmXFmE/edit?usp=sharing


5. Estructura del proyecto
DDS-COG
│
├── backend
│   ├── config
│   ├── controllers
│   ├── dao
│   ├── domain
│   ├── middlewares
│   ├── routes
│   ├── services
│   ├── utils
│   └── main.go
│
├── frontend
│   ├── public
│   ├── src
│   │   ├── assets
│   │   ├── pages
│   │   ├── router
│   │   ├── services
│   │   └── utils
│   ├── package.json
│   └── vite.config.js
│
├── README.md
└── .gitignore

6. Herramientas 
El backend fue desarrollado en Go utilizando Gin y GORM, siguiendo una arquitectura en capas:
- **config/**: configuración de la aplicación y conexión a la base de datos.
- **controllers/**: manejo de requests y responses HTTP.
- **dao/**: acceso y persistencia de datos mediante GORM.
- **domain/**: definición de las entidades del sistema.
- **middlewares/**: autenticación JWT y validaciones.
- **routes/**: definición de endpoints de la API.
- **services/**: lógica de negocio.
- **utils/**: utilidades auxiliares, como generación y validación de JWT.
- **main.go**: punto de entrada de la aplicación.

### Frontend
El frontend fue desarrollado con React y Vite.
- **assets/**: imágenes y recursos estáticos.
- **components/**: componentes reutilizables de la interfaz.
- **pages/**: pantallas principales de la aplicación.
- **router/**: configuración de navegación y rutas.
- **services/**: comunicación con la API REST.
- **utils/**: funciones auxiliares utilizadas por la interfaz.