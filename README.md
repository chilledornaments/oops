# OOPS

OOPS One-time Password Sharing is a Go tool written to provide short-lived (one hour), one-time access to a secret.

## Important Notes

This version has been deprecated in favor of v2, which drops support for MySQL and replaces it with DynamoDB. V2 still supports SQLite.

Check out the `v2` branch!

## Usage

### Database

You can either use MySQL/MariaDB or SQLite.

If you're using MySQL, you'll need to create the database first:

```mysql
create database oops;

grant all on oops.* to 'oopsuser'@'%' identified by 'oopsPASS';
```

Careful, `%` allows the user to connect to the database from anywhere. Use proper, layered access control to the database.

### .env File

Set the `OOPS_ENV_FILE` environment variable to the path of your config file for the application (this can be anywhere). You can reference the `.env.example` file for examples for mysql and sqlite.

#### Database Connection

##### SQLITE3

If you're using SQLITE3 (only recommended for development and small sites), you only need to set `DB_DRIVER` and `DB_PATH`.

`DB_DRIVER` *must* by `sqlite3`

`DB_PATH` is the path to where you want the SQLite database to live

##### MySQL 

If you're using MySQL/MariaDB you'll need to set `DB_DRIVER`, `DB_USERNAME`, `DB_PASSWORD`, `DB_HOST`, `DB_PORT`, and `DB_NAME`.

`DB_DRIVER` *must* by `mysql`

`DB_USERNAME` is the name of a user in the database that can select, insert, and delete

`DB_PASSWORD` is the password of the above user

`DB_HOST` is the hostname or IP address of the database server

`DB_PORT` is the port used to connect to the database

`DB_NAME` is the name of the database you created

#### Site Info

You need to provide values for `SITE_URL` and `WEB_SERVER_PORT`.

`SITE_URL` is used to secrets links. It gets templated into `templates/create.html.tmpl`. If that value is incorrect, your links won't work.

`WEB_SERVER_PORT` defines what port the server listens on. Define this even if you're using a standard web server port.

You can set `LINK_EXPIRATION_TIME` to specify how long (in seconds) a link is valid for. If you don't specify a value, the app will default to 3600 seconds (one hour)

#### TLS

If you want to serve the site over TLS (and you really should), set `SERVE_TLS=true`.

Then, point `TLS_CERTIFICATE` and `TLS_KEY` to your public certificate and private key, respectively.

Per the `net/http` documentation:

```text
If the certificate is signed by a certificate authority, the certFile should be the concatenation of the server's certificate, any intermediates, and the CA's certificate.
```

### Building

Run `rice embed-go`

Run `go build`

## Creating Secrets Programmatically

`POST` a JSON to the `/create` endpoint. 

Example JSON: `{"secret": "hello world"}`

You'll receive a JSON response; grab the value of the `url` key

## Why

I use [One Time Secret](https://onetimesecret.com/) pretty frequently. Creating a self-destructing password sharing tool seemed like a fun problem to solve. It also has the added benefit of controlling the tool managing your secrets, not relying on a third party.

This tool is intended to be customized, forked, altered, and mangled. PRs welcome.
