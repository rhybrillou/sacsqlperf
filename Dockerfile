FROM golang:1.22

WORKDIR /performance

USER 12345

COPY ./perftest /bin/perftest

ENTRYPOINT [ "/bin/perftest" ]
CMD [ "/bin/perftest" ]
