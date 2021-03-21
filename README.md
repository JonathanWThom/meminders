# Meminders

A little reminders program. Meminders will send you text messages to remind you
of things.

### Development

There are a few primary dependencies to Meminders:
1. Go (though you can run it with Docker as well), and all the Go packages
   specified in [go.mod](https://github.com/JonathanWThom/meminders/blob/main/go.mod).
2. SQLite (though this could easily be swapped out).
3. Twilio (for sending text messages).

Check out the [test.env file](https://github.com/JonathanWThom/meminders/blob/main/.env.test) to see what environment variables you'll need. The program will also blow up if any are lacking, so that helps.

There is a [Makefile](https://github.com/JonathanWThom/meminders/blob/main/Makefile) which won't necessarily set everything up for you, but provides some convenience methods for development.

### API

Right now, Meminders has only one endpoint, `POST /meminders`. An example
request to a local server might look like this:

```
curl --location --request POST 'localhost:8080/reminders' \
--header 'Content-Type: application/json' \
--header 'Authorization: Basic YWRtaW46cGFzc3dvcmQ=' \
--data-raw '{"message": "This is a test. This is only a test.", "frequency": "once", "year": 2021, "month": "March", "hour": 15, "minute": 30, "day": 21, "zone": "America/Los_Angeles"}'
```

The program supports frequencies of "once", "daily", "weekly", and "monthly".
The API does not (yet!) do a great job of validating inputs, so make sure if you
set something for a given frequency, you send it a sensical set of other
parameters (e.g. if you want a weekly reminder, tell it what day of the week).

### Deployment

There's a nice little [GitHub Actions workflow](https://github.com/JonathanWThom/meminders/blob/main/.github/workflows/go.yml) set up that you could use to put this thing on a server, if you were to fork the repo. You'd probably want to sub out any reference to `jonathanwthom` for your own username.

### License

MIT
