# AutoCC

**AutoCC** *(Automatic Closed Captions)* works by utilizing Youtube Data API and DeepL API

## About

Google Translate which is the translation service used by YouTube is known to be inaccurate. This is why I decided to create this project.
DeepL's translations are much more natural sounding which makes them perfect for reaching a wider audience.

Unfortunately, YouTube's API quota is very limited and any increases have to manually requested and approved by Google. This is why AutoCC is not publicly hosted.
You can of course run it yourself and provide your own API keys. 

**AutoCC supports API credentials rotation with the purpose of overcoming Google's strict limits.**

## Tech stack
- Go
- TypeScript (SvelteKit)
- Redis
- Docker

## Action flow (simplified)

### Closed Captions translation

```mermaid
sequenceDiagram
    actor User
    participant AutoCC
    participant Youtube Data API
    participant DeepL API
    
    AutoCC ->> Youtube Data API: Get videos
    activate AutoCC
    activate Youtube Data API
    Youtube Data API ->> AutoCC: Return user's videos
    deactivate Youtube Data API

    User ->> AutoCC: Request translation of a video
    activate User
    AutoCC ->> Youtube Data API: Get a video's captions list
    activate Youtube Data API
    Youtube Data API ->> AutoCC: Return a video's captions list
    AutoCC ->> Youtube Data API: Get CC in source language
    Youtube Data API ->> AutoCC: Return CC in source language
    deactivate Youtube Data API

    loop Repeat for all target languages
        AutoCC ->> DeepL API: Translate captions
        activate DeepL API
        DeepL API ->> AutoCC: Return translated captions
        deactivate DeepL API
        AutoCC ->> Youtube Data API: Upload captions
    end

    AutoCC ->> User: Show success message
    deactivate AutoCC
    deactivate User
    
```

### Metadata translation

```mermaid
sequenceDiagram
    actor User
    participant AutoCC
    participant Youtube Data API
    participant DeepL API
    
    User ->> AutoCC: Expand metadata editor
    activate AutoCC
    activate User
    AutoCC ->> Youtube Data API: Get a video's metadata
    activate Youtube Data API
    Youtube Data API ->> AutoCC: Return a video's metadata
    
    AutoCC ->> User: Show original metadata
    User ->> AutoCC: Add (optional) separators

    deactivate Youtube Data API

    loop Repeat for all target languages
        AutoCC ->> DeepL API: Translate metadata
        activate DeepL API
        DeepL API ->> AutoCC: Return translated metadata
        deactivate DeepL API
    end
    AutoCC ->> Youtube Data API: Upload new metadata

    AutoCC ->> User: Show success message
    deactivate AutoCC
    deactivate User
    
```

## How to run

### Prerequisites

- Docker
- Google and DeepL API credentials

### Steps

0. Enable YouTube Data API v3 if you haven't done so already
1. Clone this repository
2. Create a `.env` file with:
- `REDIS_URL` - Your Redis URL (e.g. `redis://redis:6379` if you're using this project's docker-compose file)
- `GOOGLE_REDIRECT_URI` - Your API address followed by the callback endpoint path (e.g. `http://localhost:3001/youtube/callback`)
- `DEEPL_API_KEY` - Your DeepL API key
- `PORT` - Port to serve the API on
- `API_URL` - The app connects to the backend using this address

3. Run `docker compose up -d`
4. Open `localhost:3000` in your browser

## License

This project is licensed under the GNU v3 License - see the [LICENSE](LICENSE) file for details

## Screenshots

### Main page
![Main page](screenshots/1.png)

### Metadata translation
![Main page](screenshots/2.png)
