# avito-shop

avito-shop/
├── Dockerfile
├── .env
├── README.md
├── app/
│   ├── cmd/
│   │   └── main.go
│   ├── docker-compose.yaml
│   ├── internal/
│   │   ├── config/
│   │   │   └── config.go
│   │   ├── handlers/
│   │   │   ├── auth.go
│   │   │   ├── coin.go
│   │   │   ├── inventory.go
│   │   │   └── transaction.go
│   │   ├── models/
│   │   │   ├── item.go
│   │   │   ├── transaction.go
│   │   │   └── user.go
│   │   ├── repository/
│   │   │   ├── item.go
│   │   │   ├── transaction.go
│   │   │   └── user.go
│   │   └── service/
│   │       ├── auth.go
│   │       ├── coin.go
│   │       ├── inventory.go
│   │       └── transaction.go
│   └── migrations/
│       └── init.sql
└── go.mod