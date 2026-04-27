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
                       Redis (Geo + State)
```

### Components:

* **Gin** → HTTP server
* **Service Layer** → business logic
* **Repository Layer** → DB abstraction
* **PostgreSQL** → transactional data (rides)
* **Redis** → real-time driver location, availability & active ride mapping
* **New Relic** → performance monitoring

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

### 🔄 Set Driver Availability

```
POST /v1/drivers/:id/availability
```

---

### ✅ Accept Ride

```
POST /v1/drivers/:id/accept
```

---

### 🏁 End Ride (Driver Based)

```
POST /v1/drivers/:id/end
```

Ends the active ride assigned to the driver using Redis mapping.

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

### 👥 Driver List

```
GET /v1/drivers
```

---

## ⚡ Key Features

### 🚀 Real-Time Driver Matching

* Redis GEO used for nearest driver lookup
* Radius-based matching
* Filters only AVAILABLE drivers

---

### 🔄 Ride Lifecycle

```
MATCHED → ONGOING → COMPLETED
```

* Enforced transitions
* Prevents invalid state changes

---

### 🔗 Driver → Ride Mapping (Important)

* Stored in Redis:

  ```
  driver_active_ride:<driver_id> → ride_id
  ```
* Enables:

  * O(1) lookup of active ride
  * Driver-based ride actions (end ride)
* Avoids reliance on frontend state

---

### 💰 Pricing

* Flat fare (mock implementation)
* Easily extendable to distance/surge pricing

---

### 📊 Monitoring (New Relic)

* Middleware instrumentation
* Tracks API latency & performance

---

### 🧪 Unit Testing

* Service layer tested with mock repositories
* Covers:

  * Create ride
  * Accept ride
  * End ride

---

### 🧠 Scalability Design

* Stateless APIs → horizontal scaling
* Redis for low-latency lookups
* DB used only for persistence
* Separation of:

  * Real-time state (Redis)
  * Persistent data (Postgres)

---


### ⚡ API Latency Optimizations
- Redis used for real-time driver matching (avoids DB hits in hot path)
- Added caching for GET /v1/rides/:id to reduce repeated polling load
- Implemented cache invalidation on ride state updates (accept/end)
- Database indexing on frequently queried fields (driver_id, status)
- State validation ensures consistency under concurrent requests

### 🔒 Data Consistency

* Driver–ride assignment validation
* Availability checks before matching
* Prevents:

  * double assignment
  * invalid accept/end flows

---

## 🎨 Frontend

Simple interactive UI:

* Add multiple drivers
* Update driver location
* Toggle availability
* Create ride with custom coordinates
* Accept ride (only assigned driver)
* End ride (driver-based)
* Live updates via polling (2s)

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

1. Add multiple drivers (different locations)
2. Mark one driver unavailable
3. Create ride near a specific driver
4. System matches nearest available driver
5. Driver accepts ride → status becomes ONGOING
6. End ride via driver → status becomes COMPLETED
7. Driver becomes AVAILABLE again

---

## 🧰 Tech Stack

* Go (Golang)
* Gin
* PostgreSQL
* Redis
* Docker
* New Relic

---

## ✨ Bonus Features

* Driver availability toggle
* Driver-based ride termination
* Real-time driver status tracking

---

## ⚠️ Notes

* Redis is required for matching
* PostgreSQL stores ride data
* New Relic is optional

---

## 🎯 Future Improvements

* WebSockets (remove polling)
* Dynamic surge pricing
* Kafka-based async processing
* Multi-region deployment
* Driver ETA estimation

---

## 👨‍💻 Author

Aadith S
