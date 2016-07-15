FROM jeanblanchard/alpine-glibc
COPY de-job-killer /bin/de-job-killer
ENTRYPOINT ["de-job-killer"]
CMD ["--help"]
