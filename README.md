Backend for expense tracking application.

# How to build and run locally

Generate keys for signing JWT:

- Generate private key\
  ` openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048`\
- Generate public key\
  `openssl pkey -in private_key.pem -pubout -out public_key.pem`

- Make an `.env` at the root of this project. Change the one with -- with appropriate values:
<pre><code>export MONGODB_URI=--connection string for mongodb database
export MONGODB_DB=expense
export EVENT_COLL=expense
export LOG_COLL=log
export AUTH_COLL=auth
export CONNECT_TIMEOUT=5 #timeout for connecting to db (in seconds)
export FRONTEND_ADDRESS=--address of frontend application (for cors config)
export MAX_EVENT_COUNT=10 #the max amount of events in a given day
export MAX_MONTH_RANGE=12 #how far back/advance in months can you go
export PRIVATE_KEY=--private key file
export PUBLIC_KEY=--public key file
export TOKEN_EXPIRY=86400 #(in second) how long token should last</code></pre>

(Optional) import the sample data with

<pre><code>mongoimport --db=expense --collection=expense --username=user --password=pass --authenticationDatabase=authDB --file=/data/backup/expense_coll.json</code></pre>

Then simply `go run .`\
\
This application was build using `go 1.23.6`.\
Minimum `go 1.23.0` needed.
