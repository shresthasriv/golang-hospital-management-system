# Hospital Management System

A RESTful API-based hospital management system built with Go, featuring separate portals for receptionists and doctors with role-based access control.

## 🎥 Demo & Documentation

- **Video Demo**: [https://www.loom.com/share/796baea39af0417ea2460290059f6d47](https://www.loom.com/share/796baea39af0417ea2460290059f6d47?sid=3853c498-461a-4728-babe-a10d85c444e8)
- **API Documentation**: [Postman Collection](https://www.postman.com/cryosat-explorer-49065860/makerble-assessment-api/collection/e8pp30m/hospital-management-system-api)

## 🚀 Quick Start

### Option 1: Using Makefile
```bash
make deps
make run
```
*Requires PostgreSQL database. Configure `DATABASE_URL` in `.env` (see `.env.example`)*

### Option 2: Using Docker
```bash
docker compose up
```
*No external dependencies required. Includes PostgreSQL container.*
