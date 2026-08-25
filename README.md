
# Orienmod

Es un gestor educativo intermedio diseñado para liberar a los maestros de las complicaciones técnicas al momento de evaluar. Funciona como un puente dinámico y optimizado entre Google Classroom y el docente, transformando los datos complejos en una experiencia fluida, intuitiva y sin preocupaciones.

A través de una interfaz móvil equipada con herramientas prácticas y fáciles de usar, el sistema permite calificar, hacer seguimiento de entregas y organizar a los estudiantes de forma ágil. Al automatizar la sincronización y resolver el desorden en segundo plano, Orienflow le devuelve el control al maestro para que pueda enfocarse en lo que realmente importa: enseñar, evaluar de manera dinámica y disfrutar de su labor sin estrés.


## Documentation



 - [Work architecture](#Work-architecture)



## Work architecture
```
orienflow/
├── cmd/
│   └── api/
│       └── main.go          # Punto de entrada de la app (arranque del servidor)
├── internal/
│   ├── core/                # EL NÚCLEO (No depende de nada externo)
│   │   ├── domain/          # Entidades puras (Teacher, Student, Grade, Task)
│   │   │   ├── studen.go
│   │   │   ├── grade.go
|   |   |   ├── event.go
|   |   |   ├── notification.go
|   |   |   ├── course.go
|   |   |   └── task.go
│   │   └── ports/           # Interfaces (Contratos que definen qué hace el sistema)
│   │   |   ├── incoming.go  # Qué comandos recibe el núcleo (Casos de uso)
│   │   |   ├── outgoing.go  # Qué necesita el núcleo de afuera (ej: Classroom)
|   |   |   └── token.go
|   |   └── services/
|   |       └── notification_service.go
|   |   
│   ├──  adapters/            # LOS ADAPTADORES (Implementan los puertos)
│   |    ├── input/           # Entrada: Quien manipula el núcleo
│   |    │   ├── http/        # Controlador de la API para la App Móvil
│   |    │   |    └── routes.go
|   |    |   |
|   |    |   └── webhook.go
│   |    └── output/          # Salida: Lo que el núcleo manipula hacia afuera
│   |       ├── google/      # Conector optimizado a la API de Google Classroom
│   |       |   ├── classroom.go
|   |       |   ├── google_client_adapter.go
|   |       |   └── google_limit.go
|   |       |
|   |       └── storege/
|   |              └── token_storage.go
|   |
|   ├── graph/  # Actualizacion a Graph para optimizacion de busquedas
|   |     ├── resolveres.go
|   |     ├── sche,a.go
|   |     └── server.go
|   |
|   └── shared/
|         ├── error/ # Manejador de errores
|         |      └── errors.go
|         ├── logger/ # Registros de cambios
|         |      └── logger.go
|         |
|         └── validation/ # Validacion
|                └── validation.go
|
├── go.mod                   # Definición del módulo (github.com/tu-usuario/orienflow)
└── go.sum
```

## Authors

- [@tegnoword](https://www.github.com/tegnoword)

