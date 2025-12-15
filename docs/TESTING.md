# TESTING.md – PowerShell testavimo scenarijai

Šitas failas parodo **kaip praktiškai patikrinti**, kad:

- veikia registracija / login / session
- veikia rolės (admin / moderator / user)
- veikia moderator-only endpoint (DELETE review)
- veikia external API integracija + cache (CheapShark)

**Base URL:**
```
http://localhost:3000
```

> Visus testus vykdyk **PowerShell** terminale (Windows aplinkoje).

---

## 0) Serverio paleidimas

1. Paleisk PostgreSQL
2. Paleisk migracijas:
```powershell
goose up
```

3. Sugeneruok sqlc:
```powershell
sqlc generate
```

4. Paleisk serverį:
```powershell
go run main.go
```

---

## 1) Registracija – sukurti user

```powershell
$body = @{
  name="Jonas"
  surname="User"
  username="jonasuser"
  email="jonas@example.com"
  password="123456"
} | ConvertTo-Json

Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/auth/register" `
  -Method POST `
  -ContentType "application/json" `
  -Body $body
```

Tikėtina: 200 OK ir grąžinamas sukurtas user objektas

---

## 2) Login (user) + session cookie

```powershell
$sessionUser = New-Object Microsoft.PowerShell.Commands.WebRequestSession

$body = @{ email="jonas@example.com"; password="123456" } | ConvertTo-Json

Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/auth/login" `
  -Method POST `
  -ContentType "application/json" `
  -Body $body `
  -WebSession $sessionUser
```

Tikėtina: 200 OK ir `Set-Cookie: session_token=...`

---

## 3) /auth/me testas

### 3.1 Be sesijos – turi neveikti

```powershell
Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/auth/me" `
  -Method GET
```

Tikėtina: 401 Unauthorized

### 3.2 Su sesija – turi veikti

```powershell
(Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/auth/me" `
  -Method GET `
  -WebSession $sessionUser).Content
```

Tikėtina: grąžina vartotoją „Jonas“

---

## 4) Admin login + admin-only endpoint (/user)

### 4.1 Admin login

```powershell
$sessionAdmin = New-Object Microsoft.PowerShell.Commands.WebRequestSession

$body = @{ email="admin@example.com"; password="password123" } | ConvertTo-Json

Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/auth/login" `
  -Method POST `
  -ContentType "application/json" `
  -Body $body `
  -WebSession $sessionAdmin
```

### 4.2 Admin pasiekia /user – turi veikti

```powershell
(Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/user" `
  -Method GET `
  -WebSession $sessionAdmin).Content
```

Tikėtina: 200 OK ir vartotojų sąrašas

### 4.3 Paprastas user pasiekia /user – turi NEveikti

```powershell
Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/user" `
  -Method GET `
  -WebSession $sessionUser
```

Tikėtina: 403 Forbidden

---

## 5) Review testui: reikia reviewID

### 5.1 Gauti game_id

```powershell
(Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/games" `
  -Method GET).Content
```

Nukopijuok vieno žaidimo `id` (UUID)

### 5.2 Gauti user_id

```powershell
(Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/auth/me" `
  -Method GET `
  -WebSession $sessionUser).Content
```

Nukopijuok `id` (UUID)

### 5.3 Sukurti review (jei turi POST /reviews endpoint)

```powershell
$gameId = "PASTE_GAME_UUID_HERE"
$userId = "PASTE_USER_UUID_HERE"

$body = @{
  game_id = $gameId
  user_id = $userId
  rating  = 5
  comment = "Test review"
} | ConvertTo-Json

Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/reviews" `
  -Method POST `
  -ContentType "application/json" `
  -Body $body `
  -WebSession $sessionUser
```

Atsakyme turi būti `review_id` → išsisaugok jį

---

## 6) Moderator-only test: DELETE review

### 6.1 Paprastas user bando trinti (turi NEveikti)

```powershell
$reviewId = "PASTE_REVIEW_UUID_HERE"

Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/reviews/$reviewId" `
  -Method DELETE `
  -WebSession $sessionUser
```

Tikėtina: 403 Forbidden

---

## 7) Suteikiam „Jonui“ moderator rolių per DB

```sql
INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u, roles r
WHERE u.email='jonas@example.com' AND r.name='moderator'
ON CONFLICT DO NOTHING;
```

### Patikrinimas:

```sql
SELECT u.email, r.name
FROM users u
JOIN user_roles ur ON ur.user_id = u.id
JOIN roles r ON r.id = ur.role_id
WHERE u.email='jonas@example.com';
```

Turi rodyti `moderator`

---

## 8) Moderator login + DELETE review (turi veikti)

```powershell
$sessionMod = New-Object Microsoft.PowerShell.Commands.WebRequestSession

$body = @{ email="jonas@example.com"; password="123456" } | ConvertTo-Json

Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/auth/login" `
  -Method POST `
  -ContentType "application/json" `
  -Body $body `
  -WebSession $sessionMod
```

```powershell
$reviewId = "PASTE_REVIEW_UUID_HERE"

Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/reviews/$reviewId" `
  -Method DELETE `
  -WebSession $sessionMod
```

Tikėtina: 204 No Content arba 200 OK

---

## 9) External API (CheapShark) + cache test

### 9.1 Pirmas kvietimas – lėtesnis (per API)

```powershell
Measure-Command {
  Invoke-WebRequest -UseBasicParsing `
    -Uri "http://localhost:3000/external/deals?title=witcher" | Out-Null
}
```

### 9.2 Antras kvietimas – greitesnis (iš cache)

```powershell
Measure-Command {
  Invoke-WebRequest -UseBasicParsing `
    -Uri "http://localhost:3000/external/deals?title=witcher" | Out-Null
}
```

Tikėtina: 1 kartas lėtesnis, 2 kartas greitesnis

### Žiūrėti turinį

```powershell
(Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/external/deals?title=witcher").Content
```

---

## 10) Logout test

```powershell
Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/auth/logout" `
  -Method POST `
  -WebSession $sessionUser
```

Po logout /auth/me turi grąžinti 401:

```powershell
Invoke-WebRequest -UseBasicParsing `
  -Uri "http://localhost:3000/auth/me" `
  -Method GET `
  -WebSession $sessionUser
```

Tikėtina: 401 Unauthorized
