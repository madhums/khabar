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


