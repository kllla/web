# kll.la/web

This is a simple suite of web services for running an invitation only posting community with 
shared services among users. All hosted on Google Cloud Platform using Firebase and Storage Buckets. 
Authentication his handled in app with bcrypt salted passwords while Authentication is via a bearer session token.
All the html rendered server side. This is to have 0 reliance on javascript so to be able to be run over TOR with no client side javascript.

- multi user 
- bearer auth sessions
- public / private posts
- invite system 
- link shrinker 

## prerequisites

- google GCP account

### env vars

these are required to successfully start ther service.
 - BUCKET_NAME
 - SERVICE_ACCOUNT_ID
 - PROJECT_ID 