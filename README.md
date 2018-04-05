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
 
