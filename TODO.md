# Sober-API

- [x] creation of tables

-[] OnBoardingFLow:
        -[] User can add a reason and date of start for being sober
        -[] User can update the reason and date of start for being sober
        -[] update the userId with the reason and date
        -[] User can see a list of reasons for being sober

``` json {
  "userId":13,
  "sobriety":{
  "reason":"hello",
  "soberDate":"2023/10/10"
}}
```
```

version: '3.8'

services:
  psql:
    image: postgres:latest
    environment:
      POSTGRES_DB: ${DB_DATABASE}
      POSTGRES_USER: ${DB_USERNAME}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    ports:
      - "${DB_PORT}:5432"
    volumes:
      - psql_volume:/var/lib/postgresql/data

volumes:

  psql_volume:

- Default audience
  https://iam.googleapis.com/projects/868110718852/locations/global/workloadIdentityPools/sober-api/providers/github-actions
assertion.repository=="Fordjour12/sober-api"

- service account 
sober-go-api@com-thephantomdev.iam.gserviceaccount.com

CLE
assertion.repository=="Fordjour12/sober-api"
