## How to start?

1. Create kickstart.json file. Here is an example:

```json
{
  "apiKeys": [
    {
      "key": "oF3x8auXsmHgyWAWogYKbMJtNrKW_c520n1Yyf4sq7e4tZkhzITWsHyr"
    },
    {
      "key": "0ad00c5f140c43f8bc2a155dd0edf3b7d61f5dd4011843d8a28475a0d812c033",
      "description": "Restricted API key - only for Retrieve User in POC-Auth",
      "permissions": {
        "endpoints": {
          "/api/user": ["GET", "POST", "PATCH", "DELETE"],
          "/api/login": ["GET", "POST", "PATCH", "DELETE"],
          "/api/register": ["GET", "POST", "PATCH", "DELETE"],
          "/api/forgot-password": ["GET", "POST", "PATCH", "DELETE"],
          "/api/reset-password": ["GET", "POST", "PATCH", "DELETE"],
          "/api/verify-email": ["GET", "POST", "PATCH", "DELETE"]
        }
      },
      "tenantId": "#{secondTenantId}"
    }
  ],
  "requests": [
    {
      "method": "POST",
      "url": "/api/user/registration",
      "body": {
        "user": {
          "email": "admin@example.com",
          "password": "password"
        },
        "registration": {
          "applicationId": "4a434d08-00fe-4e29-aebd-c78fba6dc708"
        },
        "sendSetPasswordEmail": false,
        "skipVerification": true
      }
    },
    {
      "method": "POST",
      "url": "/api/application",
      "body": {
        "application": {
          "name": "poc-auth",
          "roles": [
            {
              "name": "User"
            }
          ]
        }
      }
    },
    {
      "method": "PATCH",
      "url": "/api/tenant",
      "body": {
        "tenant": {
          "id": "4a434d08-00fe-4e29-aebd-c78fba6dc709",
          "emailConfiguration": {
            "host": "smtp.gmail.com",
            "port": 587,
            "username": "admin@example.com",
            "password": "password",
            "security": "TLS"
          }
        }
      }
    }
  ]
}
```

2. Create config.yaml file. Here is an example:

```yaml
fusion_auth:
  host: http://fusionauth:9011
  app_id: 4a434d08-00fe-4e29-aebd-c78fba6dc708
  api_key: oF3x8auXsmHgyWAWogYKbMJtNrKW_c520n1Yyf4sq7e4tZkhzITWsHyr

database:
  mongodb:
    uri: mongodb://root:example@poc-auth-mongodb:27017/auth?authSource=admin&ssl=false
```

3. Run the following command:

```bash
docker-compose up
```

4. yay! You have a running POC-Auth instance. Now you can use it to authenticate your users.
