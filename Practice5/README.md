# Practice 5 — Go + PostgreSQL

## Setup

```bash
# 1. Create database
createdb practice5

# 2. Run migration (creates tables + seeds 20 users + friendships)
psql -d practice5 -f migrations/init.sql

# 3. Install dependencies
go mod tidy

# 4. Run server
go run main.go
```

Server starts on **http://localhost:8080**

---

## Environment Variables (optional)

| Variable    | Default     |
|-------------|-------------|
| DB_HOST     | localhost   |
| DB_PORT     | 5432        |
| DB_USER     | postgres    |
| DB_PASSWORD | postgres    |
| DB_NAME     | practice5   |

---

## Postman Examples

### 1. Pagination with order_by
```
GET http://localhost:8080/users?page=1&page_size=5&order_by=name&order_dir=asc
```

### 2. Filter by ID
```
GET http://localhost:8080/users?id=3
```

### 3. Filter by Name + Email
```
GET http://localhost:8080/users?name=alice&email=mail.com
```

### 4. Filter by 3 fields + pagination + order
```
GET http://localhost:8080/users?gender=female&order_by=birthdate&order_dir=desc&page=1&page_size=3
```

### 5. GetCommonFriends — Alice(1) and Bob(2) → returns Carol, David, Eva
```
GET http://localhost:8080/users/common-friends?user1=1&user2=2
```

---

## Common Friends Logic (no N+1)

Single JOIN query:
```sql
SELECT u.id, u.name, u.email, u.gender, u.birthdate
FROM user_friends uf1
JOIN user_friends uf2 ON uf1.friend_id = uf2.friend_id
JOIN users u          ON u.id = uf1.friend_id
WHERE uf1.user_id = $1
  AND uf2.user_id = $2
```
