# URL Shortener App

---

### Table of Contents

- [Description](#description)
- [How To Use](#how-to-use)
- [Author Info](#author-info)

---

## Description

Aim of this project is taking a URL and custom short name for it with keeping duration in Redis cache.

## Technologies

### Main Technologies

- [Go](https://go.dev/)
- [Fiber Framework](https://github.com/gofiber/fiber/v2)
- [Redis](https://redis.io/)
- [Docker](https://www.docker.com/)

### Libraries

- [go-redis/redis](https://github.com/go-redis/redis/v8)
- [google/uuid](https://github.com/google/uuid)
- [joho/godotenv](https://github.com/joho/godotenv)
- [asaskevich/govalidator](https://github.com/asaskevich/govalidator)

[Back To The Top](#URL-Shortener-App)

---

## How To Use

### Tools

- [Go](https://go.dev/dl/)
- [Redis](https://redis.io/download/)
- [Docker](https://www.docker.com/get-started/)

### Run the App

#### After installing necessary tools

- Login docker and run this command in terminal

```
docker-compose up -d
```

### Give it a try

| Request Name | Request Type | URL          | JSON                                  |
| ------------ | ------------ | ------------ | ------------------------------------- |
| Resolve URL  | GET          | :3000/{url}  |                                       |
| Shorten URL  | POST         | :3000/api/v1 | {"url": "", "short": "", "expiry": 0} |

## Author Info

- Twitter - [@dev_bck](https://twitter.com/dev_bck)

[Back To The Top](#URL-Shortener-App)
