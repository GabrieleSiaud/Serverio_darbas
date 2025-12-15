# Serverio_darbas – Server-Side Web Development Final Project

## Projekto aprašymas

REST API serveris sukurtas naudojant **Go (Golang)** ir **PostgreSQL**. Projektas realizuoja:

- Autentifikaciją (registracija / prisijungimas / sesijos valdymas)
- OAuth2 prisijungimą per **Battle.net**
- Role-based access control (RBAC) su 3 rolėmis: `user`, `moderator`, `admin`
- Išorinės API integraciją su **CheapShark** ir atsakymų kešavimą DB
- DB migracijas (Goose) ir automatiškai generuotą `repository` sluoksnį (`sqlc`)

## Naudotos technologijos

- **Programavimo kalba:** Go (Golang)
- **Duomenų bazė:** PostgreSQL
- **Router:** Chi
- **Migracijos:** Goose
- **SQL → Go generavimas:** sqlc
- **PostgreSQL driveris:** pgx
- **OAuth2 biblioteka:** Goth (Battle.net)

## Reikalavimų atitikimas

- Custom REST API (users, games, reviews, sessions)
- OAuth2 Authentication (Battle.net)
- External API Integration (CheapShark)
- Role-Based Access (admin/moderator/user)
- PostgreSQL schema + migracijos
- Dokumentacija (README + API dokumentacija + Testing)

## Konfigūracija

Projektas naudoja `.env` failą lokaliai (neįtrauktas į git).

### Kur dėti `.env`?

**Šaknyje**: `Serverio_darbas/.env.example`

### Ką įrašyti?

```env
# PostgreSQL
DATABASE_URL=postgres://postgres:root@localhost:5432/serverio_duomenubaze?sslmode=disable

# Battle.net OAuth2
BATTLE_CLIENT_ID=PASTE_CLIENT_ID
BATTLE_CLIENT_SECRET=PASTE_CLIENT_SECRET
BATTLE_NET_REGION=eu
```

## Paleidimas (Setup)

### 1) PostgreSQL

- Įsitikink, kad PostgreSQL paleistas.
- Sukurk duomenų bazę `serverio_duomenubaze`.

### 2) Migracijos

```bash
goose up
```

### 3) SQL → Go generavimas

```bash
sqlc generate
```

### 4) Serverio paleidimas

```bash
go run main.go
```

### 5) Serveris veikia

`http://localhost:3000`

## API dokumentacija

API dokumentacija saugoma faile: `Serverio_darbas/docs/API.md`

### Base URL

`http://localhost:3000`

### Autentifikacija

#### POST /auth/register

Sukuria vartotoją ir automatiškai priskiria rolę `user`.

**Body (JSON):**

```json
{
  "name": "Jonas",
  "surname": "User",
  "username": "jonasuser",
  "email": "jonas@example.com",
  "password": "123456"
}
```

**Response (200):**

Sėkmingas registravimas.

## External API

Naudojama public API: **CheapShark**

- Endpoint: `GET /external/deals?title=witcher`
- Atsakymai cache’inami DB (TTL ~ 10 min.)

