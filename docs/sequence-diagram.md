# Promotion System Sequence Diagram

## User Registration and Campaign Participation Flow

```mermaid
sequenceDiagram
    participant M as Mobile App
    participant A as Auth Service
    participant P as Promotion Service
    participant R as Redis Cache
    participant DB as PostgreSQL

    Note over M,DB: 1. User Registration Flow
    M->>A: POST /register
    A->>DB: Create User
    DB-->>A: User Created
    A-->>M: Registration Success

    Note over M,DB: 2. User Login Flow
    M->>A: POST /login
    A->>DB: Validate Credentials
    DB-->>A: User Data
    A-->>M: JWT Token

    Note over M,DB: 3. Campaign Registration Flow
    M->>P: POST /api/v1/register (with JWT)
    P->>A: Validate Token
    A-->>P: Token Valid + User ID
    P->>R: Check Cache for Campaign
    R-->>P: Campaign Data (if cached)
    P->>DB: Get Campaign (if not cached)
    DB-->>P: Campaign Data
    P->>R: Cache Campaign Data
    P->>R: Acquire Distributed Lock
    R-->>P: Lock Acquired
    P->>DB: Check & Increment Campaign Users (Atomic)
    DB-->>P: Success/Full
    alt Campaign Not Full
        P->>DB: Create Voucher
        DB-->>P: Voucher Created
        P->>DB: Create User Registration
        DB-->>P: Registration Created
        P->>R: Cache Voucher & Registration
        P->>R: Release Lock
        P-->>M: Success + Voucher Code
    else Campaign Full
        P->>R: Release Lock
        P-->>M: Campaign Full Error
    end

    Note over M,DB: 4. Voucher Usage Flow
    M->>P: POST /api/v1/voucher/use
    P->>A: Validate Token
    A-->>P: Token Valid
    P->>R: Check Voucher Cache
    R-->>P: Voucher Data (if cached)
    P->>DB: Get Voucher (if not cached)
    DB-->>P: Voucher Data
    alt Voucher Valid
        P->>DB: Mark Voucher as Used
        DB-->>P: Voucher Updated
        P->>R: Invalidate Voucher Cache
        P-->>M: Voucher Used Successfully
    else Voucher Invalid
        P-->>M: Voucher Error
    end
```

## High-Concurrency Campaign Registration

```mermaid
sequenceDiagram
    participant U1 as User 1
    participant U2 as User 2
    participant U3 as User 3
    participant P as Promotion Service
    participant R as Redis Cache
    participant DB as PostgreSQL

    Note over U1,DB: Concurrent Registration Attempts
    U1->>P: Register for Campaign
    U2->>P: Register for Campaign
    U3->>P: Register for Campaign

    P->>R: Acquire Lock (User 1)
    R-->>P: Lock Acquired
    P->>R: Acquire Lock (User 2)
    R-->>P: Lock Denied
    P->>R: Acquire Lock (User 3)
    R-->>P: Lock Denied

    P->>DB: Atomic Check & Increment (User 1)
    DB-->>P: Success (User 1)
    P->>DB: Create Voucher & Registration (User 1)
    P->>R: Release Lock (User 1)
    P-->>U1: Success + Voucher

    P->>R: Acquire Lock (User 2)
    R-->>P: Lock Acquired
    P->>DB: Atomic Check & Increment (User 2)
    DB-->>P: Success (User 2)
    P->>DB: Create Voucher & Registration (User 2)
    P->>R: Release Lock (User 2)
    P-->>U2: Success + Voucher

    P->>R: Acquire Lock (User 3)
    R-->>P: Lock Acquired
    P->>DB: Atomic Check & Increment (User 3)
    DB-->>P: Campaign Full
    P->>R: Release Lock (User 3)
    P-->>U3: Campaign Full Error
```

## Rate Limiting and Caching Strategy

```mermaid
sequenceDiagram
    participant C as Client
    participant LB as Load Balancer
    participant P1 as Promotion Service 1
    participant P2 as Promotion Service 2
    participant R as Redis Cache
    participant DB as PostgreSQL

    Note over C,DB: Rate Limiting Flow
    C->>LB: Request
    LB->>P1: Route Request
    P1->>R: Check Rate Limit
    R-->>P1: Rate Limit Status
    alt Rate Limit OK
        P1->>R: Check Cache
        R-->>P1: Cache Miss
        P1->>DB: Query Database
        DB-->>P1: Data
        P1->>R: Cache Data
        P1-->>LB: Response
        LB-->>C: Response
    else Rate Limit Exceeded
        P1-->>LB: 429 Too Many Requests
        LB-->>C: Rate Limit Error
    end

    Note over C,DB: Cache Hit Flow
    C->>LB: Request
    LB->>P2: Route Request
    P2->>R: Check Cache
    R-->>P2: Cache Hit
    P2-->>LB: Response (from cache)
    LB-->>C: Fast Response
```

## System Health Monitoring

```mermaid
sequenceDiagram
    participant M as Monitoring
    participant A as Auth Service
    participant P as Promotion Service
    participant R as Redis
    participant DB as PostgreSQL

    Note over M,DB: Health Check Flow
    M->>A: GET /health
    A->>DB: Check Connection
    DB-->>A: Connection OK
    A-->>M: Health Status

    M->>P: GET /health
    P->>R: Check Connection
    R-->>P: Connection OK
    P->>DB: Check Connection
    DB-->>P: Connection OK
    P->>A: Check Connection
    A-->>P: Connection OK
    P-->>M: Health Status

    Note over M,DB: Metrics Collection
    M->>P: Collect Metrics
    P-->>M: Request Count, Latency, Error Rate
    M->>R: Collect Metrics
    R-->>M: Memory Usage, Hit Rate
    M->>DB: Collect Metrics
    DB-->>M: Connection Count, Query Performance
```
