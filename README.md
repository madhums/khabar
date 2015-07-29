[![Build Status](https://travis-ci.org/bulletind/khabar.svg?branch=develop)](https://travis-ci.org/bulletind/khabar)
[![Gitter](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/bulletind/khabar?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge)

## Khabar

Notifications engine.

It means

> *the latest information; news.*

[google it](https://www.google.com/search?q=define+khabar)

## Table of contents

- [Concept and idea](#concept-and-idea)
- [How does it work?](#how-does-it-work)
- [Development](#development)
- [Usage](#usage)
- [API](#api)
  - **[Channels](#channels)**
    1. [Add a new channel](#add-a-new-channel)
    2. [Remove a channel](#remove-a-channel)
    3. [Get all channels](#get-all-channels)
  - **[Topics](#topics)**
    1. [Add a new topic](#add-a-new-topic)
    2. [Remove a topic](#remove-a-topic)
    3. [Get all topics](#get-all-topics)
  - **[Notifications](#notifications)**
    1. [Get all notifications](#get-all-notifications)
    2. [Mark a single notification as read](#mark-a-single-notification-as-read)
    3. [Mark all unread notifications as read](#mark-all-unread-notifications-as-read)
  - [Some general conventions](#some-general-conventions)
- [Todo and future plans](#todo-and-future-plans)

## Concept and idea

**Channels**

**Ident** or **topic**

## How does it work?

## Development

```sh
$ go get github.com/codegangsta/gin
$ go get github.com/tools/godep
```

- [gin](http://github.com/codegangsta/gin) is used to to automatically compile files while you are developing
- [godep](http://github.com/tools/godep) is used to manage dependencies

Then run

```sh
$ mkdir -p $GOPATH/src/github.com/bulletind
$ cd $GOPATH/src/github.com/bulletind
$ git clone https://github.com/bulletind/khabar.git # or your fork
$ cd khabar
$ DEBUG=* go get && go install && gin -p 8911 -i
```

Now you should be able to access the below API's on port `8911`.

MongoDB config is stored in `config/conf.go`.

After you make the changes (if you import any new deps), don't forget to run

```sh
$ godep save ./...
```

## Usage

```sh
$ go get github.com/bulletind/khabar
$ khabar
```

## API

1. #### Channels

  1. ##### Add a new channel

    ```
    POST /topic/<notification_ident>/channel/<channel_ident>
    ```

    Request:
    ```js
    {
        org: "",
        user: "",
        ident: ""
    }
    ```

  2. ##### Remove a channel

    ```
    DELETE /topic/<notification_ident>/channel/<channel_ident>
    ```

    Request:
    ```js
    {
      org: "",
      user: "",
      ident: ""
    }
    ```

  3. ##### Get all channels

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

2. #### Topics

  1. ##### Add a new topic

    ```
    POST /topic
    ```

    Request:
    ```js
    {
        org: "123",
        user: "",
        ident: "",
        app_name: "myapp"
    }
    ```

    - `ident` is requred
    - `org` is required
    - `app_name` is required

  2. ##### Remove a topic

    ```
    DELETE /topic/<notification_ident>
    ```

    Request:
    ```js
    {
        org: "",
        user: "",
        ident: ""
    }
    ```

  3. ##### Get all topics

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

3. #### Notifications

  1. ##### Get all notifications

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
        user: "",
        destination_uri: "",
        text: "",
        topic: "",
        destination_uri: "",
        is_read:false,
        created_on: <milliseconds_since_epoch>
       },
       // and so on...
    ]
    ```

    This can be polled periodically

  2. ##### Mark a single notification as read

    ```
    PUT /notification/:_id
    ```

  3. ##### Mark all unread notifications as read

    ```
    PUT /notifications
    ```

    - `destination_uri`: Link to relevant entity. #### (i.e action, incident)
    - `text`: Notification text
    - `topic`: Notification topic

#### Some general conventions:

- For all of the above request you must pass atleast one of the `org` or `user` or both
- For all the listings, you get a status code of `200`
- When you create a resource you get a status code of `201`
- When you modify/delete a resource you get a status code of `204`
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

## Todo

- Verify if the api's listed above are correct and update them.
- Spin up a nice demo.
- Ability to use `MONGODB_URL` from environment variable.
- Ability to listen on a specified port (`--port` from command line?) or via `PORT` environment variable.
- Deployment and hosting.
