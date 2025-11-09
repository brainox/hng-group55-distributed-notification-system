# hng-group55-distributed-notification-system


notification-system/
│
├── api_gateway/            
│
├── services/
│   ├── user_service/       
│   ├── template_service/   
│   ├── email_service/      
│   ├── push_service/       
│
├── infra/
│   ├── kafka/              
│   ├── redis/              
│   ├── postgres/           
│   ├── nginx/              
│
├── shared/                 
│   └── libs/
│       ├── circuit_breaker/
│       ├── idempotency/
│       ├── retry/
│       └── logging/             
│
├── observability/          
│   ├── prometheus/
│   ├── grafana/
│   ├── loki/
│   ├── jaeger/
│   └── alerting/
│
├── deployments/           
│   ├── docker/
│   ├── staging/
│   └── production/
│
├── .github/
│   └── workflows/          
│
└── docs/
    ├── architecture_diagram/
    ├── openapi_specs/
    └── readmes/