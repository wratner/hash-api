# hash-api
[![Build Status](https://travis-ci.org/wratner/hash-api.svg?branch=master)](https://travis-ci.org/wratner/hash-api)
[![Coverage Status](https://coveralls.io/repos/github/wratner/hash-api/badge.svg?branch=master)](https://coveralls.io/github/wratner/hash-api?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/wratner/hash-api)](https://goreportcard.com/report/github.com/wratner/hash-api)

# Purpose

This project allows a user to hash and encode a password over http. It also allows the user to gracefully shutdown the server and to obtain statistics about the requests that have been made to the server. 

# APIs

* `/hash`
  * POST
  * The 'hash' endpoint takes in a password as a form value and returns a base64 encoded SHA512 string of that password.
    * cURL examples:
      * `curl -d "password=angryMonkey" http://localhost:8080/hash`
      * `curl -X POST http://localhost:8080/hash?password=angryMonkey`
        * Response: `ZEHhWB65gUlzdVwtDQArEyx+KVLzp/aTaRaPlBzYRIFj6vjFdqEb0Q5B8zVKCZ0vKbZPZklJz0Fd7su2A+gf7Q==` 
* `/stats`
  * GET
  * The 'stats' endpoint returns a JSON response in format ```{"total":0,"average":0}```
   * total is the total number of requests made to the `/hash` endpoint during the life of the server.
   * average is the time it has taken to process all of those requests to the `/hash` endpoint in milliseconds.
     * This includes the 5 second delayed response. Because of this, nearly all responses are ~5001-5003 milliseconds. 
* `/shutdown`
  * GET
  * The 'shutdown' endpoint will gracefully shutdown the server. All remaining requests are allowed to complete and no additional requests will be allowed to be processed. The server will stop as soon as any existing work is completed.
  
## DISCLAIMER

This application is using SHA512 as the algorithm to hash the passwords. This is extremely unsafe and should not be used in production. This was used due to a requirement of the assignment in which I was only allowed to use the standard library. Should you not have this restriction, I highly recommened that you use [bcrypt](https://godoc.org/golang.org/x/crypto/bcrypt). If you are for some reason restricted to the Go standard library then I would suggest using HMAC-SHA-512/256. While not ideal for password hashing, it at least protects against length-extension attacks which SHA512 is vulnerable to. 
