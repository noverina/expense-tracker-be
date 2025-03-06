Backend for expense tracking application.

# How to build and run locally

- Generate private key\
  ` openssl genpkey -algorithm RSA -out private_key.pem -pkeyopt rsa_keygen_bits:2048`\
- Generate public key\
  `openssl pkey -in private_key.pem -pubout -out public_key.pem`

- Make an `.env` in the root of this project with these values
<pre><code>export MONGODB_URI=connection string for mongodb database
export MONGODB_DB=database name
export MONGODB_COLL=collection name
export CONNECT_TIMEOUT=timeout for connecting to db
export FRONTEND_ADDRESS=address of frontend application
export MAX_EVENT_COUNT=the max amount of events in a given day
export PRIVATE_KEY=file name of private key (path shouldnt be included, file name only)
export PUBLIC_KEY=file name of public key (path shouldnt be included, file name only)
export ACCESS_EXPIRY=(in hour) how long should access token last
export REFRESH_EXPIRY=(in hour) how long should refresh token last</code></pre>

Then simply `go run .`\
\
This application was build using `go 1.23.6`.\
Minimum `go 1.23.0` needed.

# Database information
