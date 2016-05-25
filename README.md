[![wercker status](https://app.wercker.com/status/2c3f0a92ad1c221b4666bb8abb6d238c/m "wercker status")](https://app.wercker.com/project/bykey/2c3f0a92ad1c221b4666bb8abb6d238c)

# Websockets
Web sockets sample using third party messaging provider PubNub.

To run this example:

* `glide install` to create/update your vendor directory
* Create a file called `env` in the `local_config` directory
* Run using Wercker:  `./runlocal local_config/env`

This file creates environment variables with the Wercker `X_` prefix so that they will be injected into the Docker container when run. Your file should look something like this:

```
# Auth0 Credentials
X_AUTHZERO_ID=(your auth0 id)
X_AUTHZERO_SECRET=(your auth0 secret)
X_AUTHZERO_DOMAIN=(your auth0 domain)
X_AUTHZERO_CALLBACK=http://192.168.99.100/callback

X_PUBNUB_KEY_PUBLISH=(your pubnub publish key)
X_PUBNUB_KEY_SUBSCRIBE=(your pubnub subscribe key)
```

Make sure to change the IP address to match that of the docker machine that you're using for your Wercker builds. You can determine this with the command:

```
docker-machine ip default
```

You will need *both* a pubnub account and an auth0 account to see this application work properly.
