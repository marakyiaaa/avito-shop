#### Project Structure
.
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── main.go
├── docker-compose.yml
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   ├── config.go
│   │   └── config_test.go
│   ├── handler
│   │   ├── auth.go
│   │   ├── handlers_test
│   │   │   ├── auth_test.go
│   │   │   ├── info_test.go
│   │   │   └── store_test.go
...............................................



#### Run  

1) make run - запуск Docker контейнеров / make go - запуск main  
2) с помощью Postman подаются запросы для наглядного примера, как работает программа  
3) make test - запуск unit тестов  
4) make cover - покрытие кода  
5) ---------------------------------  
5) make check - запуск golangci.yaml  
6) make test_k6 - запуск нагрузочного тестирования  


P.S:  
- в данных из README, данного для выполнения тестового задания было указано, что "/api/buy/:item" - GET запрос,
я посчитала, что запрос выполняет функционал POST запроса.  
- скрытые файлы намеренно запушены, для лучшего понимания и управления программой   