# bd-07-gp-chirpy

Guided project:
> In this course, we'll be working on a product called "Chirpy". Chirpy is a social network similar to Twitter.

# chirpy

This console application implements a server with API endpoints to manage users and messages (called "chirps") with authentication and webhook functionality. User, message and authentication data is stored locally on a JSON file.

## Installation and use

Clone this repository with:

```bash
git clone https://github.com/danilogalisteu/bd-07-gp-chirpy.git
```

The application uses the Go language. After [installing](https://go.dev/doc/install) the specific runtime for your platform, start the application with:

```bash
go run .
```

To make a self-contained executable file called `chirpy`, run:

```bash
go build -o chirpy
```

The required external dependencies should be downloaded and installed automatically on the first run or build.
