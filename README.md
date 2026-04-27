# 🚖 Ride Hailing System (GoComet DAW)

A scalable ride-hailing backend (Uber/Ola–like) built in Go, supporting real-time driver matching, ride lifecycle management, and performance monitoring.

---

## 🚀 Quick Start

### 1. Clone the repository

```bash
git clone https://github.com/<your-username>/ride-hailing-system.git
cd ride-hailing-system
```

---

### 2. Set up environment variables

Create a `.env` file in the root:

```env
DB_HOST=localhost
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=ride_hailing
DB_PORT=5432

REDIS_ADDR=localhost:6379

NEW_RELIC_LICENSE_KEY=your_key_here
```

---

### 3. Run the application

#### Option A — Local

```bash
go run cmd/server/main.go
```

#### Option B — Docker

```bash
docker compose up --build
```

---

## 🏗️ Architecture

```
Handler → Service → Repository → PostgreSQL
                         ↓
                       Redis (Geo)
```

### Components:

- **Gin** → HTTP server  
- **Service Layer** → business logic  
- **Repository Layer** → DB abstraction  
- **PostgreSQL** → transactional data (rides)  
- **Redis** → real-time driver location + matching  
- **New Relic** → performance monitoring  

---

## 📡 Core APIs

### 🚕 Create Ride

```
POST /v1/rides
```

```json
{
  "rider_id": "r1",
  "pickup_lat": 12.9716,
  "pickup_lng": 77.5946
}
```

---

### 📍 Update Driver Location

```
POST /v1/drivers/:id/location
```

---

### ✅ Accept Ride

```
POST /v1/drivers/:id/accept
```

---

### 🏁 End Trip

```
POST /v1/trips/:id/end
```

---

### 💳 Payment (Mock)

```
POST /v1/payments
```

---

### 🔍 Get Ride

```
GET /v1/rides/:id
```

---

### 👥 Driver Count

```
GET /v1/drivers/count
```

---

## ⚡ Key Features

### 🚀 Real-Time Driver Matching
- Redis GEO used for nearest driver lookup  
- Handles high-frequency location updates  

---

### 🔄 Ride Lifecycle

```
MATCHED → ONGOING → COMPLETED
```

- Enforced state transitions  
- Prevents double assignment  

---

### 💰 Pricing
- Simple fare calculation (mock)  
- Easily extendable  

---

### 📊 Monitoring (New Relic)
- Middleware-based instrumentation  
- Tracks latency, throughput, and API performance  

---

### 🧪 Unit Testing
- Service layer tested with mock repositories  
- Covers:
  - Create ride  
  - Accept ride  
  - End ride  

---

### 🧠 Scalability Design

- Stateless APIs → horizontal scaling  
- Redis for fast lookups  
- Minimal DB usage in hot paths  
- Designed for multi-region extension  

---

### 🔒 Data Consistency

- Ride assignment validation  
- Status transition enforcement  
- Prevents race conditions  

---

## 🎨 Frontend

Simple HTML UI included:

- Create ride  
- Accept ride  
- End ride  
- Live updates (polling every 2 seconds)  
- Driver availability toggle  

Run:

```bash
python3 -m http.server 3000
```

Open:

```
http://localhost:3000
```

---

## 🎬 Demo Flow

1. Add driver location  
2. Request ride  
3. System matches driver  
4. Driver accepts ride  
5. Trip ends  
6. Payment processed  

---

## 🧰 Tech Stack

- Go (Golang)
- Gin
- PostgreSQL
- Redis
- Docker
- New Relic

---

## ✨ Bonus Feature

### Driver Availability Toggle
Drivers can mark themselves available/unavailable, and matching respects availability.

---

## ⚠️ Notes

- Redis required for matching  
- PostgreSQL stores ride data  
- New Relic optional but included  

---

## 🎯 Future Improvements

- WebSockets for real-time updates  
- Surge pricing  
- Kafka-based event system  
- Multi-region deployment  

---

## 👨‍💻 Author

Aadith S