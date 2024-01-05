# Speed Typing auth service on Golang

speed typing auth service written on golang with MongoDB and gin framework

endpoints:

- **POST /api/v1/auth/registration/**

```json
{
  "nickname": "user",
  "email": "user@gmail.com",
  "password": "password"
}
```

- **POST /api/v1/auth/login/** (login by nickname or email)

```json
{
  "login": "user",
  "password": "password"
}
```

- **GET /api/v1/auth/refresh/**
- **DELETE /api/v1/auth/logout/**
