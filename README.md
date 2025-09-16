## Microservices-Ecommerce
---
This project is an event-driven microservices-based e-commerce system built with Go (Golang) and deployed on Kubernetes.

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

### Choreography Saga Pattern
---
![](https://terraform-tfstate-backup-test.s3.ap-southeast-1.amazonaws.com/Screenshot%202025-09-16%20112546.png?X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Content-Sha256=UNSIGNED-PAYLOAD&X-Amz-Credential=ASIAUHQ3FK4P4TZYAMRZ%2F20250916%2Fap-southeast-1%2Fs3%2Faws4_request&X-Amz-Date=20250916T043117Z&X-Amz-Expires=300&X-Amz-Security-Token=IQoJb3JpZ2luX2VjEA0aDmFwLXNvdXRoZWFzdC0xIkgwRgIhAMbmDoE9KWzPP4k3aoLR9mQWKUgg%2F%2BTDf23y8GIAjYF%2BAiEAy9HmJeaRiNXqNIzUvAcCLnalt4Zp7n5F52AdePsIRwEqgAMIhv%2F%2F%2F%2F%2F%2F%2F%2F%2F%2FARAAGgwyOTEwNDEwMDczOTEiDObAYzg3eqR2uP7MuirUAriBuezvL0%2BDMui8Pt%2BId%2FlCtHyQtwk8xlcFRcw7yuifGvWF0wc5RkfI%2FKwjghXG0CrZCVmSUxFQk%2FYb1%2FnB4LAB9jeUDzPFrGIPBc0ZSIliwa%2F2%2BLuxvcTOc3OWbUEY3cuC03tZelKnISvTK%2FP85zrPdl63a2sT4zUgXCGqAJ%2BUopk06lT4y7w2CWMxtZUaqcd5K%2FDQQc2kEEyJmjW27kp%2BmlbyxAAFNS0FQRAy4nLdp%2FmeU4jR9kQ8lZRgSctgoz%2F%2FwOnlOB8aXOhbgjkIeWVZjcn4dZDF05bUF3VNxjruaM5gOjZvmbJkYpcL03SpTbXi60l%2BH53t5V02gXAkLne478i2HxAniJDxNgz5E3u7MZqynlVN6mt%2FNwcV%2B%2BDSNXOXAm%2B%2BHr8lqdj3wviLyCReF8L2UWTv0VVof94Z1p3OAF1EnmIVy0EQi5oZLQonRHUmwHcwy5ajxgY6rAINGQiP%2BbcxuACz%2FcjzjJfoQdVcO0q6ZmthHaUb4Qt3VvBamOwCj66fnxZERok9V9iMo2coUPSmTJa533NorepdOX3D4k9p35Ir3sgIP5Vo3Yq%2FMSlauEVPsAIoEe%2Ba1TfaY%2FvjuQvsqWL82dkhcFFUHARqvtxItGPePHdXaDkJOh4M%2FUX1LUcS9geetmZjIhpbCQscRdVNINzhq%2FjZV%2FsfnXX8MMIeQDmeynOSiL%2BWrSUpDDaz8e0hYp0NWxoPp1ch09cTlJXsDsBXu2lkfe76XyhyR5%2FqVMu2kJvA86Z%2FIps9eD6j%2FMSzHaL68Bm5rGtvWn%2FWuTivyNHewIVbOgdiiWepXHTcbcvjYgw9eNdnCPbkv5UI8j5jiUzvxwR0kBCLcp8J3QGkluvanRM%3D&X-Amz-Signature=c8f1c2fee7670f982322f0e92a7851d7584cbce834a511107a166e1240aab3cd&X-Amz-SignedHeaders=host&response-content-disposition=inline)


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