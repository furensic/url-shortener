# url-shortener

![GitHub language count](https://img.shields.io/github/languages/count/furensic/url-shortener?style=plastic)
![GitHub go.mod Go version (branch)](https://img.shields.io/github/go-mod/go-version/furensic/url-shortener/main?style=plastic)
![GitHub commit activity](https://img.shields.io/github/commit-activity/w/furensic/url-shortener?style=plastic)
![GitHub License](https://img.shields.io/github/license/furensic/url-shortener?style=plastic)

A simple server application that offers HTTP redirection via an REST API.

## whats next?

i need to check if a user with username already exists. Otherwise the DB returns an error `ERROR: duplicate key value violates unique constraint "users_username_key" (SQLSTATE 23505)` but the index (id) increases anyways. I suppose i could check if the username exists inside 'service/user.go'? done

i need to create a login endpoint that is able to check user data and then maybe use a middleware or something to inject a token, cookie or something?
