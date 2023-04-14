# gspt

<div class="info" align="center">
    <br><img src="./extras/mascot-res.png" alt="mascot" width="200" class="mascot"/><br>
    Spotify for terminal written in Go.<br>
    with builtin <b>cover-art view</b> and <b>much more.</b> <br>
    <a href="https://aditya-K2.github.io/gspt/"> Documentation </a> |
    <a href="https://github.com/aditya-K2/gspt/discussions">Discussion</a>
</div>

---

![](./extras/screenshot.png)

---

***In a very experimental stage.***

## Setup

### How to Generate an API Key from Spotify Dashboard

If you want to use Spotify's API to create applications that interact with their music streaming service, you will need an API key. Here's how you can generate one from the Spotify Dashboard:

1. Go to the [Spotify Developer Dashboard](https://developer.spotify.com/dashboard/) and log in with your Spotify account credentials.

2. Click on the "Create an App" button to create a new application.

   ![Create an App](./extras/create.png)

3. Give your application a name and description, and agree to the terms of service. In the `Redirect URI` section add `http://localhost:8080/callback` as a callback URL. This is necessary for the OAuth 2.0 authentication flow to work. Click on "Create" to proceed.
   ![Create an App Form](./extras/create_form.png)

4. On the next page, you'll see the details of your newly created application. In the Settings Look for the section labeled "Client ID" and click on the "Show Client Secret" button.

5. You will now see your Client ID and Client Secret. You will need both of these to use the Spotify API in `gspt`

```bash
$ export SPOTIFY_ID= # client id
$ export SPOTIFY_SECRET= # client secret
```

7. After this you can just run `gspt`. And follow the link that it generates, and Login.

```bash
$ gspt
```

---
