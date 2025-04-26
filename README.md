# bd-07-gp-chirpy

Guided project:
> In this course, we'll be working on a product called "Chirpy". Chirpy is a social network similar to Twitter.

# chirpy

This console application implements a server with API endpoints to manage users and messages (called "chirps") with authentication and webhook functionality. User, message and authentication data is stored locally on a PostgreSQL database.

## Installation

Clone this repository with:

```bash
git clone https://github.com/danilogalisteu/bd-07-gp-chirpy.git
```

The application uses the Go language and a local PostgreSQL database.
Following the [instructions here](https://go.dev/doc/install), install the specific Go language runtime for your platform.
Following the [instructions here](https://www.postgresql.org/download/), install a recent version of the specific PostgreSQL server for your platform.

To make a self-contained executable file called `chirpy`, run:

```bash
go build -o chirpy
```

To install the application on your system, to be used from any path as the command `chirpy`, run:

```bash
go install -o chirpy
```

The required external dependencies should be downloaded and installed automatically on the first run or build.

## Configuration

### Database

The local PostgreSQL server should be run as a service, under a specific user such as `postgres`, with a password.
The Linux commands to set the user password and access the server shell are:
```bash
sudo passwd postgres
sudo service postgresql start
sudo -u postgres psql
```

In the psql shell, under the chosen user, you should create a database named `chirpy`, usig the SQL query
```sql
CREATE DATABASE chirpy;
```
and then connect to the database using the command `\c chirpy`.

Finally, the following command sets the database user and password for connection:
```sql
ALTER USER postgres PASSWORD 'postgres';
```
The database user (referred to as PG_USER) should be the same as the system user, and the database password (referred to as PG_PASS) should be different from the system password.

### Application

The application expects specific environment variables that could alternatively be defined in a text file named `.env` in your current directory.
The variables should be defined as follows:
```bash
DB_URL="postgres://<PG_USER>:<PG_PASS>@localhost:5432/chirpy?sslmode=disable"
JWT_SECRET="<random 64-character string>"
POLKA_KEY="<API key from payment service>"
```

## Use

After installing the requirements, start the application with:

```bash
go run .
```

or run the built executable `chirpy` directly.
