# gs-onboarding

[Architecture Diagram](https://lucid.app/lucidchart/invitations/accept/inv_b0ee37ba-6710-41c4-be4e-b4992c2d6e0d)

## Services

### Consumer

The consumer service periodically fetches the top stories from hacker news and stores all non dead nor deleted items in the database

### API

The API service is a gRPC server that offers a interface to fetched the stored hacker news stories

### Gateway

The gateway service is main entry point for third parties to access all other systems. Currently, it is responsible for proxying requests to the API service
