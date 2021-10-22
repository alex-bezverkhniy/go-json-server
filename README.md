# go-json-server
Dummy simple full fake REST API with zero coding.


## Project status

- [x] Load DB frmo `.json` file 
- [x] Generates endpoints by high level json properties
from JSON:
```json

{
    "users": [
        {
            "id": 0,
            "name": "Bob"
        },
        {
            "id": 1,
            "name": "Joe"
        }
    ]
}
```
to:
```
GET    /users
GET    /users/1
POST   /users
```
- [ ] CLI usage. Use [Cobra CLI](https://cobra.dev).
```
go-json-server [options] <source>
```
- [ ] CRUD operatons
```
GET    /posts
GET    /posts/1
POST   /posts
PUT    /posts/1
PATCH  /posts/1
DELETE /posts/1
```
- [ ] Pagination. Use `_page` and optionally `_limit` to paginate returned data. 
```
GET /posts?_page=3&_limit=10
```
- [ ] Full-text search. Add q 
```
GET /posts?q=[programming]
```
- [ ] Filtering. Use . to access deep properties 
```
GET /posts?author.name=<string>
```
- [ ] Swagger UI. Use `/swagger` endpoint to open [Swagger UI](https://swagger.io/tools/swagger-ui/) page
- [ ] Static file server. Use `./public` directory to serve your HTML, JS and CS. 
- [ ] Use `--static` to set a different static files directory.