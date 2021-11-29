 #!/bin/bash

# Change to the directory with our code that we plan to work from
# NOTE: DO NOT INCLUDE THIS LINE IF YOU ARE USING GO MODULES!
# old: cd "$GOPATH/src/lenslocked.com"
cd "$GOPATH/src/gitlab.com/go-courses/lenslocked.com"

echo "==== Releasing lenslocked.com ===="
echo "  Deleting the local binary if it exists (so it isn't uploaded)..."
rm lenslocked.com
echo "  Done!"

echo "  Deleting existing code..."
ssh root@143.110.237.111 "rm -rf /root/go/src/gitlab.com/go-courses/lenslocked.com"
echo "  Code deleted successfully!"

echo "  Uploading code..."
# The \ at the end of the line tells bash that our
# command isn't done and wraps to the next line.
rsync -avr --exclude '.git/*' --exclude 'tmp/*' \
  --exclude 'images/*' ./ \
  root@143.110.237.111:/root/go/src/gitlab.com/go-courses/lenslocked.com/
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

echo "  Building the code on remote server..."
ssh root@143.110.237.111 'export GOPATH=/root/go; \
  cd /root/app; \
  /usr/local/go/bin/go build -o ./server \
    $GOPATH/src/gitlab.com/go-courses/lenslocked.com/*.go'
echo "  Code built successfully!"

echo "  Moving assets..."
ssh root@143.110.237.111 "cd /root/app; \
  cp -R /root/go/src/gitlab.com/go-courses/lenslocked.com/assets ."
echo "  Assets moved successfully!"

echo "  Moving views..."
ssh root@143.110.237.111 "cd /root/app; \
  cp -R /root/go/src/gitlab.com/go-courses/lenslocked.com/views ."
echo "  Views moved successfully!"

echo "  Moving Caddyfile..."
ssh root@gallery.faulkners.io "cp /root/go/src/gitlab.com/go-courses/lenslocked.com/Caddyfile /etc/caddy/Caddyfile"
echo "  Caddyfile moved successfully!"

echo "  Restarting the server..."
ssh root@143.110.237.111 "sudo service lenslocked.com restart"
echo "  Server restarted successfully!"

echo "  Restarting Caddy server..."
ssh root@143.110.237.111 "sudo service caddy restart"
echo "  Caddy restarted successfully!"

echo "==== Done releasing lenslocked.com ===="
