# Markdown Progress ![](https://geps.dev/progress/100)

Progress bars for markdown.

Have you ever wanted to track some progress in your markdown documents?
Well, I do, and I used `progressed.io` before but it was shutted down.

So I decided to recreate it.

## Usage

Add it as an image in your favorite markdown document, like this github readme, and change the progress number at the end.

    ![](https://geps.dev/progress/10)

## Examples

![](https://geps.dev/progress/10)

![](https://geps.dev/progress/50)

![](https://geps.dev/progress/75)

## Deploy

### Google Cloud

Login and set the project in `gcloud` if you are not already logged in.

    gcloud auth login
    gcloud config set project THE_PROJECT_NAME

Deploy it as an HTTP Cloud Function with the `Progress` entrypoint.

    gcloud functions deploy progress --runtime go116 --entry-point Progress --trigger-http --memory 128MB --allow-unauthenticated
