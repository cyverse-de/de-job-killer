FROM golang:1.6-alpine

ARG git_commit=unknown
LABEL org.cyverse.git-ref="$git_commit"

COPY . /go/src/github.com/cyverse-de/de-job-killer
RUN go install github.com/cyverse-de/de-job-killer

ENTRYPOINT ["de-job-killer"]
CMD ["--help"]
