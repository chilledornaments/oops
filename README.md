# OOPS

OOPS One-time Password Sharing is a Go tool written to provide short-lived, one-time access to a secret.

## Usage

### Database

You can either use DynamoDB or SQLite

### .env File

Set the `OOPS_ENV_FILE` environment variable to the path of your config file for the application (this can be anywhere). You can reference the `.env.example` file for examples for DynamoDB and SQLite.

#### Database Connection

##### SQLITE3

If you're using SQLITE3 (recommended for development and low traffic sites), you only need to set `DB_DRIVER` and `DB_PATH`.

`DB_DRIVER` *must* by `sqlite3`

`DB_PATH` is the path to where you want the SQLite database to live

##### DynamoDB

You can use AWS DynamoDB as a datastore for OOPS. 

`DB_DRIVER` *must* be `dynamo`

`DYNAMO_TABLE_NAME` should be set to the name of the table to store records in

When using DynamoDB, the application doesn't check for secret expiration. Instead, use DynamoDB TTL. It's important to note that DynamoDB TTL's don't always expire a secret exactly at expiration time. 

There is sample Terraform code to stand up the DynamoDB table along with an IAM role (for an EC2 instance) and an IAM policy to allow the role to interact with DynamoDB.


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

Run `rice -v -i ./cmd embed-go`

Run `go build`

## Creating Secrets Programmatically

`POST` a JSON to the `/create` endpoint. 

Example JSON: `{"secret": "hello world"}`

You'll receive a JSON response; grab the value of the `url` key

## Why

I use [One Time Secret](https://onetimesecret.com/) pretty frequently. Creating a self-destructing password sharing tool seemed like a fun problem to solve. It also has the added benefit of controlling the tool managing your secrets, not relying on a third party.

This tool is intended to be customized, forked, altered, and mangled. PRs welcome.
