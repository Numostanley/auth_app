# Docker compose Installation

Clone the repository:

```
git clone https://github.com/Numostanley/auth_app.git
```

Enter the root directory.
```
cd auth_app
```

Run docker-compose

```
docker-compose up --build
```


# Manual Installation

Clone the repository:

```
git clone https://github.com/Numostanley/auth_app.git
```

Enter the root directory.
```
cd auth_app
```

Change the directory to the src folder:
```
cd src
```

Install go dependencies:
```
go mod tidy
```

Load go dependencies
```
go mod vendor
```

Change the PostgreSQL host settings in the env/dev/.env or env/prod/.env file to localhost
```
PG_HOST="localhost"  // or your database host
```

Running the Project
```
go run main.go
```

# NOTE:
Ensure your email settings in the env/dev/.env or env/prod/.env file are properly set up:
```
EMAIL_USER="example@gmail.com"  // change to your email user

EMAIL_PASSWORD="ExamplePassword@1"  // change to your email password

EMAIL_HOST="smtp.gmail.com"  // change to your email host
```

Also change the ISSUER settings in the env file to your application url or host
```
ISSUER="https://api.example.com/"  # change to your application url
```

Change the redirect_uris in the src/extras/clients.json file to your application uri
```
"redirect_uris": "https://api.example.com/v1/admin"
```


## Next:
Ensure you have PostgreSQL properly installed and configured if you are using the manual installation method.
<br/>
You can also find the Postman API documentation 
to this project at: <br/>
https://www.postman.com/orange-capsule-84916/workspace/auth-app
