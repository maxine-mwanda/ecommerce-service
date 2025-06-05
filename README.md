E-Commerce Service API
Go
MySQL
Docker
Kubernetes

A robust e-commerce backend service built with Go, featuring product catalog management, order processing, customer authentication, and notification systems.

Table of Contents
Features

Technologies

Architecture

Setup

Prerequisites

Local Development

Docker Setup

Kubernetes Deployment

API Documentation

Testing

Monitoring

Troubleshooting

Contributing

Features
Product Management

Hierarchical category system

Product CRUD operations

Price analytics by category

Order Processing

Order creation with validation

Automatic total calculation

Inventory checks

Customer System

JWT/OIDC authentication

Customer profiles

Order history

Notifications

SMS order confirmations (Africa's Talking API)

Admin email notifications

Infrastructure

Docker containerization

Kubernetes deployment

CI/CD pipeline

Technologies
Backend: Go 1.19+

Database: MySQL 5.7+

Authentication: JWT + OIDC

Containerization: Docker 20.10+

Orchestration: Kubernetes 1.20+

Logging: Zerolog with rotation

Testing: Go test + testify

Architecture
Diagram
Code











Setup
Prerequisites
Go 1.19+

MySQL 5.7+ or Docker

(Optional) Kubernetes cluster

Local Development
Clone the repository:

bash
git clone https://github.com/yourusername/ecommerce-service.git
cd ecommerce-service
Set up environment variables:

bash
cp .env.example .env
# Edit .env with your configuration
Install dependencies:

bash
go mod download
Start MySQL:

bash
docker run --name ecommerce-mysql -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=ecommerce -p 3306:3306 -d mysql:5.7
Run the application:

bash
go run cmd/main.go
Docker Setup
Build and run:

bash
docker-compose up --build
Verify containers:

bash
docker-compose ps
Kubernetes Deployment
Set up Kubernetes cluster

Deploy resources:

bash
kubectl apply -f kubernetes/configmap.yaml
kubectl apply -f kubernetes/secrets.yaml
kubectl apply -f kubernetes/deployment.yaml
kubectl apply -f kubernetes/service.yaml
Check deployment:

bash
kubectl get pods
kubectl get service ecommerce-service
API Documentation
Base URL
http://localhost:8080/api/v1

Endpoints
Authentication
POST /auth/login - Get JWT token

POST /auth/refresh - Refresh token

Products
GET /products - List all products

POST /products - Create new product

GET /products/{id} - Get product details

GET /categories/{id}/products - Get products by category

Orders
POST /orders - Create new order

GET /orders/{id} - Get order details

GET /customers/{id}/orders - Get customer's orders

Customers
POST /customers - Register new customer

GET /customers/{id} - Get customer profile

View complete API documentation

Testing
Unit Tests
bash
go test -v ./tests/unit/...
Integration Tests
bash
# Start test database
docker run --name test-mysql -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=testdb -p 3306:3306 -d mysql:5.7

# Run tests
export $(cat .env.test | xargs)
go test -v -tags=integration ./tests/integration/...
Test Coverage
bash
go test -coverprofile=coverage.out -tags=integration ./...
go tool cover -html=coverage.out
Monitoring
The service includes built-in monitoring endpoints:

GET /health - Service health check

GET /metrics - Prometheus metrics

GET /debug/pprof - Performance profiling

Troubleshooting
Common Issues
Database connection failures:

Verify MySQL is running

Check credentials in .env

Test connection manually: mysql -u root -p -h 127.0.0.1 -P 3306 ecommerce

Migration errors:

Ensure schema matches models

Check foreign key constraints

Authentication problems:

Validate JWT secret

Check token expiration