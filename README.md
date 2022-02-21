# Mini Project

Mini Project Authentication

## Authors

- [@apryonoboman](https://github.com/apryono)


# Information

## Mini Project
### Server
- Location of Routers, Handlers
- Bootstrap location register all routers and depedencies
- Handlers location of all Handler
- Location of main.go

### Middlewares
- Location of jwt and others

### Configs
- Location of load all configurations

## Usecase 
- Location of all Usecase, View Model, Requests


## How To Run

First use should git clone the project in folder github.com

```bash
  cd go/src/github.com

  git clone git@github.com:apryono/k-digital-project.git
```

After that git pull from develop

```bash
  git pull origin dev
```

Then go to folder server

```bash
  cd server

  go run .
```

And dont forget run your redis server 
```bash
    run redis-server
```

Dont forget to fill in your database

``` bash
APP_HOST=127.0.0.1:3000
APP_CORS_DOMAIN=http://127.0.0.1

TOKEN_SECRET={{Your_token_secret}}
TOKEN_REFRESH_SECRET={{Your_token_refresh_secret}}
TOKEN_EXP_SECRET=72
TOKEN_EXP_REFRESH_SECRET=720

AES_KEY={{Your_aes_key}}
AES_FRONT_KEY={{Your_front_aes_key}}
AES_FRONT_IV={{Your_front_aes_iv}}

REDIS_HOST=127.0.0.1:6379
REDIS_PASSWORD=

APP_PRIVATE_KEY_LOCATION=../key/id_rsa
APP_PRIVATE_KEY_PASSPHRASE=

APP_TIMEOUT=60

# LOCAL
 DATABASE_HOST={{your_database_host}}
 DATABASE_DB={{Your_database_db}}
 DATABASE_USER={{Your_database_user}}
 DATABASE_PASSWORD={{Your_database_password}}
 DATABASE_PORT=5432
 DATABASE_SSL_MODE=disable
 DATABASE_MAX_CONNECTION=5
 DATABASE_MAX_IDLE_CONNECTION=5
 DATABASE_MAX_LIFETIME_CONNECTION=10
```
