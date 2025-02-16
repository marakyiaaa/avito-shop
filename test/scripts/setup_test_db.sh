#!/bin/bash

psql -h localhost -U postgres <<EOF
CREATE DATABASE shop_test;
CREATE USER testuser WITH PASSWORD 'testpassword';
GRANT ALL PRIVILEGES ON DATABASE shop_test TO testuser;
EOF

#echo "Тестовая база данных и пользователь созданы."