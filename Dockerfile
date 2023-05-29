FROM golang:1.20-alpine

WORKDIR /usr/src/app

#TODO: change to import latest built app from artifactory
ADD ./AuthServer /usr/src/app

# RUN go build .

# CMD [ "./AuthServer" ] 



