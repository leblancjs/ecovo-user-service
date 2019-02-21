# User Service
## Introduction
The user service implements the user REST API. It makes it possible to access a user's details, such as it's profile, as well as create and update a user.

## To-Do
* Document errors codes/responses more cleanly
* Refactor handlers to reduce the amount of business logic they have
* Add validation to the user struct
* Add table of contents to README

## Configuration
The application's database connection and Auth0 domain are configured using environment variables. To avoid having to define them every time the service is run, they are kept in the `.env` file at the root of the repository.

The table below enumerates the different environment variables.

|Name|Required|Description|
|---|---|---|
|AUTH_DOMAIN|Yes|Domain where the user info endpoint is hosted (ex. my.domain.com)|
|DB_HOST|Yes|URI to where the database is hosted|
|DB_USERNAME|Yes|Username to use to to establish the database connection|
|DB_PASSWORD|Yes|Password to use to establish the database connection|
|DB_NAME|Yes|Name of the database to use on the server|
|DB_CONNECTION_TIMEOUT|No|Time to wait before giving up on connecting to the database|

## Build and Test
### Prerequisites
#### Docker
Docker is used to simplify the build and test processes. It makes it possible
to build and run the application without needing to install Go, and also makes it much easier to define environment variables to use to configure the service (see the next section).

Please download and install [Docker Desktop](https://www.docker.com/products/docker-desktop), and make sure that it is running on your machine before you proceed.

### Step 1 - Build an Image
In order to run the application locally to test it, we need to build an image
using Docker.

To do so, run following command in a terminal:

```
docker build --tag=user-service .
```

You will need to rebuild the image every time a change is made in the code, or when new changes are pulled.

Don't worry, it doesn't take that long.

### Step 2 - Run the Image in a Container
To run the service, we need to run the image we built in the previous step in a
container using Docker.

To do so, run the following command in a terminal and replace `<PORT>` with the port you want to use to access the API:

```
docker run -it -p <PORT>:8080 --env-file .env user-service
```

It is important to note that the `--env-file` argument is used to tell Docker
to define the environment variables found in the `.env` file in the Docker
container. Otherwise, the service will not start.

## Deploy
The service can be deployed to [Heroku](https://heroku.com) by pushing a Docker
image to its container registry, and releasing it in a Heroku application.

### Environment Variables
It is important to note that the service still needs those environment
variables! On Heroku, they need to be defined in the dashboard as Config Vars.
Without them, the service will fail to start.

### Prerequisites
The same prerequisites defined in the Build and Test section apply here.

#### Heroku CLI
The [Heroku CLI](https://devcenter.heroku.com/articles/heroku-cli#download-and-install)
is used to deploy the application to Heroku. Please download and install it on your machine.

##### Login
To log in to Heroku, enter the following command in a terminal:

```
heroku login
```

It should open a web browser in which you can log in using the Ecovo account credentials, which can be found on Google Drive.

This step only needs to be done once, after you've installed the Heroku CLI.

##### Link the Git Repository to the Heroku Application
To make sure that we deploy the service to the right application on Heroku, we
need to link the Git repository to the application.

To do so, run the following command in a terminal:

```
heroku create ecovo-user-service
```

This step only needs to be done once, after you've cloned the Git repository.

#### Step 1 - Push the Image to the Container Registry
To build and push the image to the Heroku container registry, use the following command:

```
heroku container:push web
```

#### Step 2 - Release the Container
To release the container that was pushed in the previous step, use the following command:

```
heroku container:release web
```

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