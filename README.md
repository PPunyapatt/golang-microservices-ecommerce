## Golang-Microservices-Ecommerce
---
This project is a microservices-based e-commerce system built with Go (Golang). It follows an event-driven architecture using RabbitMQ as the message broker and includes a payment service integrated with Stripe for payment processing. The project is deployed on Kubernetes.

The system is secured using HashiCorp Vault for secrets management and follows best practices for protecting sensitive data such as database credentials, API keys, and certificates.

The entire infrastructure is provisioned and managed on Amazon Web Services (AWS) using Terraform.

### Technical Stack
---
- Backed
    - Golang
    - JWT
    - Opentelemetry
- Infrastructure
    - AWS
    - Kubernetes
    - Hashicorp Vault
    - Postgres
    - MongoDB
    - RabbitMQ
    - Jaeger
    - Stripe
    - Debezium

### Choreography Saga Pattern
---
![](https://github.com/PPunyapatt/golang-microservices-ecommerce/blob/dev/Image/diagram.png)


### Services
---
1. Gateway
2. Order
3. Cart
4. Payment
5. Inventory
6. Auth

### Starting Project
---
First, you need to prepare the infrastructure services such as Postgres, RabbitMQ, and others to support the application.

Navigate to the `Infrastructure/Infra` directory and apply `base.yml` to create the namespace and storage class:

```
> kubectl apply -f base.yml
```
Then, navigate to the `Infrastructure/Services` directory and apply all Kubernetes manifest files to deploy the services:
```
> kubectl apply -f ./
```

Finally, navigate to the `Infrastructure/nginx` directory and apply the ClusterIssuer and nginx-ingress:
```
> kubectl apply -f ./
```

After you prepare infrastructure you have to navigate to `Infrastructure/Services` for apply all services
```
> kubectl apply -f ./
```