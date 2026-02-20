Here is the fully formatted, raw Markdown text. You can copy everything inside the code block below and paste it directly into your `README.md` file.


# 📝 Gin + MySQL + Elasticsearch Blog Application

A full-stack backend blog application built with:
* ⚡ **Gin** (HTTP Web Framework)
* 🗄 **MySQL** (Primary Database via GORM)
* 🔎 **Elasticsearch 8** (Search Engine Integration)
* 🧩 **GORM** (ORM for Go)

This project demonstrates how to:
* Build REST routes using Gin
* Store blog data in MySQL
* Index blog data into Elasticsearch
* Create search-ready infrastructure



## 🚀 Features

* List all blogs (`/blogs`)
* View single blog (`/blogs/:id`)
* Index blogs into Elasticsearch (`/blogs/index`)
* Auto database migration with GORM
* Secure Elasticsearch 8 connection

---

## 🏗 Project Structure

```text
ElasticSearch/
│
├── controllers/
├── models/
│   └── setup.go
├── templates/
│   ├── blogs/
│   └── layouts/
├── main.go
└── go.mod

```

---

## ⚙️ Requirements

* Go 1.20+
* MySQL 8+
* Elasticsearch 8+
* Windows / Mac / Linux

---

## 🗄 MySQL Setup

1. **Login to MySQL:**
```bash
mysql -u root -p

```


2. **Run the following queries:**
```sql
CREATE DATABASE gin_elasticsearch;
CREATE USER 'gin_elasticsearch'@'localhost' IDENTIFIED BY 'tmp_pwd';
GRANT ALL PRIVILEGES ON gin_elasticsearch.* TO 'gin_elasticsearch'@'localhost';
FLUSH PRIVILEGES;

```



---

## 🔎 Elasticsearch Setup (Windows ZIP Method)

1. **Download Elasticsearch** from: [https://www.elastic.co/downloads/elasticsearch](https://www.elastic.co/downloads/elasticsearch)
2. **Extract** the downloaded file.
3. **Start the service:**
```bash
bin\elasticsearch.bat

```


4. **Reset password (if needed):**
```bash
elasticsearch-reset-password -u elastic

```


5. **Test the connection:** Visit `https://localhost:9200`

### 🔐 Elasticsearch Config Used in Project

```go
cfg := elasticsearch.Config{
    Addresses: []string{
        "https://localhost:9200",
    },
    Username: "elastic",
    Password: "YOUR_PASSWORD",
    Transport: &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true, // DEV only
        },
    },
}

```

---

## 📦 Install Dependencies

```bash
go mod tidy

```

---

## ▶️ Run the Application

```bash
go run .

```

**Server runs on:** `http://localhost:8080`

---

## 🌐 Available Routes

| Route | Description |
| --- | --- |
| `/blogs` | List all blogs |
| `/blogs/:id` | View single blog |
| `/blogs/index` | Index all blogs into Elasticsearch |

---

## 🧪 Insert Sample Data

Run this in your MySQL console to populate the database:

```sql
INSERT INTO blogs (title, content) VALUES
('Getting Started with Go', 'Go is fast and simple.'),
('Understanding Gin', 'Gin makes building APIs easy.');

```

---

## 🔍 Elasticsearch Search Example

After hitting the `/blogs/index` endpoint, test your search functionality:

* **Endpoint:** `https://localhost:9200/blogs/_search`
* **Query Search:** `https://localhost:9200/blogs/_search?q=Go`

---

## 🧠 Architecture Overview

```text
       Browser
          ↓
Gin (Router + Controllers)
          ↓
MySQL (Primary Storage via GORM)
          ↓
Elasticsearch (Search Index)

```

* **MySQL** stores persistent data.
* **Elasticsearch** stores searchable copies.

---

## 📌 Notes

* Elasticsearch 8 uses HTTPS + security by default.
* `InsecureSkipVerify` is used for local development only.
* Calling `/blogs/index` pushes MySQL data into Elasticsearch.
* Logs printed in the terminal represent backend operations, not browser responses.

---

## 📚 Learning Goals

This project helps understand:

* Go backend architecture
* ORM usage with GORM
* Search engine integration
* RESTful design
* Service initialization order

---

## 🧑‍💻 Author

Developed as a learning project integrating:

* Go backend development
* Database management
* Distributed search systems

