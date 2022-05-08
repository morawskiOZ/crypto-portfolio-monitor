# Camp Me Go API GQL
Graphql API for Camp Me App

### Technologies
#### Prisma Go Client
ORM client for postgres DB
#### gqlgen
Builds gql server

### Setup locally
1. Project needs `.env` in root folder and in `./internal/prisma`. Copy `.env.examples` in both respective folders as `.env` files and fill all variables.

2. When you first start or you have made changes to your model, migrate your database and re-generate your prisma and gql code. From the root run:

    `make prepare`

It will prepare run migration and generate graphql server code.

3. Run gql server

    `make dev`


Resolver example
```
func (r *queryResolver) GetUserVote(ctx context.Context, input api.PriceVoteCheckInput) (bool, error) {
	a , err := r.client.PriceVote.Query().Where(pricevote.And(pricevote.HasUserWith(user.IDEQ(input.UserID)), pricevote.HasPriceWith(price.IDEQ(input.VoteID)))).Only(ctx)
}
```