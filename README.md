# To run locally:

- clone repository
- ensure that there are no services running on localhost:8080 or on port 5432(check postgres)
- in terminal run docker-compose up -d . This will create two containers.
  1. postgres is a locally running instance of a postgres database.
  2. loanpro_api is a container that will run the GoLang rest api.
- in terminal run $ curl -X POST http://localhost:8080/api/v1/operation/seedthedb to run a basic seeding function for loading the operations.
- you can now use the api through the frontend container (see readme on loanprofe)
