# url-shortener

![GitHub language count](https://img.shields.io/github/languages/count/furensic/url-shortener?style=plastic)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/furensic/url-shortener/main?style=plastic)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/furensic/url-shortener?style=plastic)
![GitHub License](https://img.shields.io/github/license/furensic/url-shortener?style=plastic)

A simple server application that offers HTTP redirection via an REST API.

## whats next?

i need to check if a user with username already exists. Otherwise the DB returns an error `ERROR: duplicate key value violates unique constraint "users_username_key" (SQLSTATE 23505)` but the index (id) increases anyways. I suppose i could check if the username exists inside 'service/user.go'? done

i need to create a login endpoint that is able to check user data and then maybe use a middleware or something to inject a token, cookie or something?

## to check

/health has about 19ms of delay when using basic auth with argon2id. Have to check if JWT is more performant so i dont have to Verify login credential each time.

```mermaid
flowchart LR
    logger[Initialize logger]
    application[Initialize application struct]
    logger --> application
    app_config[Initialize application config struct]
    app.config_Struct@{ shape: lean-r, label: "app.config"}
    appConfig --- appConfigStruct
    NewPostgresDatabase@{ shape: subproc, label "NewPostgresDatabase" }
    db <-- calls with postgres params --- NewPostgresDatabase
    shortenedUriRepo -- calls with db_conn parameter --- db --> NewShortenedUriPostgresAdapter@{ shape: subproc }
    userRepo -- calls with db_conn parameter --- db --> NewUserPostgresAdapter@{ shape: subproc }
    repositories <--- shortenedUriRepo
    repositories <--- userRepo


```
