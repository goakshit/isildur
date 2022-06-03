#### Isildur
The subscription service exposes api to fetch products, subscriptions and create subscriptions.

#### Steps to run the code:
1. Install docker.
2. Run `docker-compose --env-file ./.env -f ./build/docker/docker-compose.yaml up --build`.
3. You can use Postman collection for quickly testing the APIs.

#### Whats covered:
- User story 1 - Done
- Optional story 1 - NA
- Optional story 2 - NA

#### Tests (Unit):
I have tried to add some tests, but there can be a lot more. Lot of edge cases in the story.
Run `go test ./...`

#### Areas of improvement:
1. Logging: We can use a good logging library like zerolog to log at various level. Maybe we can log function entry and exit as well.
2. Handling env/config: We can use powerful libraries like envconfig or viper for better management.
3. Docs: I am just adding comments for documentation which can be seen using `go doc`, but we can use open api specification and build and serve docs locally or in production as well.
4. DB design can be cleaner.