# CRM Backend Project for Udacity Golang Course

This is my submission for the CRM Backend project,
which is the final part of the Udacity course on the
Go language (Golang).

The project provides a simple API to interact with a
mock database of customer data on a server. The server
can be started via `go run main.go` and accessed at
`localhost:3000`.

What follows is a description of how to use the API.
There are five basic tasks that can be performed:

- get the entire database of all customers
- get the data for a single customer
- create a new customer in the database
- update a customer's information in the database
- delete a customer from the database

Each task can be performed with API requests as
described below. All API requests accept JSON as input
and respond with JSON as output. These requests can be
made using various GUI programs such as
[Postman](https://www.postman.com/) or using simple
command-line utilities such as `curl`. Note that the
examples below are provided as `curl` commands.

We begin with a brief description of the structure of
an individual customer record.

## Customer Record Structure

Each customer record consists of six fields:

- id, a string of nine digits that must be
  unique for each customer
- name, a string of the customer's name
- role, a string of the customer's role
- email, a string of the customer's email
- phone, an integer of the customer's phone number
> Note: phone is an integer because it is required by
> the tests provided with the project. However, phone
> would perhaps be better represented as a string,
> since a phone "number" is not actually numeric, and
> since phone numbers with leading zeros (unlikely in
> the United States, but perhaps more common
> elsewhere) could be problematic to represent as
> integers.
- contacted, a boolean (true/false) value of whether
  the customer has been contacted

Note that missing values for the fields are handled
differently according to the fields' types:

- A missing value for each of the string fields
  **except id** is represented by the `""`, the empty
  string.
- A missing value for the phone field is represented
  by `0`.
- A missing value for the contacted field is
  represented by `false`.

Note that the id field cannot be missing, as will be
discussed further below.

In JSON data sent to or received from the server, the
customer is represented as a JSON object with the
field names as keys.

## Errors

If an erroneous request to the server is made, it will
respond with an appropriate HTTP status code and the
`null` JSON value. More information about specific
status codes is discussed below for each request.

## Getting All Customers

To get all customers, simply make a `GET` request to
the `/customers` path, e.g.:
```
curl -X GET localhost:3000/customers
```
The server will respond with the `200: OK` status code
and a JSON object of the entire customer database. The
object contains each customer represented through a
key of the customer's id mapped to a value of an
object containing the customer's information
(including the id).

## Getting a Single Customer

To get a single customer, make a `GET` request to the
path `/customers/{id}`, where `{id}` should be
replaced with the id of the customer desired, e.g.:
```
curl -X GET localhost:3000/customers/023004163
```
The server will respond with the `200: OK` status code
and a JSON object representation of the customer.

If the id provided in the path does not correspond to
any customer in the database, the server will instead
respond with the `404: Not Found` status code.

## Creating a New Customer

To add a customer, make a `POST` request to the path
`/customers` and send JSON data containing the
information for the new customer, e.g.:
```
curl -X POST -d '{
	"id": "923392655",
	"name": "Ed Regis",
	"role": "Publicist",
	"email": "eregis@ingen.com",
	"phone": 2126364756,
	"contacted": false
}' localhost:3000/customers
```
Any of the six fields can be missing from the data
provided. For every field **except id**, a missing
field for a newly created customer will be replaced
with the type-appropriate missing value as described
above.

However, if the id field is missing, a random id will
be automatically generated for the newly created
customer.
> Note: This behavior was incorporated to satisfy the
> constraints of both the project rubric and the
> provided tests. The project rubric states that each
> customer `struct` must have an id field. However,
> the provided tests attempt to create a new customer
> using data that is missing the id field.

If the request is successful, the server will respond
with the `201: Created` status code and a JSON object
of the entire customer database, including the newly
created customer.

If the data for the new customer is malformed (as
opposed to merely having missing fields), the server
will respond with the `400: Bad Request` status code.

If an id is supplied in the the data that is already
in use for another customer, the server will respond
with the `409: Conflict` status code. The server will
also respond with this status code in the extremely
unlikely event that 1000 attempts to randomly generate
an id fail because they all conflict with an existing
customer id. (This is effectively impossible unless
the number of customers in the database has begun to
approach one billion.)

## Updating a Customer

To update a customer's information, make a `PUT`
request to the path `/customers/{id}`, where `{id}`
should be replaced with the id of the customer
desired, and send JSON data containing the customer's
updated information, e.g.:
```
curl -X PUT -d '{
	"id": "023004163",
	"name": "Alan Grant",
	"role": "Paleontologist",
	"email": "agrant@montana.edu",
	"phone": 4068672096,
	"contacted" true
}' localhost:3000/customers/023004163
```
Note that the customer's entire information (**except
id**, as described below) must be supplied in
the data, **including fields that you do not wish to
change**. Fields that are missing from the updated
data will be set to the type-appropriate missing
value, just as they would be if they were missing when
creating a new customer.

If the id field is missing from the updated customer
data, the customer's id will be unchanged.
> Note: The purpose of this behavior is to be
> consistent with the behavior when creating a new
> customer in which the id field is allowed to be
> missing without setting it to the empty string.

To change a customer's id, simply provide a new id for
the cusomter in the data. Note that the final portion
of the **path** must still be the **current** id. For
example, in the update request shown above, the
customer's id could also be changed from 023004163
to 470318456 (along with the other changes) as follows:
```
curl -X PUT -d '{
	"id": "470318456",
	"name": "Alan Grant",
	"role": "Paleontologist",
	"email": "agrant@montana.edu",
	"phone": 4068672096,
	"contacted" true
}' localhost:3000/customers/023004163
```
After making this change, any subsequent requests for
this customer would be made using the new path
`/customers/470318456`.

If the update of the customer is successful, the
server will respond with the `200: OK` status code and
a JSON object representation of the newly updated
customer.

If the id provided in the path does not correspond to
any customer in the database, the server will respond
with the `404: Not Found` status code.

If the updated data for the customer is malformed, the
server will respond with the `400: Bad Request` status
code.

If an id is supplied in the the data that is already
in use for another customer, the server will respond
with the `409: Conflict` status code.

## Deleting a Customer

To delete a customer from the database, simply make a
`DELETE` request to the path `/customers/{id}`, where
`{id}` should be replaced with the id of the customer
desired, e.g.:
```
curl -X DELETE localhost:3000/customers/023004163
```
The server will respond with the `200: OK` status code
and a JSON object representation of the entire
customer database, which no longer contains the
deleted customer.

If the id provided in the path does not correspond to
any customer in the database, the server will instead
respond with the `404: Not Found` status code.
