# API Documentation

## Base URL

```
http://localhost:3000
```

---

## Authentication

### POST /auth/register

Registruoja naują vartotoją.

**Request Body (JSON):**
```json
{
  "name": "Jonas",
  "surname": "User",
  "username": "jonasuser",
  "email": "jonas@example.com",
  "password": "123456"
}
```

**Behavior:**

- Sukuriamas naujas vartotojas
- Slaptažodis suhashinamas
- Automatiškai priskiriama rolė `user`

---

### POST /auth/login

Prisijungimas naudojant el. paštą ir slaptažodį.

**Request Body (JSON):**
```json
{
  "email": "jonas@example.com",
  "password": "123456"
}
```

**Response:**

- Sukuriama vartotojo sesija
- `session_token` grąžinamas per HTTP cookie

---

### GET /auth/me

Grąžina prisijungusio vartotojo informaciją.

**Reikalavimai:**

- Vartotojas turi būti prisijungęs

---

### POST /auth/logout

Atjungia vartotoją.

**Behavior:**

- Vartotojas turi būti prisijungęs
- Sesija ištrinama iš DB

---

## OAuth2 Authentication (Battle.net)

### GET /auth/battlenet/login

Nukreipia vartotoją į Battle.net prisijungimo puslapį.

---

### GET /auth/battlenet/callback

OAuth2 callback endpoint.

**Behavior:**

- Gaunami vartotojo duomenys iš Battle.net
- Sukuriamas arba atnaujinamas vartotojas DB
- Sukuriama sesija

---

## Users

### GET /user

Grąžina visų vartotojų sąrašą.

**Access:** Tik `admin`

---

### POST /user

Sukuria naują vartotoją (naudojama testavimui / development).

---

## Games

### GET /games

Grąžina žaidimų sąrašą.

---

## Reviews

### DELETE /reviews/{reviewID}

Ištrina pasirinktą review.

**Path Parameters:**

- `reviewID` – review UUID

**Access:** Tik `moderator`

---

## External API Integration

### GET /external/deals?title={title}

Grąžina žaidimų pasiūlymus iš CheapShark API.

**Query Parameters:**

- `title` – žaidimo pavadinimas

**Example:**
```
GET /external/deals?title=witcher
```

**Behavior:**

- Jei atsakymas jau yra cache DB – grąžinamas iš DB
- Jei cache nėra arba pasibaigęs – kviečiamas CheapShark API
- Atsakymas cache’inamas DB (~10 min)

---

## Roles & Permissions

Sistemoje naudojamos šios rolės:

- `user`
- `moderator`
- `admin`

### Role examples:

| Endpoint                        | Required Role       |
|---------------------------------|---------------------|
| `GET /user`                     | admin               |
| `DELETE /reviews/{reviewID}`    | moderator           |
| `GET /auth/me`                  | authenticated user  |

---

## Error Handling

- `401 Unauthorized` – vartotojas neprisijungęs
- `403 Forbidden` – nepakankamos teisės
- `400 Bad Request` – neteisingi duomenys
- `404 Not Found` – endpoint nerastas
