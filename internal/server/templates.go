package server

const malauthtpl string = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>shinkro</title>
    <link
      rel="icon"
      href="https://raw.githubusercontent.com/varoOP/shinkro/main/.github/images/logo.png"
      type="image/png"
    />
    <style>
      body {
        font-family: "Arial", sans-serif;
        background-color: #2d2d2d;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100vh;
        margin: 0;
        color: #ebaf01;
        text-align: center;
      }
      h1 {
        margin-bottom: 20px;
      }
      form {
        background-color: #2e51a1;
        padding: 30px;
        border-radius: 15px;
        box-shadow: 0px 4px 15px rgba(0, 0, 0, 0.2);
        color: #ebaf01;
        max-width: 600px;
        width: 80%;
      }
      label {
        margin-bottom: 5px;
        display: block;
      }
      input {
        padding: 10px;
        border: 1px solid #ebaf01;
        border-radius: 5px;
        width: 50%;
        margin-bottom: 15px;
        box-sizing: border-box;
        font-size: 16px;
      }
      input[type="submit"] {
        background-color: #ebaf01;
        border: none;
        color: #2e51a1;
        padding: 10px 20px;
        font-size: 16px;
        cursor: pointer;
        border-radius: 5px;
        transition: background-color 0.3s ease;
        width: 50%;
      }
      input[type="submit"]:hover {
        background-color: #c98e00;
      }
      input:focus {
        border-color: #2e51a1;
        outline: none;
        box-shadow: 0 0 5px #2e51a1;
      }
      .logo-container {
        background-color: #2d2d2d;
        display: inline-block;
        border-radius: 5px;
      }
      img {
        max-width: 120px;
      }
    </style>
  </head>
  <body>
    <form id="authForm" action="{{.ActionURL}}" method="post">
      <div class="logo-container">
        <img
          src="https://raw.githubusercontent.com/varoOP/shinkro/main/.github/images/logo.png"
          alt="Logo"
        />
      </div>
      <h1>shinkro</h1>
      <h3>Authenticate with myanimelist.net</h3>
      <label for="clientID">Client ID:</label>
      <input type="password" id="clientID" name="clientID" required />
      <label for="clientSecret">Client Secret:</label>
      <input type="password" id="clientSecret" name="clientSecret" required />
      <input type="submit" value="Start OAuth" />
    </form>
  </body>
</html>
`

const malauth_statustpl string = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1.0" />
    <title>shinkro</title>
    <link
      rel="icon"
      href="https://raw.githubusercontent.com/varoOP/shinkro/main/.github/images/logo.png"
      type="image/png"
    />
    <style>
      body {
        font-family: "Arial", sans-serif;
        background-color: #2d2d2d;
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100vh;
        margin: 0;
        color: #ebaf01;
        text-align: center;
      }
      h1 {
        margin-bottom: 20px;
      }
      .logo-container {
        background-color: #2d2d2d;
        display: inline-block;
        border-radius: 5px;
        margin-bottom: 20px;
      }
      img {
        max-width: 120px;
      }
      button {
        background-color: #ebaf01;
        border: none;
        color: #2e51a1;
        padding: 10px 20px;
        font-size: 16px;
        cursor: pointer;
        border-radius: 5px;
        transition: background-color 0.3s ease;
      }
      button:hover {
        background-color: #c98e00;
      }
    </style>
  </head>
  <body>
    <div class="logo-container">
      <img
        src="https://raw.githubusercontent.com/varoOP/shinkro/main/.github/images/logo.png"
        alt="Logo"
      />
    </div>
    <h1>shinkro</h1>
	{{if .IsAuthenticated}}
    <h2 id="authStatus">Authentication Success</h2>
    <p id="authMessage">You can close this window now.</p>
	{{else}}
	<h2 id="authStatus">Authentication Error</h2>
	<a href="{{.RetryURL}}" id="retryButtonLink">
    	<button id="retryButton">Retry</button>
    </a>
	{{end}}
  </body>
</html>
`
