# User Service
>For the moment, this service stores users in memory. All data will be lost
>once it is shutdown.
>
>Harold will still manually reply to `GET /users/me` requests as fast as he
>can, though. Please be patient with him!
>
>![](https://hungarytoday.hu/wp-content/uploads/2018/02/18ps27.jpg)

## Introduction
The user service implements the user REST API. It makes it possible to access a user's details, such as it's profile, as well as create and update a user.

## To-Do
* Document errors codes
* Document how to deploy to Heroku
* Find out why we get 404s when deployed on Heroku

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
200 OK upon success

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
        "animals": "{0|1|2}",
        "conversation": "{0|1|2}",
        "music": "{0|1|2}"
    },
    "signUpPhase": "{0|1}"
}
```

### GET /users/{id}
#### URL Parameters
##### id
The user's unique identifier obtained from Auth0 when sign-in/sign-up is completed.

#### Request
##### Headers
```
Content-Type: application/json
```

#### Response
##### Status Code(s)
200 OK upon success

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
        "animals": "{0|1|2}",
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
        "animals": "{0|1|2}",
        "conversation": "{0|1|2}",
        "music": "{0|1|2}"
    }
}
```

#### Response
##### Status Code
201 CREATED upon success

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
        "animals": {0|1|2},
        "conversation": {0|1|2},
        "music": {0|1|2}
    },
    "signUpPhase": {0|1}
}
```

### PATCH /users/{id}
#### URL Parameters
##### id
The user's unique identifier obtained from Auth0 when sign-in/sign-up is completed.

#### Request
##### Headers
```
Content-Type: application/json
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
        "animals": {0|1|2},
        "conversation": {0|1|2},
        "music": {0|1|2}
    },
    "signUpPhase": {0|1}
}
```

#### Response
##### Status Code(s)
200 OK upon success

##### Headers
```
Content-Type: application/json
```