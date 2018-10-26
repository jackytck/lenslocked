#!/bin/bash

# Change to the directory with our code that we plan to work from
cd "$GOPATH/src/github.com/jackytck/lenslocked"

echo "==== Releasing lenslocked.jackytck.com ===="
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm server lenslocked
echo "  Done!"

echo "  Syncing code..."
ssh jacky@lenslocked.jackytck.com "mkdir -p /home/jacky/go/src/github.com/jackytck"
rsync -avr --exclude '.git/*' --exclude 'tmp/*' --exclude 'images/*' --delete ./ jacky@lenslocked.jackytck.com:/home/jacky/go/src/github.com/jackytck/lenslocked/
echo "  Code uploaded successfully!"

echo "  Go getting deps..."
ssh jacky@lenslocked.jackytck.com "export GOPATH=/home/jacky/go; /usr/local/go/bin/go get golang.org/x/crypto/bcrypt "
ssh jacky@lenslocked.jackytck.com "export GOPATH=/home/jacky/go; /usr/local/go/bin/go get github.com/gorilla/mux"
ssh jacky@lenslocked.jackytck.com "export GOPATH=/home/jacky/go; /usr/local/go/bin/go get github.com/gorilla/schema"
ssh jacky@lenslocked.jackytck.com "export GOPATH=/home/jacky/go; /usr/local/go/bin/go get github.com/lib/pq"
ssh jacky@lenslocked.jackytck.com "export GOPATH=/home/jacky/go; /usr/local/go/bin/go get github.com/jinzhu/gorm"
ssh jacky@lenslocked.jackytck.com "export GOPATH=/home/jacky/go; /usr/local/go/bin/go get github.com/gorilla/csrf"
ssh jacky@lenslocked.jackytck.com "export GOPATH=/home/jacky/go; /usr/local/go/bin/go get gopkg.in/mailgun/mailgun-go.v1"

echo "  Building the code on remote server..."
ssh jacky@lenslocked.jackytck.com 'export GOPATH=/home/jacky/go; cd /home/jacky/app; /usr/local/go/bin/go build -o ./server $GOPATH/src/github.com/jackytck/lenslocked/*.go'
echo "  Code built successfully!"

echo "  Moving assets..."
ssh jacky@lenslocked.jackytck.com "cd /home/jacky/app; cp -R /home/jacky/go/src/github.com/jackytck/lenslocked/assets ."
echo "  Assets moved successfully!"

echo "  Moving views..."
ssh jacky@lenslocked.jackytck.com "cd /home/jacky/app; cp -R /home/jacky/go/src/github.com/jackytck/lenslocked/views ."
echo "  Views moved successfully!"

echo "  Moving Caddyfile..."
ssh jacky@lenslocked.jackytck.com "cd /home/jacky/app; cp /home/jacky/go/src/github.com/jackytck/lenslocked/deploy/Caddyfile ."
echo "  Views moved successfully!"

echo "  Restarting the server..."
ssh jacky@lenslocked.jackytck.com "sudo systemctl restart lenslocked"
echo "  Server restarted successfully!"

echo "  Restarting Caddy server..."
ssh jacky@lenslocked.jackytck.com "sudo systemctl restart caddy"
echo "  Caddy restarted successfully!"

echo "==== Done releasing lenslocked.jackytck.com ===="
