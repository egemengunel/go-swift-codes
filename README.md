# SWIFT Codes API

A Go-based microservice that imports a spreadsheet of bank SWIFT (BIC) codes into SQLite and exposes four RESTful endpoints for querying, adding, and removing codes.

## Table of Contents

1. [Prerequisites](#prerequisites)  
2. [Installation & Setup](#installation--setup)  
3. [Running the Application](#running-the-application)  
4. [API Endpoints](#api-endpoints)  
5. [Running Tests](#running-tests)
6. [Note to Remitly Team](#note-to-remitly-team)

---

## Prerequisites

- **Go** ≥ 1.22 installed and in your `PATH`  
- **SQLite** CLI (optional, for manual DB inspection)  
- **Git** (to clone the repository)  
- **curl** or **Postman** or a **browser** (to exercise the API)  

---

## Installation & Setup

Follow these steps to prepare your environment and get the project running:

1. **Install Go** (if not already installed)
   - Download and install Go 1.22 or later from the official site:
     https://golang.org/dl/
   - Verify the installation:
     ```bash
     go version
     ```
     You should see output like `go version go1.24.2 darwin/arm64`.

2. **Verify Git is installed** (if not)
   ```bash
   git --version
   ```

3. **Clone the repository**
   ```bash
   git clone https://github.com/egemengunel/go-swift-codes.git
   ```

4. **Install project dependencies**
   ```bash
   go mod download
   ```

5. **The spreadsheet file is already present in the root**

   The `data/SWIFT_CODES.xlsx` excel data file already exists in the project root. So do not add the excel file to the project root again.

---

## Running the Application

On first run, the app will:

- Create (or open) `swift_codes.db` in the project root  
- Create the `swift_codes` table if missing  
- Parse and import all rows from `data/SWIFT_CODES.xlsx`  

Start the server:

```bash
go run main.go
```

You should see:

```
Server is starting on :8080...
```

The API is now listening on **http://localhost:8080**.

---

## API Endpoints

All requests and responses use **JSON**.

### 1) Get a single SWIFT code

**Request**  
```
GET http://localhost:8080/v1/swift-codes/{swiftCode}
```

**Response**  
- **200 OK** with JSON either:
  - **Head-office** with nested `branches` array  
  - **Branch** (no `branches` field)

**Example (Head-office)**
```json
{
  "address": "23 BOULEVARD PRINCESSE CHARLOTTE …",
  "bankName": "CREDIT AGRICOLE MONACO …",
  "countryISO2": "MC",
  "countryName": "MONACO",
  "isHeadquarter": true,
  "swiftCode": "AGRIMCM1XXX",
  "branches": [
    {
      "address": "…",
      "bankName": "…",
      "countryISO2": "MC",
      "countryName": "MONACO",
      "isHeadquarter": false,
      "swiftCode": "AGRIMCM1…"
    }
  ]
}
```

### 2) List all SWIFT codes for a country

**Request**  
```
GET http://localhost:8080/v1/swift-codes/country/{ISO2}
```

**Response**  
- **200 OK**

**Example**
```json
{
  "countryISO2": "PL",
  "countryName": "POLAND",
  "swiftCodes": [
    {
      "address": "…",
      "bankName": "…",
      "countryISO2": "PL",
      "countryName": "POLAND",
      "isHeadquarter": true,
      "swiftCode": "ALBPPLPWXXX"
    },
    {
      "address": "…",
      "bankName": "…",
      "countryISO2": "PL",
      "countryName": "POLAND",
      "isHeadquarter": false,
      "swiftCode": "ALBPPLPWCUS"
    }
  ]
}
```

### 3) Add a new SWIFT code

**Request**  
```
POST http://localhost:8080/v1/swift-codes
Content-Type: application/json
```

**Body**  
```json
{
  "address":       "123 TEST AVENUE",
  "bankName":      "MY TEST BANK",
  "countryISO2":   "ZZ",
  "countryName":   "ZELAND",
  "isHeadquarter": false,
  "swiftCode":     "ZZTESTBANKXXX"
}
```

**Response**  
- **201 Created**  
```json
{ "message": "swift code created" }
```

### 4) Delete a SWIFT code

**Request**  
```
DELETE http://localhost:8080/v1/swift-codes/{swiftCode}
```

**Response**  
- **200 OK**  
```json
{ "message": "swift code deleted" }
```

---

## Running Tests

Two test suites cover core logic:

1. **Service (repository) tests**  
   ```bash
   go test ./service
   ```
2. **Handler tests**  
   ```bash
   go test ./handler
   ```

Or run **all** at once:

```bash
go test ./...
```

---
## Note to Remitly Team

Thank you for reviewing my submission! A few personal notes:

- I am primarily an **iOS/Swift** developer and this was my **first time working with Go**.
- I focused on clear, idiomatic Go code while ensuring I could meet all of the exercise requirements.
- I have tried learn and understand Go and its syntax as much as i can in the limited amount of time to reduce my reliance on AI Tools for coding.
- Writing unit tests and HTTP handler tests in Go was a new challenge, and I learned a lot about interfaces, dependency injection, REST API's and testing with `httptest`.
- I appreciate your patience and hope this solution demonstrates my ability to learn new technologies quickly.

Thank you for your time and consideration!

---