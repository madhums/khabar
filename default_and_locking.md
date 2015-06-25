## 1. User level notification settings

**Url:** `/topics`

**Method:** `GET`

**Query parameters:** `app_name` (Compulsory) `user` (Optional) `org` (Optional)

**Example Response:**

```json
{
    "log_used": {
        "email": {
            "locked": false,
            "value": "true"
        },
        "push": {
            "locked": false,
            "value": "disabled"
        },
        "web": {
            "locked": false,
            "value": "disabled"
        }
    },
    "log_used2": {
        "email": {
            "locked": false,
            "value": "disabled"
        },
        "push": {
            "locked": false,
            "value": "disabled"
        },
        "web": {
            "locked": false,
            "value": "disabled"
        }
    }
}
```

**Explanation:** 
`value` denotes whether the channel for particular notification is enabled or not. Possible values : `"true"`, `"false"`, `"disabled"`

`locked` denotes whether the channel for particular notification is locked or not.  Possible values: `true`, `false`


## 2. Organization level notification settings

**Url:** `topics`

**Method:** `GET`

**Query parameters:** `app_name` (Compulsory) `org` (Optional)

**Example Response:**

```json
{
    "log_used": {
        "email": {
            "locked": false,
            "value": "true"
        },
        "push": {
            "locked": false,
            "value": "disabled"
        },
        "web": {
            "locked": false,
            "value": "disabled"
        }
    },
    "log_used2": {
        "email": {
            "locked": false,
            "value": "disabled"
        },
        "push": {
            "locked": false,
            "value": "disabled"
        },
        "web": {
            "locked": false,
            "value": "disabled"
        }
    }
}
```

**Explanation:** 
`value` denotes whether the channel for particular notification is enabled or not. Possible values : `"true"`, `"false"`, `"disabled"`

`locked` denotes whether the channel for particular notification is locked or not.  Possible values: `true`, `false`


## 3. Create defaults:

**Url:** `/topics/defaults/<topic>/channels/<channel>`

**Method:** `POST`

**Example Body:**

```json
{
      "org": <org_id>,
      "enabled": true
}
```

**Response Code:** `201`

## 4. Delete defaults:

**Url:** `/topics/defaults/<topic>/channels/<channel>`

**Method:** `DELETE`

**Example Body:**

```json
{
      "org": <org_id>,
      "enabled": true
}
```

**Response Code:** `204`

## 5. Create Lock:

**Url:** `/topics/locked/<topic>/channels/<channel>`

**Method:** `POST`

**Example Body:**

```json
{
      "org": <org_id>,
      "enabled": true
}
```

**Response Code:** `201`

## 6. Delete lock:

**Url:** `/topics/locked/<topic>/channels/<channel>`

**Method:** `POST`

**Example Body:**

```json
{
      "org": <org_id>,
      "enabled": true
}
```

**Response code:** `201`
