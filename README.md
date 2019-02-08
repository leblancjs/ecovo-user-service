# User Service
>For the moment, this service is just a mock.
>
>Harold will manually reply to each request as
>fast as he can. Please be patient with him!
>
>![](https://hungarytoday.hu/wp-content/uploads/2018/02/18ps27.jpg)

## Introduction
The user service implements the user REST API. It makes it possible to access a user's details, such as it's profile, as well as create and update a user.

## To-Do
* Document how to deploy to Heroku
* Find out why we get 404s when deployed on Heroku
* Put this bad boy in a Docker container

## Build and Test
To build and test the service, use the following command in the terminal (Linux and macOS) or the command prompt (Windows):

`go run main.go`

The service will listen for requests on port 8080.

## Deploy
TODO: Fill out this section.

## Endpoints
### GET /users/me
TODO: Fill out this section.

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
    "id": "{id}", // MANDATORY
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