docker run --name my-postgres-container -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=mypassword -e POSTGRES_DB=postgres -p 5432:5432 -d postgres:latest

#TODO

1. Change the Image and name of ap when running the server
2. Prettify the CLI
3. Connect postgres using cli
4. Implement Go routine for concurrency