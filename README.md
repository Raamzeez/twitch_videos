# twitch_videos
A GoLang server that returns a list of a twitch user's videos in a JSON format

How To Run:
1. Clone the Repository - `git clone https://github.com/Raamzeez/twitch_videos`
2. Edit the .env file provided with the correct information
   >NOTE: Please only edit the CLIENT_ID, CLIENT_SECRET, AND BEARER_HEADER fields. 
   >You will have to create a Twitch Developer account to retrieve a Client ID, Client Secret Key, and OAuth token.
   >Please visit [https://dev.twitch.tv/docs/api] for more information
3. Install all necessary packages - `go get`
4. Run the server - `go run main.go`
5. Access the API Endpoint
   >1. Go to your web browser and type `http://localhost:8080/user/{TWITCH_USERNAME_HERE}`
   >2. Use an application such as Postman and access the API Endpoint.
