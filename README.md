# Ale WMS

Modern open-source WMS (Warehouse Management System) focused on simplicity, scalability and API-first architecture.

## Goals

Ale WMS aims to provide a modern warehouse management platform with:

* API-first architecture
* Modular monolith design
* Domain-Driven Design (DDD)
* Event-driven ready architecture
* PostgreSQL-first approach
* OpenAPI specification
* Easy self-hosting
* Docker-based development
* Future support for web, mobile and collector devices

The project starts simple and evolves incrementally into a complete WMS platform.

---

# Tech Stack

## Backend

* Go
* PostgreSQL
* Chi Router
* PGX
* SQLC
* Golang Migrate

## Frontend (planned)

* Next.js
* TypeScript
* TailwindCSS

## Mobile (planned)

* PWA-first approach
* Future React Native support

---

# Architecture

The project follows a modular monolith architecture using DDD principles.

```text
apps/api/internal/
├── inventory/
├── catalog/
├── receiving/
├── outbound/
├── location/
└── identity/
```

Each module contains:

```text
domain/
application/
infrastructure/
interfaces/
```

---

# Philosophy

Ale WMS prioritizes:

* Strong transactional consistency
* Explicit business rules
* Clear domain boundaries
* SQL-first data access
* Simplicity over premature complexity
* Evolutionary architecture

The project intentionally avoids premature adoption of:

* Microservices
* Distributed transactions
* Event sourcing
* Complex CQRS
* Heavy infrastructure

---

# Project Structure

```text
wms/
├── apps/
│   ├── api/
│   ├── web/
│   └── mobile/
│
├── packages/
│   ├── openapi/
│   ├── sdk-ts/
│   └── ui/
│
├── deployments/
├── scripts/
└── README.md
```

---

# Running Locally

## Requirements

* Go
* Docker
* Docker Compose

---

## Start PostgreSQL

```bash
docker compose up -d
```

---

## Run API

```bash
cd apps/api

air
```

---

# Database Migrations

Create migration:

```bash
migrate create -ext sql -dir migrations create_products_table
```

Run migrations:

```bash
migrate -path migrations \
  -database "postgres://wms:wms@localhost:5432/wms?sslmode=disable" \
  up
```

---

# OpenAPI

The API contract will be defined using OpenAPI.

```text
packages/openapi/
```

Future SDKs and frontend integrations will be generated from the OpenAPI specification.

---

# Roadmap

## v0.1

* Products
* Warehouses
* Locations
* Stock movements

## v0.2

* Reservations
* Transfers
* Inventory balances

## v0.3

* Picking
* Packing
* Shipping

## v0.4

* Barcode support
* PWA collector
* Event-driven integrations

---

# License

Apache 2.0

---

# Status

Early development 🚧
