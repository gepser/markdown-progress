# Markdown Progress ![](https://us-central1-una-fiesta-bonita-1.cloudfunctions.net/progress/100)

Progress bars for markdown.

Have you ever wanted to track some progress in your markdown documents?
Well, I do, and I used `progressed.io` before but it was shutted down.

So I decided to recreated.

## Usage

Add it as an image in your favorite markdown document, like this github readme, and change the progress number at the end.

    ![](https://us-central1-una-fiesta-bonita-1.cloudfunctions.net/progress/10)

## Examples

![](https://us-central1-una-fiesta-bonita-1.cloudfunctions.net/progress/10)

![](https://us-central1-una-fiesta-bonita-1.cloudfunctions.net/progress/50)

![](https://us-central1-una-fiesta-bonita-1.cloudfunctions.net/progress/75)

## Deploy

### Google CLoud

Login and set the project in `gcloud` if you are nott already logged in.

    gcloud auth login
    gcloud config set project THE_PROJECT_NAME

Deploy it as an HTTP Cloud Function with the `Progress` entrypoint.

    gcloud functions deploy progress --runtime go111 --entry-point Progress --trigger-http --memory 128MB
