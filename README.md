# khabar
=======
Notifications engine

### API

1. Adding a new channel to notification setting

   **Method**: `POST`
   
  **EndPoint**: `/topic/<notification_ident>/channel/<channel_ident>`
  
  **Request body** :
  ```js
 {
      org:"",
      app_name: "",
      user:"",
      ident:""
  }
  ``` 
    **Response Code**: `200`
  
2. Removing a channel from notification setting

   **Method**: `DELETE`

  **EndPoint**: `/topic/<notification_ident>/channel/<channel_ident>`
  
  **Request body** :
  
  ```js
 {
      org:"",
      app_name: "",
      user:"",
      ident:""
  }
  ``` 
  **Response Code**: `200`

3. Removing  notification setting

  **Method** : `DELETE`
 
  **EndPoint**: `/topic/<notification_ident>`
 
  **Request Body**:
 
```js
 {
      org:"",
      app_name: "",
      user:"",
      ident:""
 }
```
  **Response Code**: `200`

4. Get all notification settings

 **Method** : `GET`
 
 **EndPoint**: `/topics`
 
 **Request Filteting**: `app_name` `org` `user`
 
 **Response Body**:
 
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
**Response Code**: `200`

5. Get all channels

 **Method** : `GET`
 
 **EndPoint**: `/channels`
 
 **Request Filteting**: `app_name` `org` `user`
 
 **Response Body**:
 
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
**Response Code**: `200`

1. Get all notifications

  This will be polled periodically 

  **Method**: `GET`
  **Endpoint**: `/notifications`
  **Query params**: `user`, `org` (this information is already known to the api via the request, so not necessarily needed)
  **Response**:
  
  ```js
  [
    {
      org:"",
      app_name: "",
      user:"",
      destination_uri:"",
      text:"",
      topic:"",
      destination_uri:"",
      is_read:false,
      created_on: <milliseconds_since_epoch>
     },
    { _id: 2,... }
  ]
  ```

2. Mark a single notification as read

  **Method**: `PUT`
  **Endpoint**: `/notification/:_id`
  **Response status**: `200`

3. Mark all unread notifications as read

  **Method**: `PUT`
  **Endpoint**: `/notifications`
  **Response status**: `200`

4. View notifications history

  **Method**: `GET`
  **Endpoint**: `/notifications`
  **Query params**: `user`, `org`, `app`
  **Response**:
  
  ```js
  [
    {
      _id: 1,
      org:"",
      app_name: "",
      user:"",
      destination_uri:"",
      text:"",
      topic:"",
      destination_uri:"",
      is_read:false,
      created_on: <milliseconds_since_epoch>
    },
    { _id: 2,... }
  ]
  ```
  
`destination_uri` is  the link to relevant entity. (i.e action, incident)
`text` is the notification text
`topic` is the topic of notification (i.e incident_title_changed)


**Note**
For all of the above request you must pass atleast one of the `org`, `app_name` and `user`.
`ident` denotes the notification type (e.g incident_title_changed)

Every response has following structure :

```js
{
    "data": {..},
    "message": "",
    "status": 200
}
```


