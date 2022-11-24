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

    gcloud functions deploy progress --runtime go119 --entry-point Progress --trigger-http --memory 128MB --allow-unauthenticated

## Test it locally

Build the project so it downloads the dependencies

    go build

Change the path of the template interchanging these two lines on the [progress file](progress.go#L57).

From this:

    //tpl, err := template.ParseFiles("progress.html")
	tpl, err := template.ParseFiles("src/progress/progress.html")

Change it to this:

    tpl, err := template.ParseFiles("progress.html")
	//tpl, err := template.ParseFiles("src/progress/progress.html")

> **Warning**
> Don't commit these changes!

Run it

    go run cmd/main.go

You can visit the endpoint in your favorite browser, for example:

    http://localhost:8080/progress/76
