## Backend written in Go utilizing rest principles for API

### Features
- Intergrated Twilio Api for SMS based authentication.
- Implemented secure data transfer through JWT creation and signing with 256-bit Hashing Algorithm.
- Inclueded CORS functionallity.
- Created access to a Mongo db database using the Mgo driver.

### API breakdown 
Domain: http://localhost:12345
#### Endpoint: /api/signup
#### Description:
1. Get firstName, lastName, studentID, phoneNumber from request body
2. Assign the Student a uniqueID 
3. Check they are not already in the database, if they are return an error
4. Else save the student to the DB and return the JSON object.

Documentation:
#### `POST /api/signup` which accepts the following data structure in JSON:

```
{
"first_name": "John",
"last_name":  "Doe",
"phone_no": "0720301929301",
"student_id": "acfb463"
}
```
and returns the following:
```
{
"first_name": "John",
"last_name":  "Doe",
"phone_no": "0720301929301",
"student_id": "acfb463",
"unique_id": "f1ed1046-738f-45c4-aa94-090fbcc2e6f5"
}
```
#### Endpoint: /api/authenticate
#### Description:
1. Authenticate the user from the DB
2. Create a JWT that takes in the students credetials as the payload
3. generate a 6 digit code to send via twilio for SMS basd authentication
4. Return the token
Documentation:
#### `POST /api/authenticate` which accepts the following data structure in JSON:

```
{
"first_name": "John",
"phone_no": "0720301929301",
"student_id": "acfb463"
}
```

and returns the following:
```
{
"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiSm9obiIsInN0dWRlbnRfaWQiOiJhY2ZiNDYzIiwicGhvbmVfbm8iOiIwNzIwMzAxOTI5MzAxIiwiaXNzIjoiZ29sYW5nIGFwaSIsImV4cCI6MTUxNjIzOTAyMn0.0y5110nIWnv5OhwXuKnHnhgK4AIBI52y7TJuDbNw-eg" 
}
```
#### Endpoint: /api/phonecode
#### Description:
1. Decode the JWT with the passcode the user has sent from the client
2. If successful send decoded student id and unique id
3. Else send and Authorization error message
Documentation:
#### `POST /api/phonecode` which accepts the following data structure in JSON:

```
{
"phoneCode": "8382010"
}
```
#### with headers:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmaXJzdF9uYW1lIjoiSm9obiIsInN0dWRlbnRfaWQiOiJhY2ZiNDYzIiwicGhvbmVfbm8iOiIwNzIwMzAxOTI5MzAxIiwiaXNzIjoiZ29sYW5nIGFwaSIsImV4cCI6MTUxNjIzOTAyMn0.0y5110nIWnv5OhwXuKnHnhgK4AIBI52y7TJuDbNw-eg
```
and returns the following:
```
{
"student_id": "acfb463",
"unique_id": "f1ed1046-738f-45c4-aa94-090fbcc2e6f5"
}
```


