## Commonly used clis

### Migration
To migrate, run:
```shell
docker-compose run --rm migrate 1 // 1 migrate up
docker-compose run --rm migrate -2 // 2 migrate down
```

### Testing
#### Auth
```shell
curl -i -H \
  "X-App-Token: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwiaWF0IjoxNTE2MjM5MDIyLCJhY2NvdW50X2lkIjoiVXNlck5vMSIsImNsaWVudF9rZXkiOiJGb2ROQ0ExeWI4OE5oQ09vckRITWNVdjFJRkp0QUxqcjhJdDMifQ.HtYQ-mV9eUcql95m8mP3yV6oiKBlfgN_4OgjjGgzrj2XSDaHE33T633i3lhIWD54OkZIgPzQ4Fl74Wfjkp6pwyAf-l7M22W2nk2vqIuETdefDmglbKvkBI4PpDiw60JpVyzev25TvA0P0qA94fVDmtXuS5HWIhCiLoy7FTb2_rCd9-txLPD2fvfCmQr7irVqzBc4R_4HXz_Oy04p-m8Co7DPkWpRSJZVYDtnO9_5luzuFwMh-bqEAZVMANRdU0VYgAnpKfFB5U3NmlAy-jXlXRm8fQNnIOuWR5igbEbwtxGPjH4PiPWRSFHxTKWwo8FuPqiNfVCExd3jNPQmX7pARw" \
  "http://localhost:3333/v1/pub-sub?message=Lorem%20ipsum%20dolor%20sit%20amet"; 
  echo "";

HTTP/1.1 200 OK
Server: nginx/1.27.3
Date: Wed, 30 Apr 2025 09:49:20 GMT
Content-Type: application/json; charset=utf-8
Content-Length: 64
Connection: keep-alive

{"message":"Message 'Lorem ipsum dolor sit amet' has been sent"}

```