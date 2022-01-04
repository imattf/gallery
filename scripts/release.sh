 #!/bin/bash

# Change to the directory with our code that we plan to work from
# NOTE: DO NOT INCLUDE THIS LINE IF YOU ARE USING GO MODULES!
# old: cd "$GOPATH/src/lenslocked.com"
cd "$GOPATH/src/github.com/imattf/go-courses/gallery"

echo "==== Releasing gallery ===="
echo "  Starting deployment from $PWD..."
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm gallery
echo "  Done!"

echo "  Deleting existing code..."
ssh root@143.110.237.111 "rm -rf /root/go/src/github.com/imattf/go-courses/gallery"
echo "  Code deleted successfully!"

echo "  Uploading code..."
# The \ at the end of the line tells bash that our
# command isn't done and wraps to the next line.
rsync -avr --exclude '.git/*' --exclude 'tmp/*' \
  --exclude 'images/*' ./ \
  root@143.110.237.111:/root/go/src/github.com/imattf/go-courses/gallery/
echo "  Code uploaded successfully!"

echo "  Go getting deps..."
ssh root@143.110.237.111 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get golang.org/x/crypto/bcrypt"
ssh root@143.110.237.111 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/mux"
ssh root@143.110.237.111 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/schema"
ssh root@143.110.237.111 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/lib/pq"
ssh root@143.110.237.111 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/jinzhu/gorm"
ssh root@143.110.237.111 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/gorilla/csrf"
ssh root@143.110.237.111 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get gopkg.in/mailgun/mailgun-go.v1"
ssh root@143.110.237.111 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox"
ssh root@143.110.237.111 "export GOPATH=/root/go; \
  /usr/local/go/bin/go get github.com/dropbox/dropbox-sdk-go-unofficial/v6/dropbox/files"

echo "  Building the code on remote server..."
ssh root@143.110.237.111 'export GOPATH=/root/go; \
  cd /root/app; \
  /usr/local/go/bin/go build -o ./server \
    $GOPATH/src/github.com/imattf/go-courses/gallery/*.go'
echo "  Code built successfully!"

echo "  Moving assets..."
ssh root@143.110.237.111 "cd /root/app; \
  cp -R /root/go/src/github.com/imattf/go-courses/gallery/assets ."
echo "  Assets moved successfully!"

echo "  Moving views..."
ssh root@143.110.237.111 "cd /root/app; \
  cp -R /root/go/src/github.com/imattf/go-courses/gallery/views ."
echo "  Views moved successfully!"

echo "  Moving Caddyfile..."
ssh root@gallery.faulkners.io "cp /root/go/src/github.com/imattf/go-courses/gallery/Caddyfile /etc/caddy/Caddyfile"
echo "  Caddyfile moved successfully!"

echo "  Restarting the server..."
ssh root@143.110.237.111 "sudo service gallery restart"
echo "  Server restarted successfully!"

echo "  Restarting Caddy server..."
ssh root@143.110.237.111 "sudo service caddy restart"
echo "  Caddy restarted successfully!"

echo "==== Done releasing gallery ===="
