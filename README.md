# User Service
## Introduction
The user service implements the user REST API. It makes it possible to access a user's details, such as it's profile, as well as create and update a user.

## To-Do
* Document errors codes/responses more cleanly
* Refactor handlers to reduce the amount of business logic they have
* Add validation to the user struct
* Refactor Auth to make it configurable via environment variables (or config file)
* Refactor main to make DB configuratble via environment variables (or config file)
* Fix bug when deploying to Heroku
* Document how to deploy to Heroku

## Build and Test
### Docker
To build the service in a Docker container, use the following command in the terminal (Linux and macOS):

```
docker build --tag=user-service .
```

To run the service, use the following command and replace `<PORT>` with the port number to open on the container:

```
docker run -p <PORT>:8080 user-service
```

### Manually
To build and test the service, use the following command in the terminal (Linux and macOS) or the command prompt (Windows):

```
go run main.go
```

The service will listen for requests on port 8080.

## Deploy
TODO: Fill out this section.

## Endpoints
### GET /users/me
#### Request
##### Headers
```
Authorization: Bearer {access_token}
```

#### Response
##### Status Code(s)
* 200 OK upon success
* 401 Not Authorized if no token is present or it is invalid
* 404 Not Found if no user exists for the given authorization

##### Headers
```
Content-Type: application/json
```

##### Body
```
{
    "id": "{id}",
    "email": "{email}",
    "firstName": "{firstName}",
    "lastName": "{lastName",
    "dateOfBirth": "{timestamp}",
    "phoneNumber": "{phoneNumber}",
    "gender": "{Male|Female}",
    "photo": "{photoUrl}",
    "description": "{description}",
    "preferences": {
        "smoking": "{0|1|2}",
        "conversation": "{0|1|2}",
        "music": "{0|1|2}"
    },
    "signUpPhase": "{0|1}"
}
```

### GET /users/{id}
#### URL Parameters
##### id
The user's unique identifier generated when it is created.

#### Request
##### Headers
```
Content-Type: application/json
Authorization: Bearer {access_token}
```

#### Response
##### Status Code(s)
* 200 OK upon success
* 401 Not Authorized if no token is present or it is invalid
* 404 Not Found if no user exists for the given ID

##### Headers
```
Content-Type: application/json
```

##### Body
```
{
    "id": "{id}",
    "email": "{email}",
    "firstName": "{firstName}",
    "lastName": "{lastName",
    "dateOfBirth": "{timestamp}",
    "phoneNumber": "{phoneNumber}",
    "gender": "{Male|Female}",
    "photo": "{photoUrl}",
    "description": "{description}",
    "preferences": {
        "smoking": "{0|1|2}",
        "conversation": "{0|1|2}",
        "music": "{0|1|2}"
    },
    "signUpPhase": "{0|1}"
}
```

### POST /users
#### Request
##### Headers
```
Content-Type: application/json
Authorization: Bearer {access_token}
```

##### Body
The following example shows all the fields that can be included:
```
{
    "email": "{email}",
    "firstName": "{firstName}",
    "lastName": "{lastName",
    "dateOfBirth": "{timestamp}",
    "phoneNumber": "{phoneNumber}",
    "gender": "{Male|Female}",
    "photo": "{photoUrl}",
    "description": "{description}",
    "preferences": {
        "smoking": "{0|1|2}",
        "conversation": "{0|1|2}",
        "music": "{0|1|2}"
    }
}
```

#### Response
##### Status Code
* 201 CREATED upon success
* 400 Bad Request if the payload is malformated
* 401 Not Authorized if no token is present or it is invalid
* 500 Internal Server Error if something else goes wrong

##### Headers
```
Content-Type: application/json
```

##### Body
```
{
    "email": "{email}",
    "firstName": "{firstName}",
    "lastName": "{lastName}",
    "dateOfBirth": "{dateOfBirth}",
    "phoneNumber": "{phoneNumber}",
    "gender": "{gender}",
    "photo": "{photo}",
    "description": "{description}",
    "preferences": {
        "smoking": {0|1|2},
        "conversation": {0|1|2},
        "music": {0|1|2}
    },
    "signUpPhase": {0|1}
}
```

### PATCH /users/{id}
#### URL Parameters
##### id
The user's unique identifier generated when it is created.

#### Request
##### Headers
```
Content-Type: application/json
Authorization: Bearer {access_token}
```

##### Body
The following example shows all the fields that can be modified:
```
{
    "firstName": "{firstName}",
    "lastName": "{lastName}",
    "dateOfBirth": "{dateOfBirth}",
    "phoneNumber": "{phoneNumber}",
    "gender": "{gender}",
    "photo": "{photo}",
    "description": "{description}",
    "preferences": {
        "smoking": {0|1|2},
        "conversation": {0|1|2},
        "music": {0|1|2}
    },
    "signUpPhase": {0|1}
}
```

#### Response
##### Status Code(s)
* 200 OK upon success
* 400 Bad Request if the payload is malformated
* 401 Not Authorized if no token is present or it is invalid
* 404 Not Found if no user exists for the given ID

##### Headers
```
Content-Type: application/json
```