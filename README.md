# OOPS

OOPS One-time Password Sharing is a Go tool written to provide short-lived (one hour), one-time access to a secret.

## Usage

### Database

You can either use MySQL/MariaDB or SQLite.

If you're using MySQL, you'll need to create the database first:

```
create database oops;

grant all on oops.* to 'oopsuser'@'%' identified by 'oopsPASS';
```

Careful, `%` allows the user to connect to the database from anywhere. Use proper, layered access control to the database.

### .env File

Set the `OOPS_ENV_FILE` environment variable to the path of your config file for the application (this can be anywhere). You can reference the `.env.example` file for examples for mysql and sqlite.

## Why

I use [One Time Secret](https://onetimesecret.com/) pretty frequently. Creating a self-destructing password sharing tool seemed like a fun problem to solve. It also has the added benefit of controlling the tool managing your secrets, not relying on a third party.

This tool is intended to be customized, forked, altered, and mangled. PRs welcome.