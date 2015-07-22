## Khabar

Notifications engine.

It means

> the latest information; news.

> -- [google](https://www.google.com/search?q=define+khabar)

## Development

Clone the repo

```sh
$ go get github.com/codegangsta/gin
```

[gin](http://github.com/codegangsta/gin) is used to to automatically compile files while you are developing

Then run

```sh
$ DEBUG=* go get && go install && gin -p 8911 -i
```

Now you should be able to access the below API's on port `8911`.

MongoDB config is stored in `config/conf.go`.

## Usage

```sh
$ go get github.com/bulletind/khabar
$ khabar
```

## Concept and idea



**Channels**

**Ident** or **topic**

## How does it work?


## API

1. Add a new channel

  ```
  POST /topic/<notification_ident>/channel/<channel_ident>
  ```

  Request:
  ```js
  {
      org: "",
      app_name: "",
      user: "",
      ident: ""
  }
  ```


2. Remove a channel

  ```
  DELETE /topic/<notification_ident>/channel/<channel_ident>
  ```

  Request:
  ```js
  {
    org: "",
    app_name: "",
    user: "",
    ident: ""
  }
  ```

3. Remove a topic

  ```
  DELETE /topic/<notification_ident>
  ```

  Request:
  ```js
  {
      org: "",
      app_name: "",
      user: "",
      ident: ""
  }
  ```


4. Get all topics

  ```
  GET /topics
  ```

  Query params:

  - `org`: organization id
  - `user`: user id

  Response:
  ```js
  [
    {
      "_id": "",
      "created_on": 1425547531188,
      "updated_on": 1425879125700,
      "user": "",
      "org": "",
      "app_name": "",
      "channels": [
        "",
        ""
      ],
      "ident": ""
    }
  ]
  ```

5. Get all channels

  ```
  GET /channels
  ```

  Query params:

  - `org`: organization id
  - `user`: user id

  Response:
  ```js
  [
    {
        "_id": "",
        "created_on": 1425545240236,
        "updated_on": 1425545240236,
        "user": "",
        "org": "",
        "app_name": "",
        "data": {
        },
        "ident": ""
    }
  ]
  ```

6. Get all notifications

  ```
  GET /notifications
  ```

  Query params:

  - `user`: user id
  - `org`: organization id

  Response:
  ```js
  [
    {
      org: "",
      app_name: "",
      user: "",
      destination_uri: "",
      text: "",
      topic: "",
      destination_uri: "",
      is_read:false,
      created_on: <milliseconds_since_epoch>
     },
    { _id: 2,... }
  ]
  ```

  This can be polled periodically

7. Mark a single notification as read

  ```
  PUT /notification/:_id
  ```

8. Mark all unread notifications as read

  ```
  PUT /notifications
  ```

  - `destination_uri`: Link to relevant entity. (i.e action, incident)
  - `text`: Notification text
  - `topic`: Notification topic

Some general conventions:

- For all of the above request you must pass atleast one of the `org` or `user`
- For a successful request you would get a response status code of `201` if the entity is created or `204` if an entity is modified.
- Response for creation of an entity

  ```js
  {
    "body": "[ID of entity created]",
    "message": "Created",
    "status": 201
  }
  ```
- Response for modifying/deleting an entity

  ```js
  {
    "body": "",
    "message": "NoContent",
    "status": 204
  }
  ```

#### What does Khabar mean?

It literally means "message" or "news" in hindi language.

## Todo

- Break the library into multiple pieces like types, api and only keep the business logic in this repo
- Spin up a nice demo
- Add the admin for managing topics etc
- Ability to use `MONGODB_URL` from environment variable
- Ability to listen on a specified port (from command line?) or via `PORT` environment variable
