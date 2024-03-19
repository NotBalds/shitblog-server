# API Reference

- ```POST /api/user/<username>``` - create user with given username if not exists and return token
Responses:
200 - user created
406 - username is already taken
400 - СРАНЫЙ ХАКЕР ПОШЁЛ НАФИГ

- ```POST /api/post/``` - create post and return it id, params:
```json
{
    "token": "user token",
    "title": "post title",
    "text": "post text"
}
```
Responses:
200 - post created
400 - incorrect request (not all fields) or no such user or incorrect token

- ```DELETE /api/user/<username>``` - delete user, params:
```json
{
    "token": "user token"
}
```
Responses:
200 - user deleted
400 - incorrect request (not all fields) or no such user
401 - incorrect token or ПОШЕЛ НАФИГ СРАНЫЙ ХАКЕР

- ```DELETE /api/post/<id>``` - delete post, params:
```json
{
    "token": "creator token"
}
```
Responses:
200 - post deleted or no such post or not deleted because you are not owner
400 - incorrect request (not all fields)
401 - incorrect token

- ```GET /api/user/``` - get all users (usernames)

- ```GET /api/post/``` - get posts, params:
```json
{
    "count": "how many posts, required param",
    "author": "Author of posts, optional param"
}
```

- ```GET /api/post/<id>``` - get post by id
